# Architecture Layers - Implementation Guide

## Overview
Backend sử dụng 4-layer architecture: Handler → Service → Repository → Database

```
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT (Frontend)                        │
├─────────────────────────────────────────────────────────────┤
│                  Echo Router & Middleware                   │
│     (CORS, Logger, Recover, JWT Auth, Role Check)          │
├─────────────────────────────────────────────────────────────┤
│   HANDLER LAYER (Internal/Handler)                          │
│   - Parse request                                           │
│   - Validate input                                          │
│   - Call service                                            │
│   - Return response                                         │
├─────────────────────────────────────────────────────────────┤
│   SERVICE LAYER (Internal/Service)                          │
│   - Implement business logic                                │
│   - Validate business rules                                 │
│   - Call repository                                         │
│   - Return domain objects                                   │
├─────────────────────────────────────────────────────────────┤
│   REPOSITORY LAYER (Internal/Repository)                    │
│   - Query builder                                           │
│   - Execute SQL                                             │
│   - Error handling & mapping                                │
│   - Return model objects                                    │
├─────────────────────────────────────────────────────────────┤
│            DATABASE (PostgreSQL)                            │
│   - Tables, indexes, constraints                            │
└─────────────────────────────────────────────────────────────┘
```

---

## Layer 1: HANDLER LAYER

**Location:** `internal/handler/`

**Responsibility:**
1. Receive HTTP request
2. Parse request body/params
3. Validate with validator
4. Call corresponding service
5. Handle service errors (convert to HTTP errors)
6. Return formatted response

### File Structure (6 files)
```
internal/handler/
├── health_handler.go ✅ (EXISTS)
├── auth_handler.go ❌ TODO
├── user_handler.go ❌ TODO
├── product_handler.go ❌ TODO
├── category_handler.go ❌ TODO
└── order_handler.go ❌ TODO
```

### Template Pattern

```go
package handler

import (
    "net/http"
    "enterprise-order-management-api/backend/internal/dto"
    "enterprise-order-management-api/backend/internal/service"
    "enterprise-order-management-api/backend/internal/util"
    "github.com/labstack/echo/v4"
)

type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register new user
// @Description Create new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} util.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
    var req dto.RegisterRequest
    
    // Bind & Validate
    if err := c.BindAndValidate(&req); err != nil {
        return err // Validator handles this
    }
    
    // Call Service
    user, err := h.authService.Register(req)
    if err != nil {
        return err // Service returns AppError
    }
    
    // Return Response
    return util.Success(c, http.StatusCreated, "User registered successfully", user)
}
```

### Rules
- ✅ Handler chỉ parse request, validate, gọi service, return response
- ❌ NO business logic
- ❌ NO database queries
- ❌ NO SQL
- ✅ All service errors are AppError → HTTP handled by HTTPErrorHandler

---

## Layer 2: SERVICE LAYER

**Location:** `internal/service/`

**Responsibility:**
1. Implement business logic
2. Validate business rules (not input format)
3. Call repository for data
4. Process/transform data
5. Return domain objects or AppError
6. Handle transactions (if needed)

### File Structure (6 files)
```
internal/service/
├── auth_service.go ❌ TODO (100+ lines)
├── user_service.go ❌ TODO (60+ lines)
├── product_service.go ❌ TODO (70+ lines)
├── category_service.go ❌ TODO (40+ lines)
├── order_service.go ❌ TODO (120+ lines - MOST COMPLEX)
└── role_service.go ❌ TODO (20+ lines)
```

### Template Pattern

