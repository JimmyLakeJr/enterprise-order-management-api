# PURCHASE PAYMENT QA REPORT

> Audit date: 2026-06-24  
> Source of truth: `cmd/api`, `internal/*`, `migrations/*`, `frontend/src/*`, `frontend/package.json`, root `docker-compose.yml`, `tonghop.md`

## 1. Final verdict

Final verdict: **FAIL for purchase/payment scope**.

Repo hiện tại vẫn là core product/user/order management system, chưa phải implementation hoàn chỉnh cho nhóm Mua hàng & Thanh toán theo checklist verify này.

Điểm tích cực:

- Core auth email/password hiện có
- Core cart localStorage hiện có
- Core order create/list/detail/status hiện có
- Backend vẫn là nguồn tính `total_amount` khi tạo order
- Test/build nền tảng đều pass

Điểm quyết định fail:

- Chưa có Auth Email/SĐT
- Chưa có cart quote backend
- Chưa có voucher
- Chưa có shipping
- Chưa có payment
- Chưa có checkout aggregate
- Chưa có các admin pages cho vouchers/payments/shipping rules

## 2. Pass/fail từng nhóm

| Nhóm | Trạng thái | Ghi chú |
|---|---|---|
| 1. Auth Email/SĐT | Fail | Chỉ có email/password; không có phone/identifier |
| 2. Cart | Partial pass | Add/remove/update local có; quote/warning backend chưa có |
| 3. Voucher | Fail | Không có migration, API, service, UI |
| 4. Shipping | Fail | Không có schema, quote, snapshot, UI |
| 5. Payment | Fail | Không có bảng `payments`, không có payment API/UI |
| 6. Checkout | Fail | Chưa có `/checkout`, chưa có voucher/shipping/payment transaction |
| 7. Frontend | Partial fail | CartPage/OrderDetail/AdminOrderDetail có; các page checkout/payment/voucher/shipping thiếu |
| 8. Security | Partial pass | Role guard/ownership/order transaction có; payment-specific security chưa áp dụng vì chưa có payment module |

## 3. Blocker

Các blocker ngăn hệ thống đạt trạng thái pass cho nhóm Mua hàng & Thanh toán:

1. Không có migration/backend/frontend cho `vouchers`.
2. Không có migration/backend/frontend cho `payments`.
3. Không có shipping schema, shipping quote, shipping snapshot trong order.
4. Không có `POST /api/v1/cart/quote`.
5. Không có `POST /api/v1/checkout`.
6. Auth chưa hỗ trợ phone/SĐT ở schema, DTO, service và UI.
7. Frontend không có `CheckoutPage`, `PaymentResultPage`, `AdminVouchersPage`, `AdminPaymentsPage`, `AdminShippingRulesPage`.

## 4. High priority bugs

Các lỗi/gap ưu tiên cao theo mục tiêu prompt:

| Mức | Vấn đề | Bằng chứng |
|---|---|---|
| High | Register không hỗ trợ phone | `internal/dto/auth.go`: `RegisterRequest` chỉ có `name`, `email`, `password` |
| High | Login không hỗ trợ phone/identifier | `internal/dto/auth.go`: `LoginRequest` chỉ có `email`, `password`; `AuthService.Login` chỉ gọi `FindByEmail` |
| High | `users` chưa có `phone` | `migrations/001_init.sql`, `002_google_oauth.sql`, `003_profile_media.sql` không thêm `phone` |
| High | Không có cart quote backend | `internal/http/server.go` không có route `/api/v1/cart/quote` |
| High | CartPage dùng giá tạm tính ở frontend | `frontend/src/pages/user/CartPage.jsx` dùng `getCartTotal()` từ `CartContext` |
| High | Không có voucher module | Không có migration/API/page/service chứa `voucher` ngoài tài liệu audit |
| High | Không có shipping module | Không có migration/API/page/service chứa `shipping` ngoài order status/UI text |
| High | Không có payment module | Không có `payments` table trong `tonghop.md`; không có route/service/repository payment |
| High | Không có checkout aggregate | `frontend/src/routes/AppRoutes.jsx` không có route checkout; backend không có `/api/v1/checkout` |
| High | Admin không có payment/voucher/shipping management | `AppRoutes.jsx` chỉ có categories/products/orders/users |

## 5. Medium priority improvements

