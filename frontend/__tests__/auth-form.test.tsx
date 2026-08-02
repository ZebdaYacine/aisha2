import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { AuthForm } from "@/components/forms/auth-form";
import { storeCopy } from "@/lib/store-copy";

describe("AuthForm validation", () => {
  it("shows accessible translated field errors and invalid styling", async () => {
    render(<AuthForm mode="login" locale="en" copy={storeCopy("en")} />);
    fireEvent.change(screen.getByLabelText("Email address"), { target: { value: "invalid" } });
    fireEvent.click(screen.getByRole("button", { name: "Sign in" }));

    await waitFor(() => expect(screen.getByText("Please correct the highlighted fields.")).toHaveFocus());
    expect(screen.getByLabelText("Email address")).toHaveAttribute("aria-invalid", "true");
    expect(screen.getByText("Enter a valid email address.")).toBeInTheDocument();
    expect(screen.getByText("Use at least 12 characters.")).toBeInTheDocument();
  });
});
