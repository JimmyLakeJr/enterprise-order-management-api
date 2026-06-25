import { apiClient, getData } from "./apiClient";

export const paymentApi = {
  createZaloPayPayment: async (payload) => getData(await apiClient.post("/payments/zalopay/create", payload)),
  getZaloPayStatus: async (transactionId) => getData(await apiClient.get(`/payments/zalopay/status/${transactionId}`)),
};
