import Image from "next/image";
import { notFound } from "next/navigation";
import { Breadcrumbs } from "@/core/components/shared/breadcrumbs";
import { Container } from "@/core/components/layout/container";
import { ButtonLink } from "@/core/components/ui/button";
import { ProductGrid } from "@/features/product";
import { catalogue, catalogueWorkshop } from "@/features/catalogue/api";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export default async function WorkshopPage({
  params,
}: {
  params: Promise<{ locale: string; slug: string }>;
}) {
  const { locale, slug } = await params;
  if (!isLocale(locale)) notFound();
  const workshop = await catalogueWorkshop(slug, locale).catch(() => undefined);
  if (!workshop) notFound();
  const data = await catalogue(locale, { workshop: workshop.slug });
  const copy = storeCopy(locale);
  return (
    <>
      <Container className="py-8">
        <Breadcrumbs
          locale={locale}
          items={[
            { label: copy.workshop, href: `/${locale}/workshops` },
            { label: workshop.name },
          ]}
        />
        <div className="mt-10 grid gap-10 lg:grid-cols-[.9fr_1.1fr] lg:items-center">
          <div className="relative aspect-[4/3] overflow-hidden bg-secondary">
            <Image
              src={workshop.image}
              alt={workshop.name}
              fill
              sizes="(max-width:1024px) 100vw, 45vw"
              className="object-cover"
            />
          </div>
          <div>
            <p className="text-xs uppercase tracking-[.2em] text-primary">
              {workshop.wilaya} · {workshop.location}
            </p>
            <h1 className="mt-4 font-serif text-5xl">{workshop.name}</h1>
            <p className="mt-5 text-lg leading-8 text-muted-foreground">
              {workshop.description}
            </p>
            <p className="mt-4 text-sm text-muted-foreground">
              {workshop.craft} · {workshop.productCount}{" "}
              {copy.products.toLowerCase()}
            </p>
            <ButtonLink
              className="mt-7"
              variant="outline"
              href={`/${locale}/artisans/${workshop.artisanId}`}
            >
              {copy.exploreArtisan}: {workshop.artisanName}
            </ButtonLink>
          </div>
        </div>
      </Container>
      <section className="section-space bg-muted">
        <Container>
          <p className="text-xs uppercase tracking-[.2em] text-primary">
            {copy.products}
          </p>
          <h2 className="mt-4 font-serif text-4xl">{workshop.name}</h2>
          <div className="mt-10">
            <ProductGrid products={data.products} locale={locale} copy={copy} />
          </div>
        </Container>
      </section>
    </>
  );
}
