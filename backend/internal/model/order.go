package model

import "time"

const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusShipping  = "shipping"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

type Order struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	TotalAmount int64     `db:"total_amount"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	User  *User       `db:"-"`
	Items []OrderItem `db:"-"`
}

type OrderItem struct {
	ID        int64     `db:"id"`
	OrderID   int64     `db:"order_id"`
	ProductID int64     `db:"product_id"`
	Quantity  int       `db:"quantity"`
	UnitPrice int64     `db:"unit_price"`
	Subtotal  int64     `db:"subtotal"`
	CreatedAt time.Time `db:"created_at"`

	Product *Product `db:"-"`
}
