export type AuthUser = {
  id: string;
  email: string;
  displayName: string;
  status?: string;
  roles?: string[];
  customerEnabled: boolean;
  artisanStatus: string;
  artisanEnabled: boolean;
  capabilities: string[];
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
  const displayName =
    stringValue(value.displayName) || email || "AISHA customer";
  const roles = Array.isArray(value.roles)
    ? value.roles.filter((role): role is string => typeof role === "string")
    : undefined;
  const capabilities = Array.isArray(value.capabilities)
    ? value.capabilities.filter(
        (capability): capability is string => typeof capability === "string",
      )
    : [];

  return {
    id: stringValue(value.id),
    email,
    displayName,
    status: stringValue(value.status) || undefined,
    roles,
    customerEnabled: value.customerEnabled === true,
    artisanStatus: stringValue(value.artisanStatus) || "NOT_STARTED",
    artisanEnabled: value.artisanEnabled === true,
    capabilities,
  };
}

export function userFromAuthResponse(value: unknown) {
  if (!isRecord(value)) return null;
  return normalizeAuthUser((value as AuthResponse).user ?? value);
}

export function landingPathForUser(user: AuthUser, locale: string) {
  if (hasCapability(user, "admin.audit.read"))
    return `/${locale}/admin`;
  if (hasCapability(user, "admin.product_moderation.read"))
    return `/${locale}/admin/moderation`;
  if (hasCapability(user, "warehouse.read"))
    return `/${locale}/admin/warehouse`;
  if (user.artisanEnabled && hasCapability(user, "artisan.account.read"))
    return `/${locale}/artisan`;
  return `/${locale}/account`;
}

export function hasCapability(
  user: AuthUser | null | undefined,
  capability: string,
) {
  return user?.capabilities.includes(capability) ?? false;
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
