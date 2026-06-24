# PURCHASE & PAYMENT AUDIT

> Audit date: 2026-06-24  
> Source of truth: `cmd/api`, `internal/*`, `migrations/*`, root `Dockerfile`, root `docker-compose.yml`, `frontend/src/*`, `tonghop.md`  
> Out of scope as runtime source: `backend/`

## 1. Flow mua hàng hiện tại

Luồng hiện tại của project đang là:

1. Guest hoặc user xem sản phẩm qua `GET /api/v1/products` và `GET /api/v1/products/:id`.
2. Frontend thêm sản phẩm vào cart trong `CartContext`.
3. Cart chỉ tồn tại ở `localStorage`, không có cart API backend.
4. Khi user bấm tạo đơn ở `CartPage`, frontend gọi `POST /api/v1/orders` với payload chỉ gồm `items[{ product_id, quantity }]`.
5. Backend `OrderService.Create`:
   - gộp quantity theo `product_id`
   - khóa sản phẩm bằng `FindProductForUpdate`
   - lấy `price` từ DB
   - kiểm tra `is_active`
   - kiểm tra stock
   - tính `subtotal` và `total_amount`
   - tạo `orders`
   - tạo `order_items`
   - trừ stock trong transaction
6. Sau khi tạo order thành công, frontend xóa cart local và chuyển sang trang order detail hoặc danh sách order của tôi.

Kết luận:

- Backend đã là nguồn tính tiền cuối cùng cho order core.
- Flow hiện tại chưa có voucher, shipping, payment, checkout aggregate, hay payment status.
- Cart hiện tại mới phù hợp cho demo order cơ bản.

## 2. Cart hiện lưu ở đâu

Cart hiện được lưu hoàn toàn ở frontend:

- File chính: `frontend/src/contexts/CartContext.jsx`
- Storage: `localStorage`
- Mỗi item có dạng:
  - `product`
  - `quantity`

Đặc điểm hiện tại:

- Có sanitize dữ liệu từ `localStorage`
- Có clamp quantity theo stock đã biết trên frontend
- Có `addToCart`, `updateQuantity`, `removeFromCart`, `clearCart`
- `totalAmount` hiện chỉ là estimate phía client

Thiếu:

- Không có quote chuẩn từ backend
- Không có voucher
- Không có shipping estimate
- Không có sync cart phía server
- Không có warning chuẩn từ backend cho inactive/out-of-stock ngoài lúc submit order

## 3. Create order hiện nhận gì

Backend DTO hiện tại:

- File: `internal/dto/order.go`
- Request:

```json
{
  "items": [
    { "product_id": 1, "quantity": 2 }
  ]
}
```

Hiện chưa nhận:

- `voucher_code`
- `shipping_info`
- `shipping_fee`
- `payment_method`
- `payment_amount`
- `receiver_name`
- `receiver_phone`
- địa chỉ giao hàng

Response hiện tại:

- `id`
- `user_id`
- `status`
- `total_amount`
- `created_at`
- `updated_at`
- `user`
- `items`

Hiện chưa có trong response:

- `subtotal_amount`
- `discount_amount`
- `shipping_fee`
- `final_amount`
- `voucher_code`
- `payment`
- `shipping_info`

## 4. Thiếu gì so với yêu cầu mới

### Auth / user

- `users` chưa có `phone`
- chưa có `phone_verified_at`
- register chỉ cho email bắt buộc
- login chỉ cho email, chưa có `identifier`
- admin user list/update chưa hỗ trợ phone

### Cart / quote

- chưa có `POST /api/v1/cart/quote`
- frontend cart chưa dùng backend để tính subtotal/discount/shipping/final
- chưa có warnings chuẩn cho item inactive/hết hàng

### Voucher

- chưa có bảng `vouchers`
- chưa có bảng `voucher_usages`
- chưa có admin voucher API
- chưa có validate/apply/remove voucher flow
- `orders` chưa có cột liên quan voucher/tổng tiền chi tiết

### Shipping

- chưa có shipping snapshot trong order
- chưa có bảng shipping info riêng
- chưa có địa giới hành chính
- chưa có shipping quote API
- chưa có shipping rules/admin UI

### Payment

- chưa có bảng `payments`
- chưa có payment methods config endpoint
- chưa có COD/bank QR flow
- chưa có payment status hiển thị trong order
- chưa có admin payment management

