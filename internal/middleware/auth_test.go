package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"enterprise-order-management-api/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestRequireRolesRejectsUserWithoutAdminRole(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(ContextRole, "user")

	handler := RequireRoles("admin")(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	err := handler(c)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, http.StatusForbidden, appErr.StatusCode)
	require.Equal(t, "FORBIDDEN", appErr.Code)
}

func TestRequireRolesAllowsAdmin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set(ContextRole, "admin")

	handler := RequireRoles("admin")(func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	err := handler(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}
