import { NextRequest, NextResponse } from "next/server";

import {
  authenticatedBackend,
  backend,
  clearSession,
  proxyResponse,
  refreshToken,
  setSession,
} from "@/core/lib/server/backend";

const allowed = new Set(["login", "register", "forgot-password", "reset-password"]);

export async function GET(
  request: NextRequest,
  context: { params: Promise<{ action: string }> },
) {
  const { action } = await context.params;
  if (action === "me") {
    const response = await authenticatedBackend("/me");
    if (response.status === 401) await clearSession();
    return proxyResponse(response);
  }
  if (action !== "activate") {
    return NextResponse.json({ error: { code: "RESOURCE_NOT_FOUND" } }, { status: 404 });
  }

  const token = request.nextUrl.searchParams.get("token")?.trim();
  if (!token) {
    return NextResponse.json(
      { error: { code: "VALIDATION_FAILED", message: "Activation token is required." } },
      { status: 400 },
    );
  }
  const response = await backend(`/auth/activate?token=${encodeURIComponent(token)}`);
  if (!response.ok) return proxyResponse(response);
  const result = await response.json();
  return NextResponse.json({ activated: true, user: result.user }, { status: response.status });
}

export async function POST(
  request: NextRequest,
  context: { params: Promise<{ action: string }> },
) {
  const { action } = await context.params;
  if (action === "logout") {
    try {
      const token = await refreshToken();
      if (token) {
        await backend("/auth/logout", {
          method: "POST",
          body: JSON.stringify({ refreshToken: token }),
        });
      }
    } finally {
      await clearSession();
    }
    return new NextResponse(null, { status: 204 });
  }
  if (!allowed.has(action)) {
    return NextResponse.json({ error: { code: "RESOURCE_NOT_FOUND" } }, { status: 404 });
  }

  const input = await request.json();
  const payload = action === "register"
    ? { email: input.email, password: input.password, displayName: input.fullName }
    : input;
  const response = await backend(`/auth/${action}`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
  if (!response.ok) return proxyResponse(response);

  if (action === "login") {
    const result = await response.json();
    await setSession(result.tokens);
    return NextResponse.json({ user: result.user }, { status: response.status });
  }
  if (action === "register") {
    const result = await response.json();
    // Registration deliberately does not create a browser session. The user
    // must prove ownership of the email address through the activation link.
    return NextResponse.json(
      { user: result.user, activationRequired: true },
      { status: response.status },
    );
  }
  return new NextResponse(null, { status: response.status });
}
