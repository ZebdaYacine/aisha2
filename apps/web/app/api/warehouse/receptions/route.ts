import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function GET(request: NextRequest) {
  return proxyResponse(
    await authenticatedBackend(`/warehouse/receptions${request.nextUrl.search}`),
  );
}

export async function POST(request: Request) {
  return proxyResponse(
    await authenticatedBackend("/warehouse/receptions", {
      method: "POST",
      body: await request.text(),
    }),
  );
}
