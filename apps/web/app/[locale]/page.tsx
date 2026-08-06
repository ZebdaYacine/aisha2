import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { HomeHero } from "@/core/components/editorial/home-hero";
import { HomeSections } from "@/core/components/editorial/home-sections";
import { dictionary, isLocale } from "@/core/lib/i18n";
import { catalogue } from "@/features/catalogue/api";

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
  const messages = dictionary(locale);
  const data = await catalogue(locale);
  return (
    <>
      <HomeHero locale={locale} messages={messages} />
      <HomeSections locale={locale} messages={messages} {...data} />
    </>
  );
}
