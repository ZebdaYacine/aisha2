import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { WarehouseOperations } from "@/features/admin/components/warehouse-operations";

function jsonResponse(body: unknown) {
  return { ok: true, status: 200, json: async () => body } as Response;
}

describe("warehouse operations", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads receptions and opens inspection details", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue(
      jsonResponse({
        items: [
          {
            id: "reception-1",
            productId: "product-1",
            productName: "Copper bowl",
            artisanName: "Artisan",
            workshopName: "Workshop",
            receivedQuantity: 5,
            referenceKey: "parcel-1",
            status: "RECEIVED_PENDING_INSPECTION",
            evidence: [],
          },
        ],
        page: 1,
        pageSize: 10,
        total: 1,
      }),
    ) as unknown as typeof fetch;

    render(<WarehouseOperations />);
    await waitFor(() => expect(screen.getByText("Copper bowl")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Details" }));
    expect(screen.getByRole("dialog")).toBeInTheDocument();
    expect(screen.getByText("Inspect batch")).toBeInTheDocument();
    expect(screen.getByText(/must total 5/)).toBeInTheDocument();
  });
});