| Mức | Vấn đề | Ghi chú |
|---|---|---|
| Medium | Cart chỉ lưu localStorage | Phù hợp demo core, chưa có quote chuẩn/backend sync |
| Medium | Cart không có warning inactive/out-of-stock từ backend trước khi submit order | Chỉ fail ở bước create order |
| Medium | Validation error keys vẫn theo tên field Go | `internal/pkg/response/response.go` dùng `fieldErr.Field()` thay vì JSON tag |
| Medium | Không có `npm test` cho frontend | `frontend/package.json` không có script `test` |
| Medium | `tonghop.md` đã ghi rõ payment/voucher/shipping chưa implement, nhưng scope verify này cần module thật để pass | Tài liệu đúng, nhưng chức năng chưa có |

## 6. Security risks

| Mức | Risk | Trạng thái |
|---|---|---|
| Medium | Token lưu ở client-side storage | Đã biết từ `tonghop.md`, phù hợp demo hơn production |
| Medium | CORS cấu hình một origin từ env, chưa thấy ma trận env rộng hơn | `docker-compose.yml` dùng `FRONTEND_URL=http://localhost:5173` |
| Low | SQL injection login/search | Thấp; repository dùng placeholder parameterized query |
| Low | No raw card data | Hiện pass vì hoàn toàn chưa có card flow |
| High | Payment secret / webhook verification | Chưa áp dụng vì payment provider chưa implement |
| High | User tự set paid | Hiện chưa có payment module; nếu triển khai thiếu role guard sẽ là blocker bảo mật |

## 7. Manual test checklist

### Auth Email/SĐT

- [x] Register bằng email
- [ ] Register bằng phone
- [x] Login bằng email
- [ ] Login bằng phone
- [x] Duplicate email bị chặn
- [ ] Duplicate phone bị chặn
- [ ] Missing email/phone bị chặn theo rule “ít nhất một trong hai”
- [x] Inactive user bị chặn gián tiếp do `FindByEmail` chỉ lấy `is_active = TRUE`

### Cart

- [x] Add item
- [x] Remove item
- [x] Update quantity
- [ ] Quote subtotal từ backend
- [ ] Out-of-stock warning từ quote backend
- [ ] Product inactive warning từ quote backend
- [x] Cart localStorage hoạt động
- [ ] Backend compatibility cho quote/checkout

### Voucher

- [ ] Admin CRUD
- [ ] Validate voucher
- [ ] Apply voucher trong cart
- [ ] Apply voucher trong checkout
- [ ] Usage count
- [ ] Per-user limit
- [ ] Percent/fixed discount
- [ ] Max discount
- [ ] Min order amount

### Shipping

- [ ] Province/district/ward selection
- [ ] Shipping quote
- [ ] Free shipping threshold
- [ ] Shipping snapshot in order
- [ ] Order detail shipping info
- [ ] Admin shipping info

### Payment

- [ ] COD
- [ ] Bank QR
- [ ] MoMo disabled theo config thật
- [ ] ZaloPay disabled theo config thật
- [ ] VNPAY disabled theo config thật
- [ ] OnePay disabled theo config thật
- [ ] Credit card disabled theo provider config thật
- [ ] Admin confirm bank payment
- [ ] User không tự set paid
- [ ] Payment amount lấy từ backend order amount
- [ ] Webhook signature nếu có provider

### Checkout

- [ ] Checkout gom items/voucher/shipping/payment
- [x] Backend tự tính total order core
- [ ] Backend tự tính subtotal/discount/shipping/final đầy đủ
- [x] Backend trừ stock khi create order core
- [ ] Backend tạo payment
- [ ] Rollback khi payment init/voucher/shipping lỗi

## 8. API checklist

### API có thật

- [x] `POST /api/v1/auth/register`
- [x] `POST /api/v1/auth/login`
- [x] `POST /api/v1/auth/refresh-token`
- [x] `POST /api/v1/auth/logout`
- [x] `GET /api/v1/auth/me`
- [x] `POST /api/v1/orders`
- [x] `GET /api/v1/orders`
- [x] `GET /api/v1/orders/:id`
- [x] `PUT /api/v1/orders/:id/status`
- [x] `GET /api/v1/users/me/orders`

### API thiếu so với verify scope

- [ ] `POST /api/v1/cart/quote`
- [ ] `POST /api/v1/vouchers/validate`
- [ ] `GET /api/v1/admin/vouchers`
- [ ] `POST /api/v1/admin/vouchers`
- [ ] `PUT /api/v1/admin/vouchers/:id`
- [ ] `DELETE /api/v1/admin/vouchers/:id`
- [ ] `POST /api/v1/shipping/quote`
- [ ] `GET /api/v1/admin/shipping-rules`
- [ ] `POST /api/v1/admin/shipping-rules`
- [ ] `PUT /api/v1/admin/shipping-rules/:id`
- [ ] `DELETE /api/v1/admin/shipping-rules/:id`
- [ ] `GET /api/v1/payments/methods`
- [ ] `POST /api/v1/payments/create`
- [ ] `GET /api/v1/payments/:id`
- [ ] `GET /api/v1/orders/:id/payment`
- [ ] `POST /api/v1/payments/:id/cancel`
- [ ] `PUT /api/v1/admin/payments/:id/status`
- [ ] `POST /api/v1/payments/webhook/:provider`
- [ ] `POST /api/v1/checkout`

