import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { Container } from "@/core/components/layout/container";
import { Breadcrumbs } from "@/core/components/shared/breadcrumbs";
import { catalogue } from "@/features/catalogue/api";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
import { ArtisanDirectory } from "@/features/catalogue";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  if (!isLocale(locale)) return {};
  const copy = storeCopy(locale);
  return {
    title: `${copy.artisans} — AISHA`,
    description: copy.artisanDirectoryBody,
    alternates: { canonical: `/${locale}/artisans` },
  };
}
export default async function ArtisansPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const copy = storeCopy(locale);
  const { artisans } = await catalogue(locale);
  return (
    <Container className="py-12 lg:py-16">
      <Breadcrumbs locale={locale} items={[{ label: copy.artisans }]} />
      <header className="mt-10 max-w-3xl">
        <p className="text-xs uppercase tracking-[.2em] text-primary">
          {copy.artisanStory}
        </p>
        <h1 className="mt-3 font-serif text-5xl sm:text-6xl">
          {copy.artisanDirectory}
        </h1>
        <p className="mt-5 text-lg leading-8 text-muted-foreground">
          {copy.artisanDirectoryBody}
        </p>
      </header>
      <ArtisanDirectory artisans={artisans} locale={locale} copy={copy} />
    </Container>
  );
}
