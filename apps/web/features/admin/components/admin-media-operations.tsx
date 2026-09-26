/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { Modal } from "@/core/components/ui/modal";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { formatFullDateTime } from "@/core/lib/format";
import type { Locale } from "@/core/lib/i18n";
import { hasCapability } from "@/features/auth/types";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

type UserMedia = { id: string; userId: string; userDisplayName: string; userEmail: string; mediaKind: string; documentType?: string; originalFilename: string; mediaType: string; sizeBytes: number; url?: string; createdAt: string };
type ProductMedia = { id: string; productId: string; productName: string; productStatus: string; artisanName: string; mediaKind: string; originalFilename: string; mediaType: string; sizeBytes: number; altText: string; visibility: string; url?: string; createdAt: string };
type Tab = "users" | "products";

const labels: Record<Locale, Record<string, string>> = {
  en: { eyebrow: "Media oversight", title: "User and product media", users: "User media", products: "Product media", filename: "File", owner: "Owner", product: "Product", type: "Type", status: "Status", created: "Uploaded", actions: "Actions", open: "Open media", delete: "Delete", cancel: "Cancel", confirmDelete: "Delete this media permanently?", deleting: "Deleting…", deleted: "Media deleted", previous: "Previous", next: "Next", noMedia: "No media found.", loading: "Loading media…", page: "Page", of: "of", error: "Unable to load media", deleteError: "Unable to delete media" },
  fr: { eyebrow: "Contrôle des médias", title: "Médias utilisateurs et produits", users: "Médias utilisateurs", products: "Médias produits", filename: "Fichier", owner: "Propriétaire", product: "Produit", type: "Type", status: "Statut", created: "Téléversé", actions: "Actions", open: "Ouvrir le média", delete: "Supprimer", cancel: "Annuler", confirmDelete: "Supprimer définitivement ce média ?", deleting: "Suppression…", deleted: "Média supprimé", previous: "Précédent", next: "Suivant", noMedia: "Aucun média trouvé.", loading: "Chargement…", page: "Page", of: "sur", error: "Impossible de charger les médias", deleteError: "Impossible de supprimer le média" },
  ar: { eyebrow: "مراقبة الوسائط", title: "وسائط المستخدمين والمنتجات", users: "وسائط المستخدمين", products: "وسائط المنتجات", filename: "الملف", owner: "المالك", product: "المنتج", type: "النوع", status: "الحالة", created: "تاريخ الرفع", actions: "الإجراءات", open: "فتح الوسائط", delete: "حذف", cancel: "إلغاء", confirmDelete: "هل تريد حذف هذه الوسائط نهائياً؟", deleting: "جارٍ الحذف…", deleted: "تم حذف الوسائط", previous: "السابق", next: "التالي", noMedia: "لا توجد وسائط.", loading: "جارٍ التحميل…", page: "الصفحة", of: "من", error: "تعذر تحميل الوسائط", deleteError: "تعذر حذف الوسائط" },
  es: { eyebrow: "Supervisión de medios", title: "Medios de usuarios y productos", users: "Medios de usuarios", products: "Medios de productos", filename: "Archivo", owner: "Propietario", product: "Producto", type: "Tipo", status: "Estado", created: "Subido", actions: "Acciones", open: "Abrir medio", delete: "Eliminar", cancel: "Cancelar", confirmDelete: "¿Eliminar este medio permanentemente?", deleting: "Eliminando…", deleted: "Medio eliminado", previous: "Anterior", next: "Siguiente", noMedia: "No se encontraron medios.", loading: "Cargando…", page: "Página", of: "de", error: "No se pudieron cargar los medios", deleteError: "No se pudo eliminar el medio" },
};

