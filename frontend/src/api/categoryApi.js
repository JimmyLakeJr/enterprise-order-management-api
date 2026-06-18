import { apiClient, getData } from "./apiClient";

export async function getCategories() {
  const response = await apiClient.get("/categories");
  return response.data?.data || [];
}

export const categoryApi = {
  list: getCategories,
  create: async (payload) => getData(await apiClient.post("/categories", payload)),
  update: async (id, payload) => getData(await apiClient.put(`/categories/${id}`, payload)),
  remove: async (id) => apiClient.delete(`/categories/${id}`),
};
