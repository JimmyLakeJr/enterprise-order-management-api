import Button from "./Button";

export default function ErrorMessage({ message, onRetry, retryLabel = "Thử lại" }) {
  if (!message) return null;
  return (
    <div className="alert alert-danger error-state" role="alert">
      <span className="state-icon" aria-hidden="true">!</span>
      <div>
        <strong>Không thể hoàn tất</strong>
        <p>{message}</p>
      </div>
      {onRetry && <Button type="button" variant="secondary" onClick={onRetry}>{retryLabel}</Button>}
    </div>
  );
}
