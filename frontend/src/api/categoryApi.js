import { apiClient, getData } from "./apiClient";

export async function getCategories() {
  const response = await apiClient.get("/categories");
  return Array.isArray(response.data?.data) ? response.data.data : [];
}

export const categoryApi = {
  list: getCategories,
  listAdmin: async (status = "all") => {
    const response = await apiClient.get("/admin/categories", { params: { status } });
    return Array.isArray(response.data?.data) ? response.data.data : [];
  },
  create: async (payload) => getData(await apiClient.post("/categories", payload)),
  update: async (id, payload) => getData(await apiClient.put(`/categories/${id}`, payload)),
  remove: async (id) => apiClient.delete(`/categories/${id}`),
  restore: async (id) => getData(await apiClient.put(`/admin/categories/${id}/restore`)),
};