### Checkout aggregate

- chưa có `POST /api/v1/checkout`
- frontend hiện gọi thẳng `POST /orders` từ `CartPage`
- chưa có transaction thống nhất cho voucher + shipping + payment + order

### Admin UI

- `AdminOrdersPage` và `AdminOrderDetailPage` chỉ quản lý order status cơ bản
- chưa có payment status
- chưa có shipping info
- chưa có voucher info

## 5. API mới cần thêm

### Nên thêm ở mức MVP an toàn

- `POST /api/v1/cart/quote`
- `POST /api/v1/checkout`
- `POST /api/v1/shipping/quote`
- `GET /api/v1/payments/methods`
- `POST /api/v1/payments/create`
- `GET /api/v1/payments/:id`
- `GET /api/v1/orders/:id/payment`
- `POST /api/v1/payments/:id/cancel`
- `PUT /api/v1/admin/payments/:id/status`

### Voucher

- `POST /api/v1/vouchers/validate`
- `GET /api/v1/admin/vouchers`
- `GET /api/v1/admin/vouchers/:id`
- `POST /api/v1/admin/vouchers`
- `PUT /api/v1/admin/vouchers/:id`
- `DELETE /api/v1/admin/vouchers/:id`

### Shipping rules

- `GET /api/v1/admin/shipping-rules`
- `POST /api/v1/admin/shipping-rules`
- `PUT /api/v1/admin/shipping-rules/:id`
- `DELETE /api/v1/admin/shipping-rules/:id`

### Auth backward-compatible expansion

