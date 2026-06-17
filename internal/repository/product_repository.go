package repository

import (
	"context"
	"fmt"
	"strings"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, id int64) (*model.Product, error)
	FindActiveByID(ctx context.Context, id int64) (*model.Product, error)
	List(ctx context.Context, query dto.ProductListQuery) ([]model.Product, int64, error)
	Update(ctx context.Context, product *model.Product) error
	SoftDelete(ctx context.Context, id int64) error
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (category_id, name, description, price, stock, image_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
		product.IsActive,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepository) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT p.id, p.category_id, p.name, p.description, p.price, p.stock, p.image_url, p.is_active, p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.is_active, c.created_at, c.updated_at
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`
	return r.findOne(ctx, query, id)
}

func (r *productRepository) FindActiveByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT p.id, p.category_id, p.name, p.description, p.price, p.stock, p.image_url, p.is_active, p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.is_active, c.created_at, c.updated_at
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1 AND p.is_active = TRUE AND c.is_active = TRUE
	`
	return r.findOne(ctx, query, id)
}

func (r *productRepository) findOne(ctx context.Context, query string, args ...any) (*model.Product, error) {
	product := &model.Product{Category: &model.Category{}}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Category.ID,
		&product.Category.Name,
		&product.Category.Description,
		&product.Category.IsActive,
		&product.Category.CreatedAt,
		&product.Category.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return product, err
}

func (r *productRepository) List(ctx context.Context, query dto.ProductListQuery) ([]model.Product, int64, error) {
	where, args := buildProductWhere(query)
	offset := (query.Page - 1) * query.Limit

	countSQL := "SELECT COUNT(*) FROM products p JOIN categories c ON c.id = p.category_id " + where
	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, query.Limit, offset)
	listSQL := fmt.Sprintf(`
		SELECT p.id, p.category_id, p.name, p.description, p.price, p.stock, p.image_url, p.is_active, p.created_at, p.updated_at,
		       c.id, c.name, c.description, c.is_active, c.created_at, c.updated_at
		FROM products p
		JOIN categories c ON c.id = p.category_id
		%s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		product := model.Product{Category: &model.Category{}}
		if err := scanProductWithCategory(rows, &product); err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, rows.Err()
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products
		SET category_id = $1, name = $2, description = $3, price = $4,
		    stock = $5, image_url = $6, is_active = $7, updated_at = NOW()
		WHERE id = $8 AND is_active = TRUE
	`
	commandTag, err := r.db.Exec(ctx, query,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
		product.IsActive,
		product.ID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *productRepository) SoftDelete(ctx context.Context, id int64) error {
	query := `UPDATE products SET is_active = FALSE, updated_at = NOW() WHERE id = $1 AND is_active = TRUE`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func buildProductWhere(query dto.ProductListQuery) (string, []any) {
	conditions := []string{"p.is_active = TRUE", "c.is_active = TRUE"}
	args := make([]any, 0)

	if query.Search != "" {
		args = append(args, "%"+query.Search+"%")
		conditions = append(conditions, fmt.Sprintf("p.name ILIKE $%d", len(args)))
	}

	if query.CategoryID > 0 {
		args = append(args, query.CategoryID)
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", len(args)))
	}

	if query.MinPrice > 0 {
		args = append(args, query.MinPrice)
		conditions = append(conditions, fmt.Sprintf("p.price >= $%d", len(args)))
	}

	if query.MaxPrice > 0 {
		args = append(args, query.MaxPrice)
		conditions = append(conditions, fmt.Sprintf("p.price <= $%d", len(args)))
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func scanProductWithCategory(rows pgx.Rows, product *model.Product) error {
	return rows.Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Category.ID,
		&product.Category.Name,
		&product.Category.Description,
		&product.Category.IsActive,
		&product.Category.CreatedAt,
		&product.Category.UpdatedAt,
	)
}
