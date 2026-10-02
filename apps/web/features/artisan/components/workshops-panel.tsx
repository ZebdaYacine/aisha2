/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useEffect, useState, type ChangeEvent, type FormEvent } from "react";
import { Eye, Pencil, Plus, Power, Trash2 } from "lucide-react";

import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { IconAction } from "@/core/components/ui/icon-action";
import { Modal } from "@/core/components/ui/modal";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { wilayaOptions } from "@/core/lib/location-options";

type Locale = "en" | "fr" | "ar" | "es";
type Profile = { membershipStatus?: string; status?: string };
type Workshop = {
  id: string;
  name: string;
  description?: string;
  wilaya?: string;
  location?: string;
  status: "ACTIVE" | "INACTIVE" | "ARCHIVED";
  isDefault: boolean;
  isPublic: boolean;
  productCount: number;
};
type WorkshopForm = { name: string; description: string; wilaya: string; location: string; isPublic: boolean };
type Dialog = "create" | "edit" | "details" | "delete" | null;

const copy = {
  en: { title: "Workshops", activate: "Activate artisan membership", activateHelp: "Create your first workshop to activate seller tools.", add: "Add workshop", update: "Update", details: "Details", save: "Save", cancel: "Cancel", name: "Name", description: "Description", wilaya: "Wilaya", location: "Location", public: "Visible in public catalogue", active: "Deactivate", inactive: "Activate", delete: "Delete", deleting: "Deleting…", empty: "No workshops yet.", error: "We could not save this change.", verification: "Verification", products: "Products", visibility: "Visibility", publicValue: "Public", privateValue: "Private", status: "Status", default: "Default", confirmDelete: "Delete this workshop?", deleteHelp: "Only an empty workshop can be permanently deleted.", close: "Close" },
  fr: { title: "Ateliers", activate: "Activer l’adhésion artisan", activateHelp: "Créez votre premier atelier pour activer les outils de vente.", add: "Ajouter un atelier", update: "Modifier", details: "Détails", save: "Enregistrer", cancel: "Annuler", name: "Nom", description: "Description", wilaya: "Wilaya", location: "Adresse", public: "Visible dans le catalogue public", active: "Désactiver", inactive: "Activer", delete: "Supprimer", deleting: "Suppression…", empty: "Aucun atelier pour le moment.", error: "Impossible d’enregistrer cette modification.", verification: "Vérification", products: "Produits", visibility: "Visibilité", publicValue: "Public", privateValue: "Privé", status: "Statut", default: "Par défaut", confirmDelete: "Supprimer cet atelier ?", deleteHelp: "Seul un atelier vide peut être supprimé définitivement.", close: "Fermer" },
  ar: { title: "الورشات", activate: "تفعيل عضوية الحرفي", activateHelp: "أنشئ ورشتك الأولى لتفعيل أدوات البيع.", add: "إضافة ورشة", update: "تحديث", details: "التفاصيل", save: "حفظ", cancel: "إلغاء", name: "الاسم", description: "الوصف", wilaya: "الولاية", location: "العنوان", public: "إظهار في الكتالوج العام", active: "تعطيل", inactive: "تفعيل", delete: "حذف", deleting: "جارٍ الحذف…", empty: "لا توجد ورشات بعد.", error: "تعذر حفظ التغيير.", verification: "التحقق", products: "المنتجات", visibility: "الرؤية", publicValue: "عام", privateValue: "خاص", status: "الحالة", default: "افتراضي", confirmDelete: "حذف هذه الورشة؟", deleteHelp: "يمكن حذف الورشة الفارغة فقط نهائياً.", close: "إغلاق" },
  es: { title: "Talleres", activate: "Activar membresía artesanal", activateHelp: "Crea tu primer taller para activar las herramientas de venta.", add: "Añadir taller", update: "Actualizar", details: "Detalles", save: "Guardar", cancel: "Cancelar", name: "Nombre", description: "Descripción", wilaya: "Wilaya", location: "Dirección", public: "Visible en el catálogo público", active: "Desactivar", inactive: "Activar", delete: "Eliminar", deleting: "Eliminando…", empty: "Aún no hay talleres.", error: "No se pudo guardar el cambio.", verification: "Verificación", products: "Productos", visibility: "Visibilidad", publicValue: "Público", privateValue: "Privado", status: "Estado", default: "Predeterminado", confirmDelete: "¿Eliminar este taller?", deleteHelp: "Solo se puede eliminar permanentemente un taller vacío.", close: "Cerrar" },
} as const;

