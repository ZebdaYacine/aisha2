import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";
export async function GET() { return proxyResponse(await authenticatedBackend("/cart")); }
