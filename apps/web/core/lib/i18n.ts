import ar from "@/messages/ar.json";
import en from "@/messages/en.json";
import es from "@/messages/es.json";
import fr from "@/messages/fr.json";

export const locales = ["en", "fr", "ar", "es"] as const;
export type Locale = (typeof locales)[number];

const messages = { ar, en, es, fr };
export type Messages = typeof en;

export function isLocale(value: string): value is Locale {
  return locales.includes(value as Locale);
}

export function dictionary(locale: Locale) {
  return messages[locale];
}

export function direction(locale: Locale) {
  return locale === "ar" ? "rtl" : "ltr";
}