const blank: WorkshopForm = { name: "", description: "", wilaya: "", location: "", isPublic: true };

function formFrom(workshop?: Workshop): WorkshopForm {
  return workshop ? { name: workshop.name, description: workshop.description ?? "", wilaya: workshop.wilaya ?? "", location: workshop.location ?? "", isPublic: workshop.isPublic } : blank;
}

async function readError(response: Response) {
  try {
    const body = await response.json();
    return body?.error?.message ?? "Request failed";
  } catch {
    return "Request failed";
  }
}

export function WorkshopsPanel({ locale, profile }: { locale: Locale; profile: Profile }) {
  const t = copy[locale];
  const [workshops, setWorkshops] = useState<Workshop[]>([]);
  const [form, setForm] = useState<WorkshopForm>(blank);
  const [focused, setFocused] = useState<Workshop>();
  const [editing, setEditing] = useState<string | null>(null);
  const [pendingDelete, setPendingDelete] = useState<Workshop>();
  const [dialog, setDialog] = useState<Dialog>(null);
  const [loaded, setLoaded] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [verification, setVerification] = useState<string>("");
  const activeMembership = profile.membershipStatus === "ACTIVE";
  const blockedMembership = profile.membershipStatus === "SUSPENDED" || profile.membershipStatus === "CLOSED";

  const load = async () => {
    const [workshopsResponse, verificationResponse] = await Promise.all([fetch("/api/artisan/workshops"), fetch("/api/artisan/verification")]);
    if (workshopsResponse.ok) setWorkshops((await workshopsResponse.json()).items ?? []);
    if (verificationResponse.ok) setVerification((await verificationResponse.json()).status ?? "NOT_SUBMITTED");
    setLoaded(true);
  };

  useEffect(() => { void load(); }, []);

  const closeDialog = (force = false) => {
    if (busy && !force) return;
    setDialog(null);
    setFocused(undefined);
    setPendingDelete(undefined);
    setEditing(null);
    setForm(blank);
  };

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setBusy(true);
    setError("");
    const endpoint = editing ? `/api/artisan/workshops/${editing}` : "/api/artisan/workshops";
    const response = await fetch(endpoint, { method: editing ? "PATCH" : "POST", headers: { "content-type": "application/json", ...(editing ? {} : { "Idempotency-Key": crypto.randomUUID() }) }, body: JSON.stringify(form) });
    if (!response.ok) setError(await readError(response));
    else { closeDialog(true); await load(); }
    setBusy(false);
  };

  const activate = async (event: FormEvent) => {
    event.preventDefault();
    setBusy(true);
    setError("");
    const response = await fetch("/api/artisan/membership/activate", { method: "POST", headers: { "content-type": "application/json", "Idempotency-Key": crypto.randomUUID() }, body: JSON.stringify(form) });
    if (!response.ok) setError(await readError(response));
    else { setForm(blank); await load(); }
    setBusy(false);
  };

  const setStatus = async (workshop: Workshop) => {
    setBusy(true);
    setError("");
    const next = workshop.status === "ACTIVE" ? "INACTIVE" : "ACTIVE";
    const response = await fetch(`/api/artisan/workshops/${workshop.id}/status`, { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify({ status: next }) });
    if (!response.ok) setError(await readError(response)); else await load();
    setBusy(false);
  };

  const remove = async () => {
    if (!pendingDelete) return;
    setBusy(true);
    setError("");
    const response = await fetch(`/api/artisan/workshops/${pendingDelete.id}`, { method: "DELETE" });
    if (!response.ok) setError(await readError(response));
    else { closeDialog(true); await load(); }
    setBusy(false);
  };

  if (!loaded) return <section className="rounded-xl border border-border p-6" aria-busy="true">{t.title}…</section>;

  return (
    <section className="mt-8 rounded-xl border border-border p-6" aria-labelledby="workshops-title">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div><h2 id="workshops-title" className="font-serif text-2xl">{t.title}</h2><p className="mt-1 text-sm text-muted-foreground">{t.verification}: {verification || "—"}</p></div>
        {!activeMembership && <span className="rounded-full bg-muted px-3 py-1 text-xs"><StatusBadge status={profile.membershipStatus ?? "NOT_STARTED"} locale={locale} /></span>}
        {activeMembership && <IconAction icon={<Plus size={18} />} label={t.add} onClick={() => { setForm(blank); setEditing(null); setDialog("create"); }} />}
      </div>
      {!activeMembership && !blockedMembership ? (
        <form className="mt-6 space-y-4" onSubmit={activate}>
          <p className="text-sm text-muted-foreground">{t.activateHelp}</p>
          <WorkshopFields value={form} onChange={setForm} t={t} />
          <Button type="submit" disabled={busy}>{busy ? "…" : t.activate}</Button>
        </form>
      ) : (
        <>
          {workshops.length === 0 ? <p className="mt-6 text-sm text-muted-foreground">{t.empty}</p> : (
            <>
            <div className="mt-6 hidden max-w-full touch-pan-x overflow-x-auto rounded-lg border border-border md:block">
              <table className="w-full min-w-[760px] text-sm">
                <thead className="bg-muted/40 text-start"><tr><th className="px-4 py-3 text-start font-medium">{t.details}</th><th className="px-4 py-3 text-start font-medium">{t.name}</th><th className="px-4 py-3 text-start font-medium">{t.wilaya}</th><th className="px-4 py-3 text-start font-medium">{t.products}</th><th className="px-4 py-3 text-start font-medium">{t.visibility}</th><th className="px-4 py-3 text-start font-medium">{t.status}</th></tr></thead>
                <tbody>{workshops.map((workshop) => <tr className="border-t border-border align-top" key={workshop.id}>
                  <td className="px-4 py-3"><div className="flex flex-wrap gap-2"><IconAction icon={<Eye size={17} />} label={t.details} onClick={() => { setFocused(workshop); setDialog("details"); }} /><IconAction icon={<Pencil size={17} />} label={t.update} disabled={busy} onClick={() => { setEditing(workshop.id); setForm(formFrom(workshop)); setDialog("edit"); }} />{!workshop.isDefault && <IconAction icon={<Trash2 size={17} />} label={t.delete} variant="destructive" disabled={busy} onClick={() => { setPendingDelete(workshop); setDialog("delete"); }} />}</div></td>
                  <td className="px-4 py-3 font-medium">{workshop.name}{workshop.isDefault && <span className="ms-2 text-xs text-muted-foreground">{t.default}</span>}</td>
                  <td className="px-4 py-3">{workshop.wilaya || "—"}</td><td className="px-4 py-3">{workshop.productCount}</td><td className="px-4 py-3">{workshop.isPublic ? t.publicValue : t.privateValue}</td><td className="px-4 py-3"><div className="flex items-center gap-2"><StatusBadge status={workshop.status} locale={locale} /><IconAction icon={<Power size={17} />} label={workshop.status === "ACTIVE" ? t.active : t.inactive} variant="ghost" disabled={busy} onClick={() => void setStatus(workshop)} /></div></td>
                </tr>)}</tbody>
              </table>
            </div>
            <div className="mt-6 grid gap-3 md:hidden">
              {workshops.map((workshop) => <article className="rounded-lg border border-border bg-card p-4" key={workshop.id}>
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <h3 className="truncate font-medium">{workshop.name}</h3>
                    <p className="mt-1 text-sm text-muted-foreground">{workshop.wilaya || workshop.location || "—"}</p>
                  </div>
                  <StatusBadge status={workshop.status} locale={locale} />
                </div>
                <dl className="mt-4 grid grid-cols-2 gap-x-4 gap-y-3 border-y border-border py-3 text-sm">
                  <div><dt className="text-muted-foreground">{t.products}</dt><dd className="mt-1">{workshop.productCount}</dd></div>
                  <div><dt className="text-muted-foreground">{t.visibility}</dt><dd className="mt-1">{workshop.isPublic ? t.publicValue : t.privateValue}</dd></div>
                </dl>
                <div className="mt-3 flex flex-wrap gap-2">
                  <IconAction icon={<Eye size={17} />} label={t.details} onClick={() => { setFocused(workshop); setDialog("details"); }} />
                  <IconAction icon={<Pencil size={17} />} label={t.update} disabled={busy} onClick={() => { setEditing(workshop.id); setForm(formFrom(workshop)); setDialog("edit"); }} />
                  {!workshop.isDefault && <IconAction icon={<Trash2 size={17} />} label={t.delete} variant="destructive" disabled={busy} onClick={() => { setPendingDelete(workshop); setDialog("delete"); }} />}
                  <IconAction icon={<Power size={17} />} label={workshop.status === "ACTIVE" ? t.active : t.inactive} variant="ghost" disabled={busy} onClick={() => void setStatus(workshop)} />
                </div>
              </article>)}
            </div>
            </>
          )}
        </>
      )}
      {error && <p className="mt-4 text-sm text-destructive" role="alert">{error || t.error}</p>}

      {(dialog === "create" || dialog === "edit") && <Modal label={dialog === "create" ? t.add : t.update} closeLabel={t.close} onClose={closeDialog} panelClassName="w-full p-6 sm:p-10 lg:max-w-3xl"><form className="min-h-full space-y-6" onSubmit={submit}><div className="border-b border-border pb-5"><p className="text-xs uppercase tracking-widest text-primary">AISHA</p><h2 className="mt-2 font-serif text-3xl">{dialog === "create" ? t.add : t.update}</h2></div><WorkshopFields value={form} onChange={setForm} t={t} /><div className="flex flex-wrap gap-3"><Button type="submit" disabled={busy}>{busy ? "…" : t.save}</Button><Button type="button" variant="outline" disabled={busy} onClick={() => closeDialog()}>{t.cancel}</Button></div>{error && <p className="text-sm text-destructive" role="alert">{error}</p>}</form></Modal>}
      {dialog === "details" && focused && <Modal label={t.details} closeLabel={t.close} onClose={closeDialog} panelClassName="w-full p-6 sm:p-10 lg:max-w-2xl"><div className="min-h-full"><p className="text-xs uppercase tracking-widest text-primary">{t.title}</p><h2 className="mt-2 font-serif text-3xl">{focused.name}</h2><dl className="mt-8 divide-y divide-border border-y border-border"><Detail label={t.wilaya} value={focused.wilaya} /><Detail label={t.location} value={focused.location} /><Detail label={t.description} value={focused.description} /><Detail label={t.products} value={String(focused.productCount)} /><Detail label={t.visibility} value={focused.isPublic ? t.publicValue : t.privateValue} /><Detail label="Status" value={focused.status} /></dl><div className="mt-6"><Button type="button" variant="outline" onClick={() => closeDialog()}>{t.close}</Button></div></div></Modal>}
      {dialog === "delete" && pendingDelete && <Modal label={t.delete} closeLabel={t.close} onClose={closeDialog} panelClassName="w-full max-w-lg p-6 sm:p-10"><div><h2 className="font-serif text-3xl">{t.confirmDelete}</h2><p className="mt-3 text-muted-foreground">{pendingDelete.name}</p><p className="mt-2 text-sm text-muted-foreground">{t.deleteHelp}</p><div className="mt-8 flex flex-wrap gap-3"><Button type="button" variant="destructive" disabled={busy} onClick={() => void remove()}>{busy ? t.deleting : t.delete}</Button><Button type="button" variant="outline" disabled={busy} onClick={() => closeDialog()}>{t.cancel}</Button></div></div></Modal>}
    </section>
  );
}

