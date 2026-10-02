/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useCallback, useEffect, useState } from "react";
import { Eye } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { IconAction } from "@/core/components/ui/icon-action";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { formatFullDateTime } from "@/core/lib/format";
import { commonCopy } from "@/core/lib/common-copy";
import type { Locale } from "@/core/lib/i18n";

type Balance = {
  productId: string;
  productCode: string;
  productName: string;
  workshopId: string;
  workshopName: string;
  artisanName: string;
  onHand: number;
  available: number;
  reserved: number;
  quarantined: number;
  damaged: number;
  rejected: number;
  shipped: number;
  updatedAt: string;
};

export function InventoryOperations({ locale = "en" }: { locale?: Locale }) {
  const common = commonCopy(locale);
  const text = {
    ledger: locale === "fr" ? "Registre des stocks" : locale === "ar" ? "سجل المخزون" : locale === "es" ? "Libro de inventario" : "Inventory ledger",
    balances: locale === "fr" ? "Soldes du stock" : locale === "ar" ? "أرصدة المخزون" : locale === "es" ? "Saldos de stock" : "Stock balances",
    accepted: locale === "fr" ? "Le stock accepté est disponible ; les réservations sont conservées pendant le paiement." : locale === "ar" ? "المخزون المقبول متاح؛ وتُحجز الكميات أثناء الدفع." : locale === "es" ? "El stock aceptado está disponible; las reservas se mantienen durante el pago." : "Accepted stock is available; reservations are held during checkout.",
    allWorkshops: locale === "fr" ? "Tous les ateliers" : locale === "ar" ? "كل الورشات" : locale === "es" ? "Todos los talleres" : "All workshops",
    inventory: locale === "fr" ? "Inventaire" : locale === "ar" ? "المخزون" : locale === "es" ? "Inventario" : "Inventory",
    available: locale === "fr" ? "Disponible" : locale === "ar" ? "متاح" : locale === "es" ? "Disponible" : "Available",
    reserved: locale === "fr" ? "Réservé" : locale === "ar" ? "محجوز" : locale === "es" ? "Reservado" : "Reserved",
    onHand: locale === "fr" ? "En stock" : locale === "ar" ? "في المخزون" : locale === "es" ? "En mano" : "On hand",
    exceptions: locale === "fr" ? "Exceptions" : locale === "ar" ? "استثناءات" : locale === "es" ? "Excepciones" : "Exceptions",
    inventoryDetails: locale === "fr" ? "Détails de l’inventaire" : locale === "ar" ? "تفاصيل المخزون" : locale === "es" ? "Detalles del inventario" : "Inventory details",
  };
  const [items, setItems] = useState<Balance[]>([]);
  const [selected, setSelected] = useState<Balance | null>(null);
  useEscapeKey(() => setSelected(null), selected !== null);
  const [workshopId, setWorkshopId] = useState("");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState({ quantityDelta: "", reason: "", referenceKey: "" });
  const pageSize = 10;
  const pageCount = Math.max(1, Math.ceil(total / pageSize));

  const load = useCallback(async (nextPage = page) => {
    setLoading(true);
    try {
      const query = new URLSearchParams({ page: String(nextPage), pageSize: String(pageSize) });
      if (workshopId) query.set("workshopId", workshopId);
      const response = await fetch(`/api/warehouse/inventory?${query}`, { cache: "no-store" });
      if (!response.ok) throw new Error(common.noResults);
      const data = (await response.json()) as { items?: Balance[]; page?: number; total?: number };
      setItems(data.items ?? []);
      setPage(data.page ?? nextPage);
      setTotal(data.total ?? 0);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : common.noResults);
    } finally {
      setLoading(false);
    }
  }, [common.noResults, page, workshopId]);

  useEffect(() => { void load(); }, [load]);

  const workshopOptions = Array.from(
    new Map(items.map((item) => [item.workshopId, item.workshopName])).entries(),
  ).map(([value, label]) => ({ value, label }));

  const adjust = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!selected) return;
    setSaving(true);
    try {
      const response = await fetch(`/api/warehouse/inventory/${selected.productId}/adjust`, {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ ...form, quantityDelta: Number(form.quantityDelta) }),
      });
      if (!response.ok) {
        const body = await response.json().catch(() => ({}));
        throw new Error(body.error?.message ?? common.noResults);
      }
      toast.success(`${text.inventory} ✓`);
      setSelected(null);
      setForm({ quantityDelta: "", reason: "", referenceKey: "" });
      await load();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : common.noResults);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-8">
      <div className="border border-border p-5">
        <p className="text-xs uppercase tracking-widest text-primary">{text.ledger}</p>
        <h2 className="mt-2 font-serif text-3xl">{text.balances}</h2>
        <p className="mt-2 text-sm text-muted-foreground">{text.accepted}</p>
      </div>
      <div className="flex flex-wrap items-end justify-between gap-3">
        <label className="block text-sm">
          <span className="mb-2 block">{common.workshop}</span>
          <Combobox className="min-w-72" options={[{ value: "", label: text.allWorkshops }, ...workshopOptions]} value={workshopId} onChange={(value) => { setWorkshopId(value); setPage(1); }} ariaLabel={common.workshop} placeholder={text.allWorkshops} emptyMessage={common.noResults} disabled={loading && !workshopOptions.length} />
        </label>
        <span className="text-sm text-muted-foreground">{total} product balance(s)</span>
      </div>
      {loading ? <p className="border border-border p-5 text-sm text-muted-foreground" role="status">{common.loading}</p> : <>
        <div className="overflow-x-auto border border-border">
          <table className="w-full min-w-[74rem] text-start text-sm">
            <thead className="border-b border-border bg-muted/40"><tr>{[common.product, common.workshop, text.available, text.reserved, text.onHand, text.exceptions, common.update, common.details].map((label) => <th key={label} className="px-4 py-3 font-medium">{label}</th>)}</tr></thead>
            <tbody className="divide-y divide-border">
              {items.map((item) => <tr key={item.productId}>
                <td className="px-4 py-4">{item.productName}<br /><span className="text-xs text-muted-foreground">{item.productCode || "Approved product"}</span></td>
                <td className="px-4 py-4">{item.workshopName}<br /><span className="text-xs text-muted-foreground">{item.artisanName}</span></td>
                <td className="px-4 py-4"><StatusBadge status={item.available > 0 ? "AVAILABLE" : "OUT_OF_STOCK"} locale={locale} /><div className="mt-1 text-lg">{item.available}</div></td>
                <td className="px-4 py-4">{item.reserved}</td><td className="px-4 py-4">{item.onHand}</td>
                <td className="px-4 py-4 text-xs">Q {item.quarantined} · D {item.damaged}<br />R {item.rejected} · S {item.shipped}</td>
                <td className="px-4 py-4">{formatFullDateTime(item.updatedAt, locale)}</td>
                <td className="px-4 py-4"><IconAction icon={<Eye size={17} />} label={common.details} onClick={() => setSelected(item)} /></td>
              </tr>)}
              {items.length === 0 && <tr><td className="px-4 py-8 text-muted-foreground" colSpan={8}>{common.noResults}</td></tr>}
            </tbody>
          </table>
        </div>
        <div className="flex items-center justify-between text-sm"><span>{common.page} {page} {common.of} {pageCount}</span><div className="flex gap-2"><Button type="button" variant="outline" disabled={page <= 1} onClick={() => void load(page - 1)}>{common.previous}</Button><Button type="button" variant="outline" disabled={page >= pageCount} onClick={() => void load(page + 1)}>{common.next}</Button></div></div>
      </>}
      {selected && <div className="fixed inset-0 z-50 flex items-center justify-center overflow-x-auto overflow-y-auto p-4" role="dialog" aria-modal="true"><div className="max-h-[calc(100svh-2rem)] w-full max-w-xl overflow-x-auto overflow-y-auto bg-background p-6 shadow-xl">
        <div className="flex items-start justify-between gap-4"><div><p className="text-xs uppercase tracking-widest text-primary">{text.inventoryDetails}</p><h3 className="mt-2 font-serif text-2xl">{selected.productName}</h3></div><Button type="button" variant="ghost" onClick={() => setSelected(null)}>{common.close}</Button></div>
        <dl className="mt-6 grid grid-cols-2 gap-4 text-sm"><div><dt className="text-muted-foreground">Available</dt><dd className="text-2xl">{selected.available}</dd></div><div><dt className="text-muted-foreground">Reserved</dt><dd className="text-2xl">{selected.reserved}</dd></div><div><dt className="text-muted-foreground">Quarantined</dt><dd>{selected.quarantined}</dd></div><div><dt className="text-muted-foreground">Damaged / rejected</dt><dd>{selected.damaged} / {selected.rejected}</dd></div></dl>
        <form onSubmit={adjust} className="mt-6 space-y-4 border-t border-border pt-6"><p className="font-medium">Append adjustment</p><div className="grid gap-4 md:grid-cols-2"><label className="block text-sm"><span className="mb-2 block">Quantity delta</span><input className="auth-input w-full" required type="number" value={form.quantityDelta} onChange={(event) => setForm({ ...form, quantityDelta: event.target.value })} placeholder="+10 or -2" /></label><label className="block text-sm"><span className="mb-2 block">Reference key</span><input className="auth-input w-full" required value={form.referenceKey} onChange={(event) => setForm({ ...form, referenceKey: event.target.value })} /></label></div><label className="block text-sm"><span className="mb-2 block">Reason</span><textarea className="auth-input min-h-24 w-full" required value={form.reason} onChange={(event) => setForm({ ...form, reason: event.target.value })} /></label><Button type="submit" disabled={saving}>{saving ? "Saving…" : "Record adjustment"}</Button></form>
      </div></div>}
    </div>
  );
}