export function AdminMediaOperations({ locale }: { locale: Locale }) {
  const text = labels[locale];
  const auth = useOptionalAuth();
  const canDelete = auth ? hasCapability(auth.user, "admin.media.write") : true;
  const [tab, setTab] = useState<Tab>("users");
  const [users, setUsers] = useState<UserMedia[]>([]);
  const [products, setProducts] = useState<ProductMedia[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [deleting, setDeleting] = useState(false);
  const [pendingDelete, setPendingDelete] = useState<{ id: string; label: string } | null>(null);
  const pageSize = 12;
  const pageCount = Math.max(1, Math.ceil(total / pageSize));

  const load = useCallback(async (nextPage: number, nextTab: Tab) => {
    setLoading(true);
    try {
      const response = await fetch(`/api/admin/media/${nextTab}?page=${nextPage}&pageSize=${pageSize}`, { cache: "no-store" });
      if (!response.ok) throw new Error(text.error);
      const data = (await response.json()) as { items?: UserMedia[] | ProductMedia[]; page?: number; total?: number };
      if (nextTab === "users") setUsers((data.items ?? []) as UserMedia[]);
      else setProducts((data.items ?? []) as ProductMedia[]);
      setPage(data.page ?? nextPage);
      setTotal(data.total ?? 0);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : text.error);
    } finally {
      setLoading(false);
    }
  }, [text.error]);

  useEffect(() => { void load(1, tab); }, [load, tab]);

  const deleteMedia = async () => {
    if (!pendingDelete) return;
    setDeleting(true);
    try {
      const response = await fetch(`/api/admin/media/${tab}/${encodeURIComponent(pendingDelete.id)}`, { method: "DELETE" });
      if (!response.ok) throw new Error(text.deleteError);
      setPendingDelete(null);
      toast.success(text.deleted);
      await load(page, tab);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : text.deleteError);
    } finally {
      setDeleting(false);
    }
  };

  const switchTab = (value: string) => { const next = value as Tab; setTab(next); setPage(1); };
  const currentUsers = tab === "users" ? users : [];
  const currentProducts = tab === "products" ? products : [];

  return <div className="space-y-8">
    <div className="border border-border p-5"><p className="text-xs uppercase tracking-widest text-primary">{text.eyebrow}</p><h2 className="mt-2 font-serif text-3xl">{text.title}</h2><p className="mt-2 text-sm text-muted-foreground">Private files are delivered through short-lived signed URLs.</p></div>
    <Combobox className="max-w-xs" options={[{ value: "users", label: text.users }, { value: "products", label: text.products }]} value={tab} onChange={switchTab} ariaLabel={text.title} />
    {loading ? <p className="border border-border p-5 text-sm text-muted-foreground" role="status">{text.loading}</p> : <>
      <div className="overflow-x-auto border border-border"><table className="w-full min-w-[74rem] text-start text-sm"><thead className="border-b border-border bg-muted/40"><tr>{[text.filename, tab === "users" ? text.owner : text.product, text.type, text.status, text.created, text.actions].map((label, index) => <th key={`${label}-${index}`} className="px-4 py-3 font-medium">{label}</th>)}</tr></thead><tbody className="divide-y divide-border">
        {tab === "users" ? currentUsers.map((item) => <tr key={item.id}><td className="px-4 py-4"><MediaPreview item={item} openLabel={text.open} /><div className="mt-2">{item.originalFilename || item.mediaKind}</div><span className="text-xs text-muted-foreground">{item.mediaType} · {formatBytes(item.sizeBytes)}</span></td><td className="px-4 py-4">{item.userDisplayName || item.userEmail}<br /><span className="text-xs text-muted-foreground">{item.userEmail}</span></td><td className="px-4 py-4">{item.documentType || item.mediaKind}</td><td className="px-4 py-4"><StatusBadge status={item.mediaKind} /></td><td className="px-4 py-4">{formatFullDateTime(item.createdAt, locale)}</td><td className="px-4 py-4"><MediaActions item={item} openLabel={text.open} deleteLabel={text.delete} canDelete={canDelete} onDelete={() => setPendingDelete({ id: item.id, label: item.originalFilename || item.mediaKind })} /></td></tr>) : currentProducts.map((item) => <tr key={item.id}><td className="px-4 py-4"><MediaPreview item={item} openLabel={text.open} /><div className="mt-2">{item.originalFilename || item.mediaKind}</div><span className="text-xs text-muted-foreground">{item.mediaType} · {formatBytes(item.sizeBytes)}</span></td><td className="px-4 py-4">{item.productName}<br /><span className="text-xs text-muted-foreground">{item.artisanName}</span></td><td className="px-4 py-4">{item.mediaKind}</td><td className="px-4 py-4"><StatusBadge status={item.productStatus} /><div className="mt-2 text-xs">{item.visibility}</div></td><td className="px-4 py-4">{formatFullDateTime(item.createdAt, locale)}</td><td className="px-4 py-4"><MediaActions item={item} openLabel={text.open} deleteLabel={text.delete} canDelete={canDelete} onDelete={() => setPendingDelete({ id: item.id, label: item.originalFilename || item.mediaKind })} /></td></tr>)}
        {((tab === "users" && currentUsers.length === 0) || (tab === "products" && currentProducts.length === 0)) && <tr><td className="px-4 py-8 text-muted-foreground" colSpan={6}>{text.noMedia}</td></tr>}
      </tbody></table></div>
      <div className="flex items-center justify-between text-sm"><span>{text.page} {page} {text.of} {pageCount}</span><div className="flex gap-2"><Button type="button" variant="outline" disabled={page <= 1} onClick={() => void load(page - 1, tab)}>{text.previous}</Button><Button type="button" variant="outline" disabled={page >= pageCount} onClick={() => void load(page + 1, tab)}>{text.next}</Button></div></div>
    </>}
    {pendingDelete && <Modal label={text.delete} closeLabel={text.cancel} onClose={() => setPendingDelete(null)} panelClassName="mx-auto my-8 h-fit max-h-[calc(100svh-2rem)] w-[min(100%-2rem,32rem)] overflow-y-auto p-6"><h3 className="font-serif text-2xl">{text.delete}</h3><p className="mt-4 text-sm text-muted-foreground">{text.confirmDelete}</p><p className="mt-2 truncate text-sm font-medium">{pendingDelete.label}</p><div className="mt-6 flex justify-end gap-2"><Button type="button" variant="outline" onClick={() => setPendingDelete(null)}>{text.cancel}</Button><Button type="button" variant="destructive" disabled={deleting} onClick={() => void deleteMedia()}>{deleting ? text.deleting : text.delete}</Button></div></Modal>}
  </div>;
}