## 9. Frontend route checklist

### Route có thật

- [x] `/`
- [x] `/products`
- [x] `/products/:id`
- [x] `/cart`
- [x] `/login`
- [x] `/register`
- [x] `/my-orders`
- [x] `/orders/:id`
- [x] `/profile`
- [x] `/admin`
- [x] `/admin/categories`
- [x] `/admin/products`
- [x] `/admin/orders`
- [x] `/admin/orders/:id`
- [x] `/admin/users`

### Route thiếu cho verify scope

- [ ] `/checkout`
- [ ] payment result route
- [ ] `/admin/vouchers`
- [ ] `/admin/payments`
- [ ] `/admin/shipping-rules`

## 10. Có thể demo chưa

**Có thể demo core order management**, nhưng **chưa thể demo purchase/payment scope hoàn chỉnh**.

Demo được:

- auth email/password
- product list/detail
- cart localStorage
- create order core
- order detail
- admin order lifecycle

Không demo được theo scope verify này:

- Email/SĐT auth
- voucher
- shipping
- payment
- checkout aggregate

## 11. Có thể beta chưa

**Chưa thể beta** cho tính năng Mua hàng & Thanh toán.

Lý do:

- chưa có voucher/shipping/payment
- chưa có checkout transaction aggregate
- chưa có payment security model, webhook, admin confirm flow
- chưa có frontend routes/pages cho purchase flow đầy đủ

## 12. Có thể production chưa

**Chưa production-ready**.

Lý do:

- feature scope chưa hoàn tất
- payment module chưa tồn tại
- shipping module chưa tồn tại
- auth email/SĐT chưa tồn tại
- frontend chưa có automated test
- payment provider, signature verification, env-hardening, audit log, rollback edge cases chưa có

## 13. Test/build result

| Lệnh | Kết quả |
|---|---|
| `go test ./cmd/... ./internal/...` | Pass |
| `go vet ./cmd/... ./internal/...` | Pass |
| `docker compose config --quiet` | Pass |
| `cd frontend && npm run build` | Pass |
| `cd frontend && npm run lint` | Pass |
| `cd frontend && npm test` | Không có script `test` trong `frontend/package.json` |

## 14. Patch đề xuất sau audit

Chưa sửa code ở bước này. Các patch nhỏ, an toàn và nên làm tiếp sau khi được xác nhận:

1. Auth backward-compatible:
   - thêm `phone` vào `users`
   - giữ endpoint cũ nhưng mở rộng DTO login/register
   - login hỗ trợ cả `identifier` và fallback `email`

2. Cart quote an toàn:
   - thêm `POST /api/v1/cart/quote`
   - giữ `CartContext` localStorage hiện tại
   - `CartPage` gọi quote backend để lấy subtotal/warnings

3. Checkout MVP:
   - thêm `POST /api/v1/checkout`
   - gom `items`, `shipping_info`, `payment_method`, `voucher_code`
   - backend là nguồn tính cuối cùng

4. Payment MVP:
   - chỉ enable `cod` và `bank_qr`
   - các provider còn lại trả disabled theo config

5. Frontend pages tối thiểu:
   - `CheckoutPage`
   - `AdminVouchersPage`
   - `AdminPaymentsPage`
   - `AdminShippingRulesPage`

6. Response validation polish:
   - map validation key sang JSON tag thay vì tên field Go

## 15. Bằng chứng chính trong code

- Auth email-only:
  - `internal/dto/auth.go`
  - `internal/service/auth_service.go`
  - `frontend/src/pages/auth/LoginPage.jsx`
  - `frontend/src/pages/auth/RegisterPage.jsx`
- Order core only:
  - `internal/dto/order.go`
  - `internal/handler/order_handler.go`
  - `internal/service/order_service.go`
  - `internal/repository/order_repository.go`
- Cart local-only:
  - `frontend/src/contexts/CartContext.jsx`
  - `frontend/src/pages/user/CartPage.jsx`
- Routes thiếu purchase/payment:
  - `frontend/src/routes/AppRoutes.jsx`
- Payment/voucher/shipping chưa implement:
  - `migrations/*`
  - `internal/http/server.go`
  - `tonghop.md`
