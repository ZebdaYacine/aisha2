import { NextRequest } from "next/server";
import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";
const paths: Record<string, string> = {
  profile: "/me/profile",
  password: "/me/password",
  addresses: "/addresses",
};
async function forward(request: NextRequest, resource: string) {
  const path = paths[resource];
  if (!path) return new Response(null, { status: 404 });
  const body =
    request.method === "GET" ? undefined : JSON.stringify(await request.json());
  return proxyResponse(
    await authenticatedBackend(path, { method: request.method, body }),
  );
}
export async function GET(
  r: NextRequest,
  c: { params: Promise<{ resource: string }> },
) {
  return forward(r, (await c.params).resource);
}
export async function POST(
  r: NextRequest,
  c: { params: Promise<{ resource: string }> },
) {
  return forward(r, (await c.params).resource);
}
export async function PATCH(
  r: NextRequest,
  c: { params: Promise<{ resource: string }> },
) {
  return forward(r, (await c.params).resource);
}
