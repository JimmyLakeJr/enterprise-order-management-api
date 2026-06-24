# API Endpoints Implementation Plan

## Overview
Danh sách đầy đủ tất cả endpoints cần triển khai với Chi tiết yêu cầu, validation, và business logic.

---

## 🔐 AUTH ENDPOINTS

### 1. POST /auth/register
**Mục đích:** Đăng ký tài khoản mới

**Request:**
```json
{
  "full_name": "string",
  "email": "string",
  "password": "string"
}
```

**Validation:**
- `full_name`: required, min=2, max=100
- `email`: required, email format, max=255, unique in DB
- `password`: required, min=6, max=72

**Response Success (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "token_type": "Bearer",
    "user": {
      "id": 1,
      "full_name": "Nguyễn A",
      "email": "user@example.com",
      "role_id": 2,
      "role_name": "user",
      "is_active": true,
      "created_at": "2026-06-20T...",
      "updated_at": "2026-06-20T..."
    }
  }
}
```

**Error Cases:**
- 400: Email already exists → `util.Conflict("Email already registered")`
- 400: Validation failed
- 500: Database error

**Business Logic:**
- Hash password with bcrypt
- Create user with role_id=2 (user)
- Generate access token (15 mins)
- Generate refresh token (7 days)
- Save refresh token hash to DB
- Return both tokens + user info

**Repository Calls:**
- CheckEmailExists()
- CreateUser()
- CreateRefreshToken()

---

### 2. POST /auth/login
**Mục đích:** Đăng nhập với email/password

**Request:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Validation:**
- `email`: required, email format
- `password`: required

**Response Success (200):** (Same as register)

**Error Cases:**
- 401: Email not found / Wrong password
- 400: Validation failed
- 500: Database error

**Business Logic:**
- Get user by email
- Compare password (bcrypt)
- Generate new access + refresh tokens
- Save refresh token to DB
- Return tokens

**Repository Calls:**
- GetUserByEmail()
- CreateRefreshToken()

---

### 3. POST /auth/refresh
**Mục đích:** Lấy access token mới từ refresh token

**Request:**
```json
{
  "refresh_token": "string"
}
```

**Validation:**
- `refresh_token`: required

**Response Success (200):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "token_type": "Bearer"
  }
}
```

**Error Cases:**
- 401: Token not found / expired / revoked
- 401: Invalid token signature
- 400: Validation failed

**Business Logic:**
- Parse & verify refresh token
- Check exists in DB
- Check not expired
- Check not revoked (revoked_at is null)
- Generate new access token
- Return new tokens

**Repository Calls:**
- GetRefreshTokenByHash()

---

### 4. POST /auth/logout
**Mục đích:** Đăng xuất (revoke refresh token)

**Auth Required:** ✅ JWT

**Request:**
```json
{
  "refresh_token": "string"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Logged out successfully",
  "data": null
}
```

**Business Logic:**
- Mark refresh_token.revoked_at = NOW()
- Response success

**Repository Calls:**
- RevokeRefreshToken()

---

### 5. GET /auth/profile
**Mục đích:** Lấy thông tin profile người dùng hiện tại

**Auth Required:** ✅ JWT

**Query Params:** None

**Response Success (200):**
```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 1,
    "full_name": "Nguyễn A",
    "email": "user@example.com",
    "role_name": "user",
    "created_at": "2026-06-20T..."
  }
}
```

**Business Logic:**
- Extract user_id from JWT
- Get user from DB
- Return profile (exclude password_hash)

**Repository Calls:**
- GetUserByID()

---

## 👥 CATEGORY ENDPOINTS

### 6. POST /categories
**Mục đích:** Tạo category mới

**Auth Required:** ✅ JWT + **Admin Role**

**Request:**
```json
{
  "name": "string",
  "description": "string"
}
```

**Validation:**
- `name`: required, min=2, max=100, unique
- `description`: optional, max=1000

**Response Success (201):** (CategoryResponse)

**Repository Calls:**
- CreateCategory()

---

### 7. GET /categories
**Mục đích:** Lấy danh sách categories

**Auth Required:** ❌ No

**Query Params:**
- `page`: number (default=1)
- `limit`: number (default=10)
- `search`: string (optional, filter by name)
- `is_active`: boolean (default=true if not admin)

