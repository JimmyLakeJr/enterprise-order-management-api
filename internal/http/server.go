package http

import (
	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/handler"
	appmiddleware "enterprise-order-management-api/internal/middleware"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/oauth"
	"enterprise-order-management-api/internal/pkg/response"
	appvalidator "enterprise-order-management-api/internal/pkg/validator"
	"enterprise-order-management-api/internal/repository"
	"enterprise-order-management-api/internal/service"
	"enterprise-order-management-api/internal/storage"

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
	e.Static("/uploads", cfg.UploadDir)

	userRepo := repository.NewUserRepository(db)
	oauthRepo := repository.NewOAuthAccountRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	fileStorage := storage.NewLocalFileStorage(cfg.UploadDir, cfg.BackendPublicURL)
	googleProvider := oauth.NewGoogleProvider(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)

	authService := service.NewAuthService(repository.NewTxBeginner(db), userRepo, oauthRepo, googleProvider, cfg)
	userService := service.NewUserService(userRepo, fileStorage)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo, fileStorage)
	orderService := service.NewOrderService(db, orderRepo)
	cartService := service.NewCartService(productRepo)
	paymentService := service.NewPaymentService(db, orderRepo, paymentRepo, cfg)

	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	cartHandler := handler.NewCartHandler(cartService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	e.GET("/health", func(c echo.Context) error {
		return response.OK(c, map[string]string{"status": "ok"})
	})

	api := e.Group("/api/v1")

	public := api.Group("")
	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)
	public.POST("/auth/refresh-token", authHandler.Refresh)
	public.GET("/auth/google/login", authHandler.GoogleLogin)
	public.GET("/auth/google/callback", authHandler.GoogleCallback)

	public.GET("/categories", categoryHandler.List)
	public.GET("/categories/:id", categoryHandler.FindByID)
	public.GET("/products", productHandler.List)
	public.GET("/products/:id", productHandler.FindByID)
	public.POST("/cart/quote", cartHandler.Quote)
	public.POST("/payments/zalopay/callback", paymentHandler.ZaloPayCallback)

	userProtected := api.Group("", appmiddleware.JWTAuth(cfg.JWTAccessSecret))
	userProtected.POST("/auth/logout", authHandler.Logout)
	userProtected.GET("/auth/me", authHandler.Me)
	userProtected.PUT("/users/me", userHandler.UpdateMe)
	userProtected.POST("/users/me/password", userHandler.ChangePassword)
	userProtected.POST("/users/me/avatar", userHandler.UploadAvatar)
	userProtected.POST("/users/me/profile-video", userHandler.UploadProfileVideo)
	userProtected.GET("/users/me/orders", orderHandler.MyOrders)

	orders := userProtected.Group("/orders")
	orders.POST("", orderHandler.Create)
	orders.GET("", orderHandler.List)
	orders.GET("/:id", orderHandler.FindByID)
	orders.PUT("/:id/status", orderHandler.UpdateStatus, appmiddleware.RequireRoles(model.RoleAdmin))

	payments := userProtected.Group("/payments")
	payments.POST("/zalopay/create", paymentHandler.CreateZaloPay)
	payments.GET("/zalopay/status/:transactionId", paymentHandler.ZaloPayStatus)

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
	products.POST("/upload-image", productHandler.UploadImage)
	products.PUT("/:id", productHandler.Update)
	products.DELETE("/:id", productHandler.Delete)

	admin := api.Group("/admin", appmiddleware.JWTAuth(cfg.JWTAccessSecret), appmiddleware.RequireRoles(model.RoleAdmin))
	admin.GET("/categories", categoryHandler.AdminList)
	admin.PUT("/categories/:id/restore", categoryHandler.Restore)
	admin.GET("/products", productHandler.AdminList)
	admin.PUT("/products/:id/restore", productHandler.Restore)

	return e
}
