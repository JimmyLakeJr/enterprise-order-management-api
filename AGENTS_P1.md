QUY CHUẨN CHUNG CỦA PROJECT

Đề tài:
“Phát triển backend API cho hệ thống quản lý sản phẩm và đơn hàng trong doanh nghiệp sử dụng Golang”.

Trọng tâm:
- Backend API là phần chính của đồ án thực tập.
- Frontend chỉ là demo client đơn giản để kiểm thử API và trình bày nghiệp vụ.
- Project phải phù hợp với chương trình thực tập 2 tháng.
- Code cần dễ đọc, dễ bảo trì, rõ kiến trúc, phù hợp trình độ intern/junior backend.

Stack bắt buộc:
- Backend: Golang 1.22+
- Framework: Echo v4
- Database: PostgreSQL
- Database access: SQL thuần với pgx/v5 và pgxpool
- Không dùng ORM như GORM
- Authentication: JWT access token + refresh token
- JWT library: github.com/golang-jwt/jwt/v5
- Password hashing: golang.org/x/crypto/bcrypt
- Validation: github.com/go-playground/validator/v10
- Config: .env, godotenv, os.Getenv
- Migration: SQL migration thủ công hoặc golang-migrate
- Testing: testing + testify
- Frontend demo: React + Vite
- Deploy local: Docker Compose
- Deploy demo production: PostgreSQL trên Supabase/Neon, backend trên Render, frontend trên Vercel

Kiến trúc backend:
- Handler: nhận request, parse request, gọi service, trả response.
- Service: xử lý business logic.
- Repository: thao tác database bằng SQL thuần.
- Model: ánh xạ dữ liệu database.
- DTO: định nghĩa request/response.
- Middleware: JWT auth, role authorization, CORS, logger, recovery.
- Config: load biến môi trường.
- Database: khởi tạo pgxpool, transaction, close connection.
- Util: helper cho JWT, password, response, pagination.

Quy tắc code:
- Không viết SQL trong handler.
- Không viết business logic trong handler.
- Không hard-code secret, database URL, JWT secret.
- Không trả password_hash ra response.
- Không expose lỗi database thô cho client.
- Không nối chuỗi SQL trực tiếp với input người dùng.
- Phải dùng parameterized query.
- Phải dùng transaction cho chức năng tạo đơn hàng.
- Phải validate input trước khi xử lý.
- Phải trả response JSON thống nhất.
- Phải có .env.example.
- Phải có README hướng dẫn chạy local và deploy.
- Phải có Dockerfile và docker-compose.yml.
- Phải có API docs hoặc Postman collection.

Quy tắc bảo mật:
- Password phải hash bằng bcrypt.
- Không lưu plain password.
- Access token có thời hạn ngắn, ví dụ 15 phút.
- Refresh token có thời hạn dài hơn, ví dụ 7 ngày.
- Refresh token phải được hash trước khi lưu vào database.
- Không lưu raw refresh token trong database.
- Logout phải revoke refresh token.
- API admin phải có middleware kiểm tra role.
- User chỉ được xem đơn hàng của chính mình.
- Admin được xem toàn bộ đơn hàng.
- JWT secret phải lấy từ biến môi trường.
- CORS phải cấu hình theo frontend domain.
- Không commit file .env.

Quy tắc database:
- Database chính là PostgreSQL.
- Không dùng ORM.
- Dùng SQL thuần với pgxpool.
- Các bảng chính:
  - roles
  - users
  - refresh_tokens
  - categories
  - products
  - orders
  - order_items
- users phải có role_id.
- refresh_tokens phải lưu token_hash, expires_at, revoked_at.
- categories nên có is_active.
- products phải có category_id, price, stock, is_active.
- orders phải có user_id, total_amount, status.
- order_items phải có order_id, product_id, quantity, unit_price, subtotal.
- order_items phải lưu unit_price tại thời điểm đặt hàng.
- Backend tự tính total_amount, không tin price từ frontend.
- Product price và stock không được âm.
- Quantity trong order phải lớn hơn 0.
- Không cho tạo order rỗng.
- Product đã inactive không hiển thị ở API public.
- DELETE product/user/category ưu tiên soft delete bằng is_active = false.
Quy tắc nghiệp vụ order:
- User đăng nhập mới được tạo đơn hàng.
- Khi tạo order phải kiểm tra product tồn tại.
- Phải kiểm tra product còn active.
- Phải kiểm tra stock đủ.
- Phải lấy unit_price từ database tại thời điểm đặt hàng.
- Phải tạo order, order_items và trừ stock trong cùng một transaction.
- Nếu bất kỳ bước nào lỗi thì rollback toàn bộ.
- Không cho user xem order của người khác.
- Admin được xem tất cả order.
- Chỉ admin được cập nhật trạng thái order.
- Trạng thái order gồm:
  - pending
  - confirmed
  - shipping
  - completed
  - cancelled
- Luồng trạng thái hợp lệ:
  - pending -> confirmed
  - pending -> cancelled
  - confirmed -> shipping
  - confirmed -> cancelled
  - shipping -> completed
- Không cho chuyển trạng thái ngược hoặc sai nghiệp vụ.

Phân quyền:
Guest:
- Xem danh sách category.
- Xem danh sách product.
- Xem chi tiết product.

User:
- Có toàn bộ quyền của Guest.
- Đăng nhập, đăng xuất.
- Tạo đơn hàng.
- Xem đơn hàng của chính mình.
- Xem thông tin cá nhân.

Admin:
- Có toàn bộ quyền của User.
- Quản lý user.
- Quản lý category.
- Quản lý product.
- Xem toàn bộ order.
- Cập nhật trạng thái order.

Response format:

Success:
{
  "success": true,
  "message": "Success",
  "data": {}
}

Error:
{
  "success": false,
  "message": "Validation failed",
  "errors": {}
}

Pagination:
{
  "success": true,
  "message": "Success",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 0,
    "total_pages": 0
  }
}

Yêu cầu khi trả lời:
1. Trình bày rõ mục tiêu của bước.
2. Liệt kê file/thư mục cần tạo.
3. Viết code đầy đủ, dễ hiểu.
4. Giải thích ngắn gọn vai trò từng file.
5. Có lệnh chạy/test.
6. Có ví dụ curl hoặc Postman.
7. Có lỗi thường gặp và cách sửa.
8. Không tạo code quá phức tạp.
9. Ưu tiên tính đúng, rõ ràng, dễ bảo trì.
