package domain

import (
	"context"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
)

type Comment struct {
	ID        string `json:"id"`
	ThreadID  string `json:"thread_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	IsActive  *bool  `json:"is_active"`
	CreatedBy string `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

// // CommentRepository represent the Comment repository contract
type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, comment *Comment) (int64, error)
	GetList(ctx context.Context, request *request.GetListCommentReq) ([]response.GetListCommentRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailCommentReq) (Comment, error)
}

// CommentUsecase represent the Comment usecase contract
type CommentUsecase interface {
	Create(ctx context.Context, request *request.CreateCommentReq) (response.CreateCommentRes, error)
	Update(ctx context.Context, request *request.UpdateCommentReq) error
	Delete(ctx context.Context, request *request.DeleteCommentReq) (int64, error)
	GetList(ctx context.Context, request *request.GetListCommentReq) ([]response.GetListCommentRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailCommentReq) (Comment, error)
}
