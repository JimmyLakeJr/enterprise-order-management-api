import { apiClient, getData } from "./apiClient";

export async function createOrder(payload) {
  return getData(await apiClient.post("/orders", payload));
}

export async function getMyOrders() {
  return getData(await apiClient.get("/users/me/orders"));
}

export async function getOrderById(id) {
  return getData(await apiClient.get(`/orders/${id}`));
}

export const orderApi = {
  list: async () => getData(await apiClient.get("/orders")),
  detail: getOrderById,
  updateStatus: async (id, status) => getData(await apiClient.put(`/orders/${id}/status`, { status })),
};
