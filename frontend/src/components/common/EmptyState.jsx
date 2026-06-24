export default function EmptyState({
  title = "Chưa có dữ liệu",
  description = "Dữ liệu sẽ hiển thị tại đây.",
  action,
}) {
  return (
    <div className="empty" role="status">
      <span className="empty-illustration" aria-hidden="true">
        <i />
        <i />
        <i />
      </span>
      <h3>{title}</h3>
      <p>{description}</p>
      {action && <div className="empty-action">{action}</div>}
    </div>
  );
}
