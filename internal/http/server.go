package http

import (
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
)

func NewServer(db *pgxpool.Pool, cfg config.Config) *echo.Echo {
	e := echo.New()
	e.Validator = appvalidator.New()
	e.HTTPErrorHandler = response.ErrorHandler

	e.Use(appmiddleware.Logger())
	e.Use(appmiddleware.Recovery())
	e.Use(appmiddleware.CORS(cfg.FrontendURL))

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

	public := api.Group("")
	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)
	public.POST("/auth/refresh-token", authHandler.Refresh)

	public.GET("/categories", categoryHandler.List)
	public.GET("/categories/:id", categoryHandler.FindByID)
	public.GET("/products", productHandler.List)
	public.GET("/products/:id", productHandler.FindByID)

	userProtected := api.Group("", appmiddleware.JWTAuth(cfg.JWTAccessSecret))
	userProtected.POST("/auth/logout", authHandler.Logout)
	userProtected.GET("/auth/me", authHandler.Me)
	userProtected.GET("/users/me/orders", orderHandler.MyOrders)

	orders := userProtected.Group("/orders")
	orders.POST("", orderHandler.Create)
	orders.GET("", orderHandler.List)
	orders.GET("/:id", orderHandler.FindByID)
	orders.PUT("/:id/status", orderHandler.UpdateStatus, appmiddleware.RequireRoles(model.RoleAdmin))

	users := api.Group("/users", appmiddleware.JWTAuth(cfg.JWTAccessSecret), appmiddleware.RequireRoles(model.RoleAdmin))
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.FindByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)

	categories := api.Group("/categories", appmiddleware.JWTAuth(cfg.JWTAccessSecret), appmiddleware.RequireRoles(model.RoleAdmin))
	categories.POST("", categoryHandler.Create)
	categories.PUT("/:id", categoryHandler.Update)
	categories.DELETE("/:id", categoryHandler.Delete)

	products := api.Group("/products", appmiddleware.JWTAuth(cfg.JWTAccessSecret), appmiddleware.RequireRoles(model.RoleAdmin))
	products.POST("", productHandler.Create)
	products.PUT("/:id", productHandler.Update)
	products.DELETE("/:id", productHandler.Delete)

	return e
}
