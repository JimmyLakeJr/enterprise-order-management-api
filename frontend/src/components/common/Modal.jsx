import { useEffect, useId } from "react";

export default function Modal({ open, title, children, actions, onClose, closeLabel = "Đóng" }) {
  const titleId = useId();

  useEffect(() => {
    if (!open) return undefined;

    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    function handleKeyDown(event) {
      if (event.key === "Escape") onClose?.();
    }

    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.body.style.overflow = previousOverflow;
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [open, onClose]);

  if (!open) return null;

  function handleOverlayClick(event) {
    if (event.target === event.currentTarget) onClose?.();
  }

  return (
    <div className="modal-overlay" role="presentation" onMouseDown={handleOverlayClick}>
      <section className="modal-panel" role="dialog" aria-modal="true" aria-labelledby={title ? titleId : undefined}>
        <header className="modal-header">
          {title && <h2 id={titleId}>{title}</h2>}
          <button type="button" className="modal-close" onClick={onClose} aria-label={closeLabel}>×</button>
        </header>
        <div className="modal-body">{children}</div>
        {actions && <footer className="modal-actions">{actions}</footer>}
      </section>
    </div>
  );
}
