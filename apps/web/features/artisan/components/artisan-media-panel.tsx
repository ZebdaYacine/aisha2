"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";

type Locale = "en" | "fr" | "ar" | "es";
type DocumentItem = { id: string; documentType: string; originalFilename: string; url?: string };
type MediaItem = { id: string; mediaKind: string; originalFilename: string; url?: string };
const labels = {
  en: { title: "Private files", document: "Application document", documentType: "Document type", upload: "Upload", profile: "Profile media", unavailable: "Unable to load private files." },
  fr: { title: "Fichiers privés", document: "Document de candidature", documentType: "Type de document", upload: "Importer", profile: "Médias du profil", unavailable: "Impossible de charger les fichiers privés." },
  ar: { title: "الملفات الخاصة", document: "وثيقة الطلب", documentType: "نوع الوثيقة", upload: "رفع", profile: "وسائط الملف", unavailable: "تعذر تحميل الملفات الخاصة." },
  es: { title: "Archivos privados", document: "Documento de solicitud", documentType: "Tipo de documento", upload: "Subir", profile: "Medios del perfil", unavailable: "No se pudieron cargar los archivos privados." },
} as const;

export function ArtisanMediaPanel({ locale }: { locale: Locale }) {
  const text = labels[locale];
  const [documents, setDocuments] = useState<DocumentItem[]>([]);
  const [media, setMedia] = useState<MediaItem[]>([]);
  const [documentType, setDocumentType] = useState("IDENTITY");
  const [loading, setLoading] = useState(true);
  useEffect(() => { void Promise.all([fetch("/api/artisan/documents"), fetch("/api/artisan/media")]).then(async ([documentsResponse, mediaResponse]) => { if (!documentsResponse.ok || !mediaResponse.ok) throw new Error(); setDocuments(await documentsResponse.json()); setMedia(await mediaResponse.json()); }).catch(() => toast.error(text.unavailable)).finally(() => setLoading(false)); }, [text.unavailable]);
  const upload = async (endpoint: string, file: File, extra?: Record<string, string>) => { const form = new FormData(); form.set("file", file); Object.entries(extra ?? {}).forEach(([key, value]) => form.set(key, value)); const response = await fetch(endpoint, { method: "POST", body: form }); if (!response.ok) { const body = await response.json().catch(() => undefined); toast.error(body?.error?.message ?? text.unavailable); return; } const item = await response.json(); if (endpoint.includes("documents")) setDocuments((current) => [...current, item]); else setMedia((current) => [...current, item]); };
  if (loading) return null;
  return <section className="mt-10 border-t border-border pt-8"><h2 className="font-serif text-2xl">{text.title}</h2><div className="mt-5 grid gap-8 lg:grid-cols-2"><div><h3 className="font-medium">{text.document}</h3><div className="mt-3 flex gap-2"><input className="auth-input" value={documentType} onChange={(event) => setDocumentType(event.target.value)} aria-label={text.documentType}/><label className="cursor-pointer"><span className="sr-only">{text.upload}</span><input className="sr-only" type="file" accept="application/pdf,image/jpeg,image/png" onChange={(event) => { const file = event.target.files?.[0]; if (file) void upload("/api/artisan/documents", file, { documentType }); }}/><Button type="button" variant="outline">{text.upload}</Button></label></div><ul className="mt-4 space-y-2 text-sm">{documents.map((item) => <li key={item.id}><a className="underline" href={item.url} target="_blank" rel="noreferrer">{item.originalFilename || item.documentType}</a></li>)}</ul></div><div><h3 className="font-medium">{text.profile}</h3><label className="mt-3 inline-block cursor-pointer"><span className="sr-only">{text.upload}</span><input className="sr-only" type="file" accept="image/jpeg,image/png,image/webp,video/mp4" onChange={(event) => { const file = event.target.files?.[0]; if (file) void upload("/api/artisan/media", file); }}/><Button type="button" variant="outline">{text.upload}</Button></label><ul className="mt-4 space-y-2 text-sm">{media.map((item) => <li key={item.id}><a className="underline" href={item.url} target="_blank" rel="noreferrer">{item.originalFilename || item.mediaKind}</a></li>)}</ul></div></div></section>;
}
