# Enterprise Order Management API

Backend API và frontend demo cho hệ thống quản lý sản phẩm và đơn hàng trong doanh nghiệp, xây dựng bằng Golang, Echo v4, PostgreSQL và React + Vite.

Project phù hợp đồ án thực tập 2 tháng với trọng tâm là backend API, đồng thời có frontend đủ dùng để demo các flow chính: xem sản phẩm, giỏ hàng, tạo đơn hàng, xem đơn hàng và quản trị admin.

## Công Nghệ

Backend:

- Golang 1.22+
- Echo v4
- PostgreSQL
- SQL thuần với `pgx/v5` và `pgxpool`
- JWT access token + refresh token
- `github.com/golang-jwt/jwt/v5`
- `golang.org/x/crypto/bcrypt`
- `github.com/go-playground/validator/v10`
- `.env`, `godotenv`, `os.Getenv`
- Docker Compose

Frontend:

- React
- Vite
- React Router DOM
- Axios
- CSS thuần với CSS variables
- LocalStorage cho demo auth token và cart

Không dùng GORM, Fiber, MySQL, MongoDB hoặc ORM khác.

## Kiến Trúc Backend

Luồng xử lý chính:

```text
HTTP Request -> Handler -> Service -> Repository -> PostgreSQL
```

Các thư mục chính:

```text
cmd/api/main.go
internal/config
internal/database
internal/dto
internal/handler
internal/http
internal/middleware
internal/model
internal/pkg
internal/repository
internal/service
migrations/001_init.sql
docs
frontend
```

Quy tắc:

- Handler chỉ parse request, validate, gọi service và trả response.
- Service xử lý business logic, authorization, transaction và status transition.
- Repository thao tác database bằng SQL thuần, dùng parameterized query.
- Không viết SQL trong handler.
- Không trả `password_hash` ra response.
- Tạo order phải dùng transaction.
- Backend tự tính `total_amount`, không tin giá từ frontend.

## Response Format

Success:

```json
{
  "success": true,
  "message": "Success",
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {}
}
```

Pagination:

```json
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
```

## Chức Năng Backend Chính

Auth:

- Register
- Login
- Refresh access token
- Logout revoke refresh token
- Password hash bằng bcrypt
- Refresh token được hash trước khi lưu database

Product và Category:

- Guest xem danh mục, danh sách sản phẩm và chi tiết sản phẩm
- Admin tạo, sửa, soft delete category/product
- Product có category, price, stock, image URL và trạng thái active

Order:

- User đăng nhập mới được tạo order
- Order item chỉ gửi `product_id` và `quantity`
- Backend lấy `unit_price` từ database tại thời điểm đặt hàng
- Kiểm tra product tồn tại, active và đủ stock
- Tạo order, order_items và trừ stock trong cùng transaction
- User chỉ xem order của mình
- Admin xem toàn bộ order
- Admin cập nhật trạng thái order

Luồng trạng thái order:

```text
pending -> confirmed
pending -> cancelled
confirmed -> shipping
confirmed -> cancelled
shipping -> completed
```

## Frontend

Frontend không chỉ là demo sơ sài. Đây là giao diện React + Vite có thể thao tác thật với backend cho các flow chính của hệ thống.

### Layout Frontend

Public layout:

- Header store
- Product list
- Product detail
- Login/Register
- Cart link

User layout:

- Cart
- My Orders
- Order Detail
- Profile
- Logout

Admin layout:

- Sidebar admin
- Header admin
- Dashboard
- Categories
- Products
- Orders
- Users
- Back to Store
- Logout

### Chức Năng Theo Role

Guest:

- Xem danh sách sản phẩm
- Tìm kiếm/lọc sản phẩm
- Xem chi tiết sản phẩm
- Thêm sản phẩm vào giỏ local
- Đăng ký
- Đăng nhập

User:

- Có toàn bộ quyền Guest
- Quản lý giỏ hàng
- Cập nhật số lượng sản phẩm
- Tạo đơn hàng
- Xem đơn hàng của tôi
- Xem chi tiết đơn hàng
- Đăng xuất

Admin:

- Dashboard thống kê cơ bản
- Quản lý danh mục
- Quản lý sản phẩm
- Quản lý đơn hàng
- Cập nhật trạng thái đơn hàng đúng luồng
- Quản lý user nếu backend user API được bật

### Cấu Trúc Thư Mục Frontend

```text
frontend
├── public
├── src
│   ├── api
│   │   ├── apiClient.js
│   │   ├── authApi.js
│   │   ├── categoryApi.js
│   │   ├── productApi.js
│   │   ├── orderApi.js
│   │   └── userApi.js
│   ├── components
│   │   ├── common
│   │   └── products
│   ├── contexts
│   │   ├── AuthContext.jsx
│   │   └── CartContext.jsx
│   ├── hooks
│   ├── layouts
│   ├── pages
│   │   ├── admin
│   │   ├── auth
│   │   ├── public
│   │   └── user
│   ├── routes
│   ├── styles
│   └── utils
├── .env.example
├── package.json
└── vite.config.js
```

