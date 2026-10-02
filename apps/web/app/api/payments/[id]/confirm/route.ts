import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST(request: Request, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return proxyResponse(await authenticatedBackend(`/payments/${id}/confirm`, {
    method: "POST",
    headers: { "Idempotency-Key": request.headers.get("Idempotency-Key") ?? "" },
  }));
}
