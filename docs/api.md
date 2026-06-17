# API Documentation

## 1. Overview

Base URL:

```text
http://localhost:8080/api/v1
```

Authentication dùng JWT access token qua header:

```http
Authorization: Bearer <ACCESS_TOKEN>
```

Refresh token được gửi qua JSON body để đơn giản cho demo thực tập. API không dùng cookie ở giai đoạn này.

## 2. Response Format

Success:

```json
{
  "success": true,
  "message": "Success",
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {}
}
```

Pagination:

```json
{
  "success": true,
  "message": "Success",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 0,
    "total_pages": 0
  }
}
```

## 3. Auth API

### 3.1 Register

Method: `POST`

URL: `/auth/register`

Auth required: No

Role required: Guest

Request body:

```json
{
  "name": "Nguyen Van A",
  "email": "user@example.com",
  "password": "123456"
}
```

Success response:

```json
{
  "success": true,
  "message": "Created successfully",
  "data": {
    "access_token": "ACCESS_TOKEN",
    "refresh_token": "REFRESH_TOKEN",
    "user": {
      "id": 2,
      "name": "Nguyen Van A",
      "email": "user@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2026-06-17T10:00:00Z",
      "updated_at": "2026-06-17T10:00:00Z"
    }
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Email already exists",
  "errors": {
    "code": "CONFLICT"
  }
}
```

### 3.2 Login

Method: `POST`

URL: `/auth/login`

Auth required: No

Role required: Guest

Request body:

```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "access_token": "ACCESS_TOKEN",
    "refresh_token": "REFRESH_TOKEN",
    "user": {
      "id": 2,
      "name": "Nguyen Van A",
      "email": "user@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2026-06-17T10:00:00Z",
      "updated_at": "2026-06-17T10:00:00Z"
    }
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Invalid email or password",
  "errors": {
    "code": "UNAUTHORIZED"
  }
}
```

### 3.3 Refresh Token

Method: `POST`

URL: `/auth/refresh-token`

Auth required: No

Role required: Guest

Request body:

```json
{
  "refresh_token": "REFRESH_TOKEN"
}
```

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "access_token": "NEW_ACCESS_TOKEN",
    "refresh_token": "NEW_REFRESH_TOKEN",
    "user": {
      "id": 2,
      "name": "Nguyen Van A",
      "email": "user@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2026-06-17T10:00:00Z",
      "updated_at": "2026-06-17T10:00:00Z"
    }
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Refresh token was revoked",
  "errors": {
    "code": "UNAUTHORIZED"
  }
}
```

### 3.4 Logout

Method: `POST`

URL: `/auth/logout`

Auth required: Yes

Role required: User or Admin

Request body:

```json
{
  "refresh_token": "REFRESH_TOKEN"
}
```

Success response:

```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

Error response:

```json
{
  "success": false,
  "message": "Invalid refresh token",
  "errors": {
    "code": "UNAUTHORIZED"
  }
}
```

### 3.5 Get Current User

Method: `GET`

URL: `/auth/me`

Auth required: Yes

Role required: User or Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 2,
    "name": "Nguyen Van A",
    "email": "user@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2026-06-17T10:00:00Z",
    "updated_at": "2026-06-17T10:00:00Z"
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "invalid access token",
  "errors": {
    "code": "UNAUTHORIZED"
  }
}
```

## 4. User API

### 4.1 List Users

Method: `GET`

URL: `/users?page=1&limit=10&search=admin`

Auth required: Yes

Role required: Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Admin",
      "email": "admin@example.com",
      "role": "admin",
      "is_active": true,
      "created_at": "2026-06-17T10:00:00Z",
      "updated_at": "2026-06-17T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "you do not have permission to access this resource",
  "errors": {
    "code": "FORBIDDEN"
  }
}
```

### 4.2 Get User Detail

Method: `GET`

URL: `/users/{id}`

Auth required: Yes

Role required: Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 2,
    "name": "Nguyen Van A",
    "email": "user@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2026-06-17T10:00:00Z",
    "updated_at": "2026-06-17T10:00:00Z"
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "User not found",
  "errors": {
    "code": "NOT_FOUND"
  }
}
```

### 4.3 Update User

Method: `PUT`

URL: `/users/{id}`

Auth required: Yes

Role required: Admin

Request body:

```json
{
  "name": "Nguyen Van B",
  "email": "user2@example.com",
  "role": "user"
}
```

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 2,
    "name": "Nguyen Van B",
    "email": "user2@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2026-06-17T10:00:00Z",
    "updated_at": "2026-06-17T10:05:00Z"
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Email already exists",
  "errors": {
    "code": "CONFLICT"
  }
}
```

