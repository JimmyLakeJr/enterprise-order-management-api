package response

import (
	"errors"
	"net/http"

	"enterprise-order-management-api/internal/pkg/apperror"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
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

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		appErr := mapDatabaseError(pgErr)
		_ = c.JSON(appErr.StatusCode, Body{
			Success: false,
			Message: appErr.Message,
			Errors:  map[string]string{"code": appErr.Code},
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
		field := fieldErr.Field()
		result[field] = "failed validation rule: " + fieldErr.Tag()
	}
	return result
}

func mapDatabaseError(pgErr *pgconn.PgError) *apperror.AppError {
	switch pgErr.Code {
	case "23505":
		return apperror.Conflict("Data already exists")
	case "23503":
		return apperror.BadRequest("Related data is invalid or no longer exists")
	case "23514", "22P02":
		return apperror.BadRequest("Data is invalid")
	default:
		return apperror.New(http.StatusInternalServerError, "INTERNAL_ERROR", "Unexpected server error")
	}
}
