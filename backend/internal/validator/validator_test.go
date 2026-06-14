package validator

import (
	"errors"
	"testing"

	"enterprise-order-management-api/backend/internal/dto"
	"enterprise-order-management-api/backend/internal/util"

	"github.com/stretchr/testify/require"
)

func TestValidateRequiredDTOs(t *testing.T) {
	negativePrice := int64(-1)
	negativeStock := -1
	shortName := "A"
	tests := []struct {
		name          string
		request       any
		expectedField string
	}{
		{
			name:          "register",
			request:       dto.RegisterRequest{},
			expectedField: "full_name",
		},
		{
			name:          "login",
			request:       dto.LoginRequest{},
			expectedField: "email",
		},
		{
			name:          "create category",
			request:       dto.CreateCategoryRequest{},
			expectedField: "name",
		},
		{
			name:          "update category",
			request:       dto.UpdateCategoryRequest{Name: &shortName},
			expectedField: "name",
		},
		{
			name:          "create product",
			request:       dto.CreateProductRequest{},
			expectedField: "category_id",
		},
		{
			name:          "update product",
			request:       dto.UpdateProductRequest{Price: &negativePrice, Stock: &negativeStock},
			expectedField: "price",
		},
		{
			name:          "create order",
			request:       dto.CreateOrderRequest{},
			expectedField: "items",
		},
		{
			name:          "update order status",
			request:       dto.UpdateOrderStatusRequest{Status: "done"},
			expectedField: "status",
		},
	}

	validator := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.request)
			appErr := requireValidationError(t, err)
			errorsMap, ok := appErr.Errors.(map[string]string)

			require.True(t, ok)
			require.Contains(t, errorsMap, tt.expectedField)
		})
	}
}

func TestValidateDTOsSuccess(t *testing.T) {
	price := int64(0)
	stock := 0
	name := "Updated name"
	description := "Updated description"
	isActive := true
	categoryID := int64(1)

	tests := []struct {
		name    string
		request any
	}{
		{
			name: "register",
			request: dto.RegisterRequest{
				FullName: "Nguyen Van A",
				Email:    "user@example.com",
				Password: "secret123",
			},
		},
		{
			name: "login",
			request: dto.LoginRequest{
				Email:    "user@example.com",
				Password: "secret123",
			},
		},
		{
			name: "create category",
			request: dto.CreateCategoryRequest{
				Name:        "Electronics",
				Description: "Electronic products",
			},
		},
		{
			name: "update category",
			request: dto.UpdateCategoryRequest{
				Name:        &name,
				Description: &description,
				IsActive:    &isActive,
			},
		},
		{
			name: "create product with zero price and stock",
			request: dto.CreateProductRequest{
				CategoryID: categoryID,
				Name:       "Keyboard",
				Price:      &price,
				Stock:      &stock,
			},
		},
		{
			name: "update product",
			request: dto.UpdateProductRequest{
				CategoryID: &categoryID,
				Name:       &name,
				Price:      &price,
				Stock:      &stock,
				IsActive:   &isActive,
			},
		},
		{
			name: "create order",
			request: dto.CreateOrderRequest{
				Items: []dto.CreateOrderItemRequest{
					{ProductID: 1, Quantity: 2},
				},
			},
		},
		{
			name: "update order status",
			request: dto.UpdateOrderStatusRequest{
				Status: "confirmed",
			},
		},
	}

	validator := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, validator.Validate(tt.request))
		})
	}
}

func TestValidateRegisterReturnsReadableMessages(t *testing.T) {
	validator := New()

	err := validator.Validate(dto.RegisterRequest{
		FullName: "A",
		Email:    "not-email",
		Password: "123",
	})

	appErr := requireValidationError(t, err)
	errorsMap, ok := appErr.Errors.(map[string]string)

	require.True(t, ok)
	require.Equal(t, "full_name must be at least 2 characters or items", errorsMap["full_name"])
	require.Equal(t, "email must be a valid email address", errorsMap["email"])
	require.Equal(t, "password must be at least 6 characters or items", errorsMap["password"])
}

func TestValidateOrderStatusOneOf(t *testing.T) {
	validator := New()

	err := validator.Validate(dto.UpdateOrderStatusRequest{Status: "done"})

	appErr := requireValidationError(t, err)
	errorsMap, ok := appErr.Errors.(map[string]string)

	require.True(t, ok)
	require.Equal(t, "status must be one of: pending confirmed shipping completed cancelled", errorsMap["status"])
}

func requireValidationError(t *testing.T, err error) *util.AppError {
	t.Helper()

	require.Error(t, err)

	var appErr *util.AppError
	require.True(t, errors.As(err, &appErr))
	require.Equal(t, "Validation failed", appErr.Message)
	require.Equal(t, 400, appErr.StatusCode)

	return appErr
}
