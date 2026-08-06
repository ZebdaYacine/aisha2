import { Container } from "@/core/components/layout/container";
import { ProductGridSkeleton } from "@/core/components/commerce/product-grid";
export default function Loading() {
  return (
    <Container className="py-16">
      <div className="h-12 w-72 animate-pulse bg-muted" />
      <div className="mt-16">
        <ProductGridSkeleton />
      </div>
    </Container>
  );
}
