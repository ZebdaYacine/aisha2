import { render, screen } from "@testing-library/react";

import { AccountOverview } from "@/features/account/components/account-overview";
import { storeCopy } from "@/core/lib/store-copy";
import { useAuth } from "@/features/auth/viewmodel/auth-context";
import type { AuthUser } from "@/features/auth/types";

jest.mock("@/features/auth/viewmodel/auth-context", () => ({ useAuth: jest.fn() }));

const mockedUseAuth = useAuth as jest.MockedFunction<typeof useAuth>;
const copy = storeCopy("en");
const user: AuthUser = {
  id: "user-1",
  email: "customer@example.test",
  displayName: "Customer",
  status: "ACTIVE",
  customerEnabled: true,
  artisanStatus: "NOT_STARTED",
  artisanEnabled: false,
  capabilities: ["customer.account.read"],
};

describe("account overview", () => {
  afterEach(() => mockedUseAuth.mockReset());

  it("shows a loading state while backend-derived account data is pending", () => {
    mockedUseAuth.mockReturnValue({ user: null, status: "loading", refresh: jest.fn(), setUser: jest.fn(), updateUser: jest.fn(), logout: jest.fn() });
    render(<AccountOverview locale="en" copy={copy} />);
    expect(screen.getByRole("status")).toHaveTextContent(copy.accountLoading);
  });

  it("shows customer and artisan state from the authenticated user", () => {
    mockedUseAuth.mockReturnValue({ user, status: "authenticated", refresh: jest.fn(), setUser: jest.fn(), updateUser: jest.fn(), logout: jest.fn() });
    render(<AccountOverview locale="en" copy={copy} />);
    expect(screen.getByText(copy.customerAccess)).toBeInTheDocument();
    expect(screen.getByText(copy.enabled)).toBeInTheDocument();
    expect(screen.getByText(copy.notStarted)).toBeInTheDocument();
  });

  it("shows an unavailable state when the session cannot be resolved", () => {
    mockedUseAuth.mockReturnValue({ user: null, status: "unauthenticated", refresh: jest.fn(), setUser: jest.fn(), updateUser: jest.fn(), logout: jest.fn() });
    render(<AccountOverview locale="en" copy={copy} />);
    expect(screen.getByRole("alert")).toHaveTextContent(copy.accountUnavailable);
    expect(screen.getByRole("link", { name: copy.login })).toHaveAttribute("href", "/en/login");
  });
});
