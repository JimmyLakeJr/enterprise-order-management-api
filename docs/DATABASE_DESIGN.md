# Thiết kế database PostgreSQL

Project: **enterprise-order-management-api**

Đề tài: **Phát triển backend API cho hệ thống quản lý sản phẩm và đơn hàng trong doanh nghiệp sử dụng Golang**

## 1. Mục tiêu

Thiết kế cơ sở dữ liệu PostgreSQL cho hệ thống quản lý sản phẩm và đơn hàng, đáp ứng các yêu cầu:

- Không dùng ORM.
- Dùng SQL thuần.
- Phù hợp với Go backend sử dụng `pgxpool`.
- Có primary key, foreign key, unique, not null và check constraint rõ ràng.
- Có index cho các truy vấn thường dùng.
- Có dữ liệu seed ban đầu.
- Hỗ trợ xác thực bằng refresh token.
- Hỗ trợ quản lý sản phẩm, danh mục, đơn hàng và chi tiết đơn hàng.
- Hỗ trợ transaction khi tạo đơn hàng.

## 2. Cấu trúc file/thư mục

```text
migrations/
└── 001_init.sql

docs/
├── DATABASE_DESIGN.md
└── ERD.md
```

## 3. Code hoàn chỉnh

SQL migration chính nằm tại:

```text
migrations/001_init.sql
```

Nội dung migration:

```sql
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role_id BIGINT NOT NULL REFERENCES roles(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    name VARCHAR(150) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price BIGINT NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    image_url TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    total_amount BIGINT NOT NULL CHECK (total_amount >= 0),
    status VARCHAR(30) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'confirmed', 'shipping', 'completed', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price BIGINT NOT NULL CHECK (unit_price >= 0),
    subtotal BIGINT NOT NULL CHECK (subtotal >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);

INSERT INTO roles (name)
VALUES ('admin'), ('user')
ON CONFLICT (name) DO NOTHING;

INSERT INTO users (full_name, email, password_hash, role_id)
VALUES (
    'Admin',
    'admin@example.com',
    '$2a$10$0LCwL/15uA7zMUXusqGU6OofPjnYqvm3.jOBzYBETOrXPALCz4m9q',
    (SELECT id FROM roles WHERE name = 'admin')
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO categories (name, description)
VALUES
    ('Electronics', 'Electronic devices and accessories'),
    ('Office', 'Office supplies and equipment')
ON CONFLICT (name) DO NOTHING;
```

## 4. Giải thích ngắn gọn

### 4.1. ERD dạng text

```text
roles
  1 ──── n users

users
  1 ──── n refresh_tokens
  1 ──── n orders

categories
  1 ──── n products

orders
  1 ──── n order_items

products
  1 ──── n order_items
```

Biểu diễn đầy đủ:

```text
+-------------+        +-------------+
|   roles     | 1    n |    users    |
|-------------|--------|-------------|
| id          |        | id          |
| name        |        | full_name   |
+-------------+        | email       |
                       | role_id     |
                       +-------------+
                              |
                              | 1
                              | n
                       +----------------+
                       | refresh_tokens |
                       +----------------+

+---------------+      +-------------+      +---------------+
|  categories   | 1  n |  products   | 1  n |  order_items  |
+---------------+------|-------------|------|---------------|
| id            |      | id          |      | id            |
| name          |      | category_id |      | order_id      |
| is_active     |      | name        |      | product_id    |
+---------------+      | price       |      | quantity      |
                       | stock       |      | unit_price    |
                       +-------------+      | subtotal      |
                                            +---------------+
                                                   n
                                                   |
                                                   | 1
                                            +-------------+
                                            |   orders    |
                                            +-------------+
                                            | id          |
                                            | user_id     |
                                            | status      |
                                            | total_amount|
                                            +-------------+
```

### 4.2. Quan hệ giữa các bảng

#### roles - users

Một role có thể được gán cho nhiều user. Mỗi user bắt buộc thuộc một role thông qua `users.role_id`.

Ví dụ:

- Admin thuộc role `admin`.
- Người dùng thông thường thuộc role `user`.

#### users - refresh_tokens