function MediaActions({ item, openLabel, deleteLabel, canDelete, onDelete }: { item: { url?: string }; openLabel: string; deleteLabel: string; canDelete: boolean; onDelete: () => void }) { return <div className="flex flex-wrap gap-2">{item.url ? <a className="underline" href={item.url} target="_blank" rel="noreferrer">{openLabel}</a> : null}{canDelete && <Button type="button" variant="destructive" onClick={onDelete}>{deleteLabel}</Button>}</div>; }
function formatBytes(value: number) { if (value < 1024) return `${value} B`; if (value < 1024 * 1024) return `${Math.round(value / 1024)} KB`; return `${(value / (1024 * 1024)).toFixed(1)} MB`; }
function MediaPreview({ item, openLabel }: { item: { mediaType: string; url?: string }; openLabel: string }) { if (!item.url) return <span className="inline-flex h-20 w-28 items-center justify-center border border-border bg-muted text-xs text-muted-foreground">{item.mediaType}</span>; if (item.mediaType.startsWith("video/")) return <video className="h-20 w-28 border border-border bg-muted object-cover" controls preload="metadata" src={item.url} aria-label={openLabel} />; if (!item.mediaType.startsWith("image/")) return <span className="inline-flex h-20 w-28 items-center justify-center border border-border bg-muted text-xs text-muted-foreground">{item.mediaType}</span>; return <a href={item.url} target="_blank" rel="noreferrer" aria-label={openLabel} className="block h-20 w-28 overflow-hidden border border-border bg-muted"><span className="block h-full w-full bg-cover bg-center" style={{ backgroundImage: `url("${item.url}")` }} /></a>; }
