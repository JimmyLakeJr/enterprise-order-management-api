# Phân tích scope hiện tại của Enterprise Order Management

> Ngày kiểm tra: 2026-06-18  
> Nguồn xác minh: code runtime ở root repo, route, DTO, service/repository, migration, Docker Compose, frontend React và kiểm tra database local đang chạy.  
> Nguyên tắc: nội dung dưới đây mô tả chức năng đang có thật, không suy diễn từ tên project hoặc tài liệu định hướng.

## Lưu ý quan trọng về cấu trúc repo

Repo đang có **hai cây backend**:

- Backend chính, đang được root `Dockerfile` và root `docker-compose.yml` sử dụng: `cmd/api` + `internal/*` + `migrations/001_init.sql`.
- Thư mục `backend/` là một Go module riêng, hiện chỉ đăng ký endpoint `GET /health`. Compose và Dockerfile nằm trong thư mục này là cấu hình khác, không phải cấu hình được root Compose sử dụng.

Vì vậy, toàn bộ API nghiệp vụ trong tài liệu này lấy backend root làm source of truth. Khi chạy project theo README bằng `docker compose up --build` hoặc `go run ./cmd/api`, backend root là runtime thực tế.

Stack được xác nhận từ repo:

- Backend: Go 1.25, Echo v4.
- Database: PostgreSQL; SQL thuần qua `pgx/v5` và `pgxpool`, không có ORM.
- Auth: JWT access token + refresh token; refresh token được hash SHA-256 trước khi lưu database.
- Frontend: React 19 + Vite 8 + React Router + Axios.
- Local cart và token demo được lưu trong `localStorage`.

# 1. Project này làm về cái gì?

Đây là hệ thống demo quản lý danh mục, sản phẩm, tồn kho cơ bản, tài khoản và đơn hàng trong doanh nghiệp. Project có luồng giống một cửa hàng trực tuyến nhỏ để minh họa nghiệp vụ, nhưng **chưa phải một hệ thống e-commerce hoàn chỉnh**.

Phạm vi hiện có:

- Khách xem danh mục, tìm/lọc sản phẩm và xem chi tiết sản phẩm.
- Người dùng đăng ký, đăng nhập, lưu giỏ hàng trên frontend, tạo đơn và xem đơn của chính mình.
- Admin quản lý danh mục, sản phẩm, user và đơn hàng; cập nhật trạng thái đơn theo luồng cho phép.
- Backend tự lấy giá hiện hành, kiểm tra tồn kho, tạo order/order item và trừ kho trong một transaction.

Những thành phần thường có ở e-commerce hoàn chỉnh nhưng repo chưa có gồm địa chỉ giao hàng, phí vận chuyển, phương thức vận chuyển, thanh toán, voucher, upload/media, email, hoàn tiền và tích hợp cổng thanh toán.

Vai trò của backend:

- Cung cấp REST API tại `/api/v1`.
- Xử lý xác thực, phân quyền `user`/`admin`, validation và response envelope.
- Thực thi business rule và transaction tạo đơn.
- Đọc/ghi PostgreSQL bằng SQL thuần.

Vai trò của frontend:

- Là React SPA để thao tác thật với backend.
- Cung cấp giao diện public, user và admin.
- Quản lý cart ở localStorage; gửi lên backend chỉ `product_id` và `quantity` khi tạo order.
- Không thay backend quyết định giá, tồn kho, quyền truy cập hoặc trạng thái đơn.

# 2. Project này làm cho ai?

Repo chỉ định nghĩa hai role trong database/code là `user` và `admin`. `Guest` là trạng thái chưa đăng nhập, không phải một role trong bảng `roles`. Không có `staff` hoặc `manager`.

## Guest

- Xem danh sách category đang active.
- Xem danh sách và chi tiết product đang active, thuộc category active.
- Tìm kiếm/lọc sản phẩm theo từ khóa, category và khoảng giá.
- Thêm/sửa/xóa sản phẩm trong cart local trên trình duyệt.
- Đăng ký và đăng nhập.
- Khi tạo đơn từ cart, frontend yêu cầu đăng nhập.

## User/Customer

- Có toàn bộ khả năng của Guest.
- Xem profile hiện tại qua `/auth/me`.
- Tạo đơn từ cart.
- Xem danh sách đơn của chính mình và chi tiết từng đơn của mình.
- Đăng xuất và revoke refresh token được gửi lên.
- Không có API để user tự sửa profile, đổi mật khẩu, hủy đơn hoặc cập nhật trạng thái đơn.

