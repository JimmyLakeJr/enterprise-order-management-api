# Demo Runbook

Checklist ngắn để chạy demo với backend root và PostgreSQL thật.

## Trước khi demo

- [ ] `docker compose up -d --build` và `docker compose ps` đều ổn định.
- [ ] `GET /health` trả `200`.
- [ ] Frontend chạy tại `http://localhost:5173`.
- [ ] Database có admin và dữ liệu mẫu cần thiết.
- [ ] Xác nhận credential demo đang dùng đúng với volume hiện tại; volume hiện tại đã kiểm tra được `vu@gmail.com / 123456` là admin.
- [ ] Không dùng credential/JWT secret local trên môi trường public.

## Luồng demo

1. Mở trang public, lọc danh mục và xem chi tiết sản phẩm.
2. Đăng ký hoặc đăng nhập user, kiểm tra profile.
3. Thêm sản phẩm vào giỏ, cập nhật số lượng và tạo đơn.
4. Mở “Đơn hàng của tôi” và xem chi tiết đơn.
5. Đăng nhập admin, xem dashboard.
6. Tạo/sửa/ẩn/khôi phục danh mục và sản phẩm; xác nhận public chỉ hiện dữ liệu active.
7. Mở danh sách đơn, kiểm tra thời gian, user summary và cập nhật trạng thái hợp lệ.
8. Mở danh sách user và kiểm tra dữ liệu hiển thị.

## Browser smoke test cuối

- [ ] Public, user và admin đều chạy với backend + DB thật.
- [ ] Console không có lỗi nghiêm trọng; Network không gọi endpoint không tồn tại.
- [ ] Loading, empty, error/retry, toast và confirm dialog hiển thị đúng.
- [ ] Kiểm tra desktop, tablet và mobile; bảng admin cuộn ngang khi cần.
- [ ] Refresh trang ở các route chính không làm mất phiên hoặc vỡ giao diện.

## Warning cần nói rõ

- Hủy đơn chỉ đổi trạng thái, chưa hoàn tồn kho.
- Sản phẩm dùng `image_url`, chưa upload file.
- Profile chỉ lưu tên hiển thị; avatar/video chưa có lưu trữ backend.
- Nếu cần quay về dữ liệu seed sạch của migration, phải reset volume PostgreSQL bằng `docker compose down -v`.
- Các module Phase 2 không nằm trong luồng demo; danh sách xem tại tài liệu scope.