function Detail({ label, value }: { label: string; value?: string }) {
  return <div className="grid gap-2 py-4 sm:grid-cols-[10rem_1fr]"><dt className="text-sm text-muted-foreground">{label}</dt><dd className="whitespace-pre-wrap">{value || "—"}</dd></div>;
}

function WorkshopFields({ value, onChange, t }: { value: WorkshopForm; onChange: (value: WorkshopForm) => void; t: (typeof copy)[Locale] }) {
  const set = (field: keyof WorkshopForm, next: string | boolean) => onChange({ ...value, [field]: next });
  const inputClass = "auth-input mt-1";
  return <div className="grid gap-4 sm:grid-cols-2"><label className="space-y-1 text-sm">{t.name}<input className={inputClass} required minLength={2} value={value.name} onChange={(event: ChangeEvent<HTMLInputElement>) => set("name", event.target.value)} /></label><label className="space-y-1 text-sm">{t.wilaya}<Combobox className="mt-1" options={wilayaOptions} value={value.wilaya} onChange={(next) => set("wilaya", next)} ariaLabel={t.wilaya} placeholder={t.wilaya} /></label><label className="space-y-1 text-sm sm:col-span-2">{t.description}<textarea className={`${inputClass} min-h-24`} value={value.description} onChange={(event: ChangeEvent<HTMLTextAreaElement>) => set("description", event.target.value)} /></label><label className="space-y-1 text-sm sm:col-span-2">{t.location}<input className={inputClass} value={value.location} onChange={(event: ChangeEvent<HTMLInputElement>) => set("location", event.target.value)} /></label><label className="flex items-center gap-2 text-sm sm:col-span-2"><input type="checkbox" checked={value.isPublic} onChange={(event: ChangeEvent<HTMLInputElement>) => set("isPublic", event.target.checked)} />{t.public}</label></div>;
}
