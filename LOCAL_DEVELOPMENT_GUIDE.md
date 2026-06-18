# Hướng dẫn chạy và bảo trì môi trường local

Tài liệu này là nguồn hướng dẫn vận hành local chính của project.

- Ngày kiểm tra gần nhất: `2026-06-18`
- Hệ điều hành đã kiểm tra: Windows + PowerShell
- Thư mục project: `C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api`

> Các lệnh trong tài liệu dùng **PowerShell**, không dùng Command Prompt (`cmd`).
> Trong PowerShell, dùng `curl.exe` để tránh alias `curl` của `Invoke-WebRequest`.

## 1. Kiến trúc chạy local

Phương án chuẩn và ít lỗi nhất:

| Thành phần | Cách chạy | Địa chỉ |
| --- | --- | --- |
| PostgreSQL | Docker Compose | `localhost:5432` |
| Backend API | Docker Compose | `http://localhost:8080` |
| Frontend | Vite trên máy host | `http://localhost:5173` |

API base URL của frontend:

```text
http://localhost:8080/api/v1
```

CORS backend cho phép origin:

```text
http://localhost:5173
```

### Lưu ý về hai module backend

Backend hoàn chỉnh dùng để demo nằm ở **thư mục gốc**:

```text
go.mod
cmd/api/main.go
internal/
Dockerfile
docker-compose.yml
```

Thư mục `backend/` là một Go module/scaffold riêng. Không chạy `backend/cmd/server`
cho luồng demo đầy đủ. Khi nâng dependency, vẫn nên kiểm tra module này để tránh repo
bị lệch phiên bản.

## 2. Phiên bản hiện tại

| Thành phần | Phiên bản được khóa/kiểm tra |
| --- | --- |
| Go trên máy | `1.25.5` |
| Go trong Docker builder | `1.25.5-alpine3.23` |
| Alpine runtime | `3.23` |
| PostgreSQL | `18.3-alpine3.23` |
| Echo | `4.15.4` |
| pgx | `5.10.0` |
| React | `19.2.7` |
| Vite | `8.0.16` |
| Node.js | `24.11.1` |
| npm | `11.6.2` |

Kiểm tra công cụ trên máy:

```powershell
go version
node --version
npm --version
docker --version
docker compose version
git --version
```

Docker Desktop phải được mở và Docker Engine phải chạy trước khi dùng Compose.

## 3. Cấu hình môi trường

### Backend chạy bằng Docker Compose

Các biến backend local đã được khai báo trong `docker-compose.yml`. Backend trong
container kết nối database bằng hostname service `postgres`:

```text
postgres://postgres:postgres@postgres:5432/enterprise_order_management?sslmode=disable
```

### Backend chạy trực tiếp bằng Go

Tạo `.env` ở thư mục gốc nếu chưa có:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
Copy-Item .env.example .env
```

Nội dung local:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/enterprise_order_management?sslmode=disable
JWT_ACCESS_SECRET=change-this-access-secret
JWT_REFRESH_SECRET=change-this-refresh-secret
FRONTEND_URL=http://localhost:5173
ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_HOURS=168
```

Quy tắc hostname database:

- Backend chạy trên host bằng `go run`: dùng `localhost`.
- Backend chạy trong Compose: dùng `postgres`.

### Frontend

Tạo `frontend/.env` nếu chưa có:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api\frontend
Copy-Item .env.example .env
```

Nội dung:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

Sau khi sửa biến Vite, phải dừng và chạy lại `npm run dev`.

## 4. Chạy project lần đầu

### Terminal 1: database và backend

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
docker compose up -d --build
docker compose ps
```

Trạng thái mong đợi:

- `enterprise-order-postgres`: `Up ... (healthy)`
- `enterprise-order-api`: `Up`

Xem log khi cần:

```powershell
docker compose logs --tail=100 postgres
docker compose logs --tail=100 api
```

Theo dõi log liên tục:

```powershell
docker compose logs -f api
```

Nhấn `Ctrl+C` chỉ dừng việc xem log, không dừng container.

Kiểm tra backend:

```powershell
curl.exe http://localhost:8080/health
```

Response thành công:

```json
{"success":true,"message":"Success","data":{"status":"ok"}}
```

