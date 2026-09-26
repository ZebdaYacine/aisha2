import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST() {
  return proxyResponse(await authenticatedBackend("/artisan-applications/me/submit", { method: "POST" }));
}
