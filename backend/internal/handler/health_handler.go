package handler

import (
	"net/http"

	"enterprise-order-management-api/backend/internal/util"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c echo.Context) error {
	data := map[string]string{
		"status": "ok",
	}

	return util.Success(c, http.StatusOK, "Success", data)
}
