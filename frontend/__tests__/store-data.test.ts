import { existsSync } from "node:fs";
import { join } from "node:path";

import { artisans, categories, products } from "@/features/storefront/data";

describe("store data", () => {
  it("uses integer minor units and four complete locales", () => {
    for (const product of products) {
      expect(Number.isInteger(product.priceMinor)).toBe(true);
      expect(product.priceMinor).toBeGreaterThanOrEqual(0);
      expect(Object.keys(product.name).sort()).toEqual(["ar", "en", "es", "fr"]);
    }

    expect(categories.length).toBeGreaterThan(3);
  });

  it("uses existing craft imagery instead of the editorial hero for catalogue records", () => {
    const catalogueImages = [
      ...categories.map((category) => category.image),
      ...artisans.map((artisan) => artisan.image),
      ...products.flatMap((product) => product.images),
    ];

    for (const image of catalogueImages) {
      expect(image).toMatch(/^\/images\/aisha\//);
      expect(image).not.toBe("/images/aisha-hero-editorial.png");
      expect(existsSync(join(process.cwd(), "public", image))).toBe(true);
    }
  });

  it("gives every product imagery assigned to its category", () => {
    const categoryImages = new Map(
      categories.map((category) => [category.slug, category.image]),
    );

    for (const product of products) {
      expect(product.images[0]).toBe(categoryImages.get(product.categorySlug));
    }
  });
});
