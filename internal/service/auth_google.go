package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/oauth"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/password"
	"enterprise-order-management-api/internal/repository"
)

const googleStateTTL = 10 * time.Minute

func (s *authService) BeginGoogleLogin(ctx context.Context) (string, string, error) {
	_ = ctx
	if s.google == nil || !s.google.Enabled() {
		return "", "", apperror.New(503, "FEATURE_DISABLED", "Google login is not configured")
	}

	state, err := oauth.GenerateState(s.cfg.OAuthStateSecret, oauth.GoogleProvider, googleStateTTL)
	if err != nil {
		return "", "", err
	}

	return s.google.AuthCodeURL(state), state, nil
}

func (s *authService) CompleteGoogleLogin(ctx context.Context, code string) (*dto.AuthResponse, error) {
	if s.google == nil || !s.google.Enabled() {
		return nil, apperror.New(503, "FEATURE_DISABLED", "Google login is not configured")
	}

	userInfo, err := s.google.Exchange(ctx, code)
	if err != nil {
		return nil, apperror.Unauthorized("Google login failed")
	}
	if !userInfo.EmailVerified {
		return nil, apperror.Forbidden("Google account email is not verified")
	}
	if userInfo.ProviderUserID == "" || userInfo.Email == "" {
		return nil, apperror.BadRequest("Google account data is incomplete")
	}

	account, err := s.oauthAccounts.FindByProviderUserID(ctx, oauth.GoogleProvider, userInfo.ProviderUserID)
	if err != nil {
		return nil, err
	}
	if account != nil {
		user, err := s.users.FindByIDAny(ctx, account.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil || !user.IsActive {
			return nil, apperror.Forbidden("This user account is inactive")
		}
		return s.issueTokens(ctx, user)
	}

	existingUser, err := s.users.FindByEmailAny(ctx, userInfo.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil && !existingUser.IsActive {
		return nil, apperror.Forbidden("This user account is inactive")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	user := existingUser
	if user == nil {
		user, err = s.createGoogleUser(ctx, tx, userInfo)
		if err != nil {
			return nil, err
		}
	}

	if err := s.oauthAccounts.CreateWithQuerier(ctx, tx, &model.OAuthAccount{
		UserID:         user.ID,
		Provider:       oauth.GoogleProvider,
		ProviderUserID: userInfo.ProviderUserID,
		Email:          userInfo.Email,
		AvatarURL:      userInfo.AvatarURL,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) createGoogleUser(ctx context.Context, q repository.Queryer, userInfo *oauth.GoogleUserInfo) (*model.User, error) {
	passwordHash, err := password.Hash(randomCredentialSeed())
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         googleDisplayName(userInfo),
		Email:        userInfo.Email,
		PasswordHash: passwordHash,
		AvatarURL:    userInfo.AvatarURL,
		Role:         model.RoleUser,
	}

	if err := s.users.CreateWithQuerier(ctx, q, user); err != nil {
		return nil, err
	}

	return user, nil
}

func googleDisplayName(userInfo *oauth.GoogleUserInfo) string {
	name := strings.TrimSpace(userInfo.Name)
	if len(name) >= 2 {
		if len(name) > 100 {
			return name[:100]
		}
		return name
	}

	localPart := userInfo.Email
	if at := strings.Index(localPart, "@"); at > 0 {
		localPart = localPart[:at]
	}
	if len(localPart) < 2 {
		return "Google User"
	}
	if len(localPart) > 100 {
		return localPart[:100]
	}
	return localPart
}

func randomCredentialSeed() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return time.Now().UTC().Format(time.RFC3339Nano)
	}
	return hex.EncodeToString(raw[:])
}
