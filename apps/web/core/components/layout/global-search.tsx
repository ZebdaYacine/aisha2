"use client";
import { Search } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import type { Product, Artisan, Category, Workshop } from "@/features/catalogue/types";
import { localized } from "@/features/catalogue/format";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import { Modal } from "@/core/components/ui/modal";
export function GlobalSearch({
  locale,
  copy,
  products,
  artisans,
  categories,
  workshops,
}: {
  locale: Locale;
  copy: StoreCopy;
  products?: Product[];
  artisans?: Artisan[];
  categories?: Category[];
  workshops?: Workshop[];
}) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [catalogue, setCatalogue] = useState({
    products: products ?? [],
    artisans: artisans ?? [],
    categories: categories ?? [],
    workshops: workshops ?? [],
  });
  const [loadState, setLoadState] = useState<"idle" | "loading" | "loaded" | "failed">(
    products || artisans || categories || workshops ? "loaded" : "idle",
  );

  async function openSearch() {
    setOpen(true);
    if (loadState === "loading" || loadState === "loaded") return;
    setLoadState("loading");
    try {
      const response = await fetch(`/api/catalogue?locale=${locale}`);
      if (!response.ok) throw new Error(`Catalogue request failed: ${response.status}`);
      setCatalogue(await response.json());
      setLoadState("loaded");
    } catch {
      setLoadState("failed");
      toast.error(copy.errorBody);
    }
  }

  const { products: loadedProducts, artisans: loadedArtisans, categories: loadedCategories, workshops: loadedWorkshops } = catalogue;
  const q = query.toLowerCase();
  const searchTerms = q.split(/\s+/).filter(Boolean);
  const searchable = (values: string[]) => values.some((value) => {
    const normalized = value.toLowerCase();
    return searchTerms.some((term) => normalized.includes(term));
  });
  const productMatches = q
    ? loadedProducts.filter((product) => searchable([
        localized(product.name, locale),
        localized(product.summary, locale),
        localized(product.story, locale),
        product.artisanName ?? "",
        product.workshop ?? "",
        localized(product.region, locale),
        localized(product.materials, locale),
        localized(product.method, locale),
        product.categorySlug,
      ])).slice(0, 4)
    : loadedProducts.slice(0, 3);
  const artisanMatches = q
    ? loadedArtisans.filter((artisan) => searchable([
        artisan.name,
        artisan.workshop,
        localized(artisan.biography, locale),
        localized(artisan.craft, locale),
        localized(artisan.region, locale),
      ])).slice(0, 3)
    : [];
  const workshopMatches = q
    ? loadedWorkshops.filter((workshop) => searchable([
        workshop.name,
        workshop.description,
        workshop.wilaya,
        workshop.location,
        workshop.craft,
        workshop.artisanName,
      ])).slice(0, 3)
    : [];
  const categoryMatches = q
    ? loadedCategories.filter((category) => searchable([
        localized(category.name, locale),
        localized(category.description, locale),
        category.slug,
      ])).slice(0, 3)
    : [];
  return (
    <>
      <button
        type="button"
        className="hidden min-h-11 min-w-11 items-center justify-center sm:flex"
        onClick={openSearch}
        aria-label={copy.search}
      >
        <Search size={19} />
      </button>
      {open && (
        <Modal label={copy.searchTitle} closeLabel={copy.close} onClose={() => setOpen(false)} panelClassName="overflow-y-auto">
          <div className="mx-auto max-w-4xl px-4 py-8">
            <div className="flex items-center gap-3 border-b border-foreground">
              <Search />
              <label className="sr-only" htmlFor="global-search">
                {copy.search}
              </label>
              <input
                autoFocus
                id="global-search"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder={copy.searchHint}
                className="h-16 flex-1 bg-transparent text-xl"
              />
            </div>
            <div className="grid gap-10 py-10 md:grid-cols-[1fr_2fr]">
              <div>
                <h2 className="text-xs uppercase tracking-widest">
                  {copy.popular}
                </h2>
                <div className="mt-4 flex flex-wrap gap-2">
                  {loadedCategories.slice(0, 4).map((category) => (
                    <button
                      className="border border-border px-3 py-2 text-sm"
                      onClick={() => setQuery(localized(category.name, locale))}
                      key={category.slug}
                    >
                      {localized(category.name, locale)}
                    </button>
                  ))}
                </div>
              </div>
              <div>
                <h2 className="text-xs uppercase tracking-widest">
                  {copy.searchResults}
                </h2>
                <div className="mt-4 divide-y divide-border">
                  {productMatches.map((product) => (
                    <LocalizedLink
                      onClick={() => setOpen(false)}
                      className="block py-4 font-serif text-2xl hover:text-primary"
                      locale={locale}
                      href={`/products/${product.slug}`}
                      key={product.slug}
                    >
                      {localized(product.name, locale)} · {copy.products}
                    </LocalizedLink>
                  ))}
                  {artisanMatches.map((artisan) => (
                    <LocalizedLink onClick={() => setOpen(false)} className="block py-4" locale={locale} href={`/artisans/${artisan.slug}`} key={`artisan-${artisan.slug}`}>
                      {artisan.name} · {copy.artisan}
                    </LocalizedLink>
                  ))}
                  {workshopMatches.map((workshop) => (
                    <LocalizedLink onClick={() => setOpen(false)} className="block py-4" locale={locale} href={`/workshops/${workshop.slug}`} key={`workshop-${workshop.slug}`}>
                      {workshop.name} · {copy.workshop}
                    </LocalizedLink>
                  ))}
                  {categoryMatches.map((category) => (
                    <LocalizedLink onClick={() => setOpen(false)} className="block py-4" locale={locale} href={`/categories/${category.slug}`} key={`category-${category.slug}`}>
                      {localized(category.name, locale)} · {copy.categories}
                    </LocalizedLink>
                  ))}
                  {q && !productMatches.length && !artisanMatches.length && !workshopMatches.length && !categoryMatches.length && <p className="py-4 text-sm text-muted-foreground">{copy.noResults}</p>}
                </div>
                <LocalizedLink
                  onClick={() => setOpen(false)}
                  locale={locale}
                  href={`/search?q=${encodeURIComponent(query)}`}
                  className="mt-6 inline-block border-b border-foreground pb-1"
                >
                  {copy.viewAll}
                </LocalizedLink>
              </div>
            </div>
          </div>
        </Modal>
      )}
    </>
  );
}