## Admin

- Xem dashboard tổng hợp từ các API hiện có.
- Tạo, sửa, soft delete category và product.
- Xem, sửa role/name/email và soft delete user.
- Xem tất cả order, xem chi tiết và cập nhật trạng thái đúng luồng nghiệp vụ.
- Không có module riêng cho payment, voucher, upload/media hoặc email.

# 3. Mục đích sử dụng

Mục đích chính của project là demo một backend API có kiến trúc Handler → Service → Repository → PostgreSQL và một frontend đủ để chạy các luồng nghiệp vụ cốt lõi:

- [x] Guest xem sản phẩm.
- [x] User đăng ký, đăng nhập và refresh/logout token.
- [x] User thêm cart local, tạo đơn và xem đơn của mình.
- [x] Admin quản lý category, product, user và order.
- [x] Backend có JWT auth và role-based authorization.
- [x] Backend tạo order bằng database transaction, khóa product để kiểm tra/trừ stock và tự tính tổng tiền.
- [x] Frontend có public/user/admin route guard.
- [ ] Chưa có thanh toán, voucher, upload/media, email hoặc quy trình giao vận đầy đủ.

# 4. Các module hiện có trong backend

## Quy ước response chung

Success thông thường:

```json
{
  "success": true,
  "message": "Success",
  "data": {}
}
```

Danh sách có phân trang thêm `meta` gồm `page`, `limit`, `total`, `total_pages`. Error trả `success: false`, `message` và `errors`. Create trả HTTP 201; các thao tác khác chủ yếu trả 200.

## Auth — đã implement

| Method | Endpoint | Auth | Role | Request chính | Response `data` chính |
|---|---|---|---|---|---|
| POST | `/api/v1/auth/register` | Không | Public | `name`, `email`, `password` | `access_token`, `refresh_token`, `user` |
| POST | `/api/v1/auth/login` | Không | Public | `email`, `password` | `access_token`, `refresh_token`, `user` |
| POST | `/api/v1/auth/refresh-token` | Không cần access token | Public | `refresh_token` | Cặp token mới và `user`; token cũ bị revoke |
| POST | `/api/v1/auth/logout` | Có | user/admin | `refresh_token` | Message `Logged out successfully` |
| GET | `/api/v1/auth/me` | Có | user/admin | Không có body | User hiện tại |

`user` response gồm `id`, `name`, `email`, `role`, `is_active`, `created_at`, `updated_at`. Register luôn tạo role `user`; client không được chọn role.

## User — đã implement cho admin

| Method | Endpoint | Auth | Role | Request chính | Response `data` chính |
|---|---|---|---|---|---|
| GET | `/api/v1/users` | Có | admin | Query `page`, `limit`, `search` | Danh sách user active + pagination |
| GET | `/api/v1/users/:id` | Có | admin | Path `id` | Một user active |
| PUT | `/api/v1/users/:id` | Có | admin | `name`, `email`, `role` (`admin`/`user`) | User sau cập nhật |
| DELETE | `/api/v1/users/:id` | Có | admin | Path `id` | Message; soft delete bằng `is_active=false` |

Admin không thể tự delete chính mình. Không có endpoint admin tạo user. Không có endpoint user tự update profile hay đổi password.

## Profile — implement một phần

- `GET /api/v1/auth/me`: đã có, trả profile của tài khoản đăng nhập.
- Update profile/change password/avatar: **chưa implement**.

## Category — đã implement

| Method | Endpoint | Auth | Role | Request chính | Response `data` chính |
|---|---|---|---|---|---|
| GET | `/api/v1/categories` | Không | Public | Không có body | Danh sách category active |
| GET | `/api/v1/categories/:id` | Không | Public | Path `id` | Category active |
| POST | `/api/v1/categories` | Có | admin | `name`, `description`, `is_active` tùy chọn | Category vừa tạo |
| PUT | `/api/v1/categories/:id` | Có | admin | `name`, `description`, `is_active` tùy chọn | Category sau cập nhật |
| DELETE | `/api/v1/categories/:id` | Có | admin | Path `id` | Message; soft delete |

Category response gồm `id`, `name`, `description`, `is_active`. API list/detail dùng chung cho public nên chỉ đọc category active; chưa có admin endpoint để liệt kê cả inactive.

## Product — đã implement

