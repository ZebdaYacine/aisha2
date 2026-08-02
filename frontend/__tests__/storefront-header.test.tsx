import { render, screen } from "@testing-library/react";

import { StorefrontHeader } from "@/components/layout/storefront-header";
import { dictionary } from "@/lib/i18n";
import { CartProvider } from "@/features/cart/cart-context";

function renderHeader(locale: "en" | "ar") {
  const messages = dictionary(locale);
  return { messages, ...render(<CartProvider><StorefrontHeader locale={locale} messages={messages} /></CartProvider>) };
}

describe("StorefrontHeader", () => {
  it("renders translated Arabic navigation and an accessible cart", () => {
    const { messages } = renderHeader("ar");

    expect(screen.getByRole("link", { name: messages.homeLabel })).toBeInTheDocument();
    expect(screen.getByRole("navigation", { name: messages.navigation.primary })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: messages.navigation.cart })).toBeInTheDocument();
    expect(screen.getAllByText(messages.navigation.products).length).toBeGreaterThan(0);
  });

  it("offers every supported language", () => {
    const { container } = renderHeader("en");

    expect(container.querySelector('a[lang="ar"]')).toHaveAttribute("href", "/ar");
    expect(container.querySelector('a[lang="fr"]')).toHaveAttribute("href", "/fr");
  });
});
