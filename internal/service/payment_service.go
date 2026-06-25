package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentService interface {
	CreateZaloPayPayment(ctx context.Context, userID int64, req dto.CreateZaloPayPaymentRequest) (*dto.PaymentResponse, error)
	HandleZaloPayCallback(ctx context.Context, req dto.ZaloPayCallbackRequest) (*dto.ZaloPayCallbackResponse, error)
	GetZaloPayStatus(ctx context.Context, userID int64, transactionID string) (*dto.PaymentStatusResponse, error)
}

type paymentService struct {
	db        repository.Queryer
	cfg       config.Config
	orders    repository.OrderRepository
	payments  repository.PaymentRepository
	zaloPay   ZaloPayClient
	timeNowFn func() time.Time
}

type zaloPayCallbackPayload struct {
	AppID          int64  `json:"app_id"`
	AppTransID     string `json:"app_trans_id"`
	AppTime        int64  `json:"app_time"`
	AppUser        string `json:"app_user"`
	Amount         int64  `json:"amount"`
	EmbedData      string `json:"embed_data"`
	Item           string `json:"item"`
	ZPTransID      int64  `json:"zp_trans_id"`
	ServerTime     int64  `json:"server_time"`
	MerchantUserID string `json:"merchant_user_id"`
}

func NewPaymentService(
	db *pgxpool.Pool,
	orders repository.OrderRepository,
	payments repository.PaymentRepository,
	cfg config.Config,
) PaymentService {
	return &paymentService{
		db:        db,
		cfg:       cfg,
		orders:    orders,
		payments:  payments,
		zaloPay:   NewZaloPayClient(cfg.ZaloPay),
		timeNowFn: time.Now,
	}
}

func (s *paymentService) CreateZaloPayPayment(ctx context.Context, userID int64, req dto.CreateZaloPayPaymentRequest) (*dto.PaymentResponse, error) {
	if req.Method != model.PaymentProviderZaloPay {
		return nil, apperror.BadRequest("Unsupported payment method")
	}
	if !s.cfg.ZaloPay.Enabled {
		return nil, apperror.New(http.StatusServiceUnavailable, "PAYMENT_PROVIDER_DISABLED", "ZaloPay is disabled in the current environment")
	}
	if !s.cfg.ZaloPay.Ready() {
		return nil, apperror.New(http.StatusServiceUnavailable, "PAYMENT_PROVIDER_NOT_READY", "ZaloPay is not fully configured on the server")
	}

	order, err := s.orders.FindByID(ctx, s.db, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, apperror.NotFound("Order not found")
	}
	if order.UserID != userID {
		return nil, apperror.Forbidden("You cannot pay for another user's order")
	}
	if order.Status == model.OrderStatusCancelled {
		return nil, apperror.BadRequest("Cancelled orders cannot be paid")
	}

	existing, err := s.payments.FindLatestByOrderIDAndProvider(ctx, s.db, order.ID, model.PaymentProviderZaloPay)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		switch existing.Status {
		case model.PaymentStatusPaid:
			return nil, apperror.Conflict("This order has already been paid")
		case model.PaymentStatusPending:
			if existing.ExpiredAt == nil || existing.ExpiredAt.After(s.timeNowFn()) {
				res := toPaymentResponse(existing)
				return &res, nil
			}
		}
	}

	transactionID, err := generateInternalTransactionID()
	if err != nil {
		return nil, err
	}
	appTransID := buildZaloPayAppTransID(order.ID, s.timeNowFn())
	expiresAt := s.timeNowFn().Add(15 * time.Minute)
	payment := &model.Payment{
		TransactionID:    transactionID,
		OrderID:          order.ID,
		UserID:           userID,
		Provider:         model.PaymentProviderZaloPay,
		AppTransactionID: appTransID,
		Amount:           order.TotalAmount,
		Currency:         s.cfg.ZaloPay.Currency,
		Status:           model.PaymentStatusPending,
		ExpiredAt:        &expiresAt,
	}
	if err := s.payments.Create(ctx, s.db, payment); err != nil {
		return nil, err
	}

	redirectURL := appendQueryParams(s.cfg.ZaloPay.RedirectURL, map[string]string{
		"transactionId": payment.TransactionID,
		"orderId":       fmt.Sprintf("%d", order.ID),
	})

	items := make([]map[string]any, 0, len(order.Items))
	if len(order.Items) == 0 {
		order.Items, _ = s.orders.FindItemsByOrderID(ctx, s.db, order.ID)
	}
	for _, item := range order.Items {
		items = append(items, map[string]any{
			"itemid":       fmt.Sprintf("%d", item.ProductID),
			"itemname":     paymentItemName(item),
			"itemprice":    item.UnitPrice,
			"itemquantity": item.Quantity,
		})
	}

	createReq := ZaloPayCreateOrderRequest{
		AppUser:      fmt.Sprintf("user-%d", userID),
		AppTransID:   appTransID,
		Amount:       order.TotalAmount,
		Description:  fmt.Sprintf("Thanh toan don hang #%d", order.ID),
		Items:        items,
		RedirectURL:  redirectURL,
		CallbackURL:  s.cfg.ZaloPay.CallbackURL,
		DefaultBank:  s.cfg.ZaloPay.DefaultBankCode,
		ExpiresAfter: 15 * time.Minute,
		ReferenceData: map[string]any{
			"internal_transaction_id": payment.TransactionID,
			"order_id":                order.ID,
		},
	}

	zaloRes, rawResponse, err := s.zaloPay.CreateOrder(ctx, createReq)
	rawRequestBytes, _ := json.Marshal(createReq)
	payment.RawRequest = string(rawRequestBytes)
	payment.RawResponse = rawResponse
	if err != nil {
		payment.Status = model.PaymentStatusFailed
		_ = s.payments.UpdateGatewayInitialization(ctx, s.db, payment)
		return nil, apperror.New(http.StatusBadGateway, "PAYMENT_PROVIDER_ERROR", "Could not initialize ZaloPay payment")
	}

	if zaloRes.ReturnCode != 1 && zaloRes.ReturnCode != 3 {
		payment.Status = model.PaymentStatusFailed
		payment.PaymentURL = zaloRes.OrderURL
		_ = s.payments.UpdateGatewayInitialization(ctx, s.db, payment)
		return nil, apperror.New(http.StatusBadGateway, "PAYMENT_PROVIDER_REJECTED", zaloRes.SubReturnMessageOrFallback())
	}

	if zaloRes.OrderURL == "" && zaloRes.QRCode == "" {
		payment.Status = model.PaymentStatusFailed
		_ = s.payments.UpdateGatewayInitialization(ctx, s.db, payment)
		return nil, apperror.New(http.StatusBadGateway, "PAYMENT_PROVIDER_INVALID_RESPONSE", "ZaloPay did not return a payment URL")
	}

	payment.PaymentURL = zaloRes.OrderURL
	payment.RawResponse = rawResponse
	if err := s.payments.UpdateGatewayInitialization(ctx, s.db, payment); err != nil {
		return nil, err
	}

	res := toPaymentResponse(payment)
	return &res, nil
}

