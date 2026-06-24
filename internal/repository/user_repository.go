package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	CreateWithQuerier(ctx context.Context, q Queryer, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByEmailAny(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByIDAny(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context, query dto.UserListQuery) ([]model.User, int64, error)
	ExistsByEmailOtherUser(ctx context.Context, email string, userID int64) (bool, error)
	Update(ctx context.Context, user *model.User) error
	UpdateProfileName(ctx context.Context, id int64, name string) error
	UpdateAvatarURL(ctx context.Context, id int64, avatarURL string) error
	UpdateProfileVideoURL(ctx context.Context, id int64, profileVideoURL string) error
	SoftDelete(ctx context.Context, id int64) error
	SaveRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.CreateWithQuerier(ctx, r.db, user)
}

func (r *userRepository) CreateWithQuerier(ctx context.Context, q Queryer, user *model.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, avatar_url, profile_video_url, role_id)
		VALUES ($1, $2, $3, $4, $5, (SELECT id FROM roles WHERE name = $6))
		RETURNING id, role_id, is_active, created_at, updated_at
	`

	return q.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash, user.AvatarURL, user.ProfileVideoURL, user.Role).
		Scan(&user.ID, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.password_hash, u.avatar_url, u.profile_video_url, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1 AND u.is_active = TRUE
	`

	return r.findOne(ctx, query, email)
}

func (r *userRepository) FindByEmailAny(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.password_hash, u.avatar_url, u.profile_video_url, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1
	`

	return r.findOne(ctx, query, email)
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.password_hash, u.avatar_url, u.profile_video_url, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1 AND u.is_active = TRUE
	`

	return r.findOne(ctx, query, id)
}

func (r *userRepository) FindByIDAny(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.password_hash, u.avatar_url, u.profile_video_url, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`

	return r.findOne(ctx, query, id)
}

func (r *userRepository) List(ctx context.Context, query dto.UserListQuery) ([]model.User, int64, error) {
	where, args := buildUserWhere(query)
	offset := (query.Page - 1) * query.Limit

	countSQL := "SELECT COUNT(*) FROM users u JOIN roles r ON r.id = u.role_id " + where
	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, query.Limit, offset)
	listSQL := fmt.Sprintf(`
		SELECT u.id, u.full_name, u.email, u.password_hash, u.avatar_url, u.profile_video_url, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		%s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		if err := scanUser(rows, &user); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	return users, total, rows.Err()
}

func (r *userRepository) ExistsByEmailOtherUser(ctx context.Context, email string, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id <> $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email, userID).Scan(&exists)
	return exists, err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET full_name = $1,
		    email = $2,
		    role_id = (SELECT id FROM roles WHERE name = $3),
		    updated_at = NOW()
		WHERE id = $4 AND is_active = TRUE
	`

	commandTag, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Role, user.ID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) UpdateProfileName(ctx context.Context, id int64, name string) error {
	query := `UPDATE users SET full_name = $1, updated_at = NOW() WHERE id = $2 AND is_active = TRUE`
	commandTag, err := r.db.Exec(ctx, query, name, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) UpdateAvatarURL(ctx context.Context, id int64, avatarURL string) error {
	query := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2 AND is_active = TRUE`
	commandTag, err := r.db.Exec(ctx, query, avatarURL, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) UpdateProfileVideoURL(ctx context.Context, id int64, profileVideoURL string) error {
	query := `UPDATE users SET profile_video_url = $1, updated_at = NOW() WHERE id = $2 AND is_active = TRUE`
	commandTag, err := r.db.Exec(ctx, query, profileVideoURL, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
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

func (r *userRepository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	refreshToken := &model.RefreshToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.RevokedAt,
		&refreshToken.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return refreshToken, err
}

func (r *userRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}

func buildUserWhere(query dto.UserListQuery) (string, []any) {
	conditions := []string{"u.is_active = TRUE"}
	args := make([]any, 0)

	if query.Search != "" {
		args = append(args, "%"+query.Search+"%")
		conditions = append(conditions, fmt.Sprintf("(u.email ILIKE $%d OR u.full_name ILIKE $%d)", len(args), len(args)))
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (r *userRepository) findOne(ctx context.Context, query string, args ...any) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.ProfileVideoURL,
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
		&user.AvatarURL,
		&user.ProfileVideoURL,
		&user.RoleID,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
