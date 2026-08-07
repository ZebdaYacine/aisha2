import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function PATCH(request: NextRequest, context: { params: Promise<{ id: string }> }) {
  const { id } = await context.params;
  return proxyResponse(await authenticatedBackend(`/admin/users/${encodeURIComponent(id)}/roles`, { method: "PATCH", body: JSON.stringify(await request.json()) }));
}
