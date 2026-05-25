package controllers

import (
	"backend/filters"
	"backend/pkg/apperr"
	"backend/pkg/ctxutil"
	"backend/pkg/response"
	"errors"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) User() *ctxutil.UserContext {
	return filters.GetUser(c.Ctx)
}

func (c *BaseController) Success(data interface{}) {
	c.Data["json"] = response.OK(data)
	c.ServeJSON()
}

func (c *BaseController) Fail(err error) {
	code, msg, status := "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError
	if ae, ok := err.(*apperr.AppError); ok {
		code = ae.Code
		if ae.Message != "" {
			msg = ae.Message
		}
		err = ae.Err
	}
	switch {
	case errors.Is(err, apperr.ErrUnauthorized):
		status = 401
	case errors.Is(err, apperr.ErrForbidden):
		status = 403
	case errors.Is(err, apperr.ErrNotFound):
		status = 404
	case errors.Is(err, apperr.ErrNoImageAvailable), errors.Is(err, apperr.ErrInvalidTransition), errors.Is(err, apperr.ErrDiscardLimit):
		status = 409
	case errors.Is(err, apperr.ErrValidation):
		status = 400
	}
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = response.Err(code, msg)
	c.ServeJSON()
}

func (c *BaseController) PageQuery() (int, int) {
	page, _ := c.GetInt("page", 1)
	size, _ := c.GetInt("page_size", 20)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}

func PageData(items interface{}, total int64, page, size int) map[string]interface{} {
	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": size,
	}
}
