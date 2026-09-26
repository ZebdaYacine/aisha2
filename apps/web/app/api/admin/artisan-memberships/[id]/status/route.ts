import { NextRequest } from "next/server";
import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST(request: NextRequest, context: { params: Promise<{ id: string }> }) {
  const { id } = await context.params;
  return proxyResponse(await authenticatedBackend(`/admin/artisan-memberships/${encodeURIComponent(id)}/status`, { method: "POST", body: JSON.stringify(await request.json()) }));
}
