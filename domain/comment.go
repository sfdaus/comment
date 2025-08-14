package domain

import (
	"context"
	"prakarsa-app/entity"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
)

// // CommentRepository represent the Comment repository contract
type CommentRepository interface {
	Create(ctx context.Context, comment *entity.Comment) error
	Update(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, comment *entity.Comment) (int64, error)
	GetList(ctx context.Context, request *request.GetListCommentReq) ([]response.GetListCommentRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailCommentReq) (response.GetDetailCommentRes, error)
}

// CommentUsecase represent the Comment usecase contract
type CommentUsecase interface {
	Create(ctx context.Context, request *request.CreateCommentReq) (response.CreateCommentRes, error)
	Update(ctx context.Context, request *request.UpdateCommentReq) error
	Delete(ctx context.Context, request *request.DeleteCommentReq) (int64, error)
	GetList(ctx context.Context, request *request.GetListCommentReq) ([]response.GetListCommentRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailCommentReq) (response.GetDetailCommentRes, error)
}