Giữ endpoint cũ:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`

Nhưng mở rộng payload:

- register nhận `email` optional, `phone` optional, bắt buộc ít nhất một trong hai
- login nhận `identifier` thay vì `email`, hoặc hỗ trợ cả hai field trong giai đoạn tương thích

Khuyến nghị tương thích ngược:

- Login DTO nên hỗ trợ cả:
  - `identifier`
  - `email`
- Nếu `identifier` rỗng mà `email` có thì fallback sang `email`

## 6. Migration mới cần thêm

### Auth / user

1. Migration thêm vào `users`:
   - `phone VARCHAR(...) UNIQUE NULL`
   - `phone_verified_at TIMESTAMPTZ NULL`
2. Cần nới ràng buộc `email NOT NULL` nếu cho phép đăng ký chỉ bằng SĐT.

Lưu ý:

- Đây là thay đổi schema nhạy cảm vì `users.email` hiện đang `NOT NULL UNIQUE`.
- Cần migration cẩn thận để không phá dữ liệu cũ.

### Voucher

3. Tạo `vouchers`
4. Tạo `voucher_usages`
5. Mở rộng `orders`:
   - `voucher_code`
   - `discount_amount`
   - `subtotal_amount`
   - `shipping_fee`
   - `final_amount`

### Shipping

6. Tạo bảng snapshot shipping, khuyến nghị:
   - `order_shipping_infos`

7. Tạo bảng:
   - `shipping_fee_rules`

### Payment

8. Tạo bảng:
   - `payments`

### Optional later

9. `inventory_logs` nếu muốn audit stock movement

## 7. Frontend page/component cần thêm hoặc sửa

### Cần sửa

- `frontend/src/pages/auth/RegisterPage.jsx`
- `frontend/src/pages/auth/LoginPage.jsx`
- `frontend/src/pages/user/CartPage.jsx`
- `frontend/src/pages/user/OrderDetailPage.jsx`
- `frontend/src/pages/admin/AdminOrdersPage.jsx`
- `frontend/src/pages/admin/AdminOrderDetailPage.jsx`
- `frontend/src/api/orderApi.js`
- `frontend/src/api/userApi.js`
- `frontend/src/contexts/CartContext.jsx`

### Cần thêm

- `frontend/src/pages/user/CheckoutPage.jsx`
- `frontend/src/pages/public/PaymentResultPage.jsx` hoặc page tương đương
- `frontend/src/pages/admin/AdminVouchersPage.jsx`
- `frontend/src/pages/admin/AdminPaymentsPage.jsx`
- `frontend/src/pages/admin/AdminShippingRulesPage.jsx`
- `frontend/src/components/cart/VoucherInput.jsx`
- `frontend/src/components/checkout/ShippingForm.jsx`
- `frontend/src/components/checkout/PaymentMethodSelector.jsx`

### Data / utilities

- `frontend/src/data/vietnam_locations.json`
- helper format payment/shipping/voucher status nếu cần

## 8. Rủi ro tương thích ngược

### Cao

- Đổi `users.email` từ bắt buộc sang optional ảnh hưởng auth hiện tại, token generation, admin flows, Google OAuth assumptions, seed data.
- Đổi request login từ `email` sang `identifier` có thể làm frontend cũ lỗi nếu không giữ fallback tương thích.
- Mở rộng `orders` dễ làm frontend/admin pages cũ hiển thị sai nếu code đang giả định chỉ có `total_amount`.

### Trung bình

- Thêm `/checkout` nhưng vẫn giữ `/orders` cần thống nhất business rules để không có hai luồng lệch nhau.
- Voucher usage và payment state cần transaction/idempotency rõ, nếu không dễ double count.
- Bank QR manual confirmation cần admin flow rõ để tránh user tự đánh dấu paid.

### Thấp

- Cart quote endpoint có thể thêm song song mà không phá flow cũ.
- Shipping quote nội bộ có thể thêm dần.

## 9. Kế hoạch triển khai từng phase

### Phase 1

Mở rộng auth theo hướng backward-compatible:

- thêm `phone`
- register nhận email hoặc phone
- login hỗ trợ `identifier` và fallback `email`
- cập nhật admin user UI/API

### Phase 2

Giữ cart ở `localStorage`, thêm backend quote:

- `POST /api/v1/cart/quote`
- CartPage gọi quote mỗi khi cart đổi
- chưa sync cart server để giảm scope

### Phase 3

Thêm voucher core:

- migration vouchers + voucher_usages
- admin CRUD voucher
- validate voucher
- tích hợp vào quote + checkout

### Phase 4

Thêm shipping info snapshot:

- `order_shipping_infos`
- `CheckoutPage`
- render shipping info trong order detail/admin

### Phase 5

Thêm shipping quote nội bộ:

- `shipping_fee_rules`
- `POST /api/v1/shipping/quote`
- CheckoutPage gọi quote

### Phase 6

MVP payment an toàn:

- COD luôn enable
- Bank QR manual enable nếu có env
- payment methods config endpoint
- `payments` table
- admin confirm bank payment

Provider gateway thật như MoMo/ZaloPay/VNPAY/OnePay:

- chỉ để abstraction + disabled state nếu chưa có env
- không nên fake paid

### Phase 7

Gộp flow qua `POST /api/v1/checkout`:

- validate items
- quote subtotal
- apply voucher
- calculate shipping
- create order
- create order items
- deduct stock
- create payment
- increment voucher usage

### Phase 8

Admin UI:

- orders có shipping/payment/voucher quick view
- admin vouchers
- admin payments
- admin shipping rules

### Phase 9

Docs + seed + test:

- update `tonghop.md`
- update `README.md`
- create implementation report
- run test/build/smoke test

## 10. Khuyến nghị phạm vi triển khai an toàn

Để giữ tính demo-ready và không phá core hiện có, thứ tự an toàn nhất là:

1. Phase 1
2. Phase 2
3. Phase 3
4. Phase 4
5. Phase 5 với internal shipping engine
6. Phase 6 chỉ enable:
   - `cod`
   - `bank_qr` manual
   - các gateway khác chỉ disabled/env-gated
7. Phase 7 dùng `/checkout`
8. Phase 8 admin pages mức core

Không khuyến nghị trong đợt này:

- lưu thông tin thẻ
- fake webhook paid
- production claim cho payment gateway thật khi chưa có sandbox verify
- server-side cart sync nếu chưa cần

## 11. Trạng thái audit kết luận

Repo hiện tại đã có nền tốt cho order core, response envelope thống nhất, stock transaction cơ bản, UI admin/user/public đủ để mở rộng dần. Tuy nhiên khối lượng thay đổi của prompt là lớn và chạm cả auth, schema users, vouchers, shipping, payment, checkout aggregate, admin UI và docs.

Đợt triển khai phù hợp nhất là theo hướng:

- giữ tương thích ngược cho auth/order cũ
- thêm endpoint mới thay vì phá endpoint đang có
- ưu tiên demo payment an toàn bằng `cod` và `bank_qr` manual
- giữ provider thật ở trạng thái disabled/env-gated cho đến khi có credential sandbox hợp lệ
