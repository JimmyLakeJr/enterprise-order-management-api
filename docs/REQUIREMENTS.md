# Tài liệu yêu cầu hệ thống

Project: **enterprise-order-management-api**

Đề tài: **Phát triển backend API cho hệ thống quản lý sản phẩm và đơn hàng trong doanh nghiệp sử dụng Golang**

## 1. Mô tả bài toán

Trong doanh nghiệp, việc quản lý sản phẩm và đơn hàng thường cần theo dõi nhiều thông tin như danh sách sản phẩm, giá bán, tồn kho, người đặt hàng, trạng thái đơn hàng và lịch sử giao dịch. Nếu quản lý thủ công bằng file Excel hoặc ghi chép rời rạc, dữ liệu dễ bị sai lệch, khó tìm kiếm, khó kiểm soát tồn kho và khó mở rộng khi số lượng sản phẩm, đơn hàng tăng lên.

Hệ thống cần cung cấp một backend API để:

- Quản lý tài khoản người dùng.
- Xác thực người dùng bằng JWT.
- Phân quyền giữa Admin và User.
- Quản lý danh sách sản phẩm.
- Tìm kiếm, lọc và phân trang sản phẩm.
- Cho phép người dùng tạo đơn hàng.
- Tự động tính tổng tiền đơn hàng.
- Kiểm tra tồn kho trước khi tạo đơn.
- Giảm tồn kho sau khi tạo đơn thành công.
- Quản lý và xem danh sách đơn hàng theo quyền của từng actor.

Mục tiêu chính của project là xây dựng một RESTful API rõ ràng, dễ bảo trì, có cấu trúc backend chuẩn để phục vụ đồ án thực tập và có thể mở rộng thành hệ thống thật trong tương lai.

## 2. Actor trong hệ thống

Hệ thống có 2 actor chính:

### Admin

Admin là người quản trị hệ thống. Admin có quyền quản lý sản phẩm, xem toàn bộ đơn hàng và kiểm soát dữ liệu nghiệp vụ chính của hệ thống.

### User

User là người dùng thông thường. User có thể đăng ký, đăng nhập, xem sản phẩm, tạo đơn hàng và xem các đơn hàng của chính mình.

## 3. Chức năng của từng actor

### Admin

Admin có các chức năng:

- Đăng nhập vào hệ thống.
- Xem danh sách sản phẩm.
- Xem chi tiết sản phẩm.
- Tạo sản phẩm mới.
- Cập nhật thông tin sản phẩm.
- Xóa sản phẩm.
- Tìm kiếm sản phẩm theo tên.
- Lọc sản phẩm theo danh mục và khoảng giá.
- Xem toàn bộ danh sách đơn hàng trong hệ thống.
- Xem chi tiết đơn hàng.
- Theo dõi tồn kho sản phẩm.

### User

User có các chức năng:

- Đăng ký tài khoản.
- Đăng nhập vào hệ thống.
- Làm mới access token bằng refresh token.
- Xem danh sách sản phẩm.
- Xem chi tiết sản phẩm.
- Tìm kiếm và lọc sản phẩm.
- Tạo đơn hàng từ danh sách sản phẩm.
- Xem danh sách đơn hàng của chính mình.
- Xem chi tiết đơn hàng của chính mình.

## 4. Danh sách module

### Auth Module

Quản lý xác thực và phiên đăng nhập:

- Đăng ký tài khoản.
- Đăng nhập.
- Sinh JWT access token.
- Sinh refresh token.
- Làm mới access token.
- Mã hóa mật khẩu bằng bcrypt.
- Middleware kiểm tra JWT.

### User Module

Quản lý thông tin người dùng:

- Lưu thông tin user.
- Phân quyền user theo role.
- Truy xuất thông tin user khi xử lý request.

Trong phạm vi 2 tháng, module này nên giữ đơn giản, chưa cần làm quản lý user đầy đủ như khóa tài khoản, đổi role qua giao diện, reset mật khẩu.

### Product Module

Quản lý sản phẩm:

- Tạo sản phẩm.
- Cập nhật sản phẩm.
- Xóa sản phẩm.
- Xem chi tiết sản phẩm.
- Xem danh sách sản phẩm.
- Phân trang danh sách sản phẩm.
- Tìm kiếm theo tên.
- Lọc theo danh mục.
- Lọc theo khoảng giá.
- Theo dõi tồn kho.

### Order Module

Quản lý đơn hàng:

- Tạo đơn hàng.
- Thêm nhiều sản phẩm vào một đơn hàng.
- Kiểm tra tồn kho trước khi tạo đơn.
- Tính tổng tiền đơn hàng.
- Lưu chi tiết từng item trong đơn hàng.
- Trừ tồn kho sau khi tạo đơn thành công.
- Xem danh sách đơn hàng.
- Phân quyền xem đơn hàng theo actor.

### Middleware Module

Xử lý các logic dùng chung trước khi request vào handler:

