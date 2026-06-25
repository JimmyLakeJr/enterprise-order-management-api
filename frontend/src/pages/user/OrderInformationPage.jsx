import { useMemo, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Table from "../../components/common/Table";
import Textarea from "../../components/common/Textarea";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";
import { formatCurrency } from "../../utils/formatCurrency";
import { buildOrderDraftFromCart, readOrderDraft, saveOrderDraft } from "../../utils/orderDraft";

function validateDraft(form, items) {
  const errors = {};

  if (!form.customerName.trim()) errors.customerName = "Vui lòng nhập tên khách hàng.";
  if (!form.phone.trim()) {
    errors.phone = "Vui lòng nhập số điện thoại.";
  } else if (!/^[0-9+\s()-]{8,20}$/.test(form.phone.trim())) {
    errors.phone = "Số điện thoại không hợp lệ.";
  }

  if (form.email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) {
    errors.email = "Email không hợp lệ.";
  }

  if (!form.shippingAddress.trim()) errors.shippingAddress = "Vui lòng nhập địa chỉ giao hàng.";
  if (!items.length) errors.items = "Hãy thêm ít nhất một sản phẩm trước khi thanh toán.";
  if (items.some((item) => !Number.isFinite(Number(item.quantity)) || Number(item.quantity) <= 0)) {
    errors.items = "Số lượng sản phẩm phải lớn hơn 0.";
  }

  return errors;
}

export default function OrderInformationPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();
  const userID = user?.id || null;
  const { items, updateQuantity, removeFromCart } = useCart();
  const [form, setForm] = useState(() => {
    const existingDraft = readOrderDraft(user?.id || null);
    return existingDraft || buildOrderDraftFromCart([], user);
  });
  const [errors, setErrors] = useState({});
  const [banner] = useState(location.state?.message || "");

  const cartDraftItems = useMemo(
    () =>
      items.map((item) => ({
        product_id: item.product.id,
        name: item.product.name,
        price: Number(item.product.price) || 0,
        stock: item.product.stock ?? null,
        image_url: item.product.image_url || "",
        quantity: item.quantity,
      })),
    [items]
  );

  function handleChange(field, value) {
    setForm((current) => ({ ...current, [field]: value }));
    setErrors((current) => ({ ...current, [field]: "" }));
  }

  function handleContinue(event) {
    event.preventDefault();
    const nextErrors = validateDraft(form, cartDraftItems);
    setErrors(nextErrors);

    if (Object.keys(nextErrors).length > 0) return;

    saveOrderDraft(
      {
        ...form,
        items: cartDraftItems,
      },
      userID
    );
    navigate("/orders/new/payment");
  }

  if (!cartDraftItems.length) {
    return (
      <Card className="orders-card">
        <div className="page-header">
          <div>
            <span className="eyebrow">Bước 2/3</span>
            <h1>Thông tin đơn hàng</h1>
            <p className="muted">Bạn cần thêm sản phẩm vào giỏ hàng trước khi tiếp tục tạo đơn.</p>
          </div>
        </div>
        <EmptyState
          title="Chưa có sản phẩm để tạo đơn"
          description="Hãy chọn sản phẩm và thêm vào giỏ hàng, sau đó quay lại bước nhập thông tin."
        />
        <div className="actions">
          <Link className="btn btn-secondary" to="/orders">Quay lại đơn hàng</Link>
          <Link className="btn btn-primary" to="/products">Chọn sản phẩm</Link>
        </div>
      </Card>
    );
  }

  const totalAmount = cartDraftItems.reduce((sum, item) => sum + item.price * item.quantity, 0);

  return (
    <div className="grid order-flow-layout">
      <Card className="orders-card">
        <div className="page-header">
          <div>
            <span className="eyebrow">Bước 2/3</span>
            <h1>Nhập thông tin đơn hàng</h1>
            <p className="muted">Xác nhận khách hàng, địa chỉ giao hàng và sản phẩm trước khi sang bước thanh toán.</p>
          </div>
          <div className="order-step-chip">Thông tin</div>
        </div>

        {banner && <div className="alert alert-success order-flow-alert">{banner}</div>}
        <ErrorMessage message={errors.items} />

        <form className="order-flow-form" onSubmit={handleContinue}>
          <div className="order-flow-grid">
            <Input
              label="Tên khách hàng"
              value={form.customerName}
              error={errors.customerName}
              onChange={(event) => handleChange("customerName", event.target.value)}
            />
            <Input
              label="Số điện thoại"
              value={form.phone}
              error={errors.phone}
              onChange={(event) => handleChange("phone", event.target.value)}
            />
            <Input
              label="Email"
              type="email"
              value={form.email}
              error={errors.email}
              onChange={(event) => handleChange("email", event.target.value)}
            />
            <div className="field">
              <span>Tổng tạm tính</span>
              <div className="order-total-pill">{formatCurrency(totalAmount)}</div>
            </div>
          </div>

          <Textarea
            label="Địa chỉ giao hàng"
            value={form.shippingAddress}
            error={errors.shippingAddress}
            onChange={(event) => handleChange("shippingAddress", event.target.value)}
          />

          <Textarea
            label="Ghi chú đơn hàng"
            value={form.note}
            onChange={(event) => handleChange("note", event.target.value)}
            placeholder="Ví dụ: giao giờ hành chính, gọi trước khi giao..."
          />

          <div className="section-heading">
            <div>
              <h2>Sản phẩm trong đơn</h2>
              <p className="muted">Bạn có thể chỉnh số lượng ngay tại đây trước khi sang bước thanh toán.</p>
            </div>
          </div>

          <Table
            caption={`${cartDraftItems.length} dòng sản phẩm`}
            rows={cartDraftItems}
            getRowKey={(item) => item.product_id}
            columns={[
              {
                key: "name",
                title: "Sản phẩm",
                render: (item) => (
                  <div className="cart-product">
                    <strong>{item.name}</strong>
                    <small className="muted">Mã sản phẩm: {item.product_id}</small>
                  </div>
                ),
              },
              { key: "price", title: "Đơn giá", render: (item) => formatCurrency(item.price) },
              {
                key: "quantity",
                title: "Số lượng",
                render: (item) => (
                  <input
                    className="quantity-input"
                    type="number"
                    min="1"
                    step="1"
                    max={item.stock ?? undefined}
                    value={item.quantity}
                    onChange={(event) => updateQuantity(item.product_id, event.target.value)}
                  />
                ),
              },
              {
                key: "subtotal",
                title: "Thành tiền",
                render: (item) => <strong>{formatCurrency(item.price * item.quantity)}</strong>,
              },
              {
                key: "action",
                title: "Thao tác",
                render: (item) => (
                  <Button type="button" variant="danger" onClick={() => removeFromCart(item.product_id)}>
                    Xóa
                  </Button>
                ),
              },
            ]}
          />

          <div className="actions order-flow-actions">
            <Link className="btn btn-secondary" to="/orders">Quay lại</Link>
            <Button type="submit">Tiếp tục thanh toán</Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
