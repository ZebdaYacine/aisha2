import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { WorkshopsPanel } from "@/features/artisan/components/workshops-panel";
import { ProductWorkspace } from "@/features/product/components/product-workspace";

function jsonResponse(body: unknown) {
  return { ok: true, status: 200, json: async () => body } as Response;
}

describe("artisan management tables", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("shows workshop rows with first-column details, update, and delete actions", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue(
      jsonResponse({
        items: [{ id: "workshop-1", name: "Amina Studio", wilaya: "Algiers", location: "Kasbah", description: "Ceramics", status: "ACTIVE", isDefault: false, isPublic: true, productCount: 3 }],
      }),
    ) as unknown as typeof fetch;

    render(<WorkshopsPanel locale="en" profile={{ membershipStatus: "ACTIVE", status: "APPROVED" }} />);

    await waitFor(() => expect(screen.getByText("Amina Studio")).toBeInTheDocument());
    expect(screen.getByRole("button", { name: "Add workshop" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Details" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Update" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Delete" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Details" }));
    expect(screen.getByRole("dialog", { name: "Details" })).toHaveTextContent("Kasbah");
    fireEvent.keyDown(document, { key: "Escape" });
    expect(screen.queryByRole("dialog", { name: "Details" })).not.toBeInTheDocument();
  });

  it("shows product rows and opens product details without leaving the editor as the only index", async () => {
    globalThis.fetch = jest.fn().mockImplementation(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.includes("/api/artisan/products")) return jsonResponse({ items: [{ id: "product-1", workshopId: "workshop-1", workshopName: "Amina Studio", categoryId: "category-1", productType: "ARTISAN_SPECIFIC", status: "DRAFT", priceMinor: 2500, currency: "DZD", materials: "Clay", productionMethod: "Hand painted", intendedUse: "Decoration", dimensions: "20cm", countryOfOrigin: "Algeria", regionOfOrigin: "Algiers", ecoFriendlyVerified: false, fairTradeVerified: false, madeToOrderEligible: false, translations: [{ locale: "en", name: "Painted bowl", description: "A bowl", story: "", culturalContext: "" }], media: [] }] });
      if (url.includes("/api/artisan/workshops")) return jsonResponse({ items: [{ id: "workshop-1", name: "Amina Studio", status: "ACTIVE", productCount: 1 }] });
      if (url.includes("/categories")) return jsonResponse({ items: [{ id: "category-1", slug: "decoration", name: "Decoration" }], page: 1, pageSize: 100, total: 1 });
      return jsonResponse({ items: [], page: 1, pageSize: 100, total: 0 });
    }) as unknown as typeof fetch;

    render(<ProductWorkspace locale="en" />);

    await waitFor(() => expect(screen.getByText("Painted bowl")).toBeInTheDocument());
    expect(screen.getByRole("button", { name: "New draft" })).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "New draft" }));
    expect(screen.getByRole("dialog", { name: "New draft" })).toBeInTheDocument();
    fireEvent.keyDown(document, { key: "Escape" });
    expect(screen.queryByRole("dialog", { name: "New draft" })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Details" }));
    expect(screen.getByRole("dialog", { name: "Details" })).toHaveTextContent("Hand painted");
  });
});
