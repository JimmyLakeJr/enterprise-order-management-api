import Button from "./Button";
import Modal from "./Modal";

export default function ConfirmDialog({
  open,
  title = "Xác nhận thao tác",
  message,
  confirmLabel = "Xác nhận",
  cancelLabel = "Hủy",
  danger = false,
  loading = false,
  onConfirm,
  onCancel,
}) {
  const actions = (
    <>
      <Button variant="secondary" onClick={onCancel} disabled={loading}>{cancelLabel}</Button>
      <Button variant={danger ? "danger" : "primary"} onClick={onConfirm} disabled={loading}>
        {loading ? "Đang xử lý..." : confirmLabel}
      </Button>
    </>
  );

  return (
    <Modal open={open} title={title} onClose={loading ? undefined : onCancel} actions={actions}>
      <p>{message}</p>
    </Modal>
  );
}
