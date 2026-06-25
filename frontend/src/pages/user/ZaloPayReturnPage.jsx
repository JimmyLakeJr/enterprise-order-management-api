import { useEffect, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { getMessage } from "../../api/apiClient";
import { paymentApi } from "../../api/paymentApi";
import Card from "../../components/common/Card";
import ErrorMessage from "../../components/common/ErrorMessage";
import { formatCurrency } from "../../utils/formatCurrency";

function getStatusTitle(status) {
  switch (status) {
    case "paid":
      return "Thanh toan thanh cong";
    case "failed":
      return "Thanh toan that bai";
    case "expired":
      return "Thanh toan het han";
    default:
      return "Thanh toan dang cho xu ly";
  }
}

function getStatusMessage(status) {
  switch (status) {
    case "paid":
      return "He thong da xac nhan giao dich tu ZaloPay.";
    case "failed":
      return "Giao dich khong thanh cong. Ban co the quay lai buoc thanh toan de thu lai.";
    case "expired":
      return "Phien thanh toan da het han. Vui long tao lai giao dich moi.";
    default:
      return "Backend dang kiem tra ket qua giao dich voi ZaloPay.";
  }
}

export default function ZaloPayReturnPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const transactionId = searchParams.get("transactionId");
  const orderId = searchParams.get("orderId");
  const [status, setStatus] = useState("pending");
  const [payment, setPayment] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function checkStatus() {
      if (!transactionId) {
        setError("Khong tim thay ma giao dich thanh toan. Vui long quay lai trang don hang.");
        return;
      }

      try {
        const result = await paymentApi.getZaloPayStatus(transactionId);
        if (cancelled) return;
        setPayment(result);
        setStatus(result?.status || "pending");
      } catch (err) {
        if (cancelled) return;
        setError(getMessage(err));
      }
    }

    void checkStatus();
    return () => {
      cancelled = true;
    };
  }, [transactionId]);

  return (
    <div className="grid order-flow-layout">
      <Card className="orders-card">
        <div className="page-header">
          <div>
            <span className="eyebrow">ZaloPay</span>
            <h1>{getStatusTitle(status)}</h1>
            <p className="muted">{getStatusMessage(status)}</p>
          </div>
          <div className={`order-step-chip payment-result-chip payment-result-${status}`}>{status}</div>
        </div>

        <ErrorMessage message={error} />

        {payment ? (
          <div className="gap-feature-panel payment-result-panel">
            <div className="profile-readonly-grid order-customer-grid">
              <div>
                <span>Ma giao dich</span>
                <strong>{payment.transactionId}</strong>
              </div>
              <div>
                <span>Don hang</span>
                <strong>#{payment.orderId}</strong>
              </div>
              <div>
                <span>So tien</span>
                <strong>{formatCurrency(payment.amount)}</strong>
              </div>
              <div>
                <span>Trang thai</span>
                <strong>{payment.status}</strong>
              </div>
            </div>
          </div>
        ) : !error ? (
          <div className="preview-only-banner">Dang kiem tra trang thai thanh toan voi backend...</div>
        ) : null}

        <div className="actions order-flow-actions">
          {(status === "failed" || status === "expired") && orderId ? (
            <button className="btn btn-secondary" type="button" onClick={() => navigate(`/orders/${orderId}`, { replace: true })}>
              Về đơn hàng
            </button>
          ) : null}
          {payment?.orderId ? (
            <Link className="btn btn-primary" to={`/orders/${payment.orderId}`}>
              Xem chi tiet don hang
            </Link>
          ) : (
            <Link className="btn btn-primary" to="/orders">
              Ve danh sach don hang
            </Link>
          )}
        </div>
      </Card>
    </div>
  );
}
