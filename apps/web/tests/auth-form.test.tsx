import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { AuthForm } from "@/features/auth/components/forms/auth-form";
import { storeCopy } from "@/core/lib/store-copy";

describe("AuthForm validation", () => {
  it("shows accessible translated field errors and invalid styling", async () => {
    render(<AuthForm mode="login" locale="en" copy={storeCopy("en")} />);
    fireEvent.click(screen.getByRole("button", { name: "Sign in" }));

    await waitFor(() => expect(screen.getByText("Please correct the highlighted fields.")).toHaveFocus());
    expect(screen.getByLabelText("Email or phone number")).toHaveAttribute("aria-invalid", "true");
    expect(screen.getByText("This field is required.")).toBeInTheDocument();
    expect(screen.getByText("Use at least 12 characters.")).toBeInTheDocument();
  });
});
