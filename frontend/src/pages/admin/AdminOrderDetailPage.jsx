import { useState } from "react";
import { Link, useParams } from "react-router-dom";
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
import { ORDER_STATUS_TRANSITIONS } from "../../constants/domain";
import { useAsync } from "../../hooks/useAsync";
import { formatDate } from "../../utils/format";
import { formatCurrency } from "../../utils/formatCurrency";
import { getOrderStatus, getOrderStatusLabel } from "../../utils/orderStatus";
import { useConfirm } from "../../hooks/useConfirm";

export default function AdminOrderDetailPage() {
  const { confirm } = useConfirm();
  const { id } = useParams();
  const { data: order, loading, error: loadError, reload } = useAsync(() => orderApi.detail(id), [id]);
  const [nextStatus, setNextStatus] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [actionError, setActionError] = useState("");

  if (loading) return <Loading label="Đang tải chi tiết đơn hàng..." variant="detail" count={2} />;
  if (loadError && !order) return <ErrorMessage message={`Không tải được đơn hàng: ${loadError}`} />;
  if (!order) return <ErrorMessage message="Không tìm thấy đơn hàng." />;

  const items = Array.isArray(order.items) ? order.items : [];
  const availableStatuses = ORDER_STATUS_TRANSITIONS[order.status] || [];
  const status = getOrderStatus(order.status);

  async function handleUpdateStatus(event) {
    event.preventDefault();
    if (!availableStatuses.includes(nextStatus)) return;
    if (nextStatus === "cancelled") {
      const accepted = await confirm({
        title: `Hủy đơn hàng #${order.id}?`,
        message: "Đơn hàng sẽ chuyển sang trạng thái cuối và tồn kho không tự động hoàn lại.",
        confirmLabel: "Hủy đơn hàng",
        danger: true,
      });
      if (!accepted) return;
    }

    setSubmitting(true);
    setActionError("");
    try {
      await orderApi.updateStatus(order.id, nextStatus);
      setNextStatus("");
      await reload();
    } catch (err) {
      setActionError(getMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="grid admin-order-detail">
      <Link to="/admin/orders" className="back-link">← Quay lại quản lý đơn hàng</Link>

      <Card className="admin-panel-card">
        <div className="page-header order-detail-header">
          <div>
            <span className="eyebrow">Chi tiết đơn hàng</span>
            <h1>Đơn hàng #{order.id}</h1>
            <p className="muted">Thông tin sản phẩm, trạng thái và giá trị đơn hàng.</p>
          </div>
          <Badge tone={status.tone} className="order-status-badge">{status.label}</Badge>
        </div>

        <ErrorMessage message={actionError} />

        <div className="order-info-grid">
          <div><span>Người dùng</span><strong>{order.user?.name || `Mã #${order.user_id}`}</strong>{order.user?.email && <small>{order.user.email}</small>}</div>
          <div><span>Tổng tiền</span><strong>{formatCurrency(order.total_amount)}</strong></div>
          <div><span>Trạng thái</span><strong>{status.label}</strong></div>
          <div><span>Ngày tạo</span><strong>{formatDate(order.created_at) || "—"}</strong></div>
        </div>

        {availableStatuses.length === 0 ? (
          <div className="api-gap-note terminal-note">Đơn hàng đang ở trạng thái cuối và không thể chuyển tiếp.</div>
        ) : (
          <form className="status-form" onSubmit={handleUpdateStatus}>
            <Select label="Trạng thái tiếp theo" value={nextStatus} onChange={(event) => setNextStatus(event.target.value)}>
              <option value="">Chọn theo luồng hợp lệ</option>
              {availableStatuses.map((next) => <option key={next} value={next}>{getOrderStatusLabel(next)}</option>)}
            </Select>
            <Button type="submit" disabled={!nextStatus || submitting}>{submitting ? "Đang cập nhật..." : "Cập nhật trạng thái"}</Button>
          </form>
        )}
      </Card>

      <Card className="admin-panel-card">
        <div className="page-header compact-header">
          <div>
            <h2>Sản phẩm trong đơn</h2>
            <p className="muted">Đơn giá và thành tiền được lưu tại thời điểm đặt hàng.</p>
          </div>
        </div>

        {items.length === 0 ? (
          <EmptyState title="Đơn hàng chưa có sản phẩm" description="Chưa có sản phẩm nào trong đơn hàng này." />
        ) : (
          <Table
            rows={items}
            getRowKey={(item, index) => `${item.product_id}-${index}`}
            columns={[
              { key: "product", title: "Sản phẩm", render: (item) => item.name || `Sản phẩm #${item.product_id}` },
              { key: "product_id", title: "Product ID", render: (item) => `#${item.product_id}` },
              { key: "unit_price", title: "Đơn giá", render: (item) => formatCurrency(item.unit_price) },
              { key: "quantity", title: "Số lượng" },
              { key: "subtotal", title: "Thành tiền", render: (item) => <strong>{formatCurrency(item.subtotal)}</strong> },
            ]}
          />
        )}

        <div className="order-total-row"><span>Tổng tiền đơn hàng</span><strong>{formatCurrency(order.total_amount)}</strong></div>
      </Card>
    </div>
  );
}
