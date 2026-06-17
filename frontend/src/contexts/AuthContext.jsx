import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { authApi } from "../api/authApi";
import { clearAuthStorage, getAccessToken, getRefreshToken, saveAuthTokens } from "../api/apiClient";

const AuthContext = createContext(null);

function readStoredUser() {
  const raw = localStorage.getItem("user");
  return raw ? JSON.parse(raw) : null;
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(readStoredUser);
  const [accessToken, setAccessToken] = useState(getAccessToken);
  const [refreshToken, setRefreshToken] = useState(getRefreshToken);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadMe();

    function handleForcedLogout() {
      setUser(null);
      setAccessToken(null);
      setRefreshToken(null);
    }

    window.addEventListener("auth:logout", handleForcedLogout);
    return () => window.removeEventListener("auth:logout", handleForcedLogout);
  }, []);

  function applyAuthData(data) {
    const nextAccessToken = data?.access_token || getAccessToken();
    const nextRefreshToken = data?.refresh_token || getRefreshToken();
    const nextUser = data?.user || data;

    saveAuthTokens({
      accessToken: nextAccessToken,
      refreshToken: nextRefreshToken,
    });

    if (nextUser) {
      localStorage.setItem("user", JSON.stringify(nextUser));
      setUser(nextUser);
    }

    setAccessToken(nextAccessToken);
    setRefreshToken(nextRefreshToken);
  }

  function clearAuthData() {
    clearAuthStorage();
    setUser(null);
    setAccessToken(null);
    setRefreshToken(null);
  }

  async function loadMe() {
    const token = getAccessToken();
    if (!token) {
      setLoading(false);
      return null;
    }

    setLoading(true);
    try {
      const me = await authApi.getMe();
      localStorage.setItem("user", JSON.stringify(me));
      setUser(me);
      setAccessToken(getAccessToken());
      setRefreshToken(getRefreshToken());
      return me;
    } catch {
      clearAuthData();
      return null;
    } finally {
      setLoading(false);
    }
  }

  async function login(payload) {
    setLoading(true);
    try {
      const data = await authApi.login(payload);
      applyAuthData(data);
      return data;
    } finally {
      setLoading(false);
    }
  }

  async function register(payload) {
    setLoading(true);
    try {
      const data = await authApi.register(payload);
      applyAuthData(data);
      return data;
    } finally {
      setLoading(false);
    }
  }

  async function logout() {
    const currentRefreshToken = getRefreshToken();
    if (currentRefreshToken) {
      await authApi.logout(currentRefreshToken).catch(() => {});
    }
    clearAuthData();
  }

  const value = useMemo(
    () => ({
      user,
      accessToken,
      refreshToken,
      loading,
      isAuthenticated: Boolean(accessToken && user),
      isAdmin: user?.role === "admin",
      register,
      login,
      logout,
      loadMe,
    }),
    [user, accessToken, refreshToken, loading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  return useContext(AuthContext);
}
