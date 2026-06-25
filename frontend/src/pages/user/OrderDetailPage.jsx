import { Link, useLocation, useParams } from "react-router-dom";
import { getOrderById } from "../../api/orderApi";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Loading from "../../components/common/Loading";
import Table from "../../components/common/Table";
import { useAsync } from "../../hooks/useAsync";
import { formatCurrency } from "../../utils/formatCurrency";
import { getOrderStatus } from "../../utils/orderStatus";

export default function OrderDetailPage() {
  const location = useLocation();
  const { id } = useParams();
  const { data: order, loading, error, reload } = useAsync(() => getOrderById(id), [id]);

  if (loading) return <Loading label="Đang tải chi tiết đơn hàng..." variant="detail" count={2} />;

  if (error) {
    return (
      <div className="detail-error">
        <ErrorMessage message={`Không tải được đơn hàng: ${error}`} />
        <div className="actions">
          <Button onClick={reload}>Thử lại</Button>
          <Link className="btn btn-secondary" to="/orders">Về danh sách</Link>
        </div>
      </div>
    );
  }

  if (!order) return <ErrorMessage message="Không tìm thấy đơn hàng." />;

  const items = Array.isArray(order.items) ? order.items : [];
  const status = getOrderStatus(order.status);

  return (
    <div className="grid order-detail-layout">
      <Link to="/orders" className="back-link">← Quay lại đơn hàng của tôi</Link>

      {location.state?.successMessage && (
        <div className="alert alert-success order-flow-alert">{location.state.successMessage}</div>
      )}

      <Card className="order-overview-card">
        <div className="page-header order-detail-header">
          <div>
            <span className="eyebrow">Chi tiết đơn hàng</span>
            <h1>Đơn hàng #{order.id}</h1>
            <p className="muted">Theo dõi sản phẩm, trạng thái và giá trị đơn hàng.</p>
          </div>
          <Badge tone={status.tone} className="order-status-badge">{status.label}</Badge>
        </div>

        <div className="order-info-grid">
          <div>
            <span>Trạng thái</span>
            <strong>{status.label}</strong>
          </div>
          <div>
            <span>Số dòng sản phẩm</span>
            <strong>{items.length}</strong>
          </div>
          <div>
            <span>Tổng tiền</span>
            <strong>{formatCurrency(order.total_amount)}</strong>
          </div>
        </div>
      </Card>

      <Card className="order-items-card">
        <div className="section-heading">
          <div>
            <h2>Sản phẩm trong đơn</h2>
            <p className="muted">Đơn giá và thành tiền được lưu tại thời điểm tạo đơn.</p>
          </div>
        </div>

        {items.length === 0 ? (
          <EmptyState title="Đơn hàng chưa có sản phẩm" description="Chưa có sản phẩm nào trong đơn hàng này." />
        ) : (
          <Table
            rows={items}
            getRowKey={(item, index) => `${item.product_id}-${index}`}
            columns={[
              {
                key: "product",
                title: "Sản phẩm",
                render: (item) => (
                  <div className="cart-product">
                    <strong>{item.name || `Sản phẩm #${item.product_id}`}</strong>
                    <small className="muted">Mã sản phẩm: {item.product_id}</small>
                  </div>
                ),
              },
              { key: "unit_price", title: "Đơn giá", render: (item) => formatCurrency(item.unit_price) },
              { key: "quantity", title: "Số lượng" },
              { key: "subtotal", title: "Thành tiền", render: (item) => <strong>{formatCurrency(item.subtotal)}</strong> },
            ]}
          />
        )}

        <div className="order-total-row">
          <span>Tổng tiền đơn hàng</span>
          <strong>{formatCurrency(order.total_amount)}</strong>
        </div>
      </Card>
    </div>
  );
}
