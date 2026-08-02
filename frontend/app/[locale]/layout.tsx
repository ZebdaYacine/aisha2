import { notFound } from "next/navigation";

import { StorefrontFooter } from "@/components/layout/storefront-footer";
import { StorefrontHeader } from "@/components/layout/storefront-header";
import { dictionary, direction, isLocale, locales } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";
import { CartProvider } from "@/features/cart/cart-context";
import { CartDrawer } from "@/components/commerce/cart-drawer";
import { Toaster } from "sonner";

export function generateStaticParams() {
  return locales.map((locale) => ({ locale }));
}

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) {
    notFound();
  }
  const messages = dictionary(locale);
  const copy = storeCopy(locale);
  return (
    <html lang={locale} dir={direction(locale)} data-scroll-behavior="smooth">
      <body>
        <a
          className="fixed start-4 top-4 z-[100] -translate-y-24 bg-foreground px-4 py-3 text-sm text-background focus:translate-y-0"
          href="#main-content"
        >
          {messages.navigation.skipToContent}
        </a>
        <CartProvider>
          <StorefrontHeader locale={locale} messages={messages} />
          <main id="main-content">{children}</main>
          <StorefrontFooter locale={locale} messages={messages} />
          <CartDrawer locale={locale} copy={copy}/>
          <Toaster position="top-center" richColors closeButton />
        </CartProvider>
      </body>
    </html>
  );
}
