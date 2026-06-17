package handler

import (
	"enterprise-order-management-api/internal/dto"
	appmiddleware "enterprise-order-management-api/internal/middleware"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orders service.OrderService
}

func NewOrderHandler(orders service.OrderService) *OrderHandler {
	return &OrderHandler{orders: orders}
}

func (h *OrderHandler) Create(c echo.Context) error {
	var req dto.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.orders.Create(c.Request().Context(), appmiddleware.CurrentUserID(c), req)
	if err != nil {
		return err
	}
	return response.Created(c, res)
}

func (h *OrderHandler) List(c echo.Context) error {
	res, err := h.orders.List(
		c.Request().Context(),
		appmiddleware.CurrentUserID(c),
		appmiddleware.CurrentRole(c),
	)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *OrderHandler) MyOrders(c echo.Context) error {
	res, err := h.orders.List(
		c.Request().Context(),
		appmiddleware.CurrentUserID(c),
		model.RoleUser,
	)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *OrderHandler) FindByID(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	res, err := h.orders.FindByID(
		c.Request().Context(),
		id,
		appmiddleware.CurrentUserID(c),
		appmiddleware.CurrentRole(c),
	)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *OrderHandler) UpdateStatus(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.orders.UpdateStatus(c.Request().Context(), id, req.Status)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}
