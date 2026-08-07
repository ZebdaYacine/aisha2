import type { Locale } from "@/core/lib/i18n";
import type { Artisan, Category, LocalizedText, Product } from "./types";
const baseURL = process.env.API_BASE_URL ?? "http://localhost:8080/api/v1";
const local = (value: string): LocalizedText => ({
  en: value,
  fr: value,
  ar: value,
  es: value,
});
const categoryImages: Record<string, string> = {
  pottery: "/images/aisha/A1.jpg",
  "pottery-and-porcelain": "/images/aisha/A1.jpg",
  textiles: "/images/aisha/M1.jpg",
  "carpets-and-textiles": "/images/aisha/M1.jpg",
  jewellery: "/images/aisha/O5.jpg",
  "jewellery-and-beauty-accessories": "/images/aisha/O5.jpg",
  leather:
    "/images/aisha/Snapinsta.app_470984486_18032514275416001_5509346141223573205_n_1080.jpg",
  "leather-goods":
    "/images/aisha/Snapinsta.app_470984486_18032514275416001_5509346141223573205_n_1080.jpg",
  metalwork: "/images/aisha/Screenshot 2025-01-04 002802.png",
  "traditional-copper-products": "/images/aisha/Screenshot 2025-01-04 002802.png",
  basketry: "/images/aisha/CP1.jpg",
  "bamboo-and-halfa-products": "/images/aisha/CP1.jpg",
  "decoration-and-art": "/images/aisha/Algeria art.jpg",
  "domestic-use": "/images/aisha/plateu en bois 1.jpg",
  embroidery: "/images/aisha/Snapinsta.app_447895980_18152353162314221_3689087632318507403_n_1080.jpg",
};
const fallback = "/images/aisha/Algeria art.jpg";
type Page<T> = { items: T[]; page: number; pageSize: number; total: number };
type CategoryDTO = { id: string; slug: string; name: string };
type ProductDTO = {
  id: string;
  artisanId: string;
  artisanName: string;
  categorySlug: string;
  name: string;
  description: string;
  story: string;
  materials: string;
  productionMethod: string;
  region: string;
  currency: string;
  priceMinor: number;
  media: string[];
};
type ArtisanDTO = {
  id: string;
  name: string;
  workshop: string;
  wilaya: string;
  location: string;
  biography: string;
  media: string[];
  productCount: number;
};
async function get<T>(path: string): Promise<T> {
  const response = await fetch(`${baseURL}${path}`, { cache: "no-store" });
  if (!response.ok)
    throw new Error(
      `Catalogue ${path.split("?")[0]} request failed: ${response.status}`,
    );
  return response.json();
}
const image = (keys: string[] | undefined, category?: string) =>
  keys?.find((key) => key.startsWith("/") || key.startsWith("http")) ??
  categoryImages[category ?? ""] ??
  fallback;
const product = (item: ProductDTO): Product => ({
  slug: item.id,
  name: local(item.name),
  summary: local(item.description),
  story: local(item.story),
  artisanSlug: item.artisanId,
  artisanName: item.artisanName,
  categorySlug: item.categorySlug,
  region: local(item.region),
  materials: local(item.materials),
  method: local(item.productionMethod),
  priceMinor: item.priceMinor,
  currency: item.currency,
  availability: "out_of_stock",
  images: item.media?.length
    ? item.media.map((key) => image([key], item.categorySlug))
    : [image(undefined, item.categorySlug)],
  rating: 0,
  reviewCount: 0,
});
const artisan = (item: ArtisanDTO): Artisan => ({
  slug: item.id,
  name: item.name,
  workshop: item.workshop,
  region: local(item.location || item.wilaya),
  craft: local(""),
  biography: local(item.biography),
  image: image(item.media),
  verified: true,
  productCount: item.productCount,
});
export async function catalogue(locale: Locale) {
  const query = `locale=${locale}&pageSize=100`;
  const [categoryPage, productPage, artisanPage] = await Promise.all([
    get<Page<CategoryDTO>>(`/categories?${query}`),
    get<Page<ProductDTO>>(`/products?${query}`),
    get<Page<ArtisanDTO>>(`/artisans?${query}`),
  ]);
  return {
    categories: categoryPage.items.map(
      (item) =>
        ({
          id: item.id,
          slug: item.slug,
          name: local(item.name),
          description: local(""),
          image: image(undefined, item.slug),
        }) satisfies Category,
    ),
    products: productPage.items.map(product),
    artisans: artisanPage.items.map(artisan),
  };
}
export type Catalogue = Awaited<ReturnType<typeof catalogue>>;
export async function catalogueProduct(id: string, locale: Locale) {
  return product(
    await get<ProductDTO>(
      `/products/${encodeURIComponent(id)}?locale=${locale}`,
    ),
  );
}
export async function catalogueArtisan(id: string, locale: Locale) {
  return artisan(
    await get<ArtisanDTO>(
      `/artisans/${encodeURIComponent(id)}?locale=${locale}`,
    ),
  );
}
