import { apiClient, getData } from "./apiClient";

export const authApi = {
  register: async (payload) => getData(await apiClient.post("/auth/register", payload)),
  login: async (payload) => getData(await apiClient.post("/auth/login", payload)),
  refreshToken: async (refreshToken) => getData(await apiClient.post("/auth/refresh-token", { refresh_token: refreshToken })),
  logout: async (refreshToken) => apiClient.post("/auth/logout", { refresh_token: refreshToken }),
  getMe: async () => getData(await apiClient.get("/auth/me")),
};
