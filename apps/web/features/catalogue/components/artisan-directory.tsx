"use client";

import { useMemo, useState } from "react";

import { EmptyState } from "@/core/components/feedback/empty-state";
import { ArtisanCard } from "@/features/artisan";
import { localized } from "@/features/catalogue/format";
import type { Artisan } from "@/features/catalogue/types";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";

export function ArtisanDirectory({ artisans, locale, copy }: { artisans: Artisan[]; locale: Locale; copy: StoreCopy }) {
  const [query, setQuery] = useState("");
  const [region, setRegion] = useState("");
  const [craft, setCraft] = useState("");
  const regions = useMemo(() => [...new Set(artisans.map((artisan) => localized(artisan.region, locale)).filter(Boolean))].sort(), [artisans, locale]);
  const crafts = useMemo(() => [...new Set(artisans.map((artisan) => localized(artisan.craft, locale)).filter(Boolean))].sort(), [artisans, locale]);
  const visible = useMemo(() => {
    const normalized = query.trim().toLowerCase();
    return artisans.filter((artisan) => {
      const searchable = [artisan.name, artisan.workshop, localized(artisan.biography, locale), localized(artisan.craft, locale), localized(artisan.region, locale)].join(" ").toLowerCase();
      return (!normalized || searchable.includes(normalized)) && (!region || localized(artisan.region, locale) === region) && (!craft || localized(artisan.craft, locale) === craft);
    });
  }, [artisans, craft, locale, query, region]);

  return (
    <>
      <form className="mt-10 grid gap-3 border-y border-border py-5 sm:grid-cols-3" onSubmit={(event) => event.preventDefault()}>
        <label>
          <span className="sr-only">{copy.search}</span>
          <input className="h-12 w-full border border-border bg-background px-4" placeholder={copy.search} type="search" value={query} onChange={(event) => setQuery(event.target.value)} />
        </label>
        <label>
          <span className="sr-only">{copy.region}</span>
          <select className="h-12 w-full border border-border bg-background px-4" value={region} onChange={(event) => setRegion(event.target.value)}>
            <option value="">{copy.region}</option>
            {regions.map((item) => <option key={item} value={item}>{item}</option>)}
          </select>
        </label>
        <label>
          <span className="sr-only">{copy.craft}</span>
          <select className="h-12 w-full border border-border bg-background px-4" value={craft} onChange={(event) => setCraft(event.target.value)}>
            <option value="">{copy.craft}</option>
            {crafts.map((item) => <option key={item} value={item}>{item}</option>)}
          </select>
        </label>
      </form>
      {visible.length ? <div className="mt-12 grid gap-10 md:grid-cols-2 lg:grid-cols-3">{visible.map((artisan) => <ArtisanCard artisan={artisan} locale={locale} copy={copy} key={artisan.slug} />)}</div> : <EmptyState title={copy.noResults} body={copy.tryAgain} />}
    </>
  );
}
