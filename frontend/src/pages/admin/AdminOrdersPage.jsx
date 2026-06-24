import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getMessage } from "../../api/apiClient";
import { orderApi } from "../../api/orderApi";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Loading from "../../components/common/Loading";
import Select from "../../components/common/Select";
import Table from "../../components/common/Table";
import { ORDER_STATUSES, ORDER_STATUS_TRANSITIONS } from "../../constants/domain";
import { formatDate } from "../../utils/format";
import { formatCurrency } from "../../utils/formatCurrency";
import { getOrderStatus, getOrderStatusLabel } from "../../utils/orderStatus";
import { useConfirm } from "../../hooks/useConfirm";

const INITIAL_QUERY = {
  page: 1,
  limit: 10,
  status: "",
};

export default function AdminOrdersPage() {
  const { confirm } = useConfirm();
  const [orders, setOrders] = useState([]);
  const [meta, setMeta] = useState(null);
  const [query, setQuery] = useState(INITIAL_QUERY);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState("");
  const [updatingId, setUpdatingId] = useState(null);
  const [actionError, setActionError] = useState("");

  async function loadOrders(params = query) {
    setLoading(true);
    setLoadError("");
    try {
      const result = await orderApi.list(params);
      setOrders(result.data);
      setMeta(result.meta);
    } catch (err) {
      setOrders([]);
      setMeta(null);
      setLoadError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    const timer = window.setTimeout(() => void loadOrders(query), 0);
    return () => window.clearTimeout(timer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query]);

  async function updateStatus(order, nextStatus) {
    const allowed = ORDER_STATUS_TRANSITIONS[order.status] || [];
    if (!allowed.includes(nextStatus)) return;
    if (nextStatus === "cancelled") {
      const accepted = await confirm({
        title: `Hủy đơn hàng #${order.id}?`,
        message: "Đơn hàng sẽ chuyển sang trạng thái cuối và tồn kho không tự động hoàn lại.",
        confirmLabel: "Hủy đơn hàng",
        danger: true,
      });
      if (!accepted) return;
    }

    setUpdatingId(order.id);
    setActionError("");
    try {
      await orderApi.updateStatus(order.id, nextStatus);
      await loadOrders(query);
    } catch (err) {
      setActionError(getMessage(err));
    } finally {
      setUpdatingId(null);
    }
  }

  const currentPage = Number(meta?.page || query.page);
  const totalPages = Number(meta?.total_pages || 0);
  const totalOrders = Number(meta?.total || orders.length);

  return (
    <Card className="admin-list-card">
      <div className="page-header admin-list-header">
        <div>
          <span className="eyebrow">Vận hành đơn hàng</span>
          <h1>Quản lý đơn hàng</h1>
          <p className="muted">{loading ? "Đang tải..." : `${totalOrders} đơn hàng phù hợp bộ lọc hiện tại.`}</p>
        </div>
        <Select
          label="Lọc trạng thái"
          value={query.status}
          onChange={(event) => setQuery((current) => ({ ...current, page: 1, status: event.target.value }))}
        >
          <option value="">Tất cả trạng thái</option>
          {ORDER_STATUSES.map((status) => <option key={status} value={status}>{getOrderStatusLabel(status)}</option>)}
        </Select>
      </div>

      <div className="status-flow" aria-label="Luồng trạng thái đơn hàng">
        <span>Chờ xác nhận</span><b>→</b><span>Đã xác nhận</span><b>→</b><span>Đang giao</span><b>→</b><span>Hoàn thành</span>
        <small>Đơn đang chờ hoặc đã xác nhận cũng có thể chuyển sang trạng thái đã hủy.</small>
      </div>

      <ErrorMessage message={loadError ? `Không tải được đơn hàng: ${loadError}` : actionError} />

      {loading ? (
        <Loading label="Đang tải tất cả đơn hàng..." variant="table" count={6} />
      ) : orders.length === 0 ? (
        <EmptyState title="Không có đơn hàng" description="Không có order phù hợp với trạng thái đã chọn." />
      ) : (
        <Table
          rows={orders}
          columns={[
            { key: "id", title: "Mã đơn", render: (order) => <Link className="order-link" to={`/admin/orders/${order.id}`}>#{order.id}</Link> },
            {
              key: "user",
              title: "Người dùng",
              render: (order) => (
                <div className="user-summary">
                  <strong>{order.user?.name || `Mã #${order.user_id}`}</strong>
                  {order.user?.email && <small>{order.user.email}</small>}
                </div>
              ),
            },
            { key: "items", title: "Sản phẩm", render: (order) => `${order.items?.length || 0} dòng` },
            { key: "created_at", title: "Ngày tạo", render: (order) => formatDate(order.created_at) || "—" },
            { key: "total_amount", title: "Tổng tiền", render: (order) => <strong>{formatCurrency(order.total_amount)}</strong> },
            {
              key: "status",
              title: "Trạng thái",
              render: (order) => {
                const status = getOrderStatus(order.status);
                return <Badge tone={status.tone}>{status.label}</Badge>;
              },
            },
            {
              key: "action",
              title: "Thao tác",
              render: (order) => {
                const options = ORDER_STATUS_TRANSITIONS[order.status] || [];
                return (
                  <div className="actions table-actions order-actions">
                    <Link className="btn btn-secondary" to={`/admin/orders/${order.id}`}>Chi tiết</Link>
                    {options.length ? (
                      <Select
                        aria-label={`Cập nhật trạng thái đơn hàng ${order.id}`}
                        value=""
                        disabled={updatingId === order.id}
                        onChange={(event) => updateStatus(order, event.target.value)}
                      >
                        <option value="">{updatingId === order.id ? "Đang cập nhật..." : "Chuyển trạng thái"}</option>
                        {options.map((status) => <option key={status} value={status}>{getOrderStatusLabel(status)}</option>)}
                      </Select>
                    ) : (
                      <span className="terminal-status">Trạng thái cuối</span>
                    )}
                  </div>
                );
              },
            },
          ]}
        />
      )}

      {totalPages > 1 && (
        <div className="pagination">
          <Button type="button" variant="secondary" disabled={currentPage <= 1 || loading} onClick={() => setQuery((current) => ({ ...current, page: currentPage - 1 }))}>
            Trang trước
          </Button>
          <span>Trang {currentPage} / {totalPages}</span>
          <Button type="button" variant="secondary" disabled={currentPage >= totalPages || loading} onClick={() => setQuery((current) => ({ ...current, page: currentPage + 1 }))}>
            Trang sau
          </Button>
        </div>
      )}
    </Card>
  );
}
