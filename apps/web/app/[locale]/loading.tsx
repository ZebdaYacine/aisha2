import { Container } from "@/core/components/layout/container";
import { Skeleton } from "@/core/components/ui/skeleton";

export default function Loading() {
  return <Container className="py-16" aria-busy="true">
    <Skeleton className="h-10 w-64" />
    <Skeleton className="mt-8 aspect-[16/7]" />
    <div className="mt-10 grid grid-cols-2 gap-4 md:grid-cols-4">
      {Array.from({ length: 4 }, (_, index) => <Skeleton className="aspect-[4/5]" key={index} />)}
    </div>
  </Container>;
}
