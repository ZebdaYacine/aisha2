import { BadgeCheck, MapPin, PackageCheck, Share2 } from "lucide-react";
import type { Metadata } from "next";
import Image from "next/image";
import { notFound } from "next/navigation";
import { ProductGallery, ProductGrid, ProductPrice, PurchaseControls } from "@/features/product";
import { Container } from "@/core/components/layout/container";
import { Breadcrumbs } from "@/core/components/shared/breadcrumbs";
import { SectionHeading } from "@/core/components/shared/section-heading";
import { ButtonLink } from "@/core/components/ui/button";
import { catalogue, catalogueProduct } from "@/features/catalogue/api";
import { localized } from "@/features/catalogue/format";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export async function generateMetadata({ params }: { params: Promise<{ locale: string; slug: string }> }): Promise<Metadata> {
  const { locale, slug } = await params;
  if (!isLocale(locale)) return {};
  const product = await catalogueProduct(slug, locale).catch(() => undefined);
  if (!product) return {};
  return { title: `${localized(product.name, locale)} — AISHA`, description: localized(product.summary, locale), alternates: { canonical: `/${locale}/products/${slug}` } };
}
export default async function ProductPage({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { locale, slug } = await params;
  if (!isLocale(locale)) notFound();
  const product = await catalogueProduct(slug, locale).catch(() => undefined);
  if (!product) notFound();
  const data = await catalogue(locale);
  const copy = storeCopy(locale),
    artisan = data.artisans.find((item) => item.slug === product.artisanSlug),
    category = data.categories.find((item) => item.slug === product.categorySlug);
  return (
    <>
      <Container className="py-8">
        <Breadcrumbs
          locale={locale}
          items={[
            { label: copy.products, href: "/products" },
            { label: localized(product.name, locale) },
          ]}
        />
        <div className="mt-8 grid gap-10 lg:grid-cols-[1.35fr_.85fr] lg:gap-16">
          <ProductGallery
            images={product.images}
            alt={localized(product.name, locale)}
            previous={copy.previous}
            next={copy.next}
            zoomLabel={copy.zoomImage}
            closeLabel={copy.close}
          />
          <aside className="lg:sticky lg:top-28 lg:self-start">
            <p className="text-xs uppercase tracking-[.16em] text-primary">
              {localized(product.region, locale)} ·{" "}
              {category && localized(category.name, locale)}
            </p>
            <h1 className="mt-4 font-serif text-4xl leading-tight sm:text-5xl">
              {localized(product.name, locale)}
            </h1>
            {artisan && (
              <ButtonLink
                variant="ghost"
                className="mt-3 h-auto min-h-0 border-0 p-0 text-muted-foreground"
                href={`/${locale}/artisans/${artisan.slug}`}
              >
                {copy.artisan}: {artisan.name}
                <BadgeCheck size={16} />
              </ButtonLink>
            )}
            <div className="mt-5 flex items-center gap-3 text-sm">
              <span aria-label={`${product.rating} ${copy.ratingLabel}`}>
                ★ {product.rating}
              </span>
              <span className="text-muted-foreground">
                ({product.reviewCount})
              </span>
            </div>
            <ProductPrice
              {...product}
              locale={locale}
              className="mt-6 text-xl"
            />
            <p className="mt-6 leading-7 text-muted-foreground">
              {localized(product.summary, locale)}
            </p>
            <div className="mt-7 border-y border-border py-4">
              <p className="flex items-center gap-2 text-sm">
                <PackageCheck size={17} />
                {product.availability === "in_stock"
                  ? copy.inStock
                  : product.availability === "low_stock"
                    ? copy.lowStock
                    : product.availability === "made_to_order"
                      ? copy.madeToOrder
                      : copy.outOfStock}
              </p>
            </div>
            <div className="mt-7">
              <PurchaseControls
                slug={product.slug}
                availability={product.availability}
                copy={copy}
              />
            </div>
            <div className="mt-7 space-y-3 text-sm text-muted-foreground">
              <p className="flex items-center gap-2">
                <BadgeCheck size={17} />
                {copy.quality}
              </p>
              <p className="flex items-center gap-2">
                <MapPin size={17} />
                {copy.origin}
              </p>
              <p>{copy.delivery}</p>
            </div>
            <button
              className="mt-7 flex min-h-11 items-center gap-2 text-sm"
              type="button"
            >
              <Share2 aria-hidden="true" size={16} /> {copy.share}
            </button>
          </aside>
        </div>
      </Container>
      <section className="section-space bg-muted">
        <Container className="grid gap-12 lg:grid-cols-[1.2fr_.8fr]">
          <div>
            <p className="text-xs uppercase tracking-[.2em] text-primary">
              {copy.story}
            </p>
            <h2 className="mt-4 font-serif text-4xl">
              {localized(product.name, locale)}
            </h2>
            <p className="mt-6 max-w-2xl text-lg leading-8 text-muted-foreground">
              {localized(product.story, locale)}
            </p>
          </div>
          <dl className="divide-y divide-border border-y border-border">
            <Detail term={copy.materials}>
              {localized(product.materials, locale)}
            </Detail>
            <Detail term={copy.method}>
              {localized(product.method, locale)}
            </Detail>
            <Detail term={copy.region}>
              {localized(product.region, locale)}
            </Detail>
            <Detail term={copy.shipping}>{copy.delivery}</Detail>
          </dl>
        </Container>
      </section>
      {artisan && (
        <section className="section-space">
          <Container className="grid items-center gap-10 lg:grid-cols-2">
            <div className="relative aspect-[4/3] overflow-hidden bg-secondary"><Image src={artisan.image} alt={`${artisan.name}, ${artisan.workshop}`} fill sizes="(max-width:1024px) 100vw,50vw" className="object-cover" /></div>
            <div>
              <p className="text-xs uppercase tracking-widest text-primary">
                {copy.verified}
              </p>
              <h2 className="mt-4 font-serif text-4xl">{artisan.name}</h2>
              <p className="mt-2 text-muted-foreground">
                {artisan.workshop} · {localized(artisan.region, locale)}
              </p>
              <p className="mt-6 leading-7 text-muted-foreground">
                {localized(artisan.biography, locale)}
              </p>
              <ButtonLink
                className="mt-7"
                variant="outline"
                href={`/${locale}/artisans/${artisan.slug}`}
              >
                {copy.exploreArtisan}
              </ButtonLink>
            </div>
          </Container>
        </section>
      )}
      <section className="section-space bg-muted">
        <Container>
          <SectionHeading title={copy.related} />
          <div className="mt-10">
            <ProductGrid
              products={data.products
                .filter((p) => p.slug !== product.slug)
                .slice(0, 4)}
              locale={locale}
              copy={copy}
            />
          </div>
        </Container>
      </section>
    </>
  );
}
function Detail({
  term,
  children,
}: {
  term: string;
  children: React.ReactNode;
}) {
  return (
    <div className="grid grid-cols-2 gap-4 py-5">
      <dt className="text-sm text-muted-foreground">{term}</dt>
      <dd className="text-sm">{children}</dd>
    </div>
  );
}
