import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getMessage } from "../../api/apiClient";
import { createOrder } from "../../api/orderApi";
import { paymentApi } from "../../api/paymentApi";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import ErrorMessage from "../../components/common/ErrorMessage";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";
import { formatCurrency } from "../../utils/formatCurrency";
import { clearOrderDraft, readOrderDraft, saveOrderDraft } from "../../utils/orderDraft";

const ZALOPAY_ENABLED = String(import.meta.env.VITE_ZALOPAY_ENABLED || "false").toLowerCase() === "true";

const PAYMENT_METHODS = [
  { value: "cod", label: "Thanh toán khi nhận hàng (COD)", enabled: true },
  { value: "bank_transfer", label: "Chuyển khoản ngân hàng", enabled: true },
  {
    value: "zalopay",
    label: "ZaloPay",
    enabled: ZALOPAY_ENABLED,
    helper: ZALOPAY_ENABLED
      ? "Thanh toán qua cổng ZaloPay sandbox/production theo cấu hình backend."
      : "ZaloPay đang tắt ở frontend local dev. Bật VITE_ZALOPAY_ENABLED=true sau khi cấu hình backend.",
  },
  { value: "card_wallet", label: "Thẻ / Ví điện tử", enabled: false, helper: "Chưa tích hợp API trong local dev." },
];

