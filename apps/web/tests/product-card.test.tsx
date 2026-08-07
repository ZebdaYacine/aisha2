import { render, screen } from "@testing-library/react";
import { ProductCard } from "@/features/product/components/product-card";
import { products } from "@/features/catalogue/data";
import { storeCopy } from "@/core/lib/store-copy";
describe("ProductCard", () => {
  it("shows localized identity, artisan, price and availability", () => {
    render(
      <ProductCard product={products[0]} locale="ar" copy={storeCopy("ar")} />,
    );
    expect(
      screen.getAllByRole("link", { name: products[0].name.ar }),
    ).toHaveLength(2);
    expect(screen.getByText("Lila Aït Mansour")).toBeInTheDocument();
    expect(screen.getByText(storeCopy("ar").newArrivals)).toBeInTheDocument();
  });
});
