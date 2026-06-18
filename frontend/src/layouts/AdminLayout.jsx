import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

const menuItems = [
  { to: "/admin", label: "Admin Dashboard", end: true },
  { to: "/admin/categories", label: "Categories" },
  { to: "/admin/products", label: "Products" },
  { to: "/admin/orders", label: "Orders" },
  { to: "/admin/users", label: "Users" },
  { to: "/profile", label: "Profile" },
];

export default function AdminLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  async function handleLogout() {
    await logout();
    navigate("/login");
  }

  return (
    <div className="admin-layout">
      <aside className="sidebar">
        <NavLink to="/admin" className="brand">
          Admin Panel
        </NavLink>

        <nav className="admin-menu">
          {menuItems.map((item) => (
            <NavLink key={item.to} to={item.to} end={item.end}>
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="sidebar-footer">
          <button type="button" className="btn btn-secondary" onClick={() => navigate("/")}>
            Back to Store
          </button>
          <button type="button" className="btn btn-danger" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </aside>

      <section className="admin-main">
        <header className="admin-header">
          <div>
            <strong>Enterprise Order Management</strong>
            <p className="muted">Product, order and user administration</p>
          </div>
          <div className="admin-user">
            <span>{user?.name || "Admin"}</span>
            <small className="muted">{user?.email}</small>
          </div>
        </header>

        <main className="admin-content">
          <Outlet />
        </main>
      </section>
    </div>
  );
}
