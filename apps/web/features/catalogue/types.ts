import type { Locale } from "@/core/lib/i18n";

export type LocalizedText = Record<Locale, string>;
export type Availability =
  "in_stock" | "low_stock" | "made_to_order" | "out_of_stock" | "unknown";

export interface Category {
  id?: string;
  slug: string;
  name: LocalizedText;
  description: LocalizedText;
  image: string;
}
export interface Artisan {
  slug: string;
  name: string;
  workshopId?: string;
  workshop: string;
  region: LocalizedText;
  craft: LocalizedText;
  biography: LocalizedText;
  image: string;
  verified: boolean;
  productCount: number;
}
export interface Product {
  slug: string;
  name: LocalizedText;
  summary: LocalizedText;
  story: LocalizedText;
  artisanSlug: string;
  artisanName?: string;
  categorySlug: string;
  workshopId?: string;
  workshop?: string;
  region: LocalizedText;
  materials: LocalizedText;
  method: LocalizedText;
  priceMinor: number;
  previousPriceMinor?: number;
  currency: string;
  availability: Availability;
  images: string[];
  featured?: boolean;
  new?: boolean;
  rating: number;
  reviewCount: number;
}
export interface Workshop {
  slug: string;
  name: string;
  description: string;
  wilaya: string;
  location: string;
  craft: string;
  artisanId: string;
  artisanName: string;
  image: string;
  productCount: number;
}
export interface CartItem {
  productSlug: string;
  quantity: number;
  productName?: string;
  artisanName?: string;
  workshopName?: string;
  priceMinor?: number;
  currency?: string;
  active?: boolean;
  warning?: string;
}
export interface Address {
  fullName: string;
  email: string;
  phone: string;
  line1: string;
  line2?: string;
  city: string;
  postalCode: string;
  country: string;
}
export interface Order {
  id: string;
  placedAt: string;
  status: "processing" | "shipped" | "delivered";
  items: CartItem[];
  totalMinor: number;
  currency: "EUR";
  address: Address;
  trackingEvents: Array<{ label: string; date: string; complete: boolean }>;
}
