package dto

import "time"

type CreateZaloPayPaymentRequest struct {
	OrderID int64  `json:"orderId" validate:"required,gt=0"`
	Method  string `json:"method" validate:"required,oneof=zalopay"`
}

type PaymentResponse struct {
	TransactionID string     `json:"transactionId"`
	OrderID       int64      `json:"orderId"`
	Provider      string     `json:"provider"`
	Status        string     `json:"status"`
	Amount        int64      `json:"amount"`
	Currency      string     `json:"currency"`
	PaymentURL    string     `json:"paymentUrl,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
}

type PaymentStatusResponse struct {
	TransactionID         string     `json:"transactionId"`
	OrderID               int64      `json:"orderId"`
	Provider              string     `json:"provider"`
	Status                string     `json:"status"`
	Amount                int64      `json:"amount"`
	Currency              string     `json:"currency"`
	PaymentURL            string     `json:"paymentUrl,omitempty"`
	AppTransactionID      string     `json:"appTransactionId,omitempty"`
	ProviderTransactionID string     `json:"providerTransactionId,omitempty"`
	ExpiresAt             *time.Time `json:"expiresAt,omitempty"`
	PaidAt                *time.Time `json:"paidAt,omitempty"`
}

type ZaloPayCallbackRequest struct {
	Data string `json:"data"`
	MAC  string `json:"mac"`
	Type int    `json:"type"`
}

type ZaloPayCallbackResponse struct {
	ReturnCode    int    `json:"return_code"`
	ReturnMessage string `json:"return_message"`
}
