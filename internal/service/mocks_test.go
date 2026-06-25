package service

import (
	"context"
	"time"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/model"
	"enterprise-order-management-api/internal/oauth"
	"enterprise-order-management-api/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockUserRepository struct {
	createFunc                 func(ctx context.Context, user *model.User) error
	createWithQuerierFunc      func(ctx context.Context, q repository.Queryer, user *model.User) error
	findByEmailFunc            func(ctx context.Context, email string) (*model.User, error)
	findByPhoneFunc            func(ctx context.Context, phone string) (*model.User, error)
	findByIdentifierFunc       func(ctx context.Context, identifier string) (*model.User, error)
	findByEmailAnyFunc         func(ctx context.Context, email string) (*model.User, error)
	findByPhoneAnyFunc         func(ctx context.Context, phone string) (*model.User, error)
	findByIDFunc               func(ctx context.Context, id int64) (*model.User, error)
	findByIDAnyFunc            func(ctx context.Context, id int64) (*model.User, error)
	listFunc                   func(ctx context.Context, query dto.UserListQuery) ([]model.User, int64, error)
	existsByEmailOtherUserFunc func(ctx context.Context, email string, userID int64) (bool, error)
	existsByPhoneOtherUserFunc func(ctx context.Context, phone string, userID int64) (bool, error)
	updateFunc                 func(ctx context.Context, user *model.User) error
	updateProfileFunc          func(ctx context.Context, id int64, name string, phone string) error
	updatePasswordHashFunc     func(ctx context.Context, id int64, passwordHash string) error
	updateAvatarURLFunc        func(ctx context.Context, id int64, avatarURL string) error
	updateProfileVideoURLFunc  func(ctx context.Context, id int64, profileVideoURL string) error
	softDeleteFunc             func(ctx context.Context, id int64) error
	saveRefreshTokenFunc       func(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	findRefreshTokenByHashFunc func(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	revokeRefreshTokenFunc     func(ctx context.Context, tokenHash string) error
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) CreateWithQuerier(ctx context.Context, q repository.Queryer, user *model.User) error {
	if m.createWithQuerierFunc != nil {
		return m.createWithQuerierFunc(ctx, q, user)
	}
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.findByEmailFunc(ctx, email)
}

func (m *mockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	if m.findByPhoneFunc != nil {
		return m.findByPhoneFunc(ctx, phone)
	}
	return nil, nil
}

func (m *mockUserRepository) FindByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	if m.findByIdentifierFunc != nil {
		return m.findByIdentifierFunc(ctx, identifier)
	}
	if m.findByEmailFunc != nil {
		return m.findByEmailFunc(ctx, identifier)
	}
	return nil, nil
}

func (m *mockUserRepository) FindByEmailAny(ctx context.Context, email string) (*model.User, error) {
	if m.findByEmailAnyFunc != nil {
		return m.findByEmailAnyFunc(ctx, email)
	}
	return m.findByEmailFunc(ctx, email)
}

func (m *mockUserRepository) FindByPhoneAny(ctx context.Context, phone string) (*model.User, error) {
	if m.findByPhoneAnyFunc != nil {
		return m.findByPhoneAnyFunc(ctx, phone)
	}
	return m.FindByPhone(ctx, phone)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	if m.findByIDFunc == nil {
		return nil, nil
	}
	return m.findByIDFunc(ctx, id)
}

func (m *mockUserRepository) FindByIDAny(ctx context.Context, id int64) (*model.User, error) {
	if m.findByIDAnyFunc != nil {
		return m.findByIDAnyFunc(ctx, id)
	}
	return m.findByIDFunc(ctx, id)
}

func (m *mockUserRepository) List(ctx context.Context, query dto.UserListQuery) ([]model.User, int64, error) {
	return m.listFunc(ctx, query)
}

func (m *mockUserRepository) ExistsByEmailOtherUser(ctx context.Context, email string, userID int64) (bool, error) {
	return m.existsByEmailOtherUserFunc(ctx, email, userID)
}

func (m *mockUserRepository) ExistsByPhoneOtherUser(ctx context.Context, phone string, userID int64) (bool, error) {
	if m.existsByPhoneOtherUserFunc != nil {
		return m.existsByPhoneOtherUserFunc(ctx, phone, userID)
	}
	return false, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	return m.updateFunc(ctx, user)
}

func (m *mockUserRepository) UpdateProfile(ctx context.Context, id int64, name string, phone string) error {
	if m.updateProfileFunc == nil {
		return nil
	}
	return m.updateProfileFunc(ctx, id, name, phone)
}

func (m *mockUserRepository) UpdatePasswordHash(ctx context.Context, id int64, passwordHash string) error {
	if m.updatePasswordHashFunc == nil {
		return nil
	}
	return m.updatePasswordHashFunc(ctx, id, passwordHash)
}

func (m *mockUserRepository) UpdateAvatarURL(ctx context.Context, id int64, avatarURL string) error {
	if m.updateAvatarURLFunc == nil {
		return nil
	}
	return m.updateAvatarURLFunc(ctx, id, avatarURL)
}

func (m *mockUserRepository) UpdateProfileVideoURL(ctx context.Context, id int64, profileVideoURL string) error {
	if m.updateProfileVideoURLFunc == nil {
		return nil
	}
	return m.updateProfileVideoURLFunc(ctx, id, profileVideoURL)
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
	restoreFunc        func(ctx context.Context, id int64) error
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

func (m *mockProductRepository) Restore(ctx context.Context, id int64) error {
	if m.restoreFunc == nil {
		return nil
	}
	return m.restoreFunc(ctx, id)
}

type mockCategoryRepository struct {
	createFunc                    func(ctx context.Context, category *model.Category) error
	findByIDFunc                  func(ctx context.Context, id int64) (*model.Category, error)
	findActiveByIDFunc            func(ctx context.Context, id int64) (*model.Category, error)
	listActiveFunc                func(ctx context.Context) ([]model.Category, error)
	listAdminFunc                 func(ctx context.Context, status string) ([]model.Category, error)
	existsByNameOtherCategoryFunc func(ctx context.Context, name string, categoryID int64) (bool, error)
	hasActiveProductsFunc         func(ctx context.Context, categoryID int64) (bool, error)
	updateFunc                    func(ctx context.Context, category *model.Category) error
	softDeleteFunc                func(ctx context.Context, id int64) error
	restoreFunc                   func(ctx context.Context, id int64) error
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

func (m *mockCategoryRepository) ListAdmin(ctx context.Context, status string) ([]model.Category, error) {
	if m.listAdminFunc == nil {
		return nil, nil
	}
	return m.listAdminFunc(ctx, status)
}

func (m *mockCategoryRepository) ExistsByNameOtherCategory(ctx context.Context, name string, categoryID int64) (bool, error) {
	return m.existsByNameOtherCategoryFunc(ctx, name, categoryID)
}

func (m *mockCategoryRepository) HasActiveProducts(ctx context.Context, categoryID int64) (bool, error) {
	if m.hasActiveProductsFunc == nil {
		return false, nil
	}
	return m.hasActiveProductsFunc(ctx, categoryID)
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	return m.updateFunc(ctx, category)
}

func (m *mockCategoryRepository) SoftDelete(ctx context.Context, id int64) error {
	return m.softDeleteFunc(ctx, id)
}

func (m *mockCategoryRepository) Restore(ctx context.Context, id int64) error {
	if m.restoreFunc == nil {
		return nil
	}
	return m.restoreFunc(ctx, id)
}

type mockOrderRepository struct {
	products                map[int64]*model.Product
	listByUserIDFunc        func(ctx context.Context, db repository.Queryer, userID int64, query dto.OrderListQuery) ([]model.Order, int64, error)
	listAllFunc             func(ctx context.Context, db repository.Queryer, query dto.OrderListQuery) ([]model.Order, int64, error)
	findByIDFunc            func(ctx context.Context, db repository.Queryer, orderID int64) (*model.Order, error)
	updateStatusFunc        func(ctx context.Context, db repository.Queryer, orderID int64, status string) error
	findItemsByOrderIDFunc  func(ctx context.Context, db repository.Queryer, orderID int64) ([]model.OrderItem, error)
	findItemsByOrderIDsFunc func(ctx context.Context, db repository.Queryer, orderIDs []int64) (map[int64][]model.OrderItem, error)
}

type mockPaymentRepository struct {
	createFunc                         func(ctx context.Context, db repository.Queryer, payment *model.Payment) error
	findLatestByOrderIDAndProviderFunc func(ctx context.Context, db repository.Queryer, orderID int64, provider string) (*model.Payment, error)
	findByTransactionIDFunc            func(ctx context.Context, db repository.Queryer, transactionID string) (*model.Payment, error)
	findByAppTransactionIDFunc         func(ctx context.Context, db repository.Queryer, appTransactionID string) (*model.Payment, error)
	updateGatewayInitializationFunc    func(ctx context.Context, db repository.Queryer, payment *model.Payment) error
	updateStatusFunc                   func(ctx context.Context, db repository.Queryer, payment *model.Payment) error
	updateRawCallbackFunc              func(ctx context.Context, db repository.Queryer, paymentID int64, rawCallback string) error
}

type mockZaloPayClient struct {
	createOrderFunc    func(ctx context.Context, req ZaloPayCreateOrderRequest) (*ZaloPayCreateOrderResponse, string, error)
	queryOrderFunc     func(ctx context.Context, appTransID string) (*ZaloPayQueryOrderResponse, string, error)
	verifyCallbackFunc func(data string, mac string) bool
}

type mockOAuthAccountRepository struct {
	createWithQuerierFunc    func(ctx context.Context, q repository.Queryer, account *model.OAuthAccount) error
	findByProviderUserIDFunc func(ctx context.Context, provider string, providerUserID string) (*model.OAuthAccount, error)
}

func (m *mockOAuthAccountRepository) CreateWithQuerier(ctx context.Context, q repository.Queryer, account *model.OAuthAccount) error {
	if m.createWithQuerierFunc != nil {
		return m.createWithQuerierFunc(ctx, q, account)
	}
	return nil
}

func (m *mockOAuthAccountRepository) FindByProviderUserID(ctx context.Context, provider string, providerUserID string) (*model.OAuthAccount, error) {
	if m.findByProviderUserIDFunc != nil {
		return m.findByProviderUserIDFunc(ctx, provider, providerUserID)
	}
	return nil, nil
}

type mockGoogleProvider struct {
	enabledFunc     func() bool
	authCodeURLFunc func(state string) string
	exchangeFunc    func(ctx context.Context, code string) (*oauth.GoogleUserInfo, error)
}

func (m *mockGoogleProvider) Enabled() bool {
	if m.enabledFunc != nil {
		return m.enabledFunc()
	}
	return true
}

func (m *mockGoogleProvider) AuthCodeURL(state string) string {
	if m.authCodeURLFunc != nil {
		return m.authCodeURLFunc(state)
	}
	return ""
}

func (m *mockGoogleProvider) Exchange(ctx context.Context, code string) (*oauth.GoogleUserInfo, error) {
	if m.exchangeFunc != nil {
		return m.exchangeFunc(ctx, code)
	}
	return nil, nil
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

func (m *mockOrderRepository) ListByUserID(ctx context.Context, db repository.Queryer, userID int64, query dto.OrderListQuery) ([]model.Order, int64, error) {
	if m.listByUserIDFunc != nil {
		return m.listByUserIDFunc(ctx, db, userID, query)
	}
	return nil, 0, nil
}

func (m *mockOrderRepository) ListAll(ctx context.Context, db repository.Queryer, query dto.OrderListQuery) ([]model.Order, int64, error) {
	if m.listAllFunc != nil {
		return m.listAllFunc(ctx, db, query)
	}
	return nil, 0, nil
}

func (m *mockOrderRepository) FindByID(ctx context.Context, db repository.Queryer, orderID int64) (*model.Order, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, db, orderID)
	}
	return nil, nil
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, db repository.Queryer, orderID int64, status string) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, db, orderID, status)
	}
	return nil
}

