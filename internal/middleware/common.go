package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func CORS(frontendURL string) echo.MiddlewareFunc {
	return echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{frontendURL},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	})
}

func Logger() echo.MiddlewareFunc {
	return echomiddleware.LoggerWithConfig(echomiddleware.LoggerConfig{
		Format: "${time_rfc3339} method=${method} path=${path} status=${status} latency=${latency_human}\n",
	})
}

func Recovery() echo.MiddlewareFunc {
	return echomiddleware.Recover()
}
