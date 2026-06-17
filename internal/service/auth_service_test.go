package service

import (
	"context"
	"testing"
	"time"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/password"

	"github.com/stretchr/testify/require"
)

func TestAuthService_RegisterSuccess(t *testing.T) {
	repo := &mockUserRepository{}
	repo.findByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
		return nil, nil
	}
	repo.createFunc = func(ctx context.Context, user *model.User) error {
		user.ID = 1
		user.Role = model.RoleUser
		user.IsActive = true
		return nil
	}
	repo.saveRefreshTokenFunc = func(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
		require.Equal(t, int64(1), userID)
		require.NotEmpty(t, tokenHash)
		return nil
	}

	service := NewAuthService(repo, testConfig())

	res, err := service.Register(context.Background(), dto.RegisterRequest{
		Name:     "Nguyen Van A",
		Email:    "user@example.com",
		Password: "123456",
	})

	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)
	require.Equal(t, "user@example.com", res.User.Email)
}

func TestAuthService_RegisterDuplicateEmail(t *testing.T) {
	repo := &mockUserRepository{}
	repo.findByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
		return &model.User{ID: 1, Email: email}, nil
	}

	service := NewAuthService(repo, testConfig())

	res, err := service.Register(context.Background(), dto.RegisterRequest{
		Name:     "Nguyen Van A",
		Email:    "user@example.com",
		Password: "123456",
	})

	require.Error(t, err)
	require.Nil(t, res)
}

func TestAuthService_LoginSuccess(t *testing.T) {
	hashedPassword, err := password.Hash("123456")
	require.NoError(t, err)

	repo := &mockUserRepository{}
	repo.findByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
		return &model.User{
			ID:           1,
			Name:         "Nguyen Van A",
			Email:        email,
			PasswordHash: hashedPassword,
			Role:         model.RoleUser,
			IsActive:     true,
		}, nil
	}
	repo.saveRefreshTokenFunc = func(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
		return nil
	}

	service := NewAuthService(repo, testConfig())

	res, err := service.Login(context.Background(), dto.LoginRequest{
		Email:    "user@example.com",
		Password: "123456",
	})

	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
	require.Equal(t, "user@example.com", res.User.Email)
}

func TestAuthService_LoginWrongPassword(t *testing.T) {
	hashedPassword, err := password.Hash("123456")
	require.NoError(t, err)

	repo := &mockUserRepository{}
	repo.findByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
		return &model.User{
			ID:           1,
			Email:        email,
			PasswordHash: hashedPassword,
			Role:         model.RoleUser,
		}, nil
	}

	service := NewAuthService(repo, testConfig())

	res, err := service.Login(context.Background(), dto.LoginRequest{
		Email:    "user@example.com",
		Password: "wrong-password",
	})

	require.Error(t, err)
	require.Nil(t, res)
}

func testConfig() config.Config {
	return config.Config{
		JWTAccessSecret:        "test-access-secret",
		JWTRefreshSecret:       "test-refresh-secret",
		AccessTokenExpiration:  15 * time.Minute,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
	}
}