func (m *mockOrderRepository) FindItemsByOrderID(ctx context.Context, db repository.Queryer, orderID int64) ([]model.OrderItem, error) {
	if m.findItemsByOrderIDFunc != nil {
		return m.findItemsByOrderIDFunc(ctx, db, orderID)
	}
	return nil, nil
}

func (m *mockOrderRepository) FindItemsByOrderIDs(ctx context.Context, db repository.Queryer, orderIDs []int64) (map[int64][]model.OrderItem, error) {
	if m.findItemsByOrderIDsFunc != nil {
		return m.findItemsByOrderIDsFunc(ctx, db, orderIDs)
	}
	itemsByOrderID := make(map[int64][]model.OrderItem, len(orderIDs))
	for _, orderID := range orderIDs {
		items, err := m.FindItemsByOrderID(ctx, db, orderID)
		if err != nil {
			return nil, err
		}
		itemsByOrderID[orderID] = items
	}
	return itemsByOrderID, nil
}

func (m *mockPaymentRepository) Create(ctx context.Context, db repository.Queryer, payment *model.Payment) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, db, payment)
	}
	payment.ID = 1
	return nil
}

func (m *mockPaymentRepository) FindLatestByOrderIDAndProvider(ctx context.Context, db repository.Queryer, orderID int64, provider string) (*model.Payment, error) {
	if m.findLatestByOrderIDAndProviderFunc != nil {
		return m.findLatestByOrderIDAndProviderFunc(ctx, db, orderID, provider)
	}
	return nil, nil
}