### 4.4 Delete User

Method: `DELETE`

URL: `/users/{id}`

Auth required: Yes

Role required: Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

Error response:

```json
{
  "success": false,
  "message": "Admin cannot delete own account",
  "errors": {
    "code": "BAD_REQUEST"
  }
}
```

## 5. Category API

### 5.1 List Categories

Method: `GET`

URL: `/categories`

Auth required: No

Role required: Guest, User, Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic devices",
      "is_active": true
    }
  ]
}
```

Error response:

```json
{
  "success": false,
  "message": "Unexpected server error",
  "errors": {
    "code": "INTERNAL_ERROR"
  }
}
```

### 5.2 Get Category Detail

Method: `GET`

URL: `/categories/{id}`

Auth required: No

Role required: Guest, User, Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices",
    "is_active": true
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Category not found",
  "errors": {
    "code": "NOT_FOUND"
  }
}
```

### 5.3 Create Category

Method: `POST`

URL: `/categories`

Auth required: Yes

Role required: Admin

Request body:

```json
{
  "name": "Furniture",
  "description": "Office furniture",
  "is_active": true
}
```

Success response:

```json
{
  "success": true,
  "message": "Created successfully",
  "data": {
    "id": 3,
    "name": "Furniture",
    "description": "Office furniture",
    "is_active": true
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Category name already exists",
  "errors": {
    "code": "CONFLICT"
  }
}
```

### 5.4 Update Category

Method: `PUT`

URL: `/categories/{id}`

Auth required: Yes

Role required: Admin

Request body:

```json
{
  "name": "Furniture Updated",
  "description": "Updated description",
  "is_active": true
}
```

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 3,
    "name": "Furniture Updated",
    "description": "Updated description",
    "is_active": true
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Category not found",
  "errors": {
    "code": "NOT_FOUND"
  }
}
```

### 5.5 Delete Category

Method: `DELETE`

URL: `/categories/{id}`

Auth required: Yes

Role required: Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Category deleted successfully"
}
```

Error response:

```json
{
  "success": false,
  "message": "Category not found",
  "errors": {
    "code": "NOT_FOUND"
  }
}
```

## 6. Product API

### 6.1 List Products

Method: `GET`

URL: `/products?page=1&limit=10&keyword=laptop&category_id=1&min_price=100000&max_price=50000000`

Auth required: No

Role required: Guest, User, Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "category_id": 1,
      "name": "Laptop Dell",
      "description": "Business laptop",
      "price": 15000000,
      "stock": 20,
      "image_url": "https://example.com/laptop.jpg",
      "is_active": true,
      "category": {
        "id": 1,
        "name": "Electronics",
        "description": "Electronic devices",
        "is_active": true
      }
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Unexpected server error",
  "errors": {
    "code": "INTERNAL_ERROR"
  }
}
```

### 6.2 Get Product Detail

Method: `GET`

URL: `/products/{id}`

Auth required: No

Role required: Guest, User, Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 1,
    "category_id": 1,
    "name": "Laptop Dell",
    "description": "Business laptop",
    "price": 15000000,
    "stock": 20,
    "image_url": "https://example.com/laptop.jpg",
    "is_active": true,
    "category": {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic devices",
      "is_active": true
    }
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Product not found",
  "errors": {
    "code": "NOT_FOUND"
  }
}
```

### 6.3 Create Product

Method: `POST`

URL: `/products`

Auth required: Yes

Role required: Admin

Request body:

```json
{
  "category_id": 1,
  "name": "Laptop Dell",
  "description": "Business laptop",
  "price": 15000000,
  "stock": 20,
  "image_url": "https://example.com/laptop.jpg",
  "is_active": true
}
```

Success response:

```json
{
  "success": true,
  "message": "Created successfully",
  "data": {
    "id": 1,
    "category_id": 1,
    "name": "Laptop Dell",
    "description": "Business laptop",
    "price": 15000000,
    "stock": 20,
    "image_url": "https://example.com/laptop.jpg",
    "is_active": true
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Category is invalid or inactive",
  "errors": {
    "code": "BAD_REQUEST"
  }
}
```

