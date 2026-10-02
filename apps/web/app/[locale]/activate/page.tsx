import { notFound } from "next/navigation";

import { EmailActivation } from "@/core/components/forms/email-activation";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export default async function ActivatePage({
  params,
  searchParams,
}: {
  params: Promise<{ locale: string }>;
  searchParams: Promise<{ token?: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const { token = "" } = await searchParams;
  return <EmailActivation locale={locale} copy={storeCopy(locale)} token={token} />;
}
