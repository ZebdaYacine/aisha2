"use client";

import { usePathname } from "next/navigation";
import { ChevronRight } from "lucide-react";

import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { LocalizedLink } from "@/core/components/shared/localized-link";
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
  const pathname = usePathname();
  const { user } = useAuth();
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
  return (
    <aside className="w-full min-w-0">
      <div className="relative min-w-0">
        <div className="horizontal-tab-scroll w-full min-w-0 max-w-full touch-pan-x overflow-x-auto overflow-y-hidden lg:overflow-visible">
          <nav
            aria-label={copy.account}
            className="flex w-max min-w-full snap-x gap-2 border-b border-border pb-3 lg:w-full lg:snap-none lg:flex-col lg:border-b-0 lg:border-e lg:pe-8"
          >
        {customerLinks.map(([label, href]) => (
          <LocalizedLink
            className={`min-h-11 shrink-0 snap-start whitespace-nowrap rounded px-3 py-3 transition-colors hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring ${pathname === `/${locale}${href}` ? "bg-secondary font-medium text-foreground" : ""}`}
            locale={locale}
            href={href}
            key={href}
          >
            {label}
          </LocalizedLink>
        ))}

        {mode === "customer" && !canUseArtisanArea && canApplyForArtisan && (
          <LocalizedLink
            className="min-h-11 shrink-0 snap-start whitespace-nowrap border border-foreground bg-foreground px-4 py-3 text-sm text-background hover:border-primary hover:bg-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:mt-2"
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
            className="min-h-11 shrink-0 snap-start whitespace-nowrap px-3 py-3 hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            locale={locale}
            href="/artisan#workshops"
          >
            {copy.manageWorkshops}
          </LocalizedLink>
        )}
        {mode === "artisan" && !canUseArtisanArea && canApplyForArtisan && (
          <LocalizedLink
            className="min-h-11 shrink-0 snap-start whitespace-nowrap border-t border-border px-3 py-3 text-sm hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:mt-2"
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
            className="shrink-0 snap-start border-t border-border px-3 py-3 text-sm text-muted-foreground lg:mt-2"
            role="status"
          >
            {copy.artisanSuspended}
          </span>
        )}

          </nav>
        </div>
        <div className="mt-2 flex h-1.5 w-full overflow-hidden rounded-full bg-muted lg:hidden" aria-hidden="true">
          <span className="h-full w-1/3 rounded-full bg-primary" />
        </div>
        <span className="pointer-events-none absolute end-0 top-0 flex h-14 w-10 items-center justify-end bg-gradient-to-l from-background via-background/90 to-transparent pe-1 text-primary rtl:bg-gradient-to-r lg:hidden" aria-hidden="true">
          <ChevronRight size={16} />
        </span>
      </div>
    </aside>
  );
}
