package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// CreateCommentReq represent create request body
type CreateCommentReq struct {
	ThreadID string `json:"thread_id"`
	Content  string `json:"content"`
}

func (request CreateCommentReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ThreadID, validation.Required),
		validation.Field(&request.Content, validation.Required),
	)
}

// Update request body
type UpdateCommentReq struct {
	ID      string `param:"id"`
	Content string `json:"content"`
}

func (request UpdateCommentReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// Delete request body
type DeleteCommentReq struct {
	ID string `param:"id"`
}

func (request DeleteCommentReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// GetList request body
type GetListCommentReq struct {
	ThreadID string `query:"thread_id"`
	IsActive *bool  `query:"is_active"`
	PerPage  int64  `query:"per_page"`
	Page     int64  `query:"page"`
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