Một user có thể có nhiều refresh token. Quan hệ này hỗ trợ đăng nhập trên nhiều thiết bị hoặc nhiều phiên đăng nhập.

Khi user logout, refresh token tương ứng được revoke bằng cách cập nhật `revoked_at`.

#### users - orders

Một user có thể tạo nhiều đơn hàng. Mỗi đơn hàng thuộc về một user thông qua `orders.user_id`.

Quy tắc phân quyền:

- User chỉ được xem đơn hàng của chính mình.
- Admin được xem toàn bộ đơn hàng.

#### categories - products

Một danh mục có thể chứa nhiều sản phẩm. Mỗi sản phẩm bắt buộc thuộc một danh mục thông qua `products.category_id`.

Danh mục inactive không nên hiển thị ở API public.

#### orders - order_items

Một đơn hàng có nhiều dòng chi tiết đơn hàng. Mỗi dòng trong `order_items` đại diện cho một sản phẩm được mua trong đơn hàng.

Nếu order bị xóa vật lý, các order item tương ứng sẽ bị xóa theo nhờ `ON DELETE CASCADE`. Trong nghiệp vụ thực tế, order thường không nên xóa vật lý.

#### products - order_items

Một sản phẩm có thể xuất hiện trong nhiều dòng chi tiết đơn hàng khác nhau. Quan hệ này giúp truy xuất sản phẩm đã được đặt trong từng đơn hàng.

### 4.3. Giải thích các bảng

#### roles

Lưu vai trò trong hệ thống, gồm:

- `admin`
- `user`

#### users

Lưu thông tin tài khoản người dùng.

Các điểm quan trọng:

- `email` là duy nhất.
- `password_hash` lưu mật khẩu đã hash, không lưu plain password.
- `role_id` dùng để phân quyền.
- `is_active` dùng cho soft delete.

#### refresh_tokens

Lưu refresh token đã hash.

Các điểm quan trọng:

- Không lưu raw refresh token.
- `expires_at` kiểm tra token hết hạn.
- `revoked_at` dùng cho logout.
- `token_hash` unique để tránh trùng token.

#### categories

Lưu danh mục sản phẩm.

Các điểm quan trọng:

- `name` là duy nhất.
- `is_active` dùng để ẩn danh mục khỏi API public.

#### products

Lưu sản phẩm.

Các điểm quan trọng:

- `price >= 0`.
- `stock >= 0`.
- `image_url` lưu đường dẫn ảnh sản phẩm.
- `is_active` dùng cho soft delete.

#### orders

Lưu thông tin đơn hàng.

Các trạng thái hợp lệ:

- `pending`
- `confirmed`
- `shipping`
- `completed`
- `cancelled`

#### order_items

Lưu chi tiết từng sản phẩm trong đơn hàng.

Các điểm quan trọng:

- `quantity > 0`.
- `unit_price >= 0`.
- `subtotal >= 0`.
- Có `created_at` để biết thời điểm item được tạo.

### 4.4. Vì sao order_items cần lưu unit_price?

`order_items` cần lưu `unit_price` vì giá sản phẩm có thể thay đổi sau khi đơn hàng được tạo.

Ví dụ:

- Ngày 01, user mua sản phẩm A với giá 100.000.
- Ngày 05, admin cập nhật giá sản phẩm A thành 120.000.

Nếu order item không lưu `unit_price`, khi xem lại đơn hàng cũ hệ thống có thể hiển thị sai giá. Vì vậy, `unit_price` phải lưu giá tại thời điểm đặt hàng để đảm bảo lịch sử đơn hàng chính xác.

### 4.5. Vì sao cần transaction khi tạo order?

Tạo order liên quan nhiều thao tác:

- Kiểm tra sản phẩm.
- Kiểm tra tồn kho.
- Tạo bản ghi `orders`.
- Tạo các bản ghi `order_items`.
- Trừ tồn kho trong `products`.

Các thao tác này phải thành công cùng nhau. Nếu một bước lỗi, toàn bộ dữ liệu phải rollback.

Nếu không dùng transaction, có thể xảy ra lỗi:

