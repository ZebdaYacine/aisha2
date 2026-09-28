"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { formatFullDateTime } from "@/core/lib/format";
import type { Locale } from "@/core/lib/i18n";

type Category = { id: string; slug: string; displayName: string; translations: Record<string, string>; benefitRateBasisPoints: number; isActive: boolean; updatedAt: string };
const locales = ["en", "fr", "ar", "es"] as const;
const emptyForm = { slug: "", displayName: "", translations: { en: "", fr: "", ar: "", es: "" }, benefitRate: "0", isActive: "true" };

export function CategoryManagement({ locale = "en" }: { locale?: Locale }) {
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
      if (!response.ok) throw new Error("Unable to load categories");
      setItems(((await response.json()) as { items?: Category[] }).items ?? []);
    } catch (error) { toast.error(error instanceof Error ? error.message : "Unable to load categories"); }
    finally { setLoading(false); }
  }, []);
  useEffect(() => { void load(); }, [load]);

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
      toast.success(selected ? "Category updated" : "Category created"); setSelected(null); setEditorOpen(false); await load();
    } catch (error) { toast.error(error instanceof Error ? error.message : "Unable to save category"); }
    finally { setSaving(false); }
  };
  const deactivate = async (item: Category) => {
    if (!window.confirm(`Deactivate ${item.displayName}?`)) return;
    const response = await fetch(`/api/admin/categories/${item.id}`, { method: "DELETE" });
    if (!response.ok) { toast.error("Unable to deactivate category"); return; }
    toast.success("Category deactivated"); await load();
  };

  return <section className="space-y-6">
    <div className="flex flex-wrap items-end justify-between gap-4"><div><p className="text-xs uppercase tracking-widest text-primary">Catalogue configuration</p><h2 className="mt-2 font-serif text-3xl">Categories and benefit rates</h2><p className="mt-2 text-sm text-muted-foreground">Maintain the four public language labels and the operational benefit rate.</p></div><Button type="button" onClick={() => open()}>Add category</Button></div>
    {loading ? <p className="border border-border p-5 text-sm text-muted-foreground" role="status">Loading categories…</p> : <div className="overflow-x-auto border border-border"><table className="w-full min-w-[62rem] text-start text-sm"><thead className="border-b border-border bg-muted/40"><tr>{["Category", "English", "French", "Arabic", "Spanish", "Benefit", "Status", "Updated", "Actions"].map((label) => <th className="px-4 py-3 font-medium" key={label}>{label}</th>)}</tr></thead><tbody className="divide-y divide-border">{items.map((item) => <tr key={item.id}><td className="px-4 py-4"><strong>{item.displayName}</strong><span className="block text-xs text-muted-foreground">{item.slug}</span></td><td className="px-4 py-4">{item.translations.en}</td><td className="px-4 py-4">{item.translations.fr}</td><td className="px-4 py-4">{item.translations.ar}</td><td className="px-4 py-4">{item.translations.es}</td><td className="px-4 py-4">{(item.benefitRateBasisPoints / 100).toFixed(2)}%</td><td className="px-4 py-4"><StatusBadge status={item.isActive ? "ACTIVE" : "INACTIVE"} /></td><td className="whitespace-nowrap px-4 py-4">{formatFullDateTime(item.updatedAt, locale)}</td><td className="px-4 py-4"><div className="flex gap-2"><Button type="button" variant="outline" onClick={() => open(item)}>Update</Button>{item.isActive && <Button type="button" variant="destructive" onClick={() => void deactivate(item)}>Deactivate</Button>}</div></td></tr>)}{items.length === 0 && <tr><td className="px-4 py-8 text-muted-foreground" colSpan={9}>No categories found.</td></tr>}</tbody></table></div>}
    {editorOpen ? <div className="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-foreground/50 p-4" role="dialog" aria-modal="true"><form onSubmit={save} className="max-h-[calc(100svh-2rem)] w-full max-w-2xl overflow-y-auto border border-border bg-background p-6"><div className="flex items-start justify-between gap-4"><div><p className="text-xs uppercase tracking-widest text-primary">Category editor</p><h3 className="mt-2 font-serif text-2xl">{selected ? "Update category" : "Add category"}</h3></div><Button type="button" variant="ghost" onClick={() => { setSelected(null); setEditorOpen(false); }}>Close</Button></div><div className="mt-6 grid gap-4 sm:grid-cols-2"><label className="block text-sm"><span className="mb-2 block">Slug</span><input className="auth-input w-full" required value={form.slug} onChange={(event) => setForm({ ...form, slug: event.target.value })} /></label><label className="block text-sm"><span className="mb-2 block">Default display name</span><input className="auth-input w-full" required value={form.displayName} onChange={(event) => setForm({ ...form, displayName: event.target.value })} /></label><label className="block text-sm"><span className="mb-2 block">Benefit rate (%)</span><input className="auth-input w-full" type="number" min="0" max="100" step="0.01" required value={form.benefitRate} onChange={(event) => setForm({ ...form, benefitRate: event.target.value })} /></label><label className="block text-sm"><span className="mb-2 block">Status</span><Combobox options={[{ value: "true", label: "Active" }, { value: "false", label: "Inactive" }]} value={form.isActive} onChange={(value) => setForm({ ...form, isActive: value })} ariaLabel="Status" /></label></div><div className="mt-5 grid gap-4 sm:grid-cols-2">{locales.map((item) => <label className="block text-sm" key={item}><span className="mb-2 block">Name ({item})</span><input className="auth-input w-full" required value={form.translations[item]} onChange={(event) => setForm({ ...form, translations: { ...form.translations, [item]: event.target.value } })} /></label>)}</div><div className="mt-6 flex justify-end gap-2"><Button type="button" variant="outline" onClick={() => { setSelected(null); setEditorOpen(false); }}>Cancel</Button><Button type="submit" disabled={saving}>{saving ? "Saving…" : "Save category"}</Button></div></form></div> : null}
  </section>;
}
