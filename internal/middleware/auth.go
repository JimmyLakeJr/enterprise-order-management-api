package middleware

import (
	"strings"

	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/token"

	"github.com/labstack/echo/v4"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

func JWTAuth(accessSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return apperror.Unauthorized("missing authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return apperror.Unauthorized("invalid authorization header")
			}

			claims, err := token.Parse(parts[1], accessSecret)
			if err != nil {
				return apperror.Unauthorized("invalid access token")
			}

			c.Set(ContextUserID, claims.UserID)
			c.Set(ContextRole, claims.Role)
			return next(c)
		}
	}
}

func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, _ := c.Get(ContextRole).(string)
			if !allowed[role] {
				return apperror.Forbidden("you do not have permission to access this resource")
			}
			return next(c)
		}
	}
}

func CurrentUserID(c echo.Context) int64 {
	userID, _ := c.Get(ContextUserID).(int64)
	return userID
}

func CurrentRole(c echo.Context) string {
	role, _ := c.Get(ContextRole).(string)
	return role
}
