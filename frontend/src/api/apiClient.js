import axios from "axios";
import { AUTH_EVENTS, STORAGE_KEYS } from "../constants/domain";

export const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1").replace(/\/$/, "");
export const API_ORIGIN = new URL(API_BASE_URL).origin;

const AUTH_ENDPOINTS_WITHOUT_REFRESH = [
  "/auth/register",
  "/auth/login",
  "/auth/refresh-token",
  "/auth/logout",
];

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: { "Content-Type": "application/json" },
});

const refreshClient = axios.create({
  baseURL: API_BASE_URL,
  headers: { "Content-Type": "application/json" },
});

let refreshPromise = null;

export function getAccessToken() {
  return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
}

export function getRefreshToken() {
  return localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
}

export function saveAuthTokens({ accessToken, refreshToken }) {
  if (accessToken) localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken);
  if (refreshToken) localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken);
}

export function clearAuthStorage() {
  localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
  localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
  localStorage.removeItem(STORAGE_KEYS.USER);
}

function forceLogout() {
  clearAuthStorage();
  window.dispatchEvent(new Event(AUTH_EVENTS.LOGOUT));

  if (window.location.pathname !== "/login") {
    window.location.assign("/login");
  }
}

function isRefreshExcludedRequest(config) {
  if (config?.skipAuthRefresh) return true;
  const url = config?.url || "";
  return AUTH_ENDPOINTS_WITHOUT_REFRESH.some((endpoint) => url.startsWith(endpoint));
}

async function performRefresh() {
  const currentRefreshToken = getRefreshToken();
  if (!currentRefreshToken) throw new Error("Missing refresh token");

  const response = await refreshClient.post("/auth/refresh-token", {
    refresh_token: currentRefreshToken,
  });
  const data = getData(response);

  if (!data?.access_token || !data?.refresh_token) {
    throw new Error("Invalid refresh token response");
  }

  saveAuthTokens({
    accessToken: data.access_token,
    refreshToken: data.refresh_token,
  });

  if (data.user) {
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(data.user));
  }

  window.dispatchEvent(new CustomEvent(AUTH_EVENTS.REFRESHED, { detail: data }));
  return data.access_token;
}

function refreshAccessTokenOnce() {
  if (!refreshPromise) {
    refreshPromise = performRefresh()
      .catch((error) => {
        forceLogout();
        throw error;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }

  return refreshPromise;
}

apiClient.interceptors.request.use((config) => {
  const accessToken = getAccessToken();
  if (accessToken && !config.skipAuth) {
    config.headers.Authorization = `Bearer ${accessToken}`;
    config._accessToken = accessToken;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    const shouldRefresh =
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retry &&
      !isRefreshExcludedRequest(originalRequest);

    if (!shouldRefresh) return Promise.reject(error);

    originalRequest._retry = true;
    const latestAccessToken = getAccessToken();
    const newAccessToken =
      latestAccessToken && originalRequest._accessToken !== latestAccessToken
        ? latestAccessToken
        : await refreshAccessTokenOnce();
    originalRequest.headers ||= {};
    originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
    return apiClient(originalRequest);
  }
);

export function getData(response) {
  return response.data?.data;
}

export function getMessage(error) {
  return error?.response?.data?.message || error?.message || "Request failed";
}
