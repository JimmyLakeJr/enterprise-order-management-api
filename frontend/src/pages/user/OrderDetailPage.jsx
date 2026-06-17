import { Link, useParams } from "react-router-dom";
import { getOrderById } from "../../api/orderApi";
import Badge from "../../components/common/Badge";
import Card from "../../components/common/Card";
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

export default function OrderDetailPage() {
  const { id } = useParams();
  const { data: order, loading, error } = useAsync(() => getOrderById(id), [id]);

  if (loading) return <Loading />;
  if (error) return <ErrorMessage message={error} />;
  if (!order) return <ErrorMessage message="Không tìm thấy đơn hàng." />;

  const items = order.items || order.order_items || [];

  return (
    <Card>
      <div className="page-header">
        <div>
          <Link to="/my-orders" className="muted">
            Quay lại đơn hàng
          </Link>
          <h1>Đơn hàng #{order.code || order.id}</h1>
          <div className="order-meta">
            <Badge tone={getStatusTone(order.status)}>{order.status}</Badge>
            {order.created_at && <span className="muted">{formatDate(order.created_at)}</span>}
          </div>
        </div>
        <strong className="cart-total">{formatCurrency(order.total_amount)}</strong>
      </div>

      <Table
        rows={items}
        getRowKey={(item, index) => item.id || `${item.product_id}-${index}`}
        columns={[
          {
            key: "product",
            title: "Product",
            render: (item) => item.name || item.product_name || item.product?.name || `#${item.product_id}`,
          },
          { key: "unit_price", title: "Unit Price", render: (item) => formatCurrency(item.unit_price) },
          { key: "quantity", title: "Quantity" },
          { key: "subtotal", title: "Subtotal", render: (item) => formatCurrency(item.subtotal) },
        ]}
      />

      <div className="cart-summary">
        <span className="muted">Tổng tiền</span>
        <strong>{formatCurrency(order.total_amount)}</strong>
      </div>
    </Card>
  );
}
