import { apiClient, getData } from "./apiClient";

export const authApi = {
  register: async (payload) => getData(await apiClient.post("/auth/register", payload, { skipAuth: true, skipAuthRefresh: true })),
  login: async (payload) => getData(await apiClient.post("/auth/login", payload, { skipAuth: true, skipAuthRefresh: true })),
  refreshToken: async (refreshToken) => getData(await apiClient.post(
    "/auth/refresh-token",
    { refresh_token: refreshToken },
    { skipAuth: true, skipAuthRefresh: true }
  )),
  logout: async (refreshToken) => apiClient.post(
    "/auth/logout",
    { refresh_token: refreshToken },
    { skipAuthRefresh: true }
  ),
  getMe: async () => getData(await apiClient.get("/auth/me")),
};
