import { createContext, useContext, useMemo, useState } from "react";

const CART_STORAGE_KEY = "cart";
const CartContext = createContext(null);

function readCartFromStorage() {
  try {
    const raw = localStorage.getItem(CART_STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

function normalizeQuantity(quantity) {
  const parsed = Number(quantity);
  if (!Number.isFinite(parsed)) return 1;
  return Math.floor(parsed);
}

function clampQuantity(quantity, stock) {
  const nextQuantity = normalizeQuantity(quantity);
  if (nextQuantity <= 0) return 0;
  if (Number.isFinite(Number(stock)) && Number(stock) >= 0) {
    return Math.min(nextQuantity, Number(stock));
  }
  return nextQuantity;
}

export function CartProvider({ children }) {
  const [items, setItems] = useState(readCartFromStorage);

  function saveCart(nextItems) {
    setItems(nextItems);
    localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(nextItems));
  }

  function addToCart(product, quantity = 1) {
    const productId = product?.id;
    if (!productId) return;

    const existing = items.find((item) => item.product.id === productId);
    const currentQuantity = existing?.quantity || 0;
    const nextQuantity = clampQuantity(currentQuantity + normalizeQuantity(quantity), product.stock);

    if (nextQuantity <= 0) return;

    const nextItems = existing
      ? items.map((item) => (item.product.id === productId ? { ...item, product, quantity: nextQuantity } : item))
      : [...items, { product, quantity: nextQuantity }];

    saveCart(nextItems);
  }

  function updateQuantity(productId, quantity) {
    const existing = items.find((item) => item.product.id === productId);
    if (!existing) return;

    const nextQuantity = clampQuantity(quantity, existing.product.stock);
    if (nextQuantity <= 0) {
      removeFromCart(productId);
      return;
    }

    saveCart(items.map((item) => (item.product.id === productId ? { ...item, quantity: nextQuantity } : item)));
  }

  function removeFromCart(productId) {
    saveCart(items.filter((item) => item.product.id !== productId));
  }

  function clearCart() {
    saveCart([]);
  }

  function getCartTotal() {
    return items.reduce((sum, item) => sum + Number(item.product.price || 0) * item.quantity, 0);
  }

  function getCartCount() {
    return items.reduce((sum, item) => sum + item.quantity, 0);
  }

  const totalAmount = getCartTotal();
  const totalItems = getCartCount();

  const value = useMemo(
    () => ({
      items,
      totalAmount,
      totalItems,
      addToCart,
      updateQuantity,
      removeFromCart,
      clearCart,
      getCartTotal,
      getCartCount,
      addItem: addToCart,
      removeItem: removeFromCart,
    }),
    [items, totalAmount, totalItems]
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useCart() {
  return useContext(CartContext);
}
