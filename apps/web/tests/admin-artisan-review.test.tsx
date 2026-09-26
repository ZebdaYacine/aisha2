import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { ArtisanReview } from "@/core/components/admin/artisan-review";

function jsonResponse(body: unknown) {
  return { ok: true, status: 200, json: async () => body } as Response;
}

describe("admin artisan review media", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads application documents and profile media in the details dialog", async () => {
    globalThis.fetch = jest.fn().mockImplementation(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.includes("/documents")) {
        return jsonResponse([{ id: "document-1", documentType: "IDENTITY", originalFilename: "identity.jpg", mediaType: "image/jpeg", sizeBytes: 100, url: "/signed/identity" }]);
      }
      if (url.includes("/media")) {
        return jsonResponse([{ id: "media-1", mediaKind: "IMAGE", originalFilename: "workshop.jpg", mediaType: "image/jpeg", sizeBytes: 200, url: "/signed/workshop" }]);
      }
      return jsonResponse({ items: [{ id: "application-1", publicDisplayName: "Test Artisan", workshopName: "Test Workshop", wilaya: "Algiers", status: "SUBMITTED" }], page: 1, pageSize: 10, total: 1 });
    }) as unknown as typeof fetch;

    render(<ArtisanReview />);
    await waitFor(() => expect(screen.getByText("Test Artisan")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Details" }));

    expect(await screen.findByText("identity.jpg")).toBeInTheDocument();
    expect(screen.getByText("workshop.jpg")).toBeInTheDocument();
    expect(globalThis.fetch).toHaveBeenCalledWith("/api/admin/artisan-applications/application-1/documents");
    expect(globalThis.fetch).toHaveBeenCalledWith("/api/admin/artisan-applications/application-1/media");
  });
});
