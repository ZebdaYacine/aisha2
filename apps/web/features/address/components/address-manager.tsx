/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { EmptyState } from "@/core/components/feedback/empty-state";
import { ErrorState } from "@/core/components/feedback/error-state";
import { Button } from "@/core/components/ui/button";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";

type Address = { id: string; fullName: string; phone: string; line1: string; line2: string; city: string; postalCode: string; country: string; isDefault: boolean };
type Input = Omit<Address, "id">;
type LoadState = "loading" | "ready" | "error";

const empty: Input = { fullName: "", phone: "", line1: "", line2: "", city: "", postalCode: "", country: "", isDefault: false };
const fields: (keyof Omit<Input, "isDefault">)[] = ["fullName", "phone", "line1", "line2", "city", "postalCode", "country"];
const labels = {
  en: { default: "Default", edit: "Edit", delete: "Delete", update: "Update", add: "Add", load: "Unable to load addresses", save: "Unable to save address", remove: "Unable to delete address" },
  fr: { default: "Par défaut", edit: "Modifier", delete: "Supprimer", update: "Mettre à jour", add: "Ajouter", load: "Impossible de charger les adresses", save: "Impossible d’enregistrer l’adresse", remove: "Impossible de supprimer l’adresse" },
  ar: { default: "افتراضي", edit: "تعديل", delete: "حذف", update: "تحديث", add: "إضافة", load: "تعذر تحميل العناوين", save: "تعذر حفظ العنوان", remove: "تعذر حذف العنوان" },
  es: { default: "Predeterminada", edit: "Editar", delete: "Eliminar", update: "Actualizar", add: "Añadir", load: "No se pudieron cargar las direcciones", save: "No se pudo guardar la dirección", remove: "No se pudo eliminar la dirección" },
} as const;

export function AddressManager({ copy, locale }: { copy: StoreCopy; locale: Locale }) {
  const text = labels[locale];
  const [items, setItems] = useState<Address[]>([]);
  const [form, setForm] = useState<Input>(empty);
  const [editing, setEditing] = useState<string>();
  const [loadState, setLoadState] = useState<LoadState>("loading");
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async () => {
    setLoadState("loading");
    try {
      const response = await fetch("/api/customer/addresses", { cache: "no-store" });
      if (!response.ok) throw new Error("address request failed");
      setItems(await response.json());
      setLoadState("ready");
    } catch {
      setLoadState("error");
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  const save = async (event: React.FormEvent) => {
    event.preventDefault();
    setSubmitting(true);
    try {
      const response = await fetch(editing ? `/api/customer/addresses/${editing}` : "/api/customer/addresses", {
        method: editing ? "PATCH" : "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify(form),
      });
      if (!response.ok) {
        const body = await response.json().catch(() => ({}));
        toast.error(body.error?.message ?? text.save);
        return;
      }
      setForm(empty);
      setEditing(undefined);
      toast.success(copy.addresses);
      await load();
    } catch {
      toast.error(text.save);
    } finally {
      setSubmitting(false);
    }
  };

  const remove = async (id: string) => {
    try {
      const response = await fetch(`/api/customer/addresses/${id}`, { method: "DELETE" });
      if (!response.ok) {
        toast.error(text.remove);
        return;
      }
      if (editing === id) {
        setEditing(undefined);
        setForm(empty);
      }
      await load();
    } catch {
      toast.error(text.remove);
    }
  };

  if (loadState === "loading") {
    return <p className="mt-8 border border-border p-6 text-sm text-muted-foreground" role="status">{copy.accountLoading}</p>;
  }

  if (loadState === "error") {
    return <ErrorState title={copy.accountUnavailable} body={text.load} action={<Button variant="outline" type="button" onClick={() => void load()}>{copy.retry}</Button>} />;
  }

  return (
    <div className="mt-8 grid gap-8 xl:grid-cols-2">
      <div className="space-y-4">
        {items.length ? items.map((item) => (
          <article key={item.id} className="border border-border bg-card p-5">
            <div className="flex justify-between gap-4"><strong>{item.fullName}</strong>{item.isDefault && <span className="text-xs text-primary">{text.default}</span>}</div>
            <address className="mt-2 not-italic text-muted-foreground">{item.line1}<br />{item.line2 && <>{item.line2}<br /></>}{item.postalCode} {item.city}<br />{item.country}</address>
            <div className="mt-4 flex gap-2"><Button variant="outline" type="button" onClick={() => { setEditing(item.id); setForm(item); }}>{text.edit}</Button><Button variant="outline" type="button" onClick={() => void remove(item.id)}>{text.delete}</Button></div>
          </article>
        )) : <EmptyState title={copy.emptyAddresses} body={copy.emptyAddressesBody} />}
      </div>
      <form onSubmit={save} className="space-y-3 border border-border p-5" aria-busy={submitting}>
        {fields.map((field) => <label className="block" key={field}><span className="mb-1 block text-sm">{copy[field]}</span><input required={field !== "phone" && field !== "line2"} className="auth-input" value={form[field]} onChange={(event) => setForm({ ...form, [field]: event.target.value })} /></label>)}
        <label className="flex gap-2"><input type="checkbox" checked={form.isDefault} onChange={(event) => setForm({ ...form, isDefault: event.target.checked })} />{text.default}</label>
        <Button disabled={submitting} type="submit">{editing ? text.update : text.add} {copy.address}</Button>
      </form>
    </div>
  );
}
