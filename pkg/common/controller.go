package common

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
}

type APIError struct {
	Code     string      `json:"code,omitempty"`
	Field    string      `json:"field,omitempty"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"`
}

func NewAPIError(code, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func (ae *APIError) Error() string {
	if ae.Internal == nil {
		return fmt.Sprintf("code=%s, message=%s", ae.Code, ae.Message)
	}
	return fmt.Sprintf("code=%s, message=%s, internal=%v", ae.Code, ae.Message, ae.Internal)
}

func (ae *APIError) SetInternal(err error) *APIError {
	ae.Internal = err
	return ae
}

type APIResponse struct {
	Code    int         `json:"code,omitempty"      extensions:"x-nullable,x-omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"      extensions:"x-nullable,x-omitempty"`
	Errors  interface{} `json:"errors,omitempty"    extensions:"x-nullable,x-omitempty"`
	Error   interface{} `json:"error,omitempty"     extensions:"x-nullable,x-omitempty"`
}

func (c *Controller) Ok(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	})
}

func (c *Controller) Unauthorized(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusUnauthorized, APIResponse{
		Code:    http.StatusUnauthorized,
		Message: err.Error(),
	})
}

func (c *Controller) WrapBadRequest(ctx echo.Context, response *APIResponse) error {
	if response.Code == 0 {
		response.Code = http.StatusBadRequest
	}
	return ctx.JSON(http.StatusBadRequest, &response)
}

func (c *Controller) InternalServerError(ctx echo.Context, response *APIResponse) error {
	if response.Code == 0 {
		response.Code = http.StatusInternalServerError
	}
	return ctx.JSON(http.StatusInternalServerError, &response)
}

func (c *Controller) DataNotFound(ctx echo.Context, response *APIResponse) error {
	if response.Code == 0 {
		response.Code = http.StatusNotFound
	}
	return ctx.JSON(http.StatusNotFound, &response)
}
