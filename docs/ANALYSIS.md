# Tài liệu phân tích yêu cầu hệ thống

Project: **enterprise-order-management-api**

Đề tài: **Phát triển backend API cho hệ thống quản lý sản phẩm và đơn hàng trong doanh nghiệp sử dụng Golang**

## 1. Mô tả bài toán

Trong hoạt động kinh doanh của doanh nghiệp, việc quản lý sản phẩm và đơn hàng là một nghiệp vụ quan trọng. Doanh nghiệp cần lưu trữ thông tin sản phẩm, danh mục sản phẩm, giá bán, tồn kho, thông tin người dùng và các đơn hàng phát sinh trong quá trình mua bán. Nếu quản lý thủ công bằng file Excel hoặc các công cụ rời rạc, dữ liệu dễ bị sai lệch, khó kiểm soát tồn kho, khó theo dõi trạng thái đơn hàng và khó mở rộng khi quy mô dữ liệu tăng lên.

Đề tài này tập trung xây dựng một backend API phục vụ hệ thống quản lý sản phẩm và đơn hàng. Hệ thống cho phép người dùng xem sản phẩm, tạo đơn hàng, theo dõi đơn hàng của mình; đồng thời cho phép quản trị viên quản lý danh mục, sản phẩm, người dùng và trạng thái đơn hàng.

Backend API là phần trọng tâm của đồ án. Frontend nếu có chỉ đóng vai trò demo client đơn giản để kiểm thử API và trình bày nghiệp vụ chính.

## 2. Mục tiêu của hệ thống

Mục tiêu của hệ thống là xây dựng một backend API có cấu trúc rõ ràng, dễ bảo trì và phù hợp với phạm vi thực tập 2 tháng.

Các mục tiêu cụ thể:

- Xây dựng RESTful API bằng Golang và Echo Framework.
- Sử dụng PostgreSQL làm hệ quản trị cơ sở dữ liệu.
- Tổ chức code theo kiến trúc `handler/service/repository`.
- Xây dựng chức năng xác thực bằng JWT access token và refresh token.
- Phân quyền người dùng theo vai trò Guest, User và Admin.
- Quản lý danh mục sản phẩm.
- Quản lý sản phẩm, giá bán và tồn kho.
- Cho phép người dùng tạo đơn hàng.
- Đảm bảo tạo đơn hàng và trừ tồn kho trong cùng một transaction.
- Chuẩn hóa response JSON và xử lý lỗi tập trung.
- Cung cấp tài liệu hướng dẫn chạy local, API docs, ERD và Docker Compose.

## 3. Actor trong hệ thống

### Guest

Guest là người chưa đăng nhập vào hệ thống. Guest chỉ có quyền xem các thông tin công khai như danh mục và sản phẩm đang hoạt động.

### User

User là người dùng đã đăng ký và đăng nhập vào hệ thống. User có thể xem sản phẩm, tạo đơn hàng và quản lý các đơn hàng của chính mình.

### Admin

Admin là người quản trị hệ thống. Admin có quyền quản lý dữ liệu chính của hệ thống như người dùng, danh mục, sản phẩm và đơn hàng.

## 4. Chức năng của từng actor

### Guest

Guest có các chức năng:

- Xem danh sách danh mục đang hoạt động.
- Xem danh sách sản phẩm đang hoạt động.
- Xem chi tiết sản phẩm.
- Tìm kiếm sản phẩm theo tên.
- Lọc sản phẩm theo danh mục và khoảng giá.
- Đăng ký tài khoản.
- Đăng nhập vào hệ thống.

### User

User có toàn bộ chức năng của Guest và có thêm các chức năng:

- Xem thông tin cá nhân.
- Đăng xuất khỏi hệ thống.
- Làm mới access token bằng refresh token.
- Tạo đơn hàng.
- Xem danh sách đơn hàng của chính mình.
- Theo dõi thông tin chi tiết các sản phẩm trong đơn hàng.

### Admin

Admin có toàn bộ chức năng của User và có thêm các chức năng:

- Xem danh sách người dùng.
- Khóa hoặc vô hiệu hóa người dùng bằng soft delete.
- Tạo, cập nhật, vô hiệu hóa danh mục.
- Tạo, cập nhật, vô hiệu hóa sản phẩm.
- Xem toàn bộ đơn hàng trong hệ thống.
- Cập nhật trạng thái đơn hàng theo luồng nghiệp vụ hợp lệ.

## 5. Danh sách module cần xây dựng

### Auth Module

Module xác thực và quản lý phiên đăng nhập:

- Đăng ký tài khoản.
- Đăng nhập.
- Sinh JWT access token.
- Sinh refresh token.
- Hash refresh token trước khi lưu database.
- Làm mới access token.
- Đăng xuất và revoke refresh token.
- Hash mật khẩu bằng bcrypt.

