package util

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AppError struct {
	StatusCode int
	Message    string
	Errors     any
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, message string, errors any) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}
}

func BadRequest(message string, errors any) *AppError {
	return NewAppError(http.StatusBadRequest, message, errors)
}

func Unauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, nil)
}

func Forbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, nil)
}

func NotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, nil)
}

func Conflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message, nil)
}

func InternalServerError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, nil)
}

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		_ = Error(c, appErr.StatusCode, appErr.Message, appErr.Errors)
		return
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		message := http.StatusText(echoErr.Code)
		if echoErr.Code == http.StatusBadRequest {
			message = "Bad request"
		}
		_ = Error(c, echoErr.Code, message, nil)
		return
	}

	c.Logger().Error(err)
	_ = Error(c, http.StatusInternalServerError, "Internal server error", nil)
}
