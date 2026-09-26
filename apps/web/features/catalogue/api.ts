import type { Locale } from "@/core/lib/i18n";
import type {
  Artisan,
  Category,
  LocalizedText,
  Product,
  Workshop,
} from "./types";
const baseURL = process.env.API_BASE_URL ?? "http://localhost:8088/api/v1";
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
  "traditional-copper-products":
    "/images/aisha/Screenshot 2025-01-04 002802.png",
  basketry: "/images/aisha/CP1.jpg",
  "bamboo-and-halfa-products": "/images/aisha/CP1.jpg",
  "decoration-and-art": "/images/aisha/Algeria art.jpg",
  "domestic-use": "/images/aisha/plateu en bois 1.jpg",
  embroidery:
    "/images/aisha/Snapinsta.app_447895980_18152353162314221_3689087632318507403_n_1080.jpg",
};
const fallback = "/images/aisha/Algeria art.jpg";
type Page<T> = { items: T[]; page: number; pageSize: number; total: number };
type CategoryDTO = { id: string; slug: string; name: string };
type ProductDTO = {
  id: string;
  artisanId: string;
  artisanName: string;
  workshopId?: string;
  workshop?: string;
  categorySlug: string;
  name: string;
  description: string;
  story: string;
  materials: string;
  productionMethod: string;
  region: string;
  currency: string;
  priceMinor: number;
  availability?: string;
  availableQuantity?: number;
  media: string[];
};
type ArtisanDTO = {
  id: string;
  name: string;
  workshopId?: string;
  workshop: string;
  wilaya: string;
  location: string;
  biography: string;
  media: string[];
  productCount: number;
  craft?: string;
};
type WorkshopDTO = {
  id: string;
  name: string;
  description: string;
  wilaya: string;
  location: string;
  craft: string;
  artisanId: string;
  artisanName: string;
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
const availability = (value: string | undefined): Product["availability"] =>
  (
    ({
      IN_STOCK: "in_stock",
      LOW_STOCK: "low_stock",
      MADE_TO_ORDER: "made_to_order",
      OUT_OF_STOCK: "out_of_stock",
    }) as Record<string, Product["availability"]>
  )[value ?? ""] ?? "unknown";
const product = (item: ProductDTO): Product => ({
  slug: item.id,
  name: local(item.name),
  summary: local(item.description),
  story: local(item.story),
  artisanSlug: item.artisanId,
  artisanName: item.artisanName,
  workshopId: item.workshopId,
  workshop: item.workshop,
  categorySlug: item.categorySlug,
  region: local(item.region),
  materials: local(item.materials),
  method: local(item.productionMethod),
  priceMinor: item.priceMinor,
  currency: item.currency,
  availability: availability(item.availability),
  images: item.media?.length
    ? item.media.map((key) => image([key], item.categorySlug))
    : [image(undefined, item.categorySlug)],
  rating: 0,
  reviewCount: 0,
});
const artisan = (item: ArtisanDTO): Artisan => ({
  slug: item.id,
  name: item.name,
  workshopId: item.workshopId,
  workshop: item.workshop,
  region: local(item.location || item.wilaya),
  craft: local(item.craft ?? ""),
  biography: local(item.biography),
  image: image(item.media),
  verified: false,
  productCount: item.productCount,
});
export type CatalogueOptions = {
  page?: number;
  pageSize?: number;
  category?: string;
  query?: string;
  workshop?: string;
};

export async function catalogue(
  locale: Locale,
  options: CatalogueOptions = {},
) {
  const params = new URLSearchParams({
    locale,
    pageSize: String(options.pageSize ?? 100),
  });
  if (options.page) params.set("page", String(options.page));
  if (options.category) params.set("category", options.category);
  if (options.query?.trim()) params.set("q", options.query.trim());
  if (options.workshop) params.set("workshop", options.workshop);
  const query = params.toString();
  const [categoryPage, productPage, artisanPage, workshopPage] =
    await Promise.all([
      get<Page<CategoryDTO>>(`/categories?${query}`),
      get<Page<ProductDTO>>(`/products?${query}`),
      get<Page<ArtisanDTO>>(`/artisans?${query}`),
      get<Page<WorkshopDTO>>(`/workshops?${query}`),
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
    workshops: workshopPage.items.map((item): Workshop => ({
      slug: item.id,
      name: item.name,
      description: item.description,
      wilaya: item.wilaya,
      location: item.location,
      craft: item.craft,
      artisanId: item.artisanId,
      artisanName: item.artisanName,
      image: image(item.media),
      productCount: item.productCount,
    })),
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

export async function catalogueWorkshop(id: string, locale: Locale) {
  const item = await get<WorkshopDTO>(
    `/workshops/${encodeURIComponent(id)}?locale=${locale}`,
  );
  return {
    slug: item.id,
    name: item.name,
    description: item.description,
    wilaya: item.wilaya,
    location: item.location,
    craft: item.craft,
    artisanId: item.artisanId,
    artisanName: item.artisanName,
    image: image(item.media),
    productCount: item.productCount,
  } satisfies Workshop;
}