```go
package service

import (
    "enterprise-order-management-api/backend/internal/dto"
    "enterprise-order-management-api/backend/internal/model"
    "enterprise-order-management-api/backend/internal/repository"
    "enterprise-order-management-api/backend/internal/util"
)

type AuthService interface {
    Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
    Login(req dto.LoginRequest) (*dto.AuthResponse, error)
    RefreshToken(req dto.RefreshTokenRequest) (*dto.AuthResponse, error)
    Logout(req dto.LogoutRequest) error
    GetProfile(userID int64) (*dto.ProfileResponse, error)
}

type authService struct {
    userRepo repository.UserRepository
    refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthService(
    userRepo repository.UserRepository,
    refreshTokenRepo repository.RefreshTokenRepository,
) AuthService {
    return &authService{
        userRepo: userRepo,
        refreshTokenRepo: refreshTokenRepo,
    }
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
    // 1. Validate business rules
    exists, err := s.userRepo.GetByEmail(req.Email)
    if err != nil && err != repository.ErrNotFound {
        return nil, util.InternalServerError("Database error")
    }
    if exists != nil {
        return nil, util.Conflict("Email already registered")
    }
    
    // 2. Hash password
    hashedPassword, err := util.HashPassword(req.Password)
    if err != nil {
        return nil, util.InternalServerError("Failed to hash password")
    }
    
    // 3. Create user
    user := &model.User{
        FullName:     req.FullName,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        RoleID:       2, // user role
        IsActive:     true,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, util.InternalServerError("Failed to create user")
    }
    
    // 4. Generate tokens
    accessToken, err := util.GenerateToken(user.ID, "user", config.JWTAccessSecret, ...)
    if err != nil {
        return nil, util.InternalServerError("Failed to generate token")
    }
    
    refreshToken, err := util.GenerateToken(user.ID, "user", config.JWTRefreshSecret, ...)
    if err != nil {
        return nil, util.InternalServerError("Failed to generate token")
    }
    
    // 5. Save refresh token (hashed)
    refreshTokenHash := util.HashToken(refreshToken) // TODO: implement
    if err := s.refreshTokenRepo.Create(...); err != nil {
        return nil, util.InternalServerError("Failed to save refresh token")
    }
    
    // 6. Return response
    return &dto.AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        User:         mapUserToResponse(user),
    }, nil
}
```

### Rules
- ✅ Service chứa tất cả business logic
- ✅ Service validate business rules
- ✅ Service call repository
- ✅ Service return AppError hoặc domain object
- ❌ NO HTTP status codes in service (only in handler)
- ❌ NO echo.Context in service
- ❌ NO direct database calls (use repository)

---

## Layer 3: REPOSITORY LAYER

**Location:** `internal/repository/`

**Responsibility:**
1. Build SQL queries
2. Execute queries against database
3. Map database results to models
4. Handle database errors
5. Return models or repository errors

### File Structure (7 files)
```
internal/repository/
├── user_repository.go ❌ TODO (70+ lines)
├── product_repository.go ❌ TODO (80+ lines)
├── category_repository.go ❌ TODO (50+ lines)
├── order_repository.go ❌ TODO (100+ lines - MOST COMPLEX)
├── refresh_token_repository.go ❌ TODO (40+ lines)
├── role_repository.go ❌ TODO (20+ lines)
└── errors.go ✅ TODO (custom error types)
```

### Template Pattern

```go
package repository

import (
    "context"
    "errors"
    "enterprise-order-management-api/backend/internal/model"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Custom errors
var (
    ErrNotFound = errors.New("not found")
    ErrConflict = errors.New("conflict")
)

type UserRepository interface {
    Create(user *model.User) error
    GetByID(id int64) (*model.User, error)
    GetByEmail(email string) (*model.User, error)
    Update(user *model.User) error
    Delete(id int64) error
    List(page, limit int, search string) ([]*model.User, int, error)
}

type userRepository struct {
    pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
    return &userRepository{pool: pool}
}

func (r *userRepository) Create(user *model.User) error {
    ctx := context.Background()
    
    sql := `
        INSERT INTO users (full_name, email, password_hash, role_id, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    err := r.pool.QueryRow(ctx, sql,
        user.FullName,
        user.Email,
        user.PasswordHash,
        user.RoleID,
        user.IsActive,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        // Check for unique constraint violation
        if pgErr := err.(*pgconn.PgError); pgErr != nil && pgErr.Code == "23505" {
            return ErrConflict
        }
        return err
    }
    
    return nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
    ctx := context.Background()
    
    sql := `
        SELECT id, full_name, email, password_hash, role_id, is_active, created_at, updated_at
        FROM users
        WHERE email = $1 AND is_active = true
    `
    
    user := &model.User{}
    err := r.pool.QueryRow(ctx, sql, email).Scan(
        &user.ID,
        &user.FullName,
        &user.Email,
        &user.PasswordHash,
        &user.RoleID,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, err
    }
    
    return user, nil
}

