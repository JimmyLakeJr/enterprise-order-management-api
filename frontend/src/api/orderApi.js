import { apiClient, getData } from "./apiClient";

export async function createOrder(payload) {
  const items = Array.isArray(payload) ? payload : payload?.items;
  const orderPayload = {
    items: (items || []).map((item) => ({
      product_id: item.product_id,
      quantity: item.quantity,
    })),
  };

  return getData(await apiClient.post("/orders", orderPayload));
}

export async function getMyOrders() {
  return orderApi.myOrders();
}

export async function getOrderById(id) {
  return getData(await apiClient.get(`/orders/${id}`));
}

export const orderApi = {
  create: createOrder,
  myOrders: async (params = {}) => {
    const response = await apiClient.get("/users/me/orders", { params: removeEmptyParams(params) });
    return {
      data: Array.isArray(response.data?.data) ? response.data.data : [],
      meta: response.data?.meta || null,
    };
  },
  list: async (params = {}) => {
    const response = await apiClient.get("/orders", { params: removeEmptyParams(params) });
    return {
      data: Array.isArray(response.data?.data) ? response.data.data : [],
      meta: response.data?.meta || null,
    };
  },
  detail: getOrderById,
  updateStatus: async (id, status) => getData(await apiClient.put(`/orders/${id}/status`, { status })),
};

function removeEmptyParams(params) {
  return Object.fromEntries(
    Object.entries(params).filter(([, value]) => value !== "" && value !== null && value !== undefined)
  );
}
