package response

import (
	"errors"
	"net/http"
	"strings"

	"enterprise-order-management-api/internal/pkg/apperror"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Body struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Success(c echo.Context, statusCode int, message string, data any) error {
	if message == "" {
		message = "Success"
	}
	return c.JSON(statusCode, Body{Success: true, Message: message, Data: data})
}

func OK(c echo.Context, data any) error {
	return Success(c, http.StatusOK, "Success", data)
}

func Created(c echo.Context, data any) error {
	return Success(c, http.StatusCreated, "Created successfully", data)
}

func Message(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, Body{Success: true, Message: message})
}

func Paginated(c echo.Context, data any, meta Meta) error {
	return c.JSON(http.StatusOK, Body{
		Success: true,
		Message: "Success",
		Data:    data,
		Meta:    &meta,
	})
}

func ErrorHandler(err error, c echo.Context) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		_ = c.JSON(appErr.StatusCode, Body{
			Success: false,
			Message: appErr.Message,
			Errors:  map[string]string{"code": appErr.Code},
		})
		return
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		_ = c.JSON(http.StatusBadRequest, Body{
			Success: false,
			Message: "Validation failed",
			Errors:  formatValidationErrors(validationErrors),
		})
		return
	}

	if echoErr, ok := err.(*echo.HTTPError); ok {
		_ = c.JSON(echoErr.Code, Body{
			Success: false,
			Message: http.StatusText(echoErr.Code),
			Errors:  map[string]string{"code": "HTTP_ERROR"},
		})
		return
	}

	_ = c.JSON(http.StatusInternalServerError, Body{
		Success: false,
		Message: "Unexpected server error",
		Errors:  map[string]string{"code": "INTERNAL_ERROR"},
	})
}

func formatValidationErrors(validationErrors validator.ValidationErrors) map[string]string {
	result := make(map[string]string, len(validationErrors))
	for _, fieldErr := range validationErrors {
		field := strings.ToLower(fieldErr.Field())
		result[field] = "failed validation rule: " + fieldErr.Tag()
	}
	return result
}