- Logging request.
- Recover khi panic.
- CORS.
- JWT authentication.
- Role-based authorization.

### Error Response Module

Chuẩn hóa response lỗi:

- Validation error.
- Unauthorized.
- Forbidden.
- Not found.
- Conflict.
- Internal server error.

### Database Module

Quản lý kết nối PostgreSQL:

- Mở kết nối database.
- Ping kiểm tra kết nối.
- Cấu hình connection pool.
- Migration tạo bảng.

### Documentation Module

Tài liệu phục vụ báo cáo và demo:

- README hướng dẫn chạy local.
- API docs.
- ERD.
- File `.env.example`.
- Dockerfile.
- Docker Compose.

## 5. Business rules

### Quy tắc tài khoản và xác thực

- Email của user phải là duy nhất.
- Mật khẩu phải được hash, không lưu plain text trong database.
- Access token dùng cho các API cần đăng nhập.
- Refresh token dùng để lấy access token mới khi access token hết hạn.
- API quản trị sản phẩm chỉ dành cho Admin.
- User chỉ được xem và tạo đơn hàng của chính mình.
- Admin được xem toàn bộ đơn hàng.

### Quy tắc sản phẩm

- Mỗi sản phẩm phải có tên, giá, tồn kho, ảnh sản phẩm và danh mục.
- Giá sản phẩm phải lớn hơn 0.
- Tồn kho không được nhỏ hơn 0.
- Sản phẩm không active thì không nên cho đặt hàng.
- Khi xóa sản phẩm, cần cân nhắc nếu sản phẩm đã nằm trong đơn hàng. Với đồ án 2 tháng, có thể dùng hard delete đơn giản hoặc chuyển sang `is_active = false` nếu muốn an toàn hơn.

### Quy tắc đơn hàng

- Một đơn hàng phải có ít nhất 1 sản phẩm.
- Số lượng mỗi sản phẩm trong đơn hàng phải lớn hơn 0.
- Khi tạo đơn hàng, hệ thống phải kiểm tra sản phẩm tồn tại.
- Khi tạo đơn hàng, hệ thống phải kiểm tra sản phẩm còn active.
- Khi tạo đơn hàng, hệ thống phải kiểm tra tồn kho đủ số lượng.
- Tổng tiền đơn hàng được tính từ `quantity * unit_price` của từng item.
- Đơn hàng phải lưu `unit_price` tại thời điểm mua để không bị thay đổi khi giá sản phẩm thay đổi sau này.
- Khi tạo đơn hàng thành công, tồn kho sản phẩm phải giảm tương ứng.
- Việc tạo đơn hàng và trừ tồn kho phải nằm trong cùng một transaction.
- Nếu bất kỳ bước nào khi tạo đơn hàng thất bại, toàn bộ transaction phải rollback.

### Quy tắc phân quyền

- API công khai:
  - Đăng ký.
  - Đăng nhập.
  - Refresh token.
  - Xem danh sách sản phẩm.
  - Xem chi tiết sản phẩm.
- API cần đăng nhập:
  - Tạo đơn hàng.
  - Xem đơn hàng.
- API cần quyền Admin:
  - Tạo sản phẩm.
  - Cập nhật sản phẩm.
  - Xóa sản phẩm.
  - Xem toàn bộ đơn hàng.

## 6. Non-functional requirements

### Bảo mật

- Mật khẩu phải được hash bằng bcrypt.
- Không trả password hash ra API response.
- JWT secret phải đặt trong biến môi trường, không hard-code trong source code khi deploy thật.
- Các API cần đăng nhập phải có middleware JWT.
- Các API quản trị phải có role-based authorization.
- Validate input để tránh dữ liệu sai hoặc request thiếu field.
- Không viết SQL bằng cách nối chuỗi trực tiếp với dữ liệu người dùng.
- Dùng parameterized query để giảm rủi ro SQL injection.
- File `.env` không được commit lên Git.

### Hiệu năng

- Sử dụng connection pool cho PostgreSQL.
- Danh sách sản phẩm phải có phân trang.
- Các trường hay tìm kiếm như `name`, `category_id` nên có index phù hợp.
- API không nên trả toàn bộ dữ liệu lớn trong một lần request.
- Handler nên xử lý nhanh, không chứa logic nặng.
- Transaction chỉ nên bao quanh phần cần đảm bảo toàn vẹn dữ liệu, ví dụ tạo đơn hàng và trừ tồn kho.

### Maintainability

- Code chia theo cấu trúc `handler/service/repository`.
- Handler chỉ nhận request, validate và trả response.
- Service chứa business logic.
- Repository chứa SQL.
- Error response được chuẩn hóa tại một nơi.
- DTO tách riêng với model database.
- Tên file, tên hàm, tên biến phải rõ nghĩa.
- README cần đủ hướng dẫn chạy local, chạy Docker và test API.
- Code không nên quá trừu tượng trong giai đoạn đồ án, ưu tiên dễ hiểu và dễ giải thích.

