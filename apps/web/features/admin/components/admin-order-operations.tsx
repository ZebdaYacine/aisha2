"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { formatFullDateTime } from "@/core/lib/format";
import type { Locale } from "@/core/lib/i18n";

type Order = { orderNumber: string; customerName: string; customerEmail: string; status: string; currency: string; totalMinor: number; createdAt: string };
type Page = { items: Order[]; page: number; pageSize: number; total: number };

export function AdminOrderOperations({ locale }: { locale: Locale }) {
  const [items, setItems] = useState<Order[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [selected, setSelected] = useState<Order | null>(null);
  const [loading, setLoading] = useState(true);
  const load = useCallback(async (nextPage: number) => {
    setLoading(true);
    try {
      const response = await fetch(`/api/admin/orders?page=${nextPage}&pageSize=20`, { cache: "no-store" });
      if (!response.ok) throw new Error("Unable to load orders");
      const body = (await response.json()) as Page;
      setItems(body.items ?? []); setPage(body.page ?? nextPage); setTotal(body.total ?? 0);
    } catch (error) { toast.error(error instanceof Error ? error.message : "Unable to load orders"); }
    finally { setLoading(false); }
  }, []);
  useEffect(() => { void load(1); }, [load]);
  useEscapeKey(() => setSelected(null), selected !== null);
  const pages = Math.max(1, Math.ceil(total / 20));
  return <section className="space-y-6"><div><p className="text-xs uppercase tracking-widest text-primary">Commerce operations</p><h2 className="mt-2 font-serif text-3xl">Orders</h2><p className="mt-2 text-sm text-muted-foreground">Review customer orders and open a sanitized detail view.</p></div>{loading ? <p className="border border-border p-5 text-sm text-muted-foreground" role="status">Loading orders…</p> : <><div className="overflow-x-auto border border-border"><table className="w-full min-w-[56rem] text-start text-sm"><thead className="border-b border-border bg-muted/40"><tr>{["Order", "Customer", "Status", "Total", "Created", "Details"].map((label) => <th className="px-4 py-3 text-start font-medium" key={label}>{label}</th>)}</tr></thead><tbody className="divide-y divide-border">{items.map((item) => <tr key={item.orderNumber}><td className="px-4 py-4 font-medium">{item.orderNumber}</td><td className="px-4 py-4"><span>{item.customerName || "—"}</span><span className="block text-xs text-muted-foreground">{item.customerEmail}</span></td><td className="px-4 py-4"><StatusBadge status={item.status} /></td><td className="px-4 py-4">{(item.totalMinor / 100).toFixed(2)} {item.currency}</td><td className="whitespace-nowrap px-4 py-4 text-xs text-muted-foreground">{formatFullDateTime(item.createdAt, locale)}</td><td className="px-4 py-4"><Button type="button" variant="outline" onClick={() => setSelected(item)}>Details</Button></td></tr>)}{items.length === 0 && <tr><td className="px-4 py-8 text-muted-foreground" colSpan={6}>No orders found.</td></tr>}</tbody></table></div><div className="flex items-center justify-between gap-3"><Button type="button" variant="outline" disabled={page <= 1} onClick={() => void load(page - 1)}>Previous</Button><span className="text-sm text-muted-foreground">Page {page} of {pages}</span><Button type="button" variant="outline" disabled={page >= pages} onClick={() => void load(page + 1)}>Next</Button></div></>}{selected && <div className="fixed inset-0 z-50 overflow-y-auto bg-foreground/50 p-4" role="dialog" aria-modal="true"><div className="mx-auto mt-8 max-h-[calc(100svh-4rem)] w-full max-w-xl overflow-y-auto border border-border bg-background p-6"><div className="flex items-start justify-between gap-4"><div><p className="text-xs uppercase tracking-widest text-primary">Order details</p><h3 className="mt-2 font-serif text-2xl">{selected.orderNumber}</h3></div><Button type="button" variant="ghost" onClick={() => setSelected(null)}>Close</Button></div><dl className="mt-6 grid gap-4 text-sm sm:grid-cols-2"><div><dt className="text-muted-foreground">Customer</dt><dd className="mt-1">{selected.customerName || "—"}</dd></div><div><dt className="text-muted-foreground">Email</dt><dd className="mt-1 break-all">{selected.customerEmail}</dd></div><div><dt className="text-muted-foreground">Status</dt><dd className="mt-1"><StatusBadge status={selected.status} /></dd></div><div><dt className="text-muted-foreground">Total</dt><dd className="mt-1">{(selected.totalMinor / 100).toFixed(2)} {selected.currency}</dd></div><div><dt className="text-muted-foreground">Created</dt><dd className="mt-1">{formatFullDateTime(selected.createdAt, locale)}</dd></div></dl></div></div>}</section>;
}
