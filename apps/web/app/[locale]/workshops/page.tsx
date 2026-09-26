import Image from "next/image";
import { notFound } from "next/navigation";
import { Breadcrumbs } from "@/core/components/shared/breadcrumbs";
import { Container } from "@/core/components/layout/container";
import { ButtonLink } from "@/core/components/ui/button";
import { catalogue } from "@/features/catalogue/api";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export default async function WorkshopsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const data = await catalogue(locale);
  const copy = storeCopy(locale);
  return (
    <Container className="py-8">
      <Breadcrumbs locale={locale} items={[{ label: copy.workshop }]} />
      <header className="section-space pb-10">
        <p className="text-xs uppercase tracking-[.2em] text-primary">
          {copy.workshop}
        </p>
        <h1 className="mt-4 max-w-3xl font-serif text-5xl">
          {copy.artisanDirectory}
        </h1>
        <p className="mt-5 max-w-2xl text-lg leading-8 text-muted-foreground">
          {copy.artisanDirectoryBody}
        </p>
      </header>
      <div className="grid gap-10 md:grid-cols-2">
        {data.workshops.map((workshop) => (
          <article className="border border-border bg-card" key={workshop.slug}>
            <div className="relative aspect-[16/10] overflow-hidden bg-secondary">
              <Image
                src={workshop.image}
                alt={workshop.name}
                fill
                sizes="(max-width:768px) 100vw, 50vw"
                className="object-cover"
              />
            </div>
            <div className="p-6">
              <p className="text-xs uppercase tracking-[.16em] text-primary">
                {workshop.wilaya}
              </p>
              <h2 className="mt-3 font-serif text-3xl">{workshop.name}</h2>
              <p className="mt-3 text-sm text-muted-foreground">
                {workshop.craft} · {workshop.productCount}{" "}
                {copy.products.toLowerCase()}
              </p>
              <p className="mt-4 leading-7 text-muted-foreground">
                {workshop.description}
              </p>
              <ButtonLink
                className="mt-6"
                href={`/${locale}/workshops/${workshop.slug}`}
              >
                {copy.exploreArtisan}
              </ButtonLink>
            </div>
          </article>
        ))}
      </div>
    </Container>
  );
}
