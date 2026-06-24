const STATUS_CONFIG = Object.freeze({
  pending: { label: "Chờ xác nhận", tone: "warning" },
  confirmed: { label: "Đã xác nhận", tone: "primary" },
  shipping: { label: "Đang giao", tone: "info" },
  completed: { label: "Hoàn thành", tone: "success" },
  cancelled: { label: "Đã hủy", tone: "danger" },
});

export function getOrderStatus(status) {
  return STATUS_CONFIG[status] || { label: status || "Không xác định", tone: "default" };
}

export function getOrderStatusLabel(status) {
  return getOrderStatus(status).label;
}

export function getOrderStatusTone(status) {
  return getOrderStatus(status).tone;
}
