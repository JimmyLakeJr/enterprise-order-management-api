import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";

const guestLinks = [
  { to: "/", label: "Home", end: true },
  { to: "/products", label: "Products" },
  { to: "/login", label: "Login", action: true },
  { to: "/register", label: "Register", action: true },
];

export default function AppHeader() {
  const { user, isAuthenticated, isAdmin, logout } = useAuth();
  const { totalItems } = useCart();
  const navigate = useNavigate();

  const userLinks = [
    { to: "/products", label: "Products" },
    { to: "/cart", label: `Cart (${totalItems})` },
    { to: "/my-orders", label: "My Orders" },
    { to: "/profile", label: "Profile", action: true },
  ];

  if (isAdmin) {
    userLinks.unshift({ to: "/admin", label: "Admin Dashboard", action: true });
  }

  async function handleLogout() {
    await logout();
    navigate("/login", { replace: true });
  }

  const links = isAuthenticated ? userLinks : guestLinks;

  return (
    <header className="header">
      <div className="container header-inner">
        <NavLink to="/" className="brand">Enterprise OMS</NavLink>
        <nav className="nav" aria-label="Main navigation">
          {links.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              end={link.end}
              className={({ isActive }) => [link.action && "nav-action", isActive && "active"].filter(Boolean).join(" ")}
            >
              {link.label}
            </NavLink>
          ))}
          {isAuthenticated && (
            <>
              <span className="nav-user">{user?.name}</span>
              <button type="button" className="nav-action nav-action-danger" onClick={handleLogout}>
                Logout
              </button>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
