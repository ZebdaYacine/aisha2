"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { LockKeyhole } from "lucide-react";
import { useEffect, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import type { FormSubmitter } from "@/core/forms/types";
import { formCopy } from "@/core/lib/form-copy";
import { countryOptions } from "@/core/lib/location-options";
import { useCart } from "@/features/cart/viewmodel/cart-context";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { FormErrorSummary } from "./form-error-summary";

type Values = {
  email: string;
  fullName: string;
  phone: string;
  line1: string;
  line2?: string;
  city: string;
  postalCode: string;
  country: string;
};
export function CheckoutForm({
  copy,
  locale,
  submit,
}: {
  copy: StoreCopy;
  locale: Locale;
  submit?: FormSubmitter<Values>;
}) {
  const messages = formCopy(locale);
  const schema = z.object({
    email: z.email(messages.email),
    fullName: z.string().trim().min(2, messages.name),
    phone: z.string().trim().min(6, messages.phone),
    line1: z.string().trim().min(4, messages.address),
    line2: z.string().optional(),
    city: z.string().trim().min(2, messages.required),
    postalCode: z.string().trim().min(2, messages.required),
    country: z.string().trim().min(2, messages.required),
  });
  const summaryId = "checkout-form-errors";
  const [summary, setSummary] = useState<string>();
  useEffect(() => {
    if (summary) document.getElementById(summaryId)?.focus();
  }, [summary]);
  const {
    register,
    control,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<Values>({
    resolver: zodResolver(schema),
    shouldFocusError: false,
  });
  const cart = useCart();
  const submitCheckout: FormSubmitter<Values> =
    submit ??
    (async (values) => {
      const address = await fetch("/api/customer/addresses", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          fullName: values.fullName,
          phone: values.phone,
          line1: values.line1,
          line2: values.line2 ?? "",
          city: values.city,
          postalCode: values.postalCode,
          country: values.country,
          isDefault: false,
        }),
      });
      if (!address.ok) return { ok: false, code: "VALIDATION_FAILED" as const };
      const saved = (await address.json()) as { id: string };
      const response = await fetch("/api/checkout", {
        method: "POST",
        headers: {
          "content-type": "application/json",
          "Idempotency-Key": crypto.randomUUID(),
        },
        body: JSON.stringify({
          addressId: saved.id,
          items: cart.items.map((item) => ({
            productId: item.productSlug,
            quantity: item.quantity,
          })),
        }),
      });
      if (!response.ok)
        return {
          ok: false,
          code:
            response.status === 503
              ? ("SERVICE_UNAVAILABLE" as const)
              : ("UNKNOWN" as const),
        };
      const result = (await response.json()) as { order: { id: string } };
      cart.clear();
      window.location.assign(`/${locale}/account/orders/${result.order.id}`);
      return { ok: true };
    });
  const focusSummary = (message: string) => setSummary(message);
  const valid = async (values: Values) => {
    setSummary(undefined);
    try {
      const result = await submitCheckout(values);
      if (result.ok) return;
      Object.entries(result.fieldErrors ?? {}).forEach(([name, message]) =>
        setError(name as keyof Values, { type: "server", message }),
      );
      const message =
        result.code === "SERVICE_UNAVAILABLE"
          ? messages.unavailable
          : messages.backendError;
      focusSummary(message);
      toast.error(message);
    } catch {
      focusSummary(messages.backendError);
      toast.error(messages.backendError);
    }
  };
  const fields: [keyof Values, string, string][] = [
    ["fullName", copy.fullName, "text"],
    ["email", copy.email, "email"],
    ["phone", copy.phone, "tel"],
    ["country", copy.country, "text"],
    ["line1", copy.line1, "text"],
    ["line2", copy.line2, "text"],
    ["city", copy.city, "text"],
    ["postalCode", copy.postalCode, "text"],
  ];
  return (
    <form
      onSubmit={handleSubmit(valid, () => focusSummary(messages.formInvalid))}
      noValidate
    >
      <FormErrorSummary message={summary} id={summaryId} />
      <Step number="01" title={copy.contact} />
      <div className="mt-6 grid gap-5 sm:grid-cols-2">
        {fields.slice(0, 2).map((field) => (
          <Field
            key={field[0]}
            field={field}
            error={errors[field[0]]?.message}
            register={register}
          />
        ))}
      </div>
      <Step number="02" title={copy.address} />
      <div className="mt-6 grid gap-5 sm:grid-cols-2">
        {fields.slice(2).map((field) => (
          field[0] === "country" ? (
            <Controller
              key={field[0]}
              name="country"
              control={control}
              render={({ field: countryField }) => (
                <div>
                  <label htmlFor="checkout-country" className="mb-2 block text-sm">
                    {field[1]}
                  </label>
                  <Combobox
                    options={countryOptions}
                    value={countryField.value ?? ""}
                    onChange={countryField.onChange}
                    id="checkout-country"
                    placeholder={field[1]}
                    ariaLabel={field[1]}
                  />
                  {errors.country?.message && (
                    <span className="mt-1 block border-b border-destructive pb-1 text-xs text-destructive" role="alert">
                      {errors.country.message}
                    </span>
                  )}
                </div>
              )}
            />
          ) : (
            <Field
              key={field[0]}
              field={field}
              error={errors[field[0]]?.message}
              register={register}
            />
          )
        ))}
      </div>
      <Step number="03" title={copy.shippingMethod} />
      <div className="mt-6 border border-border p-5">
        <p className="text-sm">{copy.shippingPending}</p>
      </div>
      <Step number="04" title={copy.paymentMethod} />
      <div className="mt-6 border border-border p-5">
        <p className="flex gap-3 text-sm">
          <LockKeyhole aria-hidden size={18} />
          {copy.paymentPending}
        </p>
      </div>
      <Step number="05" title={copy.review} />
      <Button className="mt-7 w-full" disabled={isSubmitting} type="submit">
        {copy.placeOrder}
      </Button>
    </form>
  );
}
function Step({ number, title }: { number: string; title: string }) {
  return (
    <div className="mt-12 flex items-center gap-4 border-b border-border pb-4 first:mt-0">
      <span className="text-xs text-primary">{number}</span>
      <h2 className="font-serif text-2xl">{title}</h2>
    </div>
  );
}
function Field({
  field: [name, label, type],
  error,
  register,
}: {
  field: [keyof Values, string, string];
  error?: string;
  register: ReturnType<typeof useForm<Values>>["register"];
}) {
  const id = `checkout-${name}`;
  return (
    <div
      className={name === "line1" || name === "line2" ? "sm:col-span-2" : ""}
    >
      <label htmlFor={id} className="mb-2 block text-sm">
        {label}
      </label>
      <input
        id={id}
        {...register(name)}
        type={type}
        aria-invalid={!!error}
        aria-describedby={error ? `${id}-error` : undefined}
        className="auth-input"
      />
      {error && (
        <span
          id={`${id}-error`}
          role="alert"
          className="mt-1 block border-b border-destructive pb-1 text-xs text-destructive"
        >
          {error}
        </span>
      )}
    </div>
  );
}
