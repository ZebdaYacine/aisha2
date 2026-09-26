import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST(
  request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  const { id } = await params;
  return proxyResponse(
    await authenticatedBackend(`/warehouse/receptions/${id}/inspect`, {
      method: "POST",
      body: await request.text(),
    }),
  );
}
