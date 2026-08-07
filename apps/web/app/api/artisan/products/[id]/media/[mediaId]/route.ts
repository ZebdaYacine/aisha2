import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function DELETE(_request: NextRequest, context: { params: Promise<{ id: string; mediaId: string }> }) {
  const { id, mediaId } = await context.params;
  return proxyResponse(await authenticatedBackend(`/artisan/products/${encodeURIComponent(id)}/media/${encodeURIComponent(mediaId)}`, { method: "DELETE" }));
}
