package dto

import "time"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type LoginRequest struct {
	Identifier string `json:"identifier,omitempty" validate:"omitempty,max=255"`
	Email      string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Password   string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	AvatarURL       string    `json:"avatar_url"`
	ProfileVideoURL string    `json:"profile_video_url"`
	Role            string    `json:"role"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
