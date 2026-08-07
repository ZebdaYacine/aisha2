import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function GET() {
  return proxyResponse(await authenticatedBackend("/artisan/profile/media"));
}

export async function POST(request: NextRequest) {
  return proxyResponse(await authenticatedBackend("/artisan/profile/media", { method: "POST", body: await request.formData() }));
}
