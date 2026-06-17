import { apiClient, getData } from "./apiClient";

export const categoryApi = {
  list: async () => getData(await apiClient.get("/categories")),
  detail: async (id) => getData(await apiClient.get(`/categories/${id}`)),
  create: async (payload) => getData(await apiClient.post("/categories", payload)),
  update: async (id, payload) => getData(await apiClient.put(`/categories/${id}`, payload)),
  remove: async (id) => apiClient.delete(`/categories/${id}`),
};

export const productApi = {
  list: async (params) => {
    const response = await apiClient.get("/products", { params });
    return { data: response.data.data || [], meta: response.data.meta };
  },
  detail: async (id) => getData(await apiClient.get(`/products/${id}`)),
  create: async (payload) => getData(await apiClient.post("/products", payload)),
  update: async (id, payload) => getData(await apiClient.put(`/products/${id}`, payload)),
  remove: async (id) => apiClient.delete(`/products/${id}`),
};
