"use client";

import { ButtonLink } from "@/core/components/ui/button";
import { ErrorState } from "@/core/components/feedback/error-state";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { useAuth } from "@/features/auth/viewmodel/auth-context";
import { hasCapability } from "@/features/auth/types";

function artisanStatusLabel(status: string, copy: StoreCopy) {
  if (status === "ACTIVE") return copy.active;
  if (status === "SUSPENDED") return copy.suspended;
  if (status === "PENDING") return copy.pending;
  return copy.notStarted;
}

export function AccountOverview({ locale, copy }: { locale: Locale; copy: StoreCopy }) {
  const { status, user } = useAuth();

  if (status === "loading") {
    return <p className="mt-8 border border-border p-6 text-sm text-muted-foreground" role="status">{copy.accountLoading}</p>;
  }

  if (status !== "authenticated" || !user) {
    return (
      <ErrorState
        title={copy.accountUnavailable}
        body={copy.accountErrorBody}
        action={<ButtonLink href={`/${locale}/login`}>{copy.login}</ButtonLink>}
      />
    );
  }

  const canUseArtisanArea = user.artisanEnabled && hasCapability(user, "artisan.account.read");
  const accountStatus = user.status === "SUSPENDED" ? copy.suspended : copy.active;

  return (
    <div className="mt-8 grid gap-4 md:grid-cols-3" aria-label={copy.accountStatus}>
      <section className="border border-border p-6">
        <h2 className="font-serif text-2xl">{copy.accountStatus}</h2>
        <p className="mt-4 text-sm font-medium">{accountStatus}</p>
        <p className="mt-2 break-words text-sm text-muted-foreground">{user.email}</p>
      </section>
      <section className="border border-border p-6">
        <h2 className="font-serif text-2xl">{copy.customerAccess}</h2>
        <p className="mt-4 text-sm font-medium">{user.customerEnabled ? copy.enabled : copy.disabled}</p>
        <LocalizedLink className="mt-5 inline-flex border-b border-foreground pb-1 text-sm" locale={locale} href="/account/profile">
          {copy.profile}
        </LocalizedLink>
      </section>
      <section className="border border-border p-6">
        <h2 className="font-serif text-2xl">{copy.artisanAccess}</h2>
        <p className="mt-4 text-sm font-medium">{artisanStatusLabel(user.artisanStatus, copy)}</p>
        {canUseArtisanArea && (
          <ButtonLink className="mt-5" href={`/${locale}/artisan`} variant="outline">
            {copy.artisanArea}
          </ButtonLink>
        )}
      </section>
    </div>
  );
}
