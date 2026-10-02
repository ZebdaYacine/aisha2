/* eslint-disable @typescript-eslint/no-explicit-any */
"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { FormErrorSummary } from "@/core/components/forms/form-error-summary";
import { Button } from "@/core/components/ui/button";
import { FileDropzone } from "@/core/components/ui/file-dropzone";
import { UploadProgress } from "@/core/components/ui/upload-progress";

const copy = {
  en: { title: "Artisan application", submit: "Submit application", name: "Public display name", workshop: "Workshop name", wilaya: "Wilaya", location: "Workshop address", email: "Contact email", phone: "Contact phone", categories: "Craft categories", bio: "Biography", private: "Keep contact private", public: "Make contact public", files: "Application files", filesHelp: "Add at least one supporting document and one profile image or video before submitting.", document: "Supporting document", documentType: "Document type", media: "Profile media", mediaType: "Media type", identity: "Identity document", registration: "Business registration", proof: "Proof of address", image: "Image", video: "Video", chooseFile: "Choose file", existing: "Already uploaded", filesRequired: "A workshop, document, and profile media are required before submission.", uploadFailed: "The files could not be uploaded. Your draft was saved; please try again.", invalid: "Please correct the highlighted fields.", categoriesUnavailable: "Craft categories are temporarily unavailable." },
  fr: { title: "Candidature artisan", submit: "Envoyer la candidature", name: "Nom public", workshop: "Nom de l’atelier", wilaya: "Wilaya", location: "Adresse de l’atelier", email: "E-mail de contact", phone: "Téléphone", categories: "Catégories artisanales", bio: "Biographie", private: "Garder le contact privé", public: "Rendre le contact public", files: "Fichiers de candidature", filesHelp: "Ajoutez au moins un document justificatif et une image ou vidéo de profil avant l’envoi.", document: "Document justificatif", documentType: "Type de document", media: "Média du profil", mediaType: "Type de média", identity: "Pièce d’identité", registration: "Registre de commerce", proof: "Justificatif de domicile", image: "Image", video: "Vidéo", chooseFile: "Choisir un fichier", existing: "Déjà importé", filesRequired: "Un atelier, un document et un média de profil sont requis avant l’envoi.", uploadFailed: "Les fichiers n’ont pas pu être importés. Votre brouillon est conservé ; réessayez.", invalid: "Corrigez les champs indiqués.", categoriesUnavailable: "Les catégories artisanales sont temporairement indisponibles." },
  ar: { title: "طلب الانضمام كحرفي", submit: "إرسال الطلب", name: "الاسم العلني", workshop: "اسم الورشة", wilaya: "الولاية", location: "عنوان الورشة", email: "البريد الإلكتروني", phone: "الهاتف", categories: "فئات الحرفة", bio: "نبذة", private: "إبقاء بيانات الاتصال خاصة", public: "إظهار بيانات الاتصال", files: "ملفات الطلب", filesHelp: "أضف وثيقة داعمة واحدة ووسائط للملف الشخصي قبل الإرسال.", document: "وثيقة داعمة", documentType: "نوع الوثيقة", media: "وسائط الملف", mediaType: "نوع الوسائط", identity: "وثيقة الهوية", registration: "السجل التجاري", proof: "إثبات العنوان", image: "صورة", video: "فيديو", chooseFile: "اختر ملفاً", existing: "تم رفعه مسبقاً", filesRequired: "يجب إدخال ورشة ووثيقة ووسائط للملف الشخصي قبل الإرسال.", uploadFailed: "تعذر رفع الملفات. تم حفظ المسودة؛ حاول مرة أخرى.", invalid: "يرجى تصحيح الحقول المحددة.", categoriesUnavailable: "فئات الحرفة غير متاحة مؤقتًا." },
  es: { title: "Solicitud de artesano", submit: "Enviar solicitud", name: "Nombre público", workshop: "Nombre del taller", wilaya: "Provincia", location: "Dirección del taller", email: "Correo de contacto", phone: "Teléfono", categories: "Categorías artesanales", bio: "Biografía", private: "Mantener el contacto privado", public: "Mostrar el contacto", files: "Archivos de solicitud", filesHelp: "Añade al menos un documento y una imagen o vídeo de perfil antes de enviar.", document: "Documento justificativo", documentType: "Tipo de documento", media: "Medio del perfil", mediaType: "Tipo de medio", identity: "Documento de identidad", registration: "Registro comercial", proof: "Comprobante de domicilio", image: "Imagen", video: "Vídeo", chooseFile: "Elegir archivo", existing: "Ya subido", filesRequired: "Se requieren un taller, un documento y un medio de perfil antes de enviar.", uploadFailed: "No se pudieron subir los archivos. El borrador se guardó; inténtalo de nuevo.", invalid: "Corrige los campos indicados.", categoriesUnavailable: "Las categorías artesanales no están disponibles temporalmente." },
} as const;

