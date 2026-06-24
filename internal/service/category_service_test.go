package service

import (
	"context"
	"testing"

	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCategoryService_DeleteRejectsCategoryWithActiveProducts(t *testing.T) {
	softDeleteCalled := false
	repo := &mockCategoryRepository{
		hasActiveProductsFunc: func(context.Context, int64) (bool, error) {
			return true, nil
		},
		softDeleteFunc: func(context.Context, int64) error {
			softDeleteCalled = true
			return nil
		},
	}
	service := &categoryService{categories: repo}

	err := service.Delete(context.Background(), 1)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, "CONFLICT", appErr.Code)
	require.False(t, softDeleteCalled)
}

func TestCategoryService_DeleteSuccess(t *testing.T) {
	deletedID := int64(0)
	repo := &mockCategoryRepository{
		hasActiveProductsFunc: func(context.Context, int64) (bool, error) {
			return false, nil
		},
		softDeleteFunc: func(_ context.Context, id int64) error {
			deletedID = id
			return nil
		},
	}
	service := &categoryService{categories: repo}

	err := service.Delete(context.Background(), 7)

	require.NoError(t, err)
	require.Equal(t, int64(7), deletedID)
}

func TestCategoryService_DeleteNotFound(t *testing.T) {
	repo := &mockCategoryRepository{
		hasActiveProductsFunc: func(context.Context, int64) (bool, error) {
			return false, nil
		},
		softDeleteFunc: func(context.Context, int64) error {
			return pgx.ErrNoRows
		},
	}
	service := &categoryService{categories: repo}

	err := service.Delete(context.Background(), 99)

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, "NOT_FOUND", appErr.Code)
}

func TestCategoryService_RestoreSuccess(t *testing.T) {
	restored := false
	repo := &mockCategoryRepository{
		restoreFunc: func(context.Context, int64) error {
			restored = true
			return nil
		},
		findByIDFunc: func(context.Context, int64) (*model.Category, error) {
			return &model.Category{ID: 1, Name: "Office", IsActive: true}, nil
		},
	}
	service := &categoryService{categories: repo}

	res, err := service.Restore(context.Background(), 1)

	require.NoError(t, err)
	require.True(t, restored)
	require.True(t, res.IsActive)
}
