package handler

import (
	"net/http"

	"enterprise-order-management-api/internal/dto"
	appmiddleware "enterprise-order-management-api/internal/middleware"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	auth service.AuthService
}

func NewAuthHandler(auth service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.auth.Register(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Created(c, res)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.auth.Login(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.auth.Refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	var req dto.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	if err := h.auth.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return err
	}
	return response.Message(c, http.StatusOK, "Logged out successfully")
}

func (h *AuthHandler) Me(c echo.Context) error {
	res, err := h.auth.Me(c.Request().Context(), appmiddleware.CurrentUserID(c))
	if err != nil {
		return err
	}
	return response.OK(c, res)
}
