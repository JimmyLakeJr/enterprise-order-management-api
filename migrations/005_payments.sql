CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL UNIQUE,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(30) NOT NULL,
    provider_transaction_id VARCHAR(64),
    app_transaction_id VARCHAR(64),
    amount BIGINT NOT NULL CHECK (amount >= 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'VND',
    status VARCHAR(30) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'paid', 'failed', 'cancelled', 'expired')),
    payment_url TEXT,
    raw_request TEXT,
    raw_response TEXT,
    raw_callback TEXT,
    paid_at TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_provider ON payments(provider);
CREATE UNIQUE INDEX IF NOT EXISTS uq_payments_app_transaction_id ON payments(app_transaction_id) WHERE app_transaction_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_payments_provider_transaction_id ON payments(provider_transaction_id) WHERE provider_transaction_id IS NOT NULL;
