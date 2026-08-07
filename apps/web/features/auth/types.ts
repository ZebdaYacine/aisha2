export type AuthUser = {
  id: string;
  email: string;
  displayName: string;
  status?: string;
  roles?: string[];
};

type AuthResponse = {
  user?: unknown;
};

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function stringValue(value: unknown) {
  return typeof value === "string" ? value : "";
}

export function normalizeAuthUser(value: unknown): AuthUser | null {
  if (!isRecord(value) || !stringValue(value.id)) return null;

  const email = stringValue(value.email);
  const displayName = stringValue(value.displayName) || email || "AISHA customer";
  const roles = Array.isArray(value.roles)
    ? value.roles.filter((role): role is string => typeof role === "string")
    : undefined;

  return {
    id: stringValue(value.id),
    email,
    displayName,
    status: stringValue(value.status) || undefined,
    roles,
  };
}

export function userFromAuthResponse(value: unknown) {
  if (!isRecord(value)) return null;
  return normalizeAuthUser((value as AuthResponse).user ?? value);
}

export function landingPathForUser(user: AuthUser, locale: string) {
  if (user.roles?.includes("administrator")) return `/${locale}/admin/artisan-applications`;
  if (user.roles?.includes("artisan")) return `/${locale}/artisan`;
  return `/${locale}/account`;
}

export function userInitials(user: Pick<AuthUser, "displayName" | "email">) {
  const source = user.displayName.trim() || user.email.trim() || "A";
  const words = source.split(/\s+/).filter(Boolean);
  return words
    .slice(0, 2)
    .map((word) => word[0])
    .join("")
    .toUpperCase();
}
