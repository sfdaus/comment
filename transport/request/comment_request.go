package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// CreateCommentReq represent create request body
type CreateCommentReq struct {
	ThreadID string `json:"thread_id"`
	Content  string `json:"content"`
	UserID   string
}

func (request CreateCommentReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ThreadID, validation.Required),
		validation.Field(&request.Content, validation.Required),
		validation.Field(&request.UserID, validation.Required),
	)
}

// Update request body
type UpdateCommentReq struct {
	ID      string `param:"id"`
	Content string `json:"content"`
	UserID  string
}

func (request UpdateCommentReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.UserID, validation.Required),
	)
}

// Delete request body
type DeleteCommentReq struct {
	ID     string `param:"id"`
	UserID string
}

func (request DeleteCommentReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.UserID, validation.Required),
	)
}

// GetList request body
type GetListCommentReq struct {
	ThreadID string `query:"thread_id"`
	IsActive *bool  `query:"is_active"`
	PerPage  int64  `query:"per_page"`
	Page     int64  `query:"page"`
	UserID   string
}

func (request GetListCommentReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ThreadID, validation.Required),
	)
}

// GetDetail request body
type GetDetailCommentReq struct {
	ID string `param:"id"`
}

func (request GetDetailCommentReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// CommentReportReq represent create request body
type CommentReportReq struct {
	ID          string `param:"id"`
	ReasonID    string `json:"reason_id"`
	Description string `json:"description"`
	UserID      string
}

func (request CommentReportReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
		validation.Field(&request.ReasonID, validation.Required),
		validation.Field(&request.UserID, validation.Required),
	)
}
