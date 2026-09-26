import { render, screen, waitFor } from "@testing-library/react";

import { AccountSidebar } from "@/core/components/account/account-sidebar";
import { direction } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
import { useAuth } from "@/features/auth/viewmodel/auth-context";
import type { AuthUser } from "@/features/auth/types";

jest.mock("next/navigation", () => ({
  useRouter: () => ({ replace: jest.fn(), refresh: jest.fn() }),
}));
jest.mock("@/features/auth/viewmodel/auth-context", () => ({ useAuth: jest.fn() }));

const mockedUseAuth = useAuth as jest.MockedFunction<typeof useAuth>;

function renderSidebar(user: AuthUser, locale: "en" | "ar" = "en", mode: "customer" | "artisan" = "customer") {
  mockedUseAuth.mockReturnValue({
    user,
    status: "authenticated",
    refresh: jest.fn(),
    setUser: jest.fn(),
    updateUser: jest.fn(),
    logout: jest.fn(),
  });
  render(<AccountSidebar locale={locale} copy={storeCopy(locale)} mode={mode} />);
}

const baseUser: AuthUser = {
  id: "same-user",
  email: "artisan@example.test",
  displayName: "Artisan Buyer",
  status: "ACTIVE",
  customerEnabled: true,
  artisanStatus: "NOT_STARTED",
  artisanEnabled: false,
  capabilities: ["customer.account.read", "customer.profile.write"],
};

describe("combined account navigation", () => {
  afterEach(() => {
    mockedUseAuth.mockReset();
  });

  it("keeps customer navigation visible for a pending artisan", async () => {
    renderSidebar({ ...baseUser, artisanStatus: "PENDING", artisanEnabled: false });
    const copy = storeCopy("en");

    await waitFor(() => expect(screen.getByText(copy.artisanPending)).toBeInTheDocument());
    expect(screen.getByRole("link", { name: copy.profile })).toHaveAttribute("href", "/en/account/profile");
    expect(screen.getByRole("link", { name: copy.addresses })).toHaveAttribute("href", "/en/account/addresses");
  });

  it("shows the seller mode only from backend artisan capabilities", async () => {
    renderSidebar({
      ...baseUser,
      artisanStatus: "ACTIVE",
      artisanEnabled: true,
      capabilities: [...baseUser.capabilities, "artisan.account.read"],
    });
    const copy = storeCopy("en");

    await waitFor(() => expect(screen.getByRole("link", { name: copy.artisanArea })).toHaveAttribute("href", "/en/artisan"));
    expect(screen.getByRole("link", { name: copy.purchases })).toHaveAttribute("href", "/en/account/orders");
  });

  it("shares the customer switch, workshop management, and logout in artisan mode", async () => {
    renderSidebar({
      ...baseUser,
      artisanStatus: "ACTIVE",
      artisanEnabled: true,
      capabilities: [...baseUser.capabilities, "artisan.account.read"],
    }, "en", "artisan");
    const copy = storeCopy("en");
    await waitFor(() => expect(screen.getByRole("link", { name: copy.customerArea })).toHaveAttribute("href", "/en/account"));
    expect(screen.getByRole("link", { name: copy.manageWorkshops })).toHaveAttribute("href", "/en/artisan#workshops");
    expect(screen.getByRole("button", { name: copy.logout })).toBeInTheDocument();
  });

  it("keeps customer mode and disables seller navigation for a suspended artisan", async () => {
    renderSidebar({ ...baseUser, artisanStatus: "SUSPENDED", artisanEnabled: false });
    const copy = storeCopy("en");

    await waitFor(() => expect(screen.getByText(copy.artisanSuspended)).toBeInTheDocument());
    expect(screen.queryByRole("link", { name: copy.artisanArea })).not.toBeInTheDocument();
    expect(screen.getByRole("link", { name: copy.tracking })).toHaveAttribute("href", "/en/account/tracking");
  });

  it("renders translated Arabic navigation with responsive overflow classes", async () => {
    renderSidebar({ ...baseUser, artisanStatus: "NOT_STARTED", artisanEnabled: false }, "ar");
    const copy = storeCopy("ar");

    await waitFor(() => expect(screen.getByRole("link", { name: copy.purchases })).toBeInTheDocument());
    const navigation = screen.getByRole("navigation", { name: copy.account });
    expect(direction("ar")).toBe("rtl");
    expect(navigation.className).toContain("overflow-x-auto");
    expect(navigation.className).toContain("lg:flex-col");
  });

  it("does not treat a forged local role as seller authorization", async () => {
    renderSidebar({ ...baseUser, roles: ["artisan"], artisanStatus: "NOT_STARTED", artisanEnabled: false });
    const copy = storeCopy("en");

    await waitFor(() => expect(screen.getByRole("link", { name: copy.profile })).toBeInTheDocument());
    expect(screen.queryByRole("link", { name: copy.artisanArea })).not.toBeInTheDocument();
  });
});
