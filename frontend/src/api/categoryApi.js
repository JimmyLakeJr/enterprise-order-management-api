import { apiClient } from "./apiClient";

export async function getCategories() {
  const response = await apiClient.get("/categories");
  return response.data?.data || [];
}
