import { createContext, useContext, useEffect, useState } from "react";
import { STORAGE_KEYS } from "../constants/domain";

const CartContext = createContext(null);

function normalizeQuantity(quantity) {
  const parsed = Number(quantity);
  if (!Number.isFinite(parsed)) return null;

  const normalized = Math.floor(parsed);
  return normalized > 0 ? normalized : null;
}

function getKnownStock(stock) {
  if (stock === null || stock === undefined || stock === "") return null;

  const parsed = Number(stock);
  if (!Number.isFinite(parsed) || parsed < 0) return null;
  return Math.floor(parsed);
}

function clampToStock(quantity, stock) {
  const normalized = normalizeQuantity(quantity);
  if (!normalized) return null;

  const knownStock = getKnownStock(stock);
  if (knownStock === null) return normalized;
  return Math.min(normalized, knownStock) || null;
}

function sanitizeStoredCart(value) {
  if (!Array.isArray(value)) return [];

  return value.flatMap((item) => {
    if (!item?.product?.id) return [];
    const quantity = clampToStock(item.quantity, item.product.stock);
    return quantity ? [{ product: item.product, quantity }] : [];
  });
}

function readCartFromStorage() {
  try {
    const raw = localStorage.getItem(STORAGE_KEYS.CART);
    return raw ? sanitizeStoredCart(JSON.parse(raw)) : [];
  } catch {
    localStorage.removeItem(STORAGE_KEYS.CART);
    return [];
  }
}

export function CartProvider({ children }) {
  const [items, setItems] = useState(readCartFromStorage);

  useEffect(() => {
    localStorage.setItem(STORAGE_KEYS.CART, JSON.stringify(items));
  }, [items]);

  function addToCart(product, quantity = 1) {
    if (!product?.id || !normalizeQuantity(quantity)) return false;

    setItems((currentItems) => {
      const existing = currentItems.find((item) => item.product.id === product.id);
      const requestedQuantity = (existing?.quantity || 0) + normalizeQuantity(quantity);
      const nextQuantity = clampToStock(requestedQuantity, product.stock);
      if (!nextQuantity) return currentItems;

      if (existing) {
        return currentItems.map((item) =>
          item.product.id === product.id ? { product, quantity: nextQuantity } : item
        );
      }

      return [...currentItems, { product, quantity: nextQuantity }];
    });

    return true;
  }

  function updateQuantity(productId, quantity) {
    const normalized = normalizeQuantity(quantity);
    if (!normalized) return false;

    setItems((currentItems) =>
      currentItems.map((item) => {
        if (item.product.id !== productId) return item;

        const nextQuantity = clampToStock(normalized, item.product.stock);
        if (!nextQuantity) return item;
        return { ...item, quantity: nextQuantity };
      })
    );
    return true;
  }

  function removeFromCart(productId) {
    setItems((currentItems) => currentItems.filter((item) => item.product.id !== productId));
  }

  function clearCart() {
    setItems([]);
  }

  function getCartTotal() {
    return items.reduce((sum, item) => {
      const price = Number(item.product.price);
      return sum + (Number.isFinite(price) ? price : 0) * item.quantity;
    }, 0);
  }

  function getCartCount() {
    return items.reduce((sum, item) => sum + item.quantity, 0);
  }

  const totalAmount = getCartTotal();
  const totalItems = getCartCount();

  return (
    <CartContext.Provider
      value={{
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
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useCart() {
  return useContext(CartContext);
}
