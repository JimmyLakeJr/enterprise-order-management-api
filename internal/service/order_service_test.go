package service

import (
	"context"
	"testing"
	"time"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/repository"

	"github.com/stretchr/testify/require"
)

func TestOrderService_CreateSuccess(t *testing.T) {
	tx := &mockTx{}
	txBegin := &mockTxBeginner{tx: tx}
	orderRepo := &mockOrderRepository{
		products: map[int64]*model.Product{
			1: {ID: 1, Name: "Laptop", Price: 15000000, Stock: 5, IsActive: true},
		},
	}

	service := &orderService{
		db:      &mockQueryer{},
		txBegin: txBegin,
		orders:  orderRepo,
	}

	res, err := service.Create(context.Background(), 10, dto.CreateOrderRequest{
		Items: []dto.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})

	require.NoError(t, err)
	require.Equal(t, int64(1), res.ID)
	require.Equal(t, int64(30000000), res.TotalAmount)
	require.True(t, tx.committed)
	require.False(t, tx.rolledBack)
}

func TestOrderService_CreateEmptyItems(t *testing.T) {
	service := &orderService{
		db:      &mockQueryer{},
		txBegin: &mockTxBeginner{tx: &mockTx{}},
		orders:  &mockOrderRepository{},
	}

	res, err := service.Create(context.Background(), 10, dto.CreateOrderRequest{})

	require.Error(t, err)
	require.Nil(t, res)
}

func TestOrderService_CreateQuantityLessThanOrEqualZero(t *testing.T) {
	tx := &mockTx{}
	service := &orderService{
		db:      &mockQueryer{},
		txBegin: &mockTxBeginner{tx: tx},
		orders:  &mockOrderRepository{},
	}

	res, err := service.Create(context.Background(), 10, dto.CreateOrderRequest{
		Items: []dto.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 0},
		},
	})

	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, tx.rolledBack)
}

func TestOrderService_CreateNotEnoughStock(t *testing.T) {
	tx := &mockTx{}
	orderRepo := &mockOrderRepository{
		products: map[int64]*model.Product{
			1: {ID: 1, Name: "Laptop", Price: 15000000, Stock: 1, IsActive: true},
		},
	}
	service := &orderService{
		db:      &mockQueryer{},
		txBegin: &mockTxBeginner{tx: tx},
		orders:  orderRepo,
	}

	res, err := service.Create(context.Background(), 10, dto.CreateOrderRequest{
		Items: []dto.CreateOrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})

	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, tx.rolledBack)
}

func TestOrderService_CreateProductNotFound(t *testing.T) {
	tx := &mockTx{}
	service := &orderService{
		db:      &mockQueryer{},
		txBegin: &mockTxBeginner{tx: tx},
		orders:  &mockOrderRepository{products: map[int64]*model.Product{}},
	}

	res, err := service.Create(context.Background(), 10, dto.CreateOrderRequest{
		Items: []dto.CreateOrderItemRequest{
			{ProductID: 99, Quantity: 1},
		},
	})

	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, tx.rolledBack)
}

func TestCanChangeOrderStatus(t *testing.T) {
	tests := []struct {
		name    string
		current string
		next    string
		allowed bool
	}{
		{name: "pending to confirmed", current: model.OrderStatusPending, next: model.OrderStatusConfirmed, allowed: true},
		{name: "pending to cancelled", current: model.OrderStatusPending, next: model.OrderStatusCancelled, allowed: true},
		{name: "confirmed to shipping", current: model.OrderStatusConfirmed, next: model.OrderStatusShipping, allowed: true},
		{name: "confirmed to cancelled", current: model.OrderStatusConfirmed, next: model.OrderStatusCancelled, allowed: true},
		{name: "shipping to completed", current: model.OrderStatusShipping, next: model.OrderStatusCompleted, allowed: true},
		{name: "pending to completed", current: model.OrderStatusPending, next: model.OrderStatusCompleted},
		{name: "completed is terminal", current: model.OrderStatusCompleted, next: model.OrderStatusPending},
		{name: "cancelled is terminal", current: model.OrderStatusCancelled, next: model.OrderStatusPending},
		{name: "cancelled to shipping", current: model.OrderStatusCancelled, next: model.OrderStatusShipping},
		{name: "same status", current: model.OrderStatusPending, next: model.OrderStatusPending},
		{name: "unknown status", current: model.OrderStatusPending, next: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.allowed, canChangeOrderStatus(tt.current, tt.next))
		})
	}
}

