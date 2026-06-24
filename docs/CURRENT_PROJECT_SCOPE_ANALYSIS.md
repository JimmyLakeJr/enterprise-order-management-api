# Current Project Scope

## Mô tả

Project là backend API Golang/Echo quản lý sản phẩm, tồn kho cơ bản, người dùng và đơn hàng. Frontend React phục vụ demo các luồng thật; project không phải hệ thống e-commerce hoàn chỉnh.

Backend source of truth: `cmd/api`, `internal/*`, `migrations/001_init.sql`, root `Dockerfile` và `docker-compose.yml`. Thư mục `backend/` không phải runtime nghiệp vụ.

## Trạng thái module

| Module | Trạng thái hiện tại |
|---|---|
| Auth | Register, login, refresh, logout và `/auth/me` đã hoạt động |
| Profile | User cập nhật tên; response nhất quán với `/auth/me` |
| Category | Public active list; admin CRUD, inactive filter và restore |
| Product | Public active list/detail; admin CRUD, filter, inactive list và restore |
| Cart | Lưu `localStorage`, tạo order bằng `product_id` và `quantity` |
| Order | Transaction tạo đơn và giảm stock; user xem đơn; admin xem và đổi trạng thái |
| User admin | Danh sách và quản lý trạng thái/role trong scope hiện có |
| Frontend | Public, user, admin core; responsive liquid-glass UI |

## API/module hiện có

- `/api/v1/auth/*`: auth và session.
- `/api/v1/users/me`: đọc/cập nhật tên profile.
- `/api/v1/categories`, `/api/v1/products`: public chỉ trả dữ liệu active.
- `/api/v1/orders`: tạo, liệt kê có pagination/filter và xem đơn theo ownership.
- `/api/v1/categories`, `/products`, `/orders`, `/users`: CRUD/quản trị cốt lõi theo role admin.
- `/api/v1/admin/categories`, `/admin/products`: danh sách active/inactive và restore.
- Admin category/product hỗ trợ `status=all|active|inactive` và restore.
- Order admin có `created_at`, `updated_at`, user summary, pagination và status filter.

Chi tiết contract xem tại [api.md](api.md).

## Database

Schema cơ bản gồm:

- `roles`
- `users`
- `refresh_tokens`
- `categories`
- `products`
- `orders`
- `order_items`

Quan hệ và constraint xem tại [DATABASE_DESIGN.md](DATABASE_DESIGN.md) và [ERD.md](ERD.md).

## Đã xác nhận gần nhất

- Backend test/vet hẹp đã pass.
- Frontend lint/build đã pass.
- Docker runtime và browser smoke test public/user/admin đã pass.
- UI responsive, loading/error/empty state, toast và dialog đã có trong source.

Các kết quả này là snapshot tại thời điểm report; trước demo cần chạy lại checklist ngắn trong [DEMO_RUNBOOK.md](DEMO_RUNBOOK.md).

## Backlog thật

### Demo verification

- Smoke test lại public/user/admin với backend và DB cuối cùng.
- Kiểm tra console, network và responsive mobile/tablet.

### Frontend cleanup

- Dọn component/API alias/CSS không còn dùng.
- Kiểm tra toàn bộ UI dùng custom dialog thay native confirm.
- Chuẩn hóa nốt copy Việt/Anh và metadata `index.html`.
- Bổ sung automated test/E2E khi scope cho phép.

### Backend tối ưu nhỏ

- Chốt nghiệp vụ cancel/restock, idempotency và concurrency trước khi hoàn kho.

### Phase 2

Payment, Voucher, Upload/Media, Email, Shipping/address, Staff/Manager và avatar/profile video upload.

## Giới hạn demo

- Cancel chưa hoàn stock.
- Product chỉ lưu `image_url`.
- Profile chỉ cập nhật tên, không upload media.
- Credential admin seed sạch chỉ chắc chắn sau khi reset volume DB; volume hiện tại đã xác minh `vu@gmail.com / 123456` có role `admin`.
- Không dùng secret hoặc credential mặc định cho public production.
