"use client";
import { useEffect, useState } from "react";
import { EmptyState } from "@/core/components/feedback/empty-state";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { formatFullDate } from "@/core/lib/format";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { formatMoney } from "@/features/catalogue/format";
type Order = {
  id: string;
  orderNumber: string;
  status: string;
  totalMinor: number;
  currency: string;
  createdAt: string;
};
export function LiveOrders({
  locale,
  copy,
}: {
  locale: Locale;
  copy: StoreCopy;
}) {
  const [data, setData] = useState<Order[] | null>(null);
  useEffect(() => {
    fetch("/api/orders")
      .then((r) => (r.ok ? r.json() : Promise.reject()))
      .then((v) => setData(v.items ?? []))
      .catch(() => setData([]));
  }, []);
  if (data === null)
    return <p className="mt-10 text-sm text-muted-foreground">Loading…</p>;
  if (!data.length)
    return <EmptyState title={copy.emptyOrders} body={copy.emptyOrdersBody} />;
  return (
    <div className="mt-10 divide-y divide-border border-y border-border">
      {data.map((o) => (
        <article
          className="grid gap-4 py-6 sm:grid-cols-[1fr_auto] sm:items-center"
          key={o.id}
        >
          <div>
            <p className="font-medium">
              {copy.order} {o.orderNumber}
            </p>
            <p className="mt-2 text-sm text-muted-foreground">
              {formatFullDate(o.createdAt, locale)} · <StatusBadge status={o.status} locale={locale} /> · {formatMoney(o.totalMinor, o.currency, locale)}
            </p>
          </div>
          <LocalizedLink
            locale={locale}
            href={`/account/orders/${o.id}`}
            className="border-b border-foreground pb-1 text-sm"
          >
            {copy.track}
          </LocalizedLink>
        </article>
      ))}
    </div>
  );
}
