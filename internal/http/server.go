package http

import (
	"net/http"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/handler"
	appmiddleware "enterprise-order-management-api/internal/middleware"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/response"
	appvalidator "enterprise-order-management-api/internal/pkg/validator"
	"enterprise-order-management-api/internal/repository"
	"enterprise-order-management-api/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer(db *pgxpool.Pool, cfg config.Config) *echo.Echo {
	e := echo.New()
	e.Validator = appvalidator.New()
	e.HTTPErrorHandler = response.ErrorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.FrontendURL},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	orderService := service.NewOrderService(db, orderRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	e.GET("/health", func(c echo.Context) error {
		return response.OK(c, map[string]string{"status": "ok"})
	})

	api := e.Group("/api/v1")

	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/refresh", authHandler.Refresh)

	api.GET("/categories", categoryHandler.List)
	api.GET("/products", productHandler.List)
	api.GET("/products/:id", productHandler.FindByID)

	auth := api.Group("", appmiddleware.JWTAuth(cfg.JWTAccessSecret))
	auth.POST("/auth/logout", authHandler.Logout)
	auth.GET("/me", userHandler.Me)

	orders := auth.Group("/orders")
	orders.POST("", orderHandler.Create)
	orders.GET("", orderHandler.List)
	orders.PATCH("/:id/status", orderHandler.UpdateStatus, appmiddleware.RequireRoles(model.RoleAdmin))

	admin := auth.Group("/admin", appmiddleware.RequireRoles(model.RoleAdmin))
	admin.GET("/users", userHandler.List)
	admin.DELETE("/users/:id", userHandler.Delete)
	admin.POST("/categories", categoryHandler.Create)
	admin.PUT("/categories/:id", categoryHandler.Update)
	admin.DELETE("/categories/:id", categoryHandler.Delete)
	admin.POST("/products", productHandler.Create)
	admin.PUT("/products/:id", productHandler.Update)
	admin.DELETE("/products/:id", productHandler.Delete)

	return e
}
