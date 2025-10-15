package http

import (
	"net/http"
	"prakarsa-app/delivery/middleware"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	CommentUC domain.CommentUsecase
}

// NewCommentHandler will initialize the todo resources endpoint
func NewCommentHandler(e *echo.Echo, middleware *middleware.Middleware, commentUC domain.CommentUsecase) {
	handler := &CommentHandler{
		CommentUC: commentUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/comments", handler.Create)
	apiV1.PATCH("/comments/:id", handler.Update)
	apiV1.DELETE("/comments/:id", handler.Delete)
	apiV1.GET("/comments", handler.GetList)
	apiV1.GET("/comments/:id", handler.GetDetail)
	apiV1.POST("/comments/:id/report", handler.CommentReport)
}

func (h *CommentHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateCommentReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, err := h.CommentUC.Create(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Comment successfully created",
			"data":    res,
		})
	}
}

func (h *CommentHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.UpdateCommentReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.CommentUC.Update(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Comment successfully updated",
	})
}

func (h *CommentHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.DeleteCommentReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if rowsAffected, err := h.CommentUC.Delete(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Comment successfully deleted",
			"data": map[string]int64{
				"rows_affected": rowsAffected,
			},
		})
	}
}

func (h *CommentHandler) GetList(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.GetListCommentReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}
	
	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, meta, err := h.CommentUC.GetList(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Comments successfully retrieved",
			"data":    res,
			"meta":    meta,
		})
	}
}

func (h *CommentHandler) GetDetail(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.GetDetailCommentReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, err := h.CommentUC.GetDetail(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Comment successfully retrieved",
			"data":    res,
		})
	}
}

func (h *CommentHandler) CommentReport(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CommentReportReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.CommentUC.CommentReport(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Comment report successfully created",
		})
	}
}