| Method | Endpoint | Auth | Role | Request chính | Response `data` chính |
|---|---|---|---|---|---|
| GET | `/api/v1/products` | Không | Public | Query `page`, `limit`, `keyword`, `category_id`, `min_price`, `max_price` | Product active + pagination |
| GET | `/api/v1/products/:id` | Không | Public | Path `id` | Product active và category |
| POST | `/api/v1/products` | Có | admin | `category_id`, `name`, `description`, `price`, `stock`, `image_url`, `is_active` tùy chọn | Product vừa tạo |
| PUT | `/api/v1/products/:id` | Có | admin | Cùng cấu trúc create | Product sau cập nhật |
| DELETE | `/api/v1/products/:id` | Có | admin | Path `id` | Message; soft delete |

Product response gồm `id`, `category_id`, `name`, `description`, `price`, `stock`, `image_url`, `is_active`, `category`. `price` là integer (`BIGINT`), không dùng số thực. List/detail public chỉ trả product active thuộc category active; chưa có admin endpoint đọc cả inactive.

## Order — đã implement

| Method | Endpoint | Auth | Role | Request chính | Response `data` chính |
|---|---|---|---|---|---|
| POST | `/api/v1/orders` | Có | user/admin | `items: [{product_id, quantity}]` | Order vừa tạo |
| GET | `/api/v1/orders` | Có | user/admin | Không có body | Admin: mọi order; user: order của mình |
| GET | `/api/v1/users/me/orders` | Có | user/admin | Không có body | Order của user ID đang đăng nhập |
| GET | `/api/v1/orders/:id` | Có | user/admin | Path `id` | Admin xem bất kỳ; user chỉ xem order của mình |
| PUT | `/api/v1/orders/:id/status` | Có | admin | `status` | Order sau cập nhật status |

Order response hiện gồm `id`, `user_id`, `status`, `total_amount`, `items`. Mỗi item gồm `product_id`, `name`, `quantity`, `unit_price`, `subtotal`. DTO **không trả** `created_at`, `updated_at` hoặc thông tin chi tiết user, dù bảng database có timestamp; một số cột ngày/user trên frontend vì thế hiển thị `N/A` hoặc chỉ hiện user ID.

Luồng status thật trong service:

```text
pending -> confirmed
pending -> cancelled
confirmed -> shipping
confirmed -> cancelled
shipping -> completed
```

Không có chuyển ngược trạng thái. Khi cancel, code hiện không hoàn stock. User không có endpoint tự cancel.

## Health — đã implement

| Method | Endpoint | Auth | Response |
|---|---|---|---|
| GET | `/health` | Không | `{ "status": "ok" }` trong response envelope |

## Các module chưa implement

- Upload/media: **chưa implement**; product chỉ nhận chuỗi `image_url`.
- Payment: **chưa implement**; không có table, model, route hoặc UI quản lý payment.
- Voucher: **chưa implement**.
- Email: **chưa implement**.
- Staff/Manager: **chưa implement** và không có role tương ứng.

# 5. Database chạy ở đâu khi local?

## Cấu hình runtime chính ở root repo

Database chạy bằng root Docker Compose:

- Service name/hostname nội bộ Compose: `postgres`.
- Container name: `enterprise-order-postgres`.
- Image: `postgres:18.3-alpine3.23`.
- Host port: `5432` → container port `5432`.
- Database: `enterprise_order_management`.
- User: `postgres`.
- Password local demo: `postgres`.
- Volume: `postgres_data` mount tại `/var/lib/postgresql`.
- Migration root `./migrations` được mount vào `/docker-entrypoint-initdb.d` và chỉ tự chạy khi volume database được khởi tạo lần đầu.

Root Compose không dùng các biến `DB_USER`, `DB_PASSWORD`, `DB_NAME`. Nó khai báo trực tiếp `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, rồi truyền cho API một `DATABASE_URL` hoàn chỉnh.

Khi backend chạy trực tiếp trên host bằng `go run ./cmd/api`, `.env.example` dùng:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/enterprise_order_management?sslmode=disable
```

Do đó DB host phải là `localhost`.

Khi backend chạy trong root Compose, host phải là service name `postgres`:

```env
DATABASE_URL=postgres://postgres:postgres@postgres:5432/enterprise_order_management?sslmode=disable
```

## Cấu hình trong thư mục `backend/`

`backend/docker-compose.yml` dùng `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` lấy từ file env (`backend/.env` theo mặc định). Đây là module/skeleton riêng hiện chỉ có health route, không phải runtime nghiệp vụ chính. Không nên trộn bộ biến `DB_*` này với backend root vốn đọc `DATABASE_URL`.

# 6. Xem data trong database như thế nào?

## Dùng psql trong container

Từ root repo, vào interactive psql:

```powershell
docker compose exec postgres psql -U postgres -d enterprise_order_management
```

Một số lệnh psql hữu ích:

```sql
\dt
\d roles
\d users
\d refresh_tokens
\d categories
\d products
\d orders
\d order_items
```

Có thể chạy thẳng từng câu lệnh từ PowerShell:

```powershell
docker compose exec postgres psql -U postgres -d enterprise_order_management -c "SELECT * FROM roles ORDER BY id;"
docker compose exec postgres psql -U postgres -d enterprise_order_management -c "SELECT id, full_name, email, role_id, is_active, created_at FROM users ORDER BY id;"
docker compose exec postgres psql -U postgres -d enterprise_order_management -c "SELECT id, user_id, expires_at, revoked_at, created_at FROM refresh_tokens ORDER BY id DESC;"
docker compose exec postgres psql -U postgres -d enterprise_order_management -c "SELECT * FROM categories ORDER BY id;"
docker compose exec postgres psql -U postgres -d enterprise_order_management -c "SELECT id, category_id, name, price, stock, image_url, is_active FROM products ORDER BY id;"
docker compose exec postgres psql -U postgres -d enterprise_order_management -c "SELECT * FROM orders ORDER BY id DESC;"
docker compose exec postgres psql -U postgres -d enterprise_order_management -c "SELECT * FROM order_items ORDER BY order_id, id;"
```

Không nên select hoặc chia sẻ `password_hash` và `token_hash` khi không cần thiết.

Query xem order đầy đủ:

```sql
SELECT
    o.id AS order_id,
    u.email,
    o.status,
    o.total_amount,
    oi.product_id,
    p.name AS product_name,
    oi.quantity,
    oi.unit_price,
    oi.subtotal,
    o.created_at
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id
ORDER BY o.id DESC, oi.id;
```

Kiểm tra các bảng chưa có trong schema:

```sql
SELECT to_regclass('public.payments') AS payments;
SELECT to_regclass('public.vouchers') AS vouchers;
SELECT to_regclass('public.uploads') AS uploads;
SELECT to_regclass('public.media') AS media;
```

Các câu trên hiện trả `NULL`, vì migration chỉ tạo 7 bảng: `roles`, `users`, `refresh_tokens`, `categories`, `products`, `orders`, `order_items`.

## Dùng DBeaver, TablePlus hoặc pgAdmin

Dùng thông tin kết nối:

```text
Host: localhost
Port: 5432
Database: enterprise_order_management
Username: postgres
Password: postgres
SSL mode: disable
```

Container phải đang chạy và host port 5432 không bị ứng dụng PostgreSQL khác chiếm.

# 7. Frontend cần thể hiện những chức năng nào?

## Guest pages

- [x] Home/Product List (`/`): gọi `GET /products`, có keyword, category và pagination; backend còn hỗ trợ min/max price dù UI hiện chưa gửi hai filter này.
- [x] Product Detail (`/products/:id`): gọi `GET /products/:id`, thêm vào cart local.
- [x] Login (`/login`): gọi `POST /auth/login`.
- [x] Register (`/register`): gọi `POST /auth/register`.
- [x] Cart (`/cart`): hiện là public page, dữ liệu localStorage; chưa gọi backend cho đến khi tạo order.

## User pages

- [x] Profile (`/profile`): chỉ đọc user hiện tại; thiếu backend API để edit profile/change password.
- [x] Cart (`/cart`): thêm/xóa/sửa quantity local.
- [~] Checkout: chưa có route/page riêng. Nút “Tạo đơn hàng” trong Cart gọi trực tiếp `POST /orders`. Không có địa chỉ, shipping hay payment để tạo checkout đầy đủ.
- [x] My Orders (`/my-orders`): gọi `GET /users/me/orders`.
- [x] Order Detail (`/orders/:id`): gọi `GET /orders/:id`.

## Admin pages

- [x] Dashboard (`/admin`): tổng hợp client-side từ category/product/order/user APIs; chưa có dashboard/statistics endpoint riêng.
- [x] Category Management (`/admin/categories`): create/update/delete/list.
- [x] Product Management (`/admin/products`): create/update/delete/list.
- [x] Order Management (`/admin/orders`, `/admin/orders/:id`): list/detail/update status.
- [x] User Management (`/admin/users`): list/search/update/soft delete; không có create user.
- [ ] Voucher Management: thiếu backend API và frontend page.
- [ ] Payment Management: thiếu backend API và frontend page.
- [ ] Media/Upload Management: thiếu backend API và frontend page; admin nhập `image_url` thủ công.