export default function OrderPaymentPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const userID = user?.id || null;
  const { clearCart } = useCart();
  const [draft, setDraft] = useState(() => {
    const storedDraft = readOrderDraft(user?.id || null);
    if (!storedDraft) return storedDraft;

    const selectedMethod = PAYMENT_METHODS.find((method) => method.value === storedDraft.paymentMethod);
    if (selectedMethod && !selectedMethod.enabled) {
      return saveOrderDraft(
        {
          ...storedDraft,
          paymentMethod: "cod",
        },
        user?.id || null
      );
    }

    return storedDraft;
  });
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!draft?.items?.length) {
      navigate("/orders/new", {
        replace: true,
        state: { message: "Thông tin đơn hàng đã thiếu hoặc hết hạn. Vui lòng nhập lại trước khi thanh toán." },
      });
    }
  }, [draft, navigate]);

  const totalAmount = useMemo(
    () => (draft?.items || []).reduce((sum, item) => sum + (Number(item.price) || 0) * item.quantity, 0),
    [draft]
  );

  function handlePaymentMethodChange(event) {
    const nextDraft = saveOrderDraft(
      {
        ...draft,
        paymentMethod: event.target.value,
      },
      userID
    );
    setDraft(nextDraft);
  }

  async function handleSubmit() {
    if (!draft?.items?.length || submitting) return;

    setSubmitting(true);
    setError("");

    let createdOrder = null;

    try {
      createdOrder = await createOrder({
        items: draft.items.map((item) => ({
          product_id: item.product_id,
          quantity: item.quantity,
        })),
      });

      if (draft.paymentMethod === "zalopay") {
        const payment = await paymentApi.createZaloPayPayment({
          orderId: createdOrder.id,
          method: "zalopay",
        });

        clearOrderDraft(userID);
        clearCart();

        if (payment?.paymentUrl) {
          window.location.assign(payment.paymentUrl);
          return;
        }

        navigate(`/payment/zalopay/return?transactionId=${payment.transactionId}&orderId=${payment.orderId}`, {
          replace: true,
        });
        return;
      }

      clearOrderDraft(userID);
      clearCart();
      navigate(createdOrder?.id ? `/orders/${createdOrder.id}` : "/orders", {
        replace: true,
        state: {
          successMessage:
            draft.paymentMethod === "bank_transfer"
              ? "Đơn hàng đã được tạo thành công. Thanh toán chuyển khoản đang ở chế độ local dev preview."
              : "Đơn hàng đã được tạo thành công.",
        },
      });
    } catch (err) {
      if (createdOrder?.id && draft.paymentMethod === "zalopay") {
        clearOrderDraft(userID);
        clearCart();
        navigate(`/orders/${createdOrder.id}`, {
          replace: true,
          state: {
            successMessage:
              "Đơn hàng đã được tạo nhưng chưa khởi tạo được phiên thanh toán ZaloPay. Bạn có thể kiểm tra và thanh toán lại sau.",
          },
        });
        return;
      }

      setError(`${getMessage(err)}. Vui lòng kiểm tra tồn kho, cấu hình thanh toán hoặc quay lại bước trước để cập nhật thông tin.`);
    } finally {
      setSubmitting(false);
    }
  }

  if (!draft?.items?.length) return null;

  return (
    <div className="grid order-flow-layout">
      <Card className="orders-card">
        <div className="page-header">
          <div>
            <span className="eyebrow">Bước 3/3</span>
            <h1>Thanh toán và hoàn tất đơn hàng</h1>
            <p className="muted">Xem lại thông tin đơn hàng, chọn phương thức thanh toán và xác nhận tạo đơn.</p>
          </div>
          <div className="order-step-chip">Thanh toán</div>
        </div>

        <ErrorMessage message={error} />

        <div className="checkout-gap-grid order-payment-grid">
          <section className="gap-feature-panel">
            <legend>Tóm tắt đơn hàng</legend>
            <div className="order-summary-list">
              {draft.items.map((item) => (
                <div key={item.product_id} className="order-summary-item">
                  <div>
                    <strong>{item.name}</strong>
                    <small className="muted">
                      {item.quantity} x {formatCurrency(item.price)}
                    </small>
                  </div>
                  <strong>{formatCurrency(item.price * item.quantity)}</strong>
                </div>
              ))}
            </div>
            <div className="order-total-row order-payment-total">
              <span>Tổng tiền</span>
              <strong>{formatCurrency(totalAmount)}</strong>
            </div>
          </section>

          <section className="gap-feature-panel">
            <legend>Thông tin khách hàng</legend>
            <div className="profile-readonly-grid order-customer-grid">
              <div>
                <span>Khách hàng</span>
                <strong>{draft.customerName}</strong>
              </div>
              <div>
                <span>Số điện thoại</span>
                <strong>{draft.phone}</strong>
              </div>
              <div>
                <span>Email</span>
                <strong>{draft.email || "Chưa cung cấp"}</strong>
              </div>
            </div>
            <div className="profile-readonly-grid order-customer-grid single-column">
              <div>
                <span>Địa chỉ giao hàng</span>
                <strong>{draft.shippingAddress}</strong>
              </div>
              <div>
                <span>Ghi chú</span>
                <strong>{draft.note || "Không có ghi chú thêm"}</strong>
              </div>
            </div>
          </section>
        </div>

        <fieldset className="gap-feature-panel order-payment-methods">
          <legend>Phương thức thanh toán</legend>
          {PAYMENT_METHODS.map((method) => (
            <label key={method.value} className={`payment-method-option ${!method.enabled ? "payment-method-disabled" : ""}`}>
              <input
                type="radio"
                name="paymentMethod"
                value={method.value}
                checked={draft.paymentMethod === method.value}
                disabled={!method.enabled || submitting}
                onChange={handlePaymentMethodChange}
              />
              <div>
                <strong>{method.label}</strong>
                {method.helper ? <small>{method.helper}</small> : null}
              </div>
            </label>
          ))}

          {draft.paymentMethod === "bank_transfer" ? (
            <div className="preview-only-banner">
              Chuyển khoản hiện là chế độ giả lập cho local dev. Khi xác nhận, hệ thống vẫn tạo đơn hàng bằng API hiện có.
            </div>
          ) : null}

          {draft.paymentMethod === "zalopay" ? (
            <div className="preview-only-banner">
              Khi xác nhận, frontend sẽ tạo order trước, sau đó gọi backend khởi tạo giao dịch ZaloPay và chuyển bạn sang trang thanh toán của ZaloPay.
            </div>
          ) : null}
        </fieldset>

        <div className="actions order-flow-actions">
          <Link className="btn btn-secondary" to="/orders/new">
            Quay lại thông tin
          </Link>
          <Button type="button" onClick={handleSubmit} disabled={submitting}>
            {submitting
              ? draft.paymentMethod === "zalopay"
                ? "Đang khởi tạo thanh toán..."
                : "Đang hoàn tất đơn hàng..."
              : draft.paymentMethod === "zalopay"
                ? "Thanh toán với ZaloPay"
                : "Hoàn tất đơn hàng"}
          </Button>
        </div>
      </Card>
    </div>
  );
}
