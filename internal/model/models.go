package model

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusShipping  = "shipping"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

const (
	PaymentProviderZaloPay = "zalopay"
)

const (
	PaymentStatusPending   = "pending"
	PaymentStatusPaid      = "paid"
	PaymentStatusFailed    = "failed"
	PaymentStatusCancelled = "cancelled"
	PaymentStatusExpired   = "expired"
)

type User struct {
	ID              int64
	Name            string
	Email           string
	Phone           string
	PasswordHash    string
	AvatarURL       string
	ProfileVideoURL string
	RoleID          int64
	Role            string
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type OAuthAccount struct {
	ID             int64
	UserID         int64
	Provider       string
	ProviderUserID string
	Email          string
	AvatarURL      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Category struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Product struct {
	ID          int64
	CategoryID  int64
	Name        string
	Description string
	Price       int64
	Stock       int
	ImageURL    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Category    *Category
}

type Order struct {
	ID          int64
	UserID      int64
	Status      string
	TotalAmount int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	User        *User
	Items       []OrderItem
}

type OrderItem struct {
	ID        int64
	OrderID   int64
	ProductID int64
	Quantity  int
	UnitPrice int64
	Subtotal  int64
	Product   *Product
}

type Payment struct {
	ID                    int64
	TransactionID         string
	OrderID               int64
	UserID                int64
	Provider              string
	ProviderTransactionID string
	AppTransactionID      string
	Amount                int64
	Currency              string
	Status                string
	PaymentURL            string
	RawRequest            string
	RawResponse           string
	RawCallback           string
	PaidAt                *time.Time
	ExpiredAt             *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
