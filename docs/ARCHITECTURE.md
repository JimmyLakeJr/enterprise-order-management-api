# Thiết kế kiến trúc tổng thể hệ thống

Project: **enterprise-order-management-api**

Đề tài: **Phát triển backend API cho hệ thống quản lý sản phẩm và đơn hàng trong doanh nghiệp sử dụng Golang**

## 1. Mục tiêu

Mục tiêu của bước này là thiết kế kiến trúc tổng thể cho hệ thống quản lý sản phẩm và đơn hàng, đảm bảo:

- Backend là phần trọng tâm của đồ án.
- Kiến trúc rõ ràng, dễ bảo trì, phù hợp sinh viên thực tập backend.
- Tuân thủ stack đã thống nhất: Go, Echo v4, PostgreSQL, pgxpool, JWT, Docker.
- Tách rõ trách nhiệm giữa các lớp: handler, service, repository, model, dto, middleware.
- Không viết SQL trong handler.
- Không viết business logic trong handler.
- Dễ mở rộng thêm frontend demo React + Vite và deploy lên môi trường demo.

## 2. Cấu trúc file/thư mục

### Backend

```text
enterprise-order-management-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── dto/
│   │   ├── auth.go
│   │   ├── order.go
│   │   └── product.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── category_handler.go
│   │   ├── order_handler.go
│   │   ├── product_handler.go
│   │   └── user_handler.go
│   ├── http/
│   │   └── server.go
│   ├── middleware/
│   │   └── auth.go
│   ├── model/
│   │   └── models.go
│   ├── pkg/
│   │   ├── apperror/
│   │   ├── hasher/
│   │   ├── password/
│   │   ├── response/
│   │   ├── token/
│   │   └── validator/
│   ├── repository/
│   │   ├── category_repository.go
│   │   ├── order_repository.go
│   │   ├── product_repository.go
│   │   └── user_repository.go
│   └── service/
│       ├── auth_service.go
│       ├── category_service.go
│       ├── order_service.go
│       ├── product_service.go
│       └── user_service.go
├── migrations/
│   └── 001_init.sql
├── docs/
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

### Frontend demo

Frontend chỉ dùng để demo client đơn giản, không phải trọng tâm của đồ án.

```text
frontend/
├── src/
│   ├── api/
│   │   ├── authApi.js
│   │   ├── categoryApi.js
│   │   ├── productApi.js
│   │   └── orderApi.js
│   ├── components/
│   ├── pages/
│   │   ├── LoginPage.jsx
│   │   ├── ProductListPage.jsx
│   │   ├── ProductDetailPage.jsx
│   │   ├── CartPage.jsx
│   │   ├── MyOrdersPage.jsx
│   │   └── AdminDashboardPage.jsx
│   ├── routes/
│   ├── stores/
│   ├── utils/
│   ├── App.jsx
│   └── main.jsx
├── .env.example
├── package.json
└── vite.config.js
```

## 3. Code hoàn chỉnh

Bước này là bước thiết kế kiến trúc nên không thêm source code xử lý nghiệp vụ mới. Các file code hiện tại đã được tổ chức theo kiến trúc này.

Các file đại diện:

- Backend entrypoint: `cmd/api/main.go`
- Router/server: `internal/http/server.go`
- Config: `internal/config/config.go`
- Database connection: `internal/database/postgres.go`
- Handler: `internal/handler`
- Service: `internal/service`
- Repository: `internal/repository`
- Model: `internal/model`
- DTO: `internal/dto`
- Middleware: `internal/middleware`
- Utility package: `internal/pkg`

## 4. Giải thích ngắn gọn

### 4.1. Kiến trúc tổng quan FE - BE - DB

Hệ thống được chia thành 3 phần chính:

```text
React + Vite Frontend
        |
        | HTTP/JSON REST API
        v
Go + Echo Backend API
        |
        | SQL through pgxpool
        v
PostgreSQL Database
```

Vai trò từng phần:

- **Frontend demo**: giao diện đơn giản để đăng nhập, xem sản phẩm, tạo đơn hàng và demo nghiệp vụ.
- **Backend API**: xử lý xác thực, phân quyền, nghiệp vụ sản phẩm, danh mục, người dùng và đơn hàng.
- **PostgreSQL**: lưu dữ liệu roles, users, refresh_tokens, categories, products, orders, order_items.

Khi deploy demo:

```text
Vercel Frontend
        |
        v
Render Backend API
        |
        v
