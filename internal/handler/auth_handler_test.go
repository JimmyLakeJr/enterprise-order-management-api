package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/oauth"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type stubAuthService struct {
	beginGoogleLoginFunc    func(ctx context.Context) (string, string, error)
	completeGoogleLoginFunc func(ctx context.Context, code string) (*dto.AuthResponse, error)
}

func (s *stubAuthService) Register(context.Context, dto.RegisterRequest) (*dto.AuthResponse, error) {
	return nil, nil
}

func (s *stubAuthService) Login(context.Context, dto.LoginRequest) (*dto.AuthResponse, error) {
	return nil, nil
}

func (s *stubAuthService) BeginGoogleLogin(ctx context.Context) (string, string, error) {
	return s.beginGoogleLoginFunc(ctx)
}

func (s *stubAuthService) CompleteGoogleLogin(ctx context.Context, code string) (*dto.AuthResponse, error) {
	return s.completeGoogleLoginFunc(ctx, code)
}

func (s *stubAuthService) Refresh(context.Context, string) (*dto.AuthResponse, error) {
	return nil, nil
}

func (s *stubAuthService) Logout(context.Context, string) error {
	return nil
}

func (s *stubAuthService) Me(context.Context, int64) (*dto.UserResponse, error) {
	return nil, nil
}

func TestAuthHandler_GoogleCallbackRejectsInvalidState(t *testing.T) {
	e := echo.New()
	cfg := config.Config{
		FrontendAuthCallbackURL: "http://localhost:5173/auth/google/callback",
		OAuthStateSecret:        "state-secret",
	}
	service := &stubAuthService{
		beginGoogleLoginFunc: func(ctx context.Context) (string, string, error) { return "", "", nil },
		completeGoogleLoginFunc: func(ctx context.Context, code string) (*dto.AuthResponse, error) {
			t.Fatal("complete login should not be called when state is invalid")
			return nil, nil
		},
	}
	handler := NewAuthHandler(service, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?state=bad&code=oauth-code", nil)
	req.AddCookie(&http.Cookie{Name: googleOAuthStateCookie, Value: "other"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GoogleCallback(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	require.Contains(t, rec.Header().Get("Location"), "status=error")
	require.Contains(t, rec.Header().Get("Location"), "Invalid+OAuth+state")
}

func TestAuthHandler_GoogleCallbackSuccessRedirectsToFrontend(t *testing.T) {
	e := echo.New()
	state, err := oauth.GenerateState("state-secret", oauth.GoogleProvider, time.Minute)
	require.NoError(t, err)

	cfg := config.Config{
		FrontendAuthCallbackURL: "http://localhost:5173/auth/google/callback",
		OAuthStateSecret:        "state-secret",
	}
	service := &stubAuthService{
		beginGoogleLoginFunc: func(ctx context.Context) (string, string, error) { return "", "", nil },
		completeGoogleLoginFunc: func(ctx context.Context, code string) (*dto.AuthResponse, error) {
			require.Equal(t, "oauth-code", code)
			return &dto.AuthResponse{
				AccessToken:  "access",
				RefreshToken: "refresh",
			}, nil
		},
	}
	handler := NewAuthHandler(service, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?state="+state+"&code=oauth-code", nil)
	req.AddCookie(&http.Cookie{Name: googleOAuthStateCookie, Value: state})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = handler.GoogleCallback(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	location := rec.Header().Get("Location")
	require.Contains(t, location, "status=success")
	require.Contains(t, location, "access_token=access")
	require.Contains(t, location, "refresh_token=refresh")
}
