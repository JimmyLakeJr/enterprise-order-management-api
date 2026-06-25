package service

import (
	"context"
	"math"
	"mime/multipart"
	"strings"
	"unicode"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/password"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/repository"
	"enterprise-order-management-api/internal/storage"

	"github.com/jackc/pgx/v5"
)

type UserService interface {
	List(ctx context.Context, query dto.UserListQuery) ([]dto.UserResponse, response.Meta, error)
	FindByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, id int64, req dto.UpdateProfileRequest) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, id int64, req dto.ChangePasswordRequest) error
	UploadAvatar(ctx context.Context, id int64, file *multipart.FileHeader) (*dto.UserResponse, error)
	UploadProfileVideo(ctx context.Context, id int64, file *multipart.FileHeader) (*dto.UserResponse, error)
	Delete(ctx context.Context, id int64, currentUserID int64) error
}

type userService struct {
	users   repository.UserRepository
	storage *storage.LocalFileStorage
}

func NewUserService(users repository.UserRepository, fileStorage *storage.LocalFileStorage) UserService {
	return &userService{users: users, storage: fileStorage}
}

func (s *userService) List(ctx context.Context, query dto.UserListQuery) ([]dto.UserResponse, response.Meta, error) {
	query = normalizeUserListQuery(query)

	users, total, err := s.users.List(ctx, query)
	if err != nil {
		return nil, response.Meta{}, err
	}
	res := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		res = append(res, ToUserResponse(&users[i]))
	}

	meta := response.Meta{
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(query.Limit))),
	}
	return res, meta, nil
}

func (s *userService) FindByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}
	res := ToUserResponse(user)
	return &res, nil
}

func (s *userService) Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}

	exists, err := s.users.ExistsByEmailOtherUser(ctx, req.Email, id)
	email := normalizeEmail(req.Email)
	phone := normalizePhone(req.Phone)
	if err := validateProfilePhone(req.Phone, phone); err != nil {
		return nil, err
	}
	if email == "" && phone == "" {
		return nil, apperror.BadRequest("Email or phone is required")
	}

	if email != "" {
		exists, err = s.users.ExistsByEmailOtherUser(ctx, email, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperror.Conflict("Email already exists")
		}
	}

	if phone != "" {
		exists, err = s.users.ExistsByPhoneOtherUser(ctx, phone, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperror.Conflict("Phone already exists")
		}
	}

	user.Name = strings.TrimSpace(req.Name)
	user.Email = email
	user.Phone = phone
	user.Role = req.Role

	if err := s.users.Update(ctx, user); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.NotFound("User not found")
		}
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *userService) UpdateProfile(ctx context.Context, id int64, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}

	name := strings.TrimSpace(req.Name)
	if len(name) < 2 {
		return nil, apperror.BadRequest("Tên hiển thị phải có ít nhất 2 ký tự")
	}

	phone := normalizePhone(req.Phone)
	if err := validateProfilePhone(req.Phone, phone); err != nil {
		return nil, err
	}
	if phone != "" {
		exists, err := s.users.ExistsByPhoneOtherUser(ctx, phone, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperror.Conflict("Phone already exists")
		}
	}

	if err := s.users.UpdateProfile(ctx, id, name, phone); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.NotFound("User not found")
		}
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *userService) ChangePassword(ctx context.Context, id int64, req dto.ChangePasswordRequest) error {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return apperror.NotFound("User not found")
	}

	currentPassword := strings.TrimSpace(req.CurrentPassword)
	newPassword := strings.TrimSpace(req.NewPassword)

	if !password.Check(currentPassword, user.PasswordHash) {
		return apperror.BadRequest("Current password is incorrect")
	}
	if currentPassword == newPassword {
		return apperror.BadRequest("New password must be different from current password")
	}

	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.users.UpdatePasswordHash(ctx, id, hashedPassword); err != nil {
		if err == pgx.ErrNoRows {
			return apperror.NotFound("User not found")
		}
		return err
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id int64, currentUserID int64) error {
	if id == currentUserID {
		return apperror.BadRequest("Admin cannot delete own account")
	}

	if err := s.users.SoftDelete(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return apperror.NotFound("User not found")
		}
		return err
	}
	return nil
}

func (s *userService) UploadAvatar(ctx context.Context, id int64, file *multipart.FileHeader) (*dto.UserResponse, error) {
	if s.storage == nil {
		return nil, apperror.New(500, "UPLOAD_DISABLED", "Upload storage is not configured")
	}
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}

	uploadedURL, err := s.storage.SaveImage(file, "profile/avatars", 5*1024*1024)
	if err != nil {
		return nil, err
	}

	if err := s.users.UpdateAvatarURL(ctx, id, uploadedURL); err != nil {
		return nil, err
	}

	if oldURL := user.AvatarURL; oldURL != "" && oldURL != uploadedURL {
		_ = s.storage.DeleteManagedFile(oldURL)
	}

	return s.FindByID(ctx, id)
}

func (s *userService) UploadProfileVideo(ctx context.Context, id int64, file *multipart.FileHeader) (*dto.UserResponse, error) {
	if s.storage == nil {
		return nil, apperror.New(500, "UPLOAD_DISABLED", "Upload storage is not configured")
	}
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}

	uploadedURL, err := s.storage.SaveVideo(file, "profile/videos", 20*1024*1024)
	if err != nil {
		return nil, err
	}

	if err := s.users.UpdateProfileVideoURL(ctx, id, uploadedURL); err != nil {
		return nil, err
	}

	if oldURL := user.ProfileVideoURL; oldURL != "" && oldURL != uploadedURL {
		_ = s.storage.DeleteManagedFile(oldURL)
	}

	return s.FindByID(ctx, id)
}

func normalizeUserListQuery(query dto.UserListQuery) dto.UserListQuery {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	return query
}

func validateProfilePhone(raw string, normalized string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}

	for i, r := range trimmed {
		switch {
		case unicode.IsDigit(r):
		case r == '+' && i == 0:
		case r == ' ' || r == '(' || r == ')' || r == '-':
		default:
			return apperror.BadRequest("Phone format is invalid")
		}
	}

	digitsOnly := normalized
	if strings.HasPrefix(digitsOnly, "+") {
		digitsOnly = digitsOnly[1:]
	}
	if len(digitsOnly) < 9 || len(digitsOnly) > 15 {
		return apperror.BadRequest("Phone format is invalid")
	}

	return nil
}
