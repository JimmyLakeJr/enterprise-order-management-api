# Project Standards

## Phạm vi

- Giữ project là API quản lý sản phẩm, tồn kho cơ bản, user và order với frontend demo.
- Không tự mở rộng module Phase 2 hoặc inventory ledger; danh sách ngoài scope được duy trì tại tài liệu scope.
- Ưu tiên patch nhỏ, không rewrite module đang hoạt động.

## Source of truth

- Backend: `cmd/api`, `internal/*`, `migrations/001_init.sql`.
- Runtime: root `Dockerfile`, `docker-compose.yml`.
- Frontend demo: `frontend/src`.
- Không dùng `backend/` skeleton làm runtime nghiệp vụ.

## Backend

- Giữ luồng `router -> middleware -> handler -> service -> repository`.
- Handler xử lý HTTP; service giữ business rule; repository giữ SQL.
- Validate input, dùng error/response helper chung và không lộ lỗi nội bộ.
- Bảo vệ route bằng JWT/role; kiểm tra ownership ở resource của user.
- Dùng transaction cho create order/stock update; tránh N+1 khi thay đổi query mới.
- Migration rủi ro phải được đánh giá riêng, không sửa chỉ để làm đẹp demo.

## Frontend

- API call tập trung ở API layer; không hard-code endpoint trong page/component.
- Không gửi giá order từ client và không gọi contract chưa tồn tại.
- Copy tiếng Việt nhất quán; có loading/error/empty state cho request chính.
- Dùng token/CSS hiện có, responsive và reduced-motion; tránh thêm UI library nặng không cần thiết.
- Ưu tiên custom confirm/toast thay browser dialog trong flow demo.

## Kiểm tra

```bash
go test ./cmd/... ./internal/...
go vet ./cmd/... ./internal/...

cd frontend
npm run lint
npm run build
```

Không chạy `go test ./...` khi nó có thể quét nhầm `frontend/node_modules`. Chỉ chạy test/build liên quan đến phase vừa sửa.

## Tài liệu và Git

- README định hướng; local guide hướng dẫn chạy; runbook chứa checklist demo; scope/gap giữ trạng thái hiện tại.
- Prompt/report cũ chỉ là archive, không dùng làm source of truth.
- Không commit secret, token, dữ liệu cá nhân hoặc credential production.
- Giữ thay đổi tập trung, không ghi đè phần việc không liên quan trong worktree.
