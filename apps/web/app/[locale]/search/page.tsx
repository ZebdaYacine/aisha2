import { notFound } from "next/navigation";
import { ProductGrid } from "@/features/product";
import { EmptyState } from "@/core/components/feedback/empty-state";
import { Container } from "@/core/components/layout/container";
import { Breadcrumbs } from "@/core/components/shared/breadcrumbs";
import { ButtonLink } from "@/core/components/ui/button";
import { catalogue } from "@/features/catalogue/api";
import { localized } from "@/features/catalogue/format";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
export default async function SearchPage({
  params,
  searchParams,
}: {
  params: Promise<{ locale: string }>;
  searchParams: Promise<{ q?: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const { q = "" } = await searchParams;
  const { products } = await catalogue(locale);
  const copy = storeCopy(locale),
    query = q.trim().toLowerCase(),
    matches = query
      ? products.filter((p) =>
          [
            localized(p.name, locale),
            localized(p.summary, locale),
            localized(p.region, locale),
            localized(p.materials, locale),
          ].some((value) => value.toLowerCase().includes(query)),
        )
      : [];
  return (
    <Container className="py-12">
      <Breadcrumbs locale={locale} items={[{ label: copy.search }]} />
      <h1 className="mt-10 font-serif text-5xl">{copy.searchTitle}</h1>
      <form className="mt-8 flex max-w-3xl border-b border-foreground">
        <label className="sr-only" htmlFor="search-page">
          {copy.search}
        </label>
        <input
          id="search-page"
          name="q"
          defaultValue={q}
          className="h-16 flex-1 bg-transparent text-xl"
          placeholder={copy.searchHint}
        />
        <button className="px-6 font-medium" type="submit">
          {copy.search}
        </button>
      </form>
      <div className="mt-14">
        {matches.length ? (
          <>
            <p className="mb-8 text-sm text-muted-foreground">
              {matches.length} {copy.results}
            </p>
            <ProductGrid products={matches} locale={locale} copy={copy} />
          </>
        ) : (
          <EmptyState
            title={query ? copy.noResults : copy.searchHint}
            body={copy.tryAgain}
            action={
              <ButtonLink href={`/${locale}/products`}>
                {copy.products}
              </ButtonLink>
            }
          />
        )}
      </div>
    </Container>
  );
}
