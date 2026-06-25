import { apiClient, getData } from "./apiClient";

export const cartApi = {
  quote: async (payload) => getData(await apiClient.post("/cart/quote", payload, { skipAuthRefresh: true })),
};
