export default function Select({ label, children, error, className = "", fieldClassName = "", ...props }) {
  return (
    <label className={`field ${fieldClassName}`.trim()}>
      {label && <span>{label}</span>}
      <select className={className} aria-invalid={Boolean(error)} {...props}>{children}</select>
      {error && <small className="field-error">{error}</small>}
    </label>
  );
}
