import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST(request: NextRequest, context: { params: Promise<{ id: string }> }) {
  const { id } = await context.params;
  return proxyResponse(await authenticatedBackend(`/artisan/products/${encodeURIComponent(id)}/media`, { method: "POST", body: await request.formData() }));
}
