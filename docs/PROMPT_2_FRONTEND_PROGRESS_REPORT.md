# Archive — Frontend Progress Report

> Snapshot lịch sử của sprint frontend. Không dùng file này làm source of truth; xem [CURRENT_PROJECT_SCOPE_ANALYSIS.md](CURRENT_PROJECT_SCOPE_ANALYSIS.md).

## Kết quả đã hoàn thành

- Auth, public product/category, cart `localStorage`, create order và user order flow đã nối API thật.
- Admin dashboard, category, product, order và user core đã có.
- UI liquid glass/responsive, loading/error/empty state, toast và custom confirm đã được bổ sung.
- Copy dev note thô đã được dọn; frontend không gọi các Phase 2 API chưa tồn tại.

## Xác nhận tại thời điểm report

- Frontend lint/build pass.
- Browser smoke test public/user/admin với backend + PostgreSQL pass.
- Responsive và console/network đã được kiểm tra trong phạm vi sprint.

## Việc còn theo dõi

- Smoke test lại trên dữ liệu cuối trước demo.
- Dọn dead component/API alias/CSS, metadata và copy còn sót.
- Bổ sung automated test/E2E khi scope cho phép.

Backlog tích hợp hiện tại được duy trì duy nhất tại [frontend-backend-gap.md](frontend-backend-gap.md).
