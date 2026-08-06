import type { Product } from "@/features/catalogue/types";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { ProductCard, ProductCardSkeleton } from "./product-card";
export function ProductGrid({
  products,
  locale,
  copy,
  className = "",
}: {
  products: Product[];
  locale: Locale;
  copy: StoreCopy;
  className?: string;
}) {
  return (
    <div
      className={`grid grid-cols-2 gap-x-3 gap-y-10 sm:gap-x-5 md:grid-cols-3 lg:grid-cols-4 ${className}`}
    >
      {products.map((product) => (
        <ProductCard
          key={product.slug}
          product={product}
          locale={locale}
          copy={copy}
        />
      ))}
    </div>
  );
}
export function ProductGridSkeleton() {
  return (
    <div className="grid grid-cols-2 gap-3 md:grid-cols-3 lg:grid-cols-4">
      {Array.from({ length: 8 }, (_, i) => (
        <ProductCardSkeleton key={i} />
      ))}
    </div>
  );
}
