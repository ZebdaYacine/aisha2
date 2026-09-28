import { cookies } from "next/headers";
import { NextResponse } from "next/server";
const baseURL = process.env.API_BASE_URL ?? "http://localhost:8088/api/v1";
const secureCookies = process.env.SESSION_COOKIE_SECURE === "true" || (process.env.SESSION_COOKIE_SECURE !== "false" && process.env.NODE_ENV === "production");
const options = (maxAge?: number) => ({ httpOnly: true, secure: secureCookies, sameSite: "lax" as const, path: "/", ...(maxAge ? { maxAge } : {}) });
export async function backend(path: string, init?: RequestInit) {
  const headers = new Headers(init?.headers);
  if (!headers.has("content-type") && !(typeof FormData !== "undefined" && init?.body instanceof FormData)) {
    headers.set("content-type", "application/json");
  }
  return fetch(`${baseURL}${path}`, { ...init, headers, cache: "no-store" });
}
export async function setSession(tokens: { accessToken: string; refreshToken: string; expiresAt: string }) { const jar = await cookies(); jar.set("aisha_access", tokens.accessToken, options(Math.max(60, Math.floor((Date.parse(tokens.expiresAt) - Date.now()) / 1000)))); jar.set("aisha_refresh", tokens.refreshToken, options(2592000)); }
export async function clearSession() { const jar = await cookies(); jar.delete("aisha_access"); jar.delete("aisha_refresh"); }
export async function refreshToken() { return (await cookies()).get("aisha_refresh")?.value; }
export async function authenticatedBackend(path: string, init?: RequestInit) { const jar = await cookies(); const call = (token?: string) => backend(path, { ...init, headers: { ...init?.headers, authorization: `Bearer ${token ?? ""}` } }); const response = await call(jar.get("aisha_access")?.value); if (response.status !== 401) return response; const refresh = jar.get("aisha_refresh")?.value; if (!refresh) return response; const renewed = await backend("/auth/refresh", { method: "POST", body: JSON.stringify({ refreshToken: refresh }) }); if (!renewed.ok) return response; const tokens = await renewed.json(); await setSession(tokens); return call(tokens.accessToken); }
export async function proxyResponse(response: Response) { const body = response.status === 204 ? null : await response.text(); return new NextResponse(body, { status: response.status, headers: body ? { "content-type": response.headers.get("content-type") ?? "application/json" } : undefined }); }
