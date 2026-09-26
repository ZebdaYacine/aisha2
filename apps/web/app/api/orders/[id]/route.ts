import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";
export async function GET(_: Request, { params }: { params: Promise<{ id: string }> }) { const { id } = await params; return proxyResponse(await authenticatedBackend(`/orders/${id}`)); }
