# Phân tích nghiệp vụ

## Bài toán

Project minh họa một backend quản lý sản phẩm, tồn kho cơ bản, người dùng và đơn hàng. React frontend dùng để chứng minh API hoạt động qua flow public, user và admin; không mở rộng thành e-commerce hoàn chỉnh.

## Vai trò

- Public: xem category/product active.
- User: đăng ký, đăng nhập, sửa tên profile, quản lý giỏ local và tạo/xem đơn của mình.
- Admin: quản lý category, product, order và user trong scope demo.

## Luồng cốt lõi

1. Admin chuẩn bị category/product và stock.
2. User chọn sản phẩm; cart chỉ giữ `product_id`, quantity và dữ liệu hiển thị tạm.
3. Backend kiểm tra sản phẩm/stock, tính giá và tạo order trong transaction.
4. User xem đơn theo ownership; admin xem user summary và cập nhật trạng thái hợp lệ.
5. Public list tiếp tục chỉ hiển thị dữ liệu active.

## Quy tắc nghiệp vụ

- Email user là duy nhất; password được hash.
- User không tự sửa role hoặc `is_active`.
- Product thuộc category; product restore yêu cầu category active.
- Category có product active không được vô hiệu hóa theo cách làm kẹt dữ liệu.
- Giá và tổng tiền do backend tính, không tin giá từ frontend.
- Tạo order và giảm stock phải cùng transaction.
- User chỉ đọc order của mình; admin có quyền quản trị.
- Cancel hiện chỉ đổi trạng thái, chưa hoàn stock.

## Trạng thái

Auth, profile name, category/product public/admin, cart, order và admin UI core đã hoạt động. Database có bảy bảng nghiệp vụ cơ bản; giao diện responsive/liquid glass đã có. Report gần nhất xác nhận test/vet backend, lint/build frontend và browser smoke test đã pass.

## Rủi ro/backlog

- Xác minh lại console/network và mobile/tablet trước demo cuối.
- Dọn frontend dead code, copy và metadata còn sót.
- Thêm order pagination và batch item query.
- Chỉ thiết kế restock sau khi chốt rule, idempotency và concurrency.
- Giữ các module Phase 2 ngoài sprint hiện tại; danh sách duy nhất nằm trong tài liệu scope.

Chi tiết trạng thái duy nhất xem tại [CURRENT_PROJECT_SCOPE_ANALYSIS.md](CURRENT_PROJECT_SCOPE_ANALYSIS.md); gap còn lại xem tại [frontend-backend-gap.md](frontend-backend-gap.md).
