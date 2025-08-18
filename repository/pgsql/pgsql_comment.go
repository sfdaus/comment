package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/entity"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
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

func (r *pgsqlCommentRepository) Create(ctx context.Context, comment *entity.Comment) (err error) {
	query := `INSERT INTO comments (id, thread_id, user_id, content, is_active, created_by, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err = r.db.ExecContext(ctx, query, comment.ID, comment.ThreadID, comment.UserID, comment.Content, comment.IsActive,
		comment.CreatedBy, comment.CreatedAt, comment.UpdatedAt); err != nil {
		return err
	}

	return
}

func (r *pgsqlCommentRepository) Update(ctx context.Context, comment *entity.Comment) (err error) {
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

func (r *pgsqlCommentRepository) Delete(ctx context.Context, comment *entity.Comment) (rowsAffected int64, err error) {
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

	if request.ThreadID != "" {
		wheres = append(wheres, fmt.Sprintf("thread_id = $%d", idx))
		args = append(args, request.ThreadID)
		idx++
	}

	if request.IsActive != nil {
		wheres = append(wheres, fmt.Sprintf("c.is_active = $%d", idx))
		args = append(args, request.IsActive)
		idx++
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}

	// 2) Hitung total data (tanpa LIMIT/OFFSET)
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM comments c
		JOIN profiles p ON p.user_id = c.user_id
		LEFT JOIN institutions i ON i.id = p.institution_id
		%s
	`, whereSQL)

	if err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&meta.TotalData); err != nil {
		return nil, meta, err
	}

	// 3) Pagination setup
	perPage := request.PerPage
	if perPage <= 0 {
		perPage = 10
	}
	page := request.Page
	if page <= 0 {
		page = 1
	}

	meta.Page = page
	meta.PerPage = perPage
	meta.TotalPages = (meta.TotalData + perPage - 1) / perPage

	offset := (page - 1) * perPage

	// add LIMIT & OFFSET to args
	args = append(args, perPage, offset)
	limitPos, offsetPos := idx, idx+1

	// ---------- 4) Query data ----------
	// NOTE: kolom institutions.alias diasumsikan ada (sesuai struct lo).
	// Jika di DB belum ada, hapus kolom & mapping alias di SELECT/Scan.
	query := fmt.Sprintf(`
		SELECT
			-- comment
			c.id, c.thread_id, c.user_id, c.content,
			c.is_active, c.created_at, c.updated_at,
			-- profile
			p.name, p.name_alias, p.avatar,
			-- institution
			i.name, i.alias, i.type
		FROM comments c
		JOIN profiles p ON p.user_id = c.user_id
		LEFT JOIN institutions i ON i.id = p.institution_id
		%s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, limitPos, offsetPos)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()

	// ---------- 5) Scan ----------
	for rows.Next() {
		var (
			// comments
			cID, cThreadID, cUserID string
			cContent                string
			cIsActive               sql.NullBool
			cCreatedAt, cUpdatedAt  int64

			// profile
			pName           string
			pAlias, pAvatar sql.NullString

			// institution
			iName, iAlias, iType sql.NullString
		)

		if err := rows.Scan(
			&cID, &cThreadID, &cUserID, &cContent,
			&cIsActive, &cCreatedAt, &cUpdatedAt,
			&pName, &pAlias, &pAvatar,
			&iName, &iAlias, &iType,
		); err != nil {
			return nil, meta, err
		}

		var isActivePtr *bool
		if cIsActive.Valid {
			v := cIsActive.Bool
			isActivePtr = &v
		}

		item := response.GetListCommentRes{
			ID:        cID,
			ThreadID:  cThreadID,
			Content:   cContent,
			IsActive:  *isActivePtr,
			CreatedAt: cCreatedAt,
			UpdatedAt: cUpdatedAt,
			Profile: entity.Profile{
				Name:      pName,
				NameAlias: pAlias.String,
				Avatar:    pAvatar.String,
				Institution: entity.Institution{
					Name:  iName.String,
					Alias: iAlias.String,
					Type:  iType.String,
				},
			},
		}

		res = append(res, item)
	}
	if err := rows.Err(); err != nil {
		return nil, meta, err
	}

	return
}

func (r *pgsqlCommentRepository) GetDetail(ctx context.Context, request *request.GetDetailCommentReq) (res response.GetDetailCommentRes, err error) {

	const query = `
					SELECT
						-- comment
						c.id, c.thread_id, c.user_id, c.content,
						c.is_active, c.created_at, c.updated_at,
						-- profile
						p.name, p.name_alias, p.avatar,
						-- institution
						i.name, i.alias, i.type
					FROM comments c
					JOIN profiles p ON p.user_id = c.user_id
					LEFT JOIN institutions i ON i.id = p.institution_id
					WHERE c.id = $1
					LIMIT 1
					`

	// 1. QueryRowContext untuk ambil satu baris
	row := r.db.QueryRowContext(ctx, query, request.ID)

	// ---------- 5) Scan ----------
	var (
		// comments
		cID, cThreadID, cUserID string
		cContent                string
		cIsActive               sql.NullBool
		cCreatedAt, cUpdatedAt  int64

		// profile
		pName           string
		pAlias, pAvatar sql.NullString

		// institution
		iName, iAlias, iType sql.NullString
	)

	if err := row.Scan(
		&cID, &cThreadID, &cUserID, &cContent,
		&cIsActive, &cCreatedAt, &cUpdatedAt,
		&pName, &pAlias, &pAvatar,
		&iName, &iAlias, &iType,
	); err != nil {
		return res, err
	}

	var isActivePtr *bool
	if cIsActive.Valid {
		v := cIsActive.Bool
		isActivePtr = &v
	}

	res = response.GetDetailCommentRes{
		ID:        cID,
		ThreadID:  cThreadID,
		Content:   cContent,
		IsActive:  *isActivePtr,
		CreatedAt: cCreatedAt,
		UpdatedAt: cUpdatedAt,
		Profile: entity.Profile{
			Name:      pName,
			NameAlias: pAlias.String,
			Avatar:    pAvatar.String,
			Institution: entity.Institution{
				Name:  iName.String,
				Alias: iAlias.String,
				Type:  iType.String,
			},
		},
	}

	if err := row.Err(); err != nil {
		return res, err
	}

	return
}
