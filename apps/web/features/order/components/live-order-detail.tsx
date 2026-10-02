"use client";
import { useEffect, useState } from "react";
import { Check, Clock3 } from "lucide-react";
import { toast } from "sonner";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { formatMoney } from "@/features/catalogue/format";
type Item = {
  productName: string;
  quantity: number;
  unitPriceMinor: number;
  subtotalMinor: number;
  currency: string;
};
type Order = {
  id: string;
  orderNumber: string;
  status: string;
  address: {
    fullName: string;
    line1: string;
    postalCode: string;
    city: string;
    country: string;
  };
  items: Item[];
  totalMinor: number;
  currency: string;
  shipmentEvents: Array<{ id: string; status: string; trackingReference?: string; occurredAt: string }>;
  payment?: { id: string; provider: string; status: string; amountMinor: number; currency: string };
};
export function LiveOrderDetail({
  locale,
  copy,
  id,
}: {
  locale: Locale;
  copy: StoreCopy;
  id: string;
}) {
  const [o, setO] = useState<Order | null>(null);
  const [cancelling, setCancelling] = useState(false);
  const [paying, setPaying] = useState(false);
  useEffect(() => {
    fetch(`/api/orders/${id}`)
      .then((r) => (r.ok ? r.json() : Promise.reject()))
      .then(setO)
      .catch(() => setO(null));
  }, [id]);
  if (!o) return <p className="text-sm text-muted-foreground">Loading…</p>;
  const confirmPayment = async () => {
    if (!o.payment?.id) return;
    setPaying(true);
    try {
      const response = await fetch(`/api/payments/${o.payment.id}/confirm`, {
        method: "POST",
        headers: { "Idempotency-Key": crypto.randomUUID() },
      });
      if (!response.ok) throw new Error();
      setO(await response.json() as Order);
      toast.success(copy.paymentConfirmed);
    } catch {
      toast.error(copy.paymentError);
    } finally {
      setPaying(false);
    }
  };
  return (
    <>
      <p className="text-xs uppercase tracking-widest text-primary">
        {copy.order}
      </p>
      <div className="mt-3 flex flex-wrap items-end justify-between gap-4">
        <h1 className="font-serif text-4xl">{o.orderNumber}</h1>
        <div className="flex flex-wrap items-center gap-4"><StatusBadge status={o.status} locale={locale} />{o.status === "PENDING_PAYMENT" && o.payment?.status === "PENDING" && <button type="button" disabled={paying} className="rounded bg-primary px-4 py-2 text-sm text-primary-foreground disabled:opacity-60" onClick={() => void confirmPayment()}>{paying ? "…" : copy.payNow}</button>}{(o.status === "PENDING_PAYMENT" || o.status === "PAID") && <button type="button" disabled={cancelling || paying} className="text-sm underline" onClick={async () => { setCancelling(true); const response = await fetch(`/api/orders/${o.id}/cancel`, { method: "POST" }); if (response.ok) setO(await response.json() as Order); setCancelling(false); }}>{cancelling ? "…" : copy.cancelOrder}</button>}</div>
      </div>
      <div className="mt-10 grid gap-10 xl:grid-cols-[1fr_20rem]">
        <div>
          <h2 className="font-serif text-2xl">{copy.tracking}</h2>
          <ol className="mt-6 border-s border-border ps-7">
            {(o.shipmentEvents ?? []).map((event, index) => <li className="relative pb-8" key={event.id}><Check className="absolute -start-[2.4rem] top-0 rounded-full bg-success p-1 text-white" size={24} /><p className="font-medium"><StatusBadge status={event.status} locale={locale} /></p><p className="mt-1 break-words text-sm text-muted-foreground">{new Intl.DateTimeFormat(locale, { dateStyle: "medium", timeStyle: "short" }).format(new Date(event.occurredAt))}{event.trackingReference ? ` · ${event.trackingReference}` : ""}</p>{index === (o.shipmentEvents ?? []).length - 1 && <Clock3 className="absolute -start-[2.4rem] top-8 bg-background p-1" size={24} />}</li>)}
            {!(o.shipmentEvents ?? []).length && <li className="relative"><Clock3 className="absolute -start-[2.4rem] top-0 bg-background p-1" size={24} /><p className="text-sm text-muted-foreground">Payment and fulfilment updates appear here.</p></li>}
          </ol>
          <h2 className="mt-6 font-serif text-2xl">{copy.orders}</h2>
          <div className="mt-5 divide-y divide-border">
            {o.items.map((i, n) => (
              <div
                className="flex min-w-0 flex-wrap justify-between gap-3 py-5"
                key={`${i.productName}-${n}`}
              >
                <div className="min-w-0 flex-1">
                  <p className="break-words">{i.productName}</p>
                  <p className="text-sm text-muted-foreground">
                    {copy.quantity}: {i.quantity}
                  </p>
                </div>
                <p className="shrink-0">{formatMoney(i.subtotalMinor, i.currency, locale)}</p>
              </div>
            ))}
          </div>
        </div>
        <aside className="border border-border p-6">
          <h2 className="font-serif text-2xl">{copy.address}</h2>
          <address className="mt-5 text-sm not-italic leading-6 text-muted-foreground">
            {o.address.fullName}
            <br />
            {o.address.line1}
            <br />
            {o.address.postalCode} {o.address.city}
            <br />
            {o.address.country}
          </address>
          <p className="mt-7 border-t border-border pt-5 font-medium">
            {formatMoney(o.totalMinor, o.currency, locale)}
          </p>
        </aside>
      </div>
    </>
  );
}
