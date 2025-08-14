package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
	"prakarsa-app/utils"
	"strings"
)

type pgsqlCommentRepository struct {
	db *sql.DB
}

// NewPgsqlCommentRepository will create new an todoRepository object representation of CommentRepository interface
func NewPgsqlCommentRepository(db *sql.DB) *pgsqlCommentRepository {
	return &pgsqlCommentRepository{
		db: db,
	}
}

func (r *pgsqlCommentRepository) Create(ctx context.Context, comment *domain.Comment) (err error) {
	query := `INSERT INTO comments (id, thread_id, user_id, content, is_active, created_by, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err = r.db.ExecContext(ctx, query, comment.ID, comment.ThreadID, comment.UserID, comment.Content, comment.IsActive,
		comment.CreatedBy, comment.CreatedAt, comment.UpdatedAt); err != nil {
		return err
	}

	return
}

func (r *pgsqlCommentRepository) Update(ctx context.Context, comment *domain.Comment) (err error) {
	// Build dynamic SET clauses from Comment struct
	sets := []string{}
	args := []interface{}{}
	idx := 1

	if comment.Content != "" {
		sets = append(sets, fmt.Sprintf("content = $%d", idx))
		args = append(args, comment.Content)
		idx++
	}

	// kalau ada sesuatu untuk di‐update, commit ke SQL
	if len(sets) > 0 {
		// Update stamp
		sets = append(sets, fmt.Sprintf("updated_at = $%d", idx))
		args = append(args, comment.UpdatedAt)
		idx++

		sets = append(sets, fmt.Sprintf("updated_by = $%d", idx))
		args = append(args, comment.UpdatedBy)
		idx++

		// tambahkan WHERE id = $idx
		args = append(args, comment.ID)
		query := fmt.Sprintf(
			"UPDATE comments SET %s WHERE id = $%d",
			strings.Join(sets, ", "),
			idx,
		)

		if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
			return
		}
	}

	return
}

func (r *pgsqlCommentRepository) Delete(ctx context.Context, comment *domain.Comment) (rowsAffected int64, err error) {
	query := "DELETE FROM comments WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, comment.ID)
	if err != nil {
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	return
}

func (r *pgsqlCommentRepository) GetList(ctx context.Context, request *request.GetListCommentReq) (res []response.GetListCommentRes, meta response.MetaRes, err error) {
	// 1. Build WHERE clauses
	wheres := []string{}
	args := []interface{}{}
	idx := 1

	if request.Name != "" {
		wheres = append(wheres, fmt.Sprintf("name ILIKE $%d", idx))
		args = append(args, "%"+request.Name+"%")
		idx++
	}

	if request.IsActive != nil {
		wheres = append(wheres, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, request.IsActive)
		idx++
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}

	// --- 2. Hitung totalCount dulu (tanpa LIMIT/OFFSET) ---
	countQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM comments %s",
		whereSQL,
	)
	if err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&meta.TotalData); err != nil {
		return nil, meta, err
	}

	// 2. Calculate LIMIT & OFFSET
	perPage := request.PerPage
	if perPage <= 0 {
		perPage = 10
	}
	page := request.Page
	if page <= 0 {
		page = 1
	}

	// total pages = ceil(total / perPage)
	meta.Page = page
	meta.PerPage = perPage
	meta.TotalPages = (meta.TotalData + perPage - 1) / perPage

	offset := (page - 1) * perPage

	// add LIMIT & OFFSET to args
	args = append(args, perPage, offset)
	limitPos, offsetPos := idx, idx+1

	// 3. Final query
	query := fmt.Sprintf(`
        SELECT
            id, name, description,
            is_active
        FROM Comments
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereSQL, limitPos, offsetPos)

	// 4. Execute
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()

	// 5. Scan results
	for rows.Next() {
		var item response.GetListCommentRes

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.IsActive,
		); err != nil {
			return nil, meta, err
		}

		res = append(res, item)
	}
	if errRow := rows.Err(); errRow != nil {
		return nil, meta, errRow
	}

	return
}

func (r *pgsqlCommentRepository) GetDetail(ctx context.Context, request *request.GetDetailCommentReq) (res domain.Comment, err error) {

	const query = `
					SELECT
					  id,
					  name,
					  description,
					  is_active,
					  created_at,
					  created_by,
					  updated_at,
					  updated_by,
					  deleted_at
					FROM comments
					WHERE id = $1
					LIMIT 1
					`

	// 1. QueryRowContext untuk ambil satu baris
	row := r.db.QueryRowContext(ctx, query, request.ID)

	// 2. Scan kolom ke field di domain.Comment
	// since created_at is NOT NULL int8:
	var createdAt int64
	// updated_at/deleted_at can be NULL, so use NullInt64:
	var updatedAt, deletedAt sql.NullInt64
	var updatedBy sql.NullString

	err = row.Scan(
		&res.ID,
		// &res.Name,
		// &res.Description,
		&res.IsActive,
		&createdAt,
		&res.CreatedBy,
		&updatedAt,
		&updatedBy,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, utils.NewNotFoundError("comment not found")
		}
		return res, err
	}

	// assign into your domain fields
	res.CreatedAt = createdAt
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.Int64
	}
	if deletedAt.Valid {
		res.DeletedAt = deletedAt.Int64
	}
	if updatedBy.Valid {
		res.UpdatedBy = updatedBy.String
	}

	return
}
