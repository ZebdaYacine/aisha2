import { NextRequest } from "next/server";

import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function GET(request: NextRequest) {
  return proxyResponse(await authenticatedBackend(`/notifications${request.nextUrl.search}`));
}
