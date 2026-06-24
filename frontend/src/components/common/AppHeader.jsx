import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";

const guestLinks = [
  { to: "/", label: "Trang chủ", end: true },
  { to: "/products", label: "Sản phẩm" },
  { to: "/cart", label: "Giỏ hàng", action: true },
  { to: "/login", label: "Đăng nhập", action: true },
  { to: "/register", label: "Đăng ký", action: true },
];

export default function AppHeader() {
  const { user, isAuthenticated, isAdmin, logout } = useAuth();
  const { totalItems } = useCart();
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const userLinks = [
    { to: "/", label: "Trang chủ", end: true },
    { to: "/products", label: "Sản phẩm" },
    { to: "/cart", label: `Giỏ hàng (${totalItems})` },
    { to: "/my-orders", label: "Đơn hàng của tôi" },
    { to: "/profile", label: "Hồ sơ", action: true },
  ];

  if (isAdmin) {
    userLinks.unshift({ to: "/admin", label: "Quản trị", action: true });
  }

  async function handleLogout() {
    await logout();
    setMenuOpen(false);
    navigate("/login", { replace: true });
  }

  const links = isAuthenticated ? userLinks : guestLinks;

  return (
    <header className="header">
      <div className="container header-inner">
        <NavLink to="/" className="brand" onClick={() => setMenuOpen(false)}>
          Enterprise OMS
        </NavLink>
        <button
          type="button"
          className="nav-toggle"
          aria-label={menuOpen ? "Đóng menu" : "Mở menu"}
          aria-expanded={menuOpen}
          aria-controls="main-navigation"
          onClick={() => setMenuOpen((current) => !current)}
        >
          <span />
          <span />
          <span />
        </button>
        <nav id="main-navigation" className={`nav ${menuOpen ? "nav-open" : ""}`} aria-label="Điều hướng chính">
          {links.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              end={link.end}
              className={({ isActive }) => [link.action && "nav-action", isActive && "active"].filter(Boolean).join(" ")}
              onClick={() => setMenuOpen(false)}
            >
              {link.label}
            </NavLink>
          ))}
          {isAuthenticated && (
            <>
              <span className="nav-user">{user?.name}</span>
              <button type="button" className="nav-action nav-action-danger" onClick={handleLogout}>
                Đăng xuất
              </button>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
