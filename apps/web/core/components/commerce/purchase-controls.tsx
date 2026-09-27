"use client";
import { Heart, ShoppingBag } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { Button } from "@/core/components/ui/button";
import { useCart } from "@/features/cart/viewmodel/cart-context";
import type { Availability } from "@/features/catalogue/types";
import type { StoreCopy } from "@/core/lib/store-copy";
import { QuantitySelector } from "./quantity-selector";

export function PurchaseControls({ slug, availability, copy, productName, artisanName, workshopName, image, priceMinor, currency }: { slug: string; availability: Availability; copy: StoreCopy; productName: string; artisanName?: string; workshopName?: string; image?: string; priceMinor: number; currency: string }) {
  const [quantity, setQuantity] = useState(1);
  const [saved, setSaved] = useState(false);
  const { add, setOpen } = useCart();
  const disabled = !["in_stock", "low_stock", "made_to_order"].includes(availability);
  const toggleWishlist = async () => {
    try {
      const response = await fetch(`/api/wishlist/items/${slug}`, { method: saved ? "DELETE" : "POST" });
      if (response.status === 401) { toast.info(copy.wishlistSignIn); return; }
      if (response.ok) { setSaved(!saved); return; }
      toast.error(copy.wishlistError);
    } catch {
      toast.error(copy.wishlistError);
    }
  };
  return <div>
    <div className="flex items-center justify-between gap-4"><QuantitySelector value={quantity} onChange={setQuantity} label={copy.quantity} /><button type="button" onClick={() => void toggleWishlist()} aria-pressed={saved} aria-label={copy.wishlist} className="grid size-12 place-items-center border border-border"><Heart fill={saved ? "currentColor" : "none"} size={19} /></button></div>
    <Button data-testid="product-add-to-cart" className="mt-5 w-full" disabled={disabled} onClick={() => { add(slug, quantity, { productName, artisanName, workshopName, image, priceMinor, currency }); setOpen(true); }}><ShoppingBag size={18} />{disabled ? (availability === "unknown" ? copy.availabilityPending : copy.outOfStock) : availability === "made_to_order" ? copy.madeToOrder : copy.addToCart}</Button>
    <Button className="mt-3 w-full" variant="outline" disabled={disabled}>{copy.buyNow}</Button>
  </div>;
}
