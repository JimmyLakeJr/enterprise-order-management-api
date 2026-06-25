import { STORAGE_KEYS } from "../constants/domain";

function getOrderDraftStorageKey(userID) {
  if (!userID) return null;
  return `${STORAGE_KEYS.ORDER_DRAFT}:${userID}`;
}

function normalizeDraftItem(item) {
  if (!item?.product_id || !item?.name) return null;

  const quantity = Number(item.quantity);
  if (!Number.isFinite(quantity) || quantity <= 0) return null;

  return {
    product_id: item.product_id,
    name: item.name,
    price: Number(item.price) || 0,
    stock: item.stock ?? null,
    image_url: item.image_url || "",
    quantity: Math.floor(quantity),
  };
}

export function buildOrderDraftFromCart(cartItems, user = null) {
  return {
    customerName: user?.full_name || user?.name || "",
    phone: user?.phone || "",
    email: user?.email || "",
    shippingAddress: user?.address || "",
    note: "",
    paymentMethod: "cod",
    items: Array.isArray(cartItems)
      ? cartItems
          .map((item) => {
            if (!item?.product?.id) return null;

            return normalizeDraftItem({
              product_id: item.product.id,
              name: item.product.name,
              price: item.product.price,
              stock: item.product.stock,
              image_url: item.product.image_url,
              quantity: item.quantity,
            });
          })
          .filter(Boolean)
      : [],
  };
}

export function readOrderDraft(userID) {
  const storageKey = getOrderDraftStorageKey(userID);
  if (!storageKey) return null;

  try {
    const raw = localStorage.getItem(storageKey);
    if (!raw) return null;

    const parsed = JSON.parse(raw);
    return {
      customerName: parsed?.customerName || "",
      phone: parsed?.phone || "",
      email: parsed?.email || "",
      shippingAddress: parsed?.shippingAddress || "",
      note: parsed?.note || "",
      paymentMethod: parsed?.paymentMethod || "cod",
      items: Array.isArray(parsed?.items) ? parsed.items.map(normalizeDraftItem).filter(Boolean) : [],
    };
  } catch {
    localStorage.removeItem(storageKey);
    return null;
  }
}

export function saveOrderDraft(draft, userID) {
  const storageKey = getOrderDraftStorageKey(userID);
  if (!storageKey) return null;

  const normalized = {
    customerName: draft?.customerName || "",
    phone: draft?.phone || "",
    email: draft?.email || "",
    shippingAddress: draft?.shippingAddress || "",
    note: draft?.note || "",
    paymentMethod: draft?.paymentMethod || "cod",
    items: Array.isArray(draft?.items) ? draft.items.map(normalizeDraftItem).filter(Boolean) : [],
  };

  localStorage.setItem(storageKey, JSON.stringify(normalized));
  return normalized;
}

export function clearOrderDraft(userID) {
  const storageKey = getOrderDraftStorageKey(userID);
  if (!storageKey) return;
  localStorage.removeItem(storageKey);
}
