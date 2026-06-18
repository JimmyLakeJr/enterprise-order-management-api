import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getMessage } from "../../api/apiClient";
import { orderApi } from "../../api/orderApi";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import ErrorMessage from "../../components/common/ErrorMessage";
import Loading from "../../components/common/Loading";
import Select from "../../components/common/Select";
import Table from "../../components/common/Table";
import { getOrderStatusTone, ORDER_STATUS_TRANSITIONS } from "../../constants/domain";
import { formatCurrency } from "../../utils/format";

export default function AdminOrderDetailPage() {
  const { id } = useParams();
  const [order, setOrder] = useState(null);
  const [nextStatus, setNextStatus] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    loadOrder();
  }, [id]);

  async function loadOrder() {
    setLoading(true);
    setError("");
    try {
      const data = await orderApi.detail(id);
      setOrder(data);
      setNextStatus("");
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }

  async function handleUpdateStatus(event) {
    event.preventDefault();
    if (!nextStatus) return;

    setSubmitting(true);
    setError("");
    try {
      await orderApi.updateStatus(id, nextStatus);
      await loadOrder();
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) return <Loading />;
  if (error && !order) return <ErrorMessage message={error} />;
  if (!order) return <ErrorMessage message="Không tìm thấy đơn hàng." />;

  const items = order.items || [];
  const availableStatuses = ORDER_STATUS_TRANSITIONS[order.status] || [];

  return (
    <div className="grid">
      <Card>
        <div className="page-header">
          <div>
            <Link to="/admin/orders" className="muted">
              Quay lại Orders
            </Link>
            <h1>Order #{order.id}</h1>
            <div className="order-meta">
              <Badge tone={getOrderStatusTone(order.status)}>{order.status}</Badge>
            </div>
          </div>
          <strong className="cart-total">{formatCurrency(order.total_amount)}</strong>
        </div>

        <ErrorMessage message={error} />

        <div className="order-info-grid">
          <div>
            <span className="muted">User/Customer</span>
            <strong>User #{order.user_id}</strong>
          </div>
          <div>
            <span className="muted">Total Amount</span>
            <strong>{formatCurrency(order.total_amount)}</strong>
          </div>
          <div>
            <span className="muted">Status</span>
            <strong>{order.status}</strong>
          </div>
        </div>

        {availableStatuses.length === 0 ? (
          <p className="muted">Đơn hàng ở trạng thái cuối, không thể cập nhật tiếp.</p>
        ) : (
          <form className="status-form" onSubmit={handleUpdateStatus}>
            <Select label="Cập nhật trạng thái" value={nextStatus} onChange={(event) => setNextStatus(event.target.value)}>
              <option value="">Chọn trạng thái tiếp theo</option>
              {availableStatuses.map((status) => (
                <option key={status} value={status}>
                  {status}
                </option>
              ))}
            </Select>
            <Button type="submit" disabled={!nextStatus || submitting}>
              {submitting ? "Updating..." : "Save Status"}
            </Button>
          </form>
        )}
      </Card>

      <Card>
        <h2>Order Items</h2>
        <Table
          rows={items}
          getRowKey={(item, index) => item.id || `${item.product_id}-${index}`}
          columns={[
            {
              key: "product",
              title: "Product",
              render: (item) => item.name || `Product #${item.product_id}`,
            },
            { key: "unit_price", title: "Unit Price", render: (item) => formatCurrency(item.unit_price) },
            { key: "quantity", title: "Quantity" },
            { key: "subtotal", title: "Subtotal", render: (item) => formatCurrency(item.subtotal) },
          ]}
        />
      </Card>
    </div>
  );
}
