package main

import (
	"context"
	"log"
	"time"

	"enterprise-order-management-api/backend/internal/config"
	"enterprise-order-management-api/backend/internal/database"
	"enterprise-order-management-api/backend/internal/route"
	"enterprise-order-management-api/backend/internal/util"
	appvalidator "enterprise-order-management-api/backend/internal/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := database.ConnectDB(ctx, cfg)
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}
	defer database.CloseDB(dbPool)

	e := echo.New()
	e.Validator = appvalidator.New()
	e.HTTPErrorHandler = util.HTTPErrorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	route.Register(e)

	if err := e.Start(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
