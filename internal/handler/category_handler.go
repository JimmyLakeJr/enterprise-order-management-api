package handler

import (
	"net/http"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categories service.CategoryService
}

func NewCategoryHandler(categories service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categories: categories}
}

func (h *CategoryHandler) Create(c echo.Context) error {
	var req dto.CategoryRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.categories.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Created(c, res)
}

func (h *CategoryHandler) List(c echo.Context) error {
	res, err := h.categories.List(c.Request().Context())
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *CategoryHandler) AdminList(c echo.Context) error {
	res, err := h.categories.AdminList(c.Request().Context(), c.QueryParam("status"))
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *CategoryHandler) FindByID(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	res, err := h.categories.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.CategoryRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.categories.Update(c.Request().Context(), id, req)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}
	if err := h.categories.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return response.Message(c, http.StatusOK, "Category deleted successfully")
}

func (h *CategoryHandler) Restore(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}
	res, err := h.categories.Restore(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}
