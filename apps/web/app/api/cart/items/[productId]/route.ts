import { NextRequest } from "next/server";
import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";
async function path(id: string) { return `/cart/items/${encodeURIComponent(id)}`; }
export async function PATCH(request: NextRequest, context: { params: Promise<{ productId: string }> }) { const { productId } = await context.params; return proxyResponse(await authenticatedBackend(await path(productId), { method: "PATCH", headers: { "content-type": "application/json" }, body: JSON.stringify(await request.json()) })); }
export async function DELETE(_: NextRequest, context: { params: Promise<{ productId: string }> }) { const { productId } = await context.params; return proxyResponse(await authenticatedBackend(await path(productId), { method: "DELETE" })); }