Supabase/Neon PostgreSQL
```

### 4.2. Luồng request từ frontend đến database

Ví dụ User tạo đơn hàng:

```text
1. User bấm "Tạo đơn hàng" trên frontend.
2. Frontend gửi POST /api/v1/orders kèm access token.
3. Echo middleware kiểm tra JWT.
4. OrderHandler nhận request và validate dữ liệu.
5. OrderHandler gọi OrderService.
6. OrderService xử lý nghiệp vụ:
   - kiểm tra order không rỗng
   - lấy product từ database
   - kiểm tra product active
   - kiểm tra stock đủ
   - tính unit_price, subtotal, total_amount
   - mở transaction
7. OrderRepository thực hiện SQL:
   - tạo orders
   - tạo order_items
   - trừ stock products
8. PostgreSQL commit transaction.
9. Backend trả JSON response cho frontend.
10. Frontend hiển thị kết quả cho User.
```

### 4.3. Vai trò của từng thành phần backend

#### Handler

Handler là lớp tiếp nhận HTTP request.

Nhiệm vụ:

- Parse path param, query param, request body.
- Gọi validator để kiểm tra input.
- Lấy thông tin user từ context nếu API cần đăng nhập.
- Gọi service tương ứng.
- Trả JSON response.

Handler không chứa SQL và không xử lý business logic phức tạp.

#### Service

Service là lớp xử lý nghiệp vụ chính.

Nhiệm vụ:

- Kiểm tra business rules.
- Điều phối nhiều repository nếu cần.
- Xử lý transaction.
- Tính toán nghiệp vụ như tổng tiền đơn hàng.
- Kiểm tra quyền nghiệp vụ nếu cần.
- Trả dữ liệu đã xử lý cho handler.

Ví dụ:

- AuthService xử lý login, refresh token, logout.
- ProductService kiểm tra category active trước khi tạo product.
- OrderService xử lý transaction tạo order và trừ stock.

#### Repository

Repository là lớp thao tác database.

Nhiệm vụ:

- Chứa SQL thuần.
- Dùng `pgxpool` hoặc `pgx.Tx`.
- Dùng parameterized query.
- Map dữ liệu database sang model.

Repository không xử lý nghiệp vụ như kiểm tra role, tính tổng tiền hoặc quyết định luồng trạng thái đơn hàng.

#### Model

Model ánh xạ dữ liệu trong database.

Ví dụ:

- User
- Role
- Category
- Product
- Order
- OrderItem

Model thường gần với cấu trúc bảng database.

#### DTO

DTO định nghĩa dữ liệu request và response.

Ví dụ:

- RegisterRequest
- LoginRequest
- ProductRequest
- ProductResponse
- CreateOrderRequest
- OrderResponse

DTO giúp không expose trực tiếp model database ra client. Ví dụ `password_hash` nằm trong model User nhưng không xuất hiện trong UserResponse.

#### Middleware

Middleware xử lý logic trước khi request vào handler.

Nhiệm vụ:

- CORS.
- Logger.
- Recovery.
- JWT authentication.
- Role-based authorization.

Ví dụ:

- API tạo order cần JWT middleware.
- API admin cần thêm role middleware.

#### Config

Config chịu trách nhiệm load biến môi trường.

Nguồn cấu hình:

- `.env`
- biến môi trường thật khi deploy
- `os.Getenv`

Các biến quan trọng:

- `PORT`
- `DATABASE_URL`
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET`
- `FRONTEND_URL`
- `ACCESS_TOKEN_MINUTES`
- `REFRESH_TOKEN_HOURS`

#### Database

Database package chịu trách nhiệm khởi tạo kết nối PostgreSQL bằng `pgxpool`.

Nhiệm vụ:

- Parse database URL.
- Tạo connection pool.
- Ping database.
- Close connection khi server shutdown.

#### Util / pkg

`internal/pkg` chứa các helper dùng chung.

Ví dụ:

- `password`: hash/check password bằng bcrypt.
- `token`: generate/parse JWT.
- `response`: chuẩn hóa JSON response.
- `apperror`: định nghĩa lỗi ứng dụng.
- `validator`: custom validator cho Echo.
- `hasher`: hash refresh token trước khi lưu database.

### 4.4. Lý do không viết SQL trong handler

Không viết SQL trong handler vì:

- Handler sẽ bị quá nhiều trách nhiệm.
- Code khó đọc và khó bảo trì.
- Khó tái sử dụng query ở nhiều nơi.
- Khó test riêng phần xử lý database.
- Dễ trộn lẫn HTTP logic với database logic.
- Khi thay đổi schema database, phải sửa nhiều handler.

