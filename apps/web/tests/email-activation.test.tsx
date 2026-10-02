import { render, screen, waitFor } from "@testing-library/react";

import { EmailActivation } from "@/core/components/forms/email-activation";
import { storeCopy } from "@/core/lib/store-copy";

function jsonResponse(body: unknown, status = 200) {
  return { ok: status >= 200 && status < 300, status, json: async () => body } as Response;
}

describe("email activation", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("shows activation success and a sign-in link without creating a session", async () => {
    const fetchMock = jest.fn().mockResolvedValue(jsonResponse({ activated: true, user: {} }));
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<EmailActivation locale="en" copy={storeCopy("en")} token="valid-token" />);

    expect(await screen.findByText("Your account is activated.")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Sign in" })).toHaveAttribute("href", "/en/login");
    expect(fetchMock).toHaveBeenCalledWith("/api/auth/activate?token=valid-token");
  });

  it("shows the expiration message for an invalid or expired token", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue(jsonResponse({ error: {} }, 401)) as unknown as typeof fetch;

    render(<EmailActivation locale="en" copy={storeCopy("en")} token="expired-token" />);

    await waitFor(() => expect(screen.getByRole("alert")).toHaveTextContent("This activation link has expired."));
  });
});
