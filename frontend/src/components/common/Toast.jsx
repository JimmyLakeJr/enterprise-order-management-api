import { useEffect } from "react";

export default function Toast({ message, tone = "info", duration = 3500, onDismiss }) {
  useEffect(() => {
    if (!message || !onDismiss || duration <= 0) return undefined;
    const timer = window.setTimeout(onDismiss, duration);
    return () => window.clearTimeout(timer);
  }, [duration, message, onDismiss]);

  if (!message) return null;

  return (
    <div className={`toast toast-${tone}`} role={tone === "danger" ? "alert" : "status"} aria-live="polite">
      {message}
    </div>
  );
}
