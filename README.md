# Enterprise Order Management API

Project thực tập gồm backend API Golang/Echo để quản lý sản phẩm, tồn kho cơ bản, người dùng và đơn hàng; frontend React chỉ phục vụ demo các luồng thật với backend và PostgreSQL.

## Công nghệ

- Backend: Go, Echo, pgx, JWT, bcrypt.
- Database: PostgreSQL, migration SQL.
- Frontend: React, JavaScript (JSX), Vite.
- Runtime local: Docker Compose.

## Source of truth

Backend nghiệp vụ chạy từ root repository:

- `cmd/api`
- `internal/*`
- `migrations/001_init.sql`
- `Dockerfile`
- `docker-compose.yml`

Thư mục `backend/` chỉ là skeleton cũ, không dùng để mô tả hoặc chạy runtime nghiệp vụ.

## Chức năng chính

- Auth: đăng ký, đăng nhập, đăng nhập bằng Google OAuth, refresh token, đăng xuất và lấy thông tin hiện tại.
- Profile: người dùng cập nhật tên hiển thị.
- Danh mục/sản phẩm: public chỉ thấy dữ liệu active; admin quản lý, xem inactive và khôi phục.
- Giỏ hàng: lưu trên `localStorage`, gửi đúng payload khi tạo đơn.
- Đơn hàng: tạo đơn theo transaction, giảm tồn kho, xem đơn cá nhân, admin cập nhật trạng thái.
- Quản trị: dashboard và các màn danh mục, sản phẩm, đơn hàng, người dùng.
- Giao diện: responsive, liquid glass, loading/error/empty states, toast và confirm dialog.

## Chạy nhanh

Yêu cầu: Docker Desktop và Node.js/npm.

```bash
docker compose up -d --build
curl http://localhost:8080/health

cd frontend
npm install
npm run dev
```

- Frontend: `http://localhost:5173`
- Backend: `http://localhost:8080`
- API prefix: `http://localhost:8080/api/v1`
- PostgreSQL: `localhost:5432`

Tài khoản admin local được seed bởi migration:

- Email: `admin@example.com`
- Password: `Admin@123`

Credential đã xác minh trên volume local hiện tại:

- Email: `vu@gmail.com`
- Password: `123456`
- Role: `admin`

Lưu ý:

- Migration chỉ seed `admin@example.com` với bcrypt hash, không seed plaintext password.
- Nếu muốn quay về dữ liệu seed sạch theo migration, cần reset volume:

```bash
docker compose down -v
docker compose up -d --build
```

- Frontend không chạy trong `docker-compose.yml` hiện tại; chỉ có `postgres` và `api`.

## Google OAuth local

Google login đã có route backend thật:

- `GET /api/v1/auth/google/login`
- `GET /api/v1/auth/google/callback`

Biến môi trường cần cấu hình khi muốn bật Google OAuth:

```env
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
FRONTEND_AUTH_CALLBACK_URL=http://localhost:5173/auth/google/callback
OAUTH_STATE_SECRET=change-this-google-oauth-state-secret
```

Redirect URI local cần khai báo trong Google Cloud Console:

- Backend callback: `http://localhost:8080/api/v1/auth/google/callback`
- Frontend callback: `http://localhost:5173/auth/google/callback`

Nếu chưa cấu hình credential sandbox/local, nút Google login sẽ không hoàn tất đăng nhập thật; backend sẽ redirect về callback frontend với lỗi cấu hình thay vì giả lập provider.

Chỉ dùng credential và JWT secret mặc định trong môi trường local; phải thay trước khi triển khai công khai.

## Kiểm tra chất lượng

Chỉ chạy test backend tại các package root, tránh quét `frontend/node_modules`:

```bash
go test ./cmd/... ./internal/...
go vet ./cmd/... ./internal/...

cd frontend
npm run lint
npm run build
```

Report gần nhất xác nhận các lệnh trên, Docker runtime và browser smoke test đã pass. Trước buổi demo vẫn nên chạy checklist trong runbook với database thật.

## Tài liệu

- [Hướng dẫn local](LOCAL_DEVELOPMENT_GUIDE.md)
- [API](docs/api.md)
- [Kiến trúc](docs/ARCHITECTURE.md)
- [Thiết kế database](docs/DATABASE_DESIGN.md)
- [Phạm vi hiện tại](docs/CURRENT_PROJECT_SCOPE_ANALYSIS.md)
- [Demo runbook](docs/DEMO_RUNBOOK.md)
- [Gap còn lại](docs/frontend-backend-gap.md)

## Giới hạn bản demo

Project không được mô tả như một hệ thống e-commerce hoàn chỉnh. Warning vận hành được gom tại [Demo runbook](docs/DEMO_RUNBOOK.md); backlog và danh sách Phase 2 nằm tại [Phạm vi hiện tại](docs/CURRENT_PROJECT_SCOPE_ANALYSIS.md).
