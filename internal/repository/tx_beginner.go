package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type poolTxBeginner struct {
	pool *pgxpool.Pool
}

func NewTxBeginner(pool *pgxpool.Pool) TxBeginner {
	return &poolTxBeginner{pool: pool}
}

func (b *poolTxBeginner) Begin(ctx context.Context) (Tx, error) {
	return b.pool.Begin(ctx)
}
