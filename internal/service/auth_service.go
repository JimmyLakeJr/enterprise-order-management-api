package service

import (
	"context"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/hasher"
	"enterprise-order-management-api/internal/pkg/password"
	"enterprise-order-management-api/internal/pkg/token"
	"enterprise-order-management-api/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	users repository.UserRepository
	cfg   config.Config
}

func NewAuthService(users repository.UserRepository, cfg config.Config) AuthService {
	return &authService{users: users, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existing, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperror.Conflict("Email already exists")
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         model.RoleUser,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || !password.Check(req.Password, user.PasswordHash) {
		return nil, apperror.Unauthorized("Invalid email or password")
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	claims, err := token.Parse(refreshToken, s.cfg.JWTRefreshSecret)
	if err != nil {
		return nil, apperror.Unauthorized("Invalid refresh token")
	}

	tokenHash := hasher.SHA256(refreshToken)
	userID, err := s.users.FindActiveRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if userID == 0 || userID != claims.UserID {
		return nil, apperror.Unauthorized("Refresh token was revoked or expired")
	}

	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.Unauthorized("User no longer exists")
	}

	if err := s.users.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.users.RevokeRefreshToken(ctx, hasher.SHA256(refreshToken))
}

func (s *authService) issueTokens(ctx context.Context, user *model.User) (*dto.AuthResponse, error) {
	accessToken, _, err := token.Generate(user.ID, user.Email, user.Role, s.cfg.JWTAccessSecret, s.cfg.AccessTokenExpiration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := token.Generate(user.ID, user.Email, user.Role, s.cfg.JWTRefreshSecret, s.cfg.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	if err := s.users.SaveRefreshToken(ctx, user.ID, hasher.SHA256(refreshToken), refreshExpiresAt); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         ToUserResponse(user),
	}, nil
}

func ToUserResponse(user *model.User) dto.UserResponse {
	if user == nil {
		return dto.UserResponse{}
	}
	return dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}
