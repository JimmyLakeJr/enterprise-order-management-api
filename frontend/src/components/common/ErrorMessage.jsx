export default function ErrorMessage({ message }) {
  if (!message) return null;
  return <div className="alert alert-danger">{message}</div>;
}
