package service

import (
	"context"
	"fmt"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService interface {
	Create(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	List(ctx context.Context, userID int64, role string) ([]dto.OrderResponse, error)
	FindByID(ctx context.Context, orderID int64, userID int64, role string) (*dto.OrderResponse, error)
	UpdateStatus(ctx context.Context, orderID int64, status string) (*dto.OrderResponse, error)
}

type orderService struct {
	db      repository.Queryer
	txBegin txBeginner
	orders  repository.OrderRepository
}

type txBeginner interface {
	Begin(ctx context.Context) (repository.Tx, error)
}

type pgxPoolTxBeginner struct {
	pool *pgxpool.Pool
}

func (b pgxPoolTxBeginner) Begin(ctx context.Context) (repository.Tx, error) {
	return b.pool.Begin(ctx)
}

func NewOrderService(db *pgxpool.Pool, orders repository.OrderRepository) OrderService {
	return &orderService{
		db:      db,
		txBegin: pgxPoolTxBeginner{pool: db},
		orders:  orders,
	}
}

func (s *orderService) Create(ctx context.Context, userID int64, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if len(req.Items) == 0 {
		return nil, apperror.BadRequest("Order must have at least one item")
	}

	tx, err := s.txBegin.Begin(ctx)
	if err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	quantities := mergeOrderItems(req.Items)
	orderItems := make([]model.OrderItem, 0, len(quantities))
	var total int64

	for productID, quantity := range quantities {
		if productID <= 0 {
			return nil, apperror.BadRequest("Product id must be greater than 0")
		}
		if quantity <= 0 {
			return nil, apperror.BadRequest("Quantity must be greater than 0")
		}

		product, err := s.orders.FindProductForUpdate(ctx, tx, productID)
		if err != nil {
			return nil, err
		}
		if product == nil || !product.IsActive {
			return nil, apperror.NotFound(fmt.Sprintf("Product %d not found", productID))
		}
		if product.Stock < quantity {
			return nil, apperror.BadRequest(fmt.Sprintf("Product %s does not have enough stock", product.Name))
		}

		subtotal := product.Price * int64(quantity)
		total += subtotal
		orderItems = append(orderItems, model.OrderItem{
			ProductID: product.ID,
			Quantity:  quantity,
			UnitPrice: product.Price,
			Subtotal:  subtotal,
			Product:   product,
		})
	}

	order := &model.Order{
		UserID:      userID,
		Status:      model.OrderStatusPending,
		TotalAmount: total,
	}
	if err := s.orders.CreateOrder(ctx, tx, order); err != nil {
		return nil, err
	}

	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := s.orders.CreateOrderItem(ctx, tx, &orderItems[i]); err != nil {
			return nil, err
		}
		if err := s.orders.DecreaseStock(ctx, tx, orderItems[i].ProductID, orderItems[i].Quantity); err != nil {
			if err == pgx.ErrNoRows {
				return nil, apperror.BadRequest("Product does not have enough stock")
			}
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	committed = true

	order.Items = orderItems
	response := ToOrderResponse(order)
	return &response, nil
}

func (s *orderService) List(ctx context.Context, userID int64, role string) ([]dto.OrderResponse, error) {
	var (
		orders []model.Order
		err    error
	)

	if role == model.RoleAdmin {
		orders, err = s.orders.ListAll(ctx, s.db)
	} else {
		orders, err = s.orders.ListByUserID(ctx, s.db, userID)
	}
	if err != nil {
		return nil, err
	}

	responses := make([]dto.OrderResponse, 0, len(orders))
	for i := range orders {
		items, err := s.orders.FindItemsByOrderID(ctx, s.db, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
		responses = append(responses, ToOrderResponse(&orders[i]))
	}

	return responses, nil
}

func (s *orderService) FindByID(ctx context.Context, orderID int64, userID int64, role string) (*dto.OrderResponse, error) {
	order, err := s.orders.FindByID(ctx, s.db, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, apperror.NotFound("Order not found")
	}
	if role != model.RoleAdmin && order.UserID != userID {
		return nil, apperror.Forbidden("You cannot view another user's order")
	}

	items, err := s.orders.FindItemsByOrderID(ctx, s.db, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	res := ToOrderResponse(order)
	return &res, nil
}

func (s *orderService) UpdateStatus(ctx context.Context, orderID int64, status string) (*dto.OrderResponse, error) {
	order, err := s.orders.FindByID(ctx, s.db, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, apperror.NotFound("Order not found")
	}
	if !canChangeOrderStatus(order.Status, status) {
		return nil, apperror.BadRequest("Invalid order status transition")
	}

	if err := s.orders.UpdateStatus(ctx, s.db, orderID, status); err != nil {
		return nil, err
	}

	order.Status = status
	items, err := s.orders.FindItemsByOrderID(ctx, s.db, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	res := ToOrderResponse(order)
	return &res, nil
}

func canChangeOrderStatus(current string, next string) bool {
	if !isValidOrderStatus(next) {
		return false
	}

	allowed := map[string][]string{
		model.OrderStatusPending:   {model.OrderStatusConfirmed, model.OrderStatusCancelled},
		model.OrderStatusConfirmed: {model.OrderStatusShipping, model.OrderStatusCancelled},
		model.OrderStatusShipping:  {model.OrderStatusCompleted},
	}

	for _, status := range allowed[current] {
		if status == next {
			return true
		}
	}
	return false
}

func isValidOrderStatus(status string) bool {
	switch status {
	case model.OrderStatusPending,
		model.OrderStatusConfirmed,
		model.OrderStatusShipping,
		model.OrderStatusCompleted,
		model.OrderStatusCancelled:
		return true
	default:
		return false
	}
}

func mergeOrderItems(items []dto.CreateOrderItemRequest) map[int64]int {
	quantities := make(map[int64]int)
	for _, item := range items {
		quantities[item.ProductID] += item.Quantity
	}
	return quantities
}

func ToOrderResponse(order *model.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		response := dto.OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		}
		if item.Product != nil {
			response.Name = item.Product.Name
		}
		items = append(items, response)
	}

	return dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Items:       items,
	}
}
