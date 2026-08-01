import { BadgeCheck, MapPin } from "lucide-react";
import type { Metadata } from "next";
import Image from "next/image";
import { notFound } from "next/navigation";
import { ProductGrid } from "@/components/commerce/product-grid";
import { Container } from "@/components/layout/container";
import { Breadcrumbs } from "@/components/shared/breadcrumbs";
import { SectionHeading } from "@/components/shared/section-heading";
import { catalogue, catalogueArtisan } from "@/features/storefront/api";
import { localized } from "@/features/storefront/format";
import { isLocale } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";

export async function generateMetadata({ params }: { params: Promise<{ locale: string; slug: string }> }): Promise<Metadata> {
  const { locale, slug } = await params;
  if (!isLocale(locale)) return {};
  const artisan = await catalogueArtisan(slug, locale).catch(() => undefined);
  if (!artisan) return {};
  return { title: `${artisan.name} — AISHA`, description: localized(artisan.biography, locale), alternates: { canonical: `/${locale}/artisans/${slug}` } };
}
export default async function ArtisanPage({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { locale, slug } = await params;
  if (!isLocale(locale)) notFound();
  const artisan = await catalogueArtisan(slug, locale).catch(() => undefined);
  if (!artisan) notFound();
  const copy = storeCopy(locale),
    collection = (await catalogue(locale)).products.filter((p) => p.artisanSlug === slug);
  return (
    <>
      <Container className="py-8">
        <Breadcrumbs
          locale={locale}
          items={[
            { label: copy.artisans, href: "/artisans" },
            { label: artisan.name },
          ]}
        />
        <div className="relative mt-8 min-h-[28rem] overflow-hidden lg:min-h-[38rem]">
          <Image
            src={artisan.image}
            alt={`${artisan.name} — ${artisan.workshop}`}
            fill
            priority
            sizes="100vw"
            className="object-cover"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-foreground/85 via-transparent to-transparent" />
          <div className="absolute inset-x-6 bottom-7 text-background sm:inset-x-10 sm:bottom-10">
            <p className="flex items-center gap-2 text-sm">
              <MapPin size={16} />
              {localized(artisan.region, locale)}
            </p>
            <h1 className="mt-3 font-serif text-5xl sm:text-7xl">
              {artisan.name}
            </h1>
            <p className="mt-2 flex items-center gap-2">
              {artisan.workshop}
              {artisan.verified && (
                <BadgeCheck aria-label={copy.verified} size={18} />
              )}
            </p>
          </div>
        </div>
      </Container>
      <section className="section-space">
        <Container className="grid gap-12 lg:grid-cols-[.7fr_1.3fr]">
          <div>
            <p className="text-xs uppercase tracking-widest text-primary">
              {copy.craft}
            </p>
            <p className="mt-3 font-serif text-3xl">
              {localized(artisan.craft, locale)}
            </p>
          </div>
          <div>
            <p className="max-w-3xl font-serif text-3xl leading-relaxed">
              {localized(artisan.biography, locale)}
            </p>
            <p className="mt-7 max-w-2xl leading-7 text-muted-foreground">
              {copy.artisanDirectoryBody}
            </p>
          </div>
        </Container>
      </section>
      <section className="section-space bg-muted">
        <Container>
          <SectionHeading eyebrow={copy.workshop} title={copy.collection} />
          <div className="mt-12">
            {collection.length ? (
              <ProductGrid products={collection} locale={locale} copy={copy} />
            ) : (
              <p>{copy.noResults}</p>
            )}
          </div>
        </Container>
      </section>
      <section className="section-space">
        <Container>
          <SectionHeading title={copy.workshop} />
          <div className="mt-10 grid grid-cols-2 gap-3 md:grid-cols-3">
            {[0, 1, 2].map((i) => (
              <div
                className="relative aspect-[4/3] overflow-hidden bg-muted"
                key={i}
              >
                <Image
                  src={artisan.image}
                  alt=""
                  fill
                  sizes="33vw"
                  className="object-cover"
                  style={{ objectPosition: `${30 + i * 20}% center` }}
                />
              </div>
            ))}
          </div>
        </Container>
      </section>
    </>
  );
}
