# Database Design

PostgreSQL schema ban đầu nằm tại `migrations/001_init.sql`. Đây là schema demo, không phải thiết kế e-commerce đầy đủ.

## Bảng chính

| Bảng | Mục đích | Field/constraint đáng chú ý |
|---|---|---|
| `roles` | Danh mục vai trò | Giá trị role dùng cho phân quyền |
| `users` | Tài khoản | Email unique, password hash, role, `is_active`, timestamps |
| `refresh_tokens` | Phiên refresh | Token/user/expiry/revocation phục vụ refresh và logout |
| `categories` | Nhóm sản phẩm | Name, description, `is_active`, timestamps |
| `products` | Sản phẩm và stock cơ bản | Category FK, price, stock, `image_url`, `is_active`, timestamps |
| `orders` | Header đơn hàng | User FK, status, total, timestamps |
| `order_items` | Snapshot dòng hàng | Order/product FK, quantity, unit price/subtotal theo schema |

## Quan hệ

```text
roles 1 --- n users
users 1 --- n refresh_tokens
users 1 --- n orders
categories 1 --- n products
orders 1 --- n order_items
products 1 --- n order_items
```

Sơ đồ chi tiết xem [ERD.md](ERD.md).

## Quy tắc dữ liệu

- Email user không trùng.
- Price, stock và quantity không nhận giá trị âm/không hợp lệ theo constraint/service hiện có.
- Product luôn tham chiếu category hợp lệ.
- Order item lưu giá tại thời điểm đặt để lịch sử không đổi theo giá product sau này.
- Public query lọc `is_active`; inactive không đồng nghĩa xóa vật lý.
- Restore product yêu cầu category active.

## Transaction order

Tạo order cần là một transaction:

1. Kiểm tra user và danh sách item.
2. Đọc product, giá và stock hiện tại.
3. Tạo order và order items.
4. Trừ stock.
5. Commit; bất kỳ lỗi nào phải rollback.

Cancel hiện không cộng stock. Không bổ sung phép cộng đơn giản trước khi có quyết định về trạng thái được restock, retry, idempotency và concurrent update.

## Migration và dữ liệu local

- Migration init chạy khi PostgreSQL volume được tạo lần đầu qua Docker Compose.
- Sửa `001_init.sql` không tự cập nhật volume đang tồn tại.
- `docker compose down -v` xóa dữ liệu local; phải chủ động sao lưu nếu cần.
- Admin seed chỉ phục vụ local/demo; credential phải đổi ở môi trường công khai.

## Phần chưa có

Schema không chứa các module Phase 2 hoặc inventory ledger; không thêm chúng vào migration demo hiện tại.
