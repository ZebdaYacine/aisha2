import { fireEvent, render, screen } from "@testing-library/react";
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

  it("zooms the image toward the cursor while preserving the product link", () => {
    const { container } = render(
      <ProductCard product={products[0]} locale="en" copy={storeCopy("en")} />,
    );
    const media = container.querySelector(".product-media") as HTMLElement;
    const image = container.querySelector(".product-card-image") as HTMLElement;

    fireEvent.pointerMove(media, { pointerType: "mouse", clientX: 20, clientY: 30 });

    expect(image).toHaveStyle({ transform: "scale(1.2)" });
    expect(image.style.transformOrigin).toMatch(/^\d+% \d+%$/);
    expect(container.querySelector('a[href^="/en/products/"]')).toBeInTheDocument();
  });
});
