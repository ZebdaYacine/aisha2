import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function GET() {
  return proxyResponse(await authenticatedBackend("/artisan/workshops"));
}

export async function POST(request: NextRequest) {
  return proxyResponse(
    await authenticatedBackend("/artisan/workshops", {
      method: "POST",
      headers: { "Idempotency-Key": request.headers.get("Idempotency-Key") ?? crypto.randomUUID() },
      body: JSON.stringify(await request.json()),
    }),
  );
}
