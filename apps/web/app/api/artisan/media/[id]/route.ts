import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

type RouteContext = { params: Promise<{ id: string }> };

export async function PATCH(request: NextRequest, context: RouteContext) {
  const { id } = await context.params;
  return proxyResponse(
    await authenticatedBackend(`/artisan/profile/media/${encodeURIComponent(id)}`, {
      method: "PATCH",
      body: await request.formData(),
    }),
  );
}

export async function DELETE(_request: NextRequest, context: RouteContext) {
  const { id } = await context.params;
  return proxyResponse(
    await authenticatedBackend(`/artisan/profile/media/${encodeURIComponent(id)}`, { method: "DELETE" }),
  );
}
