package handler

import (
	"enterprise-order-management-api/internal/dto"
	appmiddleware "enterprise-order-management-api/internal/middleware"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	payments service.PaymentService
}

func NewPaymentHandler(payments service.PaymentService) *PaymentHandler {
	return &PaymentHandler{payments: payments}
}

func (h *PaymentHandler) CreateZaloPay(c echo.Context) error {
	var req dto.CreateZaloPayPaymentRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.payments.CreateZaloPayPayment(c.Request().Context(), appmiddleware.CurrentUserID(c), req)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *PaymentHandler) ZaloPayCallback(c echo.Context) error {
	var req dto.ZaloPayCallbackRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := h.payments.HandleZaloPayCallback(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}

func (h *PaymentHandler) ZaloPayStatus(c echo.Context) error {
	res, err := h.payments.GetZaloPayStatus(c.Request().Context(), appmiddleware.CurrentUserID(c), c.Param("transactionId"))
	if err != nil {
		return err
	}
	return response.OK(c, res)
}