func (r *userRepository) List(page, limit int, search string) ([]*model.User, int, error) {
    ctx := context.Background()
    
    // Get total count
    var total int
    countSQL := `SELECT COUNT(*) FROM users WHERE is_active = true`
    if search != "" {
        countSQL += ` AND (full_name ILIKE $1 OR email ILIKE $1)`
        r.pool.QueryRow(ctx, countSQL, "%"+search+"%").Scan(&total)
    } else {
        r.pool.QueryRow(ctx, countSQL).Scan(&total)
    }
    
    // Get paginated results
    offset := (page - 1) * limit
    sql := `
        SELECT id, full_name, email, password_hash, role_id, is_active, created_at, updated_at
        FROM users
        WHERE is_active = true
    `
    
    if search != "" {
        sql += ` AND (full_name ILIKE $1 OR email ILIKE $1)`
    }
    
    sql += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
    
    var rows pgx.Rows
    var err error
    
    if search != "" {
        rows, err = r.pool.Query(ctx, sql, "%"+search+"%", limit, offset)
    } else {
        rows, err = r.pool.Query(ctx, sql, limit, offset)
    }
    
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    users := make([]*model.User, 0, limit)
    for rows.Next() {
        user := &model.User{}
        err := rows.Scan(
            &user.ID,
            &user.FullName,
            &user.Email,
            &user.PasswordHash,
            &user.RoleID,
            &user.IsActive,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, 0, err
        }
        users = append(users, user)
    }
    
    return users, total, nil
}
```

### Rules
- ✅ Repository chỉ query database
- ✅ Repository use parameterized queries ($1, $2, ...)
- ✅ Repository map rows to models
- ✅ Repository handle database errors
- ✅ Repository return custom errors (ErrNotFound, ErrConflict)
- ❌ NO business logic
- ❌ NO HTTP status codes
- ❌ NO string concatenation in SQL

---

## Layer 4: MIDDLEWARE LAYER

**Location:** `internal/middleware/`

**Responsibility:**
1. Intercept requests
2. Perform cross-cutting concerns
3. Extract information (JWT claims)
4. Validate authorization
5. Pass to next handler or return error

### File Structure (3 files)
```
internal/middleware/
├── jwt_auth_middleware.go ❌ TODO (40+ lines)
├── role_authorization_middleware.go ❌ TODO (30+ lines)
└── cors_middleware.go ❌ TODO (20+ lines)
```

### JWT Auth Middleware Example

```go
package middleware

import (
    "strings"
    "enterprise-order-management-api/backend/internal/util"
    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

func JWTAuthMiddleware(secret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Get token from Authorization header
            auth := c.Request().Header.Get("Authorization")
            if auth == "" {
                return util.Unauthorized("Missing authorization header")
            }
            
            // Check Bearer prefix
            parts := strings.Split(auth, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                return util.Unauthorized("Invalid authorization header format")
            }
            
            token := parts[1]
            
            // Parse & verify token
            claims := &util.Claims{}
            parsedToken, err := jwt.ParseWithClaims(token, claims,
                func(token *jwt.Token) (any, error) {
                    return []byte(secret), nil
                })
            
            if err != nil || !parsedToken.Valid {
                return util.Unauthorized("Invalid token")
            }
            
            // Store claims in context
            c.Set("user_id", claims.UserID)
            c.Set("role", claims.Role)
            
            return next(c)
        }
    }
}
```

### Role Authorization Middleware Example

```go
func RoleAuthorization(requiredRole string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // JWT middleware must run first
            userRole, ok := c.Get("role").(string)
            if !ok {
                return util.Unauthorized("User not authenticated")
            }
            
            if userRole != requiredRole && requiredRole != "*" {
                return util.Forbidden("You don't have permission to access this resource")
            }
            
            return next(c)
        }
    }
}
```

---

## COMPLETE REQUEST FLOW EXAMPLE

### POST /auth/register

```
1. CLIENT sends request
   POST /auth/register
   {
     "full_name": "Nguyễn A",
     "email": "user@example.com",
     "password": "password123"
   }