### Terminal 2: frontend

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api\frontend
npm install
npm run dev
```

Mở trình duyệt tại:

```text
http://localhost:5173
```

Tài khoản admin được seed bởi migration:

```text
Email: admin@example.com
Password: 123456
```

## 5. Chạy project hằng ngày

### Terminal 1

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
docker compose up -d
docker compose ps
curl.exe http://localhost:8080/health
```

### Terminal 2

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api\frontend
npm run dev
```

Không cần chạy `docker compose down -v` hằng ngày. Lệnh đó xóa database.

## 6. Kiểm tra PostgreSQL và migration

Kiểm tra phiên bản:

```powershell
docker exec enterprise-order-postgres postgres --version
```

Kiểm tra danh sách bảng:

```powershell
docker exec enterprise-order-postgres psql `
  -U postgres `
  -d enterprise_order_management `
  -c "\dt"
```

Các bảng mong đợi:

```text
roles
users
refresh_tokens
categories
products
orders
order_items
```

Kiểm tra admin seed:

```powershell
docker exec enterprise-order-postgres psql `
  -U postgres `
  -d enterprise_order_management `
  -c "SELECT id, full_name, email, role_id, is_active FROM users;"
```

> Cột tên người dùng là `full_name`, không phải `name`.

Kiểm tra role và category seed:

```powershell
docker exec enterprise-order-postgres psql `
  -U postgres `
  -d enterprise_order_management `
  -c "SELECT * FROM roles; SELECT * FROM categories;"
```

Mở phiên `psql` tương tác:

```powershell
docker exec -it enterprise-order-postgres psql `
  -U postgres `
  -d enterprise_order_management
```

Thoát bằng:

```sql
\q
```

### Cơ chế migration

Project không dùng `golang-migrate`. File `migrations/001_init.sql` được official
PostgreSQL image chạy tự động khi volume database được tạo lần đầu và còn rỗng.

Migration sẽ không tự chạy lại trên volume đã có dữ liệu.

## 7. Reset database local

> Cảnh báo: lệnh này xóa toàn bộ dữ liệu database local.

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
docker compose down -v
docker compose up -d --build
docker compose ps
```

Sau đó kiểm tra migration:

```powershell
docker exec enterprise-order-postgres psql `
  -U postgres `
  -d enterprise_order_management `
  -c "\dt"
```

### Quy tắc volume PostgreSQL 18+

Trong `docker-compose.yml`, volume phải mount tại:

```yaml
- postgres_data:/var/lib/postgresql
```

Không đổi lại thành `/var/lib/postgresql/data`; đó là layout cũ và làm official
PostgreSQL 18 image thoát với exit code `1`.

## 8. Chạy backend trực tiếp bằng Go

Phương án này dùng khi debug backend bằng IDE. Chỉ chạy database bằng Docker:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
docker compose up -d postgres
go mod download
go run ./cmd/api
```

Không chạy đồng thời API container và `go run`, vì cả hai cùng dùng port `8080`.
Nếu API container đang chạy:

```powershell
docker compose stop api
go run ./cmd/api
```

### Xung đột PostgreSQL Windows ở port 5432

Máy này từng có service `postgresql-x64-18` chạy trên Windows. Nếu `go run` báo
`password authentication failed` dù container healthy, kiểm tra:

```powershell
Get-Service postgresql-x64-18
Get-NetTCPConnection -LocalPort 5432 -State Listen
```

Mở PowerShell bằng quyền Administrator và dừng service Windows nếu không dùng:

```powershell
Stop-Service postgresql-x64-18
Set-Service postgresql-x64-18 -StartupType Manual
```

Sau đó:

```powershell
docker compose restart postgres
go run ./cmd/api
```

Nếu không muốn dừng PostgreSQL Windows, hãy chạy backend bằng Docker Compose để
backend kết nối qua network nội bộ tới hostname `postgres`.

## 9. Dừng môi trường

Dừng nhưng giữ dữ liệu:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
docker compose down
```

Dừng và xóa database:

```powershell
docker compose down -v
```

Dừng frontend bằng `Ctrl+C` tại terminal chạy Vite.

## 10. Kiểm tra chất lượng trước khi commit

Backend chính:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
go test ./cmd/... ./internal/...
```

Module `backend/`:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api\backend
go test ./...
```

Frontend:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api\frontend
npm run lint
npm run build
npm audit
```

Docker:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
docker compose config --quiet
docker compose build api
```

Hiện frontend có một số ESLint warning về dependency array của React Hook nhưng
không có lint error và production build vẫn thành công.

