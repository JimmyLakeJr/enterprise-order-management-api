import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { getMessage } from "../../api/apiClient";
import { orderApi } from "../../api/orderApi";
import Badge from "../../components/common/Badge";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Loading from "../../components/common/Loading";
import Select from "../../components/common/Select";
import Table from "../../components/common/Table";
import { getOrderStatusTone, ORDER_STATUSES, ORDER_STATUS_TRANSITIONS } from "../../constants/domain";
import { formatCurrency } from "../../utils/format";

export default function AdminOrdersPage() {
  const [orders, setOrders] = useState([]);
  const [statusFilter, setStatusFilter] = useState("all");
  const [loading, setLoading] = useState(true);
  const [updatingId, setUpdatingId] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    loadOrders();
  }, []);

  async function loadOrders() {
    setLoading(true);
    setError("");
    try {
      setOrders(await orderApi.list());
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }

  const filteredOrders = useMemo(() => {
    if (statusFilter === "all") return orders || [];
    return (orders || []).filter((order) => order.status === statusFilter);
  }, [orders, statusFilter]);

  async function updateStatus(order, status) {
    if (!status) return;

    setUpdatingId(order.id);
    setError("");
    try {
      await orderApi.updateStatus(order.id, status);
      await loadOrders();
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setUpdatingId(null);
    }
  }

  if (loading) return <Loading />;

  return (
    <Card>
      <div className="page-header">
        <div>
          <h1>Orders</h1>
          <p className="muted">Xem toàn bộ đơn hàng và cập nhật trạng thái theo đúng luồng nghiệp vụ.</p>
        </div>
        <Select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value)}>
          <option value="all">All status</option>
          {ORDER_STATUSES.map((status) => (
            <option key={status} value={status}>{status}</option>
          ))}
        </Select>
      </div>

      <ErrorMessage message={error} />

      {filteredOrders.length === 0 ? (
        <EmptyState title="Không có đơn hàng" description="Thử đổi status filter hoặc tạo đơn hàng mới từ store." />
      ) : (
        <Table
          rows={filteredOrders}
          columns={[
            {
              key: "id",
              title: "Order ID",
              render: (order) => <Link to={`/admin/orders/${order.id}`}>#{order.id}</Link>,
            },
            {
              key: "user",
              title: "User/Customer",
              render: (order) => `User #${order.user_id}`,
            },
            { key: "total_amount", title: "Total", render: (order) => formatCurrency(order.total_amount) },
            { key: "status", title: "Status", render: (order) => <Badge tone={getOrderStatusTone(order.status)}>{order.status}</Badge> },
            {
              key: "action",
              title: "Action",
              render: (order) => {
                const options = ORDER_STATUS_TRANSITIONS[order.status] || [];
                return (
                  <div className="actions">
                    <Link className="btn btn-secondary" to={`/admin/orders/${order.id}`}>
                      View
                    </Link>
                    {options.length > 0 && (
                      <Select
                        value=""
                        disabled={updatingId === order.id}
                        onChange={(event) => updateStatus(order, event.target.value)}
                      >
                        <option value="">Update status</option>
                        {options.map((status) => (
                          <option key={status} value={status}>
                            {status}
                          </option>
                        ))}
                      </Select>
                    )}
                  </div>
                );
              },
            },
          ]}
        />
      )}
    </Card>
  );
}
