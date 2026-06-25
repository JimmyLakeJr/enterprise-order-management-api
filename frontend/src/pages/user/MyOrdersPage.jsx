import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { orderApi } from "../../api/orderApi";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Loading from "../../components/common/Loading";
import Select from "../../components/common/Select";
import Table from "../../components/common/Table";
import { ORDER_STATUSES } from "../../constants/domain";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";
import { formatCurrency } from "../../utils/formatCurrency";
import { buildOrderDraftFromCart, saveOrderDraft } from "../../utils/orderDraft";
import { getOrderStatus, getOrderStatusLabel } from "../../utils/orderStatus";

const INITIAL_QUERY = {
  page: 1,
  limit: 10,
  status: "",
};

export default function MyOrdersPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const { user } = useAuth();
  const { items } = useCart();
  const [orders, setOrders] = useState([]);
  const [meta, setMeta] = useState(null);
  const [query, setQuery] = useState(INITIAL_QUERY);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [successMessage, setSuccessMessage] = useState(location.state?.successMessage || "");

  async function loadOrders(params = query) {
    setLoading(true);
    setError("");
    try {
      const result = await orderApi.myOrders(params);
      setOrders(result.data);
      setMeta(result.meta);
    } catch (err) {
      setOrders([]);
      setMeta(null);
      setError(err?.response?.data?.message || err?.message || "Không tải được đơn hàng.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    const timer = window.setTimeout(() => void loadOrders(query), 0);
    return () => window.clearTimeout(timer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query]);

  function handleStartOrder() {
    saveOrderDraft(buildOrderDraftFromCart(items, user), user?.id || null);
    navigate("/orders/new");
  }

  const currentPage = Number(meta?.page || query.page);
  const totalPages = Number(meta?.total_pages || 0);
  const totalOrders = Number(meta?.total || orders.length);

  return (
    <Card className="orders-card">
      <div className="page-header">
        <div>
          <span className="eyebrow">Lịch sử mua hàng</span>
          <h1>Đơn hàng của tôi</h1>
          <p className="muted">{loading ? "Đang tải..." : `${totalOrders} đơn hàng phù hợp với bộ lọc hiện tại.`}</p>
        </div>
        <div className="actions">
          <Select
            label="Lọc trạng thái"
            value={query.status}
            onChange={(event) => setQuery((current) => ({ ...current, page: 1, status: event.target.value }))}
          >
            <option value="">Tất cả trạng thái</option>
            {ORDER_STATUSES.map((status) => <option key={status} value={status}>{getOrderStatusLabel(status)}</option>)}
          </Select>
          <Button type="button" onClick={handleStartOrder}>Tạo đơn hàng</Button>
          <Link className="btn btn-secondary" to="/products">Tiếp tục mua sắm</Link>
        </div>
      </div>

      {successMessage && (
        <div className="alert alert-success order-flow-alert">
          <span>{successMessage}</span>
          <button type="button" className="inline-dismiss" onClick={() => setSuccessMessage("")}>Đóng</button>
        </div>
      )}

      {loading ? (
        <Loading label="Đang tải đơn hàng..." variant="table" count={5} />
      ) : error ? (
        <div>
          <ErrorMessage message={`Không tải được đơn hàng: ${error}`} />
          <Button onClick={() => loadOrders(query)}>Thử lại</Button>
        </div>
      ) : !orders?.length ? (
        <div className="order-empty-stack">
          <EmptyState title="Chưa có đơn hàng" description="Đơn hàng bạn tạo sẽ xuất hiện tại đây." />
          <div className="actions">
            <Button type="button" onClick={handleStartOrder}>Tạo đơn hàng</Button>
            <Link className="btn btn-secondary" to="/products">Chọn sản phẩm</Link>
          </div>
        </div>
      ) : (
        <Table
          caption={`${orders.length} đơn hàng`}
          rows={orders}
          columns={[
            {
              key: "id",
              title: "Mã đơn",
              render: (order) => <Link className="order-link" to={`/orders/${order.id}`}>#{order.id}</Link>,
            },
            {
              key: "items",
              title: "Sản phẩm",
              render: (order) => `${order.items?.length || 0} dòng sản phẩm`,
            },
            {
              key: "total_amount",
              title: "Tổng tiền",
              render: (order) => <strong>{formatCurrency(order.total_amount)}</strong>,
            },
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
              render: (order) => (
                <Link className="btn btn-secondary" to={`/orders/${order.id}`}>Xem chi tiết</Link>
              ),
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
