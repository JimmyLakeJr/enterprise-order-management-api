import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { cartApi } from "../../api/cartApi";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import Table from "../../components/common/Table";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";
import { formatCurrency } from "../../utils/formatCurrency";
import { buildOrderDraftFromCart, saveOrderDraft } from "../../utils/orderDraft";

export default function CartPage() {
  const { items, getCartTotal, updateQuantity, removeFromCart, clearCart } = useCart();
  const { isAuthenticated, user } = useAuth();
  const navigate = useNavigate();
  const [quote, setQuote] = useState(null);
  const [quoteLoading, setQuoteLoading] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function loadQuote() {
      if (!items.length) {
        setQuote(null);
        return;
      }

      setQuoteLoading(true);
      try {
        const result = await cartApi.quote({
          items: items.map((item) => ({
            product_id: item.product.id,
            quantity: item.quantity,
          })),
        });
        if (!cancelled) setQuote(result);
      } catch {
        if (!cancelled) setQuote(null);
      } finally {
        if (!cancelled) setQuoteLoading(false);
      }
    }

    void loadQuote();
    return () => {
      cancelled = true;
    };
  }, [items]);

  function handleCreateOrder() {
    if (!isAuthenticated) {
      navigate("/login", { state: { from: "/cart" } });
      return;
    }

    if (!items.length) return;

    saveOrderDraft(buildOrderDraftFromCart(items, user), user?.id || null);
    navigate("/orders/new");
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
  const backendEstimate = quote?.final_amount ?? frontendEstimate;

  return (
    <Card className="cart-card">
      <div className="page-header">
        <div>
          <span className="eyebrow">Giỏ hàng</span>
          <h1>Kiểm tra sản phẩm</h1>
          <p className="muted">Kiểm tra số lượng và giá tạm tính trước khi sang luồng tạo đơn hàng 3 bước.</p>
        </div>
        <div className="cart-header-total">
          <span>Tạm tính</span>
          <strong className="cart-total">{formatCurrency(backendEstimate)}</strong>
        </div>
      </div>

      {Array.isArray(quote?.warnings) && quote.warnings.length > 0 && (
        <div className="api-gap-note warning-note">
          {quote.warnings.map((warning) => <div key={warning}>{warning}</div>)}
        </div>
      )}

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
          <span className="muted">Tổng tiền ước tính</span>
          <strong className="cart-total">{formatCurrency(backendEstimate)}</strong>
          <small className="muted">
            {quoteLoading ? "Đang tính lại giá từ backend..." : "Giá trị chính thức được backend xác nhận ở bước thanh toán."}
          </small>
        </div>
        <div className="actions cart-actions">
          <Link className="btn btn-secondary" to="/products">Tiếp tục mua</Link>
          <Button type="button" variant="secondary" onClick={clearCart}>Xóa giỏ hàng</Button>
          <Button type="button" onClick={handleCreateOrder}>
            {isAuthenticated ? "Tạo đơn hàng" : "Đăng nhập để tạo đơn"}
          </Button>
        </div>
      </div>
    </Card>
  );
}
