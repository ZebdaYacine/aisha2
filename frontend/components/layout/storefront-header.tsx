"use client";
import { Menu, ShoppingBag, UserRound, X } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

import { Container } from "@/components/layout/container";
import type { Locale, Messages } from "@/lib/i18n";
import { locales } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";
import { useCart } from "@/features/cart/cart-context";
import { GlobalSearch } from "./global-search";
import type { Artisan, Category, Product } from "@/features/storefront/types";

const navItems = [
  ["new", "/products?sort=newest"],
  ["products", "/products"],
  ["categories", "/products"],
  ["artisans", "/artisans"],
  ["story", "/#story"],
] as const;

export function StorefrontHeader({
  locale,
  messages,
  catalogue,
}: {
  locale: Locale;
  messages: Messages;
  catalogue?: { products: Product[]; artisans: Artisan[]; categories: Category[] };
}) {
  const { count, setOpen: setCartOpen } = useCart();
  const [menuOpen, setMenuOpen] = useState(false);
  const copy = storeCopy(locale);
  const cartCount = new Intl.NumberFormat(locale).format(count);
  return (
    <>
      <div className="bg-foreground py-2 text-center text-[0.6875rem] tracking-[0.12em] text-background">
        <Container>{messages.announcement}</Container>
      </div>
      <header className="sticky top-0 z-40 border-b border-border/80 bg-background/95 backdrop-blur-sm">
        <Container className="grid min-h-16 grid-cols-[1fr_auto_1fr] items-center gap-3 lg:min-h-18">
          <nav className="hidden items-center gap-5 text-xs lg:flex" aria-label={messages.navigation.primary}>
            {navItems.map(([key, path]) => (
              <Link
                className="group relative py-2"
                href={`/${locale}${path}`}
                key={key}
              >
                {messages.navigation[key]}
                <span className="absolute inset-x-0 bottom-0 h-px origin-end scale-x-0 bg-current transition-transform duration-200 group-hover:origin-start group-hover:scale-x-100" />
              </Link>
            ))}
          </nav>

          <div className="justify-self-start lg:hidden">
            <button
              type="button"
              onClick={() => setMenuOpen(!menuOpen)}
              className="flex min-h-11 min-w-11 items-center justify-center"
              aria-label={messages.navigation.openMenu}
            >
              {menuOpen ? <X aria-hidden size={21}/> : <Menu aria-hidden="true" size={21} strokeWidth={1.5} />}
            </button>
            {menuOpen && <div className="fixed inset-x-0 top-[6.5rem] z-50 min-h-[calc(100svh-6.5rem)] border-b border-border bg-background px-4 py-8 shadow-sm">
              <nav className="flex flex-col" aria-label={messages.navigation.mobile}>
                {navItems.map(([key, path]) => (
                  <Link
                    className="border-b border-border py-4 font-serif text-3xl"
                    href={`/${locale}${path}`}
                    key={key}
                    onClick={() => setMenuOpen(false)}
                  >
                    {messages.navigation[key]}
                  </Link>
                ))}
              </nav>
            </div>}
          </div>

          <Link
            href={`/${locale}`}
            className="justify-self-center font-serif text-2xl tracking-[0.16em]"
            aria-label={messages.homeLabel}
          >
            {messages.brand}
          </Link>

          <div className="flex items-center justify-end gap-1 sm:gap-2">
            <GlobalSearch locale={locale} copy={copy} {...catalogue}/>
            <Link
              className="hidden min-h-11 min-w-11 items-center justify-center sm:flex"
              href={`/${locale}/login`}
              aria-label={messages.navigation.account}
            >
              <UserRound aria-hidden="true" size={19} strokeWidth={1.5} />
            </Link>
            <button
              type="button"
              onClick={() => setCartOpen(true)}
              className="relative flex min-h-11 min-w-11 items-center justify-center"
              aria-label={messages.navigation.cart}
            >
              <ShoppingBag aria-hidden="true" size={19} strokeWidth={1.5} />
              <span className="absolute end-0 top-1 text-[0.625rem]" aria-label={messages.navigation.cartCount}>
                {cartCount}
              </span>
            </button>
            <details className="relative">
              <summary className="latin-tracking flex min-h-11 cursor-pointer list-none items-center px-2 text-[0.6875rem] font-medium uppercase tracking-[0.12em] [&::-webkit-details-marker]:hidden">
                {locale}
              </summary>
              <div className="absolute end-0 top-full min-w-32 border border-border bg-background p-2">
                {locales.map((item) => (
                  <Link
                    key={item}
                    className="latin-tracking block px-3 py-2 text-xs uppercase tracking-[0.1em] hover:bg-muted"
                    href={`/${item}`}
                    hrefLang={item}
                    lang={item}
                  >
                    {messages.locales[item]}
                  </Link>
                ))}
              </div>
            </details>
          </div>
        </Container>
      </header>
    </>
  );
}
