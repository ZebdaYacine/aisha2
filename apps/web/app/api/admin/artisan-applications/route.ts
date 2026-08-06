import { NextRequest } from "next/server";
import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";
export async function GET(r: NextRequest) {
  return proxyResponse(
    await authenticatedBackend(
      `/admin/artisan-applications${r.nextUrl.search}`,
    ),
  );
}
