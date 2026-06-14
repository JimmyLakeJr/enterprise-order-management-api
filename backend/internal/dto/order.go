package dto

import "time"

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
	Quantity  int   `json:"quantity" validate:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed shipping completed cancelled"`
}

type OrderResponse struct {
	ID          int64               `json:"id"`
	UserID      int64               `json:"user_id"`
	TotalAmount int64               `json:"total_amount"`
	Status      string              `json:"status"`
	Items       []OrderItemResponse `json:"items,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        int64            `json:"id"`
	OrderID   int64            `json:"order_id"`
	ProductID int64            `json:"product_id"`
	Product   *ProductResponse `json:"product,omitempty"`
	Quantity  int              `json:"quantity"`
	UnitPrice int64            `json:"unit_price"`
	Subtotal  int64            `json:"subtotal"`
	CreatedAt time.Time        `json:"created_at"`
}