type Locale = keyof typeof copy;
const uploadingLabels: Record<Locale, string> = {
  en: "Uploading files…",
  fr: "Importation des fichiers…",
  ar: "جارٍ رفع الملفات…",
  es: "Subiendo archivos…",
};
type Category = { id: string; name: string };
type Values = {
  publicDisplayName: string;
  workshopName: string;
  wilaya: string;
  location: string;
  contactEmail: string;
  contactPhone: string;
  contactVisibility: "PRIVATE" | "PUBLIC";
  categoryIds: string[];
  biography: string;
};

type CatalogueResponse = {
  categories?: Array<{ id: string; name?: string | Record<string, string> }>;
};

function categoryName(name: string | Record<string, string> | undefined, locale: Locale, fallback: string) {
  if (typeof name === "string") return name;
  return name?.[locale] ?? name?.en ?? fallback;
}

export function ApplicationForm({ locale, profile }: { locale: Locale; profile?: any }) {
  const t = copy[locale];
  const [categories, setCategories] = useState<Category[]>([]);
  const [summary, setSummary] = useState<string>();
  const [documentFile, setDocumentFile] = useState<File>();
  const [mediaFile, setMediaFile] = useState<File>();
  const [documentType, setDocumentType] = useState("IDENTITY");
  const [mediaKind, setMediaKind] = useState("IMAGE");
  const [existingDocuments, setExistingDocuments] = useState(0);
  const [existingMedia, setExistingMedia] = useState(0);
  const schema = z.object({
    publicDisplayName: z.string().trim().min(2, t.invalid),
    workshopName: z.string().trim().min(2, t.invalid),
    wilaya: z.string().trim().min(1, t.invalid),
    location: z.string(),
    contactEmail: z.union([z.literal(""), z.email(t.invalid)]),
    contactPhone: z.string(),
    contactVisibility: z.enum(["PRIVATE", "PUBLIC"]),
    categoryIds: z.array(z.string()).min(1, t.invalid),
    biography: z.string().max(5000),
  });
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: {
      publicDisplayName: profile?.publicDisplayName ?? "",
      workshopName: profile?.workshopName ?? "",
      wilaya: profile?.wilaya ?? "",
      location: profile?.location ?? "",
      contactEmail: profile?.contactEmail ?? "",
      contactPhone: profile?.contactPhone ?? "",
      contactVisibility: profile?.contactVisibility ?? "PRIVATE",
      categoryIds: profile?.categoryIds ?? [],
      biography: profile?.translations?.find((x: any) => x.locale === locale)?.biography ?? "",
    },
  });

  useEffect(() => {
    let active = true;
    void fetch(`/api/catalogue?locale=${locale}`)
      .then(async (response) => {
        if (!response.ok) throw new Error("categories request failed");
        return (await response.json()) as CatalogueResponse;
      })
      .then((data) => {
        if (!active) return;
        setCategories(
          (data.categories ?? []).map((category) => ({
            id: category.id,
            name: categoryName(category.name, locale, category.id),
          })),
        );
      })
      .catch(() => {
        if (active) setSummary(t.categoriesUnavailable);
      });
    return () => {
      active = false;
    };
  }, [locale, t.categoriesUnavailable]);

  useEffect(() => {
    if (!profile || profile.status === "APPROVED") return;
    void Promise.all([fetch("/api/artisan/documents"), fetch("/api/artisan/media")]).then(async ([documentsResponse, mediaResponse]) => {
      if (documentsResponse.ok) setExistingDocuments((await documentsResponse.json()).length);
      if (mediaResponse.ok) setExistingMedia((await mediaResponse.json()).length);
    }).catch(() => undefined);
  }, [profile]);

  const submit = async (values: Values) => {
    setSummary(undefined);
    const body = { ...values, translations: [{ locale, biography: values.biography }] };
    delete (body as Partial<typeof body>).biography;
    const approved = profile?.status === "APPROVED";
    if (
      !approved &&
      ((!documentFile && existingDocuments === 0) || (!mediaFile && existingMedia === 0))
    ) {
      setSummary(t.filesRequired);
      toast.error(t.filesRequired);
      return;
    }
    let response: Response;
    let uploadError = false;
    if (approved) {
      response = await fetch("/api/artisan/profile", { method: "PATCH", headers: { "content-type": "application/json" }, body: JSON.stringify(body) });
    } else {
      response = await fetch("/api/artisan/application/draft", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(body) });
      if (response.ok && documentFile) {
        const documentForm = new FormData();
        documentForm.set("documentType", documentType);
        documentForm.set("file", documentFile);
        response = await fetch("/api/artisan/documents", { method: "POST", body: documentForm });
        uploadError = !response.ok;
        if (response.ok) setExistingDocuments((count) => count + 1);
      }
      if (response.ok && mediaFile) {
        const mediaForm = new FormData();
        mediaForm.set("mediaKind", mediaKind);
        mediaForm.set("file", mediaFile);
        response = await fetch("/api/artisan/media", { method: "POST", body: mediaForm });
        uploadError = !response.ok;
        if (response.ok) setExistingMedia((count) => count + 1);
      }
      if (response.ok) response = await fetch("/api/artisan/application/submit", { method: "POST" });
    }
    if (response.ok) {
      toast.success(t.submit);
      location.reload();
      return;
    }
    const data = await response.json().catch(() => ({}));
    Object.entries(data.error?.fields ?? {}).forEach(([field, message]) => {
      setError(field as keyof Values, { message: String(message) });
    });
    const message = uploadError ? t.uploadFailed : data.error?.message ?? t.invalid;
    setSummary(message);
    toast.error(message);
  };

  return (
    <form onSubmit={handleSubmit(submit, () => setSummary(t.invalid))} className="space-y-5" noValidate>
      <FormErrorSummary id="artisan-errors" message={summary} />
      {(["publicDisplayName", "workshopName", "wilaya", "location", "contactEmail", "contactPhone"] as const).map((field) => {
        const labels = { publicDisplayName: t.name, workshopName: t.workshop, wilaya: t.wilaya, location: t.location, contactEmail: t.email, contactPhone: t.phone };
        return (
          <label className="block" key={field}>
            <span className="mb-2 block text-sm">{labels[field]}</span>
            <input className="auth-input" {...register(field)} aria-invalid={!!errors[field]} />
            {errors[field] && <span className="border-b border-destructive text-xs text-destructive">{errors[field]?.message}</span>}
          </label>
        );
      })}
      <fieldset>
        <legend className="mb-2 text-sm">{t.categories}</legend>
        <div className="grid gap-2 sm:grid-cols-2">
          {categories.map((category) => (
            <label key={category.id} className="flex gap-2">
              <input type="checkbox" value={category.id} {...register("categoryIds")} />
              {category.name}
            </label>
          ))}
        </div>
        {errors.categoryIds && <span className="text-xs text-destructive">{errors.categoryIds.message}</span>}
      </fieldset>
      <label className="block">
        <span className="mb-2 block text-sm">{t.bio}</span>
        <textarea className="auth-input min-h-32" {...register("biography")} />
      </label>
      <label className="block">
        <span className="mb-2 block text-sm">{t.email}</span>
        <select className="auth-input" {...register("contactVisibility")}>
          <option value="PRIVATE">{t.private}</option>
          <option value="PUBLIC">{t.public}</option>
        </select>
      </label>
      {!profile?.status || profile.status !== "APPROVED" ? <fieldset className="space-y-4 rounded-xl border border-border p-4 sm:p-5">
        <legend className="px-1 font-medium">{t.files}</legend>
        <p className="text-sm text-muted-foreground">{t.filesHelp}</p>
        <label className="block"><span className="mb-2 block text-sm">{t.document}</span><select className="auth-input mb-3" value={documentType} onChange={(event) => setDocumentType(event.target.value)}><option value="IDENTITY">{t.identity}</option><option value="BUSINESS_REGISTRATION">{t.registration}</option><option value="PROOF_OF_ADDRESS">{t.proof}</option></select><FileDropzone accept="application/pdf,image/jpeg,image/png" label={t.chooseFile} file={documentFile} onFileChange={setDocumentFile} />{existingDocuments > 0 && <span className="mt-2 block text-xs text-muted-foreground">{t.existing}: {existingDocuments}</span>}</label>
        <label className="block"><span className="mb-2 block text-sm">{t.media}</span><select className="auth-input mb-3" value={mediaKind} onChange={(event) => setMediaKind(event.target.value)}><option value="IMAGE">{t.image}</option><option value="VIDEO">{t.video}</option></select><FileDropzone accept="image/jpeg,image/png,image/webp,video/mp4" label={t.chooseFile} file={mediaFile} onFileChange={setMediaFile} />{existingMedia > 0 && <span className="mt-2 block text-xs text-muted-foreground">{t.existing}: {existingMedia}</span>}</label>
      </fieldset> : null}
      {isSubmitting && <UploadProgress label={uploadingLabels[locale]} />}
      <Button disabled={isSubmitting} type="submit">{t.submit}</Button>
    </form>
  );
}
