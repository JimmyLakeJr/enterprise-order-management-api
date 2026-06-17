import { apiClient } from "./apiClient";

export async function getProducts(params = {}) {
  const response = await apiClient.get("/products", {
    params: removeEmptyParams(params),
  });

  return {
    data: response.data?.data || [],
    meta: response.data?.meta || null,
  };
}

export async function getProductById(id) {
  const response = await apiClient.get(`/products/${id}`);
  return response.data?.data;
}

function removeEmptyParams(params) {
  return Object.fromEntries(
    Object.entries(params).filter(([, value]) => value !== "" && value !== null && value !== undefined)
  );
}
