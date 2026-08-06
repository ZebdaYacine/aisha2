import { Heart, Plus } from "lucide-react";
import Image from "next/image";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import { localized } from "@/features/catalogue/format";
import type { Product } from "@/features/catalogue/types";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { ProductPrice } from "./product-price";
import { Badge } from "@/core/components/ui/badge";
import { Skeleton } from "@/core/components/ui/skeleton";

export function ProductCard({ product, locale, copy }: { product: Product; locale: Locale; copy: StoreCopy }) {
  return <article className="group min-w-0">
    <div className="product-media relative isolate aspect-[4/5] overflow-hidden bg-muted">
      <LocalizedLink className="block h-full w-full" locale={locale} href={`/products/${product.slug}`} aria-label={localized(product.name, locale)}>
        <Image src={product.images[0]} alt={localized(product.name, locale)} fill sizes="(max-width:640px) 50vw,(max-width:1024px) 33vw,25vw" className="product-card-image pointer-events-none object-cover transition duration-500 group-hover:scale-[1.025]" />
      </LocalizedLink>
      <div className="absolute start-3 top-3 flex flex-col items-start gap-1">{product.new && <Badge>{copy.newArrivals}</Badge>}{product.availability === "made_to_order" && <Badge className="bg-foreground text-background">{copy.madeToOrder}</Badge>}</div>
      <button type="button" aria-label={copy.wishlist} className="absolute end-2 top-2 grid size-11 place-items-center bg-background/90 transition hover:bg-background"><Heart size={18} /></button>
      {product.availability !== "out_of_stock" && <button type="button" aria-label={copy.addToCart} className="absolute inset-x-3 bottom-3 hidden min-h-11 items-center justify-center gap-2 bg-background text-xs font-medium opacity-0 transition group-hover:opacity-100 focus:opacity-100 lg:flex"><Plus size={16} />{copy.addToCart}</button>}
    </div>
    <div className="pt-4"><p className="text-[.68rem] uppercase tracking-[.13em] text-muted-foreground">{localized(product.region, locale)}</p><LocalizedLink locale={locale} href={`/products/${product.slug}`} className="mt-1 block font-serif text-xl leading-tight hover:text-primary">{localized(product.name, locale)}</LocalizedLink><LocalizedLink locale={locale} href={`/artisans/${product.artisanSlug}`} className="mt-1 block text-xs text-muted-foreground hover:text-foreground">{product.artisanName}</LocalizedLink><ProductPrice className="mt-3" {...product} locale={locale} /></div>
  </article>;
}

export function ProductCardSkeleton() { return <div aria-hidden><Skeleton className="aspect-[4/5]" /><Skeleton className="mt-4 h-3 w-20" /><Skeleton className="mt-3 h-5 w-4/5" /><Skeleton className="mt-3 h-4 w-24" /></div>; }
