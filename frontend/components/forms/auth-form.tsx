"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { LocalizedLink } from "@/components/shared/localized-link";
import { Button } from "@/components/ui/button";
import type { FormSubmitter } from "@/features/forms/types";
import { submitJSON } from "@/features/forms/api";
import { formCopy } from "@/lib/form-copy";
import type { Locale } from "@/lib/i18n";
import type { StoreCopy } from "@/lib/store-copy";
import { FormErrorSummary } from "./form-error-summary";

type Mode = "login" | "register" | "forgot";
type Values = {
  email: string;
  password?: string;
  fullName?: string;
  confirmPassword?: string;
  terms?: boolean;
};

export function AuthForm({
  mode,
  locale,
  copy,
  submit,
}: {
  mode: Mode;
  locale: Locale;
  copy: StoreCopy;
  submit?: FormSubmitter<Values>;
}) {
  const messages = formCopy(locale);
  const schema = z
    .object({
      email: z.email(messages.email),
      password: z.string().optional(),
      fullName: z.string().optional(),
      confirmPassword: z.string().optional(),
      terms: z.boolean().optional(),
    })
    .superRefine((value, context) => {
      if (mode !== "forgot" && (!value.password || value.password.length < 12))
        context.addIssue({
          code: "custom",
          path: ["password"],
          message: messages.password,
        });
      if (mode === "register") {
        if (!value.fullName || value.fullName.trim().length < 2)
          context.addIssue({
            code: "custom",
            path: ["fullName"],
            message: messages.name,
          });
        if (value.password !== value.confirmPassword)
          context.addIssue({
            code: "custom",
            path: ["confirmPassword"],
            message: messages.passwordMatch,
          });
        if (!value.terms)
          context.addIssue({
            code: "custom",
            path: ["terms"],
            message: messages.terms,
          });
      }
    });
  const summaryId = `auth-${mode}-errors`;
  const [summary, setSummary] = useState<string>();
  useEffect(() => { if (summary) document.getElementById(summaryId)?.focus(); }, [summary, summaryId]);
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<Values>({
    resolver: zodResolver(schema),
    shouldFocusError: false,
    defaultValues: {
      email: "",
      password: "",
      fullName: "",
      confirmPassword: "",
      terms: false,
    },
  });
  const focusSummary = (message: string) => {
    setSummary(message);
  };
  const onInvalid = () => focusSummary(messages.formInvalid);
  const onValid = async (values: Values) => {
    setSummary(undefined);
    try {
      const result = await (submit ?? ((input) => submitJSON(`/api/auth/${mode === "forgot" ? "forgot-password" : mode}`, input)))(values);
      if (result.ok) {
        toast.success(mode === "forgot" ? copy.forgot : copy.account);
        if (mode !== "forgot") window.location.assign(`/${locale}/account`);
        return;
      }
      Object.entries(result.fieldErrors ?? {}).forEach(([name, message]) =>
        setError(name as keyof Values, { type: "server", message }),
      );
      const message =
        result.code === "SERVICE_UNAVAILABLE"
          ? messages.unavailable
          : result.code === "RATE_LIMITED"
            ? messages.rateLimited
          : messages.backendError;
      focusSummary(message);
      toast.error(message);
    } catch {
      focusSummary(messages.backendError);
      toast.error(messages.backendError);
    }
  };
  return (
    <form
      onSubmit={handleSubmit(onValid, onInvalid)}
      className="space-y-5"
      noValidate
    >
      <FormErrorSummary message={summary} id={summaryId} />
      {mode === "register" && (
        <Field
          id="full-name"
          label={copy.fullName}
          error={errors.fullName?.message}
        >
          <input
            id="full-name"
            {...register("fullName")}
            className="auth-input"
            aria-invalid={!!errors.fullName}
            aria-describedby={errors.fullName ? "full-name-error" : undefined}
            autoComplete="name"
          />
        </Field>
      )}
      <Field id="email" label={copy.email} error={errors.email?.message}>
        <input
          id="email"
          type="email"
          autoComplete="email"
          {...register("email")}
          className="auth-input"
          aria-invalid={!!errors.email}
          aria-describedby={errors.email ? "email-error" : undefined}
        />
      </Field>
      {mode !== "forgot" && (
        <Field
          id="password"
          label={copy.password}
          error={errors.password?.message}
        >
          <input
            id="password"
            type="password"
            autoComplete={
              mode === "login" ? "current-password" : "new-password"
            }
            {...register("password")}
            className="auth-input"
            aria-invalid={!!errors.password}
            aria-describedby={errors.password ? "password-error" : undefined}
          />
        </Field>
      )}
      {mode === "register" && (
        <>
          <Field
            id="confirm-password"
            label={copy.confirmPassword}
            error={errors.confirmPassword?.message}
          >
            <input
              id="confirm-password"
              type="password"
              autoComplete="new-password"
              {...register("confirmPassword")}
              className="auth-input"
              aria-invalid={!!errors.confirmPassword}
              aria-describedby={
                errors.confirmPassword ? "confirm-password-error" : undefined
              }
            />
          </Field>
          <div>
            <label className="flex min-h-11 items-center gap-3 text-sm">
              <input
                type="checkbox"
                {...register("terms")}
                aria-invalid={!!errors.terms}
              />
              {copy.terms}
            </label>
            {errors.terms && (
              <span
                id="terms-error"
                className="block border-b border-destructive pb-1 text-xs text-destructive"
                role="alert"
              >
                {errors.terms.message}
              </span>
            )}
          </div>
        </>
      )}
      {mode === "login" && (
        <LocalizedLink
          className="block text-end text-sm underline"
          locale={locale}
          href="/forgot-password"
        >
          {copy.forgot}
        </LocalizedLink>
      )}
      <Button disabled={isSubmitting} className="w-full" type="submit">
        {mode === "login"
          ? copy.login
          : mode === "register"
            ? copy.register
            : copy.forgot}
      </Button>
      {mode !== "forgot" ? (
        <p className="text-center text-sm text-muted-foreground">
          {mode === "login" ? copy.noAccount : copy.hasAccount}{" "}
          <LocalizedLink
            locale={locale}
            href={mode === "login" ? "/register" : "/login"}
            className="text-foreground underline"
          >
            {mode === "login" ? copy.register : copy.login}
          </LocalizedLink>
        </p>
      ) : (
        <p className="text-center text-sm">
          <LocalizedLink locale={locale} href="/login" className="underline">
            {copy.login}
          </LocalizedLink>
        </p>
      )}
    </form>
  );
}

function Field({
  id,
  label,
  error,
  children,
}: {
  id: string;
  label: string;
  error?: string;
  children: React.ReactNode;
}) {
  return (
    <div>
      <label htmlFor={id} className="mb-2 block text-sm">
        {label}
      </label>
      {children}
      {error && (
        <span
          id={`${id}-error`}
          className="mt-1 block border-b border-destructive pb-1 text-xs text-destructive"
          role="alert"
        >
          {error}
        </span>
      )}
    </div>
  );
}
