"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import { Button } from "@/core/components/ui/button";
import { hasCapability } from "@/features/auth/types";
import { useAuth } from "@/features/auth/viewmodel/auth-context";

type AccountSidebarMode = "customer" | "artisan";

export function AccountSidebar({
  locale,
  copy,
  mode = "customer",
}: {
  locale: Locale;
  copy: StoreCopy;
  mode?: AccountSidebarMode;
}) {
  const router = useRouter();
  const { logout, status, user } = useAuth();
  const [signingOut, setSigningOut] = useState(false);
  const customerLinks = [
    [copy.overview, "/account"],
    [copy.profile, "/account/profile"],
    [copy.addresses, "/account/addresses"],
    [copy.purchases, "/account/orders"],
    [copy.tracking, "/account/tracking"],
    [copy.wishlist, "/account/wishlist"],
  ] as const;
  const canUseArtisanArea =
    user?.artisanEnabled === true &&
    hasCapability(user, "artisan.account.read");
  const canApplyForArtisan =
    user?.artisanStatus === "NOT_STARTED" || user?.artisanStatus === "PENDING";

  const signOut = async () => {
    if (signingOut) return;
    setSigningOut(true);
    try {
      await logout();
    } finally {
      router.replace(`/${locale}/login`);
      router.refresh();
      setSigningOut(false);
    }
  };

  return (
    <aside>
      <nav
        aria-label={copy.account}
        className="flex gap-2 overflow-x-auto border-b border-border pb-4 lg:flex-col lg:border-b-0 lg:border-e lg:pe-8"
      >
        {customerLinks.map(([label, href]) => (
          <LocalizedLink
            className="min-h-11 whitespace-nowrap px-3 py-3 hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            locale={locale}
            href={href}
            key={href}
          >
            {label}
          </LocalizedLink>
        ))}

        {status === "loading" && (
          <span
            className="px-3 py-3 text-sm text-muted-foreground"
            aria-live="polite"
          >
            {copy.account}…
          </span>
        )}
        {mode === "artisan" && (
          <LocalizedLink
            className="min-h-11 whitespace-nowrap border-t border-border px-3 py-3 font-medium hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:mt-2"
            locale={locale}
            href="/account"
          >
            {copy.customerArea}
          </LocalizedLink>
        )}
        {mode === "customer" && canUseArtisanArea && (
          <LocalizedLink
            className="min-h-11 whitespace-nowrap border-t border-border px-3 py-3 font-medium hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:mt-2"
            locale={locale}
            href="/artisan"
          >
            {copy.artisanArea}
          </LocalizedLink>
        )}
        {mode === "customer" && !canUseArtisanArea && canApplyForArtisan && (
          <LocalizedLink
            className="min-h-11 whitespace-nowrap border border-foreground bg-foreground px-4 py-3 text-sm text-background hover:border-primary hover:bg-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:mt-2"
            locale={locale}
            href="/artisan"
          >
            {user?.artisanStatus === "PENDING"
              ? copy.artisanPending
              : copy.becomeArtisan}
          </LocalizedLink>
        )}
        {mode === "artisan" && canUseArtisanArea && (
          <LocalizedLink
            className="min-h-11 whitespace-nowrap px-3 py-3 hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            locale={locale}
            href="/artisan"
          >
            {copy.artisanArea}
          </LocalizedLink>
        )}
        {mode === "artisan" && canUseArtisanArea && (
          <LocalizedLink
            className="min-h-11 whitespace-nowrap px-3 py-3 hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            locale={locale}
            href="/artisan#workshops"
          >
            {copy.manageWorkshops}
          </LocalizedLink>
        )}
        {mode === "artisan" && !canUseArtisanArea && canApplyForArtisan && (
          <LocalizedLink
            className="min-h-11 whitespace-nowrap border-t border-border px-3 py-3 text-sm hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:mt-2"
            locale={locale}
            href="/artisan"
          >
            {user?.artisanStatus === "PENDING"
              ? copy.artisanPending
              : copy.becomeArtisan}
          </LocalizedLink>
        )}
        {!canUseArtisanArea && user?.artisanStatus === "SUSPENDED" && (
          <span
            className="border-t border-border px-3 py-3 text-sm text-muted-foreground lg:mt-2"
            role="status"
          >
            {copy.artisanSuspended}
          </span>
        )}

        <Button
          type="button"
          variant="ghost"
          disabled={signingOut}
          aria-busy={signingOut}
          onClick={() => void signOut()}
        >
          {signingOut ? `${copy.logout}…` : copy.logout}
        </Button>
      </nav>
    </aside>
  );
}