### Frontend API Base URL

Frontend gọi backend qua biến môi trường:

```text
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

File ví dụ:

```text
frontend/.env.example
```

Nội dung:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

Nếu muốn test production build local, có thể tạo:

```text
frontend/.env.production
```

Ví dụ:

```env
VITE_API_BASE_URL=https://your-backend-domain.onrender.com/api/v1
```

Lưu ý: Vite chỉ expose biến môi trường có prefix `VITE_`.

## Chạy Local

### Backend

Chạy bằng Docker Compose:

```bash
docker compose up --build
```

Hoặc chạy backend local:

```bash
cp .env.example .env
go mod tidy
go run ./cmd/api
```

Health check:

```bash
curl http://localhost:8080/health
```

Chạy test backend:

```bash
go test ./...
```

### Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Mặc định Vite chạy tại:

```text
http://localhost:5173
```

Kiểm tra build local:

```bash
cd frontend
npm run build
npm run preview
```

`npm run preview` dùng để kiểm tra bản build production ở local trước khi deploy.

## Tài Khoản Admin Mặc Định

```text
email: admin@example.com
password: 123456
```

## Một Số API Ví Dụ

Login:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"123456"}'
```

Xem sản phẩm:

```bash
curl "http://localhost:8080/api/v1/products?page=1&limit=12&keyword=&category_id=&min_price=&max_price="
```

Tạo order:

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -d '{"items":[{"product_id":1,"quantity":2}]}'
```

Admin cập nhật trạng thái order:

```bash
curl -X PUT http://localhost:8080/api/v1/orders/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -d '{"status":"confirmed"}'
```

## Deploy Backend Và Database

Gợi ý triển khai demo production:

- Database: Supabase hoặc Neon PostgreSQL
- Backend: Render
- Frontend: Vercel

Backend cần các biến môi trường:

```text
PORT
DATABASE_URL
JWT_ACCESS_SECRET
JWT_REFRESH_SECRET
FRONTEND_URL
ACCESS_TOKEN_MINUTES
REFRESH_TOKEN_HOURS
```

`FRONTEND_URL` phải là domain frontend production, ví dụ:

```text
FRONTEND_URL=https://your-frontend-domain.vercel.app
```

Biến này dùng cho CORS để frontend trên Vercel gọi được backend trên Render.

## Deploy Frontend React + Vite Lên Vercel

### 1. Chuẩn Bị Frontend Trước Deploy

Kiểm tra frontend chạy local:

```bash
cd frontend
npm install
npm run dev
```

Kiểm tra production build:

```bash
npm run build
npm run preview
```

Đảm bảo các flow chính chạy được:

- Home/product list
- Product detail
- Login/register
- Cart
- My orders
- Admin dashboard
- Admin products/orders

### 2. Cấu Hình API URL

Frontend dùng:

```text
VITE_API_BASE_URL=https://your-backend-domain.onrender.com/api/v1
```

Không dùng `localhost` trong production.

Nếu muốn test local production build với backend Render, tạo file:

```text
frontend/.env.production
```

Nội dung ví dụ:

```env
VITE_API_BASE_URL=https://your-backend-domain.onrender.com/api/v1
```

### 3. Deploy Bằng GitHub + Vercel

Các bước:

1. Push code lên GitHub.
2. Vào Vercel, chọn `Add New Project`.
3. Import repository GitHub.
4. Chọn thư mục frontend làm root:

```text
frontend
```

5. Framework Preset:

```text
Vite
```

6. Build Command:

```text
npm run build
```

7. Output Directory:

```text
dist
```

8. Thêm Environment Variable trên Vercel:

```text
VITE_API_BASE_URL=https://your-backend-domain.onrender.com/api/v1
```

9. Deploy.

### 4. Cấu Hình CORS Backend

Sau khi Vercel deploy xong, lấy domain frontend:

```text
https://your-frontend-domain.vercel.app
```

Trên Render backend, cấu hình:

```text
FRONTEND_URL=https://your-frontend-domain.vercel.app
```

Sau đó redeploy backend trên Render.

### 5. Kiểm Tra Sau Deploy

Kiểm tra theo thứ tự:

1. Mở frontend Vercel.
2. Trang home/product list tải được sản phẩm.
3. Search/filter sản phẩm chạy được.
4. Product detail chạy được.
5. Register user mới.
6. Login user.
7. Thêm sản phẩm vào cart.
8. Tạo order.
9. Xem My Orders.
10. Login admin.
11. Vào Admin Dashboard.
12. Quản lý category/product/order.
13. Cập nhật trạng thái order.

## Sửa Lỗi Reload Route Trên Vercel

Nếu reload trực tiếp route như:

```text
https://your-frontend-domain.vercel.app/admin/orders/1
```

mà bị 404, thêm file:

```text
frontend/vercel.json
```

Nội dung:

```json
{
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ]
}
```

Sau đó commit, push và redeploy Vercel.

## Lỗi Deploy Frontend Thường Gặp

Frontend gọi nhầm localhost:

- Nguyên nhân: `VITE_API_BASE_URL` trên Vercel chưa cấu hình hoặc vẫn là `http://localhost:8080/api/v1`.
- Cách sửa: vào Vercel Project Settings -> Environment Variables, đặt lại `VITE_API_BASE_URL`.