func TestOrderService_UpdateStatusRejectsInvalidTransition(t *testing.T) {
	updateCalled := false
	orderRepo := &mockOrderRepository{
		findByIDFunc: func(context.Context, repository.Queryer, int64) (*model.Order, error) {
			return &model.Order{ID: 1, UserID: 10, Status: model.OrderStatusCompleted}, nil
		},
		updateStatusFunc: func(context.Context, repository.Queryer, int64, string) error {
			updateCalled = true
			return nil
		},
	}
	service := &orderService{db: &mockQueryer{}, orders: orderRepo}

	res, err := service.UpdateStatus(context.Background(), 1, model.OrderStatusPending)

	require.Error(t, err)
	require.Nil(t, res)
	require.False(t, updateCalled)
}

func TestOrderService_FindByIDEnforcesOwnership(t *testing.T) {
	orderRepo := &mockOrderRepository{
		findByIDFunc: func(context.Context, repository.Queryer, int64) (*model.Order, error) {
			return &model.Order{ID: 1, UserID: 20, Status: model.OrderStatusPending}, nil
		},
	}
	service := &orderService{db: &mockQueryer{}, orders: orderRepo}

	res, err := service.FindByID(context.Background(), 1, 10, model.RoleUser)

	require.Error(t, err)
	require.Nil(t, res)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, "FORBIDDEN", appErr.Code)
}

func TestOrderService_FindByIDAllowsAdmin(t *testing.T) {
	orderRepo := &mockOrderRepository{
		findByIDFunc: func(context.Context, repository.Queryer, int64) (*model.Order, error) {
			return &model.Order{ID: 1, UserID: 20, Status: model.OrderStatusPending}, nil
		},
	}
	service := &orderService{db: &mockQueryer{}, orders: orderRepo}

	res, err := service.FindByID(context.Background(), 1, 10, model.RoleAdmin)

	require.NoError(t, err)
	require.Equal(t, int64(1), res.ID)
}

func TestOrderService_ListForAdminReturnsMetaAndBatchedItems(t *testing.T) {
	orderRepo := &mockOrderRepository{
		listAllFunc: func(context.Context, repository.Queryer, dto.OrderListQuery) ([]model.Order, int64, error) {
			return []model.Order{
				{ID: 1, UserID: 11, Status: model.OrderStatusPending, TotalAmount: 100, User: &model.User{ID: 11, Name: "A", Email: "a@example.com"}},
				{ID: 2, UserID: 12, Status: model.OrderStatusConfirmed, TotalAmount: 200, User: &model.User{ID: 12, Name: "B", Email: "b@example.com"}},
			}, 2, nil
		},
		findItemsByOrderIDsFunc: func(context.Context, repository.Queryer, []int64) (map[int64][]model.OrderItem, error) {
			return map[int64][]model.OrderItem{
				1: {{OrderID: 1, ProductID: 1, Quantity: 1, UnitPrice: 100, Subtotal: 100, Product: &model.Product{Name: "Mouse"}}},
				2: {{OrderID: 2, ProductID: 2, Quantity: 2, UnitPrice: 100, Subtotal: 200, Product: &model.Product{Name: "Keyboard"}}},
			}, nil
		},
	}
	service := &orderService{db: &mockQueryer{}, orders: orderRepo}

	res, meta, err := service.List(context.Background(), 0, model.RoleAdmin, dto.OrderListQuery{Page: 1, Limit: 10})

	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Equal(t, 1, meta.TotalPages)
	require.EqualValues(t, 2, meta.Total)
	require.Equal(t, "Mouse", res[0].Items[0].Name)
}

