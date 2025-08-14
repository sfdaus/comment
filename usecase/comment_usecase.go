package usecase

import (
	"context"
	"prakarsa-app/entity"
	"prakarsa-app/transport/response"
	"time"

	"prakarsa-app/domain"
	"prakarsa-app/repository/redis"
	"prakarsa-app/transport/request"

	"github.com/google/uuid"
)

type CommentUsecase struct {
	CommentRepo domain.CommentRepository
	RedisRepo   redis.RedisRepository
	CtxTimeout  time.Duration
}

// NewCommentUsecase will create new an CommentUsecase object representation of ThreadUsecase interface
func NewCommentUsecase(commentRepo domain.CommentRepository, redisRepo redis.RedisRepository, ctxTimeout time.Duration) *CommentUsecase {
	return &CommentUsecase{
		CommentRepo: commentRepo,
		RedisRepo:   redisRepo,
		CtxTimeout:  ctxTimeout,
	}
}

func (u *CommentUsecase) Create(c context.Context, request *request.CreateCommentReq) (res response.CreateCommentRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.CtxTimeout)
	defer cancel()

	// Create Payload
	CommentID := uuid.NewString()
	t := true
	CommentPayload := &entity.Comment{
		ID:        CommentID,
		ThreadID:  request.ThreadID,
		UserID:    "TODO_user_id",
		Content:   request.Content,
		IsActive:  &t,
		CreatedBy: "TODO_created_by",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Response Payload
	res.ID = CommentID

	err = u.CommentRepo.Create(ctx, CommentPayload)
	return
}

func (u *CommentUsecase) Update(c context.Context, request *request.UpdateCommentReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.CtxTimeout)
	defer cancel()

	// Update Payload
	commentPayload := &entity.Comment{
		ID:        request.ID,
		UpdatedBy: "TODO_updated_by",
		UpdatedAt: time.Now().Unix(),
	}

	if request.Content != "" {
		commentPayload.Content = request.Content
	}

	err = u.CommentRepo.Update(ctx, commentPayload)
	return
}
func (u *CommentUsecase) Delete(c context.Context, request *request.DeleteCommentReq) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(c, u.CtxTimeout)
	defer cancel()

	commentPayload := &entity.Comment{
		ID: request.ID,
	}

	rowsAffected, err = u.CommentRepo.Delete(ctx, commentPayload)
	return
}

func (u *CommentUsecase) GetList(c context.Context, request *request.GetListCommentReq) (res []response.GetListCommentRes, meta response.MetaRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.CtxTimeout)
	defer cancel()

	t := true
	request.IsActive = &t

	res, meta, err = u.CommentRepo.GetList(ctx, request)
	return
}

func (u *CommentUsecase) GetDetail(c context.Context, request *request.GetDetailCommentReq) (res entity.Comment, err error) {
	ctx, cancel := context.WithTimeout(c, u.CtxTimeout)
	defer cancel()

	res, err = u.CommentRepo.GetDetail(ctx, request)
	return
}
