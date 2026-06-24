# UI UX Optimization Audit

## 1. Tong quan UI hien tai

### Layout hien co

- `frontend/src/layouts/PublicLayout.jsx`: layout cho guest va user thong thuong, dung `AppHeader` va footer chung.
- `frontend/src/layouts/UserLayout.jsx`: layout cho khu vuc user, giu flow mua hang va tai khoan.
- `frontend/src/layouts/AdminLayout.jsx`: layout dashboard quan tri voi sidebar, header va noi dung admin.

### Man hinh public

- `frontend/src/pages/public/ProductListPage.jsx`
- `frontend/src/pages/public/ProductDetailPage.jsx`
- `frontend/src/pages/auth/LoginPage.jsx`
- `frontend/src/pages/auth/RegisterPage.jsx`
- `frontend/src/pages/auth/GoogleAuthCallbackPage.jsx`

### Man hinh user

- `frontend/src/pages/user/CartPage.jsx`
- `frontend/src/pages/user/MyOrdersPage.jsx`
- `frontend/src/pages/user/OrderDetailPage.jsx`
- `frontend/src/pages/user/ProfilePage.jsx`

### Man hinh admin

- `frontend/src/pages/admin/AdminDashboardPage.jsx`
- `frontend/src/pages/admin/AdminCategoriesPage.jsx`
- `frontend/src/pages/admin/AdminProductsPage.jsx`
- `frontend/src/pages/admin/AdminOrdersPage.jsx`
- `frontend/src/pages/admin/AdminOrderDetailPage.jsx`
- `frontend/src/pages/admin/AdminUsersPage.jsx`

### Component dung chung

- `frontend/src/components/common/Button.jsx`
- `frontend/src/components/common/Input.jsx`
- `frontend/src/components/common/Select.jsx`
- `frontend/src/components/common/Textarea.jsx`
- `frontend/src/components/common/Badge.jsx`
- `frontend/src/components/common/Card.jsx`
- `frontend/src/components/common/GlassCard.jsx`
- `frontend/src/components/common/Table.jsx`
- `frontend/src/components/common/Loading.jsx`
- `frontend/src/components/common/ErrorMessage.jsx`
- `frontend/src/components/common/EmptyState.jsx`
- `frontend/src/components/common/Modal.jsx`
- `frontend/src/components/common/ConfirmDialog.jsx`
- `frontend/src/components/common/Toast.jsx`
- `frontend/src/components/common/AppHeader.jsx`
- `frontend/src/components/products/ProductCard.jsx`

### CSS file chinh

- `frontend/src/styles/global.css`
- `frontend/src/styles/layout.css`
- `frontend/src/styles/components.css`
- `frontend/src/styles/pages.css`
- `frontend/src/styles/admin.css`

### Diem manh hien tai

- Da co du flow public, user, admin va route guard ro rang.
- Common component co ban da du de tao mot design system nho.
- Da co `focus-visible`, `prefers-reduced-motion`, loading, error, empty state.
- Admin pages da co phan trang cho orders va users, khong con qua so khai ve flow.
- Product, cart, profile, order status da co tang UX co ban va copy kha ro.

### Diem yeu hien tai

- Token mau, shadow, radius dang nam lon trong `global.css`, chua tach thanh design tokens rieng.
- Visual direction hien tai nghieng ve pastel/liquid glass showcase, chua that su data-dense va enterprise.
- Typography dang dung `Inter`/system, hierarchy chua tao du cam giac control room cho dashboard.
- Hover scale, blur, gradient va shadow dang hoi nhieu o mot so khu vuc.
- Header public, product cards, auth shell, admin cards va bang du lieu chua cung mot ngon ngu thiet ke.
- Table, filter, card va form da dep nhung chua du nghiem tuc cho demo he thong quan tri doanh nghiep.

## 2. Van de UI/UX can sua

### Mau sac va style

- Mau hien tai co do dong bo co ban, nhung phoi mau lilac/xanh mint/pastel tao cam giac showcase hon la van hanh.
- Nhieu gradient va blur cung xuat hien o header, page hero, card, sidebar, modal, toast.
- Can chuyen sang palette slate/navy/ice voi 1 accent xanh la de hop dashboard ton kho va don hang.

### Typography

- Heading, label, table header va so lieu dashboard chua tach bach manh.
- Can bo sung cap chu ro hon cho:
  - display metric
  - page title
  - eyebrow/section label
  - body/caption/meta

### Components

- Button dang co hover scale va glow kha manh.
- Input/select/textarea chua co mot surface enterprise that su on dinh.
- Table dep nhung header va row hover van hoi decorative.
- Badge/status can thong nhat mau va border ro hon de de doc nhanh.
- Modal/toast/card dang glass kha nang va bong.

### Layout va navigation

- `AppHeader` du dung nhung chua to chuc theo huong san pham doanh nghiep.
- Mobile nav on ve ky thuat, nhung visual state chua that ro.
- `AdminLayout` tot ve cau truc, nhung sidebar va stat cards van nhieu hieu ung showcase.

