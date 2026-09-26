import type { Locale } from "@/core/lib/i18n";

export function formatFullDate(value: string | Date, locale: Locale) {
  return new Intl.DateTimeFormat(locale, {
    dateStyle: "long",
  }).format(new Date(value));
}

export function formatFullDateTime(value: string | Date, locale: Locale) {
  return new Intl.DateTimeFormat(locale, {
    dateStyle: "long",
    timeStyle: "short",
  }).format(new Date(value));
}
