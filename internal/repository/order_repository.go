package repository

import (
	"context"

	"enterprise-order-management-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx pgx.Tx, order *model.Order) error
	CreateOrderItem(ctx context.Context, tx pgx.Tx, item *model.OrderItem) error
	FindProductForUpdate(ctx context.Context, tx pgx.Tx, productID int64) (*model.Product, error)
	DecreaseStock(ctx context.Context, tx pgx.Tx, productID int64, quantity int) error
	ListByUserID(ctx context.Context, db Queryer, userID int64) ([]model.Order, error)
	ListAll(ctx context.Context, db Queryer) ([]model.Order, error)
	FindByID(ctx context.Context, db Queryer, orderID int64) (*model.Order, error)
	UpdateStatus(ctx context.Context, db Queryer, orderID int64, status string) error
	FindItemsByOrderID(ctx context.Context, db Queryer, orderID int64) ([]model.OrderItem, error)
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, tx pgx.Tx, order *model.Order) error {
	query := `
		INSERT INTO orders (user_id, status, total_amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return tx.QueryRow(ctx, query, order.UserID, order.Status, order.TotalAmount).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *orderRepository) CreateOrderItem(ctx context.Context, tx pgx.Tx, item *model.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return tx.QueryRow(ctx, query, item.OrderID, item.ProductID, item.Quantity, item.UnitPrice, item.Subtotal).
		Scan(&item.ID)
}

func (r *orderRepository) FindProductForUpdate(ctx context.Context, tx pgx.Tx, productID int64) (*model.Product, error) {
	query := `
		SELECT id, category_id, name, description, price, stock, image_url, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
		FOR UPDATE
	`
	product := &model.Product{}
	err := tx.QueryRow(ctx, query, productID).Scan(
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
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return product, err
}

func (r *orderRepository) DecreaseStock(ctx context.Context, tx pgx.Tx, productID int64, quantity int) error {
	query := `UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`
	_, err := tx.Exec(ctx, query, quantity, productID)
	return err
}

func (r *orderRepository) ListByUserID(ctx context.Context, db Queryer, userID int64) ([]model.Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	return r.list(ctx, db, query, userID)
}

func (r *orderRepository) ListAll(ctx context.Context, db Queryer) ([]model.Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`
	return r.list(ctx, db, query)
}

func (r *orderRepository) FindByID(ctx context.Context, db Queryer, orderID int64) (*model.Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	order := &model.Order{}
	err := db.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return order, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, db Queryer, orderID int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	commandTag, err := db.Exec(ctx, query, status, orderID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *orderRepository) FindItemsByOrderID(ctx context.Context, db Queryer, orderID int64) ([]model.OrderItem, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.subtotal,
		       p.id, p.category_id, p.name, p.description, p.price, p.stock, p.image_url, p.is_active, p.created_at, p.updated_at
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1
		ORDER BY oi.id ASC
	`
	rows, err := db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.OrderItem, 0)
	for rows.Next() {
		var item model.OrderItem
		product := &model.Product{}
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
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
		); err != nil {
			return nil, err
		}
		item.Product = product
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *orderRepository) list(ctx context.Context, db Queryer, query string, args ...any) ([]model.Order, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}
