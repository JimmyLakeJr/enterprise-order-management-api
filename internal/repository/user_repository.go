package repository

import (
	"context"
	"time"

	"enterprise-order-management-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	SoftDelete(ctx context.Context, id int64) error
	SaveRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	FindActiveRefreshToken(ctx context.Context, tokenHash string) (int64, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))
		RETURNING id, role_id, is_active, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash, model.RoleUser).
		Scan(&user.ID, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.password_hash, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1 AND u.is_active = TRUE
	`

	return r.findOne(ctx, query, email)
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.password_hash, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1 AND u.is_active = TRUE
	`

	return r.findOne(ctx, query, id)
}

func (r *userRepository) List(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.password_hash, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := scanUser(rows, &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *userRepository) SoftDelete(ctx context.Context, id int64) error {
	query := `UPDATE users SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) SaveRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *userRepository) FindActiveRefreshToken(ctx context.Context, tokenHash string) (int64, error) {
	query := `
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()
	`

	var userID int64
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&userID)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return userID, err
}

func (r *userRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}

func (r *userRepository) findOne(ctx context.Context, query string, args ...any) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func scanUser(rows pgx.Rows, user *model.User) error {
	return rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
