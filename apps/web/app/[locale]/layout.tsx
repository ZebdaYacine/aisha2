import { notFound } from "next/navigation";

import { AppProviders } from "@/core/context/app-providers";
import { StorefrontFooter } from "@/core/components/layout/storefront-footer";
import { StorefrontHeader } from "@/core/components/layout/storefront-header";
import { BrandSplash } from "@/core/components/layout/brand-splash";
import { dictionary, direction, isLocale, locales } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

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
    <html lang={locale} dir={direction(locale)} data-scroll-behavior="smooth" suppressHydrationWarning>
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html: `(() => { try { const saved = localStorage.getItem('aisha-theme'); const dark = saved === 'dark' || (!saved && window.matchMedia('(prefers-color-scheme: dark)').matches); document.documentElement.classList.toggle('dark', dark); document.documentElement.style.colorScheme = dark ? 'dark' : 'light'; } catch (_) {} })();`,
          }}
        />
      </head>
      <body>
        <a
          className="fixed start-4 top-4 z-[100] -translate-y-24 bg-foreground px-4 py-3 text-sm text-background focus:translate-y-0"
          href="#main-content"
        >
          {messages.navigation.skipToContent}
        </a>
        <AppProviders locale={locale} copy={copy}>
            <BrandSplash locale={locale} />
            <StorefrontHeader locale={locale} messages={messages} />
            <main id="main-content">{children}</main>
            <StorefrontFooter locale={locale} messages={messages} />
        </AppProviders>
      </body>
    </html>
  );
}
