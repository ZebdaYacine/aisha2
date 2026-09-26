import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { ArtisanWorkspace } from "@/core/components/artisan/workspace";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

jest.mock("@/features/auth/viewmodel/auth-context", () => ({
  useOptionalAuth: jest.fn(),
}));

function jsonResponse(body: unknown) {
  return { ok: true, status: 200, json: async () => body } as Response;
}

describe("artisan workspace tabs", () => {
  afterEach(() => {
    jest.clearAllMocks();
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("switches between workshops, private files, and product authoring", async () => {
    (useOptionalAuth as jest.Mock).mockReturnValue({
      user: { artisanEnabled: true, artisanStatus: "ACTIVE", capabilities: ["artisan.account.read"] },
    });
    globalThis.fetch = jest.fn().mockImplementation(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url === "/api/artisan/application") return jsonResponse({ status: "APPROVED", membershipStatus: "ACTIVE" });
      if (url === "/api/artisan/workshops") return jsonResponse({ items: [] });
      if (url === "/api/artisan/verification") return jsonResponse({ status: "NOT_SUBMITTED" });
      if (url === "/api/artisan/documents" || url === "/api/artisan/media") return jsonResponse([]);
      if (url === "/api/artisan/products") return jsonResponse({ items: [] });
      if (url.includes("/categories")) return jsonResponse({ items: [], page: 1, pageSize: 100, total: 0 });
      return jsonResponse({ items: [], page: 1, pageSize: 100, total: 0 });
    }) as unknown as typeof fetch;

    render(<ArtisanWorkspace locale="en" />);

    await waitFor(() => expect(screen.getByRole("tab", { name: /Workshops/ })).toHaveAttribute("aria-selected", "true"));
    expect(screen.queryByRole("heading", { name: "Application status" })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole("tab", { name: /Private files/ }));
    await waitFor(() => expect(screen.getByRole("tab", { name: /Private files/ })).toHaveAttribute("aria-selected", "true"));
    fireEvent.click(screen.getByRole("tab", { name: /Product authoring/ }));
    expect(screen.getByRole("tab", { name: /Product authoring/ })).toHaveAttribute("aria-selected", "true");
  });
});
