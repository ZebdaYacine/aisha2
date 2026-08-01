import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { GlobalSearch } from "@/components/layout/global-search";
import { storeCopy } from "@/lib/store-copy";

describe("GlobalSearch", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads catalogue data only after the search is opened", async () => {
    const fetchMock = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        categories: [],
        artisans: [],
        products: [{
          slug: "vase",
          name: { en: "Vase", fr: "Vase", ar: "Vase", es: "Vase" },
          summary: { en: "", fr: "", ar: "", es: "" },
          story: { en: "", fr: "", ar: "", es: "" },
          artisanSlug: "artisan",
          categorySlug: "pottery",
          region: { en: "", fr: "", ar: "", es: "" },
          materials: { en: "", fr: "", ar: "", es: "" },
          method: { en: "", fr: "", ar: "", es: "" },
          priceMinor: 1000,
          currency: "DZD",
          availability: "out_of_stock",
          images: [],
          rating: 0,
          reviewCount: 0,
        }],
      }),
    } as Response);
    globalThis.fetch = fetchMock;

    render(<GlobalSearch locale="en" copy={storeCopy("en")} />);
    expect(fetchMock).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole("button", { name: "Search" }));

    await waitFor(() => expect(screen.getByRole("link", { name: "Vase" })).toBeInTheDocument());
    expect(fetchMock).toHaveBeenCalledWith("/api/catalogue?locale=en");
  });
});
