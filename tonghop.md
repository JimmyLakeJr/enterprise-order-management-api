# TỔNG HỢP PROJECT ENTERPRISE ORDER MANAGEMENT

> Snapshot kiểm tra repo tại thời điểm tổng hợp: 2026-06-24.  
> Source of truth backend runtime: `cmd/api`, `internal/*`, `migrations/001_init.sql`, root `Dockerfile`, root `docker-compose.yml`.  
> Thư mục `backend/` tồn tại nhưng không được dùng làm runtime nghiệp vụ chính của project hiện tại.

## 1. Giới thiệu project

- Project này là hệ thống quản lý sản phẩm, tồn kho cơ bản, người dùng và vòng đời đơn hàng.
- Đối tượng sử dụng chính:
  - Guest: xem sản phẩm, xem chi tiết, thêm giỏ hàng local.
  - User: đăng ký, đăng nhập, xem/cập nhật hồ sơ tối thiểu, tạo đơn, xem đơn của mình.
  - Admin: quản trị danh mục, sản phẩm, đơn hàng và người dùng.
- Mục đích sử dụng:
  - Demo backend API Go/Echo + PostgreSQL.
  - Demo frontend React/Vite gọi API thật.
  - Demo luồng quản trị cơ bản cho sản phẩm, user và order lifecycle.
- Đây **không phải** hệ thống e-commerce hoàn chỉnh.
- Backend đóng vai trò:
  - Cung cấp REST API thật dưới prefix `/api/v1`.
  - Xử lý auth JWT access/refresh token.
  - Validate dữ liệu đầu vào.
  - Truy cập PostgreSQL.
  - Tạo đơn hàng theo transaction và trừ tồn kho.
- Frontend đóng vai trò:
  - Giao diện demo public/user/admin.
  - Gọi API thật bằng Axios.
  - Quản lý auth state và cart state ở phía client.
- Các role hiện có:
  - `Guest`
  - `user`
  - `admin`
- Các role chưa có:
  - `staff`
  - `manager`
  - Không thấy role backend nào khác ngoài `admin` và `user`.

## 2. Công nghệ sử dụng

| Thành phần | Công nghệ | Vị trí trong repo | Ghi chú |
|---|---|---|---|
| Backend | Go + Echo | `cmd/api`, `internal/http`, `internal/handler`, `internal/service` | Runtime API chính |
| Database | PostgreSQL | `migrations/001_init.sql`, `docker-compose.yml` | Schema và seed cơ bản |
| Database access | `pgx/v5`, `pgxpool` | `internal/database/postgres.go`, `internal/repository/*` | Không dùng ORM |
| Auth | JWT access token + refresh token | `internal/pkg/token/jwt.go`, `internal/service/auth_service.go` | Refresh token lưu hash trong DB |
| Password hash | bcrypt | `internal/pkg/password` | Dùng cho user password |
| Validation | `go-playground/validator/v10` | `internal/pkg/validator/validator.go`, DTO tags | Validate request body |
| Frontend | React + Vite | `frontend/` | Mã hiện tại là JSX/JavaScript, không phải TypeScript |
| Router | `react-router-dom` | `frontend/src/routes/*` | Public/User/Admin route |
| API client | Axios | `frontend/src/api/*` | Có interceptor refresh token |
| State phía client | React Context | `frontend/src/contexts/*` | `AuthContext`, `CartContext`, `ConfirmProvider` |
| Docker build | Root `Dockerfile` | `Dockerfile` | Chỉ build backend Go API |
| Docker runtime | Root `docker-compose.yml` | `docker-compose.yml` | Chỉ có `postgres` và `api`, không có service frontend |
| Migration | SQL thuần | `migrations/001_init.sql` | Tự chạy qua `docker-entrypoint-initdb.d` |
| Backend testing | `go test`, `go vet` | root repo | Đã chạy thực tế |
| Frontend checking | `npm run lint`, `npm run build` | `frontend/package.json` | Đã chạy thực tế |

## 3. Cấu trúc thư mục source code

### Vai trò các thư mục chính

- `cmd/api`: entrypoint khởi động app, load config, kết nối DB, boot Echo server.
- `internal/config`: load biến môi trường và default config.
- `internal/database`: khởi tạo `pgxpool` tới PostgreSQL.
- `internal/http`: wiring dependency, middleware, route registration.
- `internal/handler`: nhận HTTP request, bind/validate, gọi service, trả JSON response.
- `internal/service`: nghiệp vụ chính.
- `internal/repository`: SQL query và truy cập database.
- `internal/model`: model domain/backend.
- `internal/dto`: request/response DTO.
- `internal/middleware`: JWT auth, role guard, logger, recovery, CORS.
- `internal/pkg`: utility dùng chung như error, response format, validator, token, password hash.
- `migrations`: migration SQL thật của runtime root repo.
- `frontend/src/api`: API client frontend.
- `frontend/src/pages`: các màn hình public/user/admin.
- `frontend/src/components`: component UI dùng lại.
- `frontend/src/contexts`: auth/cart/confirm context.
- `frontend/src/routes`: route guard và route map.
- `frontend/src/styles`: CSS cho global/layout/page/component/admin.
- `backend/`: skeleton/module riêng cũ; không phải runtime nghiệp vụ chính theo repo hiện tại.

### Cây thư mục tổng quát