## 11. Cập nhật phiên bản an toàn

Không nâng tất cả và chạy production ngay. Thực hiện theo từng lớp, test sau mỗi lớp.

### Go và dependency backend

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api
go list -m -u all
go get -u all
go mod tidy
go test ./cmd/... ./internal/...
```

Lặp lại cho module scaffold:

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api\backend
go get -u all
go mod tidy
go test ./...
```

Khi đổi Go, đồng bộ ba nơi:

1. `go` và `toolchain` trong `go.mod` gốc.
2. `go` và `toolchain` trong `backend/go.mod`.
3. Tag `golang:` trong cả hai Dockerfile.

### Frontend

```powershell
cd C:\Users\thaiv\OneDrive\Documents\enterprise-order-management-api\frontend
npm outdated
npm update
npm audit
npm run lint
npm run build
```

Khi đổi Node/npm, đồng bộ:

- `engines` và `packageManager` trong `frontend/package.json`.
- `frontend/.nvmrc`.
- `frontend/package-lock.json` bằng `npm install`.

### PostgreSQL

Khi nâng **minor** trong cùng major, build và test lại bình thường.

Khi nâng **major**, ví dụ PostgreSQL 18 lên 19:

- Đọc breaking changes của official image.
- Sao lưu dữ liệu cần giữ.
- Với dữ liệu demo có thể dùng `docker compose down -v` rồi khởi tạo lại.
- Với dữ liệu cần giữ phải dùng `pg_dump`/`pg_restore` hoặc `pg_upgrade`.
- Kiểm tra lại mount path và healthcheck của image mới.

## 12. Troubleshooting nhanh

### PostgreSQL container exited (1)

```powershell
docker compose logs --tail=200 postgres
```

Với PostgreSQL 18, xác nhận mount là `/var/lib/postgresql`, sau đó reset volume lỗi:

```powershell
docker compose down -v
docker compose up -d --build
```

### Backend không khởi động

```powershell
docker compose ps
docker compose logs --tail=200 api
docker compose logs --tail=200 postgres
```

### Port 8080 bị chiếm

```powershell
Get-NetTCPConnection -LocalPort 8080 -State Listen
docker compose stop api
```

### Port 5432 bị chiếm

```powershell
Get-NetTCPConnection -LocalPort 5432 -State Listen
Get-Service postgresql-x64-18 -ErrorAction SilentlyContinue
```

### Frontend gọi sai API hoặc lỗi CORS

Kiểm tra:

```text
frontend/.env: VITE_API_BASE_URL=http://localhost:8080/api/v1
backend CORS:   FRONTEND_URL=http://localhost:5173
```

Sau khi sửa, restart Vite.

### Vite tự chuyển sang port 5174

Backend chỉ cho phép origin `http://localhost:5173`. Dừng tiến trình đang giữ port
5173 hoặc cập nhật `FRONTEND_URL` cho khớp rồi restart backend.

## 13. Ghi chú dành cho agent và maintainer

Khi agent mới tiếp nhận repo, đọc theo thứ tự:

1. `LOCAL_DEVELOPMENT_GUIDE.md`
2. `docker-compose.yml`
3. `Dockerfile`
4. `go.mod`
5. `cmd/api/main.go`
6. `internal/http/server.go`
7. `migrations/001_init.sql`
8. `frontend/package.json`
9. `frontend/src/api/apiClient.js`

Các invariant không được phá vỡ:

- Runtime backend chính nằm ở root, entrypoint `./cmd/api`.
- PostgreSQL 18 mount volume tại `/var/lib/postgresql`.
- Backend trong Compose dùng DB hostname `postgres`.
- Backend chạy trên host dùng DB hostname `localhost`.
- Frontend gọi API qua `VITE_API_BASE_URL` và path `/api/v1`.
- CORS `FRONTEND_URL` phải khớp chính xác origin Vite.
- User table dùng cột `full_name`.
- Migration init chỉ tự chạy trên volume rỗng.
- Refresh token phải được hash trước khi lưu database.
- Không commit `.env`, token, mật khẩu thật hoặc dữ liệu production.

Sau mọi thay đổi dependency hoặc Docker image, tối thiểu phải chạy:

```powershell
go test ./cmd/... ./internal/...
cd frontend
npm run lint
npm run build
cd ..
docker compose config --quiet
docker compose build api
```
