"use client";
import { Heart, Plus } from "lucide-react";
import Image from "next/image";
import { useState, type PointerEvent } from "react";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import { localized } from "@/features/catalogue/format";
import type { Product } from "@/features/catalogue/types";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { ProductPrice } from "./product-price";
import { Badge } from "@/core/components/ui/badge";
import { Skeleton } from "@/core/components/ui/skeleton";
import { useOptionalCart } from "@/features/cart/viewmodel/cart-context";
import { toast } from "sonner";

export function ProductCard({ product, locale, copy }: { product: Product; locale: Locale; copy: StoreCopy }) {
  const [saved, setSaved] = useState(false);
  const [imageOrigin, setImageOrigin] = useState("50% 50%");
  const [imageHovered, setImageHovered] = useState(false);
  const cart = useOptionalCart();
  const toggleWishlist = async () => {
    try {
      const response = await fetch(`/api/wishlist/items/${product.slug}`, { method: saved ? "DELETE" : "POST" });
      if (response.status === 401) { toast.info(copy.wishlistSignIn); return; }
      if (response.ok) { setSaved(!saved); return; }
      toast.error(copy.wishlistError);
    } catch {
      toast.error(copy.wishlistError);
    }
  };
  const trackImagePointer = (event: PointerEvent<HTMLDivElement>) => {
    if (event.pointerType && event.pointerType !== "mouse") return;
    const bounds = event.currentTarget.getBoundingClientRect();
    const x = bounds.width ? Math.min(100, Math.max(0, ((event.clientX - bounds.left) / bounds.width) * 100)) : 50;
    const y = bounds.height ? Math.min(100, Math.max(0, ((event.clientY - bounds.top) / bounds.height) * 100)) : 50;
    setImageOrigin(`${x}% ${y}%`);
    setImageHovered(true);
  };
  return <article className="group min-w-0">
    <div className="product-media relative isolate aspect-[4/5] overflow-hidden rounded-xl border border-border bg-muted" onPointerMove={trackImagePointer} onPointerLeave={() => setImageHovered(false)}>
      <LocalizedLink className="block h-full w-full" locale={locale} href={`/products/${product.slug}`} aria-label={localized(product.name, locale)}>
        <Image src={product.images[0]} alt={localized(product.name, locale)} fill sizes="(max-width:640px) 50vw,(max-width:1024px) 33vw,25vw" className="product-card-image pointer-events-none object-cover transition duration-500 group-hover:scale-[1.025]" style={{ transformOrigin: imageOrigin, transform: imageHovered ? "scale(1.2)" : undefined }} />
      </LocalizedLink>
      <div className="absolute start-3 top-3 flex flex-col items-start gap-1">{product.new && <Badge>{copy.newArrivals}</Badge>}{product.availability === "made_to_order" && <Badge className="bg-foreground text-background">{copy.madeToOrder}</Badge>}</div>
      {product.availability === "low_stock" && <Badge className="bg-destructive text-destructive-foreground">{copy.lowStock}</Badge>}
      {product.availability === "out_of_stock" && <Badge className="bg-muted text-muted-foreground">{copy.outOfStock}</Badge>}
      <button type="button" onClick={() => void toggleWishlist()} aria-pressed={saved} aria-label={copy.wishlist} className="absolute end-2 top-2 grid size-11 place-items-center rounded-md bg-background/90 transition hover:bg-background"><Heart fill={saved ? "currentColor" : "none"} size={18} /></button>
      {["in_stock", "low_stock", "made_to_order"].includes(product.availability) && <button type="button" onClick={() => cart?.add(product.slug, 1, { productName: localized(product.name, locale), artisanName: product.artisanName, workshopName: product.workshop, image: product.images[0], priceMinor: product.priceMinor, currency: product.currency })} aria-label={copy.addToCart} className="absolute inset-x-3 bottom-3 hidden min-h-11 items-center justify-center gap-2 bg-background text-xs font-medium opacity-0 transition group-hover:opacity-100 focus:opacity-100 lg:flex"><Plus size={16} />{copy.addToCart}</button>}
    </div>
    <div className="pt-4"><p className="text-[.68rem] uppercase tracking-[.13em] text-muted-foreground">{localized(product.region, locale)}</p><LocalizedLink locale={locale} href={`/products/${product.slug}`} className="mt-1 block font-serif text-xl leading-tight hover:text-primary">{localized(product.name, locale)}</LocalizedLink><LocalizedLink locale={locale} href={`/artisans/${product.artisanSlug}`} className="mt-1 block text-xs text-muted-foreground hover:text-foreground">{product.artisanName}</LocalizedLink><ProductPrice className="mt-3" {...product} locale={locale} />{typeof product.availableQuantity === "number" && product.availability !== "made_to_order" && <p className="mt-2 text-xs text-muted-foreground">{copy.available}: {product.availableQuantity}</p>}</div>
  </article>;
}

export function ProductCardSkeleton() { return <div aria-hidden><Skeleton className="aspect-[4/5]" /><Skeleton className="mt-4 h-3 w-20" /><Skeleton className="mt-3 h-5 w-4/5" /><Skeleton className="mt-3 h-4 w-24" /></div>; }
