package response

import "prakarsa-app/entity"

type CreateCommentRes struct {
	ID string `json:"id"`
}

// Get List Response
type GetListCommentRes struct {
	ID          string             `json:"id"`
	ThreadID    string             `json:"thread_id"`
	Content     string             `json:"content"`
	Profile     entity.Profile     `json:"profile"`
	Institution entity.Institution `json:"institution"`
	IsActive    bool               `json:"is_active"`
	CreatedAt   int64              `json:"created_at"`
	UpdatedAt   int64              `json:"updated_at"`
}

// Get Detail Response
type GetDetailCommentRes struct {
	ID          string             `json:"id"`
	ThreadID    string             `json:"thread_id"`
	Content     string             `json:"content"`
	Profile     entity.Profile     `json:"profile"`
	Institution entity.Institution `json:"institution"`
	IsActive    bool               `json:"is_active"`
	CreatedAt   int64              `json:"created_at"`
	UpdatedAt   int64              `json:"updated_at"`
}