func TestOrderService_ListForUserFiltersByOwner(t *testing.T) {
	orderRepo := &mockOrderRepository{
		listByUserIDFunc: func(_ context.Context, _ repository.Queryer, userID int64, query dto.OrderListQuery) ([]model.Order, int64, error) {
			require.Equal(t, int64(42), userID)
			require.Equal(t, model.OrderStatusPending, query.Status)
			return []model.Order{
				{ID: 3, UserID: 42, Status: model.OrderStatusPending, TotalAmount: 300, User: &model.User{ID: 42, Name: "Owner", Email: "owner@example.com"}},
			}, 1, nil
		},
		findItemsByOrderIDsFunc: func(context.Context, repository.Queryer, []int64) (map[int64][]model.OrderItem, error) {
			return map[int64][]model.OrderItem{
				3: {{OrderID: 3, ProductID: 9, Quantity: 3, UnitPrice: 100, Subtotal: 300}},
			}, nil
		},
	}
	service := &orderService{db: &mockQueryer{}, orders: orderRepo}

	res, meta, err := service.List(context.Background(), 42, model.RoleUser, dto.OrderListQuery{Page: 1, Limit: 10, Status: model.OrderStatusPending})

	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, 1, meta.Page)
	require.Equal(t, int64(1), meta.Total)
	require.Equal(t, int64(42), res[0].UserID)
}

func TestOrderService_ListRejectsInvalidStatusFilter(t *testing.T) {
	service := &orderService{db: &mockQueryer{}, orders: &mockOrderRepository{}}

	res, meta, err := service.List(context.Background(), 42, model.RoleUser, dto.OrderListQuery{Page: 1, Limit: 10, Status: "invalid"})

	require.Error(t, err)
	require.Nil(t, res)
	require.Equal(t, int64(0), meta.Total)
}

func TestOrderService_UpdateStatusSuccess(t *testing.T) {
	updated := false
	orderRepo := &mockOrderRepository{
		findByIDFunc: func(context.Context, repository.Queryer, int64) (*model.Order, error) {
			return &model.Order{
				ID:          1,
				UserID:      10,
				Status:      model.OrderStatusPending,
				TotalAmount: 100,
				User:        &model.User{ID: 10, Name: "Demo", Email: "demo@example.com"},
			}, nil
		},
		updateStatusFunc: func(context.Context, repository.Queryer, int64, string) error {
			updated = true
			return nil
		},
		findItemsByOrderIDFunc: func(context.Context, repository.Queryer, int64) ([]model.OrderItem, error) {
			return []model.OrderItem{{OrderID: 1, ProductID: 5, Quantity: 1, UnitPrice: 100, Subtotal: 100}}, nil
		},
	}
	service := &orderService{db: &mockQueryer{}, orders: orderRepo}

	res, err := service.UpdateStatus(context.Background(), 1, model.OrderStatusConfirmed)

	require.NoError(t, err)
	require.True(t, updated)
	require.Equal(t, model.OrderStatusConfirmed, res.Status)
}

func TestToOrderResponseIncludesTimestamps(t *testing.T) {
	createdAt := time.Date(2026, time.June, 19, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	res := ToOrderResponse(&model.Order{
		ID:        1,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})

	require.Equal(t, createdAt, res.CreatedAt)
	require.Equal(t, updatedAt, res.UpdatedAt)
}

func TestToOrderResponseIncludesUserSummary(t *testing.T) {
	res := ToOrderResponse(&model.Order{
		ID:     1,
		UserID: 9,
		User:   &model.User{ID: 9, Name: "Demo User", Email: "demo@example.com"},
	})

	require.NotNil(t, res.User)
	require.Equal(t, int64(9), res.User.ID)
	require.Equal(t, "Demo User", res.User.Name)
	require.Equal(t, "demo@example.com", res.User.Email)
}