### User Module

Module quản lý người dùng:

- Lưu thông tin người dùng.
- Gán vai trò cho người dùng thông qua `role_id`.
- User xem thông tin cá nhân.
- Admin xem danh sách người dùng.
- Admin vô hiệu hóa người dùng bằng `is_active = false`.

### Role Module

Module quản lý vai trò:

- Lưu danh sách vai trò như `admin`, `user`.
- Hỗ trợ middleware phân quyền.
- Không cần xây dựng quản lý role động trong phạm vi 2 tháng.

### Category Module

Module quản lý danh mục sản phẩm:

- Guest/User/Admin xem danh sách danh mục đang hoạt động.
- Admin tạo danh mục.
- Admin cập nhật danh mục.
- Admin vô hiệu hóa danh mục bằng `is_active = false`.

### Product Module

Module quản lý sản phẩm:

- Xem danh sách sản phẩm.
- Xem chi tiết sản phẩm.
- Tìm kiếm sản phẩm.
- Lọc sản phẩm theo danh mục và khoảng giá.
- Phân trang danh sách sản phẩm.
- Admin tạo sản phẩm.
- Admin cập nhật sản phẩm.
- Admin vô hiệu hóa sản phẩm bằng `is_active = false`.
- Quản lý giá bán và tồn kho.

### Order Module

Module quản lý đơn hàng:

- User tạo đơn hàng.
- Backend kiểm tra sản phẩm tồn tại, còn active và đủ tồn kho.
- Backend tự lấy giá sản phẩm từ database tại thời điểm đặt hàng.
- Backend tự tính `subtotal` và `total_amount`.
- Tạo order, order_items và trừ tồn kho trong cùng một transaction.
- User xem đơn hàng của chính mình.
- Admin xem toàn bộ đơn hàng.
- Admin cập nhật trạng thái đơn hàng.

### Middleware Module

Module xử lý các tác vụ trung gian:

- CORS.
- Logger.
- Recovery.
- JWT authentication.
- Role-based authorization.

### Response & Error Handling Module

Module chuẩn hóa phản hồi API:

- Response thành công.
- Response lỗi.
- Response phân trang.
- Validation error.
- Authentication error.
- Authorization error.
- Not found, conflict và internal server error.

### Database Module

Module quản lý kết nối database:

- Khởi tạo PostgreSQL connection pool bằng `pgxpool`.
- Đóng kết nối khi server shutdown.
- Hỗ trợ transaction cho nghiệp vụ tạo đơn hàng.
- Cung cấp migration SQL để tạo bảng.

### Documentation Module

Module tài liệu:

- README hướng dẫn chạy local và deploy.
- ERD.
- API docs hoặc Postman collection.
- `.env.example`.
- Dockerfile và Docker Compose.

## 6. Business rules

### Business rules cho sản phẩm

- Mỗi sản phẩm phải thuộc một danh mục.
- Sản phẩm phải có tên, giá bán, tồn kho, ảnh sản phẩm và trạng thái hoạt động.
- Giá sản phẩm phải lớn hơn 0.
- Tồn kho không được âm.
- Sản phẩm không hoạt động (`is_active = false`) không được hiển thị ở API public.
- Sản phẩm không hoạt động không được phép đặt hàng.
- Khi xóa sản phẩm, hệ thống ưu tiên soft delete bằng cách cập nhật `is_active = false`.
- Backend không tin giá sản phẩm từ frontend khi tạo đơn hàng.

### Business rules cho danh mục

- Mỗi danh mục phải có tên.
- Tên danh mục nên là duy nhất.
- Danh mục có trạng thái `is_active`.
- Danh mục không hoạt động không được hiển thị ở API public.
- Sản phẩm thuộc danh mục không hoạt động không nên hiển thị ở API public.
- Khi xóa danh mục, hệ thống ưu tiên soft delete bằng cách cập nhật `is_active = false`.

### Business rules cho người dùng

- Email người dùng phải duy nhất.
- Mật khẩu phải được hash bằng bcrypt trước khi lưu vào database.
- Không lưu plain password.
- Không trả `password_hash` ra response.
- User phải có role thông qua `role_id`.
- User bị vô hiệu hóa không được đăng nhập.
- Refresh token phải được hash trước khi lưu vào database.
- Logout phải revoke refresh token.
- User chỉ được xem đơn hàng của chính mình.
- Admin được xem toàn bộ đơn hàng.

### Business rules cho đơn hàng

