import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { ArtisanMembershipReview } from "@/features/admin/components/artisan-membership-review";
import { AdminMediaOperations } from "@/features/admin/components/admin-media-operations";
import { ProductModeration } from "@/features/admin/components/product-moderation";

function jsonResponse(body: unknown) {
  return { ok: true, status: 200, json: async () => body } as Response;
}

describe("admin moderation details", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("shows product media in the moderation details dialog without exposing the product id", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue(
      jsonResponse({
        items: [{
          submissionId: "submission-1",
          productId: "internal-product-id",
          productName: "Hand-painted bowl",
          productStatus: "PENDING_REVIEW",
          priceMinor: 2500,
          currency: "DZD",
          artisanName: "Test Artisan",
          workshopName: "Test Workshop",
          media: [{ id: "media-1", mediaKind: "IMAGE", originalFilename: "bowl.jpg", mediaType: "image/jpeg", sizeBytes: 100, altText: "Bowl", visibility: "PRIVATE", url: "/signed/bowl" }],
        }],
        page: 1,
        pageSize: 10,
        total: 1,
      }),
    ) as unknown as typeof fetch;

    render(<ProductModeration />);
    await waitFor(() => expect(screen.getByText("Hand-painted bowl")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Details" }));

    expect(await screen.findByText("bowl.jpg")).toBeInTheDocument();
    expect(screen.queryByText("internal-product-id")).not.toBeInTheDocument();
  });

  it("renders approved artisans as shop rows and hides internal ids in details", async () => {
    globalThis.fetch = jest.fn().mockImplementation(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.includes("/media")) return jsonResponse([{ id: "artisan-media-1", mediaKind: "IMAGE", originalFilename: "profile.jpg", mediaType: "image/jpeg", sizeBytes: 100, url: "/signed/profile" }]);
      return jsonResponse({ items: [{ id: "profile-internal-id", publicDisplayName: "Test Artisan", workshopName: "Test Workshop", wilaya: "Algiers", status: "APPROVED", membershipStatus: "ACTIVE" }] });
    }) as unknown as typeof fetch;

    render(<ArtisanMembershipReview />);
    await waitFor(() => expect(screen.getByText("Test Workshop")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Details" }));

    expect(screen.getAllByText("Shop").length).toBeGreaterThan(0);
    expect(await screen.findByText("profile.jpg")).toBeInTheDocument();
    expect(screen.queryByText("profile-internal-id")).not.toBeInTheDocument();
  });

  it("confirms and deletes media from the admin media tab", async () => {
    const fetchMock = jest.fn().mockImplementation(async (_input: RequestInfo | URL, init?: RequestInit) => {
      if (init?.method === "DELETE") return { ok: true, status: 204 } as Response;
      return jsonResponse({ items: [{ id: "media-1", userId: "user-1", userDisplayName: "Test Artisan", userEmail: "artisan@example.test", mediaKind: "PROFILE_MEDIA", originalFilename: "profile.jpg", mediaType: "image/jpeg", sizeBytes: 100, url: "/signed/profile", createdAt: "2026-09-26T00:00:00Z" }], page: 1, total: 1 });
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<AdminMediaOperations locale="en" />);
    await waitFor(() => expect(screen.getByText("profile.jpg")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Delete" }));
    expect(screen.getByRole("dialog", { name: "Delete" })).toBeInTheDocument();
    fireEvent.click(screen.getAllByRole("button", { name: "Delete" }).at(-1)!);

    await waitFor(() => expect(fetchMock).toHaveBeenCalledWith("/api/admin/media/users/media-1", { method: "DELETE" }));
  });
});
