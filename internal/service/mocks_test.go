package service

import (
	"context"
	"time"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockUserRepository struct {
	createFunc                 func(ctx context.Context, user *model.User) error
	findByEmailFunc            func(ctx context.Context, email string) (*model.User, error)
	findByIDFunc               func(ctx context.Context, id int64) (*model.User, error)
	listFunc                   func(ctx context.Context, query dto.UserListQuery) ([]model.User, int64, error)
	existsByEmailOtherUserFunc func(ctx context.Context, email string, userID int64) (bool, error)
	updateFunc                 func(ctx context.Context, user *model.User) error
	softDeleteFunc             func(ctx context.Context, id int64) error
	saveRefreshTokenFunc       func(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	findRefreshTokenByHashFunc func(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	revokeRefreshTokenFunc     func(ctx context.Context, tokenHash string) error
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.findByEmailFunc(ctx, email)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockUserRepository) List(ctx context.Context, query dto.UserListQuery) ([]model.User, int64, error) {
	return m.listFunc(ctx, query)
}

func (m *mockUserRepository) ExistsByEmailOtherUser(ctx context.Context, email string, userID int64) (bool, error) {
	return m.existsByEmailOtherUserFunc(ctx, email, userID)
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	return m.updateFunc(ctx, user)
}

func (m *mockUserRepository) SoftDelete(ctx context.Context, id int64) error {
	return m.softDeleteFunc(ctx, id)
}

func (m *mockUserRepository) SaveRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	return m.saveRefreshTokenFunc(ctx, userID, tokenHash, expiresAt)
}

func (m *mockUserRepository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	return m.findRefreshTokenByHashFunc(ctx, tokenHash)
}

func (m *mockUserRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return m.revokeRefreshTokenFunc(ctx, tokenHash)
}

type mockProductRepository struct {
	createFunc         func(ctx context.Context, product *model.Product) error
	findByIDFunc       func(ctx context.Context, id int64) (*model.Product, error)
	findActiveByIDFunc func(ctx context.Context, id int64) (*model.Product, error)
	listFunc           func(ctx context.Context, query dto.ProductListQuery) ([]model.Product, int64, error)
	updateFunc         func(ctx context.Context, product *model.Product) error
	softDeleteFunc     func(ctx context.Context, id int64) error
}

func (m *mockProductRepository) Create(ctx context.Context, product *model.Product) error {
	return m.createFunc(ctx, product)
}

func (m *mockProductRepository) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockProductRepository) FindActiveByID(ctx context.Context, id int64) (*model.Product, error) {
	return m.findActiveByIDFunc(ctx, id)
}

func (m *mockProductRepository) List(ctx context.Context, query dto.ProductListQuery) ([]model.Product, int64, error) {
	return m.listFunc(ctx, query)
}

func (m *mockProductRepository) Update(ctx context.Context, product *model.Product) error {
	return m.updateFunc(ctx, product)
}

func (m *mockProductRepository) SoftDelete(ctx context.Context, id int64) error {
	return m.softDeleteFunc(ctx, id)
}

type mockCategoryRepository struct {
	createFunc                    func(ctx context.Context, category *model.Category) error
	findByIDFunc                  func(ctx context.Context, id int64) (*model.Category, error)
	findActiveByIDFunc            func(ctx context.Context, id int64) (*model.Category, error)
	listActiveFunc                func(ctx context.Context) ([]model.Category, error)
	existsByNameOtherCategoryFunc func(ctx context.Context, name string, categoryID int64) (bool, error)
	updateFunc                    func(ctx context.Context, category *model.Category) error
	softDeleteFunc                func(ctx context.Context, id int64) error
}

func (m *mockCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	return m.createFunc(ctx, category)
}

func (m *mockCategoryRepository) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockCategoryRepository) FindActiveByID(ctx context.Context, id int64) (*model.Category, error) {
	return m.findActiveByIDFunc(ctx, id)
}

func (m *mockCategoryRepository) ListActive(ctx context.Context) ([]model.Category, error) {
	return m.listActiveFunc(ctx)
}

func (m *mockCategoryRepository) ExistsByNameOtherCategory(ctx context.Context, name string, categoryID int64) (bool, error) {
	return m.existsByNameOtherCategoryFunc(ctx, name, categoryID)
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	return m.updateFunc(ctx, category)
}

func (m *mockCategoryRepository) SoftDelete(ctx context.Context, id int64) error {
	return m.softDeleteFunc(ctx, id)
}

type mockOrderRepository struct {
	products map[int64]*model.Product
}

func (m *mockOrderRepository) CreateOrder(ctx context.Context, tx repository.Tx, order *model.Order) error {
	order.ID = 1
	return nil
}

func (m *mockOrderRepository) CreateOrderItem(ctx context.Context, tx repository.Tx, item *model.OrderItem) error {
	item.ID = item.ProductID
	return nil
}

func (m *mockOrderRepository) FindProductForUpdate(ctx context.Context, tx repository.Tx, productID int64) (*model.Product, error) {
	product := m.products[productID]
	if product == nil {
		return nil, nil
	}
	copyProduct := *product
	return &copyProduct, nil
}

func (m *mockOrderRepository) DecreaseStock(ctx context.Context, tx repository.Tx, productID int64, quantity int) error {
	product := m.products[productID]
	if product == nil || product.Stock < quantity {
		return pgx.ErrNoRows
	}
	product.Stock -= quantity
	return nil
}

func (m *mockOrderRepository) ListByUserID(ctx context.Context, db repository.Queryer, userID int64) ([]model.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) ListAll(ctx context.Context, db repository.Queryer) ([]model.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) FindByID(ctx context.Context, db repository.Queryer, orderID int64) (*model.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, db repository.Queryer, orderID int64, status string) error {
	return nil
}

func (m *mockOrderRepository) FindItemsByOrderID(ctx context.Context, db repository.Queryer, orderID int64) ([]model.OrderItem, error) {
	return nil, nil
}

type mockTxBeginner struct {
	tx *mockTx
}

func (m *mockTxBeginner) Begin(ctx context.Context) (repository.Tx, error) {
	return m.tx, nil
}

type mockTx struct {
	committed  bool
	rolledBack bool
}

func (m *mockTx) Commit(ctx context.Context) error {
	m.committed = true
	return nil
}

func (m *mockTx) Rollback(ctx context.Context) error {
	m.rolledBack = true
	return nil
}

func (m *mockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (m *mockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

type mockQueryer struct{}

func (m *mockQueryer) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockQueryer) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (m *mockQueryer) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
