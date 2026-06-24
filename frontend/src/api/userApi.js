import { apiClient, getData } from "./apiClient";

export const userApi = {
  updateMe: async (payload) => getData(await apiClient.put("/users/me", payload)),
  uploadAvatar: async (file) => {
    const formData = new FormData();
    formData.append("avatar", file);
    return getData(await apiClient.post("/users/me/avatar", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    }));
  },
  uploadProfileVideo: async (file) => {
    const formData = new FormData();
    formData.append("video", file);
    return getData(await apiClient.post("/users/me/profile-video", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    }));
  },
  list: async (params) => {
    const response = await apiClient.get("/users", { params });
    return {
      data: Array.isArray(response.data?.data) ? response.data.data : [],
      meta: response.data?.meta || null,
    };
  },
  update: async (id, payload) => getData(await apiClient.put(`/users/${id}`, payload)),
  remove: async (id) => apiClient.delete(`/users/${id}`),
};
