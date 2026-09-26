/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { formatFullDateTime } from "@/core/lib/format";
import type { Locale } from "@/core/lib/i18n";

type Balance = {
  productId: string;
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
      if (!response.ok) throw new Error("Unable to load inventory");
      const data = (await response.json()) as { items?: Balance[]; page?: number; total?: number };
      setItems(data.items ?? []);
      setPage(data.page ?? nextPage);
      setTotal(data.total ?? 0);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Unable to load inventory");
    } finally {
      setLoading(false);
    }
  }, [page, workshopId]);

  useEffect(() => { void load(); }, [load]);

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
        throw new Error(body.error?.message ?? "Unable to adjust inventory");
      }
      toast.success("Inventory adjustment recorded");
      setSelected(null);
      setForm({ quantityDelta: "", reason: "", referenceKey: "" });
      await load();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Unable to adjust inventory");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-8">
      <div className="border border-border p-5">
        <p className="text-xs uppercase tracking-widest text-primary">Inventory ledger</p>
        <h2 className="mt-2 font-serif text-3xl">Stock balances</h2>
        <p className="mt-2 text-sm text-muted-foreground">Accepted stock is available; reservations are held atomically during checkout.</p>
      </div>
      <div className="flex flex-wrap items-end justify-between gap-3">
        <label className="block text-sm">
          <span className="mb-2 block">Workshop ID</span>
          <Combobox className="min-w-72" allowCustom options={[]} value={workshopId} onChange={(value) => { setWorkshopId(value); setPage(1); }} ariaLabel="Workshop ID" placeholder="All workshops or paste an ID" />
        </label>
        <span className="text-sm text-muted-foreground">{total} product balance(s)</span>
      </div>
      {loading ? <p className="border border-border p-5 text-sm text-muted-foreground" role="status">Loading inventory…</p> : <>
        <div className="overflow-x-auto border border-border">
          <table className="w-full min-w-[74rem] text-start text-sm">
            <thead className="border-b border-border bg-muted/40"><tr>{["Product", "Workshop", "Available", "Reserved", "On hand", "Exceptions", "Updated", "Details"].map((label) => <th key={label} className="px-4 py-3 font-medium">{label}</th>)}</tr></thead>
            <tbody className="divide-y divide-border">
              {items.map((item) => <tr key={item.productId}>
                <td className="px-4 py-4">{item.productName}<br /><span className="text-xs text-muted-foreground">{item.productId}</span></td>
                <td className="px-4 py-4">{item.workshopName}<br /><span className="text-xs text-muted-foreground">{item.artisanName}</span></td>
                <td className="px-4 py-4"><StatusBadge status={item.available > 0 ? "AVAILABLE" : "OUT_OF_STOCK"} /><div className="mt-1 text-lg">{item.available}</div></td>
                <td className="px-4 py-4">{item.reserved}</td><td className="px-4 py-4">{item.onHand}</td>
                <td className="px-4 py-4 text-xs">Q {item.quarantined} · D {item.damaged}<br />R {item.rejected} · S {item.shipped}</td>
                <td className="px-4 py-4">{formatFullDateTime(item.updatedAt, locale)}</td>
                <td className="px-4 py-4"><Button type="button" variant="outline" onClick={() => setSelected(item)}>Details</Button></td>
              </tr>)}
              {items.length === 0 && <tr><td className="px-4 py-8 text-muted-foreground" colSpan={8}>No inventory balances found.</td></tr>}
            </tbody>
          </table>
        </div>
        <div className="flex items-center justify-between text-sm"><span>Page {page} of {pageCount}</span><div className="flex gap-2"><Button type="button" variant="outline" disabled={page <= 1} onClick={() => void load(page - 1)}>Previous</Button><Button type="button" variant="outline" disabled={page >= pageCount} onClick={() => void load(page + 1)}>Next</Button></div></div>
      </>}
      {selected && <div className="fixed inset-0 z-50 flex items-center justify-center overflow-x-auto overflow-y-auto p-4" role="dialog" aria-modal="true"><div className="max-h-[calc(100svh-2rem)] w-full max-w-xl overflow-x-auto overflow-y-auto bg-background p-6 shadow-xl">
        <div className="flex items-start justify-between gap-4"><div><p className="text-xs uppercase tracking-widest text-primary">Inventory details</p><h3 className="mt-2 font-serif text-2xl">{selected.productName}</h3></div><Button type="button" variant="ghost" onClick={() => setSelected(null)}>Close</Button></div>
        <dl className="mt-6 grid grid-cols-2 gap-4 text-sm"><div><dt className="text-muted-foreground">Available</dt><dd className="text-2xl">{selected.available}</dd></div><div><dt className="text-muted-foreground">Reserved</dt><dd className="text-2xl">{selected.reserved}</dd></div><div><dt className="text-muted-foreground">Quarantined</dt><dd>{selected.quarantined}</dd></div><div><dt className="text-muted-foreground">Damaged / rejected</dt><dd>{selected.damaged} / {selected.rejected}</dd></div></dl>
        <form onSubmit={adjust} className="mt-6 space-y-4 border-t border-border pt-6"><p className="font-medium">Append adjustment</p><div className="grid gap-4 md:grid-cols-2"><label className="block text-sm"><span className="mb-2 block">Quantity delta</span><input className="auth-input w-full" required type="number" value={form.quantityDelta} onChange={(event) => setForm({ ...form, quantityDelta: event.target.value })} placeholder="+10 or -2" /></label><label className="block text-sm"><span className="mb-2 block">Reference key</span><input className="auth-input w-full" required value={form.referenceKey} onChange={(event) => setForm({ ...form, referenceKey: event.target.value })} /></label></div><label className="block text-sm"><span className="mb-2 block">Reason</span><textarea className="auth-input min-h-24 w-full" required value={form.reason} onChange={(event) => setForm({ ...form, reason: event.target.value })} /></label><Button type="submit" disabled={saving}>{saving ? "Saving…" : "Record adjustment"}</Button></form>
      </div></div>}
    </div>
  );
}
