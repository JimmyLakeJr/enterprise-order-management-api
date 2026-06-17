package handler

import (
	"net/http"
	"strconv"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/service"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	products service.ProductService
}

func NewProductHandler(products service.ProductService) *ProductHandler {
	return &ProductHandler{products: products}
}

func (h *ProductHandler) Create(c echo.Context) error {
	var req dto.ProductRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.products.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Created(c, res)
}

func (h *ProductHandler) FindByID(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	res, err := h.products.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *ProductHandler) List(c echo.Context) error {
	query := dto.ProductListQuery{
		Page:       parseIntQuery(c, "page", 1),
		Limit:      parseIntQuery(c, "limit", 10),
		Search:     c.QueryParam("keyword"),
		CategoryID: parseInt64Query(c, "category_id", 0),
		MinPrice:   parseInt64Query(c, "min_price", 0),
		MaxPrice:   parseInt64Query(c, "max_price", 0),
	}

	res, meta, err := h.products.List(c.Request().Context(), query)
	if err != nil {
		return err
	}
	return response.Paginated(c, res, meta)
}

func (h *ProductHandler) Update(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.ProductRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.products.Update(c.Request().Context(), id, req)
	if err != nil {
		return err
	}
	return response.OK(c, res)
}

func (h *ProductHandler) Delete(c echo.Context) error {
	id, err := parseID(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.products.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return response.Message(c, http.StatusOK, "Product deleted successfully")
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, apperror.BadRequest("Invalid id")
	}
	return id, nil
}

func parseIntQuery(c echo.Context, key string, fallback int) int {
	value := c.QueryParam(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseInt64Query(c echo.Context, key string, fallback int64) int64 {
	value := c.QueryParam(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
