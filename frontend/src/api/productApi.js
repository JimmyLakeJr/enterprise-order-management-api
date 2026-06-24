import { apiClient, getData } from "./apiClient";

export async function getProducts(params = {}) {
  const response = await apiClient.get("/products", {
    params: removeEmptyParams(params),
  });

  return {
    data: Array.isArray(response.data?.data) ? response.data.data : [],
    meta: response.data?.meta || null,
  };
}

export async function getProductById(id) {
  const response = await apiClient.get(`/products/${id}`);
  return response.data?.data;
}

export const productApi = {
  list: getProducts,
  listAdmin: async (params = {}) => {
    const response = await apiClient.get("/admin/products", { params: removeEmptyParams(params) });
    return {
      data: Array.isArray(response.data?.data) ? response.data.data : [],
      meta: response.data?.meta || null,
    };
  },
  detail: getProductById,
  create: async (payload) => getData(await apiClient.post("/products", payload)),
  uploadImage: async (file) => {
    const formData = new FormData();
    formData.append("image", file);
    return getData(await apiClient.post("/products/upload-image", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    }));
  },
  update: async (id, payload) => getData(await apiClient.put(`/products/${id}`, payload)),
  remove: async (id) => apiClient.delete(`/products/${id}`),
  restore: async (id) => getData(await apiClient.put(`/admin/products/${id}/restore`)),
};

function removeEmptyParams(params) {
  return Object.fromEntries(
    Object.entries(params).filter(([, value]) => value !== "" && value !== null && value !== undefined)
  );
}