**Response Success (200):**
```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Electronics",
      "description": "...",
      "is_active": true,
      "created_at": "...",
      "updated_at": "..."
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

**Business Logic:**
- Filter by is_active=true if user not admin
- Support pagination
- Support search by name

**Repository Calls:**
- ListCategories(page, limit, search, is_active)

---

### 8. GET /categories/:id
**Mục đích:** Lấy chi tiết category

**Auth Required:** ❌ No

**Response Success (200):** (CategoryResponse)

**Error Cases:**
- 404: Category not found

**Repository Calls:**
- GetCategoryByID(id)

---

### 9. PATCH /categories/:id
**Mục đích:** Cập nhật category

**Auth Required:** ✅ JWT + **Admin Role**

**Request:** (UpdateCategoryRequest)

**Response Success (200):** (CategoryResponse)

**Error Cases:**
- 404: Category not found

**Repository Calls:**
- UpdateCategory()

---

### 10. DELETE /categories/:id
**Mục đích:** Xóa category (soft delete)

**Auth Required:** ✅ JWT + **Admin Role**

**Response Success (204):** No content

**Business Logic:**
- Set is_active = false
- Do NOT remove from DB

**Repository Calls:**
- DeleteCategory() - actually soft delete

---

## 📦 PRODUCT ENDPOINTS

### 11. POST /products
**Mục đích:** Tạo product mới

**Auth Required:** ✅ JWT + **Admin Role**

**Request:** (CreateProductRequest)
```json
{
  "category_id": 1,
  "name": "iPhone 15",
  "description": "...",
  "price": 999999,
  "stock": 100,
  "image_url": "https://..."
}
```

**Validation:**
- `category_id`: required, gt=0
- `name`: required, min=2, max=150
- `description`: optional, max=2000
- `price`: required, gte=0
- `stock`: required, gte=0
- `image_url`: optional, valid URL

**Business Logic:**
- Verify category exists
- Create product with is_active=true
- Price & stock must be >= 0

**Repository Calls:**
- GetCategoryByID()
- CreateProduct()

---

### 12. GET /products
**Mục đích:** Lấy danh sách products

**Auth Required:** ❌ No

**Query Params:**
- `page`: number (default=1)
- `limit`: number (default=10)
- `category_id`: number (optional)
- `search`: string (optional, filter by name)
- `sort_by`: string (price, name, created_at) (optional)
- `sort_order`: asc|desc (default=desc)

**Response Success (200):**
```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "category_id": 1,
      "category": { /* CategoryResponse */ },
      "name": "iPhone 15",
      "description": "...",
      "price": 999999,
      "stock": 100,
      "image_url": "https://...",
      "is_active": true,
      "created_at": "...",
      "updated_at": "..."
    }
  ],
  "meta": { ... }
}
```

**Business Logic:**
- Filter by is_active=true
- Support pagination
- Support search & category filter
- Support sorting

**Repository Calls:**
- ListProducts(page, limit, category_id, search, sort_by)

---

### 13. GET /products/:id
**Mục đích:** Lấy chi tiết product

**Auth Required:** ❌ No

**Response Success (200):** (ProductResponse with category)

**Error Cases:**
- 404: Product not found

**Repository Calls:**
- GetProductByID(id)
- GetCategoryByID(product.category_id)

---

### 14. PATCH /products/:id
**Mục đích:** Cập nhật product

**Auth Required:** ✅ JWT + **Admin Role**

**Request:** (UpdateProductRequest)

**Response Success (200):** (ProductResponse)

**Error Cases:**
- 404: Product not found
- 400: Invalid category_id

**Repository Calls:**
- UpdateProduct()
- GetCategoryByID() if category_id changed

---

### 15. DELETE /products/:id
**Mục đích:** Xóa product (soft delete)

**Auth Required:** ✅ JWT + **Admin Role**

**Response Success (204):** No content

**Business Logic:**
- Set is_active = false

**Repository Calls:**
- DeleteProduct() - soft delete

---

## 🛒 ORDER ENDPOINTS

### 16. POST /orders
**Mục đích:** Tạo đơn hàng mới

**Auth Required:** ✅ JWT

**Request:**
```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 3,
      "quantity": 1
    }
  ]
}
```

**Validation:**
- `items`: required, min=1, non-empty array
- `product_id`: required, gt=0
- `quantity`: required, gt=0

**Response Success (201):** (OrderResponse with items)

**Error Cases:**
- 400: Items empty
- 404: Product not found
- 400: Product not active
- 400: Stock insufficient
- 500: Transaction failed

**Business Logic:** ⚠️ CRITICAL - MUST BE IN TRANSACTION
1. Validate all items exist
2. Validate all products active
3. Get current prices from DB
4. Check all stock available
5. Create order with user_id, total_amount = SUM(quantity * price)
6. Create order_items with quantity, unit_price (from current price)
7. Reduce product.stock for each item
8. **If any error → ROLLBACK ALL**
9. Return order with items

**Repository Calls:**
- GetProductByID() × N items
- CreateOrder() + CreateOrderItems() + UpdateProductStock() [IN TRANSACTION]

---

### 17. GET /orders
**Mục đích:** Lấy danh sách orders

**Auth Required:** ✅ JWT

**Query Params:**
- `page`: number (default=1)
- `limit`: number (default=10)
- `status`: string (optional, filter by status)
- `user_id`: number (optional, admin only)

**Response Success (200):**
```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "total_amount": 999999,
      "status": "pending",
      "items": [ /* OrderItemResponse[] */ ],
      "created_at": "...",
      "updated_at": "..."
    }
  ],
  "meta": { ... }
}
```

**Business Logic:**
- If user (not admin): Only return user's own orders (user_id = JWT user_id)
- If admin: Return all orders (or filter by user_id if provided)

**Repository Calls:**
- ListOrders(user_id, page, limit, status)

---

### 18. GET /orders/:id
**Mục đích:** Lấy chi tiết order

**Auth Required:** ✅ JWT

**Response Success (200):** (OrderResponse with full items + product info)

**Error Cases:**
- 404: Order not found
- 403: Forbidden (not owner + not admin)

**Business Logic:**
- If user: Check order.user_id == JWT user_id
- If admin: Allow
- Else: Return 403

**Repository Calls:**
- GetOrderByID(id)
- GetOrderItems(order_id)
- GetProductByID() × N items

---

### 19. PATCH /orders/:id/status
**Mục đích:** Cập nhật status order

**Auth Required:** ✅ JWT + **Admin Role**

**Request:**
```json
{
  "status": "confirmed"
}
```

**Validation:**
- `status`: required, oneof=[pending, confirmed, shipping, completed, cancelled]

**Response Success (200):** (OrderResponse)

**Error Cases:**
- 404: Order not found
- 400: Invalid status transition

**Business Logic:**
- Current status validation:
  - pending → confirmed ✅ | cancelled ✅ | others ❌
  - confirmed → shipping ✅ | cancelled ✅ | others ❌
  - shipping → completed ✅ | others ❌
  - completed → (nothing)
  - cancelled → (nothing)
- Update order.status + order.updated_at

**Repository Calls:**
- GetOrderByID()
- UpdateOrderStatus()

---

## 👤 USER MANAGEMENT ENDPOINTS

### 20. GET /users
**Mục đích:** Lấy danh sách users

**Auth Required:** ✅ JWT + **Admin Role**

**Query Params:**
- `page`: number (default=1)
- `limit`: number (default=10)
- `search`: string (optional, filter by email/name)
- `role_id`: number (optional)
- `is_active`: boolean (optional)

**Response Success (200):** (List of UserResponse with pagination)

**Repository Calls:**
- ListUsers(page, limit, search, role_id, is_active)

---

### 21. GET /users/:id
**Mục đích:** Lấy chi tiết user

**Auth Required:** ✅ JWT (Admin or self)

**Response Success (200):** (UserResponse)

**Error Cases:**
- 404: User not found
- 403: Forbidden (not admin + not self)

**Repository Calls:**
- GetUserByID()

---

### 22. PATCH /users/:id
**Mục đích:** Cập nhật user

**Auth Required:** ✅ JWT (Admin or self)

**Request:**
```json
{
  "full_name": "Nguyễn B",
  "is_active": true,
  "role_id": 2
}
```

**Validation:**
- `full_name`: optional, min=2, max=100
- `is_active`: optional
- `role_id`: optional, gt=0

**Business Logic:**
- If user (not admin):
  - Only can update full_name
  - Cannot change is_active or role_id
  - Can only update self
- If admin:
  - Can update all fields
  - Including role_id

**Repository Calls:**
- UpdateUser()

---

### 23. DELETE /users/:id
**Mục đích:** Xóa user (soft delete)

**Auth Required:** ✅ JWT + **Admin Role**

**Response Success (204):** No content

**Business Logic:**
- Set is_active = false

**Repository Calls:**
- DeleteUser() - soft delete

---

## 📊 SUMMARY TABLE

| # | Method | Endpoint | Auth | Role | Status |
|---|--------|----------|------|------|--------|
| 1 | POST | /auth/register | ❌ | - | ❌ TODO |
| 2 | POST | /auth/login | ❌ | - | ❌ TODO |
| 3 | POST | /auth/refresh | ❌ | - | ❌ TODO |
| 4 | POST | /auth/logout | ✅ | - | ❌ TODO |
| 5 | GET | /auth/profile | ✅ | - | ❌ TODO |
| 6 | POST | /categories | ✅ | Admin | ❌ TODO |
| 7 | GET | /categories | ❌ | - | ❌ TODO |
| 8 | GET | /categories/:id | ❌ | - | ❌ TODO |
| 9 | PATCH | /categories/:id | ✅ | Admin | ❌ TODO |
| 10 | DELETE | /categories/:id | ✅ | Admin | ❌ TODO |
| 11 | POST | /products | ✅ | Admin | ❌ TODO |
| 12 | GET | /products | ❌ | - | ❌ TODO |
| 13 | GET | /products/:id | ❌ | - | ❌ TODO |
| 14 | PATCH | /products/:id | ✅ | Admin | ❌ TODO |
| 15 | DELETE | /products/:id | ✅ | Admin | ❌ TODO |
| 16 | POST | /orders | ✅ | User | ❌ TODO |
| 17 | GET | /orders | ✅ | - | ❌ TODO |
| 18 | GET | /orders/:id | ✅ | - | ❌ TODO |
| 19 | PATCH | /orders/:id/status | ✅ | Admin | ❌ TODO |
| 20 | GET | /users | ✅ | Admin | ❌ TODO |
| 21 | GET | /users/:id | ✅ | Admin\* | ❌ TODO |
| 22 | PATCH | /users/:id | ✅ | Admin\* | ❌ TODO |
| 23 | DELETE | /users/:id | ✅ | Admin | ❌ TODO |

**Total: 23 endpoints**
- ✅ Already implemented: 1 (Health check)
- ❌ To be implemented: 22

---

*Generated: June 20, 2026*
