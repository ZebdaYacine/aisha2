"use client";

import { usePathname, useRouter } from "next/navigation";

import { Combobox } from "@/core/components/ui/combobox";
import type { Locale, Messages } from "@/core/lib/i18n";
import { locales } from "@/core/lib/i18n";

const flags: Record<Locale, string> = { en: "🇬🇧", fr: "🇫🇷", ar: "🇩🇿", es: "🇪🇸" };

export function pathForLocale(pathname: string, nextLocale: Locale) {
  const segments = pathname.split("/");
  if (locales.includes(segments[1] as Locale)) segments[1] = nextLocale;
  else segments.splice(1, 0, nextLocale);
  return segments.join("/") || `/${nextLocale}`;
}

export function LanguageSwitcher({ locale, messages, compact = false }: { locale: Locale; messages: Messages; compact?: boolean }) {
  const router = useRouter();
  const pathname = usePathname();
  const options = locales.map((item) => ({ value: item, label: `${flags[item]} ${item.toUpperCase()} · ${messages.locales[item]}` }));

  const changeLocale = (next: string) => {
    if (!locales.includes(next as Locale)) return;
    const destination = pathForLocale(pathname || `/${locale}`, next as Locale);
    router.replace(`${destination}${window.location.search}${window.location.hash}`);
  };

  return (
    <Combobox
      key={locale}
      className={compact ? "w-[4.75rem] sm:w-24" : "w-full max-w-52"}
      options={options}
      value={locale}
      ariaLabel={messages.languageLabel}
      displayValue={(option) => `${flags[option.value as Locale]} ${option.value.toUpperCase()}`}
      onChange={changeLocale}
    />
  );
}
