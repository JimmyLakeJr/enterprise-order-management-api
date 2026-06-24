package repository

import (
	"context"

	"enterprise-order-management-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	FindByID(ctx context.Context, id int64) (*model.Category, error)
	FindActiveByID(ctx context.Context, id int64) (*model.Category, error)
	ListActive(ctx context.Context) ([]model.Category, error)
	ListAdmin(ctx context.Context, status string) ([]model.Category, error)
	ExistsByNameOtherCategory(ctx context.Context, name string, categoryID int64) (bool, error)
	HasActiveProducts(ctx context.Context, categoryID int64) (bool, error)
	Update(ctx context.Context, category *model.Category) error
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO categories (name, description, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, category.Name, category.Description, category.IsActive).
		Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

func (r *categoryRepository) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`
	category := &model.Category{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return category, err
}

func (r *categoryRepository) FindActiveByID(ctx context.Context, id int64) (*model.Category, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM categories
		WHERE id = $1 AND is_active = TRUE
	`
	return r.findOne(ctx, query, id)
}

func (r *categoryRepository) ListActive(ctx context.Context) ([]model.Category, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM categories
		WHERE is_active = TRUE
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *categoryRepository) ListAdmin(ctx context.Context, status string) ([]model.Category, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM categories
	`
	switch status {
	case "active":
		query += " WHERE is_active = TRUE"
	case "inactive":
		query += " WHERE is_active = FALSE"
	}
	query += " ORDER BY name ASC"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.IsActive, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *categoryRepository) ExistsByNameOtherCategory(ctx context.Context, name string, categoryID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM categories
			WHERE LOWER(name) = LOWER($1) AND id <> $2
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, name, categoryID).Scan(&exists)
	return exists, err
}

func (r *categoryRepository) HasActiveProducts(ctx context.Context, categoryID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE category_id = $1 AND is_active = TRUE)`

	var exists bool
	err := r.db.QueryRow(ctx, query, categoryID).Scan(&exists)
	return exists, err
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	query := `
		UPDATE categories
		SET name = $1, description = $2, is_active = $3, updated_at = NOW()
		WHERE id = $4 AND is_active = TRUE
	`
	commandTag, err := r.db.Exec(ctx, query, category.Name, category.Description, category.IsActive, category.ID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *categoryRepository) SoftDelete(ctx context.Context, id int64) error {
	query := `UPDATE categories SET is_active = FALSE, updated_at = NOW() WHERE id = $1 AND is_active = TRUE`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *categoryRepository) Restore(ctx context.Context, id int64) error {
	query := `UPDATE categories SET is_active = TRUE, updated_at = NOW() WHERE id = $1 AND is_active = FALSE`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *categoryRepository) findOne(ctx context.Context, query string, args ...any) (*model.Category, error) {
	category := &model.Category{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return category, err
}
