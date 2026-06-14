package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTAccessSecret          string
	JWTRefreshSecret         string
	AccessTokenExpireMinutes int
	RefreshTokenExpireDays   int

	FrontendURL string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	if err := validateRequiredEnv(); err != nil {
		return Config{}, err
	}

	accessTokenExpireMinutes, err := parseRequiredInt("ACCESS_TOKEN_EXPIRE_MINUTES")
	if err != nil {
		return Config{}, err
	}

	refreshTokenExpireDays, err := parseRequiredInt("REFRESH_TOKEN_EXPIRE_DAYS")
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:  getRequired("APP_ENV"),
		AppPort: getRequired("APP_PORT"),

		DBHost:     getRequired("DB_HOST"),
		DBPort:     getRequired("DB_PORT"),
		DBUser:     getRequired("DB_USER"),
		DBPassword: getRequired("DB_PASSWORD"),
		DBName:     getRequired("DB_NAME"),
		DBSSLMode:  getRequired("DB_SSLMODE"),

		JWTAccessSecret:          getRequired("JWT_ACCESS_SECRET"),
		JWTRefreshSecret:         getRequired("JWT_REFRESH_SECRET"),
		AccessTokenExpireMinutes: accessTokenExpireMinutes,
		RefreshTokenExpireDays:   refreshTokenExpireDays,

		FrontendURL: getRequired("FRONTEND_URL"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validateRequiredEnv() error {
	requiredKeys := []string{
		"APP_ENV",
		"APP_PORT",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
		"JWT_ACCESS_SECRET",
		"JWT_REFRESH_SECRET",
		"ACCESS_TOKEN_EXPIRE_MINUTES",
		"REFRESH_TOKEN_EXPIRE_DAYS",
		"FRONTEND_URL",
	}

	missingKeys := make([]string, 0)
	for _, key := range requiredKeys {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missingKeys, ", "))
	}

	return nil
}

func (cfg Config) Validate() error {
	missingKeys := make([]string, 0)

	requiredValues := map[string]string{
		"APP_ENV":            cfg.AppEnv,
		"APP_PORT":           cfg.AppPort,
		"DB_HOST":            cfg.DBHost,
		"DB_PORT":            cfg.DBPort,
		"DB_USER":            cfg.DBUser,
		"DB_PASSWORD":        cfg.DBPassword,
		"DB_NAME":            cfg.DBName,
		"DB_SSLMODE":         cfg.DBSSLMode,
		"JWT_ACCESS_SECRET":  cfg.JWTAccessSecret,
		"JWT_REFRESH_SECRET": cfg.JWTRefreshSecret,
		"FRONTEND_URL":       cfg.FrontendURL,
	}

	for key, value := range requiredValues {
		if strings.TrimSpace(value) == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missingKeys, ", "))
	}

	if cfg.AccessTokenExpireMinutes <= 0 {
		return fmt.Errorf("ACCESS_TOKEN_EXPIRE_MINUTES must be greater than 0")
	}

	if cfg.RefreshTokenExpireDays <= 0 {
		return fmt.Errorf("REFRESH_TOKEN_EXPIRE_DAYS must be greater than 0")
	}

	return nil
}

func (cfg Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
}

func getRequired(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func parseRequiredInt(key string) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer", key)
	}

	return number, nil
}
