package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	DatabaseURL             string
	BackendPublicURL        string
	UploadDir               string
	JWTAccessSecret         string
	JWTRefreshSecret        string
	FrontendURL             string
	FrontendAuthCallbackURL string
	GoogleClientID          string
	GoogleClientSecret      string
	GoogleRedirectURL       string
	OAuthStateSecret        string
	AccessTokenExpiration   time.Duration
	RefreshTokenExpiration  time.Duration
	ZaloPay                 ZaloPayConfig
}

type ZaloPayConfig struct {
	Enabled             bool
	Environment         string
	AppID               int64
	Key1                string
	Key2                string
	CreateOrderEndpoint string
	QueryOrderEndpoint  string
	RefundEndpoint      string
	RedirectURL         string
	CallbackURL         string
	DefaultBankCode     string
	Currency            string
	Timeout             time.Duration
}

func (c ZaloPayConfig) Ready() bool {
	return c.Enabled &&
		c.AppID > 0 &&
		c.Key1 != "" &&
		c.Key2 != "" &&
		c.CreateOrderEndpoint != "" &&
		c.QueryOrderEndpoint != "" &&
		c.CallbackURL != "" &&
		c.RedirectURL != ""
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:                    getEnv("PORT", "8080"),
		DatabaseURL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/enterprise_order_management?sslmode=disable"),
		BackendPublicURL:        getEnv("BACKEND_PUBLIC_URL", "http://localhost:8080"),
		UploadDir:               getEnv("UPLOAD_DIR", "uploads"),
		JWTAccessSecret:         getEnv("JWT_ACCESS_SECRET", "dev-access-secret-change-me"),
		JWTRefreshSecret:        getEnv("JWT_REFRESH_SECRET", "dev-refresh-secret-change-me"),
		FrontendURL:             getEnv("FRONTEND_URL", "http://localhost:5173"),
		FrontendAuthCallbackURL: getEnv("FRONTEND_AUTH_CALLBACK_URL", "http://localhost:5173/auth/google/callback"),
		GoogleClientID:          getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:      getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:       getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
		OAuthStateSecret:        getEnv("OAUTH_STATE_SECRET", "dev-google-oauth-state-secret"),
		AccessTokenExpiration:   time.Duration(getEnvInt("ACCESS_TOKEN_MINUTES", 15)) * time.Minute,
		RefreshTokenExpiration:  time.Duration(getEnvInt("REFRESH_TOKEN_HOURS", 168)) * time.Hour,
		ZaloPay: ZaloPayConfig{
			Enabled:             getEnvBool("ZALOPAY_ENABLED", false),
			Environment:         getEnv("ZALOPAY_ENV", "sandbox"),
			AppID:               getEnvInt64("ZALOPAY_APP_ID", 0),
			Key1:                getEnv("ZALOPAY_KEY1", ""),
			Key2:                getEnv("ZALOPAY_KEY2", ""),
			CreateOrderEndpoint: getEnv("ZALOPAY_CREATE_ORDER_ENDPOINT", ""),
			QueryOrderEndpoint:  getEnv("ZALOPAY_QUERY_ORDER_ENDPOINT", ""),
			RefundEndpoint:      getEnv("ZALOPAY_REFUND_ENDPOINT", ""),
			RedirectURL:         getEnv("ZALOPAY_REDIRECT_URL", "http://localhost:5173/payment/zalopay/return"),
			CallbackURL:         getEnv("ZALOPAY_CALLBACK_URL", "http://localhost:8080/api/v1/payments/zalopay/callback"),
			DefaultBankCode:     getEnv("ZALOPAY_DEFAULT_BANK_CODE", ""),
			Currency:            getEnv("ZALOPAY_CURRENCY", "VND"),
			Timeout:             time.Duration(getEnvInt("ZALOPAY_TIMEOUT_SECONDS", 30)) * time.Second,
		},
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

func getEnvInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}
