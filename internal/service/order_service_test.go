package service

import (
	"context"
	"testing"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"

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
