"use client";
import { Menu, ShoppingBag, UserRound, X } from "lucide-react";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { usePathname, useRouter } from "next/navigation";

import { Container } from "@/core/components/layout/container";
import { ThemeToggle } from "@/core/components/layout/theme-toggle";
import { LanguageSwitcher, pathForLocale } from "@/core/components/layout/language-switcher";
import type { Locale, Messages } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
import { useCart } from "@/features/cart/viewmodel/cart-context";
import { GlobalSearch } from "./global-search";
import type { Artisan, Category, Product, Workshop } from "@/features/catalogue/types";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";
import { hasCapability, userInitials } from "@/features/auth/types";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import { NotificationCenter } from "@/features/notification/components/notification-center";

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
  catalogue?: { products: Product[]; artisans: Artisan[]; categories: Category[]; workshops?: Workshop[] };
}) {
  const { count, setOpen: setCartOpen } = useCart();
  const auth = useOptionalAuth();
  const pathname = usePathname() || `/${locale}`;
  const router = useRouter();
  const [menuOpen, setMenuOpen] = useState(false);
  const [accountMenuOpen, setAccountMenuOpen] = useState(false);
  const [loggingOut, setLoggingOut] = useState(false);
  const accountMenuRef = useRef<HTMLDivElement>(null);
  const copy = storeCopy(locale);
  const cartCount = new Intl.NumberFormat(locale).format(count);
  const artisanAccess = auth?.user?.artisanEnabled === true && hasCapability(auth.user, "artisan.account.read");

  useEffect(() => {
    if (!accountMenuOpen) return;
    const closeOnOutsideClick = (event: MouseEvent) => {
      if (!accountMenuRef.current?.contains(event.target as Node)) setAccountMenuOpen(false);
    };
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") setAccountMenuOpen(false);
    };
    document.addEventListener("mousedown", closeOnOutsideClick);
    document.addEventListener("keydown", closeOnEscape);
    return () => {
      document.removeEventListener("mousedown", closeOnOutsideClick);
      document.removeEventListener("keydown", closeOnEscape);
    };
  }, [accountMenuOpen]);

  const signOut = async () => {
    if (!auth?.logout || loggingOut) return;
    setLoggingOut(true);
    try {
      await auth.logout();
      setAccountMenuOpen(false);
      router.replace(`/${locale}/login`);
      router.refresh();
    } finally {
      setLoggingOut(false);
    }
  };
  return (
    <>
      <div className="bg-foreground py-2 text-center text-[0.6875rem] tracking-[0.12em] text-background">
        <Container>{messages.announcement}</Container>
      </div>
      <header className="relative sticky top-0 z-40 border-b border-border/80 bg-background/95 backdrop-blur-sm">
        <Container className="grid min-h-16 min-w-0 grid-cols-[auto_auto_minmax(0,1fr)] items-center gap-1 sm:grid-cols-[1fr_auto_1fr] sm:gap-3 lg:min-h-18">
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
            {menuOpen && <div className="absolute inset-x-0 top-full z-50 min-h-[calc(100svh-4rem)] overflow-y-auto border-b border-border bg-background px-4 py-8 shadow-sm">
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
              <div className="mt-8 flex items-center justify-between border-t border-border pt-5">
                <ThemeToggle locale={locale} />
              </div>
            </div>}
          </div>

          <Link
            href={`/${locale}`}
            className="justify-self-center font-serif text-xl tracking-[0.16em] sm:text-2xl"
            aria-label={messages.homeLabel}
          >
            {messages.brand}
          </Link>

          <div className="flex min-w-0 items-center justify-end gap-0 sm:gap-2">
            <GlobalSearch locale={locale} copy={copy} {...catalogue}/>
            <div className="hidden lg:block">
              <ThemeToggle locale={locale} />
            </div>
            {auth?.user ? (
              <div className="relative" ref={accountMenuRef}>
                <button
                  type="button"
                  className="flex min-h-11 min-w-11 items-center justify-center"
                  aria-label={messages.navigation.account}
                  aria-haspopup="menu"
                  aria-expanded={accountMenuOpen}
                  onClick={() => setAccountMenuOpen((open) => !open)}
                >
                  <span
                    data-testid="account-avatar"
                    className="flex size-8 items-center justify-center rounded-full bg-foreground text-[0.6875rem] font-medium tracking-[0.08em] text-background"
                    aria-hidden="true"
                  >
                    {userInitials(auth.user)}
                  </span>
                </button>
                {accountMenuOpen && (
                  <div className="absolute end-0 top-full z-50 mt-2 w-64 max-w-[calc(100vw-2rem)] border border-border bg-background p-2 shadow-lg" role="menu" aria-label={copy.account}>
                    <div className="border-b border-border px-3 py-3">
                      <p className="truncate text-sm font-medium">{auth.user.displayName}</p>
                      <p className="truncate text-xs text-muted-foreground">{auth.user.email}</p>
                    </div>
                    <LocalizedLink className="mt-2 flex min-h-11 items-center px-3 py-2 text-sm hover:bg-muted" locale={locale} href="/account" role="menuitem" onClick={() => setAccountMenuOpen(false)}>
                      {copy.customerArea}
                    </LocalizedLink>
                    {artisanAccess && <LocalizedLink className="flex min-h-11 items-center px-3 py-2 text-sm hover:bg-muted" locale={locale} href="/artisan" role="menuitem" onClick={() => setAccountMenuOpen(false)}>
                      {copy.artisanArea}
                    </LocalizedLink>}
                    <button type="button" className="mt-1 flex min-h-11 w-full items-center px-3 py-2 text-start text-sm text-destructive hover:bg-muted" role="menuitem" disabled={loggingOut} onClick={() => void signOut()}>
                      {loggingOut ? `${copy.logout}…` : copy.logout}
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <Link className="flex min-h-11 min-w-11 items-center justify-center" href={`/${locale}/login`} aria-label={messages.navigation.account}>
                <UserRound aria-hidden="true" size={19} strokeWidth={1.5} />
              </Link>
            )}
            <span className="block shrink-0"><NotificationCenter locale={locale} /></span>
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
            <LanguageSwitcher locale={locale} messages={messages} compact />
            <nav className="sr-only" aria-label="Language links">
              {(["en", "fr", "ar", "es"] as const).map((item) => (
                  <Link key={item} href={pathForLocale(pathname, item)} hrefLang={item} lang={item}>
                  {messages.locales[item]}
                </Link>
              ))}
            </nav>
          </div>
        </Container>
      </header>
    </>
  );
}
