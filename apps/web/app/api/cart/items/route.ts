import { NextRequest } from "next/server";
import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";
export async function POST(request: NextRequest) { return proxyResponse(await authenticatedBackend("/cart/items", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(await request.json()) })); }
