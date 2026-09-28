import { render, screen, waitFor } from "@testing-library/react";

import { InventoryOperations } from "@/features/admin/components/inventory-operations";

describe("inventory operations", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads balances without exposing internal product ids", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        items: [{
          productId: "internal-product-id",
          productCode: "AISHA-ABC1234567",
          productName: "Copper bowl",
          workshopId: "internal-workshop-id",
          workshopName: "Amina Atelier",
          artisanName: "Amina",
          onHand: 5,
          available: 5,
          reserved: 0,
          quarantined: 0,
          damaged: 0,
          rejected: 0,
          shipped: 0,
          updatedAt: "2026-09-27T12:00:00Z",
        }],
        page: 1,
        total: 1,
      }),
    }) as unknown as typeof fetch;

    render(<InventoryOperations />);

    await waitFor(() => expect(screen.getByText("Copper bowl")).toBeInTheDocument());
    expect(screen.getByText("AISHA-ABC1234567")).toBeInTheDocument();
    expect(screen.queryByText("internal-product-id")).not.toBeInTheDocument();
    expect(screen.getByRole("combobox", { name: "Workshop" })).toHaveValue("All workshops");
  });
});
