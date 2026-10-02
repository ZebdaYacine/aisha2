"use client";

import Link from "next/link";

import { Container } from "@/core/components/layout/container";
import { LanguageSwitcher } from "@/core/components/layout/language-switcher";
import type { Locale, Messages } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export function StorefrontFooter({
  locale,
  messages,
}: {
  locale: Locale;
  messages: Messages;
}) {
  const copy = storeCopy(locale);
  return (
    <footer className="border-t border-border bg-foreground text-background">
      <Container className="py-14 sm:py-18 lg:py-24">
        <div className="grid gap-12 md:grid-cols-2 lg:grid-cols-[1.5fr_repeat(3,1fr)]">
          <div className="max-w-sm">
            <p className="font-serif text-3xl tracking-[0.12em]">{messages.brand}</p>
            <p className="mt-5 text-sm leading-7 text-background/70">{messages.footer.story}</p>
          </div>
          <FooterGroup title={messages.footer.shop}>
            <Link href={`/${locale}/products`}>{messages.navigation.products}</Link>
            <Link href={`/${locale}/products`}>{messages.navigation.categories}</Link>
            <Link href={`/${locale}/products?sort=newest`}>{messages.navigation.new}</Link>
          </FooterGroup>
          <FooterGroup title={messages.footer.discover}>
            <Link href={`/${locale}/artisans`}>{messages.navigation.artisans}</Link>
            <Link href={`/${locale}/#story`}>{messages.navigation.story}</Link>
            <span>{messages.footer.support}</span>
            <Link href={`/${locale}/account/orders`}>{copy.orders}</Link>
          </FooterGroup>
          <FooterGroup title={messages.languageLabel}>
            <LanguageSwitcher locale={locale} messages={messages} />
            <label className="mt-3 text-xs text-background/60" htmlFor="footer-currency">
              {messages.currencyLabel}
            </label>
            <select id="footer-currency" className="h-11 max-w-52 rounded-md border border-background/30 bg-transparent px-3 text-sm text-background">
              <option className="text-foreground" value="DZD">DZD — {currencyName(locale, "DZD")}</option>
              <option className="text-foreground" value="EUR">EUR — {currencyName(locale, "EUR")}</option>
              <option className="text-foreground" value="USD">USD — {currencyName(locale, "USD")}</option>
            </select>
          </FooterGroup>
        </div>
        <div className="mt-16 flex flex-col gap-3 border-t border-background/20 pt-6 text-xs text-background/60 sm:flex-row sm:items-center sm:justify-between">
          <p>{messages.footer.copyright}</p>
          <p>{messages.footer.note}</p>
        </div>
      </Container>
    </footer>
  );
}

function currencyName(locale: Locale, currency: "DZD" | "EUR" | "USD") {
  const names = {
    en: { DZD: "Algerian dinar", EUR: "Euro", USD: "US dollar" },
    fr: { DZD: "dinar algérien", EUR: "euro", USD: "dollar américain" },
    ar: { DZD: "الدينار الجزائري", EUR: "اليورو", USD: "الدولار الأمريكي" },
    es: { DZD: "dinar argelino", EUR: "euro", USD: "dólar estadounidense" },
  } as const;
  return names[locale][currency];
}

function FooterGroup({ children, title }: { children: React.ReactNode; title: string }) {
  return (
    <div>
      <h2 className="latin-tracking text-xs font-medium uppercase tracking-[0.16em]">{title}</h2>
      <div className="mt-5 flex flex-col gap-3 text-sm text-background/70">{children}</div>
    </div>
  );
}
