package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	JWTAccessSecret        string
	JWTRefreshSecret       string
	FrontendURL            string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:                   getEnv("PORT", "8080"),
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/enterprise_order_management?sslmode=disable"),
		JWTAccessSecret:        getEnv("JWT_ACCESS_SECRET", "dev-access-secret-change-me"),
		JWTRefreshSecret:       getEnv("JWT_REFRESH_SECRET", "dev-refresh-secret-change-me"),
		FrontendURL:            getEnv("FRONTEND_URL", "http://localhost:5173"),
		AccessTokenExpiration:  time.Duration(getEnvInt("ACCESS_TOKEN_MINUTES", 15)) * time.Minute,
		RefreshTokenExpiration: time.Duration(getEnvInt("REFRESH_TOKEN_HOURS", 168)) * time.Hour,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
