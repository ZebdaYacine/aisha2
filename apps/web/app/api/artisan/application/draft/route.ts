import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST(request: NextRequest) {
  return proxyResponse(await authenticatedBackend("/artisan-applications/draft", { method: "POST", body: JSON.stringify(await request.json()) }));
}
