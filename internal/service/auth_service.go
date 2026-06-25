package service

import (
	"context"
	"strings"
	"time"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/oauth"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/hasher"
	"enterprise-order-management-api/internal/pkg/password"
	"enterprise-order-management-api/internal/pkg/token"
	"enterprise-order-management-api/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	BeginGoogleLogin(ctx context.Context) (string, string, error)
	CompleteGoogleLogin(ctx context.Context, code string) (*dto.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	Me(ctx context.Context, userID int64) (*dto.UserResponse, error)
}

type authService struct {
	db            repository.TxBeginner
	users         repository.UserRepository
	oauthAccounts repository.OAuthAccountRepository
	google        oauth.GoogleProviderClient
	cfg           config.Config
}

func NewAuthService(
	db repository.TxBeginner,
	users repository.UserRepository,
	oauthAccounts repository.OAuthAccountRepository,
	google oauth.GoogleProviderClient,
	cfg config.Config,
) AuthService {
	return &authService{
		db:            db,
		users:         users,
		oauthAccounts: oauthAccounts,
		google:        google,
		cfg:           cfg,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := normalizeEmail(req.Email)
	phone := normalizePhone(req.Phone)

	if email == "" && phone == "" {
		return nil, apperror.BadRequest("Email or phone is required")
	}
	if email != "" {
		existing, err := s.users.FindByEmailAny(ctx, email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, apperror.Conflict("Email already exists")
		}
	}
	if phone != "" {
		existing, err := s.users.FindByPhoneAny(ctx, phone)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, apperror.Conflict("Phone already exists")
		}
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		Phone:        phone,
		PasswordHash: hashedPassword,
		Role:         model.RoleUser,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	identifier := strings.TrimSpace(req.Identifier)
	if identifier == "" {
		identifier = strings.TrimSpace(req.Email)
	}
	identifier = normalizeLoginIdentifier(identifier)
	if identifier == "" {
		return nil, apperror.BadRequest("Email or phone is required")
	}

	user, err := s.users.FindByIdentifier(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if user == nil || !password.Check(req.Password, user.PasswordHash) {
		return nil, apperror.Unauthorized("Invalid email/phone or password")
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	claims, err := token.Parse(refreshToken, s.cfg.JWTRefreshSecret)
	if err != nil {
		return nil, apperror.Unauthorized("Invalid refresh token")
	}

	tokenHash := hasher.SHA256(refreshToken)
	storedToken, err := s.users.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, apperror.Unauthorized("Refresh token not found")
	}
	if storedToken.UserID != claims.UserID {
		return nil, apperror.Unauthorized("Refresh token does not match user")
	}
	if storedToken.RevokedAt != nil {
		return nil, apperror.Unauthorized("Refresh token was revoked")
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, apperror.Unauthorized("Refresh token expired")
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
	if _, err := token.Parse(refreshToken, s.cfg.JWTRefreshSecret); err != nil {
		return apperror.Unauthorized("Invalid refresh token")
	}
	return s.users.RevokeRefreshToken(ctx, hasher.SHA256(refreshToken))
}

func (s *authService) Me(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}
	res := ToUserResponse(user)
	return &res, nil
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
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Phone:           user.Phone,
		AvatarURL:       user.AvatarURL,
		ProfileVideoURL: user.ProfileVideoURL,
		Role:            user.Role,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}

	var b strings.Builder
	for i, r := range phone {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
			continue
		}
		if r == '+' && i == 0 {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalizeLoginIdentifier(identifier string) string {
	trimmed := strings.TrimSpace(identifier)
	if strings.Contains(trimmed, "@") {
		return normalizeEmail(trimmed)
	}
	return normalizePhone(trimmed)
}
