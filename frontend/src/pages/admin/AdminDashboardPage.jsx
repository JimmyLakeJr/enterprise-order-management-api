import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { categoryApi, productApi } from "../../api/catalogApi";
import { getMessage } from "../../api/apiClient";
import { orderApi } from "../../api/orderApi";
import { userApi } from "../../api/userApi";
import Badge from "../../components/common/Badge";
import Card from "../../components/common/Card";
import ErrorMessage from "../../components/common/ErrorMessage";
import Loading from "../../components/common/Loading";
import Table from "../../components/common/Table";
import { formatCurrency, formatDate } from "../../utils/format";

function getStatusTone(status) {
  const tones = {
    pending: "warning",
    confirmed: "primary",
    shipping: "info",
    completed: "success",
    cancelled: "danger",
  };
  return tones[status] || "default";
}

export default function AdminDashboardPage() {
  const [stats, setStats] = useState({
    products: 0,
    categories: 0,
    orders: 0,
    users: null,
  });
  const [recentOrders, setRecentOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function loadDashboard() {
    setLoading(true);
    setError("");

    try {
      const [categoryList, productResult, orderList, userResult] = await Promise.all([
        categoryApi.list(),
        productApi.list({ page: 1, limit: 100 }),
        orderApi.list(),
        userApi.list({ page: 1, limit: 100 }).catch(() => null),
      ]);

      setStats({
        products: productResult.meta?.total ?? productResult.data.length,
        categories: categoryList.length,
        orders: orderList?.length || 0,
        users: userResult?.meta?.total ?? userResult?.data?.length ?? null,
      });
      setRecentOrders((orderList || []).slice(0, 5));
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadDashboard();
  }, []);

  if (loading) return <Loading />;

  return (
    <>
      <div className="page-header">
        <div>
          <h1>Dashboard</h1>
          <p className="muted">Tổng quan dữ liệu quản trị sản phẩm, danh mục, đơn hàng và user.</p>
        </div>
      </div>

      <ErrorMessage message={error} />

      <div className="admin-stats">
        <Card>
          <span className="muted">Tổng sản phẩm</span>
          <strong>{stats.products}</strong>
          <Link to="/admin/products">Quản lý sản phẩm</Link>
        </Card>
        <Card>
          <span className="muted">Tổng danh mục</span>
          <strong>{stats.categories}</strong>
          <Link to="/admin/categories">Quản lý danh mục</Link>
        </Card>
        <Card>
          <span className="muted">Tổng đơn hàng</span>
          <strong>{stats.orders}</strong>
          <Link to="/admin/orders">Quản lý đơn hàng</Link>
        </Card>
        <Card>
          <span className="muted">Tổng user</span>
          <strong>{stats.users ?? "N/A"}</strong>
          <Link to="/admin/users">Quản lý user</Link>
        </Card>
      </div>

      <Card>
        <div className="page-header">
          <div>
            <h2>Đơn hàng gần đây</h2>
            <p className="muted">Tính tạm từ API danh sách đơn hàng.</p>
          </div>
        </div>

        {recentOrders.length === 0 ? (
          <p className="muted">Chưa có đơn hàng để hiển thị.</p>
        ) : (
          <Table
            rows={recentOrders}
            columns={[
              { key: "id", title: "Order", render: (order) => <Link to={`/admin/orders/${order.id}`}>#{order.id}</Link> },
              { key: "user", title: "User", render: (order) => order.user?.email || order.user_id || "N/A" },
              { key: "total", title: "Total", render: (order) => formatCurrency(order.total_amount) },
              { key: "status", title: "Status", render: (order) => <Badge tone={getStatusTone(order.status)}>{order.status}</Badge> },
              { key: "created_at", title: "Created", render: (order) => formatDate(order.created_at) || "N/A" },
            ]}
          />
        )}
      </Card>
    </>
  );
}