### 6.4 Update Product

Method: `PUT`

URL: `/products/{id}`

Auth required: Yes

Role required: Admin

Request body:

```json
{
  "category_id": 1,
  "name": "Laptop Dell Updated",
  "description": "Updated description",
  "price": 14000000,
  "stock": 15,
  "image_url": "https://example.com/laptop.jpg",
  "is_active": true
}
```

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 1,
    "category_id": 1,
    "name": "Laptop Dell Updated",
    "description": "Updated description",
    "price": 14000000,
    "stock": 15,
    "image_url": "https://example.com/laptop.jpg",
    "is_active": true
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Product not found",
  "errors": {
    "code": "NOT_FOUND"
  }
}
```

### 6.5 Delete Product

Method: `DELETE`

URL: `/products/{id}`

Auth required: Yes

Role required: Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

Error response:

```json
{
  "success": false,
  "message": "Product not found",
  "errors": {
    "code": "NOT_FOUND"
  }
}
```

## 7. Order API

### 7.1 Create Order

Method: `POST`

URL: `/orders`

Auth required: Yes

Role required: User or Admin

Request body:

```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    }
  ]
}
```

Success response:

```json
{
  "success": true,
  "message": "Created successfully",
  "data": {
    "id": 1,
    "user_id": 2,
    "status": "pending",
    "total_amount": 30000000,
    "items": [
      {
        "product_id": 1,
        "name": "Laptop Dell",
        "quantity": 2,
        "unit_price": 15000000,
        "subtotal": 30000000
      }
    ]
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Product Laptop Dell does not have enough stock",
  "errors": {
    "code": "BAD_REQUEST"
  }
}
```

### 7.2 List Orders

Method: `GET`

URL: `/orders`

Auth required: Yes

Role required: User or Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "status": "pending",
      "total_amount": 30000000,
      "items": []
    }
  ]
}
```

Error response:

```json
{
  "success": false,
  "message": "missing authorization header",
  "errors": {
    "code": "UNAUTHORIZED"
  }
}
```

### 7.3 Get Order Detail

Method: `GET`

URL: `/orders/{id}`

Auth required: Yes

Role required: User or Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 1,
    "user_id": 2,
    "status": "pending",
    "total_amount": 30000000,
    "items": [
      {
        "product_id": 1,
        "name": "Laptop Dell",
        "quantity": 2,
        "unit_price": 15000000,
        "subtotal": 30000000
      }
    ]
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "You cannot view another user's order",
  "errors": {
    "code": "FORBIDDEN"
  }
}
```

### 7.4 Get My Orders

Method: `GET`

URL: `/users/me/orders`

Auth required: Yes

Role required: User or Admin

Request body: None

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "status": "pending",
      "total_amount": 30000000,
      "items": [
        {
          "product_id": 1,
          "name": "Laptop Dell",
          "quantity": 2,
          "unit_price": 15000000,
          "subtotal": 30000000
        }
      ]
    }
  ]
}
```

Error response:

```json
{
  "success": false,
  "message": "invalid access token",
  "errors": {
    "code": "UNAUTHORIZED"
  }
}
```

### 7.5 Update Order Status

Method: `PUT`

URL: `/orders/{id}/status`

Auth required: Yes

Role required: Admin

Request body:

```json
{
  "status": "confirmed"
}
```

