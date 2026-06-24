export default function Loading({ label = "Đang tải dữ liệu...", variant = "default", count = 4 }) {
  if (variant !== "default") {
    return (
      <div className={`skeleton-layout skeleton-${variant}`} role="status" aria-live="polite" aria-label={label}>
        {Array.from({ length: count }, (_, index) => (
          <div className="skeleton-item" key={index} aria-hidden="true">
            <span className="skeleton-media" />
            <span className="skeleton-line skeleton-line-wide" />
            <span className="skeleton-line" />
            <span className="skeleton-line skeleton-line-short" />
          </div>
        ))}
        <span className="sr-only">{label}</span>
      </div>
    );
  }

  return (
    <div className="state" role="status" aria-live="polite">
      <span className="loading-spinner" aria-hidden="true" />
      <span>{label}</span>
    </div>
  );
}
