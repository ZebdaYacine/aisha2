import { normalizeAuthUser, type AuthUser } from "./types";

export async function fetchCurrentUser(signal?: AbortSignal): Promise<AuthUser | null> {
  try {
    const response = await fetch("/api/auth/me", {
      method: "GET",
      cache: "no-store",
      credentials: "same-origin",
      signal,
    });

    if (response.status === 401 || response.status === 403) return null;
    if (!response.ok) throw new Error("Unable to load the current user.");

    return normalizeAuthUser(await response.json());
  } catch (error) {
    if (error instanceof DOMException && error.name === "AbortError") return null;
    throw error;
  }
}

export async function revokeSession() {
  await fetch("/api/auth/logout", {
    method: "POST",
    credentials: "same-origin",
  });
}
