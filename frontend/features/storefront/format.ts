import type { Locale } from "@/lib/i18n";
export function formatMoney(minor: number, currency: string, locale: Locale) { return new Intl.NumberFormat(locale, { style: "currency", currency }).format(minor / 100); }
export function localized<T extends Record<Locale, string>>(value: T, locale: Locale) { return value[locale]; }
