# API Docs

Base URL: `http://localhost:8080/api/v1`

## Response format

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

## Auth

### Register

`POST /auth/register`

```json
{
  "name": "Nguyen Van A",
  "email": "user@example.com",
  "password": "123456"
}
```

### Login

`POST /auth/login`

```json
{
  "email": "admin@example.com",
  "password": "123456"
}
```

### Refresh token

`POST /auth/refresh`

```json
{
  "refresh_token": "..."
}
```

### Logout

`POST /auth/logout`

Header: `Authorization: Bearer ACCESS_TOKEN`

```json
{
  "refresh_token": "..."
}
```

## User

### Get profile

`GET /me`

Header: `Authorization: Bearer ACCESS_TOKEN`

### Admin list users

`GET /admin/users`

Role: `admin`

### Admin soft delete user

`DELETE /admin/users/{id}`

Role: `admin`

## Categories

### Public list categories

`GET /categories`

### Admin create category

`POST /admin/categories`

Role: `admin`

```json
{
  "name": "Electronics",
  "description": "Electronic devices",
  "is_active": true
}
```

### Admin update category

`PUT /admin/categories/{id}`

Role: `admin`

### Admin soft delete category

`DELETE /admin/categories/{id}`

Role: `admin`

## Products

### Public list products

`GET /products?page=1&limit=10&search=keyboard&category_id=1&min_price=1000&max_price=500000`

Only active products in active categories are returned.

### Public product detail

`GET /products/{id}`

### Admin create product

`POST /admin/products`

Role: `admin`

```json
{
  "category_id": 1,
  "name": "Mechanical Keyboard",
  "description": "Tenkeyless keyboard",
  "price": 1200000,
  "stock": 20,
  "image_url": "https://example.com/keyboard.jpg",
  "is_active": true
}
```

### Admin update product

`PUT /admin/products/{id}`

Role: `admin`

### Admin soft delete product

`DELETE /admin/products/{id}`

Role: `admin`

## Orders

### Create order

`POST /orders`

Header: `Authorization: Bearer ACCESS_TOKEN`

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

Backend tự lấy `unit_price` từ database, tự tính `subtotal` và `total_amount`.

### List orders

`GET /orders`

Header: `Authorization: Bearer ACCESS_TOKEN`

User chỉ thấy đơn hàng của chính mình. Admin thấy toàn bộ đơn hàng.

### Admin update order status

`PATCH /orders/{id}/status`

Role: `admin`

```json
{
  "status": "confirmed"
}
```

Luồng trạng thái hợp lệ:

- `pending -> confirmed`
- `pending -> cancelled`
- `confirmed -> shipping`
- `confirmed -> cancelled`
- `shipping -> completed`
