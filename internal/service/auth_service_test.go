package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/oauth"
	"enterprise-order-management-api/internal/pkg/password"
	"enterprise-order-management-api/internal/repository"

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

	service := NewAuthService(&mockTxBeginner{tx: &mockTx{}}, repo, &mockOAuthAccountRepository{}, &mockGoogleProvider{}, testConfig())

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

	service := NewAuthService(&mockTxBeginner{tx: &mockTx{}}, repo, &mockOAuthAccountRepository{}, &mockGoogleProvider{}, testConfig())

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

	service := NewAuthService(&mockTxBeginner{tx: &mockTx{}}, repo, &mockOAuthAccountRepository{}, &mockGoogleProvider{}, testConfig())

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

	service := NewAuthService(&mockTxBeginner{tx: &mockTx{}}, repo, &mockOAuthAccountRepository{}, &mockGoogleProvider{}, testConfig())

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
		OAuthStateSecret:       "test-oauth-state-secret",
		AccessTokenExpiration:  15 * time.Minute,
		RefreshTokenExpiration: 7 * 24 * time.Hour,
	}
}

func TestAuthService_BeginGoogleLoginReturnsURLAndState(t *testing.T) {
	service := NewAuthService(
		&mockTxBeginner{tx: &mockTx{}},
		&mockUserRepository{},
		&mockOAuthAccountRepository{},
		&mockGoogleProvider{
			enabledFunc: func() bool { return true },
			authCodeURLFunc: func(state string) string {
				require.NotEmpty(t, state)
				return "https://accounts.google.com/o/oauth2/auth?state=" + state
			},
		},
		testConfig(),
	)

	url, state, err := service.BeginGoogleLogin(context.Background())

	require.NoError(t, err)
	require.NotEmpty(t, state)
	require.Contains(t, url, state)
}

func TestAuthService_CompleteGoogleLoginCreatesNewUser(t *testing.T) {
	tx := &mockTx{}
	repo := &mockUserRepository{
		findByEmailAnyFunc: func(context.Context, string) (*model.User, error) {
			return nil, nil
		},
		createWithQuerierFunc: func(_ context.Context, _ repository.Queryer, user *model.User) error {
			user.ID = 15
			user.Role = model.RoleUser
			user.IsActive = true
			return nil
		},
		saveRefreshTokenFunc: func(context.Context, int64, string, time.Time) error {
			return nil
		},
	}
	oauthRepo := &mockOAuthAccountRepository{
		findByProviderUserIDFunc: func(context.Context, string, string) (*model.OAuthAccount, error) {
			return nil, nil
		},
	}
	service := NewAuthService(
		&mockTxBeginner{tx: tx},
		repo,
		oauthRepo,
		&mockGoogleProvider{
			exchangeFunc: func(context.Context, string) (*oauth.GoogleUserInfo, error) {
				return &oauth.GoogleUserInfo{
					ProviderUserID: "google-123",
					Email:          "new-google@example.com",
					Name:           "Google User",
					AvatarURL:      "https://example.com/avatar.png",
					EmailVerified:  true,
				}, nil
			},
		},
		testConfig(),
	)

	res, err := service.CompleteGoogleLogin(context.Background(), "oauth-code")

	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, tx.committed)
	require.Equal(t, "new-google@example.com", res.User.Email)
	require.Equal(t, "https://example.com/avatar.png", res.User.AvatarURL)
}

func TestAuthService_CompleteGoogleLoginLinksExistingUserByEmail(t *testing.T) {
	tx := &mockTx{}
	repo := &mockUserRepository{
		findByEmailAnyFunc: func(context.Context, string) (*model.User, error) {
			return &model.User{ID: 2, Name: "Local User", Email: "local@example.com", Role: model.RoleAdmin, IsActive: true}, nil
		},
		saveRefreshTokenFunc: func(context.Context, int64, string, time.Time) error {
			return nil
		},
	}
	linked := false
	oauthRepo := &mockOAuthAccountRepository{
		findByProviderUserIDFunc: func(context.Context, string, string) (*model.OAuthAccount, error) {
			return nil, nil
		},
		createWithQuerierFunc: func(context.Context, repository.Queryer, *model.OAuthAccount) error {
			linked = true
			return nil
		},
	}
	service := NewAuthService(
		&mockTxBeginner{tx: tx},
		repo,
		oauthRepo,
		&mockGoogleProvider{
			exchangeFunc: func(context.Context, string) (*oauth.GoogleUserInfo, error) {
				return &oauth.GoogleUserInfo{
					ProviderUserID: "google-456",
					Email:          "local@example.com",
					Name:           "Google Local",
					EmailVerified:  true,
				}, nil
			},
		},
		testConfig(),
	)

	res, err := service.CompleteGoogleLogin(context.Background(), "oauth-code")

	require.NoError(t, err)
	require.True(t, linked)
	require.True(t, tx.committed)
	require.Equal(t, model.RoleAdmin, res.User.Role)
}

func TestAuthService_CompleteGoogleLoginRejectsUnverifiedEmail(t *testing.T) {
	service := NewAuthService(
		&mockTxBeginner{tx: &mockTx{}},
		&mockUserRepository{},
		&mockOAuthAccountRepository{},
		&mockGoogleProvider{
			exchangeFunc: func(context.Context, string) (*oauth.GoogleUserInfo, error) {
				return &oauth.GoogleUserInfo{
					ProviderUserID: "google-789",
					Email:          "user@example.com",
					EmailVerified:  false,
				}, nil
			},
		},
		testConfig(),
	)

	res, err := service.CompleteGoogleLogin(context.Background(), "oauth-code")

	require.Error(t, err)
	require.Nil(t, res)
}

func TestAuthService_CompleteGoogleLoginRejectsInactiveUser(t *testing.T) {
	service := NewAuthService(
		&mockTxBeginner{tx: &mockTx{}},
		&mockUserRepository{
			findByIDAnyFunc: func(context.Context, int64) (*model.User, error) {
				return &model.User{ID: 8, Email: "user@example.com", IsActive: false}, nil
			},
		},
		&mockOAuthAccountRepository{
			findByProviderUserIDFunc: func(context.Context, string, string) (*model.OAuthAccount, error) {
				return &model.OAuthAccount{ID: 1, UserID: 8, Provider: oauth.GoogleProvider, ProviderUserID: "google-888"}, nil
			},
		},
		&mockGoogleProvider{
			exchangeFunc: func(context.Context, string) (*oauth.GoogleUserInfo, error) {
				return &oauth.GoogleUserInfo{
					ProviderUserID: "google-888",
					Email:          "user@example.com",
					EmailVerified:  true,
				}, nil
			},
		},
		testConfig(),
	)

	res, err := service.CompleteGoogleLogin(context.Background(), "oauth-code")

	require.Error(t, err)
	require.Nil(t, res)
}

func TestAuthService_CompleteGoogleLoginHandlesGoogleExchangeError(t *testing.T) {
	service := NewAuthService(
		&mockTxBeginner{tx: &mockTx{}},
		&mockUserRepository{},
		&mockOAuthAccountRepository{},
		&mockGoogleProvider{
			exchangeFunc: func(context.Context, string) (*oauth.GoogleUserInfo, error) {
				return nil, errors.New("exchange failed")
			},
		},
		testConfig(),
	)

	res, err := service.CompleteGoogleLogin(context.Background(), "oauth-code")

	require.Error(t, err)
	require.Nil(t, res)
}
