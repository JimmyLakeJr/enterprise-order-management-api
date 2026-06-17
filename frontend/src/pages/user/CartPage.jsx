import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { createOrder } from "../../api/orderApi";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Table from "../../components/common/Table";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";
import { formatCurrency } from "../../utils/format";

export default function CartPage() {
  const { items, getCartTotal, updateQuantity, removeFromCart, clearCart } = useCart();
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  async function handleCheckout() {
    if (!isAuthenticated) {
      navigate("/login", { state: { from: "/cart" } });
      return;
    }

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
      navigate(order?.id ? `/orders/${order.id}` : "/my-orders");
    } catch (err) {
      setError(err?.response?.data?.message || "Không tạo được đơn hàng. Vui lòng kiểm tra tồn kho và thử lại.");
    } finally {
      setSubmitting(false);
    }
  }

  if (items.length === 0) {
    return (
      <EmptyState
        title="Giỏ hàng trống"
        description="Hãy thêm sản phẩm từ danh sách sản phẩm trước khi tạo đơn hàng."
      />
    );
  }

  return (
    <Card>
      <div className="page-header">
        <div>
          <h1>Giỏ hàng</h1>
          <p className="muted">Kiểm tra số lượng trước khi tạo đơn. Giá chính thức sẽ được backend tính lại.</p>
        </div>
        <strong className="cart-total">{formatCurrency(getCartTotal())}</strong>
      </div>

      <ErrorMessage message={error} />

      <Table
        rows={items}
        getRowKey={(item) => item.product.id}
        columns={[
          {
            key: "product",
            title: "Product",
            render: (item) => (
              <div className="cart-product">
                <span>{item.product.name}</span>
                <small className="muted">Tồn kho: {item.product.stock}</small>
              </div>
            ),
          },
          { key: "price", title: "Price", render: (item) => formatCurrency(item.product.price) },
          {
            key: "quantity",
            title: "Quantity",
            render: (item) => (
              <input
                className="quantity-input"
                type="number"
                min="1"
                max={item.product.stock}
                value={item.quantity}
                onChange={(event) => updateQuantity(item.product.id, Number(event.target.value))}
              />
            ),
          },
          { key: "subtotal", title: "Subtotal", render: (item) => formatCurrency(item.product.price * item.quantity) },
          {
            key: "action",
            title: "Action",
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
          <strong>{formatCurrency(getCartTotal())}</strong>
        </div>
        <div className="actions">
          <Link className="btn btn-secondary" to="/">
            Tiếp tục mua
          </Link>
          <Button type="button" variant="secondary" onClick={clearCart}>
            Xóa giỏ hàng
          </Button>
          <Button type="button" onClick={handleCheckout} disabled={submitting}>
            {submitting ? "Đang tạo..." : "Tạo đơn hàng"}
          </Button>
        </div>
      </div>
    </Card>
  );
}
