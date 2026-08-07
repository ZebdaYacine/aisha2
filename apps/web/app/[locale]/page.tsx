import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { HomeView } from "@/features/home";
import { dictionary, isLocale } from "@/core/lib/i18n";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  if (!isLocale(locale)) {
    return {};
  }
  const messages = dictionary(locale);
  return {
    title: `${messages.brand} — ${messages.home.heroEyebrow}`,
    description: messages.home.heroBody,
    alternates: {
      canonical: `/${locale}`,
      languages: { en: "/en", fr: "/fr", ar: "/ar", es: "/es" },
    },
  };
}

export default async function Home({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  if (!isLocale(locale)) {
    notFound();
  }
  return <HomeView locale={locale} />;
}
