"use client";
import { X } from "lucide-react";
import Image from "next/image";
import { useCart } from "@/features/cart/viewmodel/cart-context";
import { getProduct } from "@/features/catalogue/data";
import { localized } from "@/features/catalogue/format";
import type { CartItem as CartItemType } from "@/features/catalogue/types";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { ProductPrice } from "./product-price";
import { QuantitySelector } from "./quantity-selector";

export function CartLineItem({ item, locale, copy }: { item: CartItemType; locale: Locale; copy: StoreCopy }) {
  const product = getProduct(item.productSlug);
  const { update, remove } = useCart();
  if (!product) return null;
  return <article className="grid grid-cols-[6rem_1fr] gap-4 border-b border-border py-5">
    <div className="cart-item-media bg-muted"><Image src={product.images[0]} alt="" fill sizes="96px" className="object-cover" /></div>
    <div className="min-w-0"><div className="flex justify-between gap-3"><h3 className="font-serif text-lg leading-tight">{localized(product.name, locale)}</h3><button type="button" onClick={() => remove(product.slug)} aria-label={`${copy.remove} ${localized(product.name, locale)}`} className="grid size-10 shrink-0 place-items-center"><X size={16} /></button></div><ProductPrice {...product} locale={locale} className="mt-2" /><div className="mt-4"><QuantitySelector value={item.quantity} onChange={(value) => update(product.slug, value)} label={copy.quantity} /></div></div>
  </article>;
}
