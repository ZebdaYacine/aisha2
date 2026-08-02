import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { ProductExplorer } from "@/components/commerce/product-explorer";
import { Container } from "@/components/layout/container";
import { Breadcrumbs } from "@/components/shared/breadcrumbs";
import { catalogue } from "@/features/storefront/api";
import { localized } from "@/features/storefront/format";
import { isLocale } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";

export async function generateMetadata({ params }: { params: Promise<{ locale: string; slug: string }> }): Promise<Metadata> {
  const { locale, slug } = await params;
  if (!isLocale(locale)) return {};
  const category = (await catalogue(locale)).categories.find((item) => item.slug === slug);
  if (!category) return {};
  return { title: `${localized(category.name, locale)} — AISHA`, description: localized(category.description, locale), alternates: { canonical: `/${locale}/categories/${slug}` } };
}
export default async function CategoryPage({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { locale, slug } = await params;
  if (!isLocale(locale)) notFound();
  const data = await catalogue(locale);
  const category = data.categories.find((item) => item.slug === slug);
  if (!category) notFound();
  const copy = storeCopy(locale),
    items = data.products.filter((p) => p.categorySlug === slug);
  return (
    <Container className="py-12">
      <Breadcrumbs
        locale={locale}
        items={[
          { label: copy.categories, href: "/products" },
          { label: localized(category.name, locale) },
        ]}
      />
      <header className="mt-10 max-w-3xl">
        <h1 className="font-serif text-5xl sm:text-6xl">
          {localized(category.name, locale)}
        </h1>
        <p className="mt-5 text-lg text-muted-foreground">
          {localized(category.description, locale)}
        </p>
      </header>
      <div className="mt-12">
        <ProductExplorer initialProducts={items} categories={data.categories} locale={locale} copy={copy} />
      </div>
    </Container>
  );
}
