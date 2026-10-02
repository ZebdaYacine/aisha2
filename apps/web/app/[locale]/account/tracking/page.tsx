import { notFound } from "next/navigation";

import { EmptyState } from "@/core/components/feedback/empty-state";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
import { orders } from "@/features/catalogue/data";

export default async function TrackingPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const copy = storeCopy(locale);

  return (
    <>
      <p className="text-xs uppercase tracking-widest text-primary">{copy.account}</p>
      <h1 className="mt-3 font-serif text-5xl">{copy.tracking}</h1>
      {orders.length ? (
        <div className="mt-10 divide-y divide-border border-y border-border">
          {orders.map((order) => (
            <article className="grid gap-4 py-6 sm:grid-cols-[1fr_auto] sm:items-center" key={order.id}>
              <div>
                <p className="font-medium">{copy.order} {order.id}</p>
                <p className="mt-2 text-sm text-muted-foreground">{copy[order.status]}</p>
              </div>
              <LocalizedLink locale={locale} href={`/account/orders/${order.id}`} className="inline-flex min-h-11 items-center self-start border-b border-foreground pb-1 text-sm sm:self-auto">
                {copy.track}
              </LocalizedLink>
            </article>
          ))}
        </div>
      ) : (
        <EmptyState title={copy.emptyOrders} body={copy.emptyOrdersBody} />
      )}
    </>
  );
}
