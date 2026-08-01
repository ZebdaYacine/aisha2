import type { Locale } from "@/lib/i18n";

export type LocalizedText = Record<Locale, string>;
export type Availability = "in_stock" | "low_stock" | "made_to_order" | "out_of_stock";

export interface Category { slug: string; name: LocalizedText; description: LocalizedText; image: string }
export interface Artisan { slug: string; name: string; workshop: string; region: LocalizedText; craft: LocalizedText; biography: LocalizedText; image: string; verified: boolean; productCount: number }
export interface Product {
  slug: string; name: LocalizedText; summary: LocalizedText; story: LocalizedText; artisanSlug: string; artisanName?: string;
  categorySlug: string; region: LocalizedText; materials: LocalizedText; method: LocalizedText;
  priceMinor: number; previousPriceMinor?: number; currency: string; availability: Availability;
  images: string[]; featured?: boolean; new?: boolean; rating: number; reviewCount: number;
}
export interface CartItem { productSlug: string; quantity: number }
export interface Address { fullName: string; email: string; phone: string; line1: string; line2?: string; city: string; postalCode: string; country: string }
export interface Order { id: string; placedAt: string; status: "processing" | "shipped" | "delivered"; items: CartItem[]; totalMinor: number; currency: "EUR"; address: Address; trackingEvents: Array<{ label: string; date: string; complete: boolean }> }
