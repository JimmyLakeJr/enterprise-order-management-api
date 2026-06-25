package service

import (
	"context"
	"testing"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/password"

	"github.com/stretchr/testify/require"
)

func TestUserServiceUpdateProfileChangesNameAndPhone(t *testing.T) {
	var updatedName string
	var updatedPhone string
	repo := &mockUserRepository{
		findByIDFunc: func(_ context.Context, id int64) (*model.User, error) {
			return &model.User{
				ID:           id,
				Name:         "Tên cũ",
				Email:        "user@example.com",
				Phone:        "0987000000",
				PasswordHash: "$2a$10$abcdefghijklmnopqrstuv",
				Role:         model.RoleUser,
				IsActive:     true,
			}, nil
		},
		updateProfileFunc: func(_ context.Context, id int64, name string, phone string) error {
			require.Equal(t, int64(7), id)
			updatedName = name
			updatedPhone = phone
			return nil
		},
		existsByPhoneOtherUserFunc: func(_ context.Context, phone string, userID int64) (bool, error) {
			require.Equal(t, "0987654321", phone)
			require.Equal(t, int64(7), userID)
			return false, nil
		},
	}
	service := NewUserService(repo, nil)

	res, err := service.UpdateProfile(context.Background(), 7, dto.UpdateProfileRequest{
		Name:  "Tên mới",
		Phone: "0987 654 321",
	})

	require.NoError(t, err)
	require.Equal(t, "Tên mới", updatedName)
	require.Equal(t, "0987654321", updatedPhone)
	require.Equal(t, "user@example.com", res.Email)
}

func TestUserServiceUpdateProfileRejectsDuplicatePhone(t *testing.T) {
	repo := &mockUserRepository{
		findByIDFunc: func(_ context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id, Name: "User", Email: "user@example.com", Role: model.RoleUser, IsActive: true}, nil
		},
		existsByPhoneOtherUserFunc: func(_ context.Context, phone string, userID int64) (bool, error) {
			return true, nil
		},
	}
	service := NewUserService(repo, nil)

	res, err := service.UpdateProfile(context.Background(), 7, dto.UpdateProfileRequest{Name: "Tên mới", Phone: "0987654321"})

	require.Error(t, err)
	require.Nil(t, res)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, "CONFLICT", appErr.Code)
}

func TestUserServiceUpdateProfileRejectsBlankName(t *testing.T) {
	service := NewUserService(&mockUserRepository{}, nil)

	res, err := service.UpdateProfile(context.Background(), 7, dto.UpdateProfileRequest{Name: "   "})

	require.Error(t, err)
	require.Nil(t, res)
}

func TestUserServiceChangePasswordSuccess(t *testing.T) {
	currentHash, err := password.Hash("old123456")
	require.NoError(t, err)

	updatedHash := ""
	repo := &mockUserRepository{
		findByIDFunc: func(_ context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id, PasswordHash: currentHash, Role: model.RoleUser, IsActive: true}, nil
		},
		updatePasswordHashFunc: func(_ context.Context, id int64, passwordHash string) error {
			require.Equal(t, int64(7), id)
			updatedHash = passwordHash
			return nil
		},
	}
	service := NewUserService(repo, nil)

	err = service.ChangePassword(context.Background(), 7, dto.ChangePasswordRequest{
		CurrentPassword: "old123456",
		NewPassword:     "new123456",
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedHash)
	require.True(t, password.Check("new123456", updatedHash))
}

func TestUserServiceChangePasswordRejectsWrongCurrentPassword(t *testing.T) {
	currentHash, err := password.Hash("old123456")
	require.NoError(t, err)

	repo := &mockUserRepository{
		findByIDFunc: func(_ context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id, PasswordHash: currentHash, Role: model.RoleUser, IsActive: true}, nil
		},
	}
	service := NewUserService(repo, nil)

	err = service.ChangePassword(context.Background(), 7, dto.ChangePasswordRequest{
		CurrentPassword: "wrong123",
		NewPassword:     "new123456",
	})

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, "BAD_REQUEST", appErr.Code)
}
