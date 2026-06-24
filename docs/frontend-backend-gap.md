# Frontend–Backend Gap

Chỉ ghi các gap còn đúng ở thời điểm hiện tại. Trạng thái module đầy đủ xem tại [CURRENT_PROJECT_SCOPE_ANALYSIS.md](CURRENT_PROJECT_SCOPE_ANALYSIS.md).

## Gap còn lại

| Khu vực | Hiện trạng | Hướng xử lý |
|---|---|---|
| Product media | Backend lưu `image_url`, không nhận file | Giữ URL trong demo; upload thuộc Phase 2 |
| Profile media | Backend cập nhật tên; avatar/video chưa lưu | Không giả lập đã lưu; preview local phải ghi rõ phạm vi |
| Order list | Có timestamp và user summary, chưa có pagination | Bổ sung pagination khi dữ liệu demo cần |
| Order items | Có thể còn query theo từng order | Batch query sau khi ưu tiên hiệu năng được duyệt |
| Cancel order | Chỉ đổi trạng thái, không hoàn stock | Chốt rule và idempotency/concurrency trước khi restock |
| Frontend tests | Lint/build và smoke test đã có, chưa có E2E tự động | Bổ sung test cho auth/cart/order/admin ở sprint sau |

## Quy tắc tích hợp cần giữ

- Public category/product chỉ hiển thị dữ liệu active.
- Admin dùng API inactive/restore hiện có, không suy diễn từ public list.
- Cart gửi item theo `product_id` và `quantity`; giá được backend tính lại.
- UI không gọi endpoint chưa tồn tại và không bịa dữ liệu backend chưa trả.
- Copy user-facing dùng tiếng Việt nhất quán; giới hạn demo được trình bày như trạng thái sản phẩm, không dùng dev note thô.

## Phase 2

Không mở rộng module Phase 2 trong sprint demo hiện tại; danh sách duy nhất xem tại [CURRENT_PROJECT_SCOPE_ANALYSIS.md](CURRENT_PROJECT_SCOPE_ANALYSIS.md).
