/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { ErrorState } from "@/core/components/feedback/error-state";
import { FormErrorSummary } from "@/core/components/forms/form-error-summary";
import { Button } from "@/core/components/ui/button";
import { submitJSON } from "@/core/forms/api";
import { formCopy } from "@/core/lib/form-copy";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

type Values = { fullName: string; email: string; phone: string };
type LoadState = "loading" | "ready" | "error";

export function ProfileForm({
  locale,
  copy,
}: {
  locale: Locale;
  copy: StoreCopy;
}) {
  const auth = useOptionalAuth();
  const router = useRouter();
  const messages = formCopy(locale);
  const schema = z.object({
    fullName: z.string().trim().min(2, messages.name),
    email: z.email(messages.email),
    phone: z
      .string()
      .trim()
      .refine((value) => value === "" || value.length >= 6, messages.phone),
  });
  const [summary, setSummary] = useState<string>();
  const [loadState, setLoadState] = useState<LoadState>("loading");
  const {
    register,
    handleSubmit,
    setError,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: { fullName: "", email: "", phone: "" },
  });

  const loadProfile = useCallback(async () => {
    setLoadState("loading");
    setSummary(undefined);
    try {
      const response = await fetch("/api/customer/profile", {
        cache: "no-store",
      });
      if (response.status === 401 || response.status === 403) {
        await auth?.logout();
        router.replace(`/${locale}/login`);
        return;
      }
      if (!response.ok) throw new Error("profile request failed");
      const profile = await response.json();
      reset({
        fullName: profile.displayName ?? "",
        email: profile.email ?? "",
        phone: profile.phone ?? "",
      });
      setLoadState("ready");
    } catch {
      setLoadState("error");
    }
  }, [auth, locale, reset, router]);

  useEffect(() => {
    void loadProfile();
  }, [loadProfile]);

  const valid = async (values: Values) => {
    const result = await submitJSON(
      "/api/customer/profile",
      { displayName: values.fullName, phone: values.phone },
      "PATCH",
    );
    if (result.ok) {
      auth?.updateUser({ displayName: values.fullName });
      setSummary(undefined);
      toast.success(copy.profile);
      return;
    }
    if (result.status === 401 || result.status === 403) {
      await auth?.logout();
      router.replace(`/${locale}/login`);
      return;
    }
    Object.entries(result.fieldErrors ?? {}).forEach(([name, message]) => {
      setError(name === "displayName" ? "fullName" : (name as keyof Values), {
        message,
      });
    });
    setSummary(
      result.code === "SERVICE_UNAVAILABLE"
        ? messages.unavailable
        : messages.backendError,
    );
    toast.error(messages.backendError);
  };

  if (loadState === "loading")
    return (
      <p
        className="mt-10 border border-border p-6 text-sm text-muted-foreground"
        role="status"
      >
        {copy.accountLoading}
      </p>
    );
  if (loadState === "error")
    return (
      <ErrorState
        title={copy.accountUnavailable}
        body={copy.accountErrorBody}
        action={
          <Button
            variant="outline"
            type="button"
            onClick={() => void loadProfile()}
          >
            {copy.retry}
          </Button>
        }
      />
    );

  return (
    <form
      onSubmit={handleSubmit(valid, () => setSummary(messages.formInvalid))}
      className="mt-10 space-y-5"
      noValidate
    >
      <FormErrorSummary message={summary} id="profile-errors" />
      <Field
        label={copy.fullName}
        error={errors.fullName?.message}
        errorId="profile-full-name-error"
      >
        <input
          {...register("fullName")}
          className="auth-input"
          aria-invalid={!!errors.fullName}
          aria-describedby={
            errors.fullName ? "profile-full-name-error" : undefined
          }
        />
      </Field>
      <Field
        label={copy.email}
        error={errors.email?.message}
        errorId="profile-email-error"
      >
        <input
          {...register("email")}
          className="auth-input"
          type="email"
          readOnly
          aria-invalid={!!errors.email}
          aria-describedby={errors.email ? "profile-email-error" : undefined}
        />
      </Field>
      <Field
        label={copy.phone}
        error={errors.phone?.message}
        errorId="profile-phone-error"
      >
        <input
          {...register("phone")}
          className="auth-input"
          type="tel"
          aria-invalid={!!errors.phone}
          aria-describedby={errors.phone ? "profile-phone-error" : undefined}
        />
      </Field>
      <Button disabled={isSubmitting} type="submit">
        {copy.profile}
      </Button>
    </form>
  );
}

function Field({
  label,
  error,
  errorId,
  children,
}: {
  label: string;
  error?: string;
  errorId?: string;
  children: React.ReactNode;
}) {
  return (
    <label>
      <span className="mb-2 block text-sm">{label}</span>
      {children}
      {error && (
        <span
          id={errorId}
          className="mt-1 block border-b border-destructive pb-1 text-xs text-destructive"
          role="alert"
        >
          {error}
        </span>
      )}
    </label>
  );
}
