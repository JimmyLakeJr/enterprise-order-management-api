# Archive — Backend Priority Audit

> Bản tóm tắt audit lịch sử. Route và code trong backend root mới là source of truth.

## Hạng mục đã xử lý

- Profile update tối thiểu cho `name`; `/auth/me` trả dữ liệu nhất quán.
- Admin xem active/inactive và restore category/product; public vẫn active-only.
- Guard category/product tránh trạng thái restore không hợp lệ.
- Order response có timestamps và user summary bằng query join, không tạo N+1 user.
- Validation/error mapping và test service liên quan đã được củng cố.

## Xác nhận tại thời điểm report

- `go test ./cmd/... ./internal/...` pass.
- `go vet ./cmd/... ./internal/...` pass.
- Docker runtime và frontend integration smoke test pass.

## Gap còn đúng

- Order list chưa pagination; order items còn có thể batch query.
- Cancel chỉ đổi trạng thái, chưa hoàn stock.
- Media chỉ dùng URL; không có upload/storage.

Chi tiết hiện tại xem [CURRENT_PROJECT_SCOPE_ANALYSIS.md](CURRENT_PROJECT_SCOPE_ANALYSIS.md) và [api.md](api.md).
