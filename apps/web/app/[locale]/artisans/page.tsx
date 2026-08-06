import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { ArtisanCard } from "@/core/components/editorial/artisan-card";
import { Container } from "@/core/components/layout/container";
import { Breadcrumbs } from "@/core/components/shared/breadcrumbs";
import { catalogue } from "@/features/catalogue/api";
import { localized } from "@/features/catalogue/format";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }): Promise<Metadata> {
  const { locale } = await params;
  if (!isLocale(locale)) return {};
  const copy = storeCopy(locale);
  return { title: `${copy.artisans} — AISHA`, description: copy.artisanDirectoryBody, alternates: { canonical: `/${locale}/artisans` } };
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
      <form className="mt-10 grid gap-3 border-y border-border py-5 sm:grid-cols-3">
        <label>
          <span className="sr-only">{copy.search}</span>
          <input
            className="h-12 w-full border border-border bg-background px-4"
            placeholder={copy.search}
            type="search"
          />
        </label>
        <label>
          <span className="sr-only">{copy.region}</span>
          <select className="h-12 w-full border border-border bg-background px-4">
            <option>{copy.region}</option>
            {artisans.map((artisan) => <option key={artisan.slug}>{localized(artisan.region, locale)}</option>)}
          </select>
        </label>
        <label>
          <span className="sr-only">{copy.craft}</span>
          <select className="h-12 w-full border border-border bg-background px-4">
            <option>{copy.craft}</option>
            {artisans.map((artisan) => <option key={artisan.slug}>{localized(artisan.craft, locale)}</option>)}
          </select>
        </label>
      </form>
      <div className="mt-12 grid gap-10 md:grid-cols-2 lg:grid-cols-3">
        {artisans.map((artisan) => (
          <ArtisanCard
            artisan={artisan}
            locale={locale}
            copy={copy}
            key={artisan.slug}
          />
        ))}
      </div>
    </Container>
  );
}