CORS error:

- Nguyên nhân: backend Render chưa cho phép domain Vercel.
- Cách sửa: cấu hình `FRONTEND_URL=https://your-frontend-domain.vercel.app` trên Render rồi redeploy backend.

401 do token:

- Nguyên nhân: token hết hạn, refresh token hết hạn hoặc JWT secret thay đổi sau deploy.
- Cách sửa: logout, login lại; kiểm tra `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET` trên Render.

404 khi reload route React Router:

- Nguyên nhân: Vercel không biết route client-side.
- Cách sửa: thêm `frontend/vercel.json` với rewrite về `/index.html`.

Env Vite không nhận:

- Nguyên nhân: biến môi trường thiếu prefix `VITE_`.
- Cách sửa: dùng `VITE_API_BASE_URL`, không dùng `API_BASE_URL`.

Backend chưa bật HTTPS hoặc sai API URL:

- Nguyên nhân: frontend HTTPS gọi backend HTTP hoặc URL thiếu `/api/v1`.
- Cách sửa: dùng domain HTTPS của Render và đúng path `/api/v1`.

## Ghi Chú Bảo Mật Frontend

Bản demo hiện lưu access token, refresh token và cart trong `localStorage` để đơn giản, dễ hiểu và phù hợp đồ án thực tập.

Với production thật nên cân nhắc:

- Lưu refresh token bằng `httpOnly secure cookie`.
- Giảm dữ liệu nhạy cảm lưu trong browser.
- Cấu hình CORS chặt theo domain production.
- Bật HTTPS cho toàn bộ frontend/backend.

## Ảnh Demo

Nếu có ảnh demo, nên đặt trong:

```text
docs/images
```

Gợi ý ảnh cần chụp:

- Product list
- Product detail
- Cart
- My orders
- Admin dashboard
- Admin product management
- Admin order management

Ví dụ nhúng ảnh trong README:

```md
![Product list](docs/images/product-list.png)
![Admin dashboard](docs/images/admin-dashboard.png)
```

## Checklist Frontend Hoàn Thành

- Product list hiển thị grid responsive.
- Product detail có chọn số lượng.
- Search/filter/pagination hoạt động.
- Cart lưu localStorage.
- Cart chặn quantity không hợp lệ.
- Tạo order chỉ gửi `product_id` và `quantity`.
- User xem được My Orders và Order Detail.
- Admin route chỉ cho role admin.
- Admin có sidebar, dashboard, table, form.
- Admin quản lý category/product/order/user.
- Admin cập nhật status order đúng flow.
- Loading, error, empty state đầy đủ.
- Format tiền VND.
- Format ngày tháng rõ ràng.
- Build production thành công bằng `npm run build`.
- Deploy Vercel dùng đúng `VITE_API_BASE_URL`.
- Backend Render cấu hình đúng `FRONTEND_URL`.

## Checklist Deploy Frontend Thành Công

- Backend Render hoạt động và có HTTPS.
- Database Supabase/Neon đã migrate và seed dữ liệu cần thiết.
- Frontend build local thành công.
- Vercel root directory là `frontend`.
- Vercel có biến `VITE_API_BASE_URL=https://your-backend-domain.onrender.com/api/v1`.
- Render có biến `FRONTEND_URL=https://your-frontend-domain.vercel.app`.
- Home/product list gọi API thành công.
- Login/register thành công.
- User tạo order thành công.
- Admin truy cập dashboard thành công.
- Reload route React Router không bị 404.

## Tài Liệu

- API docs: `docs/api.md`
- ERD: `docs/ERD.md`
- Phân tích yêu cầu: `docs/ANALYSIS.md`
- Kiến trúc: `docs/ARCHITECTURE.md`
- Thiết kế database: `docs/DATABASE_DESIGN.md`
- Quy chuẩn project: `AGENTS_P1.md`
