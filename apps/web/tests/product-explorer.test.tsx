import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { ProductExplorer } from "@/core/components/commerce/product-explorer";
import { storeCopy } from "@/core/lib/store-copy";

const replace = jest.fn();

jest.mock("next/navigation", () => ({
  usePathname: () => "/en/products",
  useRouter: () => ({ replace }),
  useSearchParams: () => new URLSearchParams(),
}));

describe("ProductExplorer", () => {
  beforeEach(() => replace.mockClear());

  it("does not replace an unchanged URL and updates it once after sorting", async () => {
    render(<ProductExplorer initialProducts={[]} locale="en" copy={storeCopy("en")} categories={[]} />);

    expect(replace).not.toHaveBeenCalled();
    fireEvent.change(screen.getByRole("combobox"), { target: { value: "low" } });

    await waitFor(() => expect(replace).toHaveBeenCalledWith("/en/products?sort=low", { scroll: false }));
    expect(replace).toHaveBeenCalledTimes(1);
  });
});
