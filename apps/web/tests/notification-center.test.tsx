import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { NotificationCenter } from "@/features/notification/components/notification-center";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

const mockPush = jest.fn();

jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

jest.mock("@/features/auth/viewmodel/auth-context", () => ({
  useOptionalAuth: jest.fn(),
}));

const user = {
  id: "moderator-1",
  email: "moderator@example.test",
  displayName: "Moderator",
  status: "ACTIVE",
  roles: ["moderator"],
  customerEnabled: true,
  artisanStatus: "NOT_STARTED",
  artisanEnabled: false,
  capabilities: ["admin.product_moderation.read"],
};

function jsonResponse(body: unknown, status = 200) {
  return { ok: status >= 200 && status < 300, status, json: async () => body } as Response;
}

describe("notification center", () => {
  beforeEach(() => {
    mockPush.mockReset();
    (useOptionalAuth as jest.Mock).mockReturnValue({ status: "authenticated", user });
  });

  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("shows the bell workflow and routes a clicked moderation notification", async () => {
    const fetchMock = jest.fn().mockImplementation(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.includes("/api/notifications/ws-ticket")) return jsonResponse({}, 404);
      if (url === "/api/notifications?page=1&pageSize=20") {
        return jsonResponse({
          unreadCount: 1,
          items: [{
            id: "notification-1",
            eventType: "PRODUCT_MODERATION_DECIDED",
            titleKey: "notifications.product.moderation.title",
            bodyKey: "notifications.product.body",
            payload: { productId: "product-1" },
            createdAt: "2026-10-02T12:00:00Z",
          }],
        });
      }
      return jsonResponse(null, 204);
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<NotificationCenter locale="en" />);
    const bell = await screen.findByRole("button", { name: "Notifications" });
    await waitFor(() => expect(bell).toHaveTextContent("1"));

    fireEvent.click(bell);
    const notification = await screen.findByRole("button", { name: /Product moderation decision/ });
    fireEvent.click(notification);

    await waitFor(() => expect(mockPush).toHaveBeenCalledWith("/en/admin/moderation"));
    expect(fetchMock).toHaveBeenCalledWith("/api/notifications/notification-1/read", { method: "POST" });
    expect(screen.queryByRole("dialog", { name: "Notifications" })).not.toBeInTheDocument();
  });
});
