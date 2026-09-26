import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function DELETE(_request: NextRequest, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return proxyResponse(await authenticatedBackend(`/admin/media/products/${encodeURIComponent(id)}`, { method: "DELETE" }));
}
