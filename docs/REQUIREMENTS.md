# Requirements

## Mục tiêu

Xây dựng backend API Golang/Echo quản lý sản phẩm, tồn kho cơ bản, user và đơn hàng; frontend React chứng minh các flow thật với PostgreSQL.

## Functional requirements

### Public và auth

- Public xem category/product active và chi tiết sản phẩm.
- User đăng ký, đăng nhập, refresh, logout và đọc profile.
- User cập nhật tên nhưng không tự sửa role/is_active.

### User flow

- Cart được lưu local và gửi `product_id`, `quantity` khi tạo order.
- Backend kiểm tra product/stock, tính giá và tạo order theo transaction.
- User chỉ xem danh sách/chi tiết order của mình.

### Admin flow

- Admin quản lý category và product, gồm active/inactive list và restore.
- Public list không lộ dữ liệu inactive.
- Admin xem order với timestamps/user summary và cập nhật trạng thái hợp lệ.
- Admin xem/quản lý user trong contract hiện có.

### Frontend

- Có loading, empty, error/retry, toast, custom confirm và focus state.
- Responsive cho desktop/tablet/mobile; table có horizontal scroll khi cần.
- Copy user-facing dùng tiếng Việt nhất quán và không hiển thị dev/gap note thô.
- Không gọi endpoint chưa tồn tại hoặc giả lập tính năng backend chưa có.

## Business rules

- Email unique; password hash; protected route dùng JWT/role guard.
- Product thuộc category; restore product cần category active.
- Giá order do backend tính; order creation và stock decrease cùng transaction.
- Ownership được kiểm tra khi user đọc order.
- Cancel hiện chỉ đổi trạng thái, chưa hoàn stock.

## Non-functional requirements

- Source of truth là backend root, không phải skeleton `backend/`.
- Validation/error response nhất quán, code theo handler/service/repository.
- Secret và credential local không dùng cho public production.
- Kiểm tra backend bằng package hẹp; frontend bằng lint/build và browser smoke test.
- Không rewrite hoặc thêm module lớn ngoài scope demo.

## Acceptance demo

- Auth, public product/category, cart/create order và admin core chạy với DB thật.
- Console/network không có lỗi nghiêm trọng hoặc request tới API không tồn tại.
- Giao diện dùng được trên mobile/tablet/desktop.
- Warning về cancel, media, pagination và secret được trình bày rõ.

## Backlog

- Demo verification cuối và frontend cleanup.
- Order pagination/batch item query; quyết định cancel/restock.
- Automated test/E2E.
- Các module Phase 2 được giữ ngoài sprint; danh sách duy nhất nằm trong tài liệu scope.

Trạng thái triển khai không lặp tại đây; xem [CURRENT_PROJECT_SCOPE_ANALYSIS.md](CURRENT_PROJECT_SCOPE_ANALYSIS.md).
