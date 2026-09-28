import { EmptyState } from "@/core/components/feedback/empty-state";
import { ButtonLink } from "@/core/components/ui/button";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export default async function NotFound({ params }: { params?: Promise<{ locale?: string }> }) {
  const locale = (await params)?.locale;
  const copy = locale && isLocale(locale) ? storeCopy(locale) : storeCopy("en");
  return <EmptyState title={copy.notFound} body={copy.errorBody} action={<ButtonLink href={`/${locale && isLocale(locale) ? locale : "en"}`}>{copy.returnHome}</ButtonLink>} />;
}
