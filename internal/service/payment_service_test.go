package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"enterprise-order-management-api/internal/config"
	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/repository"

	"github.com/stretchr/testify/require"
)

func TestPaymentServiceCreateZaloPayPaymentSuccess(t *testing.T) {
	now := time.Date(2026, time.June, 25, 10, 0, 0, 0, time.UTC)
	orderRepo := &mockOrderRepository{
		findByIDFunc: func(context.Context, repository.Queryer, int64) (*model.Order, error) {
			return &model.Order{
				ID:          11,
				UserID:      99,
				Status:      model.OrderStatusPending,
				TotalAmount: 250000,
				Items: []model.OrderItem{
					{ProductID: 5, Quantity: 2, UnitPrice: 125000, Product: &model.Product{Name: "Office 365"}},
				},
			}, nil
		},
	}

	paymentRepo := &mockPaymentRepository{}
	service := &paymentService{
		db:       &mockQueryer{},
		cfg:      testPaymentConfig(),
		orders:   orderRepo,
		payments: paymentRepo,
		zaloPay: &mockZaloPayClient{
			createOrderFunc: func(ctx context.Context, req ZaloPayCreateOrderRequest) (*ZaloPayCreateOrderResponse, string, error) {
				require.Equal(t, "user-99", req.AppUser)
				require.Equal(t, int64(250000), req.Amount)
				require.Contains(t, req.RedirectURL, "transactionId=")
				return &ZaloPayCreateOrderResponse{
					ReturnCode: 1,
					OrderURL:   "https://sandbox.zalopay.vn/pay/demo",
				}, `{"return_code":1,"order_url":"https://sandbox.zalopay.vn/pay/demo"}`, nil
			},
		},
		timeNowFn: func() time.Time { return now },
	}

	res, err := service.CreateZaloPayPayment(context.Background(), 99, dto.CreateZaloPayPaymentRequest{
		OrderID: 11,
		Method:  model.PaymentProviderZaloPay,
	})

	require.NoError(t, err)
	require.Equal(t, model.PaymentProviderZaloPay, res.Provider)
	require.Equal(t, model.PaymentStatusPending, res.Status)
	require.Equal(t, "https://sandbox.zalopay.vn/pay/demo", res.PaymentURL)
	require.NotEmpty(t, res.TransactionID)
}

func TestPaymentServiceCreateZaloPayPaymentRejectsOtherUsersOrder(t *testing.T) {
	service := &paymentService{
		db:  &mockQueryer{},
		cfg: testPaymentConfig(),
		orders: &mockOrderRepository{
			findByIDFunc: func(context.Context, repository.Queryer, int64) (*model.Order, error) {
				return &model.Order{ID: 11, UserID: 100, Status: model.OrderStatusPending, TotalAmount: 1000}, nil
			},
		},
		payments:  &mockPaymentRepository{},
		zaloPay:   &mockZaloPayClient{},
		timeNowFn: time.Now,
	}

	res, err := service.CreateZaloPayPayment(context.Background(), 99, dto.CreateZaloPayPaymentRequest{
		OrderID: 11,
		Method:  model.PaymentProviderZaloPay,
	})

	require.Error(t, err)
	require.Nil(t, res)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, "FORBIDDEN", appErr.Code)
}

