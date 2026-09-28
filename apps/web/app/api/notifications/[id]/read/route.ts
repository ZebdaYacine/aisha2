import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST(_: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return proxyResponse(await authenticatedBackend(`/notifications/${encodeURIComponent(id)}/read`, { method: "POST" }));
}