## 7. Scope nên làm trong 2 tháng

Trong 2 tháng, nên tập trung vào phạm vi vừa đủ để demo tốt:

### Tuần 1

- Phân tích yêu cầu.
- Thiết kế database.
- Tạo ERD.
- Khởi tạo project Golang Echo.
- Cấu hình PostgreSQL, `.env`, Docker Compose.

### Tuần 2

- Xây dựng Auth module.
- Đăng ký, đăng nhập.
- Hash password.
- JWT access token và refresh token.
- Middleware JWT.

### Tuần 3

- Xây dựng Product module.
- CRUD sản phẩm.
- Validation request.
- Phân quyền Admin cho API quản lý sản phẩm.

### Tuần 4

- Thêm phân trang, tìm kiếm, lọc sản phẩm.
- Chuẩn hóa error response.
- Hoàn thiện cấu trúc handler/service/repository.

### Tuần 5

- Xây dựng Order module.
- Tạo đơn hàng.
- Tính tổng tiền.
- Lưu order items.
- Transaction khi tạo đơn hàng.

### Tuần 6

- Hoàn thiện role-based authorization.
- User xem đơn của mình.
- Admin xem toàn bộ đơn hàng.
- Kiểm tra tồn kho và trừ tồn kho.

### Tuần 7

- Viết README.
- Viết API docs.
- Hoàn thiện ERD.
- Test API bằng curl hoặc Postman.
- Sửa lỗi và refactor nhẹ.

### Tuần 8

- Dockerize project.
- Chuẩn bị deploy lên Render, Railway hoặc VPS.
- Chuẩn bị dữ liệu demo.
- Chuẩn bị slide/báo cáo.
- Luyện demo luồng chính.

## 8. Scope không nên làm để tránh quá tải

Không nên làm các chức năng sau trong phiên bản đồ án 2 tháng:

- Frontend quá đầy đủ như dashboard phức tạp, biểu đồ nâng cao, quản trị nhiều màn hình.
- Thanh toán online.
- Tích hợp cổng vận chuyển.
- Quản lý kho nhiều chi nhánh.
- Quản lý nhà cung cấp.
- Quản lý nhập kho, xuất kho chi tiết.
- Quản lý hóa đơn điện tử.
- Báo cáo doanh thu phức tạp.
- Gửi email hoặc SMS tự động.
- Quên mật khẩu qua email.
- OAuth login bằng Google/Facebook.
- Upload ảnh sản phẩm lên cloud.
- Event-driven architecture, message queue, microservices.
- Redis cache.
- Elasticsearch.
- Audit log đầy đủ.
- Phân quyền động theo permission matrix.
- Unit test/integration test quá rộng nếu chưa đủ thời gian.

Các chức năng trên có thể ghi vào phần "hướng phát triển" trong báo cáo thay vì đưa vào scope chính.

## 9. Kết quả cần có khi demo

Khi demo, project nên có đầy đủ các kết quả sau:

### Source code

- Project Golang chạy được.
- Cấu trúc thư mục rõ ràng.
- Có handler/service/repository.
- Không viết SQL trong handler.
- Không viết business logic trong handler.

### Database

- PostgreSQL chạy được.
- Có migration tạo bảng.
- Có dữ liệu admin mẫu.
- Có bảng users, products, orders, order_items, refresh_tokens.

### API

- API health check.
- API register.
- API login.
- API refresh token.
- API CRUD sản phẩm.
- API danh sách sản phẩm có phân trang, search, filter.
- API tạo đơn hàng.
- API xem danh sách đơn hàng.
- API phân quyền theo Admin/User.

### Bảo mật

- Password được hash.
- API cần đăng nhập có JWT middleware.
- API quản trị có kiểm tra role.
- Error response không lộ thông tin nhạy cảm.

### Transaction

- Khi tạo đơn hàng, hệ thống kiểm tra tồn kho.
- Nếu đủ tồn kho, đơn hàng được tạo và stock giảm.
- Nếu không đủ tồn kho, đơn hàng không được tạo.
- Transaction rollback khi có lỗi.

### Documentation

- README hướng dẫn chạy local.
- `.env.example`.
- Dockerfile.
- `docker-compose.yml`.
- ERD.
- API docs.
- Curl hoặc Postman collection để test.

### Luồng demo đề xuất

1. Chạy project bằng Docker Compose hoặc local.
2. Gọi `/health` để kiểm tra server.
3. Login admin.
4. Admin tạo sản phẩm.
5. User đăng ký tài khoản.
6. User login.
7. User xem danh sách sản phẩm.
8. User tạo đơn hàng.
9. Kiểm tra tồn kho sản phẩm đã giảm.
10. User xem đơn hàng của mình.
11. Admin xem toàn bộ đơn hàng.
12. Demo lỗi phân quyền khi User gọi API tạo sản phẩm.
