export default function Loading({ label = "Đang tải dữ liệu..." }) {
  return (
    <div className="state" role="status" aria-live="polite">
      <span className="loading-spinner" aria-hidden="true" />
      <span>{label}</span>
    </div>
  );
}
