import { Link } from "react-router-dom";
import { getMyOrders } from "../../api/orderApi";
import Badge from "../../components/common/Badge";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Loading from "../../components/common/Loading";
import Table from "../../components/common/Table";
import { useAsync } from "../../hooks/useAsync";
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

export default function MyOrdersPage() {
  const { data: orders, loading, error } = useAsync(getMyOrders, []);

  if (loading) return <Loading />;
  if (error) return <ErrorMessage message={error} />;
  if (!orders?.length) return <EmptyState title="Chưa có đơn hàng" description="Các đơn hàng bạn tạo sẽ hiển thị tại đây." />;

  return (
    <Card>
      <div className="page-header">
        <div>
          <h1>Đơn hàng của tôi</h1>
          <p className="muted">Theo dõi trạng thái và xem chi tiết từng đơn hàng.</p>
        </div>
      </div>

      <Table
        rows={orders}
        columns={[
          {
            key: "id",
            title: "Order",
            render: (order) => <Link to={`/orders/${order.id}`}>#{order.code || order.id}</Link>,
          },
          {
            key: "total_amount",
            title: "Total",
            render: (order) => formatCurrency(order.total_amount),
          },
          {
            key: "status",
            title: "Status",
            render: (order) => <Badge tone={getStatusTone(order.status)}>{order.status}</Badge>,
          },
          {
            key: "created_at",
            title: "Created At",
            render: (order) => formatDate(order.created_at),
          },
          {
            key: "action",
            title: "Action",
            render: (order) => (
              <Link className="btn btn-secondary" to={`/orders/${order.id}`}>
                Xem chi tiết
              </Link>
            ),
          },
        ]}
      />
    </Card>
  );
}
