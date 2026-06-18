import { Navigate, Outlet, useLocation } from "react-router-dom";
import Loading from "../components/common/Loading";
import { useAuth } from "../contexts/AuthContext";

export default function AdminRoute() {
  const { isAuthenticated, isAdmin, loading } = useAuth();
  const location = useLocation();
  if (loading) return <div className="route-state"><Loading label="Đang kiểm tra quyền quản trị..." /></div>;
  if (!isAuthenticated) return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  return isAdmin ? <Outlet /> : <Navigate to="/forbidden" replace state={{ from: location.pathname }} />;
}