func (s *paymentService) HandleZaloPayCallback(ctx context.Context, req dto.ZaloPayCallbackRequest) (*dto.ZaloPayCallbackResponse, error) {
	if !s.cfg.ZaloPay.Enabled || !s.cfg.ZaloPay.Ready() {
		return &dto.ZaloPayCallbackResponse{ReturnCode: 2, ReturnMessage: "Invalid"}, nil
	}
	if req.Data == "" || req.MAC == "" {
		return &dto.ZaloPayCallbackResponse{ReturnCode: 2, ReturnMessage: "Invalid"}, nil
	}
	if !s.zaloPay.VerifyCallback(req.Data, req.MAC) {
		return &dto.ZaloPayCallbackResponse{ReturnCode: 2, ReturnMessage: "Invalid"}, nil
	}

	var payload zaloPayCallbackPayload
	if err := json.Unmarshal([]byte(req.Data), &payload); err != nil {
		return &dto.ZaloPayCallbackResponse{ReturnCode: 0, ReturnMessage: "Retry"}, nil
	}

	payment, err := s.payments.FindByAppTransactionID(ctx, s.db, payload.AppTransID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return &dto.ZaloPayCallbackResponse{ReturnCode: 2, ReturnMessage: "Invalid"}, nil
	}
	if payment.Status == model.PaymentStatusPaid {
		return &dto.ZaloPayCallbackResponse{ReturnCode: 1, ReturnMessage: "Success"}, nil
	}

	payment.RawCallback = req.Data
	payment.ProviderTransactionID = fmt.Sprintf("%d", payload.ZPTransID)
	paidAt := time.UnixMilli(payload.ServerTime)
	payment.PaidAt = &paidAt
	payment.Status = model.PaymentStatusPaid
	if err := s.payments.UpdateStatus(ctx, s.db, payment); err != nil {
		return &dto.ZaloPayCallbackResponse{ReturnCode: 0, ReturnMessage: "Retry"}, nil
	}

	order, err := s.orders.FindByID(ctx, s.db, payment.OrderID)
	if err != nil {
		return &dto.ZaloPayCallbackResponse{ReturnCode: 0, ReturnMessage: "Retry"}, nil
	}
	if order != nil && order.Status == model.OrderStatusPending {
		if err := s.orders.UpdateStatus(ctx, s.db, order.ID, model.OrderStatusConfirmed); err != nil {
			return &dto.ZaloPayCallbackResponse{ReturnCode: 0, ReturnMessage: "Retry"}, nil
		}
	}

	return &dto.ZaloPayCallbackResponse{ReturnCode: 1, ReturnMessage: "Success"}, nil
}

