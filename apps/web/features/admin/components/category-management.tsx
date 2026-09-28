"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { formatFullDateTime } from "@/core/lib/format";
import { commonCopy } from "@/core/lib/common-copy";
import type { Locale } from "@/core/lib/i18n";
import {
  AdminTablePanel,
  AdminTableScroll,
  adminTableCellClass,
  adminTableClass,
  adminTableHeadClass,
} from "@/core/components/admin/admin-table";

type Category = { id: string; slug: string; displayName: string; translations: Record<string, string>; benefitRateBasisPoints: number; isActive: boolean; updatedAt: string };
const locales = ["en", "fr", "ar", "es"] as const;
const emptyForm = { slug: "", displayName: "", translations: { en: "", fr: "", ar: "", es: "" }, benefitRate: "0", isActive: "true" };

export function CategoryManagement({ locale = "en" }: { locale?: Locale }) {
  const common = commonCopy(locale);
  const text = {
    eyebrow: locale === "fr" ? "Configuration du catalogue" : locale === "ar" ? "إعداد الكتالوج" : locale === "es" ? "Configuración del catálogo" : "Catalogue configuration",
    title: locale === "fr" ? "Catégories et taux de bénéfice" : locale === "ar" ? "الفئات ونسب الفائدة" : locale === "es" ? "Categorías y tasas de beneficio" : "Categories and benefit rates",
    description: locale === "fr" ? "Gérez les libellés publics dans les quatre langues et le taux opérationnel." : locale === "ar" ? "أدر التسميات العامة باللغات الأربع ونسبة الفائدة التشغيلية." : locale === "es" ? "Gestiona las etiquetas públicas en cuatro idiomas y la tasa operativa." : "Maintain the four public language labels and the operational benefit rate.",
    add: locale === "fr" ? "Ajouter une catégorie" : locale === "ar" ? "إضافة فئة" : locale === "es" ? "Añadir categoría" : "Add category",
    benefit: locale === "fr" ? "Bénéfice" : locale === "ar" ? "الفائدة" : locale === "es" ? "Beneficio" : "Benefit",
    updated: locale === "fr" ? "Modifiée" : locale === "ar" ? "آخر تحديث" : locale === "es" ? "Actualizada" : "Updated",
  };
  const [items, setItems] = useState<Category[]>([]);
  const [selected, setSelected] = useState<Category | null>(null);
  const [editorOpen, setEditorOpen] = useState(false);
  const [form, setForm] = useState(emptyForm);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  useEscapeKey(() => setSelected(null), selected !== null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const response = await fetch("/api/admin/categories?page=1&pageSize=100", { cache: "no-store" });
      if (!response.ok) throw new Error(common.noResults);
      setItems(((await response.json()) as { items?: Category[] }).items ?? []);
    } catch (error) { toast.error(error instanceof Error ? error.message : common.noResults); }
    finally { setLoading(false); }
  }, [common.noResults]);
  useEffect(() => {
    const timer = window.setTimeout(() => { void load(); }, 0);
    return () => window.clearTimeout(timer);
  }, [load]);

  const open = (item?: Category) => {
    setSelected(item ?? null); setEditorOpen(true);
    setForm(item ? { slug: item.slug, displayName: item.displayName, translations: { en: item.translations.en ?? "", fr: item.translations.fr ?? "", ar: item.translations.ar ?? "", es: item.translations.es ?? "" }, benefitRate: String(item.benefitRateBasisPoints / 100), isActive: String(item.isActive) } : emptyForm);
  };
  const save = async (event: React.FormEvent) => {
    event.preventDefault(); setSaving(true);
    try {
      const payload = { slug: form.slug, displayName: form.displayName, translations: form.translations, benefitRateBasisPoints: Math.round(Number(form.benefitRate) * 100), isActive: form.isActive === "true" };
      const response = await fetch(selected ? `/api/admin/categories/${selected.id}` : "/api/admin/categories", { method: selected ? "PATCH" : "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(payload) });
      const body = await response.json().catch(() => ({}));
      if (!response.ok) throw new Error(body.error?.message ?? "Unable to save category");
      toast.success(selected ? `${common.update} ✓` : `${common.save} ✓`); setSelected(null); setEditorOpen(false); await load();
    } catch (error) { toast.error(error instanceof Error ? error.message : common.noResults); }
    finally { setSaving(false); }
  };
  return <section className="space-y-6">
    <AdminTablePanel eyebrow={text.eyebrow} title={text.title} description={text.description} action={<Button type="button" onClick={() => open()}>{text.add}</Button>} summary={`${items.length} ${common.product}`}>
      {loading ? <p className="border border-border p-5 text-sm text-muted-foreground" role="status">{common.loading}</p> : <AdminTableScroll><table className={`${adminTableClass} min-w-[62rem]`}><thead className={adminTableHeadClass}><tr>{[common.product, "English", "Français", "العربية", "Español", text.benefit, common.status, text.updated, common.actions].map((label) => <th className="px-4 py-3 font-medium" key={label}>{label}</th>)}</tr></thead><tbody className="divide-y divide-border">{items.map((item) => <tr className="transition-colors hover:bg-muted/30" key={item.id}><td className={adminTableCellClass}><strong>{item.displayName}</strong><span className="block text-xs text-muted-foreground">{item.slug}</span></td><td className={adminTableCellClass}>{item.translations.en}</td><td className={adminTableCellClass}>{item.translations.fr}</td><td className={adminTableCellClass}>{item.translations.ar}</td><td className={adminTableCellClass}>{item.translations.es}</td><td className={adminTableCellClass}>{(item.benefitRateBasisPoints / 100).toFixed(2)}%</td><td className={adminTableCellClass}><StatusBadge status={item.isActive ? "ACTIVE" : "INACTIVE"} locale={locale} /></td><td className={`${adminTableCellClass} whitespace-nowrap`}>{formatFullDateTime(item.updatedAt, locale)}</td><td className={adminTableCellClass}><Button type="button" variant="outline" onClick={() => open(item)}>{common.details}</Button></td></tr>)}{items.length === 0 && <tr><td className="p-6 text-sm text-muted-foreground" colSpan={9}>{common.noResults}</td></tr>}</tbody></table></AdminTableScroll>}
    </AdminTablePanel>
    {editorOpen ? <div className="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-foreground/50 p-4" role="dialog" aria-modal="true"><form onSubmit={save} className="max-h-[calc(100svh-2rem)] w-full max-w-2xl overflow-y-auto border border-border bg-background p-6"><div className="flex items-start justify-between gap-4"><div><p className="text-xs uppercase tracking-widest text-primary">Category editor</p><h3 className="mt-2 font-serif text-2xl">{selected ? "Update category" : "Add category"}</h3></div><Button type="button" variant="ghost" onClick={() => { setSelected(null); setEditorOpen(false); }}>Close</Button></div><div className="mt-6 grid gap-4 sm:grid-cols-2"><label className="block text-sm"><span className="mb-2 block">Slug</span><input className="auth-input w-full" required value={form.slug} onChange={(event) => setForm({ ...form, slug: event.target.value })} /></label><label className="block text-sm"><span className="mb-2 block">Default display name</span><input className="auth-input w-full" required value={form.displayName} onChange={(event) => setForm({ ...form, displayName: event.target.value })} /></label><label className="block text-sm"><span className="mb-2 block">Benefit rate (%)</span><input className="auth-input w-full" type="number" min="0" max="100" step="0.01" required value={form.benefitRate} onChange={(event) => setForm({ ...form, benefitRate: event.target.value })} /></label><label className="block text-sm"><span className="mb-2 block">Status</span><Combobox options={[{ value: "true", label: "Active" }, { value: "false", label: "Inactive" }]} value={form.isActive} onChange={(value) => setForm({ ...form, isActive: value })} ariaLabel="Status" /></label></div><div className="mt-5 grid gap-4 sm:grid-cols-2">{locales.map((item) => <label className="block text-sm" key={item}><span className="mb-2 block">Name ({item})</span><input className="auth-input w-full" required value={form.translations[item]} onChange={(event) => setForm({ ...form, translations: { ...form.translations, [item]: event.target.value } })} /></label>)}</div><div className="mt-6 flex justify-end gap-2"><Button type="button" variant="outline" onClick={() => { setSelected(null); setEditorOpen(false); }}>Cancel</Button><Button type="submit" disabled={saving}>{saving ? "Saving…" : "Save category"}</Button></div></form></div> : null}
  </section>;
}
