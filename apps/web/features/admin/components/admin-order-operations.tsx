"use client";

import { useCallback, useEffect, useState } from "react";
import { Eye } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { IconAction } from "@/core/components/ui/icon-action";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { formatFullDateTime } from "@/core/lib/format";
import { commonCopy } from "@/core/lib/common-copy";
import type { Locale } from "@/core/lib/i18n";

type Order = { orderNumber: string; customerName: string; customerEmail: string; status: string; currency: string; totalMinor: number; createdAt: string };
type Page = { items: Order[]; page: number; pageSize: number; total: number };

export function AdminOrderOperations({ locale }: { locale: Locale }) {
  const common = commonCopy(locale);
  const text = {
    eyebrow: locale === "fr" ? "Opérations commerciales" : locale === "ar" ? "عمليات التجارة" : locale === "es" ? "Operaciones comerciales" : "Commerce operations",
    title: locale === "fr" ? "Commandes" : locale === "ar" ? "الطلبات" : locale === "es" ? "Pedidos" : "Orders",
    description: locale === "fr" ? "Consultez les commandes clients et ouvrez leurs détails." : locale === "ar" ? "راجع طلبات العملاء وافتح تفاصيلها." : locale === "es" ? "Revisa los pedidos de clientes y abre sus detalles." : "Review customer orders and open their details.",
    order: locale === "fr" ? "Commande" : locale === "ar" ? "الطلب" : locale === "es" ? "Pedido" : "Order",
    customer: locale === "fr" ? "Client" : locale === "ar" ? "العميل" : locale === "es" ? "Cliente" : "Customer",
    total: locale === "fr" ? "Total" : locale === "ar" ? "الإجمالي" : locale === "es" ? "Total" : "Total",
    created: locale === "fr" ? "Créée" : locale === "ar" ? "تاريخ الإنشاء" : locale === "es" ? "Creado" : "Created",
    email: locale === "fr" ? "E-mail" : locale === "ar" ? "البريد الإلكتروني" : "Email",
    orderDetails: locale === "fr" ? "Détails de la commande" : locale === "ar" ? "تفاصيل الطلب" : locale === "es" ? "Detalles del pedido" : "Order details",
  };
  const [items, setItems] = useState<Order[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [selected, setSelected] = useState<Order | null>(null);
  const [loading, setLoading] = useState(true);
  const load = useCallback(async (nextPage: number) => {
    setLoading(true);
    try {
      const response = await fetch(`/api/admin/orders?page=${nextPage}&pageSize=20`, { cache: "no-store" });
      if (!response.ok) throw new Error(common.noResults);
      const body = (await response.json()) as Page;
      setItems(body.items ?? []); setPage(body.page ?? nextPage); setTotal(body.total ?? 0);
    } catch (error) { toast.error(error instanceof Error ? error.message : common.noResults); }
    finally { setLoading(false); }
  }, [common.noResults]);
  useEffect(() => {
    const timer = window.setTimeout(() => { void load(1); }, 0);
    return () => window.clearTimeout(timer);
  }, [load]);
  useEscapeKey(() => setSelected(null), selected !== null);
  const pages = Math.max(1, Math.ceil(total / 20));
  return <section className="space-y-6"><div><p className="text-xs uppercase tracking-widest text-primary">{text.eyebrow}</p><h2 className="mt-2 font-serif text-3xl">{text.title}</h2><p className="mt-2 text-sm text-muted-foreground">{text.description}</p></div>{loading ? <p className="border border-border p-5 text-sm text-muted-foreground" role="status">{common.loading}</p> : <><div className="overflow-x-auto touch-pan-x border border-border"><table className="w-full min-w-[56rem] text-start text-sm"><thead className="border-b border-border bg-muted/40"><tr>{[text.order, text.customer, common.status, text.total, text.created, common.details].map((label) => <th className="px-4 py-3 text-start font-medium" key={label}>{label}</th>)}</tr></thead><tbody className="divide-y divide-border">{items.map((item) => <tr key={item.orderNumber}><td className="px-4 py-4 font-medium">{item.orderNumber}</td><td className="px-4 py-4"><span>{item.customerName || "—"}</span><span className="block text-xs text-muted-foreground">{item.customerEmail}</span></td><td className="px-4 py-4"><StatusBadge status={item.status} locale={locale} /></td><td className="px-4 py-4">{(item.totalMinor / 100).toFixed(2)} {item.currency}</td><td className="whitespace-nowrap px-4 py-4 text-xs text-muted-foreground">{formatFullDateTime(item.createdAt, locale)}</td><td className="px-4 py-4"><IconAction icon={<Eye size={17} />} label={common.details} onClick={() => setSelected(item)} /></td></tr>)}{items.length === 0 && <tr><td className="px-4 py-8 text-muted-foreground" colSpan={6}>{common.noResults}</td></tr>}</tbody></table></div><div className="flex items-center justify-between gap-3"><Button type="button" variant="outline" disabled={page <= 1} onClick={() => void load(page - 1)}>{common.previous}</Button><span className="text-sm text-muted-foreground">{common.page} {page} {common.of} {pages}</span><Button type="button" variant="outline" disabled={page >= pages} onClick={() => void load(page + 1)}>{common.next}</Button></div></>}{selected && <div className="fixed inset-0 z-50 overflow-y-auto bg-foreground/50 p-4" role="dialog" aria-modal="true"><div className="mx-auto mt-8 max-h-[calc(100svh-4rem)] w-full max-w-xl overflow-y-auto border border-border bg-background p-6"><div className="flex items-start justify-between gap-4"><div><p className="text-xs uppercase tracking-widest text-primary">{text.orderDetails}</p><h3 className="mt-2 font-serif text-2xl">{selected.orderNumber}</h3></div><Button type="button" variant="ghost" onClick={() => setSelected(null)}>{common.close}</Button></div><dl className="mt-6 grid gap-4 text-sm sm:grid-cols-2"><div><dt className="text-muted-foreground">{text.customer}</dt><dd className="mt-1">{selected.customerName || "—"}</dd></div><div><dt className="text-muted-foreground">{text.email}</dt><dd className="mt-1 break-all">{selected.customerEmail}</dd></div><div><dt className="text-muted-foreground">{common.status}</dt><dd className="mt-1"><StatusBadge status={selected.status} locale={locale} /></dd></div><div><dt className="text-muted-foreground">{text.total}</dt><dd className="mt-1">{(selected.totalMinor / 100).toFixed(2)} {selected.currency}</dd></div><div><dt className="text-muted-foreground">{text.created}</dt><dd className="mt-1">{formatFullDateTime(selected.createdAt, locale)}</dd></div></dl></div></div>}</section>;
}
