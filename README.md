# enterprise-order-management-api

Backend API cho hệ thống quản lý sản phẩm và đơn hàng trong doanh nghiệp, xây dựng bằng Golang, Echo v4 và PostgreSQL.

## 1. Mục tiêu của bước

Project tập trung vào backend API, phù hợp đồ án thực tập 2 tháng:

- RESTful API bằng Golang 1.22+ và Echo v4.
- PostgreSQL, SQL thuần với `pgx/v5` và `pgxpool`.
- Không dùng ORM.
- JWT access token và refresh token.
- Refresh token được hash trước khi lưu database.
- Logout revoke refresh token.
- Role authorization: `admin`, `user`.
- Product/category soft delete bằng `is_active = false`.
- Tạo order có transaction, kiểm tra stock và trừ stock.
- Response JSON thống nhất.
- Docker Compose để chạy local.

## 2. File/thư mục cần tạo

```text
cmd/api/main.go
internal/config
internal/database
internal/dto
internal/handler
internal/http
internal/middleware
internal/model
internal/pkg
internal/repository
internal/service
migrations/001_init.sql
docs/API.md
docs/ERD.md
docs/REQUIREMENTS.md
.env.example
Dockerfile
docker-compose.yml
```

## 3. Code hoàn chỉnh

Code nằm trực tiếp trong project. Các file quan trọng:

- `cmd/api/main.go`: chạy server.
- `internal/config/config.go`: load `.env` bằng `godotenv` và `os.Getenv`.
- `internal/database/postgres.go`: khởi tạo `pgxpool`.
- `internal/http/server.go`: khai báo route, middleware, dependency wiring.
- `internal/handler`: nhận request, validate, gọi service, trả response.
- `internal/service`: xử lý business logic.
- `internal/repository`: thao tác database bằng SQL thuần.
- `internal/pkg/response`: chuẩn hóa success/error/pagination response.
- `migrations/001_init.sql`: tạo bảng và seed role/admin/category.
- `docs/ANALYSIS.md`: phân tích yêu cầu hệ thống.
- `docs/ARCHITECTURE.md`: thiết kế kiến trúc tổng thể.
- `docs/DATABASE_DESIGN.md`: thiết kế cơ sở dữ liệu.

## 4. Giải thích ngắn gọn

Luồng xử lý:

```text
HTTP request -> Handler -> Service -> Repository -> PostgreSQL
```

Quy tắc tách lớp:

- Handler không viết SQL.
- Handler không xử lý business logic.
- Service chứa nghiệp vụ: kiểm tra quyền, order transaction, status transition.
- Repository chứa SQL thuần, dùng parameterized query.
- Model ánh xạ database.
- DTO định nghĩa request/response.

Response chuẩn:

```json
{
  "success": true,
  "message": "Success",
  "data": {}
}
```

Order transaction:

- User gửi `product_id` và `quantity`.
- Backend lấy giá từ database, không tin giá frontend.
- Backend kiểm tra product active và stock đủ.
- Backend tạo `orders`, `order_items`, trừ stock trong cùng transaction.
- Có lỗi thì rollback toàn bộ.

## 5. Cách chạy/test

Chạy bằng Docker Compose:

```bash
docker compose up --build
```

Chạy local không dùng Docker:

```bash
cp .env.example .env
go mod tidy
go run ./cmd/api
```

Chạy test:

```bash
go test ./...
```

Health check:

```bash
curl http://localhost:8080/health
```

Admin mặc định:

```text
email: admin@example.com
password: 123456
```

## 6. Ví dụ curl

Login admin:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"123456"}'
```

Tạo category:

```bash
curl -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -d '{"name":"Electronics","description":"Electronic devices","is_active":true}'
```

Tạo product:

```bash
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -d '{"category_id":1,"name":"Mechanical Keyboard","description":"Tenkeyless keyboard","price":1200000,"stock":20,"image_url":"https://example.com/keyboard.jpg","is_active":true}'
```

Xem product có phân trang/lọc:

```bash
curl "http://localhost:8080/api/v1/products?page=1&limit=10&search=keyboard&category_id=1"
```

Đăng ký user:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Nguyen Van A","email":"user@example.com","password":"123456"}'
```

Tạo order:

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -d '{"items":[{"product_id":1,"quantity":2}]}'
```

Admin cập nhật trạng thái order:

```bash
curl -X PATCH http://localhost:8080/api/v1/orders/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -d '{"status":"confirmed"}'
```

Logout:

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -d '{"refresh_token":"REFRESH_TOKEN"}'
```

## 7. Lỗi thường gặp và cách sửa

`connection refused`

- PostgreSQL chưa chạy.
- Kiểm tra `docker compose up` hoặc `DATABASE_URL`.

`Invalid access token`

- Thiếu header `Authorization: Bearer ACCESS_TOKEN`.
- Access token hết hạn, gọi `/auth/refresh`.

`Refresh token was revoked or expired`

- Refresh token đã logout, đã refresh trước đó, hoặc hết hạn.
- Login lại để nhận token mới.

`You do not have permission`

- Tài khoản không phải admin nhưng gọi API `/admin/*`.

`Product does not have enough stock`

- Số lượng đặt lớn hơn tồn kho.
- Giảm quantity hoặc tăng stock sản phẩm.

`Invalid order status transition`

- Chuyển trạng thái sai luồng.
- Luồng đúng: `pending -> confirmed/cancelled`, `confirmed -> shipping/cancelled`, `shipping -> completed`.

## 8. Deploy demo

Gợi ý deploy:

- Database: Supabase hoặc Neon PostgreSQL.
- Backend: Render.
- Frontend demo: Vercel.

Khi deploy backend, cần cấu hình biến môi trường:

```text
PORT
DATABASE_URL
JWT_ACCESS_SECRET
JWT_REFRESH_SECRET
FRONTEND_URL
ACCESS_TOKEN_MINUTES
REFRESH_TOKEN_HOURS
```

## 9. Tài liệu báo cáo

- Phân tích yêu cầu: `docs/ANALYSIS.md`
- Kiến trúc tổng thể: `docs/ARCHITECTURE.md`
- Thiết kế cơ sở dữ liệu: `docs/DATABASE_DESIGN.md`
- ERD: `docs/ERD.md`
- API docs: `docs/API.md`
- Quy chuẩn project: `docs/PROJECT_STANDARDS.md`
