export default function Textarea({ label, error, className = "", fieldClassName = "", ...props }) {
  return (
    <label className={`field ${fieldClassName}`.trim()}>
      {label && <span>{label}</span>}
      <textarea className={className} aria-invalid={Boolean(error)} {...props} />
      {error && <small className="field-error">{error}</small>}
    </label>
  );
}
