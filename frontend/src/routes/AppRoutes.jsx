import { Navigate, Route, Routes } from "react-router-dom";
import AdminLayout from "../layouts/AdminLayout";
import PublicLayout from "../layouts/PublicLayout";
import UserLayout from "../layouts/UserLayout";
import LoginPage from "../pages/auth/LoginPage";
import RegisterPage from "../pages/auth/RegisterPage";
import AdminCategoriesPage from "../pages/admin/AdminCategoriesPage";
import AdminDashboardPage from "../pages/admin/AdminDashboardPage";
import AdminOrderDetailPage from "../pages/admin/AdminOrderDetailPage";
import AdminOrdersPage from "../pages/admin/AdminOrdersPage";
import AdminProductsPage from "../pages/admin/AdminProductsPage";
import AdminUsersPage from "../pages/admin/AdminUsersPage";
import ForbiddenPage from "../pages/public/ForbiddenPage";
import NotFoundPage from "../pages/public/NotFoundPage";
import ProductDetailPage from "../pages/public/ProductDetailPage";
import ProductListPage from "../pages/public/ProductListPage";
import CartPage from "../pages/user/CartPage";
import MyOrdersPage from "../pages/user/MyOrdersPage";
import OrderDetailPage from "../pages/user/OrderDetailPage";
import ProfilePage from "../pages/user/ProfilePage";
import AdminRoute from "./AdminRoute";
import ProtectedRoute from "./ProtectedRoute";

export default function AppRoutes() {
  return (
    <Routes>
      <Route element={<PublicLayout />}>
        <Route index element={<ProductListPage />} />
        <Route path="products" element={<ProductListPage />} />
        <Route path="products/:id" element={<ProductDetailPage />} />
        <Route path="cart" element={<CartPage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="forbidden" element={<ForbiddenPage />} />
      </Route>

      <Route element={<ProtectedRoute />}>
        <Route element={<UserLayout />}>
          <Route path="my-orders" element={<MyOrdersPage />} />
          <Route path="orders/:id" element={<OrderDetailPage />} />
          <Route path="profile" element={<ProfilePage />} />
        </Route>
      </Route>

      <Route element={<ProtectedRoute />}>
        <Route element={<AdminRoute />}>
          <Route path="admin" element={<AdminLayout />}>
            <Route index element={<AdminDashboardPage />} />
            <Route path="categories" element={<AdminCategoriesPage />} />
            <Route path="products" element={<AdminProductsPage />} />
            <Route path="orders" element={<AdminOrdersPage />} />
            <Route path="orders/:id" element={<AdminOrderDetailPage />} />
            <Route path="users" element={<AdminUsersPage />} />
          </Route>
        </Route>
      </Route>

      <Route path="/home" element={<Navigate to="/" replace />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
