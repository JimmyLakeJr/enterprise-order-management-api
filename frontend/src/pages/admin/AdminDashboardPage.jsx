import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { categoryApi } from "../../api/categoryApi";
import { orderApi } from "../../api/orderApi";
import { productApi } from "../../api/productApi";
import { userApi } from "../../api/userApi";
import Badge from "../../components/common/Badge";
import Card from "../../components/common/Card";
import ErrorMessage from "../../components/common/ErrorMessage";
import Loading from "../../components/common/Loading";
import Table from "../../components/common/Table";
import { formatDate } from "../../utils/format";
import { formatCurrency } from "../../utils/formatCurrency";
import { getOrderStatus } from "../../utils/orderStatus";

const EMPTY_STATS = { products: null, categories: null, orders: null, users: null };

export default function AdminDashboardPage() {
  const [stats, setStats] = useState(EMPTY_STATS);
  const [recentOrders, setRecentOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let active = true;

    Promise.allSettled([
      categoryApi.list(),
      productApi.list({ page: 1, limit: 1 }),
      orderApi.list({ page: 1, limit: 5 }),
      userApi.list({ page: 1, limit: 1 }),
    ]).then(([categoryResult, productResult, orderResult, userResult]) => {
      if (!active) return;

      const categories = categoryResult.status === "fulfilled" ? categoryResult.value : null;
      const products = productResult.status === "fulfilled" ? productResult.value : null;
      const orders = orderResult.status === "fulfilled" ? orderResult.value : null;
      const users = userResult.status === "fulfilled" ? userResult.value : null;
      const failedCount = [categoryResult, productResult, orderResult, userResult].filter(
        (result) => result.status === "rejected"
      ).length;

      setStats({
        categories: categories?.length ?? null,
        products: products?.meta?.total ?? products?.data?.length ?? null,
        orders: orders?.meta?.total ?? orders?.data?.length ?? null,
        users: users?.meta?.total ?? users?.data?.length ?? null,
      });
      setRecentOrders(orders?.data || []);
      setError(failedCount ? `${failedCount} nguồn dữ liệu chưa tải được. Các số liệu còn lại vẫn được hiển thị.` : "");
      setLoading(false);
    });

    return () => {
      active = false;
    };
  }, []);

  const statCards = [
    { key: "products", label: "Sản phẩm hoạt động", to: "/admin/products", note: "Danh mục hàng hóa hiện có" },
    { key: "categories", label: "Danh mục hoạt động", to: "/admin/categories", note: "Nhóm sản phẩm đang hiển thị" },
    { key: "orders", label: "Tất cả đơn hàng", to: "/admin/orders", note: "Theo dõi trạng thái xử lý" },
    { key: "users", label: "Người dùng hoạt động", to: "/admin/users", note: "Tài khoản đang sử dụng" },
  ];

  return (
    <>
      <div className="page-header admin-page-title">
        <div>
          <span className="eyebrow">Tổng quan</span>
          <h1>Tổng quan vận hành</h1>
          <p className="muted">Theo dõi nhanh sản phẩm, danh mục, đơn hàng và người dùng từ một màn hình.</p>
        </div>
      </div>

      <ErrorMessage message={error} />

      {loading ? (
        <Loading label="Đang tổng hợp dữ liệu quản trị..." variant="dashboard" count={4} />
      ) : (
        <>
          <div className="admin-stats">
            {statCards.map((stat) => (
              <Card key={stat.key} className="admin-stat-card">
                <span>{stat.label}</span>
                <strong>{stats[stat.key] ?? "—"}</strong>
                <small>{stat.note}</small>
                <Link to={stat.to}>Mở quản lý →</Link>
              </Card>
            ))}
          </div>

          <Card className="admin-panel-card">
            <div className="page-header compact-header">
              <div>
                <h2>Đơn hàng gần đây</h2>
                <p className="muted">Tóm tắt các đơn hàng mới nhất để theo dõi tiến độ xử lý.</p>
              </div>
              <Link className="btn btn-secondary" to="/admin/orders">
                Xem tất cả
              </Link>
            </div>

            {recentOrders.length === 0 ? (
              <p className="muted">Chưa có đơn hàng để hiển thị.</p>
            ) : (
              <Table
                rows={recentOrders}
                columns={[
                  {
                    key: "id",
                    title: "Mã đơn",
                    render: (order) => (
                      <Link className="order-link" to={`/admin/orders/${order.id}`}>
                        #{order.id}
                      </Link>
                    ),
                  },
                  {
                    key: "user",
                    title: "Người dùng",
                    render: (order) => order.user?.name || `Mã #${order.user_id}`,
                  },
                  { key: "items", title: "Sản phẩm", render: (order) => `${order.items?.length || 0} dòng` },
                  { key: "created_at", title: "Ngày tạo", render: (order) => formatDate(order.created_at) || "—" },
                  { key: "total", title: "Tổng tiền", render: (order) => formatCurrency(order.total_amount) },
                  {
                    key: "status",
                    title: "Trạng thái",
                    render: (order) => {
                      const status = getOrderStatus(order.status);
                      return <Badge tone={status.tone}>{status.label}</Badge>;
                    },
                  },
                ]}
              />
            )}
          </Card>
        </>
      )}
    </>
  );
}
