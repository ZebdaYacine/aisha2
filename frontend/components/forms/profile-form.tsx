"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import { submitJSON } from "@/features/forms/api";
import { formCopy } from "@/lib/form-copy";
import type { Locale } from "@/lib/i18n";
import type { StoreCopy } from "@/lib/store-copy";
import { FormErrorSummary } from "./form-error-summary";
type Values = { fullName: string; email: string; phone: string };
export function ProfileForm({ locale, copy }: { locale: Locale; copy: StoreCopy }) { const messages = formCopy(locale); const schema = z.object({ fullName: z.string().trim().min(2, messages.name), email: z.email(messages.email), phone: z.string().trim().min(6, messages.phone) }); const [summary, setSummary] = useState<string>(); const { register, handleSubmit, setError, reset, formState: { errors, isSubmitting } } = useForm<Values>({ resolver: zodResolver(schema), defaultValues: { fullName: "", email: "", phone: "" } }); useEffect(() => { void fetch("/api/customer/profile").then(async (response) => { if (!response.ok) throw new Error(); const profile = await response.json(); reset({ fullName: profile.displayName, email: profile.email, phone: profile.phone }); }).catch(() => setSummary(messages.backendError)); }, [messages.backendError, reset]); const valid = async (values: Values) => { const result = await submitJSON("/api/customer/profile", { displayName: values.fullName, phone: values.phone }, "PATCH"); if (result.ok) { toast.success(copy.profile); return; } Object.entries(result.fieldErrors ?? {}).forEach(([name, message]) => setError(name === "displayName" ? "fullName" : name as keyof Values, { message })); setSummary(result.code === "SERVICE_UNAVAILABLE" ? messages.unavailable : messages.backendError); toast.error(messages.backendError); }; return <form onSubmit={handleSubmit(valid, () => setSummary(messages.formInvalid))} className="mt-10 space-y-5" noValidate><FormErrorSummary message={summary} id="profile-errors"/><Field label={copy.fullName} error={errors.fullName?.message}><input {...register("fullName")} className="auth-input" aria-invalid={!!errors.fullName}/></Field><Field label={copy.email} error={errors.email?.message}><input {...register("email")} className="auth-input" type="email" readOnly/></Field><Field label={copy.phone} error={errors.phone?.message}><input {...register("phone")} className="auth-input" type="tel" aria-invalid={!!errors.phone}/></Field><Button disabled={isSubmitting} type="submit">{copy.profile}</Button></form>; }
function Field({ label, error, children }: { label: string; error?: string; children: React.ReactNode }) { return <label><span className="mb-2 block text-sm">{label}</span>{children}{error && <span className="mt-1 block border-b border-destructive pb-1 text-xs text-destructive" role="alert">{error}</span>}</label>; }
