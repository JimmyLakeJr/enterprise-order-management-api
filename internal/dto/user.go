package dto

type UserListQuery struct {
	Page   int
	Limit  int
	Search string
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Phone string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Role  string `json:"role" validate:"required,oneof=admin user"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Phone string `json:"phone,omitempty" validate:"omitempty,max=20"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6,max=72"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=72"`
}

type UploadResponse struct {
	URL string `json:"url"`
}
