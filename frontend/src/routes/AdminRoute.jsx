import { Navigate, Outlet } from "react-router-dom";
import Loading from "../components/common/Loading";
import { useAuth } from "../contexts/AuthContext";

export default function AdminRoute() {
  const { isAuthenticated, isAdmin, loading } = useAuth();
  if (loading) return <Loading />;
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return isAdmin ? <Outlet /> : <Navigate to="/forbidden" replace />;
}
