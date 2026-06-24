# Local Development Guide

Tài liệu này chỉ hướng dẫn chạy project ở local. Trạng thái chức năng xem tại [CURRENT_PROJECT_SCOPE_ANALYSIS.md](docs/CURRENT_PROJECT_SCOPE_ANALYSIS.md).

## Runtime

| Thành phần | Địa chỉ |
|---|---|
| Frontend | `http://localhost:5173` |
| Backend | `http://localhost:8080` |
| Health check | `http://localhost:8080/health` |
| PostgreSQL | `localhost:5432` |

Backend source of truth nằm tại `cmd/api`, `internal/*`, `migrations/001_init.sql` và các file Docker ở root. Không chạy skeleton `backend/`.

## Biến môi trường chính

Khi backend chạy trong Docker Compose:

```env
DATABASE_URL=postgres://postgres:postgres@postgres:5432/enterprise_order_management?sslmode=disable
PORT=8080
JWT_ACCESS_SECRET=change-this-access-secret
JWT_REFRESH_SECRET=change-this-refresh-secret
FRONTEND_URL=http://localhost:5173
```

Khi chạy Go trực tiếp trên máy:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/enterprise_order_management?sslmode=disable
PORT=8080
JWT_ACCESS_SECRET=change-this-access-secret
JWT_REFRESH_SECRET=change-this-refresh-secret
FRONTEND_URL=http://localhost:5173
```

Frontend dùng:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

Không dùng credential hoặc JWT secret local cho môi trường public/production.

## Khởi động lần đầu

```bash
docker compose up -d --build
docker compose ps
curl http://localhost:8080/health
```

Migration `migrations/001_init.sql` được áp dụng khi database được tạo lần đầu. PostgreSQL dùng volume tại `/var/lib/postgresql` theo cấu hình Compose hiện tại.

Chạy frontend ở terminal khác:

```bash
cd frontend
npm install
npm run dev
```

Tài khoản admin local:

- Email: `admin@example.com`
- Password: `Admin@123`

Credential đã xác minh trên volume local hiện tại:

- Email: `vu@gmail.com`
- Password: `123456`
- Role: `admin`

Nếu cần quay về dữ liệu seed sạch của migration, hãy reset volume:

```bash
docker compose down -v
docker compose up -d --build
```

## Lệnh thường dùng

```bash
docker compose up -d
docker compose logs -f api
docker compose logs -f postgres
docker compose stop
docker compose down
```

`docker compose down -v` sẽ xóa toàn bộ dữ liệu local; chỉ dùng khi chủ động reset database.

## Chạy backend ngoài Docker

Giữ PostgreSQL đang chạy, đặt biến môi trường rồi chạy:

```powershell
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/enterprise_order_management?sslmode=disable"
$env:PORT="8080"
$env:JWT_ACCESS_SECRET="change-this-access-secret"
$env:JWT_REFRESH_SECRET="change-this-refresh-secret"
$env:FRONTEND_URL="http://localhost:5173"
go run ./cmd/api
```

## Kiểm tra tối thiểu

```bash
go test ./cmd/... ./internal/...
go vet ./cmd/... ./internal/...

cd frontend
npm run lint
npm run build
```

Không dùng `go test ./...` nếu môi trường làm việc khiến lệnh quét nhầm `frontend/node_modules`.

## Kiểm tra database

```bash
docker compose exec postgres psql -U postgres -d enterprise_order_management
```

Trong `psql`:

```sql
\dt
SELECT u.id, u.email, r.name AS role, u.is_active
FROM users u
JOIN roles r ON r.id = u.role_id
ORDER BY u.id;
SELECT id, name, stock, is_active FROM products ORDER BY id;
SELECT id, user_id, status, total_amount, created_at FROM orders ORDER BY id DESC;
```

## Lỗi thường gặp

- Port `8080`, `5173` hoặc `5432` bận: dừng process/container đang giữ port hoặc đổi cấu hình local.
- Backend không kết nối DB: kiểm tra host `postgres` khi chạy trong Compose và `localhost` khi chạy trực tiếp.
- Frontend báo network/CORS: kiểm tra backend health và `VITE_API_BASE_URL`.
- Schema không cập nhật sau khi sửa migration ban đầu: migration chỉ chạy trên volume mới; sao lưu dữ liệu trước khi reset volume.
- Đăng nhập admin thất bại sau khi đổi dữ liệu: kiểm tra user seed trong database thay vì giả định credential mặc định còn nguyên.