func (m *mockPaymentRepository) FindByTransactionID(ctx context.Context, db repository.Queryer, transactionID string) (*model.Payment, error) {
	if m.findByTransactionIDFunc != nil {
		return m.findByTransactionIDFunc(ctx, db, transactionID)
	}
	return nil, nil
}

func (m *mockPaymentRepository) FindByAppTransactionID(ctx context.Context, db repository.Queryer, appTransactionID string) (*model.Payment, error) {
	if m.findByAppTransactionIDFunc != nil {
		return m.findByAppTransactionIDFunc(ctx, db, appTransactionID)
	}
	return nil, nil
}

func (m *mockPaymentRepository) UpdateGatewayInitialization(ctx context.Context, db repository.Queryer, payment *model.Payment) error {
	if m.updateGatewayInitializationFunc != nil {
		return m.updateGatewayInitializationFunc(ctx, db, payment)
	}
	return nil
}

func (m *mockPaymentRepository) UpdateStatus(ctx context.Context, db repository.Queryer, payment *model.Payment) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, db, payment)
	}
	return nil
}

func (m *mockPaymentRepository) UpdateRawCallback(ctx context.Context, db repository.Queryer, paymentID int64, rawCallback string) error {
	if m.updateRawCallbackFunc != nil {
		return m.updateRawCallbackFunc(ctx, db, paymentID, rawCallback)
	}
	return nil
}

func (m *mockZaloPayClient) CreateOrder(ctx context.Context, req ZaloPayCreateOrderRequest) (*ZaloPayCreateOrderResponse, string, error) {
	if m.createOrderFunc != nil {
		return m.createOrderFunc(ctx, req)
	}
	return nil, "", nil
}

func (m *mockZaloPayClient) QueryOrder(ctx context.Context, appTransID string) (*ZaloPayQueryOrderResponse, string, error) {
	if m.queryOrderFunc != nil {
		return m.queryOrderFunc(ctx, appTransID)
	}
	return nil, "", nil
}

func (m *mockZaloPayClient) VerifyCallback(data string, mac string) bool {
	if m.verifyCallbackFunc != nil {
		return m.verifyCallbackFunc(data, mac)
	}
	return false
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
