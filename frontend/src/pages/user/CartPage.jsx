import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getMessage } from "../../api/apiClient";
import { createOrder } from "../../api/orderApi";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Table from "../../components/common/Table";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";
import { formatCurrency } from "../../utils/formatCurrency";

export default function CartPage() {
  const { items, getCartTotal, updateQuantity, removeFromCart, clearCart } = useCart();
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  async function handleCreateOrder() {
    if (!isAuthenticated) {
      navigate("/login", { state: { from: "/cart" } });
      return;
    }

    if (!items.length || submitting) return;

    setSubmitting(true);
    setError("");

    try {
      const order = await createOrder({
        items: items.map((item) => ({
          product_id: item.product.id,
          quantity: item.quantity,
        })),
      });

      clearCart();
      navigate(order?.id ? `/orders/${order.id}` : "/my-orders", { replace: true });
    } catch (err) {
      setError(`${getMessage(err)}. Tồn kho hoặc giá có thể đã thay đổi, vui lòng kiểm tra và thử lại.`);
    } finally {
      setSubmitting(false);
    }
  }

  if (items.length === 0) {
    return (
      <Card className="cart-empty-card">
        <EmptyState
          title="Giỏ hàng đang trống"
          description="Hãy thêm sản phẩm trước khi tạo đơn hàng."
        />
        <Link className="btn btn-primary" to="/products">Khám phá sản phẩm</Link>
      </Card>
    );
  }

  const frontendEstimate = getCartTotal();

  return (
    <Card className="cart-card">
      <div className="page-header">
        <div>
          <span className="eyebrow">Giỏ hàng</span>
          <h1>Kiểm tra sản phẩm</h1>
          <p className="muted">Kiểm tra số lượng và giá tạm tính trước khi tạo đơn hàng.</p>
        </div>
        <div className="cart-header-total">
          <span>Tạm tính</span>
          <strong className="cart-total">{formatCurrency(frontendEstimate)}</strong>
        </div>
      </div>

      <ErrorMessage message={error} />

      <Table
        caption={`${items.length} dòng sản phẩm trong giỏ`}
        rows={items}
        getRowKey={(item) => item.product.id}
        columns={[
          {
            key: "product",
            title: "Sản phẩm",
            render: (item) => (
              <div className="cart-product">
                <Link to={`/products/${item.product.id}`}>{item.product.name}</Link>
                {item.product.stock !== null && item.product.stock !== undefined && (
                  <small className="muted">Tồn kho đã biết: {item.product.stock}</small>
                )}
              </div>
            ),
          },
          { key: "price", title: "Đơn giá tạm tính", render: (item) => formatCurrency(item.product.price) },
          {
            key: "quantity",
            title: "Số lượng",
            render: (item) => (
              <input
                className="quantity-input"
                aria-label={`Số lượng ${item.product.name}`}
                type="number"
                inputMode="numeric"
                min="1"
                max={item.product.stock ?? undefined}
                step="1"
                value={item.quantity}
                onChange={(event) => updateQuantity(item.product.id, event.target.value)}
              />
            ),
          },
          {
            key: "subtotal",
            title: "Thành tiền tạm tính",
            render: (item) => <strong>{formatCurrency(Number(item.product.price || 0) * item.quantity)}</strong>,
          },
          {
            key: "action",
            title: "Thao tác",
            render: (item) => (
              <Button type="button" variant="danger" onClick={() => removeFromCart(item.product.id)}>
                Xóa
              </Button>
            ),
          },
        ]}
      />

      <div className="cart-summary">
        <div>
          <span className="muted">Tổng tiền tạm tính</span>
          <strong className="cart-total">{formatCurrency(frontendEstimate)}</strong>
          <small className="muted">Giá trị chính thức được xác nhận khi tạo đơn.</small>
        </div>
        <div className="actions cart-actions">
          <Link className="btn btn-secondary" to="/products">Tiếp tục mua</Link>
          <Button type="button" variant="secondary" onClick={clearCart} disabled={submitting}>Xóa giỏ hàng</Button>
          <Button type="button" onClick={handleCreateOrder} disabled={submitting}>
            {submitting ? "Đang tạo đơn..." : isAuthenticated ? "Tạo đơn hàng" : "Đăng nhập để tạo đơn"}
          </Button>
        </div>
      </div>
    </Card>
  );
}
