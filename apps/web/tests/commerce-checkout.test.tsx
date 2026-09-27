import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { toast } from "sonner";
import { CartProvider, useCart } from "@/features/cart/viewmodel/cart-context";
import { OrderSummary } from "@/core/components/commerce/order-summary";
import { PurchaseControls } from "@/core/components/commerce/purchase-controls";
import { WishlistList } from "@/features/wishlist/components/wishlist-list";
import { ImageLightbox } from "@/core/components/commerce/image-lightbox";
import { storeCopy } from "@/core/lib/store-copy";

jest.mock("sonner", () => ({ toast: { info: jest.fn(), error: jest.fn() } }));

function AddCartItem() {
  const { add } = useCart();
  return (
    <button
      type="button"
      onClick={() =>
        add("internal-product-id", 1, {
          productName: "Kabyle silver brooch",
          artisanName: "Oussama Atelier",
          workshopName: "Oussama Atelier",
          priceMinor: 12500,
          currency: "EUR",
        })
      }
    >
      Add test item
    </button>
  );
}

describe("commerce checkout controls", () => {
  beforeEach(() => {
    localStorage.clear();
    jest.clearAllMocks();
  });

  it("calculates checkout totals from cart metadata without exposing the product id", () => {
    render(
      <CartProvider>
        <AddCartItem />
        <OrderSummary locale="en" copy={storeCopy("en")} checkoutAction={false} />
      </CartProvider>,
    );

    fireEvent.click(screen.getByRole("button", { name: "Add test item" }));

    expect(screen.getAllByText("€125.00")).toHaveLength(2);
    expect(screen.queryByText("internal-product-id")).not.toBeInTheDocument();
  });

  it("gives unauthenticated users wishlist feedback without toggling the saved state", async () => {
    global.fetch = jest.fn().mockResolvedValue({ ok: false, status: 401 });

    render(
      <CartProvider>
        <PurchaseControls
          slug="internal-product-id"
          availability="in_stock"
          copy={storeCopy("en")}
          productName="Kabyle silver brooch"
          priceMinor={12500}
          currency="EUR"
        />
      </CartProvider>,
    );

    const wishlist = screen.getByRole("button", { name: storeCopy("en").wishlist });
    fireEvent.click(wishlist);

    await waitFor(() => expect(toast.info).toHaveBeenCalledWith(storeCopy("en").wishlistSignIn));
    expect(wishlist).toHaveAttribute("aria-pressed", "false");
  });

  it("renders wishlist product images and saved dates", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => [{ productId: "internal-product-id", productName: "Kabyle silver brooch", priceMinor: 12500, currency: "EUR", image: "/images/aisha/O5.jpg", active: true, createdAt: "2026-09-27T00:00:00Z" }],
    });

    render(<WishlistList locale="en" copy={storeCopy("en")} />);

    await waitFor(() => expect(screen.getByAltText("Kabyle silver brooch")).toBeInTheDocument());
    expect(screen.getByText(/September 27, 2026/)).toBeInTheDocument();
  });

  it("opens an image modal and zooms from the cursor position", () => {
    render(<ImageLightbox src="/images/aisha/O5.jpg" alt="Kabyle silver brooch" closeLabel="Close image" />);

    fireEvent.click(screen.getByRole("button", { name: "Kabyle silver brooch" }));
    const dialog = screen.getByRole("dialog", { name: "Kabyle silver brooch" });
    const images = screen.getAllByAltText("Kabyle silver brooch");
    const image = images[images.length - 1];
    expect(dialog).toContainElement(image);

    fireEvent.pointerMove(image?.parentElement as HTMLElement, { clientX: 20, clientY: 30 });
    fireEvent.click(image?.parentElement as HTMLElement);
    expect(image).toHaveStyle({ transform: "scale(2)" });
    expect(image?.style.transformOrigin).toMatch(/^\d+% \d+%$/);

    fireEvent.keyDown(document, { key: "Escape" });
    expect(screen.queryByRole("dialog", { name: "Kabyle silver brooch" })).not.toBeInTheDocument();
  });
});