### Page-level UX

- Product list can ro bo loc, ket qua, empty state va the san pham.
- Product detail can lam ro price, stock, quantity, CTA.
- Cart va My Orders can ro hierarchy hon giua thong tin va thao tac.
- Profile page co preview-only media, can tao cam giac "feature duoc khoa" ro hon.
- Admin orders/users/categories/products can giam nhieu visual noise de uu tien data.

### Accessibility va quality

- Co `focus-visible` va `prefers-reduced-motion`, day la nen tang tot.
- Can dam bao clickable element nao cung co `cursor: pointer` va state ro hon.
- Contrast hien tai kha on nhung mot so lop glass + muted text co nguy co yeu tren man hinh laptop.
- Animation dang nhieu o card hover/page rise/card arrive; can ha cuong do xuong 150-300ms va giam scale.
- CSS dang phan tan, co nguy co lap lai va chong cheo token/style.

## 3. Design system de xuat

### Huong thiet ke

- Product type: enterprise order management, product/order/admin dashboard, business management tool.
- Chon style: controlled liquid glass + minimal Swiss enterprise + bento dashboard cards.
- Khong dung cyberpunk, neon, y2k, chaos visuals.

### Co so tham chieu tu UI UX Pro Max Skill

- Pattern phu hop: data-dense dashboard.
- Palette de xuat: navy/slate background, slate surface, xanh la lam accent thanh cong, text sang ro, border lanh.
- Uu tien: clarity, hierarchy, responsiveness, accessibility, motion restraint.

### Token can co

- Color tokens:
  - background app
  - elevated surface
  - glass surface nhe
  - primary
  - secondary
  - success
  - warning
  - danger
  - info
  - text
  - text muted
  - border
- Typography:
  - sans stack co tinh enterprise ro hon
  - display/page/title/body/caption scale
- Spacing:
  - `--space-1` den `--space-10`
- Radius:
  - nho, vua, lon, xl
- Shadow/elevation:
  - shadow nhe cho form/table
  - shadow vua cho card
  - shadow lon cho modal/toast
- Motion:
  - `--transition-fast`
  - `--transition-normal`
  - easing mem, khong bounce

### Quy tac style

- Glass chi dung muc nhe, uu tien border va surface hon blur.
- Dashboard cards theo huong bento nhung phang, gon, ro so lieu.
- Table header va filter toolbar phai doc nhanh.
- Button:
  - primary ro rang
  - secondary trung tinh
  - danger canh bao ro
- Status badge:
  - dung nen nhat + border + text dam
- Empty/loading/error:
  - thong nhat cung mot ngon ngu hinh khoi
- Breakpoints:
  - 375
  - 768
  - 1024
  - 1440

### Accessibility rules

- Focus-visible ro va nhat quan.
- Khong chi dung mau de bieu dat trang thai.
- Contrast huong toi muc WCAG AA.
- Respect `prefers-reduced-motion`.
- Clickable phai co `cursor: pointer`.

## 4. Danh sach file se sua

### CSS se chinh

- `frontend/src/main.jsx`
- `frontend/src/styles/tokens.css` (tao moi)
- `frontend/src/styles/global.css`
- `frontend/src/styles/layout.css`
- `frontend/src/styles/components.css`
- `frontend/src/styles/pages.css`
- `frontend/src/styles/admin.css`

### Component se chuan hoa

- `frontend/src/components/common/Button.jsx`
- `frontend/src/components/common/Input.jsx`
- `frontend/src/components/common/Table.jsx`
- `frontend/src/components/common/Badge.jsx`
- `frontend/src/components/common/AppHeader.jsx`
- `frontend/src/components/products/ProductCard.jsx`

### Page/layout se toi uu

- `frontend/src/layouts/AdminLayout.jsx`
- `frontend/src/pages/auth/LoginPage.jsx`
- `frontend/src/pages/auth/RegisterPage.jsx`
- `frontend/src/pages/user/ProfilePage.jsx`
- `frontend/src/pages/admin/AdminDashboardPage.jsx`
- `frontend/src/pages/admin/AdminCategoriesPage.jsx`
- `frontend/src/pages/admin/AdminOrdersPage.jsx`
- `frontend/src/pages/admin/AdminUsersPage.jsx`

### Component co the tao moi neu can

- Co the tao mot `StatusBadge` hoac `Pagination` chung neu phat sinh lap logic UI.
- Hien tai uu tien tai su dung component da co, chi tao moi neu that su giam lap.

### Component khong nen xoa vo i

- Khong xoa component chung neu chua xac minh import path.
- Khong cham vao `AuthContext`, `CartContext`, route guards, API clients neu khong can.

### Rui ro anh huong frontend flow

- Doi class/style o `Button`, `Input`, `Table`, `AppHeader`, `AdminLayout` co the anh huong nhieu man hinh.
- Chinh page hero va card surface can giu nguyen structure JSX de tranh vo layout/pages khac.
- Khong doi route, khong doi props API, khong doi business flow.
