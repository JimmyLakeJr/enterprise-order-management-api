# Quy chuẩn chung của project

Tài liệu này là checklist kỹ thuật cho project `enterprise-order-management-api`.

## Trọng tâm

- Backend API là phần chính.
- Frontend chỉ là demo client đơn giản.
- Phạm vi phù hợp thực tập 2 tháng.
- Code rõ ràng, dễ đọc, dễ bảo trì, phù hợp intern/junior backend.

## Stack

- Backend: Golang 1.22+
- Framework: Echo v4
- Database: PostgreSQL
- Database access: SQL thuần với `pgx/v5` và `pgxpool`
- Không dùng ORM như GORM
- Authentication: JWT access token + refresh token
- JWT library: `github.com/golang-jwt/jwt/v5`
- Password hashing: `golang.org/x/crypto/bcrypt`
- Validation: `github.com/go-playground/validator/v10`
- Config: `.env`, `godotenv`, `os.Getenv`
- Testing: `testing` + `testify` khi cần assertion phức tạp
- Deploy local: Docker Compose
- Deploy demo production: Supabase/Neon, Render, Vercel

## Kiến trúc backend

- Handler: nhận request, parse request, gọi service, trả response.
- Service: xử lý business logic.
- Repository: thao tác database bằng SQL thuần.
- Model: ánh xạ dữ liệu database.
- DTO: định nghĩa request/response.
- Middleware: JWT auth, role authorization, CORS, logger, recovery.
- Config: load biến môi trường.
- Database: khởi tạo `pgxpool`, transaction, close connection.
- Util/pkg: JWT, password, response, pagination/hash helper.

## Quy tắc code

- Không viết SQL trong handler.
- Không viết business logic trong handler.
- Không hard-code secret/database URL/JWT secret khi deploy.
- Không trả `password_hash` ra response.
- Không expose lỗi database thô cho client.
- Không nối chuỗi SQL trực tiếp với input người dùng.
- Phải dùng parameterized query.
- Phải dùng transaction cho tạo đơn hàng.
- Phải validate input trước khi xử lý.
- Phải trả response JSON thống nhất.
- Phải có `.env.example`.
- Phải có README, Dockerfile, Docker Compose, API docs hoặc Postman collection.

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
