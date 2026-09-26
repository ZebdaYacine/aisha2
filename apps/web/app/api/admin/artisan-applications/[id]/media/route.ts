import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function GET(_request: Request, context: { params: Promise<{ id: string }> }) {
  const { id } = await context.params;
  return proxyResponse(await authenticatedBackend(`/admin/artisan-applications/${encodeURIComponent(id)}/media`));
}
