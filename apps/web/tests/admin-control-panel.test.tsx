import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { AdminControlPanel } from "@/features/admin/components/admin-control-panel";
import { useAuth } from "@/features/auth/viewmodel/auth-context";

jest.mock("next/navigation", () => ({
  useRouter: () => ({ replace: jest.fn(), refresh: jest.fn() }),
}));
jest.mock("@/features/auth/viewmodel/auth-context", () => ({
  useAuth: jest.fn(),
}));

const mockedUseAuth = useAuth as jest.MockedFunction<typeof useAuth>;
const admin = {
  id: "admin-1",
  email: "admin@example.test",
  displayName: "Admin",
  status: "ACTIVE",
  roles: ["administrator"],
  customerEnabled: true,
  artisanStatus: "NOT_STARTED",
  artisanEnabled: false,
  capabilities: ["admin.artisan_applications.read"],
};

function jsonResponse(body: unknown) {
  return { ok: true, status: 200, json: async () => body } as Response;
}

function user(id: string) {
  return {
    id,
    email: `${id}@example.test`,
    displayName: id,
    status: "ACTIVE",
    roles: ["customer"],
    createdAt: "2026-01-01T00:00:00Z",
  };
}

describe("admin dashboard", () => {
  afterEach(() => {
    mockedUseAuth.mockReset();
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("renders paginated rows and opens user details", async () => {
    mockedUseAuth.mockReturnValue({
      user: admin,
      status: "authenticated",
      refresh: jest.fn(),
      setUser: jest.fn(),
      updateUser: jest.fn(),
      logout: jest.fn(),
    });
    globalThis.fetch = jest
      .fn()
      .mockImplementation(async (input: RequestInfo | URL) => {
        const url = String(input);
        if (url.includes("/admin/users")) {
          const page =
            new URL(url, "http://localhost").searchParams.get("page") ?? "1";
          return jsonResponse({
            items: [user(`user-${page}`)],
            page: Number(page),
            pageSize: 10,
            total: 20,
          });
        }
        return jsonResponse({
          items: [
            {
              id: "audit-1",
              eventType: "USER_STATUS_CHANGED",
              targetType: "user",
              occurredAt: "2026-01-01T00:00:00Z",
            },
          ],
          page: 1,
          pageSize: 5,
          total: 1,
        });
      }) as unknown as typeof fetch;

    render(<AdminControlPanel locale="en" />);
    await waitFor(() => expect(screen.getByText("user-1")).toBeInTheDocument());
    expect(screen.getByText("Page 1 of 2")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Details" }));
    expect(
      screen.getByRole("heading", { name: "User details" }),
    ).toBeInTheDocument();
    expect(screen.getAllByText("user-1@example.test")).toHaveLength(2);

    fireEvent.click(screen.getByRole("button", { name: "Next" }));
    await waitFor(() => expect(screen.getByText("user-2")).toBeInTheDocument());
  });
});