Các giới hạn cần phản ánh đúng trên UI:

- Order response không có `created_at` và object user, nên UI ngày tạo/customer chi tiết chưa có đủ dữ liệu từ API.
- Admin category/product list dùng public endpoint, nên không tải được record inactive; filter inactive ở Product Management chỉ lọc trên tập active do backend trả về.
- Cart không được lưu trong database và không đồng bộ giữa thiết bị.
- Access/refresh token đang lưu localStorage để demo; production nên chuyển refresh token sang cookie `httpOnly`, `secure`, `sameSite`.

# 8. Kết luận scope hiện tại

## Đã đủ để làm và demo frontend ngay

- [x] Public product list/detail và category list.
- [x] Register, login, refresh token, logout, auth guard và admin guard.
- [x] Cart local và tạo order từ cart.
- [x] User xem danh sách/chi tiết order của mình.
- [x] Admin CRUD mềm category/product.
- [x] Admin list/update/delete mềm user.
- [x] Admin xem mọi order và cập nhật status.
- [x] Dashboard cơ bản tổng hợp từ API hiện có.

## Nên bổ sung backend trước khi hoàn thiện frontend hiện tại

Ưu tiên cao, vẫn nằm trong scope quản lý sản phẩm/đơn hàng:

- [ ] Thêm `created_at`, `updated_at` và thông tin user cần thiết vào OrderResponse, vì frontend đã có vị trí hiển thị nhưng API chưa trả.
- [ ] Thêm admin list/detail cho category/product bao gồm inactive, hoặc cho phép filter `is_active`; hiện admin không thể xem lại record đã inactive/soft delete.
- [ ] Quyết định nghiệp vụ cancel order có hoàn stock hay không và implement transaction nếu có.
- [ ] Nếu Profile cần chỉnh sửa thật, thêm endpoint user tự update profile; nếu không, đổi tên màn hình thành “Thông tin tài khoản”.
- [ ] Nếu báo cáo yêu cầu “Checkout page”, cần xác định tối thiểu dữ liệu giao nhận. Schema hiện không có địa chỉ hoặc thông tin người nhận.

## Nên để Phase 2 vì vượt scope thực tập 2 tháng

- [ ] Payment gateway, payment callback/webhook, refund.
- [ ] Voucher/promotion engine.
- [ ] Upload service, object storage và media management.
- [ ] Email xác nhận đơn/forgot password.
- [ ] Shipping provider, phí vận chuyển và tracking.
- [ ] Staff/Manager role, permission matrix chi tiết.
- [ ] Inventory ledger, nhập kho/xuất kho độc lập và audit log.
- [ ] Production-grade token storage/cookie, CSRF strategy và session/device management.

## Có thể mock để demo local

- [x] Dashboard statistics có thể tiếp tục tính từ các API list hiện có với lượng dữ liệu demo nhỏ.
- [x] Product image có thể dùng URL ảnh tĩnh/public thay vì upload.
- [x] Payment có thể hiển thị nhãn “Thanh toán khi nhận hàng / chưa tích hợp” ở mức UI, nhưng không được mô tả là payment module đã implement.
- [x] Shipping/address có thể dùng nội dung tĩnh trong prototype UI nếu cần trình bày, nhưng không gửi/lưu được bằng backend hiện tại.
- [x] Seed data local có thể dùng admin/category có sẵn từ migration và tạo thêm product/user/order qua API.

## Trạng thái kiểm tra kỹ thuật tại thời điểm phân tích

- [x] `go test ./...` chạy thành công.
- [x] `npm run build` chạy thành công.
- [x] `docker compose config --quiet` chạy thành công.
- [x] PostgreSQL container đang healthy và API container đang chạy.
- [~] `npm run lint` không có error nhưng có 7 warning về dependency của React Hooks.
- [!] Worktree đã có thay đổi chưa commit từ trước khi tạo tài liệu; tài liệu này không sửa các file code đó.

## Định nghĩa scope ngắn gọn đề xuất

> Enterprise Order Management là hệ thống web demo quản lý danh mục, sản phẩm, tồn kho cơ bản, người dùng và vòng đời đơn hàng. Backend Golang/Echo cung cấp REST API, JWT authentication, phân quyền user/admin và transaction PostgreSQL bằng SQL thuần; frontend React/Vite cung cấp giao diện cửa hàng và quản trị. Project tập trung vào nghiệp vụ backend và không tuyên bố là nền tảng e-commerce hoàn chỉnh.