func TestPaymentServiceHandleZaloPayCallbackMarksPaidIdempotently(t *testing.T) {
	now := time.Now()
	payload := zaloPayCallbackPayload{
		AppTransID: "260625_11_123",
		ZPTransID:  99887766,
		ServerTime: now.UnixMilli(),
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)

	orderUpdated := false
	paymentUpdated := false
	service := &paymentService{
		db:  &mockQueryer{},
		cfg: testPaymentConfig(),
		orders: &mockOrderRepository{
			findByIDFunc: func(context.Context, repository.Queryer, int64) (*model.Order, error) {
				return &model.Order{ID: 11, UserID: 99, Status: model.OrderStatusPending}, nil
			},
			updateStatusFunc: func(context.Context, repository.Queryer, int64, string) error {
				orderUpdated = true
				return nil
			},
		},
		payments: &mockPaymentRepository{
			findByAppTransactionIDFunc: func(context.Context, repository.Queryer, string) (*model.Payment, error) {
				return &model.Payment{
					ID:               1,
					OrderID:          11,
					UserID:           99,
					AppTransactionID: payload.AppTransID,
					Status:           model.PaymentStatusPending,
					Provider:         model.PaymentProviderZaloPay,
				}, nil
			},
			updateStatusFunc: func(context.Context, repository.Queryer, *model.Payment) error {
				paymentUpdated = true
				return nil
			},
		},
		zaloPay: &mockZaloPayClient{
			verifyCallbackFunc: func(data string, mac string) bool {
				return mac == "valid"
			},
		},
		timeNowFn: time.Now,
	}

	res, err := service.HandleZaloPayCallback(context.Background(), dto.ZaloPayCallbackRequest{
		Data: string(data),
		MAC:  "valid",
		Type: 1,
	})

	require.NoError(t, err)
	require.Equal(t, 1, res.ReturnCode)
	require.True(t, paymentUpdated)
	require.True(t, orderUpdated)
}

func TestPaymentServiceGetZaloPayStatusQueriesGateway(t *testing.T) {
	updatedPayment := false
	service := &paymentService{
		db:  &mockQueryer{},
		cfg: testPaymentConfig(),
		orders: &mockOrderRepository{
			findByIDFunc: func(context.Context, repository.Queryer, int64) (*model.Order, error) {
				return &model.Order{ID: 11, UserID: 99, Status: model.OrderStatusPending}, nil
			},
			updateStatusFunc: func(context.Context, repository.Queryer, int64, string) error {
				return nil
			},
		},
		payments: &mockPaymentRepository{
			findByTransactionIDFunc: func(context.Context, repository.Queryer, string) (*model.Payment, error) {
				expiresAt := time.Now().Add(10 * time.Minute)
				return &model.Payment{
					ID:               1,
					TransactionID:    "pay_123",
					OrderID:          11,
					UserID:           99,
					Provider:         model.PaymentProviderZaloPay,
					AppTransactionID: "260625_11_123",
					Amount:           250000,
					Currency:         "VND",
					Status:           model.PaymentStatusPending,
					ExpiredAt:        &expiresAt,
				}, nil
			},
			updateStatusFunc: func(context.Context, repository.Queryer, *model.Payment) error {
				updatedPayment = true
				return nil
			},
		},
		zaloPay: &mockZaloPayClient{
			queryOrderFunc: func(context.Context, string) (*ZaloPayQueryOrderResponse, string, error) {
				return &ZaloPayQueryOrderResponse{
					ReturnCode:   1,
					IsProcessing: false,
					ZPTransID:    123456789,
					ServerTime:   time.Now().UnixMilli(),
				}, `{"return_code":1}`, nil
			},
		},
		timeNowFn: time.Now,
	}

	res, err := service.GetZaloPayStatus(context.Background(), 99, "pay_123")

	require.NoError(t, err)
	require.Equal(t, model.PaymentStatusPaid, res.Status)
	require.True(t, updatedPayment)
}

func testPaymentConfig() config.Config {
	return config.Config{
		ZaloPay: config.ZaloPayConfig{
			Enabled:             true,
			Environment:         "sandbox",
			AppID:               2553,
			Key1:                "key1",
			Key2:                "key2",
			CreateOrderEndpoint: "https://sandbox.zalopay.vn/create",
			QueryOrderEndpoint:  "https://sandbox.zalopay.vn/query",
			CallbackURL:         "http://localhost:8080/api/v1/payments/zalopay/callback",
			RedirectURL:         "http://localhost:5173/payment/zalopay/return",
			Currency:            "VND",
			Timeout:             30 * time.Second,
		},
	}
}
