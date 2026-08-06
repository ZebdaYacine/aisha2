import { notFound } from "next/navigation";
import { ProductExplorer } from "@/core/components/commerce/product-explorer";
import { Container } from "@/core/components/layout/container";
import { Breadcrumbs } from "@/core/components/shared/breadcrumbs";
import { catalogue } from "@/features/catalogue/api";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
import type { Metadata } from "next";

export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }): Promise<Metadata> {
  const { locale } = await params;
  if (!isLocale(locale)) return {};
  const copy = storeCopy(locale);
  return { title: `${copy.products} — AISHA`, description: copy.featuredBody, alternates: { canonical: `/${locale}/products` } };
}

export default async function ProductsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const copy = storeCopy(locale);
  const { products, categories } = await catalogue(locale);
  return (
    <Container className="py-10 lg:py-16">
      <Breadcrumbs locale={locale} items={[{ label: copy.products }]} />
      <header className="mt-9 max-w-3xl">
        <p className="text-xs uppercase tracking-[.2em] text-primary">
          {copy.shop}
        </p>
        <h1 className="mt-3 font-serif text-5xl sm:text-6xl">
          {copy.products}
        </h1>
        <p className="mt-5 text-lg leading-8 text-muted-foreground">
          {copy.featuredBody}
        </p>
      </header>
      <div className="mt-12">
        <ProductExplorer
          initialProducts={products}
          locale={locale}
          copy={copy}
          categories={categories}
        />
      </div>
    </Container>
  );
}
