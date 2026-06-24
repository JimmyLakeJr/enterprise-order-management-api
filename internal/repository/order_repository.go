package repository

import (
	"context"
	"fmt"
	"strings"

	"enterprise-order-management-api/internal/dto"
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

type Tx interface {
	Queryer
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxBeginner interface {
	Begin(ctx context.Context) (Tx, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx Tx, order *model.Order) error
	CreateOrderItem(ctx context.Context, tx Tx, item *model.OrderItem) error
	FindProductForUpdate(ctx context.Context, tx Tx, productID int64) (*model.Product, error)
	DecreaseStock(ctx context.Context, tx Tx, productID int64, quantity int) error
	ListByUserID(ctx context.Context, db Queryer, userID int64, query dto.OrderListQuery) ([]model.Order, int64, error)
	ListAll(ctx context.Context, db Queryer, query dto.OrderListQuery) ([]model.Order, int64, error)
	FindByID(ctx context.Context, db Queryer, orderID int64) (*model.Order, error)
	UpdateStatus(ctx context.Context, db Queryer, orderID int64, status string) error
	FindItemsByOrderID(ctx context.Context, db Queryer, orderID int64) ([]model.OrderItem, error)
	FindItemsByOrderIDs(ctx context.Context, db Queryer, orderIDs []int64) (map[int64][]model.OrderItem, error)
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, tx Tx, order *model.Order) error {
	query := `
		INSERT INTO orders (user_id, status, total_amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return tx.QueryRow(ctx, query, order.UserID, order.Status, order.TotalAmount).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *orderRepository) CreateOrderItem(ctx context.Context, tx Tx, item *model.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return tx.QueryRow(ctx, query, item.OrderID, item.ProductID, item.Quantity, item.UnitPrice, item.Subtotal).
		Scan(&item.ID)
}

func (r *orderRepository) FindProductForUpdate(ctx context.Context, tx Tx, productID int64) (*model.Product, error) {
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

func (r *orderRepository) DecreaseStock(ctx context.Context, tx Tx, productID int64, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock - $1, updated_at = NOW()
		WHERE id = $2 AND stock >= $1
	`
	commandTag, err := tx.Exec(ctx, query, quantity, productID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *orderRepository) ListByUserID(ctx context.Context, db Queryer, userID int64, query dto.OrderListQuery) ([]model.Order, int64, error) {
	baseWhere, args := buildOrderWhere(query)
	userArgPosition := len(args) + 1

	if baseWhere == "" {
		baseWhere = fmt.Sprintf("WHERE o.user_id = $%d", userArgPosition)
	} else {
		baseWhere += fmt.Sprintf(" AND o.user_id = $%d", userArgPosition)
	}

	args = append(args, userID)
	return r.list(ctx, db, baseWhere, args, query)
}

func (r *orderRepository) ListAll(ctx context.Context, db Queryer, query dto.OrderListQuery) ([]model.Order, int64, error) {
	where, args := buildOrderWhere(query)
	return r.list(ctx, db, where, args, query)
}

func (r *orderRepository) FindByID(ctx context.Context, db Queryer, orderID int64) (*model.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.status, o.total_amount, o.created_at, o.updated_at,
		       u.id, u.full_name, u.email
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE o.id = $1
	`
	order := &model.Order{User: &model.User{}}
	err := db.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.User.ID,
		&order.User.Name,
		&order.User.Email,
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
	itemsByOrderID, err := r.FindItemsByOrderIDs(ctx, db, []int64{orderID})
	if err != nil {
		return nil, err
	}
	return itemsByOrderID[orderID], nil
}

func (r *orderRepository) FindItemsByOrderIDs(ctx context.Context, db Queryer, orderIDs []int64) (map[int64][]model.OrderItem, error) {
	itemsByOrderID := make(map[int64][]model.OrderItem, len(orderIDs))
	if len(orderIDs) == 0 {
		return itemsByOrderID, nil
	}

	placeholders := make([]string, 0, len(orderIDs))
	args := make([]any, 0, len(orderIDs))
	for i, orderID := range orderIDs {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, orderID)
	}

	query := fmt.Sprintf(`
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.subtotal,
		       p.id, p.category_id, p.name, p.description, p.price, p.stock, p.image_url, p.is_active, p.created_at, p.updated_at
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id IN (%s)
		ORDER BY oi.order_id ASC, oi.id ASC
	`, strings.Join(placeholders, ", "))

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}

	return itemsByOrderID, rows.Err()
}

func (r *orderRepository) list(ctx context.Context, db Queryer, where string, args []any, query dto.OrderListQuery) ([]model.Order, int64, error) {
	countSQL := "SELECT COUNT(*) FROM orders o JOIN users u ON u.id = o.user_id " + where
	var total int64
	if err := db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Limit
	listArgs := append(append([]any{}, args...), query.Limit, offset)
	listSQL := fmt.Sprintf(`
		SELECT o.id, o.user_id, o.status, o.total_amount, o.created_at, o.updated_at,
		       u.id, u.full_name, u.email
		FROM orders o
		JOIN users u ON u.id = o.user_id
		%s
		ORDER BY o.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(listArgs)-1, len(listArgs))

	rows, err := db.Query(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		order := model.Order{User: &model.User{}}
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.User.ID,
			&order.User.Name,
			&order.User.Email,
		); err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}
	return orders, total, rows.Err()
}

func buildOrderWhere(query dto.OrderListQuery) (string, []any) {
	conditions := make([]string, 0, 1)
	args := make([]any, 0, 1)

	if query.Status != "" {
		args = append(args, query.Status)
		conditions = append(conditions, fmt.Sprintf("o.status = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}
