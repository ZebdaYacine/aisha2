import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST() {
  return proxyResponse(await authenticatedBackend("/notifications/read-all", { method: "POST" }));
}
