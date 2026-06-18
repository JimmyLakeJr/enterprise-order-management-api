import { apiClient, getData } from "./apiClient";

export const userApi = {
  list: async (params) => {
    const response = await apiClient.get("/users", { params });
    return { data: response.data.data || [], meta: response.data.meta };
  },
  update: async (id, payload) => getData(await apiClient.put(`/users/${id}`, payload)),
  remove: async (id) => apiClient.delete(`/users/${id}`),
};