- Order được tạo nhưng order_items chưa được tạo.
- Order_items được tạo nhưng stock chưa giảm.
- Stock bị giảm nhưng order tạo thất bại.
- Dữ liệu đơn hàng và tồn kho không đồng nhất.

Do đó, chức năng tạo order bắt buộc phải dùng transaction.

## 5. Cách chạy/test

Chạy PostgreSQL và API bằng Docker Compose:

```bash
docker compose up --build
```

Nếu đã từng chạy database cũ và muốn chạy lại migration từ đầu:

```bash
docker compose down -v
docker compose up --build
```

Kiểm tra migration đã chạy bằng psql:

```bash
docker exec -it enterprise-order-postgres psql -U postgres -d enterprise_order_management
```

Một số câu SQL kiểm tra:

```sql
\dt
SELECT * FROM roles;
SELECT id, full_name, email, role_id, is_active FROM users;
SELECT * FROM categories;
```

## 6. Lỗi thường gặp

### Migration không chạy lại

Nguyên nhân:

- Docker volume của PostgreSQL đã tồn tại.
- PostgreSQL chỉ tự chạy file trong `/docker-entrypoint-initdb.d` khi database được tạo lần đầu.

Cách sửa trong môi trường development:

```bash
docker compose down -v
docker compose up --build
```

### Lỗi duplicate email

Nguyên nhân:

- `users.email` có unique constraint.

Cách sửa:

- Dùng email khác.
- Hoặc kiểm tra dữ liệu user đã tồn tại.

### Lỗi duplicate category name

Nguyên nhân:

- `categories.name` có unique constraint.

Cách sửa:

- Dùng tên danh mục khác.

### Lỗi foreign key khi tạo product

Nguyên nhân:

- `category_id` không tồn tại trong bảng `categories`.

Cách sửa:

- Tạo category trước.
- Hoặc dùng category seed có sẵn.

### Lỗi check constraint price hoặc stock

Nguyên nhân:

- `products.price < 0`.
- `products.stock < 0`.

Cách sửa:

- Đảm bảo price và stock không âm trước khi insert/update.

### Lỗi check constraint quantity

Nguyên nhân:

- `order_items.quantity <= 0`.

Cách sửa:

- Validate request tạo order, quantity phải lớn hơn 0.

## 7. Checklist hoàn thành

- [x] Có bảng `roles`.
- [x] Có bảng `users`.
- [x] Có bảng `refresh_tokens`.
- [x] Có bảng `categories`.
- [x] Có bảng `products`.
- [x] Có bảng `orders`.
- [x] Có bảng `order_items`.
- [x] Có primary key cho tất cả bảng.
- [x] Có foreign key cho các quan hệ chính.
- [x] Có unique constraint cho role name, email, token hash, category name.
- [x] Có not null constraint cho các cột bắt buộc.
- [x] Có check constraint cho price, stock, quantity, status.
- [x] Có index cho users.email.
- [x] Có index cho products.name.
- [x] Có index cho products.category_id.
- [x] Có index cho orders.user_id.
- [x] Có index cho refresh_tokens.user_id.
- [x] Có seed role `admin`, `user`.
- [x] Có seed admin mặc định.
- [x] Có seed category mẫu.
- [x] Có giải thích vì sao lưu `unit_price`.
- [x] Có giải thích vì sao cần transaction khi tạo order.

## 8. Các lỗi thiết kế database cần tránh

- Lưu plain password trong bảng users.
- Lưu raw refresh token trong database.
- Không đặt unique constraint cho email.
- Không đặt foreign key giữa các bảng liên quan.
- Cho phép price hoặc stock âm.
- Cho phép quantity bằng 0 hoặc âm.
- Tin giá sản phẩm từ frontend khi tạo order.
- Không lưu `unit_price` trong `order_items`.
- Không dùng transaction khi tạo order.
- Xóa vật lý product đã từng nằm trong order.
- Thiết kế order chỉ lưu một product, làm hệ thống không hỗ trợ đơn hàng nhiều sản phẩm.
- Không có index cho các cột thường query như email, product name, category_id, user_id.
- Trả lỗi database thô trực tiếp cho client.
