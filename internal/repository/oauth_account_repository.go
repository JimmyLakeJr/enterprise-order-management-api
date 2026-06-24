package repository

import (
	"context"

	"enterprise-order-management-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OAuthAccountRepository interface {
	CreateWithQuerier(ctx context.Context, q Queryer, account *model.OAuthAccount) error
	FindByProviderUserID(ctx context.Context, provider string, providerUserID string) (*model.OAuthAccount, error)
}

type oauthAccountRepository struct {
	db *pgxpool.Pool
}

func NewOAuthAccountRepository(db *pgxpool.Pool) OAuthAccountRepository {
	return &oauthAccountRepository{db: db}
}

func (r *oauthAccountRepository) CreateWithQuerier(ctx context.Context, q Queryer, account *model.OAuthAccount) error {
	query := `
		INSERT INTO oauth_accounts (user_id, provider, provider_user_id, email, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return q.QueryRow(ctx, query, account.UserID, account.Provider, account.ProviderUserID, account.Email, account.AvatarURL).
		Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
}

func (r *oauthAccountRepository) FindByProviderUserID(ctx context.Context, provider string, providerUserID string) (*model.OAuthAccount, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, email, avatar_url, created_at, updated_at
		FROM oauth_accounts
		WHERE provider = $1 AND provider_user_id = $2
	`

	account := &model.OAuthAccount{}
	err := r.db.QueryRow(ctx, query, provider, providerUserID).Scan(
		&account.ID,
		&account.UserID,
		&account.Provider,
		&account.ProviderUserID,
		&account.Email,
		&account.AvatarURL,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return account, err
}
