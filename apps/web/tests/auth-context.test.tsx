import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { StorefrontHeader } from "@/core/components/layout/storefront-header";
import { dictionary } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
import { CartProvider } from "@/features/cart/viewmodel/cart-context";
import { AuthProvider, useAuth } from "@/features/auth/viewmodel/auth-context";

const mockReplace = jest.fn();
jest.mock("next/navigation", () => ({
  usePathname: () => "/en/account",
  useRouter: () => ({ replace: mockReplace, refresh: jest.fn() }),
}));

const user = {
  id: "user-1",
  email: "amina@example.com",
  displayName: "Amina Sahra",
  status: "ACTIVE",
  roles: ["customer"],
  customerEnabled: true,
  artisanStatus: "NOT_STARTED",
  artisanEnabled: false,
  capabilities: ["customer.account.read"],
};

function jsonResponse(body: unknown, status = 200) {
  return { ok: status >= 200 && status < 300, status, json: async () => body } as Response;
}

describe("authentication state", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads the session and changes the account icon to the user's avatar", async () => {
    const fetchMock = jest.fn().mockResolvedValue(jsonResponse(user));
    globalThis.fetch = fetchMock as unknown as typeof fetch;
    const messages = dictionary("en");

    render(
      <AuthProvider>
        <CartProvider>
          <StorefrontHeader locale="en" messages={messages} />
        </CartProvider>
      </AuthProvider>,
    );

    await waitFor(() => expect(screen.getByTestId("account-avatar")).toHaveTextContent("AS"));
    expect(screen.getByRole("button", { name: messages.navigation.account })).toHaveAttribute("aria-expanded", "false");
    expect(screen.getByRole("button", { name: "Notifications" })).toBeInTheDocument();
    expect(screen.queryByRole("img")).not.toBeInTheDocument();
  });

  it("opens customer, artisan, and logout actions from the profile menu", async () => {
    const artisanUser = { ...user, artisanStatus: "ACTIVE", artisanEnabled: true, capabilities: [...user.capabilities, "artisan.account.read"] };
    globalThis.fetch = jest.fn().mockResolvedValue(jsonResponse(artisanUser)) as unknown as typeof fetch;
    const messages = dictionary("en");

    render(
      <AuthProvider>
        <CartProvider>
          <StorefrontHeader locale="en" messages={messages} />
        </CartProvider>
      </AuthProvider>,
    );

    await waitFor(() => expect(screen.getByTestId("account-avatar")).toHaveTextContent("AS"));
    fireEvent.click(screen.getByRole("button", { name: messages.navigation.account }));

    const copy = storeCopy("en");
    expect(screen.getByRole("menu")).toBeInTheDocument();
    expect(screen.getByRole("menuitem", { name: copy.customerArea })).toHaveAttribute("href", "/en/account");
    expect(screen.getByRole("menuitem", { name: copy.artisanArea })).toHaveAttribute("href", "/en/artisan");
    expect(screen.getByRole("menuitem", { name: copy.logout })).toBeInTheDocument();
  });

  it("clears the local session when signing out", async () => {
    window.localStorage.setItem("aisha-demo-cart", "cached-cart");
    window.sessionStorage.setItem("temporary-form", "cached-form");
    const fetchMock = jest.fn().mockImplementation(async (input: RequestInfo | URL) => {
      return String(input).endsWith("/logout") ? jsonResponse(null, 204) : jsonResponse(user);
    });
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    function Probe() {
      const { user: currentUser, logout } = useAuth();
      return (
        <>
          <span>{currentUser?.displayName}</span>
          <button type="button" onClick={() => void logout()}>
            Sign out
          </button>
        </>
      );
    }

    render(
      <AuthProvider>
        <Probe />
      </AuthProvider>,
    );

    await waitFor(() => expect(screen.getByText(user.displayName)).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Sign out" }));

    await waitFor(() => expect(screen.queryByText(user.displayName)).not.toBeInTheDocument());
    expect(window.localStorage.getItem("aisha-demo-cart")).toBeNull();
    expect(window.sessionStorage.getItem("temporary-form")).toBeNull();
    expect(fetchMock).toHaveBeenCalledWith("/api/auth/logout", {
      method: "POST",
      credentials: "same-origin",
    });
  });
});
