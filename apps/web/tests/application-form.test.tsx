import { render, screen } from "@testing-library/react";

import { ApplicationForm } from "@/core/components/artisan/application-form";

describe("artisan application form", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads localized category labels through the web catalogue proxy", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        categories: [{ id: "category-uuid", name: { ar: "الفخار", en: "Pottery" } }],
      }),
    } as Response);

    render(<ApplicationForm locale="ar" />);

    const category = await screen.findByRole("checkbox", { name: "الفخار" });
    expect(category).toHaveAttribute("value", "category-uuid");
    expect(globalThis.fetch).toHaveBeenCalledWith("/api/catalogue?locale=ar");
  });
});
