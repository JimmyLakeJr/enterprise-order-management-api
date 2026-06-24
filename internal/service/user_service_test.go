package service

import (
	"context"
	"testing"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"

	"github.com/stretchr/testify/require"
)

func TestUserServiceUpdateProfileOnlyChangesName(t *testing.T) {
	var updatedName string
	repo := &mockUserRepository{
		updateProfileNameFunc: func(_ context.Context, id int64, name string) error {
			require.Equal(t, int64(7), id)
			updatedName = name
			return nil
		},
		findByIDFunc: func(_ context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id, Name: updatedName, Email: "user@example.com", Role: model.RoleUser, IsActive: true}, nil
		},
	}
	service := NewUserService(repo, nil)

	res, err := service.UpdateProfile(context.Background(), 7, dto.UpdateProfileRequest{Name: "Tên mới"})

	require.NoError(t, err)
	require.Equal(t, "Tên mới", res.Name)
	require.Equal(t, "user@example.com", res.Email)
	require.Equal(t, model.RoleUser, res.Role)
}

func TestUserServiceUpdateProfileRejectsBlankName(t *testing.T) {
	service := NewUserService(&mockUserRepository{}, nil)

	res, err := service.UpdateProfile(context.Background(), 7, dto.UpdateProfileRequest{Name: "   "})

	require.Error(t, err)
	require.Nil(t, res)
}
