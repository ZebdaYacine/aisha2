"use client";

import { useEffect, useRef, useState } from "react";
import { Eye, Pencil, Plus, Trash2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { IconAction, IconActionLink } from "@/core/components/ui/icon-action";
import { Modal } from "@/core/components/ui/modal";
import { UploadProgress } from "@/core/components/ui/upload-progress";
import { formatFullDate } from "@/core/lib/format";

type Locale = "en" | "fr" | "ar" | "es";
type DocumentItem = { id: string; documentType: string; originalFilename: string; mediaType: string; createdAt?: string; url?: string };
type MediaItem = { id: string; mediaKind: string; originalFilename: string; mediaType: string; createdAt?: string; url?: string };
type MediaDialog = "create" | "edit" | "delete" | null;

const labels = {
  en: { title: "Private files", document: "Application documents", documentType: "Document type", media: "Profile media", mediaName: "Media name", mediaKind: "Media type", uploadDate: "Uploaded", actions: "Actions", upload: "Add media", addDocument: "Add document", uploading: "Uploading…", updating: "Updating…", deleting: "Deleting…", view: "View", update: "Update", delete: "Delete", close: "Close", save: "Save", chooseFile: "Choose a file from your computer", createMedia: "Add private media", editMedia: "Update private media", deleteMedia: "Delete private media", confirmDelete: "Delete this private media?", deleteHelp: "This removes the private file from your profile.", createDocument: "Add application document", identity: "Identity document", registration: "Business registration", proof: "Proof of address", selectType: "Select a type", selectMedia: "Select a media type", image: "Image", video: "Video", noFiles: "No private files yet.", unavailable: "Unable to load private files.", fileRequired: "Choose a file before saving." },
  fr: { title: "Fichiers privés", document: "Documents de candidature", documentType: "Type de document", media: "Médias du profil", mediaName: "Nom du média", mediaKind: "Type de média", uploadDate: "Importé le", actions: "Actions", upload: "Ajouter un média", addDocument: "Ajouter un document", uploading: "Importation…", updating: "Mise à jour…", deleting: "Suppression…", view: "Voir", update: "Modifier", delete: "Supprimer", close: "Fermer", save: "Enregistrer", chooseFile: "Choisissez un fichier sur votre ordinateur", createMedia: "Ajouter un média privé", editMedia: "Modifier le média privé", deleteMedia: "Supprimer le média privé", confirmDelete: "Supprimer ce média privé ?", deleteHelp: "Le fichier privé sera retiré de votre profil.", createDocument: "Ajouter un document de candidature", identity: "Pièce d’identité", registration: "Registre de commerce", proof: "Justificatif de domicile", selectType: "Sélectionner un type", selectMedia: "Sélectionner un type de média", image: "Image", video: "Vidéo", noFiles: "Aucun fichier privé.", unavailable: "Impossible de charger les fichiers privés.", fileRequired: "Choisissez un fichier avant d’enregistrer." },
  ar: { title: "الملفات الخاصة", document: "وثائق الطلب", documentType: "نوع الوثيقة", media: "وسائط الملف", mediaName: "اسم الوسائط", mediaKind: "نوع الوسائط", uploadDate: "تاريخ الرفع", actions: "الإجراءات", upload: "إضافة وسائط", addDocument: "إضافة وثيقة", uploading: "جارٍ الرفع…", updating: "جارٍ التحديث…", deleting: "جارٍ الحذف…", view: "عرض", update: "تحديث", delete: "حذف", close: "إغلاق", save: "حفظ", chooseFile: "اختر ملفاً من جهازك", createMedia: "إضافة وسائط خاصة", editMedia: "تحديث الوسائط الخاصة", deleteMedia: "حذف الوسائط الخاصة", confirmDelete: "حذف هذه الوسائط الخاصة؟", deleteHelp: "سيتم حذف الملف الخاص من ملفك.", createDocument: "إضافة وثيقة الطلب", identity: "وثيقة الهوية", registration: "السجل التجاري", proof: "إثبات العنوان", selectType: "اختر النوع", selectMedia: "اختر نوع الوسائط", image: "صورة", video: "فيديو", noFiles: "لا توجد ملفات خاصة بعد.", unavailable: "تعذر تحميل الملفات الخاصة.", fileRequired: "اختر ملفاً قبل الحفظ." },
  es: { title: "Archivos privados", document: "Documentos de solicitud", documentType: "Tipo de documento", media: "Medios del perfil", mediaName: "Nombre del medio", mediaKind: "Tipo de medio", uploadDate: "Subido", actions: "Acciones", upload: "Añadir medio", addDocument: "Añadir documento", uploading: "Subiendo…", updating: "Actualizando…", deleting: "Eliminando…", view: "Ver", update: "Actualizar", delete: "Eliminar", close: "Cerrar", save: "Guardar", chooseFile: "Elige un archivo de tu ordenador", createMedia: "Añadir medio privado", editMedia: "Actualizar medio privado", deleteMedia: "Eliminar medio privado", confirmDelete: "¿Eliminar este medio privado?", deleteHelp: "El archivo privado se quitará de tu perfil.", createDocument: "Añadir documento de solicitud", identity: "Documento de identidad", registration: "Registro comercial", proof: "Comprobante de domicilio", selectType: "Selecciona un tipo", selectMedia: "Selecciona un tipo de medio", image: "Imagen", video: "Vídeo", noFiles: "Aún no hay archivos privados.", unavailable: "No se pudieron cargar los archivos privados.", fileRequired: "Elige un archivo antes de guardar." },
} as const;

export function ArtisanMediaPanel({ locale }: { locale: Locale }) {
  const text = labels[locale];
  const [documents, setDocuments] = useState<DocumentItem[]>([]);
  const [media, setMedia] = useState<MediaItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [mediaDialog, setMediaDialog] = useState<MediaDialog>(null);
  const [documentDialog, setDocumentDialog] = useState<"create" | null>(null);
  const [selectedMedia, setSelectedMedia] = useState<MediaItem>();
  const [mediaName, setMediaName] = useState("");
  const [mediaKind, setMediaKind] = useState("IMAGE");
  const [documentType, setDocumentType] = useState("IDENTITY");
  const mediaInput = useRef<HTMLInputElement>(null);
  const documentInput = useRef<HTMLInputElement>(null);

  useEffect(() => {
    void Promise.all([fetch("/api/artisan/documents"), fetch("/api/artisan/media")])
      .then(async ([documentsResponse, mediaResponse]) => {
        if (!documentsResponse.ok || !mediaResponse.ok) throw new Error();
        setDocuments(await documentsResponse.json());
        setMedia(await mediaResponse.json());
      })
      .catch(() => toast.error(text.unavailable))
      .finally(() => setLoading(false));
  }, [text.unavailable]);

  const openMedia = (item?: MediaItem) => {
    setSelectedMedia(item);
    setMediaName(item?.originalFilename ?? "");
    setMediaKind(item?.mediaKind ?? "IMAGE");
    setMediaDialog(item ? "edit" : "create");
  };

  const submitMedia = async () => {
    const file = mediaInput.current?.files?.[0];
    if (!file) { toast.error(text.fileRequired); return; }
    setBusy(true);
    try {
      const form = new FormData();
      form.set("file", file);
      form.set("mediaName", mediaName.trim() || file.name);
      form.set("mediaKind", mediaKind);
      const response = await fetch(selectedMedia ? `/api/artisan/media/${selectedMedia.id}` : "/api/artisan/media", { method: selectedMedia ? "PATCH" : "POST", body: form });
      const body = await response.json().catch(() => undefined);
      if (!response.ok) throw new Error(body?.error?.message ?? text.unavailable);
      if (selectedMedia) setMedia((current) => current.map((item) => item.id === selectedMedia.id ? body as MediaItem : item));
      else setMedia((current) => [...current, body as MediaItem]);
      setMediaDialog(null);
      if (mediaInput.current) mediaInput.current.value = "";
      toast.success(selectedMedia ? text.update : text.upload);
    } catch (error) { toast.error((error as Error).message); }
    finally { setBusy(false); }
  };

  const deleteMedia = async () => {
    if (!selectedMedia) return;
    setBusy(true);
    try {
      const response = await fetch(`/api/artisan/media/${selectedMedia.id}`, { method: "DELETE" });
      if (!response.ok) { const body = await response.json().catch(() => undefined); throw new Error(body?.error?.message ?? text.unavailable); }
      setMedia((current) => current.filter((item) => item.id !== selectedMedia.id));
      setMediaDialog(null);
      toast.success(text.delete);
    } catch (error) { toast.error((error as Error).message); }
    finally { setBusy(false); }
  };

  const submitDocument = async () => {
    const file = documentInput.current?.files?.[0];
    if (!file) { toast.error(text.fileRequired); return; }
    setBusy(true);
    try {
      const form = new FormData(); form.set("file", file); form.set("documentType", documentType);
      const response = await fetch("/api/artisan/documents", { method: "POST", body: form });
      const body = await response.json().catch(() => undefined);
      if (!response.ok) throw new Error(body?.error?.message ?? text.unavailable);
      setDocuments((current) => [...current, body as DocumentItem]);
      setDocumentDialog(null);
      if (documentInput.current) documentInput.current.value = "";
      toast.success(text.addDocument);
    } catch (error) { toast.error((error as Error).message); }
    finally { setBusy(false); }
  };

  if (loading) return null;
  return (
    <section className="mt-10 min-w-0 border-t border-border pt-8">
      <div className="flex flex-wrap items-end justify-between gap-4"><div className="min-w-0"><p className="text-xs uppercase tracking-widest text-primary">AISHA</p><h2 className="mt-2 font-serif text-2xl">{text.title}</h2></div><div className="flex gap-2"><IconAction icon={<Plus size={18} />} label={text.addDocument} variant="outline" onClick={() => setDocumentDialog("create")} /><IconAction icon={<Plus size={18} />} label={text.upload} onClick={() => openMedia()} /></div></div>
      <MediaTable title={text.document} columns={[text.mediaName, text.documentType, text.uploadDate, text.actions]} empty={documents.length === 0 ? text.noFiles : undefined}>{documents.map((item) => <MediaRow key={item.id} name={item.originalFilename || item.documentType} type={item.documentType} date={item.createdAt} locale={locale} labels={{ name: text.mediaName, type: text.documentType, date: text.uploadDate, actions: text.actions }} actions={<ViewButton item={item} label={text.view} />} />)}</MediaTable>
      <MediaTable title={text.media} columns={[text.mediaName, text.mediaKind, text.uploadDate, text.actions]} empty={media.length === 0 ? text.noFiles : undefined}>{media.map((item) => <MediaRow key={item.id} name={item.originalFilename || item.mediaKind} type={item.mediaKind} date={item.createdAt} locale={locale} labels={{ name: text.mediaName, type: text.mediaKind, date: text.uploadDate, actions: text.actions }} actions={<div className="flex flex-wrap gap-2"><ViewButton item={item} label={text.view} /><IconAction icon={<Pencil size={17} />} label={text.update} variant="outline" onClick={() => openMedia(item)} /><IconAction icon={<Trash2 size={17} />} label={text.delete} variant="destructive" onClick={() => { setSelectedMedia(item); setMediaDialog("delete"); }} /></div>} />)}</MediaTable>

      {documentDialog === "create" && <Modal label={text.createDocument} closeLabel={text.close} onClose={() => !busy && setDocumentDialog(null)} panelClassName="w-full p-6 sm:p-10 lg:max-w-2xl"><MediaFormTitle title={text.createDocument} /><div className="mt-8 space-y-5"><label className="block"><span className="mb-1 block text-sm">{text.documentType}</span><Combobox options={[{ value: "IDENTITY", label: text.identity }, { value: "BUSINESS_REGISTRATION", label: text.registration }, { value: "PROOF_OF_ADDRESS", label: text.proof }]} value={documentType} onChange={setDocumentType} placeholder={text.selectType} ariaLabel={text.documentType} disabled={busy} /></label><FilePicker inputRef={documentInput} accept="application/pdf,image/jpeg,image/png" label={text.chooseFile} disabled={busy} />{busy && <UploadProgress label={text.uploading} />}<div className="flex flex-wrap gap-3"><Button type="button" disabled={busy} onClick={() => void submitDocument()}>{busy ? text.uploading : text.save}</Button><Button type="button" variant="outline" disabled={busy} onClick={() => setDocumentDialog(null)}>{text.close}</Button></div></div></Modal>}
      {mediaDialog === "create" || mediaDialog === "edit" ? <Modal label={mediaDialog === "create" ? text.createMedia : text.editMedia} closeLabel={text.close} onClose={() => !busy && setMediaDialog(null)} panelClassName="w-full p-6 sm:p-10 lg:max-w-2xl"><MediaFormTitle title={mediaDialog === "create" ? text.createMedia : text.editMedia} /><div className="mt-8 space-y-5"><label className="block"><span className="mb-1 block text-sm">{text.mediaName}</span><input className="auth-input" value={mediaName} onChange={(event) => setMediaName(event.target.value)} disabled={busy} /></label><label className="block"><span className="mb-1 block text-sm">{text.mediaKind}</span><Combobox options={[{ value: "IMAGE", label: text.image }, { value: "VIDEO", label: text.video }]} value={mediaKind} onChange={setMediaKind} placeholder={text.selectMedia} ariaLabel={text.mediaKind} disabled={busy} /></label><FilePicker inputRef={mediaInput} accept="image/jpeg,image/png,image/webp,video/mp4" label={text.chooseFile} disabled={busy} />{busy && <UploadProgress label={mediaDialog === "edit" ? text.updating : text.uploading} />}<div className="flex flex-wrap gap-3"><Button type="button" disabled={busy} onClick={() => void submitMedia()}>{busy ? (mediaDialog === "edit" ? text.updating : text.uploading) : text.save}</Button><Button type="button" variant="outline" disabled={busy} onClick={() => setMediaDialog(null)}>{text.close}</Button></div></div></Modal> : null}
      {mediaDialog === "delete" && selectedMedia ? <Modal label={text.deleteMedia} closeLabel={text.close} onClose={() => !busy && setMediaDialog(null)} panelClassName="w-full max-w-lg p-6 sm:p-10"><MediaFormTitle title={text.confirmDelete} /><p className="mt-3 text-muted-foreground">{selectedMedia.originalFilename}</p><p className="mt-2 text-sm text-muted-foreground">{text.deleteHelp}</p><div className="mt-8 flex flex-wrap gap-3"><Button type="button" variant="destructive" disabled={busy} onClick={() => void deleteMedia()}>{busy ? text.deleting : text.delete}</Button><Button type="button" variant="outline" disabled={busy} onClick={() => setMediaDialog(null)}>{text.close}</Button></div></Modal> : null}
    </section>
  );
}

function MediaTable({ title, columns, empty, children }: { title: string; columns: string[]; empty?: string; children: React.ReactNode }) { return <div className="mt-8 min-w-0"><h3 className="font-medium">{title}</h3><div className="artisan-media-table-wrap mt-3 overflow-x-auto rounded-lg border border-border"><table className="artisan-media-table w-full text-sm sm:min-w-[720px]"><thead className="bg-muted/40"><tr>{columns.map((column) => <th key={column} className="px-4 py-3 text-start font-medium">{column}</th>)}</tr></thead><tbody>{empty ? <tr><td className="px-4 py-6 text-muted-foreground" colSpan={columns.length}>{empty}</td></tr> : children}</tbody></table></div></div>; }
function MediaRow({ name, type, date, locale, labels, actions }: { name: string; type: string; date?: string; locale: Locale; labels: { name: string; type: string; date: string; actions: string }; actions: React.ReactNode }) { return <tr className="border-t border-border align-middle"><td data-label={labels.name} className="max-w-[20rem] px-4 py-4"><span className="block truncate font-medium" title={name}>{name}</span></td><td data-label={labels.type} className="px-4 py-4 text-muted-foreground">{type}</td><td data-label={labels.date} className="px-4 py-4 text-muted-foreground"><time dateTime={date}>{date ? formatFullDate(date, locale) : "—"}</time></td><td data-label={labels.actions} className="px-4 py-4">{actions}</td></tr>; }
function ViewButton({ item, label }: { item: { url?: string }; label: string }) { return item.url ? <IconActionLink href={item.url} label={label} icon={<Eye size={17} />} /> : <IconAction icon={<Eye size={17} />} label={label} variant="outline" disabled />; }
function FilePicker({ inputRef, accept, label, disabled }: { inputRef: React.RefObject<HTMLInputElement | null>; accept: string; label: string; disabled: boolean }) { return <div className="rounded-lg border border-dashed border-border p-4"><input ref={inputRef} className="block w-full text-sm" type="file" accept={accept} aria-label={label} disabled={disabled} /><p className="mt-2 text-xs text-muted-foreground">{label}</p></div>; }
function MediaFormTitle({ title }: { title: string }) { return <><p className="text-xs uppercase tracking-widest text-primary">AISHA</p><h2 className="mt-2 font-serif text-3xl">{title}</h2></>; }
