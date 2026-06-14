package dto

import "time"

type UserResponse struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"role_id"`
	RoleName  string    `json:"role_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name" validate:"omitempty,min=2,max=100"`
	IsActive *bool   `json:"is_active" validate:"omitempty"`
	RoleID   *int64  `json:"role_id" validate:"omitempty,gt=0"`
}

type ProfileResponse struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	RoleName  string    `json:"role_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