func (s *paymentService) GetZaloPayStatus(ctx context.Context, userID int64, transactionID string) (*dto.PaymentStatusResponse, error) {
	payment, err := s.payments.FindByTransactionID(ctx, s.db, transactionID)
	if err != nil {
		return nil, err
	}
	if payment == nil || payment.Provider != model.PaymentProviderZaloPay {
		return nil, apperror.NotFound("Payment not found")
	}
	if payment.UserID != userID {
		return nil, apperror.Forbidden("You cannot view another user's payment")
	}

	if payment.Status == model.PaymentStatusPending && s.cfg.ZaloPay.Ready() {
		zaloRes, rawResponse, err := s.zaloPay.QueryOrder(ctx, payment.AppTransactionID)
		if err == nil && zaloRes != nil {
			payment.RawResponse = rawResponse
			switch {
			case zaloRes.ReturnCode == 1 && !zaloRes.IsProcessing:
				paidAt := time.UnixMilli(zaloRes.ServerTime)
				payment.ProviderTransactionID = nonZeroInt64String(zaloRes.ZPTransID)
				payment.PaidAt = &paidAt
				payment.Status = model.PaymentStatusPaid
				_ = s.payments.UpdateStatus(ctx, s.db, payment)

				order, orderErr := s.orders.FindByID(ctx, s.db, payment.OrderID)
				if orderErr == nil && order != nil && order.Status == model.OrderStatusPending {
					_ = s.orders.UpdateStatus(ctx, s.db, order.ID, model.OrderStatusConfirmed)
				}
			case zaloRes.ReturnCode == 2:
				if payment.ExpiredAt != nil && payment.ExpiredAt.Before(s.timeNowFn()) {
					payment.Status = model.PaymentStatusExpired
				} else {
					payment.Status = model.PaymentStatusFailed
				}
				_ = s.payments.UpdateStatus(ctx, s.db, payment)
			}
		}
	}

	if payment.Status == model.PaymentStatusPending && payment.ExpiredAt != nil && payment.ExpiredAt.Before(s.timeNowFn()) {
		payment.Status = model.PaymentStatusExpired
		_ = s.payments.UpdateStatus(ctx, s.db, payment)
	}

	res := toPaymentStatusResponse(payment)
	return &res, nil
}

func toPaymentResponse(payment *model.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		TransactionID: payment.TransactionID,
		OrderID:       payment.OrderID,
		Provider:      payment.Provider,
		Status:        payment.Status,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentURL:    payment.PaymentURL,
		ExpiresAt:     payment.ExpiredAt,
	}
}

func toPaymentStatusResponse(payment *model.Payment) dto.PaymentStatusResponse {
	return dto.PaymentStatusResponse{
		TransactionID:         payment.TransactionID,
		OrderID:               payment.OrderID,
		Provider:              payment.Provider,
		Status:                payment.Status,
		Amount:                payment.Amount,
		Currency:              payment.Currency,
		PaymentURL:            payment.PaymentURL,
		AppTransactionID:      payment.AppTransactionID,
		ProviderTransactionID: payment.ProviderTransactionID,
		ExpiresAt:             payment.ExpiredAt,
		PaidAt:                payment.PaidAt,
	}
}

func buildZaloPayAppTransID(orderID int64, now time.Time) string {
	location := time.FixedZone("ICT", 7*60*60)
	return fmt.Sprintf("%s_%d_%d", now.In(location).Format("060102"), orderID, now.UnixMilli()%1_000_000_000)
}

func generateInternalTransactionID() (string, error) {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return "pay_" + hex.EncodeToString(randomBytes), nil
}

func appendQueryParams(rawURL string, params map[string]string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	query := parsed.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func paymentItemName(item model.OrderItem) string {
	if item.Product != nil && item.Product.Name != "" {
		return item.Product.Name
	}
	return fmt.Sprintf("Product %d", item.ProductID)
}

func nonZeroInt64String(value int64) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%d", value)
}

func (r ZaloPayCreateOrderResponse) SubReturnMessageOrFallback() string {
	if strings.TrimSpace(r.SubReturnMessage) != "" {
		return r.SubReturnMessage
	}
	if strings.TrimSpace(r.ReturnMessage) != "" {
		return r.ReturnMessage
	}
	return "ZaloPay rejected the payment request"
}
