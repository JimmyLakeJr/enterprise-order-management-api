import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { useCart } from "../contexts/CartContext";

export default function UserLayout() {
  const { user, isAdmin, logout } = useAuth();
  const { totalItems } = useCart();
  const navigate = useNavigate();

  async function handleLogout() {
    await logout();
    navigate("/login");
  }

  return (
    <>
      <header className="header">
        <div className="container header-inner">
          <NavLink to="/" className="brand">
            Enterprise OMS
          </NavLink>
          <nav className="nav">
            <NavLink to="/">Products</NavLink>
            <NavLink to="/cart">Cart ({totalItems})</NavLink>
            <NavLink to="/my-orders">My Orders</NavLink>
            <NavLink to="/profile">Account</NavLink>
            {isAdmin && <NavLink to="/admin">Admin</NavLink>}
            <span className="muted">{user?.name}</span>
            <button type="button" onClick={handleLogout}>
              Logout
            </button>
          </nav>
        </div>
      </header>
      <main className="container page">
        <Outlet />
      </main>
    </>
  );
}