Success response:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "id": 1,
    "user_id": 2,
    "status": "confirmed",
    "total_amount": 30000000,
    "items": [
      {
        "product_id": 1,
        "name": "Laptop Dell",
        "quantity": 2,
        "unit_price": 15000000,
        "subtotal": 30000000
      }
    ]
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Invalid order status transition",
  "errors": {
    "code": "BAD_REQUEST"
  }
}
```

Valid status values:

- `pending`
- `confirmed`
- `shipping`
- `completed`
- `cancelled`

Valid status transitions:

- `pending -> confirmed`
- `pending -> cancelled`
- `confirmed -> shipping`
- `confirmed -> cancelled`
- `shipping -> completed`

## 8. API Test Flow

### Step 1: Register User

```powershell
curl.exe -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"name\":\"Nguyen Van A\",\"email\":\"user@example.com\",\"password\":\"123456\"}"
```

Save `access_token` as `USER_ACCESS_TOKEN`.

### Step 2: Login User

```powershell
curl.exe -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"user@example.com\",\"password\":\"123456\"}"
```

Save `access_token` as `USER_ACCESS_TOKEN`.

### Step 3: Login Admin

```powershell
curl.exe -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"admin@example.com\",\"password\":\"123456\"}"
```

Save `access_token` as `ADMIN_ACCESS_TOKEN`.

### Step 4: Create Category

```powershell
curl.exe -X POST http://localhost:8080/api/v1/categories -H "Authorization: Bearer ADMIN_ACCESS_TOKEN" -H "Content-Type: application/json" -d "{\"name\":\"Electronics Demo\",\"description\":\"Demo category\",\"is_active\":true}"
```

Save returned category `id` as `CATEGORY_ID`.

### Step 5: Create Product

```powershell
curl.exe -X POST http://localhost:8080/api/v1/products -H "Authorization: Bearer ADMIN_ACCESS_TOKEN" -H "Content-Type: application/json" -d "{\"category_id\":CATEGORY_ID,\"name\":\"Laptop Demo\",\"description\":\"Demo product\",\"price\":15000000,\"stock\":10,\"image_url\":\"https://example.com/laptop.jpg\",\"is_active\":true}"
```

Save returned product `id` as `PRODUCT_ID`.

### Step 6: Get Products

```powershell
curl.exe -X GET "http://localhost:8080/api/v1/products?page=1&limit=10&keyword=laptop"
```

### Step 7: Create Order

```powershell
curl.exe -X POST http://localhost:8080/api/v1/orders -H "Authorization: Bearer USER_ACCESS_TOKEN" -H "Content-Type: application/json" -d "{\"items\":[{\"product_id\":PRODUCT_ID,\"quantity\":2}]}"
```

Save returned order `id` as `ORDER_ID`.

### Step 8: Get My Orders

```powershell
curl.exe -X GET http://localhost:8080/api/v1/users/me/orders -H "Authorization: Bearer USER_ACCESS_TOKEN"
```

### Step 9: Admin Update Order Status

```powershell
curl.exe -X PUT http://localhost:8080/api/v1/orders/ORDER_ID/status -H "Authorization: Bearer ADMIN_ACCESS_TOKEN" -H "Content-Type: application/json" -d "{\"status\":\"confirmed\"}"
```

## 9. Manual Postman Import Guide

Postman có thể tạo collection thủ công theo các bước sau:

1. Tạo collection mới tên `Enterprise Order Management API`.
2. Tạo environment tên `Local`.
3. Thêm biến environment:
   - `base_url`: `http://localhost:8080/api/v1`
   - `user_access_token`: để trống ban đầu
   - `admin_access_token`: để trống ban đầu
   - `user_refresh_token`: để trống ban đầu
   - `admin_refresh_token`: để trống ban đầu
   - `category_id`: để trống ban đầu
   - `product_id`: để trống ban đầu
   - `order_id`: để trống ban đầu
4. Với mỗi request, dùng URL dạng:

```text
{{base_url}}/auth/login
```

5. Với API cần đăng nhập, vào tab `Authorization`, chọn `Bearer Token`, nhập:

```text
{{user_access_token}}
```

hoặc:

```text
{{admin_access_token}}
```

6. Với request có body, chọn tab `Body`, chọn `raw`, chọn `JSON`.
7. Sau khi login, copy `access_token` và `refresh_token` từ response vào environment.
8. Chạy flow test theo thứ tự ở mục 8 để kiểm tra toàn bộ nghiệp vụ chính.

## 10. Notes For Internship Report

- Backend sử dụng Golang, Echo v4, PostgreSQL và SQL thuần với pgx/pgxpool.
- API tuân thủ kiến trúc `handler/service/repository`.
- Handler nhận request, validate input và gọi service.
- Service xử lý business rule như phân quyền, kiểm tra stock, kiểm tra trạng thái order.
- Repository thao tác database bằng parameterized query, không nối chuỗi SQL trực tiếp với input người dùng.
- Authentication dùng JWT access token và refresh token.
- Password được hash bằng bcrypt.
- Refresh token được hash trước khi lưu database.
- DELETE user/category/product là soft delete bằng `is_active = false`.
- Tạo order dùng transaction để đảm bảo tạo order, tạo order_items và trừ stock thành công hoặc rollback toàn bộ.
