"use client";
import { useEffect, useState } from "react";
import { ButtonLink } from "@/core/components/ui/button";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";

type Item = { productId: string; productName: string; priceMinor: number; currency: string; active: boolean };
export function WishlistList({ locale, copy }: { locale: Locale; copy: StoreCopy }) {
  const [items, setItems] = useState<Item[]>([]);
  const [loaded, setLoaded] = useState(false);
  useEffect(() => { void fetch("/api/wishlist", { cache: "no-store" }).then(async (response) => { if (response.ok) setItems((await response.json()) as Item[]); setLoaded(true); }).catch(() => setLoaded(true)); }, []);
  if (!loaded) return <p className="mt-8 text-sm text-muted-foreground">{copy.provisional}</p>;
  if (!items.length) return <div className="mt-8"><p className="text-muted-foreground">{copy.noResults}</p><ButtonLink href={`/${locale}/products`} className="mt-5">{copy.continueShopping}</ButtonLink></div>;
  return <div className="mt-8 grid gap-4 sm:grid-cols-2">{items.map((item) => <article key={item.productId} className="border border-border p-5"><h2 className="font-serif text-xl">{item.productName}</h2><p className="mt-2 text-sm text-muted-foreground">{item.active ? `${item.priceMinor / 100} ${item.currency}` : copy.outOfStock}</p><button className="mt-4 text-sm underline" onClick={() => void fetch(`/api/wishlist/items/${item.productId}`, { method: "DELETE" }).then(() => setItems((current) => current.filter((entry) => entry.productId !== item.productId)))}>{copy.remove}</button></article>)}</div>;
}
