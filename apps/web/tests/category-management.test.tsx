import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { CategoryManagement } from "@/features/admin/components/category-management";

const jsonResponse = (body: unknown, status = 200) =>
  ({ ok: status >= 200 && status < 300, status, json: async () => body }) as Response;

describe("category management", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("shows translated categories and benefit rate", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue(
      jsonResponse({
        items: [{
          id: "category-1",
          slug: "jewellery",
          displayName: "Jewellery",
          translations: { en: "Jewellery", fr: "Bijoux", ar: "مجوهرات", es: "Joyería" },
          benefitRateBasisPoints: 1250,
          isActive: true,
          updatedAt: "2026-09-27T10:00:00Z",
        }],
        page: 1,
        pageSize: 100,
        total: 1,
      }),
    ) as unknown as typeof fetch;

    render(<CategoryManagement />);
    await waitFor(() => expect(screen.getAllByText("Jewellery").length).toBeGreaterThan(0));
    expect(screen.getByText("12.50", { exact: false })).toBeInTheDocument();
    expect(screen.getByText("Bijoux")).toBeInTheDocument();
    expect(screen.getByText("مجوهرات")).toBeInTheDocument();
  });

  it("creates a category with all four language labels", async () => {
    const fetchMock = jest
      .fn()
      .mockResolvedValueOnce(jsonResponse({ items: [], page: 1, pageSize: 100, total: 0 }))
      .mockResolvedValueOnce(jsonResponse({ id: "category-2" }, 201))
      .mockResolvedValueOnce(jsonResponse({ items: [], page: 1, pageSize: 100, total: 0 }));
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<CategoryManagement />);
    await waitFor(() => expect(screen.getByRole("button", { name: "Add category" })).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Add category" }));
    fireEvent.change(screen.getByLabelText("Slug"), { target: { value: "jewellery" } });
    fireEvent.change(screen.getByLabelText("Default display name"), { target: { value: "Jewellery" } });
    for (const [locale, value] of [["en", "Jewellery"], ["fr", "Bijoux"], ["ar", "مجوهرات"], ["es", "Joyería"]]) {
      fireEvent.change(screen.getByLabelText(`Name (${locale})`), { target: { value } });
    }
    fireEvent.change(screen.getByLabelText("Benefit rate (%)"), { target: { value: "12.5" } });
    fireEvent.click(screen.getByRole("button", { name: "Save category" }));

    await waitFor(() => expect(fetchMock).toHaveBeenCalledWith(
      "/api/admin/categories",
      expect.objectContaining({ method: "POST", body: expect.stringContaining("\"benefitRateBasisPoints\":1250") }),
    ));
  });
});
