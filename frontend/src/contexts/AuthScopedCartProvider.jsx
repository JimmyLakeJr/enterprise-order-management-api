import { useAuth } from "./AuthContext";
import { CartProvider } from "./CartContext";

export default function AuthScopedCartProvider({ children }) {
  const { user } = useAuth();
  const userID = user?.id || null;

  return (
    <CartProvider key={userID || "guest"} userID={userID}>
      {children}
    </CartProvider>
  );
}
