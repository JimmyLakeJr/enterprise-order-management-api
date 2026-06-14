package handler

import (
	"net/http"

	appmiddleware "enterprise-order-management-api/internal/middleware"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	users service.UserService
}

func NewUserHandler(users service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Me(c echo.Context) error {
	res, err := h.users.Me(c.Request().Context(), appmiddleware.CurrentUserID(c))
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *UserHandler) List(c echo.Context) error {
	res, err := h.users.List(c.Request().Context())
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}
	if err := h.users.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return response.Message(c, http.StatusOK, "User deleted successfully")
}
