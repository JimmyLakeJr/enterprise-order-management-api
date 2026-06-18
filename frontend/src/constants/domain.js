export const ROLES = Object.freeze({
  USER: "user",
  ADMIN: "admin",
});

export const ORDER_STATUS = Object.freeze({
  PENDING: "pending",
  CONFIRMED: "confirmed",
  SHIPPING: "shipping",
  COMPLETED: "completed",
  CANCELLED: "cancelled",
});

export const ORDER_STATUSES = Object.freeze(Object.values(ORDER_STATUS));

export const ORDER_STATUS_TRANSITIONS = Object.freeze({
  [ORDER_STATUS.PENDING]: [ORDER_STATUS.CONFIRMED, ORDER_STATUS.CANCELLED],
  [ORDER_STATUS.CONFIRMED]: [ORDER_STATUS.SHIPPING, ORDER_STATUS.CANCELLED],
  [ORDER_STATUS.SHIPPING]: [ORDER_STATUS.COMPLETED],
  [ORDER_STATUS.COMPLETED]: [],
  [ORDER_STATUS.CANCELLED]: [],
});

const ORDER_STATUS_TONES = Object.freeze({
  [ORDER_STATUS.PENDING]: "warning",
  [ORDER_STATUS.CONFIRMED]: "primary",
  [ORDER_STATUS.SHIPPING]: "info",
  [ORDER_STATUS.COMPLETED]: "success",
  [ORDER_STATUS.CANCELLED]: "danger",
});

export function getOrderStatusTone(status) {
  return ORDER_STATUS_TONES[status] || "default";
}

export const STORAGE_KEYS = Object.freeze({
  ACCESS_TOKEN: "access_token",
  REFRESH_TOKEN: "refresh_token",
  USER: "user",
  CART: "cart",
});

export const AUTH_EVENTS = Object.freeze({
  LOGOUT: "auth:logout",
  REFRESHED: "auth:refreshed",
});