Cách đúng:

```text
Handler -> Service -> Repository -> Database
```

Handler chỉ gọi service, repository mới chứa SQL.

### 4.5. Lý do không viết business logic trong handler

Không viết business logic trong handler vì:

- Handler chỉ nên phụ trách HTTP layer.
- Business logic có thể được dùng lại ở nhiều API khác nhau.
- Nếu logic nằm trong handler, code sẽ khó test.
- Handler sẽ phình to khi nghiệp vụ tăng.
- Dễ gây lỗi khi nhiều handler xử lý cùng một rule nhưng viết khác nhau.

Ví dụ nghiệp vụ tạo order không nên nằm trong handler:

- Kiểm tra product active.
- Kiểm tra stock.
- Tính tổng tiền.
- Tạo transaction.
- Trừ stock.
- Rollback khi lỗi.

Các logic trên phải nằm trong service.

## 5. Cách chạy/test

Kiểm tra code backend:

```bash
go test ./...
```

Chạy backend local:

```bash
go run ./cmd/api
```

Chạy bằng Docker Compose:

```bash
docker compose up --build
```

Health check:

```bash
curl http://localhost:8080/health
```

Ví dụ login admin:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"123456"}'
```

## 6. Lỗi thường gặp

### Handler chứa quá nhiều logic

Dấu hiệu:

- Handler dài.
- Handler có nhiều `if else` nghiệp vụ.
- Handler tự tính tổng tiền hoặc tự xử lý transaction.

Cách sửa:

- Chuyển logic sang service.
- Handler chỉ parse request, validate, gọi service và trả response.

### SQL xuất hiện trong handler

Dấu hiệu:

- Trong file handler có `SELECT`, `INSERT`, `UPDATE`, `DELETE`.

Cách sửa:

- Chuyển SQL sang repository.
- Handler gọi service, service gọi repository.

### Service chứa SQL

Dấu hiệu:

- Trong service có query SQL dài.

Cách sửa:

- Chuyển SQL sang repository.
- Service chỉ điều phối nghiệp vụ.

### Frontend gọi sai endpoint

Dấu hiệu:

- API trả 404 hoặc CORS error.

Cách sửa:

- Kiểm tra base URL frontend.
- Kiểm tra `FRONTEND_URL` trong backend.
- Kiểm tra route trong `internal/http/server.go`.

### CORS lỗi khi deploy

Dấu hiệu:

- Frontend Vercel không gọi được backend Render.

Cách sửa:

- Cập nhật `FRONTEND_URL` trong môi trường deploy backend.
- Không hard-code localhost khi deploy production.

## 7. Checklist hoàn thành

### Kiến trúc tổng quan

- [x] Có mô hình FE - BE - DB.
- [x] Có định hướng deploy Vercel - Render - Supabase/Neon.
- [x] Backend là trọng tâm.
- [x] Frontend chỉ là demo client.

### Backend structure

- [x] Có `cmd/api`.
- [x] Có `internal/config`.
- [x] Có `internal/database`.
- [x] Có `internal/handler`.
- [x] Có `internal/service`.
- [x] Có `internal/repository`.
- [x] Có `internal/model`.
- [x] Có `internal/dto`.
- [x] Có `internal/middleware`.
- [x] Có `internal/pkg`.

### Quy tắc tách lớp

- [x] Handler không viết SQL.
- [x] Handler không chứa business logic.
- [x] Service xử lý business logic.
- [x] Repository chứa SQL thuần.
- [x] DTO tách khỏi model database.
- [x] Middleware xử lý JWT và role authorization.

### Công nghệ

- [x] Backend dùng Go + Echo v4.
- [x] Database dùng PostgreSQL.
- [x] Database access dùng SQL thuần với pgxpool.
- [x] Auth dùng JWT access token + refresh token.
- [x] Deploy local bằng Docker Compose.
- [x] Deploy demo phù hợp Render, Vercel, Supabase/Neon.

### Sơ đồ flow dạng text

```text
Frontend React + Vite
    |
    | HTTP Request JSON
    v
Echo Router
    |
    v
Middleware
    |
    v
Handler
    |
    v
Service
    |
    v
Repository
    |
    v
PostgreSQL
    |
    v
Repository
    |
    v
Service
    |
    v
Handler
    |
    | JSON Response
    v
Frontend React + Vite
```
