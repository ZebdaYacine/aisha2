export type ProductLocale = "ar" | "fr" | "en" | "es";

export type ProductTranslation = {
  locale: ProductLocale;
  name: string;
  description: string;
  story: string;
  culturalContext: string;
};

export type ProductDraft = {
  id: string;
  categoryId: string;
  productType: "ARTISAN_SPECIFIC" | "STANDARD_TRADITIONAL";
  status: string;
  priceMinor: number;
  currency: string;
  materials: string;
  productionMethod: string;
  intendedUse: string;
  dimensions: string;
  weightGrams?: number;
  countryOfOrigin: string;
  regionOfOrigin: string;
  ecoFriendlyVerified: boolean;
  fairTradeVerified: boolean;
  madeToOrderEligible: boolean;
  translations: ProductTranslation[];
  media: Array<{ id: string; mediaKind: string; originalFilename: string; mediaType: string; sizeBytes: number; altText: string; url?: string }>;
};

export type ProductInput = Omit<ProductDraft, "id" | "status">;

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, { ...init, cache: "no-store" });
  const body = await response.json().catch(() => undefined);
  if (!response.ok) {
    throw new Error(body?.error?.message ?? "The product request failed.");
  }
  return body as T;
}

export function listOwnedProducts(status = "") {
  const query = status ? `?status=${encodeURIComponent(status)}` : "";
  return request<{ items: ProductDraft[] }>(`/api/artisan/products${query}`);
}

export function createProduct(input: ProductInput) {
  return request<ProductDraft>("/api/artisan/products", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(input) });
}

export function updateProduct(id: string, input: ProductInput) {
  return request<ProductDraft>(`/api/artisan/products/${encodeURIComponent(id)}`, { method: "PATCH", headers: { "content-type": "application/json" }, body: JSON.stringify(input) });
}

export function submitProduct(id: string) {
  return request<ProductDraft>(`/api/artisan/products/${encodeURIComponent(id)}/submit`, { method: "POST" });
}

export function uploadProductMedia(id: string, file: File, altText: string) {
  const form = new FormData();
  form.set("file", file);
  form.set("altText", altText);
  return request<ProductDraft["media"][number]>(`/api/artisan/products/${encodeURIComponent(id)}/media`, { method: "POST", body: form });
}

export function deleteProductMedia(id: string, mediaId: string) {
  return fetch(`/api/artisan/products/${encodeURIComponent(id)}/media/${encodeURIComponent(mediaId)}`, { method: "DELETE" });
}
