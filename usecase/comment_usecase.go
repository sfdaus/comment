package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"prakarsa-app/config"
	"prakarsa-app/entity"
	"prakarsa-app/repository/s3"
	"prakarsa-app/transport/response"
	"prakarsa-app/utils"
	"time"

	"prakarsa-app/domain"
	"prakarsa-app/repository/redis"
	"prakarsa-app/transport/request"

	"github.com/google/uuid"
)

type CommentUsecase struct {
	CommentRepo domain.CommentRepository
	RedisRepo   redis.RedisRepository
	s3Repo      s3.S3Repository
	CtxTimeout  time.Duration
}

// NewCommentUsecase will create new an CommentUsecase object representation of ThreadUsecase interface
func NewCommentUsecase(commentRepo domain.CommentRepository, redisRepo redis.RedisRepository, s3Repo s3.S3Repository, ctxTimeout time.Duration) *CommentUsecase {
	return &CommentUsecase{
		CommentRepo: commentRepo,
		RedisRepo:   redisRepo,
		s3Repo:      s3Repo,
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
		UserID:    request.UserID,
		Content:   request.Content,
		IsActive:  &t,
		CreatedBy: request.UserID,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Notification Outbox Payload
	headers := map[string]string{"x-user-id": request.UserID}
	headersJSON, _ := json.Marshal(headers)

	initiatorNotificationOutbox := &entity.NotificationOutboxInsert{
		ID:            uuid.NewString(),
		Type:          utils.CREATE_COMMENT_NOTIFICATION_TYPE,
		ReferenceType: utils.COMMENT_NOTIFICATION_REFERENCE_TYPE,
		ReferenceID:   request.ThreadID,
		HeadersJSON:   headersJSON,
		Title:         utils.CommentNotificationTitle["COMMENT_NOTIFICATION_TITLE"],
		Message:       request.Content,
		Priority:      utils.CommentNotificationPriority[utils.CREATE_COMMENT_NOTIFICATION_TYPE],
		IdempotencyKey: fmt.Sprintf(
			"%s:%s:%s", utils.NotificationIdempotencyKey[utils.CREATE_COMMENT_NOTIFICATION_TYPE],
			CommentID, "[INIT_ID]",
		),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Response Payload
	res.ID = CommentID

	err = u.CommentRepo.Create(ctx, CommentPayload, initiatorNotificationOutbox)
	return
}

func (u *CommentUsecase) Update(c context.Context, request *request.UpdateCommentReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.CtxTimeout)
	defer cancel()

	// Update Payload
	commentPayload := &entity.Comment{
		ID:        request.ID,
		UpdatedBy: request.UserID,
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
		ID:     request.ID,
		UserID: request.UserID,
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

	// Handle response
	if len(res) > 0 {
		for i, comment := range res {
			res[i].Profile.Avatar, err = u.s3Repo.GetPresignedURL(c, config.LoadConfig().S3Bucket, comment.Profile.Avatar, false, time.Duration(24*time.Hour))
		}
	}

	return
}

func (u *CommentUsecase) GetDetail(c context.Context, request *request.GetDetailCommentReq) (res response.GetDetailCommentRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.CtxTimeout)
	defer cancel()

	res, err = u.CommentRepo.GetDetail(ctx, request)

	// Handler response
	res.Profile.Avatar, err = u.s3Repo.GetPresignedURL(c, config.LoadConfig().S3Bucket, res.Profile.Avatar, false, time.Duration(24*time.Hour))

	return
}

func (u *CommentUsecase) CommentReport(c context.Context, request *request.CommentReportReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.CtxTimeout)
	defer cancel()

	contentReport := &entity.ContentReport{
		ID:         uuid.NewString(),
		ReporterID: request.UserID,
		ThreadID:   "",
		CommentID:  request.ID,
		ReasonID:   request.ReasonID,
		Status:     "OPEN",
		IsActive:   true,
		CreatedAt:  time.Now().Unix(),
		CreatedBy:  request.UserID,
		UpdatedAt:  time.Now().Unix(),
	}

	err = u.CommentRepo.CommentReport(ctx, contentReport)

	return
}
