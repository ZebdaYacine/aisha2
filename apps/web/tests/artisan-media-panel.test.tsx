import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { ArtisanMediaPanel } from "@/features/artisan/components/artisan-media-panel";

function response(body: unknown, status = 200) {
  return { ok: status >= 200 && status < 300, status, json: async () => body } as Response;
}

describe("artisan private media", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
    jest.restoreAllMocks();
  });

  it("renders media rows and opens add media modal", async () => {
    const fetchMock = jest.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      if (init?.method === "POST") return response({ id: "media-2", mediaKind: "IMAGE", originalFilename: "new.png", mediaType: "image/png", createdAt: "2026-01-16T00:00:00Z", url: "https://media.test/new" }, 201);
      if (String(input).includes("/documents")) return response([]);
      return response([{ id: "media-1", mediaKind: "IMAGE", originalFilename: "lamp.png", mediaType: "image/png", createdAt: "2026-01-15T00:00:00Z", url: "https://media.test/lamp" }]);
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<ArtisanMediaPanel locale="en" />);

    await waitFor(() => expect(screen.getByText("lamp.png")).toBeInTheDocument());
    expect(screen.getByText("January 15, 2026")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "View" })).toHaveAttribute("href", "https://media.test/lamp");
    fireEvent.click(screen.getByRole("button", { name: "Add media" }));
    expect(screen.getByRole("dialog", { name: "Add private media" })).toBeInTheDocument();
    expect(screen.getByRole("combobox", { name: "Media type" })).toBeInTheDocument();
    fireEvent.keyDown(document, { key: "Escape" });
    expect(screen.queryByRole("dialog", { name: "Add private media" })).not.toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalled();
  });
});
