import { authenticatedBackend, proxyResponse } from "@/core/lib/server/backend";

export async function POST(
  request: Request,
  context: { params: Promise<{ productId: string }> },
) {
  const { productId } = await context.params;
  return proxyResponse(
    await authenticatedBackend(`/warehouse/inventory/${productId}/adjust`, {
      method: "POST",
      body: await request.text(),
    }),
  );
}
