import { fireEvent, render, screen } from "@testing-library/react";

import { StorefrontHeader } from "@/core/components/layout/storefront-header";
import { dictionary } from "@/core/lib/i18n";
import { CartProvider } from "@/features/cart/viewmodel/cart-context";

const mockReplace = jest.fn();
jest.mock("next/navigation", () => ({
  usePathname: () => "/en/products/kabyle-brooch",
  useRouter: () => ({ replace: mockReplace }),
}));

function renderHeader(locale: "en" | "ar") {
  const messages = dictionary(locale);
  return { messages, ...render(<CartProvider><StorefrontHeader locale={locale} messages={messages} /></CartProvider>) };
}

describe("StorefrontHeader", () => {
  beforeEach(() => {
    mockReplace.mockReset();
    window.history.replaceState({}, "", "/en/products/kabyle-brooch?sort=newest#story");
  });

  it("renders translated Arabic navigation and an accessible cart", () => {
    const { messages } = renderHeader("ar");

    expect(screen.getByRole("link", { name: messages.homeLabel })).toBeInTheDocument();
    expect(screen.getByRole("navigation", { name: messages.navigation.primary })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: messages.navigation.cart })).toBeInTheDocument();
    expect(screen.getAllByText(messages.navigation.products).length).toBeGreaterThan(0);
  });

  it("offers every supported language", () => {
    const { container } = renderHeader("en");

    expect(container.querySelector('a[lang="ar"]')).toHaveAttribute("href", "/ar/products/kabyle-brooch");
    expect(container.querySelector('a[lang="fr"]')).toHaveAttribute("href", "/fr/products/kabyle-brooch");
  });

  it("switches language without leaving the current route", () => {
    renderHeader("en");
    fireEvent.focus(screen.getByRole("combobox", { name: "Language" }));
    fireEvent.click(screen.getByRole("button", { name: /🇫🇷 FR/ }));
    expect(mockReplace).toHaveBeenCalledWith("/fr/products/kabyle-brooch?sort=newest#story");
  });
});
