package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadSuccess(t *testing.T) {
	setValidEnv(t)

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, "development", cfg.AppEnv)
	require.Equal(t, "8080", cfg.AppPort)
	require.Equal(t, 15, cfg.AccessTokenExpireMinutes)
	require.Equal(t, 7, cfg.RefreshTokenExpireDays)
	require.Equal(t, "postgres://postgres:postgres@localhost:5432/enterprise_order_management?sslmode=disable", cfg.DatabaseURL())
}

func TestLoadMissingRequiredConfig(t *testing.T) {
	setValidEnv(t)
	t.Setenv("JWT_ACCESS_SECRET", "")

	_, err := Load()

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "JWT_ACCESS_SECRET"))
}

func TestLoadInvalidTokenExpireConfig(t *testing.T) {
	setValidEnv(t)
	t.Setenv("ACCESS_TOKEN_EXPIRE_MINUTES", "0")

	_, err := Load()

	require.Error(t, err)
	require.EqualError(t, err, "ACCESS_TOKEN_EXPIRE_MINUTES must be greater than 0")
}

func setValidEnv(t *testing.T) {
	t.Helper()

	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("DB_NAME", "enterprise_order_management")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("JWT_ACCESS_SECRET", "test-access-secret")
	t.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret")
	t.Setenv("ACCESS_TOKEN_EXPIRE_MINUTES", "15")
	t.Setenv("REFRESH_TOKEN_EXPIRE_DAYS", "7")
	t.Setenv("FRONTEND_URL", "http://localhost:5173")
}
