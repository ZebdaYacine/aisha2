import { render, screen, waitFor } from "@testing-library/react";

import { BrandSplash } from "@/core/components/layout/brand-splash";

describe("brand splash", () => {
  afterEach(() => {
    window.localStorage.clear();
    jest.restoreAllMocks();
  });

  it("shows at the beginning when the cooldown has expired", async () => {
    jest.spyOn(window, "requestAnimationFrame").mockImplementation((callback) => {
      callback(performance.now() + 2200);
      return 1;
    });

    render(<BrandSplash locale="en" />);

    await waitFor(() => expect(screen.getByRole("status", { name: "Opening the collection" })).toBeInTheDocument());
  });

  it("stays hidden during the two-minute cooldown", () => {
    window.localStorage.setItem("aisha:splash-seen-at", String(Date.now() - 60_000));

    render(<BrandSplash locale="en" />);

    expect(screen.queryByRole("status", { name: "Opening the collection" })).not.toBeInTheDocument();
  });
});
