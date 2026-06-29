import type { Stats, Url, User } from "../types";

const BASE = "/api/v1";

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!res.ok) {
    const body = await res
      .json()
      .catch(() => ({ error: res.statusText, message: res.statusText }));
    throw body;
  }

  if (res.status === 204) return undefined as T;

  return res.json();
}

export const api = {
  // ── Users ──
  getUsers: (limit?: number, offset?: number) => {
    const q = new URLSearchParams();
    if (limit) q.set("limit", String(limit));
    if (offset) q.set("offset", String(offset));
    const s = q.toString();
    return request<User[]>(`/users${s ? "?" + s : ""}`);
  },

  createUser: (data: { nickname: string; email: string }) =>
    request<User>("/users", { method: "POST", body: JSON.stringify(data) }),

  getUser: (id: number) => request<User>(`/users/${id}`),

  patchUser: (
    id: number,
    data: { nickname?: string | null; email?: string | null },
  ) =>
    request<User>(`/users/${id}`, {
      method: "PATCH",
      body: JSON.stringify(data),
    }),

  deleteUser: (id: number) =>
    request<void>(`/users/${id}`, { method: "DELETE" }),

  // ── URLs ──
  getURLs: (userId?: number, limit?: number, offset?: number) => {
    const q = new URLSearchParams();
    if (userId) q.set("user_id", String(userId));
    if (limit) q.set("limit", String(limit));
    if (offset) q.set("offset", String(offset));
    const s = q.toString();
    return request<Url[]>(`/urls${s ? "?" + s : ""}`);
  },

  createURL: (data: { url: string; user_id: number }) =>
    request<Url>("/urls", { method: "POST", body: JSON.stringify(data) }),

  getURL: (id: number) => request<Url>(`/urls/${id}`),

  patchURL: (id: number) => request<Url>(`/urls/${id}`, { method: "PATCH" }),

  deleteURL: (id: number) => request<void>(`/urls/${id}`, { method: "DELETE" }),

  // ── Stats ──
  getStats: (shortUrl: string, from?: string, to?: string) => {
    const q = new URLSearchParams({ short_url: shortUrl });
    if (from) q.set("from", from);
    if (to) q.set("to", to);
    return request<Stats>(`/stats?${q.toString()}`);
  },
};
