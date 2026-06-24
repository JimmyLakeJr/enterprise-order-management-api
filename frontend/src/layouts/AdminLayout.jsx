import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

const menuItems = [
  { to: "/admin", label: "Tổng quan", shortLabel: "01", end: true },
  { to: "/admin/categories", label: "Danh mục", shortLabel: "02" },
  { to: "/admin/products", label: "Sản phẩm", shortLabel: "03" },
  { to: "/admin/orders", label: "Đơn hàng", shortLabel: "04" },
  { to: "/admin/users", label: "Người dùng", shortLabel: "05" },
];

export default function AdminLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  async function handleLogout() {
    await logout();
    navigate("/login", { replace: true });
  }

  return (
    <div className="admin-layout">
      <aside className="sidebar">
        <div className="sidebar-heading">
          <NavLink to="/admin" className="admin-brand">
            <span className="admin-brand-mark" aria-hidden="true">
              EO
            </span>
            <span>
              <strong>Enterprise OMS</strong>
              <small>Khu vực quản trị vận hành</small>
            </span>
          </NavLink>
        </div>

        <nav className="admin-menu" aria-label="Điều hướng quản trị">
          {menuItems.map((item) => (
            <NavLink key={item.to} to={item.to} end={item.end}>
              <span className="menu-index" aria-hidden="true">
                {item.shortLabel}
              </span>
              <span>{item.label}</span>
            </NavLink>
          ))}
        </nav>

        <div className="sidebar-footer">
          <button type="button" className="btn btn-secondary" onClick={() => navigate("/")}>
            Về giao diện người dùng
          </button>
          <button type="button" className="btn btn-danger" onClick={handleLogout}>
            Đăng xuất
          </button>
        </div>
      </aside>

      <section className="admin-main">
        <header className="admin-header">
          <div>
            <span className="eyebrow">Khu vực quản trị</span>
            <strong>Enterprise Order Management</strong>
          </div>
          <div className="admin-header-actions">
            <NavLink className="btn btn-secondary" to="/profile">
              Hồ sơ
            </NavLink>
            <div className="admin-user">
              <span>{user?.name || "Quản trị viên"}</span>
              <small className="muted">{user?.email || "—"}</small>
            </div>
          </div>
        </header>

        <main className="admin-content">
          <Outlet />
        </main>
      </section>
    </div>
  );
}
