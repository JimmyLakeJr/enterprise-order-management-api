# API Reference

Base URL local: `http://localhost:8080/api/v1`.

Protected endpoint dùng header:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

Response thành công dùng `data`; lỗi dùng thông điệp và mã lỗi theo helper chung của backend. Client phải dựa vào HTTP status, không suy diễn dữ liệu khi request thất bại.

## Health

| Method | Path | Quyền | Mục đích |
|---|---|---|---|
| `GET` | `/health` | Public | Kiểm tra API runtime; path không có prefix `/api/v1` |

## Auth và profile

| Method | Path | Quyền | Nội dung chính |
|---|---|---|---|
| `POST` | `/auth/register` | Public | Đăng ký bằng name, email, password |
| `POST` | `/auth/login` | Public | Nhận access token và refresh token |
| `POST` | `/auth/refresh-token` | Public | Đổi refresh token hợp lệ lấy token mới |
| `POST` | `/auth/logout` | User | Thu hồi refresh token theo payload hiện tại |
| `GET` | `/auth/me` | User | Lấy user hiện tại |
| `PUT` | `/users/me` | User | Cập nhật `name`; không cho sửa role/is_active |

Profile update tối thiểu:

```json
{ "name": "Nguyen Van A" }
```

Endpoint này chỉ nhận dữ liệu profile được khai báo trong DTO, không nhận file media.

## Category

| Method | Path | Quyền | Ghi chú |
|---|---|---|---|
| `GET` | `/categories` | Public | Chỉ category active |
| `GET` | `/categories/:id` | Public | Chỉ dữ liệu public hợp lệ |
| `POST` | `/categories` | Admin | Tạo category |
| `PUT` | `/categories/:id` | Admin | Cập nhật category |
| `DELETE` | `/categories/:id` | Admin | Soft delete/inactive theo nghiệp vụ hiện tại |
| `GET` | `/admin/categories?status=all\|active\|inactive` | Admin | Xem cả active/inactive |
| `PUT` | `/admin/categories/:id/restore` | Admin | Khôi phục category inactive |

Không thể làm category inactive khi còn product active nếu rule bảo vệ dữ liệu được kích hoạt.

## Product

| Method | Path | Quyền | Ghi chú |
|---|---|---|---|
| `GET` | `/products` | Public | Active only; hỗ trợ filter/pagination hiện có |
| `GET` | `/products/:id` | Public | Chi tiết product active |
| `POST` | `/products` | Admin | Tạo product |
| `PUT` | `/products/:id` | Admin | Cập nhật product |
| `DELETE` | `/products/:id` | Admin | Soft delete/inactive |
| `GET` | `/admin/products?status=all\|active\|inactive` | Admin | Danh sách quản trị; giữ filter/pagination |
| `PUT` | `/admin/products/:id/restore` | Admin | Restore khi category đang active |

Field ảnh là `image_url`; API không upload file.

## Order

| Method | Path | Quyền | Ghi chú |
|---|---|---|---|
| `POST` | `/orders` | User | Tạo đơn, backend tính giá và giảm stock trong transaction |
| `GET` | `/orders` | User/Admin | User chỉ thấy đơn của mình; admin thấy toàn bộ; hỗ trợ `page`, `limit`, `status` |
| `GET` | `/orders/:id` | User/Admin | Kiểm tra ownership hoặc quyền admin |
| `GET` | `/users/me/orders` | User | Alias danh sách đơn của user hiện tại; hỗ trợ `page`, `limit`, `status` |
| `PUT` | `/orders/:id/status` | Admin | Cập nhật trạng thái hợp lệ |

Payload tạo đơn:

```json
{
  "items": [
    { "product_id": 1, "quantity": 2 }
  ]
}
```

Client không gửi giá. Backend đọc giá hiện tại, kiểm tra stock, tạo `orders`/`order_items` và trừ stock trong cùng transaction.

Order response có `created_at`, `updated_at`; danh sách order hiện có pagination, status filter và user summary khi phù hợp. Cancel chỉ đổi trạng thái, chưa hoàn stock.

## User admin

| Method | Path | Quyền | Ghi chú |
|---|---|---|---|
| `GET` | `/users` | Admin | Danh sách user |
| `GET` | `/users/:id` | Admin | Chi tiết user |
| `PUT` | `/users/:id` | Admin | Cập nhật field quản trị được phép |
| `DELETE` | `/users/:id` | Admin | Vô hiệu hóa/xóa theo service hiện tại |

## Status thường gặp

| Status | Ý nghĩa |
|---|---|
| `200` | Đọc/cập nhật thành công |
| `201` | Tạo thành công |
| `400` | Payload hoặc chuyển trạng thái không hợp lệ |
| `401` | Thiếu/sai token |
| `403` | Không đủ quyền |
| `404` | Không tìm thấy hoặc không được phép thấy resource |
| `409` | Xung đột dữ liệu, ví dụ email trùng hoặc stock không phù hợp |
| `500` | Lỗi nội bộ |

## Nguyên tắc client

- Refresh token qua endpoint auth, không tự dựng JWT.
- Public UI không dùng admin list để lộ dữ liệu inactive.
- Không gọi API Phase 2; danh sách module ngoài scope xem trong tài liệu scope hiện tại.
- Khi contract và tài liệu khác nhau, route trong `internal/http/server.go` cùng DTO/service root là source of truth.