2. ROUTER matches route
   echo.POST("/auth/register", handler.Register)

3. MIDDLEWARE runs in order
   ✅ Logger middleware (log request)
   ✅ Recover middleware (catch panics)
   ✅ CORS middleware (check origin)
   ✓ NO JWT middleware (public endpoint)

4. HANDLER processes request
   - c.BindAndValidate(&req) → Validate input
   - handler.Register(c)
     ├─ Extract request body
     ├─ Call validator
     └─ Call service

5. SERVICE processes business logic
   - authService.Register(req)
     ├─ Check email doesn't exist
     ├─ Hash password
     ├─ Call userRepo.Create()
     ├─ Generate tokens
     ├─ Call refreshTokenRepo.Create()
     └─ Return AuthResponse

6. REPOSITORY executes database operation
   - userRepository.Create(user)
     ├─ Build SQL INSERT
     ├─ Execute with parameterized query
     ├─ Scan returned ID
     ├─ Handle errors (ErrConflict on duplicate email)
     └─ Return user with ID

7. RESPONSE returned to client
   HTTP 201 Created
   {
     "success": true,
     "message": "User registered successfully",
     "data": {
       "access_token": "eyJ...",
       "refresh_token": "eyJ...",
       "token_type": "Bearer",
       "user": { ... }
     }
   }

8. CLIENT receives response
   ✅ Store tokens
   ✅ Store user info
   ✅ Redirect to home
```

---

## DEPENDENCY INJECTION PATTERN

**In main.go or a DI container file:**

```go
// Initialize database
pool, _ := database.ConnectDB(ctx, cfg)

// Initialize repositories
userRepo := repository.NewUserRepository(pool)
productRepo := repository.NewProductRepository(pool)
categoryRepo := repository.NewCategoryRepository(pool)
orderRepo := repository.NewOrderRepository(pool)
refreshTokenRepo := repository.NewRefreshTokenRepository(pool)

// Initialize services
authService := service.NewAuthService(userRepo, refreshTokenRepo)
productService := service.NewProductService(productRepo)
categoryService := service.NewCategoryService(categoryRepo)
orderService := service.NewOrderService(orderRepo, productRepo)

// Initialize handlers
authHandler := handler.NewAuthHandler(authService)
productHandler := handler.NewProductHandler(productService)
categoryHandler := handler.NewCategoryHandler(categoryService)
orderHandler := handler.NewOrderHandler(orderService)

// Register routes
e := echo.New()
route.Register(e, authHandler, productHandler, categoryHandler, orderHandler)
```

---

## TESTING EACH LAYER

### Unit Test Repository
```go
func TestUserRepository_GetByEmail(t *testing.T) {
    // Mock database pool
    pool := setupTestDB()
    repo := repository.NewUserRepository(pool)
    
    // Test: Should return user
    user, err := repo.GetByEmail("user@example.com")
    assert.NoError(t, err)
    assert.Equal(t, "user@example.com", user.Email)
    
    // Test: Should return ErrNotFound
    user, err := repo.GetByEmail("notfound@example.com")
    assert.Equal(t, repository.ErrNotFound, err)
}
```

### Unit Test Service
```go
func TestAuthService_Register(t *testing.T) {
    // Mock repository
    mockUserRepo := &MockUserRepository{}
    mockRefreshTokenRepo := &MockRefreshTokenRepository{}
    
    service := service.NewAuthService(mockUserRepo, mockRefreshTokenRepo)
    
    // Test: Should register successfully
    req := dto.RegisterRequest{...}
    resp, err := service.Register(req)
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.AccessToken)
}
```

### Integration Test Handler
```go
func TestAuthHandler_Register(t *testing.T) {
    // Setup real database
    pool := setupRealDB()
    
    // Initialize layers
    userRepo := repository.NewUserRepository(pool)
    authService := service.NewAuthService(userRepo, ...)
    authHandler := handler.NewAuthHandler(authService)
    
    // Make request
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/auth/register", ...)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    
    err := authHandler.Register(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, rec.Code)
}
```

---

*For detailed endpoints, see: API_ENDPOINTS_PLAN.md*
*For quick reference, see: QUICK_SUMMARY.md*
