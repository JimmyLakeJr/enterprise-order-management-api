package repository

import (
	"context"
	"database/sql"

	"enterprise-order-management-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	Create(ctx context.Context, db Queryer, payment *model.Payment) error
	FindLatestByOrderIDAndProvider(ctx context.Context, db Queryer, orderID int64, provider string) (*model.Payment, error)
	FindByTransactionID(ctx context.Context, db Queryer, transactionID string) (*model.Payment, error)
	FindByAppTransactionID(ctx context.Context, db Queryer, appTransactionID string) (*model.Payment, error)
	UpdateGatewayInitialization(ctx context.Context, db Queryer, payment *model.Payment) error
	UpdateStatus(ctx context.Context, db Queryer, payment *model.Payment) error
	UpdateRawCallback(ctx context.Context, db Queryer, paymentID int64, rawCallback string) error
}

type paymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, db Queryer, payment *model.Payment) error {
	query := `
		INSERT INTO payments (
			transaction_id, order_id, user_id, provider, provider_transaction_id, app_transaction_id,
			amount, currency, status, payment_url, raw_request, raw_response, raw_callback, paid_at, expired_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`

	return db.QueryRow(
		ctx,
		query,
		payment.TransactionID,
		payment.OrderID,
		payment.UserID,
		payment.Provider,
		nullableString(payment.ProviderTransactionID),
		nullableString(payment.AppTransactionID),
		payment.Amount,
		payment.Currency,
		payment.Status,
		nullableString(payment.PaymentURL),
		nullableString(payment.RawRequest),
		nullableString(payment.RawResponse),
		nullableString(payment.RawCallback),
		payment.PaidAt,
		payment.ExpiredAt,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
}

func (r *paymentRepository) FindLatestByOrderIDAndProvider(ctx context.Context, db Queryer, orderID int64, provider string) (*model.Payment, error) {
	query := `
		SELECT id, transaction_id, order_id, user_id, provider, provider_transaction_id, app_transaction_id,
		       amount, currency, status, payment_url, raw_request, raw_response, raw_callback,
		       paid_at, expired_at, created_at, updated_at
		FROM payments
		WHERE order_id = $1 AND provider = $2
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`
	return scanPaymentRow(db.QueryRow(ctx, query, orderID, provider))
}

func (r *paymentRepository) FindByTransactionID(ctx context.Context, db Queryer, transactionID string) (*model.Payment, error) {
	query := `
		SELECT id, transaction_id, order_id, user_id, provider, provider_transaction_id, app_transaction_id,
		       amount, currency, status, payment_url, raw_request, raw_response, raw_callback,
		       paid_at, expired_at, created_at, updated_at
		FROM payments
		WHERE transaction_id = $1
	`
	return scanPaymentRow(db.QueryRow(ctx, query, transactionID))
}

func (r *paymentRepository) FindByAppTransactionID(ctx context.Context, db Queryer, appTransactionID string) (*model.Payment, error) {
	query := `
		SELECT id, transaction_id, order_id, user_id, provider, provider_transaction_id, app_transaction_id,
		       amount, currency, status, payment_url, raw_request, raw_response, raw_callback,
		       paid_at, expired_at, created_at, updated_at
		FROM payments
		WHERE app_transaction_id = $1
	`
	return scanPaymentRow(db.QueryRow(ctx, query, appTransactionID))
}

func (r *paymentRepository) UpdateGatewayInitialization(ctx context.Context, db Queryer, payment *model.Payment) error {
	query := `
		UPDATE payments
		SET provider_transaction_id = $1,
		    app_transaction_id = $2,
		    payment_url = $3,
		    raw_request = $4,
		    raw_response = $5,
		    status = $6,
		    expired_at = $7,
		    updated_at = NOW()
		WHERE id = $8
	`
	_, err := db.Exec(
		ctx,
		query,
		nullableString(payment.ProviderTransactionID),
		nullableString(payment.AppTransactionID),
		nullableString(payment.PaymentURL),
		nullableString(payment.RawRequest),
		nullableString(payment.RawResponse),
		payment.Status,
		payment.ExpiredAt,
		payment.ID,
	)
	return err
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, db Queryer, payment *model.Payment) error {
	query := `
		UPDATE payments
		SET provider_transaction_id = $1,
		    payment_url = $2,
		    raw_response = $3,
		    raw_callback = $4,
		    status = $5,
		    paid_at = $6,
		    expired_at = $7,
		    updated_at = NOW()
		WHERE id = $8
	`
	_, err := db.Exec(
		ctx,
		query,
		nullableString(payment.ProviderTransactionID),
		nullableString(payment.PaymentURL),
		nullableString(payment.RawResponse),
		nullableString(payment.RawCallback),
		payment.Status,
		payment.PaidAt,
		payment.ExpiredAt,
		payment.ID,
	)
	return err
}

func (r *paymentRepository) UpdateRawCallback(ctx context.Context, db Queryer, paymentID int64, rawCallback string) error {
	query := `UPDATE payments SET raw_callback = $1, updated_at = NOW() WHERE id = $2`
	_, err := db.Exec(ctx, query, rawCallback, paymentID)
	return err
}

func scanPaymentRow(row pgx.Row) (*model.Payment, error) {
	var payment model.Payment
	var providerTransactionID sql.NullString
	var appTransactionID sql.NullString
	var paymentURL sql.NullString
	var rawRequest sql.NullString
	var rawResponse sql.NullString
	var rawCallback sql.NullString
	err := row.Scan(
		&payment.ID,
		&payment.TransactionID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Provider,
		&providerTransactionID,
		&appTransactionID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&paymentURL,
		&rawRequest,
		&rawResponse,
		&rawCallback,
		&payment.PaidAt,
		&payment.ExpiredAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	payment.ProviderTransactionID = providerTransactionID.String
	payment.AppTransactionID = appTransactionID.String
	payment.PaymentURL = paymentURL.String
	payment.RawRequest = rawRequest.String
	payment.RawResponse = rawResponse.String
	payment.RawCallback = rawCallback.String
	return &payment, err
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
