import { API_ORIGIN } from "../api/apiClient";

export function resolveAssetUrl(value) {
  const url = String(value || "").trim();
  if (!url) return "";
  if (/^https?:\/\//i.test(url)) return url;
  if (url.startsWith("/")) return `${API_ORIGIN}${url}`;
  return `${API_ORIGIN}/${url.replace(/^\/+/, "")}`;
}
