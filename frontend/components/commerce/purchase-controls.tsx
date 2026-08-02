"use client";
import { Heart, ShoppingBag } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { useCart } from "@/features/cart/cart-context";
import type { Availability } from "@/features/storefront/types";
import type { StoreCopy } from "@/lib/store-copy";
import { QuantitySelector } from "./quantity-selector";

export function PurchaseControls({ slug, availability, copy }: { slug: string; availability: Availability; copy: StoreCopy }) {
  const [quantity, setQuantity] = useState(1);
  const [saved, setSaved] = useState(false);
  const { add, setOpen } = useCart();
  const disabled = availability === "out_of_stock";
  return <div>
    <div className="flex items-center justify-between gap-4"><QuantitySelector value={quantity} onChange={setQuantity} label={copy.quantity} /><button type="button" onClick={() => setSaved(!saved)} aria-pressed={saved} aria-label={copy.wishlist} className="grid size-12 place-items-center border border-border"><Heart fill={saved ? "currentColor" : "none"} size={19} /></button></div>
    <Button data-testid="product-add-to-cart" className="mt-5 w-full" disabled={disabled} onClick={() => { add(slug, quantity); setOpen(true); }}><ShoppingBag size={18} />{disabled ? copy.outOfStock : availability === "made_to_order" ? copy.madeToOrder : copy.addToCart}</Button>
    <Button className="mt-3 w-full" variant="outline" disabled={disabled}>{copy.buyNow}</Button>
  </div>;
}
