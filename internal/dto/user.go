package dto

type UserListQuery struct {
	Page   int
	Limit  int
	Search string
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email,max=255"`
	Role  string `json:"role" validate:"required,oneof=admin user"`
}

type UpdateProfileRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type UploadResponse struct {
	URL string `json:"url"`
}
