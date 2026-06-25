import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App.jsx";
import { AuthProvider } from "./contexts/AuthContext.jsx";
import AuthScopedCartProvider from "./contexts/AuthScopedCartProvider.jsx";
import ConfirmProvider from "./contexts/ConfirmProvider.jsx";
import "./styles/tokens.css";
import "./styles/global.css";
import "./styles/layout.css";
import "./styles/components.css";
import "./styles/pages.css";
import "./styles/admin.css";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <AuthScopedCartProvider>
          <ConfirmProvider>
            <App />
          </ConfirmProvider>
        </AuthScopedCartProvider>
      </AuthProvider>
    </BrowserRouter>
  </StrictMode>
);
