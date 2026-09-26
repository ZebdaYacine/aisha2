import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";
export async function POST(request: Request) { const body = await request.text(); const response = await authenticatedBackend("/checkout", { method: "POST", headers: { "Idempotency-Key": request.headers.get("Idempotency-Key") ?? "" }, body }); return proxyResponse(response); }
