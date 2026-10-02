"use client";

import { useEffect, useState } from "react";

import { LocalizedLink } from "@/core/components/shared/localized-link";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";

export function EmailActivation({
  locale,
  copy,
  token,
}: {
  locale: Locale;
  copy: StoreCopy;
  token: string;
}) {
  const [state, setState] = useState<"checking" | "success" | "expired" | "missing">(token ? "checking" : "missing");

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
        setState("success");
      })
      .catch(() => {
        if (!cancelled) {
          setState("expired");
        }
      });

    return () => {
      cancelled = true;
    };
  }, [copy.activationError, token]);

  const message = state === "checking"
    ? copy.activationChecking
    : state === "success"
      ? copy.activationSuccess
      : state === "missing"
        ? copy.activationMissing
        : copy.activationExpired;
  const error = state === "expired" || state === "missing";
  const success = state === "success";

  return (
    <section className="mx-auto flex min-h-[60svh] max-w-md flex-col justify-center px-5 py-16">
      <p className="text-xs uppercase tracking-widest text-primary">AISHA</p>
      <h1 className="mt-4 font-serif text-5xl">{copy.activationTitle}</h1>
      <p className={`mt-6 border p-4 text-sm ${error ? "border-destructive/30 bg-destructive/10 text-destructive" : success ? "border-emerald-300 bg-emerald-50 text-emerald-800" : "border-border bg-muted/30"}`} role={error ? "alert" : "status"}>
        {message}
      </p>
      {state === "success" && (
        <LocalizedLink locale={locale} href="/login" className="mt-6 inline-flex w-fit rounded-md bg-primary px-5 py-3 text-sm text-primary-foreground">
          {copy.login}
        </LocalizedLink>
      )}
      {error && (
        <LocalizedLink locale={locale} href="/login" className="mt-6 underline">
          {copy.login}
        </LocalizedLink>
      )}
    </section>
  );
}
