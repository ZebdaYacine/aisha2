"use client";
import { Filter } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { localized } from "@/features/catalogue/format";
import type { Category, Product } from "@/features/catalogue/types";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { Button } from "@/core/components/ui/button";
import { EmptyState } from "@/core/components/feedback/empty-state";
import { Modal } from "@/core/components/ui/modal";
import { ProductGrid } from "./product-grid";
type Sort = "featured" | "newest" | "low" | "high";
export function ProductExplorer({
  initialProducts,
  locale,
  copy,
  categories,
}: {
  initialProducts: Product[];
  locale: Locale;
  copy: StoreCopy;
  categories: Category[];
}) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const currentQuery = searchParams.toString();
  const [category, setCategory] = useState(searchParams.get("category") ?? "");
  const [availability, setAvailability] = useState(searchParams.get("availability") ?? "");
  const [sort, setSort] = useState<Sort>((searchParams.get("sort") as Sort) ?? "featured");
  const [mobile, setMobile] = useState(false);
  const visible = useMemo(
    () =>
      initialProducts
        .filter(
          (p) =>
            (!category || p.categorySlug === category) &&
            (!availability || p.availability === availability),
        )
        .sort((a, b) =>
          sort === "low"
            ? a.priceMinor - b.priceMinor
            : sort === "high"
              ? b.priceMinor - a.priceMinor
              : sort === "newest"
                ? Number(Boolean(b.new)) - Number(Boolean(a.new))
                : Number(Boolean(b.featured)) - Number(Boolean(a.featured)),
        ),
    [initialProducts, category, availability, sort],
  );
  useEffect(() => {
    const params = new URLSearchParams(currentQuery);
    const setOrDelete = (key: string, value: string) => value ? params.set(key, value) : params.delete(key);
    setOrDelete("category", category);
    setOrDelete("availability", availability);
    setOrDelete("sort", sort === "featured" ? "" : sort);
    const query = params.toString();
    const target = query ? `${pathname}?${query}` : pathname;
    const current = currentQuery ? `${pathname}?${currentQuery}` : pathname;
    if (target !== current) router.replace(target, { scroll: false });
  }, [availability, category, currentQuery, pathname, router, sort]);
  const filter = <FilterFields />;
  function FilterFields() {
    return (
      <div className="space-y-8">
        <fieldset>
          <legend className="mb-4 font-medium">{copy.category}</legend>
          <label className="flex gap-3 py-2">
            <input
              type="radio"
              name="category"
              checked={!category}
              onChange={() => setCategory("")}
            />
            {copy.allCategories}
          </label>
          {categories.map((item) => (
            <label className="flex gap-3 py-2" key={item.slug}>
              <input
                type="radio"
                name="category"
                checked={category === item.slug}
                onChange={() => setCategory(item.slug)}
              />
              {localized(item.name, locale)}
            </label>
          ))}
        </fieldset>
        <fieldset>
          <legend className="mb-4 font-medium">{copy.availability}</legend>
          {[
            ["in_stock", copy.inStock],
            ["low_stock", copy.lowStock],
            ["made_to_order", copy.madeToOrder],
            ["out_of_stock", copy.outOfStock],
          ].map(([value, label]) => (
            <label className="flex gap-3 py-2" key={value}>
              <input
                type="radio"
                name="availability"
                checked={availability === value}
                onChange={() => setAvailability(value)}
              />
              {label}
            </label>
          ))}
        </fieldset>
        <Button
          variant="outline"
          className="w-full"
          onClick={() => {
            setCategory("");
            setAvailability("");
          }}
        >
          {copy.clear}
        </Button>
      </div>
    );
  }
  return (
    <div>
      <div className="mb-8 flex flex-wrap items-center justify-between gap-4 border-y border-border py-4">
        <p className="text-sm">
          {visible.length} {copy.results}
        </p>
        <div className="flex gap-3">
          <Button
            variant="outline"
            className="lg:hidden"
            onClick={() => setMobile(true)}
          >
            <Filter size={16} />
            {copy.filters}
          </Button>
          <label className="flex items-center gap-3 text-sm">
            <span className="hidden sm:inline">{copy.sort}</span>
            <select
              className="h-12 border border-border bg-background px-3"
              value={sort}
              onChange={(e) => setSort(e.target.value as Sort)}
            >
              <option value="featured">{copy.featuredSort}</option>
              <option value="newest">{copy.newest}</option>
              <option value="low">{copy.priceLow}</option>
              <option value="high">{copy.priceHigh}</option>
            </select>
          </label>
        </div>
      </div>
      <div className="grid gap-10 lg:grid-cols-[15rem_1fr]">
        <aside className="hidden border-e border-border pe-8 lg:block">
          {filter}
        </aside>
        <div>
          {visible.length ? (
            <ProductGrid
              products={visible}
              locale={locale}
              copy={copy}
              className="lg:grid-cols-3 xl:grid-cols-4"
            />
          ) : (
            <EmptyState
              title={copy.noResults}
              body={copy.tryAgain}
              action={
                <Button
                  onClick={() => {
                    setCategory("");
                    setAvailability("");
                  }}
                >
                  {copy.clear}
                </Button>
              }
            />
          )}
        </div>
      </div>
      {mobile && (
        <Modal label={copy.filters} closeLabel={copy.close} onClose={() => setMobile(false)} panelClassName="absolute inset-y-0 start-0 w-[min(90vw,24rem)] overflow-y-auto p-6">
          <aside
            aria-label={copy.filters}
          >
            <div className="mb-8 flex justify-between pe-12">
              <h2 className="font-serif text-3xl">{copy.filters}</h2>
            </div>
            {filter}
            <Button className="mt-5 w-full" onClick={() => setMobile(false)}>
              {copy.apply}
            </Button>
          </aside>
        </Modal>
      )}
    </div>
  );
}
