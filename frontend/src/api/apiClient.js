import axios from "axios";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1";

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

const refreshClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

let refreshPromise = null;

export function getAccessToken() {
  return localStorage.getItem("access_token");
}

export function getRefreshToken() {
  return localStorage.getItem("refresh_token");
}

export function saveAuthTokens({ accessToken, refreshToken }) {
  if (accessToken) localStorage.setItem("access_token", accessToken);
  if (refreshToken) localStorage.setItem("refresh_token", refreshToken);
}

export function clearAuthStorage() {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("user");
}

function forceLogout() {
  clearAuthStorage();
  window.dispatchEvent(new Event("auth:logout"));
  if (window.location.pathname !== "/login") {
    window.location.assign("/login");
  }
}

async function refreshAccessToken() {
  if (!refreshPromise) {
    const refreshToken = getRefreshToken();
    if (!refreshToken) {
      forceLogout();
      throw new Error("Missing refresh token");
    }

    refreshPromise = refreshClient
      .post("/auth/refresh-token", { refresh_token: refreshToken })
      .then((response) => {
        const data = response.data?.data;
        saveAuthTokens({
          accessToken: data?.access_token,
          refreshToken: data?.refresh_token,
        });
        if (data?.user) {
          localStorage.setItem("user", JSON.stringify(data.user));
        }
        return data?.access_token;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }

  return refreshPromise;
}

apiClient.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    const status = error.response?.status;
    const isAuthRequest = originalRequest?.url?.startsWith("/auth/login") || originalRequest?.url?.startsWith("/auth/register") || originalRequest?.url?.startsWith("/auth/refresh-token");

    if (status !== 401 || !originalRequest || originalRequest._retry || isAuthRequest) {
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    try {
      const newAccessToken = await refreshAccessToken();
      originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
      return apiClient(originalRequest);
    } catch (refreshError) {
      forceLogout();
      return Promise.reject(refreshError);
    }
  }
);

export function getData(response) {
  return response.data?.data;
}

export function getMessage(error) {
  return error?.response?.data?.message || error?.message || "Request failed";
}
