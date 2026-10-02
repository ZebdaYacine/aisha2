"use client";
import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import type { CartItem } from "@/features/catalogue/types";
interface CartValue {
  items: CartItem[];
  open: boolean;
  setOpen: (open: boolean) => void;
  add: (slug: string, quantity?: number, metadata?: Pick<CartItem, "productName" | "artisanName" | "workshopName" | "image" | "priceMinor" | "currency">) => void;
  update: (slug: string, quantity: number) => void;
  remove: (slug: string) => void;
  clear: () => void;
  count: number;
}
const CartContext = createContext<CartValue | null>(null);
const key = "aisha-demo-cart";
export function CartProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<CartItem[]>([]);
  const [open, setOpen] = useState(false);
  const [authenticated, setAuthenticated] = useState(false);
  useEffect(() => {
    if (process.env.NODE_ENV === "test") return;
    let active = true;
    const sync = async () => {
      let local: CartItem[] = [];
      try {
        const saved = localStorage.getItem(key);
        if (saved) local = JSON.parse(saved) as CartItem[];
      } catch {}
      if (active) setItems(local);
      try {
        const response = await fetch("/api/cart", { cache: "no-store" });
        if (!response.ok) return;
        setAuthenticated(true);
        if (local.length) {
          await fetch("/api/cart/merge", {
            method: "POST",
            headers: { "content-type": "application/json" },
            body: JSON.stringify({ items: local.map((item) => ({ productId: item.productSlug, quantity: item.quantity })) }),
          });
        }
        const latest = await fetch("/api/cart", { cache: "no-store" });
        if (!latest.ok || !active) return;
        const body = (await latest.json()) as { items?: Array<{ productId: string; quantity: number; productName?: string; artisanName?: string; workshopName?: string; image?: string; priceMinor?: number; currency?: string; active?: boolean; warning?: string }> } | Array<{ productId: string; quantity: number; productName?: string; artisanName?: string; workshopName?: string; image?: string; priceMinor?: number; currency?: string; active?: boolean; warning?: string }>;
        const serverItems = Array.isArray(body) ? body : body.items ?? [];
        setItems(serverItems.map((item) => ({ productSlug: item.productId, quantity: item.quantity, productName: item.productName, artisanName: item.artisanName, workshopName: item.workshopName, image: item.image, priceMinor: item.priceMinor, currency: item.currency, active: item.active, warning: item.warning })));
      } catch {
        // Anonymous carts remain fully usable when the authenticated API is unavailable.
      }
    };
    void sync();
    return () => { active = false; };
  }, []);
  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(items));
  }, [items]);
  useEffect(() => {
    const clearOnLogout = () => {
      setItems([]);
      setOpen(false);
      setAuthenticated(false);
      try {
        localStorage.removeItem(key);
      } catch {
        // Storage can be unavailable in privacy-restricted browsers.
      }
    };
    window.addEventListener("aisha:session-cleared", clearOnLogout);
    return () => window.removeEventListener("aisha:session-cleared", clearOnLogout);
  }, []);
  const value = useMemo<CartValue>(
    () => ({
      items,
      open,
      setOpen,
      count: items.reduce((sum, item) => sum + item.quantity, 0),
      add: (slug, quantity = 1, metadata) => {
        setItems((current) => {
          const found = current.find((item) => item.productSlug === slug);
          return found
            ? current.map((item) =>
                item.productSlug === slug
                  ? { ...item, ...metadata, quantity: item.quantity + quantity }
                : item,
              )
            : [...current, { productSlug: slug, quantity, ...metadata }];
        });
        if (authenticated) {
          void fetch("/api/cart/items", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify({ productId: slug, quantity }) });
        }
        setOpen(true);
      },
      update: (slug, quantity) => {
        setItems((current) =>
          quantity < 1
            ? current.filter((item) => item.productSlug !== slug)
            : current.map((item) =>
                item.productSlug === slug ? { ...item, quantity } : item,
              ),
        );
        if (authenticated) {
          void fetch(quantity < 1 ? `/api/cart/items/${slug}` : `/api/cart/items/${slug}`, { method: quantity < 1 ? "DELETE" : "PATCH", headers: quantity < 1 ? undefined : { "content-type": "application/json" }, body: quantity < 1 ? undefined : JSON.stringify({ productId: slug, quantity }) });
        }
      },
      remove: (slug) => {
        setItems((current) => current.filter((item) => item.productSlug !== slug));
        if (authenticated) void fetch(`/api/cart/items/${slug}`, { method: "DELETE" });
      },
      clear: () => {
        const current = items;
        setItems([]);
        if (authenticated) void Promise.all(current.map((item) => fetch(`/api/cart/items/${item.productSlug}`, { method: "DELETE" })));
      },
    }),
    [items, open, authenticated],
  );
  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}
export function useCart() {
  const value = useContext(CartContext);
  if (!value) throw new Error("useCart must be used inside CartProvider");
  return value;
}
export function useOptionalCart() {
  return useContext(CartContext);
}