- Chỉ user đã đăng nhập mới được tạo đơn hàng.
- Một đơn hàng phải có ít nhất một sản phẩm.
- Số lượng mỗi sản phẩm trong đơn hàng phải lớn hơn 0.
- Khi tạo đơn hàng, hệ thống phải kiểm tra sản phẩm tồn tại.
- Khi tạo đơn hàng, hệ thống phải kiểm tra sản phẩm còn active.
- Khi tạo đơn hàng, hệ thống phải kiểm tra tồn kho đủ.
- `unit_price` phải được lấy từ database tại thời điểm đặt hàng.
- `subtotal` được tính bằng `quantity * unit_price`.
- `total_amount` được tính từ tổng `subtotal` của các item.
- Tạo order, order_items và trừ stock phải nằm trong cùng một transaction.
- Nếu bất kỳ bước nào lỗi, toàn bộ transaction phải rollback.
- Không cho user xem đơn hàng của người khác.
- Chỉ admin được cập nhật trạng thái đơn hàng.
- Trạng thái đơn hàng gồm: `pending`, `confirmed`, `shipping`, `completed`, `cancelled`.
- Luồng trạng thái hợp lệ:
  - `pending -> confirmed`
  - `pending -> cancelled`
  - `confirmed -> shipping`
  - `confirmed -> cancelled`
  - `shipping -> completed`
- Không cho chuyển trạng thái ngược hoặc sai nghiệp vụ.

## 7. Non-functional requirements

### Bảo mật

- Mật khẩu phải được hash bằng bcrypt.
- JWT secret phải lấy từ biến môi trường.
- Access token có thời hạn ngắn, ví dụ 15 phút.
- Refresh token có thời hạn dài hơn, ví dụ 7 ngày.
- Refresh token phải được hash trước khi lưu vào database.
- Không lưu raw refresh token.
- Logout phải revoke refresh token.
- API cần đăng nhập phải có JWT middleware.
- API admin phải có middleware kiểm tra role.
- Không commit file `.env`.
- Không expose lỗi database thô cho client.
- Không nối chuỗi SQL trực tiếp với input người dùng.
- Phải dùng parameterized query để giảm rủi ro SQL injection.
- CORS phải cấu hình theo frontend domain.

### Hiệu năng

- Sử dụng PostgreSQL connection pool bằng `pgxpool`.
- API danh sách sản phẩm phải có phân trang.
- Các trường thường dùng để lọc hoặc tìm kiếm nên có index phù hợp.
- Không trả toàn bộ dữ liệu lớn trong một response.
- Transaction chỉ bao quanh các thao tác cần đảm bảo toàn vẹn dữ liệu.
- Handler cần xử lý nhanh, không chứa logic nặng.

### Maintainability

- Code được chia theo kiến trúc `handler/service/repository`.
- Handler chỉ nhận request, validate, gọi service và trả response.
- Service chứa business logic.
- Repository chứa SQL thuần.
- DTO tách riêng với model database.
- Response JSON được chuẩn hóa.
- Error handling tập trung.
- Tên file, tên hàm và tên biến cần rõ nghĩa.
- Code không nên trừu tượng quá mức để phù hợp trình độ intern/junior.
- README, ERD và API docs phải được cập nhật theo code.

### Scalability cơ bản

- Tách rõ các module để dễ mở rộng thêm chức năng.
- Sử dụng connection pool để phục vụ nhiều request đồng thời.
- Thiết kế database có khóa ngoại và index cơ bản.
- Response phân trang giúp hệ thống xử lý dữ liệu lớn tốt hơn.
- Có thể deploy backend, frontend và database thành các service riêng.
- Có thể mở rộng sau này sang Redis cache, message queue hoặc background job nếu cần.

### Logging

- Hệ thống cần log các request HTTP cơ bản.
- Cần có middleware logger để ghi nhận method, path, status code và thời gian xử lý.
- Không log password, access token, refresh token hoặc thông tin nhạy cảm.
- Log lỗi server để hỗ trợ debug trong quá trình phát triển và deploy.

### Error handling

- API phải trả lỗi theo format thống nhất.
- Validation error phải trả thông tin field không hợp lệ.
- Unauthorized error dùng cho trường hợp thiếu hoặc sai token.
- Forbidden error dùng cho trường hợp không đủ quyền.
- Not found error dùng khi dữ liệu không tồn tại.
- Conflict error dùng cho dữ liệu trùng, ví dụ email hoặc tên danh mục.
- Internal server error không được trả chi tiết lỗi database thô cho client.

## 8. Scope nên làm trong 2 tháng

### Tuần 1

- Phân tích yêu cầu.
- Thiết kế database.
- Vẽ ERD.
- Khởi tạo project Golang Echo.
- Cấu hình PostgreSQL, Docker Compose và `.env.example`.

### Tuần 2

