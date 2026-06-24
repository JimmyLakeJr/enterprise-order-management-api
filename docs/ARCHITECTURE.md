# Architecture

## Tổng quan

```text
React/Vite frontend
        |
        | HTTP JSON + JWT
        v
Echo router -> middleware -> handler -> service -> repository -> PostgreSQL
```

Backend là modular monolith nhỏ, phù hợp scope thực tập. Không tách microservice hoặc thêm module lớn trong sprint demo.

## Source of truth

```text
cmd/api/                 entrypoint
internal/config/         cấu hình môi trường
internal/http/           router/server
internal/middleware/     JWT và phân quyền
internal/handler/        HTTP binding/response
internal/service/        nghiệp vụ
internal/repository/     truy cập PostgreSQL
internal/dto/            request/response contract
internal/model/          model dùng chung
internal/pkg/            error, response, validator
migrations/001_init.sql  schema ban đầu
```

Root `Dockerfile` và `docker-compose.yml` định nghĩa runtime local. `backend/` là skeleton cũ, không phải source of truth.

## Trách nhiệm layer

| Layer | Trách nhiệm |
|---|---|
| Router/middleware | Route, JWT, role guard, CORS |
| Handler | Parse/validate request, ánh xạ response/error |
| Service | Business rule, authorization theo resource, transaction orchestration |
| Repository | SQL và mapping dữ liệu |
| Database | Constraint, quan hệ và tính nhất quán dữ liệu |

Không đặt business rule quan trọng chỉ ở frontend.

## Luồng quan trọng

### Auth

Login kiểm tra password hash, phát access/refresh token; refresh token được lưu để có thể thu hồi. Middleware đưa user ID/role vào request context.

### Tạo order

Service xác thực item, khóa/đọc dữ liệu cần thiết, kiểm tra stock, tính total, tạo order/items và giảm stock trong transaction. Frontend không quyết định giá.

### Active/inactive

Public repository chỉ trả category/product active. Admin dùng route riêng để xem `all|active|inactive` và restore; product chỉ restore khi category active.

## Bảo mật và vận hành

- Password dùng bcrypt; protected route dùng JWT và role middleware.
- User không tự sửa role/is_active; order có ownership check.
- Secret, DB credential và admin seed local phải thay trước public production.
- API hiện chưa có rate limiting, observability hoặc security hardening cấp production.

## Giới hạn kiến trúc

- Order list chưa pagination; item loading còn có thể batch.
- Cancel chưa có stock compensation.
- Media chỉ là URL, không có object storage/upload.
- Phase 2 không thuộc runtime hiện tại; xem [CURRENT_PROJECT_SCOPE_ANALYSIS.md](CURRENT_PROJECT_SCOPE_ANALYSIS.md).
