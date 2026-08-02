import { NextRequest } from "next/server";
import { authenticatedBackend, proxyResponse } from "@/lib/server/backend";
export async function PATCH(r: NextRequest, c: { params: Promise<{ id: string }> }) { const { id } = await c.params; return proxyResponse(await authenticatedBackend(`/addresses/${encodeURIComponent(id)}`, { method: "PATCH", body: JSON.stringify(await r.json()) })); }
export async function DELETE(_r: NextRequest, c: { params: Promise<{ id: string }> }) { const { id } = await c.params; return proxyResponse(await authenticatedBackend(`/addresses/${encodeURIComponent(id)}`, { method: "DELETE" })); }
