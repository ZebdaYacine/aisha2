import {
  ArrowUpRight,
  Earth,
  HandHeart,
  PackageCheck,
  ShieldCheck,
} from "lucide-react";
import Image from "next/image";
import type { Artisan, Category, Product } from "@/features/storefront/types";
import { localized } from "@/features/storefront/format";
import type { Locale, Messages } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";
import { ArtisanCard } from "./artisan-card";
import { CategoryCard } from "./category-card";
import { ProductGrid } from "@/components/commerce/product-grid";
import { Container } from "@/components/layout/container";
import { LocalizedLink } from "@/components/shared/localized-link";
import { SectionHeading } from "@/components/shared/section-heading";
import { ButtonLink } from "@/components/ui/button";
export function HomeSections({
  locale,
  messages,
  categories,
  products,
  artisans,
}: {
  locale: Locale;
  messages: Messages;
  categories: Category[];
  products: Product[];
  artisans: Artisan[];
}) {
  const copy = storeCopy(locale);
  const trust = [
    [Earth, copy.trustOrigin],
    [PackageCheck, copy.trustQuality],
    [HandHeart, copy.trustPartnership],
    [ShieldCheck, copy.trustDelivery],
  ] as const;
  const regionalProducts = products.filter((product, index, items) => items.findIndex((candidate) => localized(candidate.region, locale) === localized(product.region, locale)) === index).slice(0, 4);
  return (
    <>
      <section className="section-space" id="categories">
        <Container>
          <SectionHeading
            eyebrow={messages.home.categoriesEyebrow}
            title={messages.home.categoriesTitle}
          />
          <div className="mt-10 grid grid-cols-2 gap-2 md:grid-cols-3 lg:grid-cols-6">
            {categories.map((category, index) => (
              <CategoryCard
                category={category}
                locale={locale}
                index={index}
                key={category.slug}
              />
            ))}
          </div>
        </Container>
      </section>
      <section className="section-space bg-muted" id="collection">
        <Container>
          <SectionHeading
            eyebrow={messages.home.collectionEyebrow}
            title={copy.featured}
            body={copy.featuredBody}
            action={
              <LocalizedLink
                className="border-b border-foreground pb-1 text-sm"
                locale={locale}
                href="/products"
              >
                {copy.viewAll}
              </LocalizedLink>
            }
          />
          <div className="mt-12">
            <ProductGrid
              products={products.slice(0, 4)}
              locale={locale}
              copy={copy}
            />
          </div>
        </Container>
      </section>
      <section className="section-space" id="regions">
        <Container>
          <SectionHeading eyebrow={copy.region} title={copy.discoverRegions} body={copy.featuredBody} />
          <div className="mt-12 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            {regionalProducts.map((product) => <LocalizedLink locale={locale} href={`/products/${product.slug}`} key={product.slug} className="group relative min-h-80 overflow-hidden bg-muted">
              <Image src={product.images[0]} alt="" fill sizes="(max-width:640px) 100vw,(max-width:1024px) 50vw,25vw" className="object-cover transition duration-500 group-hover:scale-[1.025]" />
              <span className="absolute inset-0 bg-gradient-to-t from-foreground/80 via-transparent to-transparent" />
              <span className="absolute inset-x-5 bottom-5 font-serif text-2xl text-background">{localized(product.region, locale)}</span>
            </LocalizedLink>)}
          </div>
        </Container>
      </section>
      {artisans[0] && <section
        className="grid bg-foreground text-background lg:grid-cols-2"
        id="artisans"
      >
        <div className="relative min-h-[28rem] lg:min-h-[42rem]">
          <Image
            src={artisans[0].image}
            alt={`${artisans[0].name}, ${artisans[0].workshop}`}
            fill
            sizes="(max-width:1024px) 100vw,50vw"
            className="object-cover opacity-85"
          />
        </div>
        <div className="flex items-center px-6 py-16 sm:px-12 lg:px-20">
          <div className="max-w-xl">
            <p className="text-xs uppercase tracking-[.2em] text-accent">
              {copy.artisanStory}
            </p>
            <h2 className="mt-5 font-serif text-5xl leading-tight">
              {artisans[0].name}
            </h2>
            <p className="mt-2 text-background/65">
              {localized(artisans[0].region, locale)} ·{" "}
              {localized(artisans[0].craft, locale)}
            </p>
            <blockquote className="mt-8 font-serif text-2xl leading-relaxed">
              “{localized(artisans[0].biography, locale)}”
            </blockquote>
            <ButtonLink
              className="mt-9 border-background text-background hover:bg-background hover:text-foreground"
              variant="outline"
              href={`/${locale}/artisans/${artisans[0].slug}`}
            >
              {copy.exploreArtisan}
              <ArrowUpRight className="rtl:-scale-x-100" size={17} />
            </ButtonLink>
          </div>
        </div>
      </section>}
      <section className="section-space">
        <Container>
          <SectionHeading
            eyebrow={copy.promotional}
            title={copy.madeToOrder}
            body={copy.madeToOrderBody}
          />
          <div className="mt-10 grid min-h-[30rem] overflow-hidden bg-secondary text-secondary-foreground lg:grid-cols-2">
            <div className="relative min-h-72">
              <Image
                src="/images/aisha/Algeria art.jpg"
                alt=""
                fill
                sizes="50vw"
                className="object-cover opacity-80"
              />
            </div>
            <div className="flex flex-col justify-center p-8 sm:p-14">
              <p className="font-serif text-4xl">{copy.promotional}</p>
              <p className="mt-5 max-w-md leading-7 text-secondary-foreground/75">
                {copy.madeToOrderBody}
              </p>
              <ButtonLink
                variant="outline"
                className="mt-8 self-start"
                href={`/${locale}/products`}
              >
                {copy.viewAll}
              </ButtonLink>
            </div>
          </div>
        </Container>
      </section>
      <section className="border-y border-border">
        <Container className="grid md:grid-cols-4">
          {trust.map(([Icon, label]) => (
            <div
              className="flex items-center gap-4 border-b border-border py-7 md:border-b-0 md:border-e md:px-6 md:first:ps-0 md:last:border-e-0"
              key={label}
            >
              <Icon size={22} strokeWidth={1.3} />
              <span className="text-sm">{label}</span>
            </div>
          ))}
        </Container>
      </section>
      <section className="section-space bg-muted">
        <Container>
          <SectionHeading
            eyebrow={copy.artisans}
            title={copy.artisanDirectory}
          />
          <div className="mt-12 grid gap-8 md:grid-cols-3">
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
      </section>
      <section className="bg-primary text-primary-foreground">
        <Container className="flex min-h-[28rem] flex-col items-center justify-center py-20 text-center">
          <h2 className="font-serif text-5xl sm:text-6xl">{copy.newsletter}</h2>
          <p className="mt-5 max-w-xl text-primary-foreground/75">
            {copy.newsletterBody}
          </p>
          <form className="mt-8 flex w-full max-w-xl flex-col gap-3 sm:flex-row">
            <label className="sr-only" htmlFor="newsletter-email">
              {copy.email}
            </label>
            <input
              id="newsletter-email"
              type="email"
              required
              placeholder={copy.email}
              className="h-12 flex-1 border border-primary-foreground/50 bg-transparent px-4 placeholder:text-primary-foreground/60"
            />
            <button
              className="h-12 bg-primary-foreground px-7 text-sm font-medium text-primary"
              type="submit"
            >
              {copy.subscribe}
            </button>
          </form>
          <p className="mt-4 text-xs text-primary-foreground/60">
            {copy.privacy}
          </p>
        </Container>
      </section>
    </>
  );
}
