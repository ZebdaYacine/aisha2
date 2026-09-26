import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";
export async function GET(request: Request) { const url = new URL(request.url); const query = url.search ? url.search : ""; return proxyResponse(await authenticatedBackend(`/orders${query}`)); }
