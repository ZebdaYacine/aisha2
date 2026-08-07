import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function GET(request: NextRequest) {
  return proxyResponse(await authenticatedBackend(`/artisan/products${request.nextUrl.search}`));
}

export async function POST(request: NextRequest) {
  return proxyResponse(await authenticatedBackend("/artisan/products", { method: "POST", body: JSON.stringify(await request.json()) }));
}
