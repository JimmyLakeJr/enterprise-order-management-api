package dto

type CartQuoteRequest struct {
	Items []CartQuoteItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CartQuoteItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
	Quantity  int   `json:"quantity" validate:"required,gt=0"`
}

type CartQuoteResponse struct {
	Items          []CartQuoteItemResponse `json:"items"`
	Subtotal       int64                   `json:"subtotal"`
	DiscountAmount int64                   `json:"discount_amount"`
	ShippingFee    int64                   `json:"shipping_fee"`
	FinalAmount    int64                   `json:"final_amount"`
	Warnings       []string                `json:"warnings,omitempty"`
}

type CartQuoteItemResponse struct {
	ProductID      int64  `json:"product_id"`
	Name           string `json:"name"`
	ImageURL       string `json:"image_url"`
	UnitPrice      int64  `json:"unit_price"`
	Quantity       int    `json:"quantity"`
	Subtotal       int64  `json:"subtotal"`
	AvailableStock int    `json:"available_stock"`
	IsAvailable    bool   `json:"is_available"`
}