- Xây dựng Auth module.
- Đăng ký, đăng nhập.
- Hash password.
- JWT access token và refresh token.
- Hash refresh token trước khi lưu database.
- Middleware JWT.

### Tuần 3

- Xây dựng User và Role module ở mức cơ bản.
- User xem thông tin cá nhân.
- Admin xem danh sách user.
- Middleware role authorization.

### Tuần 4

- Xây dựng Category module.
- Xây dựng Product module.
- CRUD category/product cho admin.
- Public API xem danh mục và sản phẩm.

### Tuần 5

- Thêm phân trang, tìm kiếm, lọc sản phẩm.
- Chuẩn hóa response format.
- Chuẩn hóa error handling.
- Hoàn thiện validation request.

### Tuần 6

- Xây dựng Order module.
- Tạo đơn hàng.
- Kiểm tra tồn kho.
- Tính tổng tiền.
- Transaction khi tạo đơn hàng.

### Tuần 7

- Admin xem toàn bộ đơn hàng.
- User xem đơn hàng của chính mình.
- Admin cập nhật trạng thái đơn hàng.
- Kiểm tra rule chuyển trạng thái.
- Test API bằng curl hoặc Postman.

### Tuần 8

- Hoàn thiện README, API docs, ERD.
- Dockerize project.
- Chuẩn bị deploy demo.
- Chuẩn bị dữ liệu demo.
- Chuẩn bị slide và báo cáo.
- Luyện demo các luồng nghiệp vụ chính.

## 9. Scope không nên làm để tránh quá tải

Trong phạm vi thực tập 2 tháng, không nên đưa các chức năng sau vào scope chính:

- Frontend dashboard phức tạp.
- Thanh toán online.
- Tích hợp cổng vận chuyển.
- Quản lý nhiều kho hoặc nhiều chi nhánh.
- Quản lý nhà cung cấp.
- Nhập kho, xuất kho chi tiết.
- Báo cáo doanh thu nâng cao.
- Biểu đồ thống kê phức tạp.
- Gửi email hoặc SMS tự động.
- Quên mật khẩu qua email.
- OAuth login bằng Google/Facebook.
- Upload ảnh sản phẩm lên cloud.
- Redis cache.
- Elasticsearch.
- Message queue.
- Microservices.
- Audit log đầy đủ.
- Phân quyền động theo permission matrix.
- Unit test và integration test bao phủ toàn bộ hệ thống nếu không đủ thời gian.

Các chức năng trên có thể đưa vào phần định hướng phát triển trong tương lai.

## 10. Kết quả cần có khi demo cuối kỳ

### Source code

- Project Golang chạy được.
- Code có cấu trúc rõ ràng.
- Có phân tách handler, service, repository.
- Không viết SQL trong handler.
- Không viết business logic trong handler.

### Database

- PostgreSQL chạy được.
- Có migration SQL.
- Có các bảng chính: roles, users, refresh_tokens, categories, products, orders, order_items.
- Có dữ liệu mẫu như admin account và category mẫu.

### API

- API health check.
- API register, login, refresh token, logout.
- API xem thông tin cá nhân.
- API quản lý user cho admin.
- API quản lý category cho admin.
- API quản lý product cho admin.
- API xem category/product public.
- API tạo order cho user.
- API xem order theo quyền.
- API cập nhật trạng thái order cho admin.

### Bảo mật

- Password được hash.
- Refresh token được hash trước khi lưu database.
- Logout revoke refresh token.
- API cần đăng nhập có JWT middleware.
- API admin có role middleware.
- User không xem được order của người khác.

### Transaction và nghiệp vụ

- Tạo order kiểm tra product tồn tại.
- Tạo order kiểm tra product active.
- Tạo order kiểm tra đủ stock.
- Backend tự tính giá và tổng tiền.
- Tạo order và trừ stock trong cùng transaction.
- Rollback khi có lỗi.
- Admin cập nhật trạng thái order theo luồng hợp lệ.

### Documentation

- README hướng dẫn chạy local và deploy.
- `.env.example`.
- Dockerfile.
- `docker-compose.yml`.
- ERD.
- API docs hoặc Postman collection.
- Tài liệu phân tích yêu cầu.

### Luồng demo đề xuất

1. Chạy project bằng Docker Compose hoặc local.
2. Gọi API `/health`.
3. Login bằng tài khoản admin.
4. Admin tạo category.
5. Admin tạo product.
6. Guest xem danh sách product.
7. User đăng ký tài khoản.
8. User đăng nhập.
9. User tạo order.
10. Kiểm tra stock sản phẩm đã giảm.
11. User xem danh sách order của mình.
12. Admin xem toàn bộ order.
13. Admin cập nhật trạng thái order.
14. Demo lỗi phân quyền khi User gọi API admin.