```text
repo/
├── cmd/
│   └── api/
├── internal/
│   ├── config/
│   ├── database/
│   ├── dto/
│   ├── handler/
│   ├── http/
│   ├── middleware/
│   ├── model/
│   ├── pkg/
│   ├── repository/
│   └── service/
├── migrations/
│   └── 001_init.sql
├── frontend/
│   ├── package.json
│   └── src/
│       ├── api/
│       ├── components/
│       ├── constants/
│       ├── contexts/
│       ├── hooks/
│       ├── layouts/
│       ├── pages/
│       ├── routes/
│       ├── styles/
│       └── utils/
├── docs/
├── backend/
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## 4. Database code nằm ở đâu?

### Runtime chính và skeleton

- Runtime backend chính:
  - `cmd/api/main.go`
  - `internal/config/config.go`
  - `internal/database/postgres.go`
  - `internal/repository/*.go`
  - `migrations/001_init.sql`
  - root `docker-compose.yml`
- Skeleton không dùng làm API nghiệp vụ chính:
  - `backend/`

### Thành phần database cụ thể

- File SQL migration chính: `migrations/001_init.sql`
- Migration bổ sung media hồ sơ local: `migrations/003_profile_media.sql`
- File tạo bảng:
  - `migrations/001_init.sql`
  - `migrations/003_profile_media.sql` thêm `users.profile_video_url`
- Seed role/admin:
  - `migrations/001_init.sql`
  - Seed `roles`: `admin`, `user`
  - Seed user `admin@example.com`
- Repository SQL:
  - `internal/repository/user_repository.go`
  - `internal/repository/category_repository.go`
  - `internal/repository/product_repository.go`
  - `internal/repository/order_repository.go`
- Model database/domain:
  - `internal/model/models.go`
- DTO request/response:
  - `internal/dto/*.go`
- Cấu hình kết nối DB:
  - `internal/config/config.go`
  - `internal/database/postgres.go`
- Docker Compose cấu hình PostgreSQL:
  - root `docker-compose.yml`
- Thư mục lưu file local:
  - `uploads/`
  - Ảnh sản phẩm: `uploads/products/images/`
  - Avatar hồ sơ: `uploads/profile/avatars/`
  - Video hồ sơ: `uploads/profile/videos/`

### Thông số DB local theo repo

| Hạng mục | Giá trị |
|---|---|
| Database name | `enterprise_order_management` |
| DB user | `postgres` |
| DB password | `postgres` |
| DB port | `5432` |
| Docker service name | `postgres` |
| Container name | `enterprise-order-postgres` |
| Database volume | `postgres_data` |
| Backend local DB host | `localhost` |
| Backend trong Docker DB host | `postgres` |

### Ghi chú quan trọng về seed admin

- Migration có seed `admin@example.com`.
- Migration chỉ seed **bcrypt hash** ở cột `password_hash`, không seed plaintext password.
- README/docs cũ có chỗ ghi password local là `Admin@123`.
- Password seed đúng cho `admin@example.com` trong repo là `123456`.
- Kiểm tra runtime hiện tại xác nhận `vu@gmail.com / 123456` đăng nhập được và tài khoản này đang có role `admin`.
- Kết luận thực tế:
  - Email admin seed có tồn tại.
  - Plaintext password hiện tại **không thể xác nhận chỉ từ code SQL**, vì migration chỉ chứa bcrypt hash.
  - Khả năng cao volume PostgreSQL hiện tại là volume cũ đã có dữ liệu trước đó hoặc password đã khác so với tài liệu.
  - Credential admin dùng được trên volume hiện tại là dữ liệu runtime hiện có, không phải credential seed sạch từ migration.
- Nếu cần môi trường sạch đúng seed, nên reset DB bằng:

```bash
docker compose down -v
docker compose up -d postgres api
```

## 5. Thiết kế database hiện tại

### Các bảng thực tế đang có trong PostgreSQL runtime

Kiểm tra trực tiếp bằng `SELECT tablename FROM pg_tables WHERE schemaname = 'public'` cho thấy chỉ có 7 bảng:

- `categories`
- `order_items`
- `oauth_accounts`
- `orders`
- `products`
- `refresh_tokens`
- `roles`
- `users`

### Bảng chi tiết

| Tên bảng | Mục đích | Các cột chính | Khóa chính | Khóa ngoại | Ràng buộc quan trọng | Ghi chú nghiệp vụ |
|---|---|---|---|---|---|---|
| `roles` | Danh sách role hệ thống | `id`, `name`, `created_at` | `id` | Không | `name` unique | Chỉ thấy `admin`, `user` |
| `users` | Tài khoản người dùng | `id`, `full_name`, `email`, `password_hash`, `avatar_url`, `profile_video_url`, `role_id`, `is_active`, `created_at`, `updated_at` | `id` | `role_id -> roles.id` | `email` unique, `is_active` default true | Soft disable user bằng `is_active = false`; `avatar_url` và `profile_video_url` đang được dùng cho media hồ sơ lưu local filesystem |
| `refresh_tokens` | Lưu hash refresh token | `id`, `user_id`, `token_hash`, `expires_at`, `revoked_at`, `created_at` | `id` | `user_id -> users.id` | `token_hash` unique | Logout/refresh dùng revoke token cũ |
| `oauth_accounts` | Liên kết tài khoản OAuth ngoài với user nội bộ | `id`, `user_id`, `provider`, `provider_user_id`, `email`, `avatar_url`, `created_at`, `updated_at` | `id` | `user_id -> users.id` | `UNIQUE(provider, provider_user_id)`, `UNIQUE(provider, email)` | Hiện dùng cho Google OAuth; không lưu Google password, không tự nâng role admin |
| `categories` | Danh mục sản phẩm | `id`, `name`, `description`, `is_active`, `created_at`, `updated_at` | `id` | Không | `name` unique | Soft delete qua `is_active` |
| `products` | Sản phẩm và tồn kho | `id`, `category_id`, `name`, `description`, `price`, `stock`, `image_url`, `is_active`, `created_at`, `updated_at` | `id` | `category_id -> categories.id` | `price >= 0`, `stock >= 0` | `image_url` hiện hỗ trợ cả URL công khai và URL local `/uploads/...`; chưa có bảng media riêng |
| `orders` | Đơn hàng | `id`, `user_id`, `total_amount`, `status`, `created_at`, `updated_at` | `id` | `user_id -> users.id` | `total_amount >= 0`, status check enum | Status: `pending`, `confirmed`, `shipping`, `completed`, `cancelled` |
| `order_items` | Dòng sản phẩm của đơn | `id`, `order_id`, `product_id`, `quantity`, `unit_price`, `subtotal`, `created_at` | `id` | `order_id -> orders.id`, `product_id -> products.id` | `quantity > 0`, `unit_price >= 0`, `subtotal >= 0` | Snapshot đơn giá tại thời điểm tạo đơn |

### Các bảng bắt buộc kiểm tra

| Bảng | Trạng thái |
|---|---|
| `roles` | Có |
| `users` | Có |
| `refresh_tokens` | Có |
| `categories` | Có |
| `products` | Có |
| `orders` | Có |
| `order_items` | Có |
| `oauth_accounts` | Có |
| `payments` | Chưa implement |
| `vouchers` | Chưa implement |
| `uploads` | Chưa implement |
| `media` | Chưa implement |
| `email_logs` | Chưa implement |
| `shipping_addresses` | Chưa implement |
| `inventory_logs` | Chưa implement |

## 6. Cách tải, cài đặt và chuẩn bị môi trường

### 6.1. Công cụ cần cài

```bash
go version
node -v
npm -v
docker --version
docker compose version
psql --version
git --version
```

### 6.2. Clone project

Không có repo URL chắc chắn trong source hiện tại, nên dùng cách an toàn:

```bash
cd <thu-muc-project-hien-tai>
```

Nếu có remote của riêng bạn thì clone theo mẫu:

```bash
git clone <repo-url>
cd <repo-folder>
```

### 6.3. Cấu hình env backend

Repo hiện tại **không còn** `.env.example` trong worktree đang kiểm tra.  
Biến môi trường thực tế được đọc từ `internal/config/config.go`, và file `.env` hiện có đang dùng các biến sau:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/enterprise_order_management?sslmode=disable
JWT_ACCESS_SECRET=change-this-access-secret
JWT_REFRESH_SECRET=change-this-refresh-secret
FRONTEND_URL=http://localhost:5173
FRONTEND_AUTH_CALLBACK_URL=http://localhost:5173/auth/google/callback
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
OAUTH_STATE_SECRET=change-this-google-oauth-state-secret
ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_HOURS=168
```

- File env backend nên tạo: `.env`
- Biến kết nối DB dùng: `DATABASE_URL`
- Nếu backend chạy local:
  - DB host: `localhost`
- Nếu backend chạy trong Docker:
  - DB host: `postgres`

### 6.4. Cấu hình env frontend

File có sẵn:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

- File env frontend: `frontend/.env.example`
- Khi chạy local, backend URL đúng là:

```text
http://localhost:8080/api/v1
```

- Google frontend callback local:

```text
http://localhost:5173/auth/google/callback
```

## 7. Cách chạy database trên môi trường SQL/PostgreSQL

### 7.1. Chạy PostgreSQL bằng Docker Compose

```bash
docker compose up -d postgres
```

### 7.2. Kiểm tra container

```bash
docker compose ps
docker compose logs postgres
```

### 7.3. Vào psql trong container

```bash
docker compose exec postgres psql -U postgres -d enterprise_order_management
```

### 7.4. Kiểm tra bảng

```sql
SELECT tablename
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY tablename;
```

### 7.5. Xem dữ liệu từng bảng

```sql
SELECT * FROM roles ORDER BY id;

SELECT id, full_name, email, role_id, is_active, created_at
FROM users
ORDER BY id;

SELECT id, user_id, expires_at, revoked_at, created_at
FROM refresh_tokens
ORDER BY id DESC;

SELECT * FROM categories ORDER BY id;

SELECT id, category_id, name, price, stock, image_url, is_active
FROM products
ORDER BY id;

SELECT * FROM orders ORDER BY id DESC;

SELECT * FROM order_items ORDER BY order_id, id;
```

### 7.6. Query kiểm tra order đầy đủ

```sql
SELECT
  o.id AS order_id,
  o.status,
  o.total_amount,
  o.created_at,
  u.id AS user_id,
  u.full_name,
  u.email,
  oi.id AS order_item_id,
  oi.product_id,
  p.name AS product_name,
  oi.quantity,
  oi.unit_price,
  oi.subtotal
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id
ORDER BY o.id DESC, oi.id ASC;
```

### 7.7. Kiểm tra module chưa có bảng

```sql
SELECT to_regclass('public.payments') AS payments;
SELECT to_regclass('public.vouchers') AS vouchers;
SELECT to_regclass('public.uploads') AS uploads;
SELECT to_regclass('public.media') AS media;
SELECT to_regclass('public.email_logs') AS email_logs;
SELECT to_regclass('public.shipping_addresses') AS shipping_addresses;
SELECT to_regclass('public.inventory_logs') AS inventory_logs;
```

Kết quả kiểm tra runtime hiện tại: tất cả đều `NULL`, tức là **chưa implement**.

## 8. Cách chạy backend

### 8.1. Chạy backend local bằng Go

Entrypoint thật:

```bash
go run ./cmd/api
```

### 8.2. Chạy backend bằng Docker Compose

```bash
docker compose up --build api
```

### 8.3. Kiểm tra backend

```bash
curl http://localhost:8080/health
```

Response kỳ vọng:

```json
{
  "success": true,
  "message": "Success",
  "data": {
    "status": "ok"
  }
}
```

## 9. Cách chạy frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend local mặc định:

```text
http://localhost:5173
```

Kiểm tra frontend gọi đúng backend:

- Xem tab Network của trình duyệt.
- Kiểm tra `VITE_API_BASE_URL`.
- Kiểm tra backend có cho phép CORS từ `FRONTEND_URL=http://localhost:5173`.

Luồng upload media hiện tại:

- Admin có thể tải ảnh sản phẩm lên local storage qua API `POST /api/v1/products/upload-image`, sau đó lưu URL trả về vào trường `image_url` của sản phẩm.
- User/admin có thể tải avatar qua `POST /api/v1/users/me/avatar`.
- User/admin có thể tải video hồ sơ ngắn qua `POST /api/v1/users/me/profile-video`.
- Backend đang public static file tại `/uploads/*`, nên các URL như `http://localhost:8080/uploads/...` hiển thị được ở trang chủ, danh sách sản phẩm, trang chi tiết và hồ sơ sau khi record đã được lưu lại.

### Cập nhật UI/UX frontend mới nhất

- Frontend đã được tối ưu UI theo hướng enterprise dashboard, giữ nguyên API và route hiện có.
- Đã tạo design token riêng tại `frontend/src/styles/tokens.css`.
- Đã chuẩn hóa lại các file style chính:
  - `frontend/src/styles/global.css`
  - `frontend/src/styles/layout.css`
  - `frontend/src/styles/components.css`
  - `frontend/src/styles/pages.css`
  - `frontend/src/styles/admin.css`
- Các khu vực UI đã được làm mới:
  - Public header/navigation
  - Product card/list/detail
  - Login/Register
  - Profile page
  - Admin layout/sidebar/header
  - Admin dashboard
- Tài liệu chi tiết về đợt tối ưu UI:
  - `docs/UI_UX_OPTIMIZATION_AUDIT.md`
  - `docs/UI_UX_OPTIMIZATION_REPORT.md`

## 10. Cách chạy toàn bộ bằng Docker

Root `docker-compose.yml` hiện tại **chỉ** start:

- `postgres`
- `api`

Frontend **không** có service Docker trong compose hiện tại.

```bash
docker compose up --build
```

Khi chạy lệnh trên:

- Service backend: `api`
- Service database: `postgres`
- Port backend: `8080`
- Port database: `5432`
- Frontend: vẫn cần chạy local bằng `npm run dev`

Xem logs:

```bash
docker compose logs -f
docker compose logs -f postgres
docker compose logs -f api
```

Dừng:

```bash
docker compose down
```

Reset database:

```bash
docker compose down -v
docker compose up --build
```

> Cảnh báo: `docker compose down -v` sẽ xóa volume `postgres_data` và làm mất toàn bộ dữ liệu hiện có.

## 11. Toàn bộ API endpoint và link test chức năng

| Nhóm chức năng | Method | URL | Auth required | Role | Request body/query | Response chính | Link test local |
|---|---|---|---|---|---|---|---|
| Health | `GET` | `/health` | Không | Guest | Không | `{status:"ok"}` | `GET http://localhost:8080/health` |
| Auth | `POST` | `/api/v1/auth/register` | Không | Guest | `name,email,password` | access token, refresh token, user | `POST http://localhost:8080/api/v1/auth/register` |
| Auth | `POST` | `/api/v1/auth/login` | Không | Guest | `email,password` | access token, refresh token, user | `POST http://localhost:8080/api/v1/auth/login` |
| Auth | `GET` | `/api/v1/auth/google/login` | Không | Guest | redirect browser | redirect sang Google hoặc callback frontend với lỗi cấu hình | `GET http://localhost:8080/api/v1/auth/google/login` |
| Auth | `GET` | `/api/v1/auth/google/callback` | Không | Guest | query `code,state` từ Google | redirect về frontend callback với `access_token`/`refresh_token` trong hash fragment hoặc lỗi | `GET http://localhost:8080/api/v1/auth/google/callback` |
| Auth | `POST` | `/api/v1/auth/refresh-token` | Không | Guest | `refresh_token` | cặp token mới + user | `POST http://localhost:8080/api/v1/auth/refresh-token` |
| Auth | `POST` | `/api/v1/auth/logout` | Có | `user/admin` | `refresh_token` | message logout | `POST http://localhost:8080/api/v1/auth/logout` |
| Profile | `GET` | `/api/v1/auth/me` | Có | `user/admin` | Không | user hiện tại | `GET http://localhost:8080/api/v1/auth/me` |
| Profile | `PUT` | `/api/v1/users/me` | Có | `user/admin` | `name` | user đã cập nhật | `PUT http://localhost:8080/api/v1/users/me` |
| Profile | `POST` | `/api/v1/users/me/avatar` | Có | `user/admin` | `multipart/form-data`, field `file` | `{ url }` + cập nhật `avatar_url` của user | `POST http://localhost:8080/api/v1/users/me/avatar` |
| Profile | `POST` | `/api/v1/users/me/profile-video` | Có | `user/admin` | `multipart/form-data`, field `file` | `{ url }` + cập nhật `profile_video_url` của user | `POST http://localhost:8080/api/v1/users/me/profile-video` |
| Category | `GET` | `/api/v1/categories` | Không | Guest | Không | danh sách category active | `GET http://localhost:8080/api/v1/categories` |
| Category | `GET` | `/api/v1/categories/:id` | Không | Guest | param `id` | category active detail | `GET http://localhost:8080/api/v1/categories/1` |
| Category admin | `POST` | `/api/v1/categories` | Có | `admin` | `name,description,is_active` | category mới | `POST http://localhost:8080/api/v1/categories` |
| Category admin | `PUT` | `/api/v1/categories/:id` | Có | `admin` | `name,description,is_active` | category cập nhật | `PUT http://localhost:8080/api/v1/categories/1` |
| Category admin | `DELETE` | `/api/v1/categories/:id` | Có | `admin` | param `id` | message xóa mềm | `DELETE http://localhost:8080/api/v1/categories/1` |
| Category admin | `GET` | `/api/v1/admin/categories` | Có | `admin` | query `status=all|active|inactive` | danh sách admin | `GET http://localhost:8080/api/v1/admin/categories?status=all` |
| Category admin | `PUT` | `/api/v1/admin/categories/:id/restore` | Có | `admin` | param `id` | category restored | `PUT http://localhost:8080/api/v1/admin/categories/1/restore` |
| Product | `GET` | `/api/v1/products` | Không | Guest | `page,limit,keyword,category_id,min_price,max_price` | danh sách + meta | `GET http://localhost:8080/api/v1/products` |
| Product | `GET` | `/api/v1/products/:id` | Không | Guest | param `id` | product active detail | `GET http://localhost:8080/api/v1/products/1` |
| Product admin | `POST` | `/api/v1/products` | Có | `admin` | `category_id,name,description,price,stock,image_url,is_active` | product mới | `POST http://localhost:8080/api/v1/products` |
| Product admin | `POST` | `/api/v1/products/upload-image` | Có | `admin` | `multipart/form-data`, field `file` | `{ url }` để gán lại vào `image_url` | `POST http://localhost:8080/api/v1/products/upload-image` |
| Product admin | `PUT` | `/api/v1/products/:id` | Có | `admin` | payload như create | product cập nhật | `PUT http://localhost:8080/api/v1/products/1` |
| Product admin | `DELETE` | `/api/v1/products/:id` | Có | `admin` | param `id` | message xóa mềm | `DELETE http://localhost:8080/api/v1/products/1` |
| Product admin | `GET` | `/api/v1/admin/products` | Có | `admin` | `page,limit,keyword,category_id,min_price,max_price,status` | danh sách + meta | `GET http://localhost:8080/api/v1/admin/products?status=all` |
| Product admin | `PUT` | `/api/v1/admin/products/:id/restore` | Có | `admin` | param `id` | product restored | `PUT http://localhost:8080/api/v1/admin/products/1/restore` |
| Order | `POST` | `/api/v1/orders` | Có | `user/admin` | `items[{product_id,quantity}]` | order mới | `POST http://localhost:8080/api/v1/orders` |
| Order | `GET` | `/api/v1/orders` | Có | `user/admin` | Không | admin: tất cả order, user: order của mình | `GET http://localhost:8080/api/v1/orders` |
| Order | `GET` | `/api/v1/orders/:id` | Có | `user/admin` | param `id` | order detail | `GET http://localhost:8080/api/v1/orders/1` |
| Order | `PUT` | `/api/v1/orders/:id/status` | Có | `admin` | `status` | order sau cập nhật status | `PUT http://localhost:8080/api/v1/orders/1/status` |
| Order | `GET` | `/api/v1/users/me/orders` | Có | `user/admin` | Không | danh sách order của current user | `GET http://localhost:8080/api/v1/users/me/orders` |
| User admin | `GET` | `/api/v1/users` | Có | `admin` | `page,limit,search` | danh sách user + meta | `GET http://localhost:8080/api/v1/users?page=1&limit=10` |
| User admin | `GET` | `/api/v1/users/:id` | Có | `admin` | param `id` | user detail | `GET http://localhost:8080/api/v1/users/1` |
| User admin | `PUT` | `/api/v1/users/:id` | Có | `admin` | `name,email,role` | user cập nhật | `PUT http://localhost:8080/api/v1/users/1` |
| User admin | `DELETE` | `/api/v1/users/:id` | Có | `admin` | param `id` | message vô hiệu hóa | `DELETE http://localhost:8080/api/v1/users/1` |

Không có endpoint thật cho:

- Payment
- Voucher
- Email
- Shipping address
- Staff/Manager riêng

Google OAuth đã có backend route và frontend callback thật, nhưng chỉ hoạt động end-to-end khi cấu hình Google OAuth credential hợp lệ.

### Đối chiếu endpoint trọng điểm với code thật

| Endpoint | Có route thật không | Handler | Service | Repository | Auth/role | Request thật | Response thật | Frontend có gọi không |
|---|---|---|---|---|---|---|---|---|
| `PUT /api/v1/users/me` | Có | `UserHandler.UpdateMe` | `UserService.UpdateProfile` | `UserRepository.UpdateProfileName`, `UserRepository.FindByID` | JWT, `user/admin` | body `{ "name": "..." }` | `response.OK` với `dto.UserResponse` | Có, `frontend/src/api/userApi.js -> updateMe`, dùng ở `ProfilePage` qua `AuthContext.updateProfile` |
| `GET /api/v1/admin/categories` | Có | `CategoryHandler.AdminList` | `CategoryService.AdminList` | `CategoryRepository.ListAdmin` | JWT + `admin` | query `status=all|active|inactive` | `response.OK` với mảng `dto.CategoryResponse` | Có, `frontend/src/api/categoryApi.js -> listAdmin`, dùng ở `AdminCategoriesPage` |
| `PUT /api/v1/admin/categories/:id/restore` | Có | `CategoryHandler.Restore` | `CategoryService.Restore` | `CategoryRepository.Restore`, `CategoryRepository.FindByID` | JWT + `admin` | param `id` | `response.OK` với `dto.CategoryResponse` | Có, `frontend/src/api/categoryApi.js -> restore`, dùng ở `AdminCategoriesPage` |
| `GET /api/v1/admin/products` | Có | `ProductHandler.AdminList` | `ProductService.AdminList` | `ProductRepository.List` | JWT + `admin` | query `page,limit,keyword,category_id,min_price,max_price,status` | `response.Paginated` với mảng `dto.ProductResponse` + `meta` | Có, `frontend/src/api/productApi.js -> listAdmin`, dùng ở `AdminProductsPage` |
| `PUT /api/v1/admin/products/:id/restore` | Có | `ProductHandler.Restore` | `ProductService.Restore` | `ProductRepository.FindByID`, `CategoryRepository.FindActiveByID`, `ProductRepository.Restore` | JWT + `admin` | param `id` | `response.OK` với `dto.ProductResponse` | Có, `frontend/src/api/productApi.js -> restore`, dùng ở `AdminProductsPage` |
| `GET /api/v1/orders` | Có | `OrderHandler.List` | `OrderService.List` | `OrderRepository.ListAll` hoặc `OrderRepository.ListByUserID`, sau đó `OrderRepository.FindItemsByOrderIDs` | JWT, `user/admin`; admin thấy tất cả, user thấy order của mình | query `page`, `limit`, `status` | `response.Paginated` với mảng `dto.OrderResponse` + `meta` | Có, `frontend/src/api/orderApi.js -> list`, dùng ở `AdminOrdersPage`, `AdminDashboardPage` |
| `GET /api/v1/users/me/orders` | Có | `OrderHandler.MyOrders` | `OrderService.List` với role ép là `user` | `OrderRepository.ListByUserID`, `OrderRepository.FindItemsByOrderIDs` | JWT, `user/admin`, nhưng luôn trả order của current user | query `page`, `limit`, `status` | `response.Paginated` với mảng `dto.OrderResponse` + `meta` | Có, `frontend/src/api/orderApi.js -> myOrders`, dùng ở `MyOrdersPage` |
| `PUT /api/v1/orders/:id/status` | Có | `OrderHandler.UpdateStatus` | `OrderService.UpdateStatus` | `OrderRepository.FindByID`, `OrderRepository.UpdateStatus`, `OrderRepository.FindItemsByOrderID` | JWT + `admin` | body `{ "status": "pending|confirmed|shipping|completed|cancelled" }`, nhưng service chỉ cho phép chuyển trạng thái hợp lệ | `response.OK` với `dto.OrderResponse` | Có, `frontend/src/api/orderApi.js -> updateStatus`, dùng ở `AdminOrdersPage`, `AdminOrderDetailPage` |

## 12. Bộ lệnh curl test chức năng

> Ghi chú quan trọng:
> - API trả JSON chuẩn dạng `success/message/data/errors/meta`.
> - Tài khoản admin seed có email `admin@example.com` và password seed chuẩn trong repo là `123456`.
> - Nếu volume DB đang chạy là dữ liệu cũ, cần reset volume để quay về đúng credential seed.
> - Credential admin đã xác minh trên volume local hiện tại là `vu@gmail.com / 123456`.

```bash
BASE=http://localhost:8080
API=http://localhost:8080/api/v1
ADMIN_EMAIL=vu@gmail.com
ADMIN_PASSWORD=123456
USER_EMAIL=user_demo_$(date +%s)@example.com
USER_PASSWORD=User@123
```

### Health check

```bash
curl "$BASE/health"
```

### Register user

```bash
curl -X POST "$API/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Demo User",
    "email": "'"$USER_EMAIL"'",
    "password": "'"$USER_PASSWORD"'"
  }'
```

### Login user

```bash
USER_LOGIN_JSON=$(curl -s -X POST "$API/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'"$USER_EMAIL"'",
    "password": "'"$USER_PASSWORD"'"
  }')

echo "$USER_LOGIN_JSON"
USER_ACCESS_TOKEN=$(echo "$USER_LOGIN_JSON" | jq -r '.data.access_token')
USER_REFRESH_TOKEN=$(echo "$USER_LOGIN_JSON" | jq -r '.data.refresh_token')
```

### Login admin nếu có seed admin

```bash
ADMIN_LOGIN_JSON=$(curl -s -X POST "$API/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'"$ADMIN_EMAIL"'",
    "password": "'"$ADMIN_PASSWORD"'"
  }')

echo "$ADMIN_LOGIN_JSON"
ADMIN_ACCESS_TOKEN=$(echo "$ADMIN_LOGIN_JSON" | jq -r '.data.access_token')
ADMIN_REFRESH_TOKEN=$(echo "$ADMIN_LOGIN_JSON" | jq -r '.data.refresh_token')
```

### Upload ảnh sản phẩm

```bash
curl -X POST "$API/products/upload-image" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -F "file=@C:/duong-dan/toi-anh-san-pham.png"
```

Lấy `url` trả về rồi gán lại vào `image_url` khi tạo hoặc cập nhật sản phẩm:

```bash
PRODUCT_IMAGE_URL="http://localhost:8080/uploads/products/images/example.png"

curl -X PUT "$API/products/$PRODUCT_ID" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "name": "Demo Product Updated",
    "description": "Updated by curl",
    "price": 120000,
    "stock": 8,
    "image_url": "'"$PRODUCT_IMAGE_URL"'",
    "is_active": true
  }'
```

### Upload avatar hồ sơ

```bash
curl -X POST "$API/users/me/avatar" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN" \
  -F "file=@C:/duong-dan/toi-avatar.png"
```

### Upload video hồ sơ ngắn

```bash
curl -X POST "$API/users/me/profile-video" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN" \
  -F "file=@C:/duong-dan/toi-video.mp4"
```

### Get me

```bash
curl "$API/auth/me" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN"
```

### Refresh token

```bash
curl -X POST "$API/auth/refresh-token" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "'"$USER_REFRESH_TOKEN"'"
  }'
```

### Logout

```bash
curl -X POST "$API/auth/logout" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "'"$USER_REFRESH_TOKEN"'"
  }'
```

### List categories

```bash
curl "$API/categories"
```

### Admin create category

```bash
curl -X POST "$API/categories" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Demo Category",
    "description": "Category for API test",
    "is_active": true
  }'
```

### Admin update category

```bash
CATEGORY_ID=<CATEGORY_ID>

curl -X PUT "$API/categories/$CATEGORY_ID" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Demo Category Updated",
    "description": "Updated from curl",
    "is_active": true
  }'
```

### Admin delete category

```bash
curl -X DELETE "$API/categories/$CATEGORY_ID" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN"
```

### List products

```bash
curl "$API/products"
```

### Search/filter/pagination products

```bash
curl "$API/products?page=1&limit=10&keyword=pin&category_id=3&min_price=0&max_price=1000"
```

### Product detail

```bash
PRODUCT_ID=<PRODUCT_ID>
curl "$API/products/$PRODUCT_ID"
```

### Admin create product

```bash
curl -X POST "$API/products" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "name": "Demo Product",
    "description": "Created by curl",
    "price": 100000,
    "stock": 5,
    "image_url": "https://example.com/demo-product.jpg",
    "is_active": true
  }'
```

### Admin update product

```bash
curl -X PUT "$API/products/$PRODUCT_ID" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "name": "Demo Product Updated",
    "description": "Updated by curl",
    "price": 120000,
    "stock": 8,
    "image_url": "https://example.com/demo-product-updated.jpg",
    "is_active": true
  }'
```

### Admin delete product

```bash
curl -X DELETE "$API/products/$PRODUCT_ID" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN"
```

### User create order

```bash
curl -X POST "$API/orders" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      { "product_id": '"$PRODUCT_ID"', "quantity": 1 }
    ]
  }'
```

### User my orders

```bash
curl "$API/users/me/orders" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN"
```

### Order detail

```bash
ORDER_ID=<ORDER_ID>

curl "$API/orders/$ORDER_ID" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN"
```

### Admin list orders

```bash
curl "$API/orders" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN"
```

### Admin update order status

```bash
curl -X PUT "$API/orders/$ORDER_ID/status" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "confirmed"
  }'
```

### Admin list users

```bash
curl "$API/users?page=1&limit=10&search=" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN"
```

### Admin update user

```bash
TARGET_USER_ID=<TARGET_USER_ID>

curl -X PUT "$API/users/$TARGET_USER_ID" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated User",
    "email": "updated.user@example.com",
    "role": "user"
  }'
```

### Admin delete user

```bash
curl -X DELETE "$API/users/$TARGET_USER_ID" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN"
```

### Endpoint chưa có

- Payment/voucher/email/shipping: chưa có endpoint

## 13. Link test frontend theo từng màn hình

| Nhóm | URL | Điều kiện đăng nhập | Chức năng cần kiểm tra | API backend được gọi |
|---|---|---|---|---|
| Guest | `http://localhost:5173/` | Không | Home/product list, filter, pagination, add to cart | `GET /products`, `GET /categories` |
| Guest | `http://localhost:5173/products` | Không | Giống home | `GET /products`, `GET /categories` |
| Guest | `http://localhost:5173/products/:id` | Không | Product detail, quantity validation, add cart | `GET /products/:id` |
| Guest | `http://localhost:5173/login` | Không | Login | `POST /auth/login`, `POST /auth/refresh-token`, `GET /auth/me` |
| Guest | `http://localhost:5173/register` | Không | Register | `POST /auth/register` |
| Guest/User | `http://localhost:5173/cart` | Không bắt buộc | Cart localStorage, create order nếu login | `POST /orders` khi submit |
| User | `http://localhost:5173/profile` | Có | Xem profile, cập nhật tên, tải avatar, tải video hồ sơ ngắn và lưu lại | `GET /auth/me`, `PUT /users/me`, `POST /users/me/avatar`, `POST /users/me/profile-video` |
| User | `http://localhost:5173/my-orders` | Có | Danh sách order của tôi | `GET /users/me/orders` |
| User | `http://localhost:5173/orders/:id` | Có | Chi tiết order của tôi | `GET /orders/:id` |
| Admin | `http://localhost:5173/admin` | Admin | Dashboard thống kê và recent orders | `GET /categories`, `GET /products`, `GET /orders`, `GET /users` |
| Admin | `http://localhost:5173/admin/categories` | Admin | CRUD category, restore inactive | `GET /admin/categories`, `POST /categories`, `PUT /categories/:id`, `DELETE /categories/:id`, `PUT /admin/categories/:id/restore` |
| Admin | `http://localhost:5173/admin/products` | Admin | CRUD product, filter, pagination, restore, tải ảnh sản phẩm rồi lưu vào `image_url` | `GET /admin/products`, `POST /products`, `POST /products/upload-image`, `PUT /products/:id`, `DELETE /products/:id`, `PUT /admin/products/:id/restore`, `GET /categories` |
| Admin | `http://localhost:5173/admin/orders` | Admin | List all orders, filter theo status ở UI, update status | `GET /orders`, `PUT /orders/:id/status` |
| Admin | `http://localhost:5173/admin/orders/:id` | Admin | Order detail, update status | `GET /orders/:id`, `PUT /orders/:id/status` |
| Admin | `http://localhost:5173/admin/users` | Admin | List/search/update/disable user | `GET /users`, `PUT /users/:id`, `DELETE /users/:id` |

## 14. Chức năng đã hoàn thành đúng

### Frontend/UI sau đợt tối ưu mới

- [x] Public/User/Admin layout vẫn hoạt động đúng route cũ
- [x] AuthContext/CartContext/ProtectedRoute/AdminRoute không bị phá API
- [x] Design token frontend đã được tách riêng tại `frontend/src/styles/tokens.css`
- [x] Header public đã rõ navigation hơn, giữ luồng guest/user/admin
- [x] Product card/list/detail đã dễ đọc hơn, CTA rõ hơn
- [x] Login/Register đã gọn hơn, vẫn giữ Google login button
- [x] Profile page đã rõ profile info và hỗ trợ upload avatar/video hồ sơ ngắn
- [x] Admin layout/sidebar/header đã chuyên nghiệp hơn
- [x] Admin dashboard đã đồng bộ visual enterprise hơn
- [x] Focus-visible và prefers-reduced-motion vẫn được giữ
- [x] Frontend vẫn build pass sau khi tối ưu UI

### Auth

- [x] Register
- [x] Login
- [x] Refresh token
- [x] Logout
- [x] Get me
- [x] JWT middleware
- [x] Role guard admin

### User

- [x] Admin list user
- [x] Admin search user
- [x] Admin update user
- [x] Admin delete/vô hiệu hóa user
- [x] User profile read
- [x] User cập nhật tên hiển thị tối thiểu
- [x] User upload avatar và lưu `avatar_url`
- [x] User upload video hồ sơ ngắn và lưu `profile_video_url`

### Category

- [x] Public list active category
- [x] Public category detail active-only
- [x] Admin create/update/delete
- [x] Soft delete
- [x] Admin list active/inactive/all
- [x] Admin restore inactive category

### Product

- [x] Public list/detail active-only
- [x] Search/filter/pagination
- [x] Admin create/update/delete
- [x] Admin list active/inactive/all
- [x] Admin restore inactive product
- [x] `image_url` preview ở frontend admin
- [x] Admin upload ảnh sản phẩm vào local storage và dùng URL trả về để lưu vào `image_url`

### Order

- [x] Create order
- [x] Backend tự tính `total_amount`
- [x] Backend lấy `price` từ DB
- [x] Backend trừ stock
- [x] Tạo order trong transaction
- [x] User chỉ xem order của mình
- [x] Admin xem toàn bộ order
- [x] Admin update status theo flow hợp lệ
- [x] `OrderResponse` hiện có `created_at`, `updated_at`
- [x] `OrderResponse` hiện có user summary (`id`, `name`, `email`)

### Frontend

- [x] Public layout
- [x] User layout
- [x] Admin layout
- [x] `AuthContext`
- [x] `CartContext`
- [x] `ProtectedRoute`
- [x] `AdminRoute`
- [x] Liquid glass UI
- [x] Loading state
- [x] Error state
- [x] Empty state
- [x] Admin dashboard

### Docker/Local

- [x] PostgreSQL chạy được bằng Docker Compose
- [x] Backend chạy được
- [x] Frontend build được
- [x] Docker compose hợp lệ
- [x] `/health` trả kết quả đúng

## 15. Chức năng chưa đúng, thiếu hoặc cần cải thiện

| Chức năng | Trạng thái | Nguyên nhân | Ảnh hưởng | Ưu tiên | Hướng xử lý |
|---|---|---|---|---|---|
| Admin seed password trong tài liệu không khớp runtime hiện tại | Đã rõ nguyên nhân | Tài liệu cũ từng ghi `Admin@123`, trong khi seed chuẩn cho `admin@example.com` là `123456`; volume DB cũ vẫn có thể chứa dữ liệu khác seed | Demo/curl admin có thể fail nếu dùng volume cũ mà không reset | Cao | Đồng bộ tài liệu về `123456` và ghi rõ cần reset volume nếu muốn quay về seed sạch |
| Cancel order không hoàn stock | Chưa đúng nghiệp vụ đầy đủ | `UpdateStatus` chỉ đổi trạng thái, không có reverse inventory | Số lượng tồn có thể bị giữ sai khi hủy đơn | Cao | Bổ sung nghiệp vụ restock có kiểm soát và transaction |
| Order list pagination/filter vừa được bổ sung | Đã cải thiện | `GET /orders` và `GET /users/me/orders` nay dùng `page`, `limit`, `status` và trả `meta` | Giảm tải dữ liệu và rõ ràng hơn cho admin/user | Thấp | Theo dõi thêm UX và backward compatibility frontend |
| Order items trước đây có N+1 query | Đã cải thiện | Service list nay batch `FindItemsByOrderIDs` thay vì lặp từng order | Giảm số query khi list nhiều order | Thấp | Theo dõi performance với dữ liệu lớn hơn |
| Admin order filter trước đây chỉ lọc ở frontend | Đã cải thiện | `AdminOrdersPage` nay gửi `status` xuống backend | Không còn cần tải toàn bộ order chỉ để lọc | Thấp | Có thể mở rộng thêm search theo user/order code sau |
| Google OAuth cần credential thật để chạy end-to-end | Đã implement theo kiểu env-gated | Backend đã có route/service/callback và frontend đã có nút + callback page, nhưng local hiện tại chưa gắn Google client thật | Nếu chưa cấu hình credential, người dùng sẽ được redirect về callback với lỗi cấu hình thay vì đăng nhập thành công | Trung bình | Khai báo Google OAuth Client trong Google Cloud Console và set đúng env local |
| Frontend chưa có automated test | Thiếu | `frontend/package.json` không có script test | Rủi ro regression UI/flow | Trung bình | Thêm Vitest/RTL hoặc E2E |
| Token lưu ở `localStorage` | Chỉ phù hợp demo | Frontend đọc/ghi token client-side | Kém an toàn hơn cookie httpOnly | Trung bình | Nếu production, chuyển chiến lược session/token |
| Payment/Voucher/Email chưa implement | Thiếu module | Không có bảng, API hay service | Không thể mô tả là e-commerce hoàn chỉnh | Cao | Triển khai Phase 2 nếu cần scope lớn hơn |
| Upload media hiện mới ở mức local filesystem, chưa thành media module hoàn chỉnh | Mới hoàn thiện core | Đã có upload API thật cho ảnh sản phẩm, avatar, video hồ sơ; nhưng chưa có media table, chưa có xoá media, chưa có cloud storage/CDN | Dùng tốt cho local demo nhưng chưa đủ cho production-scale media management | Trung bình | Nếu mở rộng Phase 2, thêm media metadata, delete flow, quota và object storage |
| Ảnh sản phẩm là luồng 2 bước: upload file rồi mới lưu `image_url` vào sản phẩm | Cần lưu ý khi demo | API upload chỉ trả `url`; record sản phẩm chỉ đổi sau `POST /products` hoặc `PUT /products/:id` | Nếu chỉ upload mà chưa bấm lưu, trang chủ/trang chi tiết vẫn chưa hiện ảnh | Trung bình | Giữ nguyên cho core hiện tại, hoặc sau này gộp thành single-flow upload + save |
| Checkout chưa phải e-commerce checkout đầy đủ | Giới hạn scope | Không có payment, shipping, address, tax, coupon | Demo được order core nhưng chưa phải checkout hoàn chỉnh | Cao | Thiết kế module checkout Phase 2 |
| Shipping/address chưa có | Thiếu | Không có bảng/API/UI thật | Không lưu địa chỉ giao hàng | Cao | Thêm bảng `shipping_addresses`, luồng chọn địa chỉ |
| Admin category/product không bị thiếu inactive view | Đã đúng | Có `/admin/categories` và `/admin/products` | Không phải lỗi hiện tại | Không áp dụng | Giữ nguyên |
| Category soft-delete làm product public bị ẩn | Đúng theo thiết kế hiện tại | Public product list yêu cầu `p.is_active = true` và `c.is_active = true` | Public không thấy product thuộc category đã ẩn; admin vẫn thấy và phải restore category trước | Trung bình | Giữ nguyên nếu đúng nghiệp vụ, hoặc thêm cảnh báo quản trị rõ hơn |
| Frontend chưa có automated UI test | Thiếu | `frontend/package.json` chưa có script test | Khó kiểm tra hồi quy UI tự động | Trung bình | Bổ sung Vitest/RTL cho route guard, context, helper UI |
| UI đã tối ưu nhưng chưa browser smoke test đầy đủ theo từng breakpoint trong đợt này | Cần kiểm tra thêm | Đợt này mới verify bằng build/lint | Còn rủi ro nhỏ về spacing/overflow ở từng màn hình | Trung bình | Kiểm tra nhanh trên 375/768/1024/1440 trước demo |

## 16. Chức năng ẩn, disabled, preview-only hoặc Phase 2

| Chức năng | Trạng thái UI | Backend có chưa | Có gọi API thật không | Ghi chú |
|---|---|---|---|---|
| Profile avatar upload | Có UI thật | Có | Có | `ProfilePage` gọi `POST /api/v1/users/me/avatar`, backend lưu file local và cập nhật `avatar_url` |
| Profile video upload | Có UI thật | Có | Có | `ProfilePage` gọi `POST /api/v1/users/me/profile-video`, backend lưu file local và cập nhật `profile_video_url` |
| Product image upload | Có upload form thật | Có | Có | `AdminProductsPage` gọi `POST /api/v1/products/upload-image`; cần bấm lưu sản phẩm để URL được persist vào `products.image_url` |
| Google login | Có nút ở `LoginPage` và có callback page | Có | Có | Cần `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`, `FRONTEND_AUTH_CALLBACK_URL`, `OAUTH_STATE_SECRET`; nếu chưa cấu hình thì callback trả lỗi cấu hình |
| Voucher | Không có màn hình thật | Chưa | Không | Không có bảng/API voucher |
| MoMo | Không có màn hình thật | Chưa | Không | Không có payment backend |
| Bank transfer | Không có màn hình thật | Chưa | Không | Không có payment backend |
| Apple Pay | Không có màn hình thật | Chưa | Không | Không có payment backend |
| Email | Không có màn hình thật | Chưa | Không | Không có `email_logs`, không có mail service |
| Upload/Media module hoàn chỉnh | Mới ở mức core local upload | Chưa hoàn chỉnh | Có một phần | Có upload API thật và thư mục `uploads/`, nhưng chưa có `uploads`/`media` table, chưa có media manager đầy đủ |
| Staff/Manager | Không có route/UI/backend | Chưa | Không | Chỉ có role `admin`, `user` |
| Shipping/address | Không có màn hình thật | Chưa | Không | Không có địa chỉ giao hàng |

## 17. Docker và môi trường local

- `Dockerfile` dùng để build backend Go API thành binary `/app/server`.
- `docker-compose.yml` có 2 service:
  - `postgres`
  - `api`
- Network:
  - Dùng default compose network.
- Volume:
  - `postgres_data`
  - `./uploads:/app/uploads`
- Port mapping:
  - `5432:5432`
  - `8080:8080`
- Env truyền vào service `api`:
  - `PORT`
  - `DATABASE_URL`
  - `JWT_ACCESS_SECRET`
  - `JWT_REFRESH_SECRET`
  - `FRONTEND_URL`
  - `BACKEND_PUBLIC_URL`
  - `UPLOAD_DIR`
  - `ACCESS_TOKEN_MINUTES`
  - `REFRESH_TOKEN_HOURS`
- Static media:
  - Backend serve file local qua `/uploads/*`
  - Local dev dùng `UPLOAD_DIR=uploads`
  - Docker dùng `UPLOAD_DIR=/app/uploads`
- Cách xem log:

```bash
docker compose logs -f
docker compose logs -f postgres
docker compose logs -f api
```

- Cách reset:

```bash
docker compose down -v
docker compose up --build
```

- Cách kiểm tra database healthy:

```bash
docker compose ps
docker compose exec postgres psql -U postgres -d enterprise_order_management -c "SELECT 1;"
```

### Lỗi thường gặp

- Port `5432` bị chiếm
- Port `8080` bị chiếm
- Database volume cũ làm migration không chạy lại
- Frontend sai `VITE_API_BASE_URL`
- CORS error do `FRONTEND_URL` không đúng
- JWT secret thay đổi làm token cũ invalid
- Tài liệu admin password cũ không khớp với volume DB đang có dữ liệu cũ

## 18. Test lỗi và bảo mật

| Test | Mục tiêu | Lệnh curl | Response kỳ vọng | Ý nghĩa bảo mật/nghiệp vụ |
|---|---|---|---|---|
| Protected API không token | Chặn truy cập trái phép | `curl "$API/auth/me"` | `401` + `missing authorization header` | Xác nhận auth bắt buộc |
| Token sai | Chặn token giả | `curl "$API/auth/me" -H "Authorization: Bearer invalid"` | `401` + `invalid access token` | Chống truy cập bằng token giả |
| User thường gọi API admin | Chặn vượt quyền | `curl "$API/users" -H "Authorization: Bearer $USER_ACCESS_TOKEN"` | `403` | Xác nhận role guard |
| Register email trùng | Chặn duplicate | gọi `POST /auth/register` 2 lần cùng email | `409` + `Email already exists` | Tránh tài khoản trùng |
| Login sai password | Chặn login sai | `POST /auth/login` password sai | `401` + `Invalid email or password` | Kiểm tra auth error |
| Product price âm | Validate giá | `POST /products` với `price:-1` | `400` | Tránh dữ liệu âm |
| Product stock âm | Validate stock | `POST /products` với `stock:-1` | `400` | Tránh tồn kho âm |
| Category name thiếu | Validate category | `POST /categories` không có `name` | `400` validation failed | Kiểm tra DTO validation |
| Order rỗng | Chặn order không item | `POST /orders` với `"items":[]` | `400` | Đúng nghiệp vụ |
| Quantity <= 0 | Chặn quantity sai | `POST /orders` với `quantity:0` | `400` | Tránh order lỗi |
| Order vượt stock | Chặn bán quá tồn | `POST /orders` quantity > stock | `400` | Đảm bảo tồn kho |
| User xem order người khác | Chặn truy cập chéo | `GET /orders/:id` của người khác | `403` | Bảo vệ ownership |
| Admin update status sai flow | Chặn transition sai | đổi `completed -> pending` | `400` + `Invalid order status transition` | Bảo vệ state machine |
| SQL injection cơ bản login/search | Query an toàn | gửi payload kiểu `' OR 1=1 --` | Không bypass auth/search | Repository dùng placeholder `$1`, `$2` |

### Ví dụ lệnh

```bash
curl "$API/auth/me"

curl "$API/auth/me" \
  -H "Authorization: Bearer invalid-token"

curl "$API/users" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN"

curl -X POST "$API/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"wrong-password"}'

curl -X POST "$API/orders" \
  -H "Authorization: Bearer $USER_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"items":[]}'
```

## 19. Cách kiểm thử bằng Postman hoặc trình duyệt

### Import vào Postman

- Có thể import trực tiếp từ các lệnh `curl` ở mục 12.

### Tạo environment

- `API_BASE_URL`
- `USER_ACCESS_TOKEN`
- `ADMIN_ACCESS_TOKEN`
- `PRODUCT_ID`
- `CATEGORY_ID`
- `ORDER_ID`

### Thứ tự test đề xuất

1. Health.
2. Login admin.
3. Create category.
4. Create product.
5. Register user.
6. Login user.
7. Add cart trên frontend.
8. Create order.
9. Check database.
10. Admin update order status.

## 20. Kết quả test/build hiện tại

### Đã chạy thực tế

```bash
go test ./cmd/... ./internal/...
go vet ./cmd/... ./internal/...
docker compose config --quiet

cd frontend
npm run build
npm run lint
```

### Kết quả

| Lệnh | Kết quả |
|---|---|
| `go test ./cmd/... ./internal/...` | Pass |
| `go vet ./cmd/... ./internal/...` | Pass, không có output lỗi |
| `docker compose config --quiet` | Pass |
| `npm run build` | Pass, Vite build thành công |
| `npm run lint` | Pass |

### Ghi chú verify UI/UX mới

- Sau đợt tối ưu UI, frontend vẫn build pass và lint pass.
- CSS bundle vẫn build thành công sau khi thêm `frontend/src/styles/tokens.css` và chuẩn hóa lại các file style chính.
- Các thay đổi UI trong đợt này không yêu cầu thay đổi migration, backend route hay business logic.

### Kiểm tra runtime thêm

| Kiểm tra | Kết quả |
|---|---|
| `docker compose up -d postgres` | Pass |
| `docker compose up -d api` | Pass |
| `GET /health` | Pass, trả `status=ok` |
| Query bảng PostgreSQL | Pass |
| `POST /api/v1/products/upload-image` | Pass, trả URL file dưới `/uploads/products/images/...` |
| `POST /api/v1/users/me/avatar` | Pass, cập nhật `users.avatar_url` và file tồn tại trong `uploads/profile/avatars/...` |
| Product public/detail render ảnh sau khi `image_url` đã được lưu | Pass, các trang dùng `resolveAssetUrl(...)` để hiển thị đúng ảnh local |

### Lỗi/phát hiện trong quá trình verify

| Vấn đề | File/liên quan | Hướng xử lý |
|---|---|---|
| `admin@example.com` có thể không login được trên volume cũ nếu dữ liệu đã lệch seed | DB volume `postgres_data`, README/docs | Password seed chuẩn cho `admin@example.com` là `123456`; nếu volume cũ không khớp seed thì reset DB volume |
| Đăng nhập liên tiếp cùng một tài khoản admin runtime hiện tại | API `/api/v1/auth/login`, `internal/pkg/token/jwt.go`, `refresh_tokens.token_hash` | Đã xác minh đăng nhập lặp lại trả `200`; token refresh hiện có `jti` ngẫu nhiên để tránh va chạm hash trên DB |
| Dữ liệu trong DB hiện tại không còn là dữ liệu seed tối thiểu | `postgres_data` | Nếu cần demo sạch, chạy `docker compose down -v` |

## 21. Checklist sẵn sàng demo

- [x] Database chạy
- [x] Backend chạy
- [x] Frontend build chạy được
- [x] Migration đã tạo bảng
- [x] Có admin/user test trong DB hiện tại
- [x] Product/category có dữ liệu trong DB hiện tại
- [x] User tạo order được theo code/runtime
- [x] Stock bị trừ sau khi order theo code hiện tại
- [x] Admin update order status được
- [x] User thường không vào admin được
- [x] API lỗi trả JSON thống nhất
- [x] Frontend không gọi API giả
- [x] Payment/Voucher/Email được ghi đúng là chưa implement
- [x] Upload media core local đã được ghi đúng là đã có một phần, chưa phải media module hoàn chỉnh
- [x] Docker compose chạy được
- [x] Build/test pass
- [x] Admin credential demo trên volume hiện tại đã xác minh được (`vu@gmail.com / 123456`)
- [x] Frontend UI đã được polish lại theo hướng enterprise demo-ready
- [x] Frontend build/lint vẫn pass sau khi đổi design tokens và CSS
- [x] Upload ảnh sản phẩm chạy được và URL lưu lại hiển thị được ở public pages sau khi lưu product
- [x] Upload avatar hồ sơ chạy được và lưu lại vào DB + filesystem
- [x] Upload video hồ sơ ngắn đã có API/backend/frontend thật cho local demo

## 22. Kịch bản demo đề xuất

1. Giới thiệu project là hệ thống quản lý sản phẩm, tồn kho cơ bản, user và order lifecycle; không gọi là e-commerce hoàn chỉnh.
2. Giới thiệu database PostgreSQL và Docker Compose root repo.
3. Chạy backend + frontend.
4. Guest xem danh sách sản phẩm và chi tiết sản phẩm.
5. User đăng ký/đăng nhập.
6. Admin tạo category/product, tải ảnh sản phẩm, lưu `image_url` và kiểm tra ảnh hiện ở trang public.
7. User cập nhật hồ sơ, tải avatar hoặc video hồ sơ ngắn.
8. User thêm vào cart và tạo order.
9. Kiểm tra database ở `orders`, `order_items`, `products.stock` và các URL media đã lưu.
10. Admin cập nhật trạng thái order theo flow.
11. Test lỗi khi user thường gọi API admin.
12. Kết luận các chức năng đã có và các module Phase 2.

## 23. Kết luận

- Project hiện tại đã đạt được:
  - Backend Go/Echo chạy ổn với PostgreSQL.
  - Có auth JWT + refresh token.
  - Có quản lý category, product, user và order lifecycle cơ bản.
  - Đã có local upload thật cho ảnh sản phẩm, avatar hồ sơ và video hồ sơ ngắn.
  - Frontend React/Vite đã nối API thật cho public/user/admin flow chính.
  - Frontend UI đã được chuẩn hóa lại theo hướng enterprise dashboard, không đổi API và không phá flow chính.
  - Docker Compose root repo chạy được cho `postgres` và `api`.
- Project chưa đạt:
  - Chưa có payment.
  - Chưa có voucher.
  - Chưa có email.
  - Chưa có shipping/address.
  - Chưa có role staff/manager.
  - Chưa có checkout hoàn chỉnh kiểu e-commerce.
  - Chưa có media module hoàn chỉnh với bảng metadata riêng, xóa media, cloud storage hay CDN.
- Scope hiện tại phù hợp để mô tả là:
  - Hệ thống quản lý sản phẩm, tồn kho cơ bản, người dùng và vòng đời đơn hàng.
- Không nên mô tả project là e-commerce hoàn chỉnh khi backend chưa có payment/shipping/voucher.
- Các việc nên làm tiếp theo:
  - Nếu muốn quay về seed sạch theo migration, reset DB volume và xác nhận lại credential `admin@example.com`.
  - Quyết định nghiệp vụ hoàn stock khi cancel theo rule idempotent rõ ràng trước khi mở rộng demo dữ liệu lớn.
  - Bổ sung automated test frontend.
  - Browser smoke test các màn hình chính sau đợt polish UI ở 375/768/1024/1440.
  - Nếu mở rộng scope, triển khai Phase 2: payment, voucher, media manager hoàn chỉnh, email, shipping/address, staff/manager.

## 24. Kết quả verify tonghop.md

| Nội dung kiểm tra | Trạng thái | Bằng chứng trong code | Cách đã sửa |
|---|---|---|---|
| Runtime backend chính là root `cmd/api` | Đúng | `cmd/api/main.go` khởi động app; `internal/http/server.go` wiring route | Giữ nguyên |
| API prefix là `/api/v1` | Đúng | `internal/http/server.go`: `api := e.Group("/api/v1")` | Giữ nguyên |
| Docker chính là root `Dockerfile` và root `docker-compose.yml` | Đúng | `Dockerfile` build `./cmd/api`; `docker-compose.yml` build `.` và mount `./migrations` | Giữ nguyên |
| `backend/` chỉ là skeleton/module cũ | Đúng | `README.md`, `docs/CURRENT_PROJECT_SCOPE_ANALYSIS.md`, `docs/ARCHITECTURE.md` đều ghi không dùng `backend/` làm runtime | Giữ nguyên |
| Migration thật ở `migrations/001_init.sql` | Đúng | `docker-compose.yml` mount `./migrations:/docker-entrypoint-initdb.d`; schema nằm trong `migrations/001_init.sql` | Giữ nguyên |
| Chỉ có 7 bảng nghiệp vụ chính | Đúng | `migrations/001_init.sql` chỉ tạo `roles`, `users`, `refresh_tokens`, `categories`, `products`, `orders`, `order_items`; query runtime PostgreSQL cũng ra 7 bảng | Giữ nguyên |
| `payments`, `vouchers`, `uploads`, `media`, `email_logs`, `shipping_addresses`, `inventory_logs` | Đúng là chưa implement | Không có trong `migrations/001_init.sql`; kiểm tra `to_regclass(...)` runtime trả `NULL` | Giữ nguyên |
| Seed admin có tồn tại | Đúng | `migrations/001_init.sql` insert `admin@example.com` | Giữ nguyên |
| Password admin demo có xác nhận được không | Đúng | `internal/pkg/password/password_test.go` xác nhận bcrypt seed trong `migrations/001_init.sql` khớp với `123456`; volume cũ có thể khác dữ liệu seed | Tài liệu hóa thống nhất password seed là `123456` và nhắc reset volume khi cần |
| Google OAuth backend route và frontend callback | Đúng | `internal/http/server.go` có `/api/v1/auth/google/login` và `/api/v1/auth/google/callback`; frontend có `LoginPage` nút Google và route `/auth/google/callback` | Bổ sung vào endpoint/docs và ghi rõ cơ chế env-gated |
| Schema OAuth mới | Đúng | `migrations/002_google_oauth.sql` thêm `users.avatar_url` và bảng `oauth_accounts`; DB runtime hiện tại đã apply và query thấy `oauth_accounts` tồn tại | Cập nhật phần database/env/docs theo schema mới |
| `PUT /api/v1/users/me` | Đúng | `internal/http/server.go`, `internal/handler/user_handler.go`, `internal/service/user_service.go`, `internal/repository/user_repository.go`; frontend gọi qua `frontend/src/api/userApi.js` | Bổ sung bảng đối chiếu endpoint trọng điểm |
| `GET /api/v1/admin/categories` | Đúng | `internal/http/server.go` + `CategoryHandler.AdminList` + `CategoryService.AdminList` + `CategoryRepository.ListAdmin`; frontend `categoryApi.listAdmin` | Bổ sung bảng đối chiếu endpoint trọng điểm |
| `PUT /api/v1/admin/categories/:id/restore` | Đúng | `internal/http/server.go` + `CategoryHandler.Restore` + `CategoryService.Restore` + `CategoryRepository.Restore` | Bổ sung bảng đối chiếu endpoint trọng điểm |
| `GET /api/v1/admin/products` | Đúng | `internal/http/server.go` + `ProductHandler.AdminList` + `ProductService.AdminList` + `ProductRepository.List`; frontend `productApi.listAdmin` | Bổ sung bảng đối chiếu endpoint trọng điểm |
| `PUT /api/v1/admin/products/:id/restore` | Đúng | `internal/http/server.go` + `ProductHandler.Restore` + `ProductService.Restore` + repository liên quan | Bổ sung bảng đối chiếu endpoint trọng điểm |
| `GET /api/v1/orders` | Đúng | `internal/http/server.go` + `OrderHandler.List` + `OrderService.List` + `OrderRepository.ListAll/ListByUserID`; frontend `orderApi.list` | Bổ sung bảng đối chiếu endpoint trọng điểm |
| `GET /api/v1/users/me/orders` | Đúng | `internal/http/server.go` + `OrderHandler.MyOrders` + `OrderService.List`; frontend `orderApi.myOrders` | Bổ sung bảng đối chiếu endpoint trọng điểm |
| `PUT /api/v1/orders/:id/status` | Đúng | `internal/http/server.go` gắn `RequireRoles(model.RoleAdmin)`; handler/service/repo đúng; frontend `orderApi.updateStatus` | Bổ sung bảng đối chiếu endpoint trọng điểm |
| User cập nhật tên hiển thị tối thiểu | Đúng | `UpdateProfileRequest` chỉ có `name`; `UserService.UpdateProfile`; `ProfilePage` gọi `updateProfile` -> `userApi.updateMe` | Giữ nguyên |
| Admin list/restore inactive category | Đúng | `CategoryService.AdminList`, `CategoryService.Restore`; frontend `AdminCategoriesPage` gọi admin endpoint thật | Giữ nguyên |
| Admin list/restore inactive product | Đúng | `ProductService.AdminList`, `ProductService.Restore`; frontend `AdminProductsPage` gọi admin endpoint thật | Giữ nguyên |
| `OrderResponse` có `created_at`, `updated_at`, user summary | Đúng | `internal/dto/order.go` và `service.ToOrderResponse` | Giữ nguyên |
| Curl health command cũ chưa tối ưu | Sai | `tonghop.md` cũ dùng `curl "$API/../health"` | Đã sửa thành `BASE=http://localhost:8080` và `curl "$BASE/health"` |
| Body field trong curl dùng `name`, `category_id`, `price`, `stock`, `refresh_token` | Đúng | `internal/dto/auth.go`, `category.go`, `product.go`, `user.go`, `order.go` | Giữ nguyên, chỉ chuẩn hóa biến `BASE/API` |
| Frontend routes mục 13 | Đúng | `frontend/src/routes/AppRoutes.jsx` | Giữ nguyên |
| Không mô tả sai các module Phase 2 như payment/voucher/upload/email/shipping/staff/manager | Đúng | Không có route/service/repository/migration tương ứng trong root runtime | Giữ nguyên |
| Test/build status hiện tại | Đúng | Đã chạy lại: `go test`, `go vet`, `docker compose config --quiet`, `npm run build`, `npm run lint` đều pass | Cập nhật lại bằng chứng verify mới nhất trong file |

## 25. Cập nhật mới nhất Auth SĐT và Cart Quote

### Đã implement

- `Auth register` hiện hỗ trợ:
  - email + password
  - phone + password
  - email + phone + password
- `Auth login` hiện hỗ trợ:
  - `identifier` là email
  - `identifier` là số điện thoại
  - backward-compatible với field `email` cũ nếu frontend cũ còn dùng
- Migration mới:
  - `migrations/004_auth_phone.sql`
  - thêm `users.phone`
  - thêm `users.phone_verified_at`
  - bỏ `NOT NULL` ở `users.email` để hỗ trợ account chỉ có SĐT
- Admin user management hiện đã hiển thị và cập nhật được:
  - `phone`
- Cart backend mới:
  - `POST /api/v1/cart/quote`
  - backend trả:
    - `items`
    - `subtotal`
    - `discount_amount`
    - `shipping_fee`
    - `final_amount`
    - `warnings`
- Frontend `CartPage` hiện:
  - vẫn giữ localStorage cart
  - tự gọi `/api/v1/cart/quote`
  - hiển thị `warnings` nếu sản phẩm inactive hoặc thiếu stock
  - ưu tiên tổng tiền backend quote thay vì chỉ tin estimate local

### Chưa implement trong phase này

- Voucher
- Shipping quote thật
- Checkout aggregate `/api/v1/checkout`
- Payment `COD`, `bank_qr`, `momo`, `zalopay`, `vnpay`, `onepay`
- Payment admin management

### Kết quả verify phase này

- `go test ./cmd/... ./internal/...`: pass
- `go vet ./cmd/... ./internal/...`: pass
- `docker compose config --quiet`: pass
- `cd frontend && npm run build`: pass
- `cd frontend && npm run lint`: pass
- Backend test mới đã bổ sung cho:
  - register bằng phone
  - register thiếu cả email và phone
  - login bằng phone
  - login backward-compatible qua `email`
- Frontend build vẫn pass sau khi đổi:
  - `RegisterPage`
  - `LoginPage`
  - `AdminUsersPage`
  - `CartPage`

### File chính đã thay đổi trong phase này

- Backend:
  - `migrations/004_auth_phone.sql`
  - `internal/dto/auth.go`
  - `internal/dto/user.go`
  - `internal/dto/cart.go`
  - `internal/model/models.go`
  - `internal/repository/user_repository.go`
  - `internal/service/auth_service.go`
  - `internal/service/user_service.go`
  - `internal/service/cart_service.go`
  - `internal/handler/cart_handler.go`
  - `internal/http/server.go`
  - `internal/service/auth_service_test.go`
  - `internal/service/mocks_test.go`
- Frontend:
  - `frontend/src/pages/auth/RegisterPage.jsx`
  - `frontend/src/pages/auth/LoginPage.jsx`
  - `frontend/src/pages/admin/AdminUsersPage.jsx`
  - `frontend/src/pages/user/CartPage.jsx`
  - `frontend/src/api/cartApi.js`

### Endpoint mới của phase này

- `POST /api/v1/cart/quote`

Ví dụ request:

```json
{
  "items": [
    { "product_id": 1, "quantity": 2 }
  ]
}
```

### Ghi chú tương thích ngược

- Frontend cũ dùng login bằng `email/password` vẫn chạy được.
- `POST /api/v1/orders` cũ vẫn giữ nguyên.
- Cart vẫn hoạt động với localStorage như trước, chỉ được tăng cường thêm quote từ backend.
