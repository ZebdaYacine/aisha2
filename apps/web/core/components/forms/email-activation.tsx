"use client";

import { useEffect, useState } from "react";

import { LocalizedLink } from "@/core/components/shared/localized-link";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { userFromAuthResponse } from "@/features/auth/types";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

export function EmailActivation({
  locale,
  copy,
  token,
}: {
  locale: Locale;
  copy: StoreCopy;
  token: string;
}) {
  const auth = useOptionalAuth();
  const [message, setMessage] = useState(token ? copy.activationChecking : copy.activationMissing);
  const [error, setError] = useState(!token);

  useEffect(() => {
    if (!token) {
      return;
    }

    let cancelled = false;
    void fetch(`/api/auth/activate?token=${encodeURIComponent(token)}`)
      .then(async (response) => {
        const result = await response.json().catch(() => ({}));
        if (!response.ok) throw new Error(result?.error?.message ?? copy.activationError);
        if (cancelled) return;
        const user = userFromAuthResponse(result);
        if (user) auth?.setUser(user);
        setMessage(copy.activationSuccess);
        window.setTimeout(() => {
          window.location.assign(`/${locale}/account`);
        }, 500);
      })
      .catch(() => {
        if (!cancelled) {
          setMessage(copy.activationError);
          setError(true);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [auth, copy.activationChecking, copy.activationError, copy.activationMissing, copy.activationSuccess, locale, token]);

  return (
    <section className="mx-auto flex min-h-[60svh] max-w-md flex-col justify-center px-5 py-16">
      <p className="text-xs uppercase tracking-widest text-primary">AISHA</p>
      <h1 className="mt-4 font-serif text-5xl">{copy.activationTitle}</h1>
      <p className={`mt-6 border p-4 text-sm ${error ? "border-destructive/30 bg-destructive/10 text-destructive" : "border-border bg-muted/30"}`} role={error ? "alert" : "status"}>
        {message}
      </p>
      {error && (
        <LocalizedLink locale={locale} href="/login" className="mt-6 underline">
          {copy.login}
        </LocalizedLink>
      )}
    </section>
  );
}
