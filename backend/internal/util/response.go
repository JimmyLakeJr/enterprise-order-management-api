package util

import "github.com/labstack/echo/v4"

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginationResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    any            `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

func Success(c echo.Context, statusCode int, message string, data any) error {
	return c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, statusCode int, message string, errors any) error {
	if errors == nil {
		errors = map[string]any{}
	}

	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func Pagination(c echo.Context, statusCode int, message string, data any, meta PaginationMeta) error {
	return c.JSON(statusCode, PaginationResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func NewPaginationMeta(page int, limit int, total int) PaginationMeta {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
