import { render, screen, waitFor } from "@testing-library/react";

import { AdminOverview } from "@/features/admin/components/admin-overview";

describe("admin overview", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads operational totals and keeps queue links available", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ total: 12 }),
    }) as unknown as typeof fetch;

    render(<AdminOverview locale="en" />);

    await waitFor(() => expect(screen.getAllByText("12")).toHaveLength(4));
    expect(screen.getByRole("heading", { name: "Good morning ✦" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /View queue/i })).toHaveAttribute("href", "/en/admin/artisan-applications");
    expect(screen.getByRole("link", { name: /View all/i })).toHaveAttribute("href", "/en/admin/audit");
  });
});
