export default function EmptyState({
  title = "Chưa có dữ liệu",
  description = "Dữ liệu sẽ hiển thị tại đây.",
}) {
  return (
    <div className="empty" role="status">
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
  );
}
