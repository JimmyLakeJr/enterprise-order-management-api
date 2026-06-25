package handler

import (
	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	carts service.CartService
}

func NewCartHandler(carts service.CartService) *CartHandler {
	return &CartHandler{carts: carts}
}

func (h *CartHandler) Quote(c echo.Context) error {
	var req dto.CartQuoteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.carts.Quote(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}
