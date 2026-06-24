import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { authApi } from "../api/authApi";
import { userApi } from "../api/userApi";
import { clearAuthStorage, getAccessToken, getRefreshToken, saveAuthTokens } from "../api/apiClient";
import { AUTH_EVENTS, ROLES, STORAGE_KEYS } from "../constants/domain";

const AuthContext = createContext(null);

function readStoredUser() {
  try {
    const raw = localStorage.getItem(STORAGE_KEYS.USER);
    return raw ? JSON.parse(raw) : null;
  } catch {
    localStorage.removeItem(STORAGE_KEYS.USER);
    return null;
  }
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(readStoredUser);
  const [accessToken, setAccessToken] = useState(getAccessToken);
  const [refreshToken, setRefreshToken] = useState(getRefreshToken);
  const [loading, setLoading] = useState(true);

  const clearAuthData = useCallback(() => {
    clearAuthStorage();
    setUser(null);
    setAccessToken(null);
    setRefreshToken(null);
  }, []);

  const applyAuthData = useCallback((data) => {
    const nextAccessToken = data?.access_token || getAccessToken();
    const nextRefreshToken = data?.refresh_token || getRefreshToken();
    const nextUser = data?.user || null;

    saveAuthTokens({
      accessToken: nextAccessToken,
      refreshToken: nextRefreshToken,
    });

    if (nextUser) {
      localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(nextUser));
      setUser(nextUser);
    }

    setAccessToken(nextAccessToken);
    setRefreshToken(nextRefreshToken);
  }, []);

  const loadMe = useCallback(async () => {
    setLoading(true);

    try {
      if (!getAccessToken()) {
        const storedRefreshToken = getRefreshToken();
        if (!storedRefreshToken) {
          clearAuthData();
          return null;
        }

        const refreshed = await authApi.refreshToken(storedRefreshToken);
        applyAuthData(refreshed);
      }

      const me = await authApi.getMe();
      localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(me));
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
  }, [applyAuthData, clearAuthData]);

  useEffect(() => {
    const loadTimer = window.setTimeout(() => void loadMe(), 0);

    function handleForcedLogout() {
      setUser(null);
      setAccessToken(null);
      setRefreshToken(null);
    }

    function handleTokenRefresh(event) {
      applyAuthData(event.detail);
    }

    window.addEventListener(AUTH_EVENTS.LOGOUT, handleForcedLogout);
    window.addEventListener(AUTH_EVENTS.REFRESHED, handleTokenRefresh);
    return () => {
      window.clearTimeout(loadTimer);
      window.removeEventListener(AUTH_EVENTS.LOGOUT, handleForcedLogout);
      window.removeEventListener(AUTH_EVENTS.REFRESHED, handleTokenRefresh);
    };
  }, [applyAuthData, loadMe]);

  const login = useCallback(async (payload) => {
    setLoading(true);
    try {
      const data = await authApi.login(payload);
      applyAuthData(data);
      return data;
    } finally {
      setLoading(false);
    }
  }, [applyAuthData]);

  const register = useCallback(async (payload) => {
    setLoading(true);
    try {
      const data = await authApi.register(payload);
      applyAuthData(data);
      return data;
    } finally {
      setLoading(false);
    }
  }, [applyAuthData]);

  const logout = useCallback(async () => {
    const currentRefreshToken = getRefreshToken();
    try {
      if (currentRefreshToken) await authApi.logout(currentRefreshToken);
    } catch {
      // Local logout must still succeed when the access token is already invalid.
    } finally {
      clearAuthData();
    }
  }, [clearAuthData]);

  const updateProfile = useCallback(async (payload) => {
    const updatedUser = await userApi.updateMe(payload);
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(updatedUser));
    setUser(updatedUser);
    return updatedUser;
  }, []);

  const uploadAvatar = useCallback(async (file) => {
    const updatedUser = await userApi.uploadAvatar(file);
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(updatedUser));
    setUser(updatedUser);
    return updatedUser;
  }, []);

  const uploadProfileVideo = useCallback(async (file) => {
    const updatedUser = await userApi.uploadProfileVideo(file);
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(updatedUser));
    setUser(updatedUser);
    return updatedUser;
  }, []);

  const completeOAuthLogin = useCallback(async ({ accessToken: nextAccessToken, refreshToken: nextRefreshToken }) => {
    saveAuthTokens({
      accessToken: nextAccessToken,
      refreshToken: nextRefreshToken,
    });
    setAccessToken(nextAccessToken);
    setRefreshToken(nextRefreshToken);
    const me = await loadMe();
    return me;
  }, [loadMe]);

  const value = useMemo(
    () => ({
      user,
      accessToken,
      refreshToken,
      loading,
      isAuthenticated: Boolean(accessToken && user),
      isAdmin: user?.role === ROLES.ADMIN,
      register,
      login,
      logout,
      loadMe,
      updateProfile,
      uploadAvatar,
      uploadProfileVideo,
      completeOAuthLogin,
    }),
    [accessToken, completeOAuthLogin, loadMe, loading, login, logout, refreshToken, register, updateProfile, uploadAvatar, uploadProfileVideo, user]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  return useContext(AuthContext);
}
