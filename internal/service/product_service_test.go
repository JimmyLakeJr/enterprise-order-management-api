package service

import (
	"context"
	"testing"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"

	"github.com/stretchr/testify/require"
)

func TestProductService_CreateSuccess(t *testing.T) {
	categoryRepo := &mockCategoryRepository{
		findActiveByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return &model.Category{ID: id, Name: "Electronics", IsActive: true}, nil
		},
	}
	productRepo := &mockProductRepository{}
	productRepo.createFunc = func(ctx context.Context, product *model.Product) error {
		product.ID = 1
		return nil
	}
	productRepo.findByIDFunc = func(ctx context.Context, id int64) (*model.Product, error) {
		return &model.Product{
			ID:         id,
			CategoryID: 1,
			Name:       "Laptop",
			Price:      15000000,
			Stock:      10,
			IsActive:   true,
			Category:   &model.Category{ID: 1, Name: "Electronics", IsActive: true},
		}, nil
	}

	service := NewProductService(productRepo, categoryRepo)

	res, err := service.Create(context.Background(), dto.ProductRequest{
		CategoryID:  1,
		Name:        "Laptop",
		Description: "Business laptop",
		Price:       15000000,
		Stock:       10,
	})

	require.NoError(t, err)
	require.Equal(t, int64(1), res.ID)
	require.Equal(t, "Laptop", res.Name)
}

func TestProductService_CreateNegativePrice(t *testing.T) {
	service := NewProductService(&mockProductRepository{}, &mockCategoryRepository{})

	res, err := service.Create(context.Background(), dto.ProductRequest{
		CategoryID: 1,
		Name:       "Laptop",
		Price:      -1,
		Stock:      10,
	})

	require.Error(t, err)
	require.Nil(t, res)
}

func TestProductService_CreateNegativeStock(t *testing.T) {
	service := NewProductService(&mockProductRepository{}, &mockCategoryRepository{})

	res, err := service.Create(context.Background(), dto.ProductRequest{
		CategoryID: 1,
		Name:       "Laptop",
		Price:      1000,
		Stock:      -1,
	})

	require.Error(t, err)
	require.Nil(t, res)
}

func TestProductService_CreateCategoryNotFound(t *testing.T) {
	categoryRepo := &mockCategoryRepository{
		findActiveByIDFunc: func(ctx context.Context, id int64) (*model.Category, error) {
			return nil, nil
		},
	}

	service := NewProductService(&mockProductRepository{}, categoryRepo)

	res, err := service.Create(context.Background(), dto.ProductRequest{
		CategoryID: 99,
		Name:       "Laptop",
		Price:      1000,
		Stock:      10,
	})

	require.Error(t, err)
	require.Nil(t, res)
}
