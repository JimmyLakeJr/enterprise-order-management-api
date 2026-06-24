package dto

import "time"

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
	Quantity  int   `json:"quantity" validate:"required,gt=0"`
}

type OrderResponse struct {
	ID          int64               `json:"id"`
	UserID      int64               `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount int64               `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	User        *OrderUserSummary   `json:"user,omitempty"`
	Items       []OrderItemResponse `json:"items"`
}

type OrderUserSummary struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed shipping completed cancelled"`
}

type OrderListQuery struct {
	Page   int
	Limit  int
	Status string
}

type OrderItemResponse struct {
	ProductID int64  `json:"product_id"`
	Name      string `json:"name,omitempty"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
	Subtotal  int64  `json:"subtotal"`
}
