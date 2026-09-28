"use client";

import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";

import type { Category } from "@/features/catalogue/types";
import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { Modal } from "@/core/components/ui/modal";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { commonCopy } from "@/core/lib/common-copy";
import { catalogue } from "@/features/catalogue/api";
import {
  archiveProduct,
  createProduct,
  deleteProductMedia,
  listOwnedProducts,
  listOwnedWorkshops,
  submitProduct,
  updateProduct,
  uploadProductMedia,
  type OwnedWorkshop,
  type ProductDraft,
  type ProductInput,
  type ProductLocale,
  type ProductTranslation,
} from "../api";

const languageOptions = [
  { value: "en", label: "English" },
  { value: "fr", label: "Français" },
  { value: "ar", label: "العربية" },
  { value: "es", label: "Español" },
] satisfies Array<{ value: ProductLocale; label: string }>;
const blankTranslation = (locale: ProductLocale): ProductTranslation => ({
  locale,
  name: "",
  description: "",
  story: "",
  culturalContext: "",
});
const blankProduct = (language: ProductLocale): ProductInput => ({
  workshopId: "",
  categoryId: "",
  productType: "ARTISAN_SPECIFIC",
  plannedQuantity: 1,
  orderTotalMinor: 1,
  priceMinor: 1,
  currency: "EUR",
  materials: "",
  productionMethod: "",
  intendedUse: "",
  dimensions: "",
  weightGrams: undefined,
  countryOfOrigin: "Algeria",
  regionOfOrigin: "",
  ecoFriendlyVerified: false,
  fairTradeVerified: false,
  madeToOrderEligible: false,
  translations: [blankTranslation(language)],
  media: [],
});
const labels = {
  en: {
    title: "Product authoring",
    create: "New draft",
    update: "Update",
    details: "Details",
    save: "Save draft",
    submit: "Submit for review",
    archive: "Archive product",
    archived: "Product archived.",
    media: "Media",
    upload: "Upload media",
    uploading: "Uploading…",
    category: "Category",
    price: "Unit price (minor units)",
    quantity: "Order quantity",
    orderTotal: "Total order price (minor units)",
    approvalCode: "Warehouse approval code",
    mediaLimit: "Up to 4 media files",
    currency: "Currency",
    type: "Product type",
    material: "Materials",
    method: "Production method",
    use: "Intended use",
    dimensions: "Dimensions",
    weight: "Weight (grams)",
    region: "Region of origin",
    locale: "Language",
    name: "Name",
    description: "Description",
    story: "Product story",
    context: "Cultural context",
    empty: "No product drafts yet.",
    saved: "Draft saved.",
    submitted: "Product submitted for review.",
    invalid: "Complete the selected language name and description before submitting.",
    workshop: "Workshop",
    selectWorkshop: "Select a workshop",
    product: "Product",
    status: "Status",
    tablePrice: "Price",
    mediaCount: "Media",
    confirmArchive: "Archive this product?",
    archiveHelp: "Archived products keep their history and are no longer available for sale.",
    close: "Close",
  },
  fr: {
    title: "Création de produits",
    create: "Nouveau brouillon",
    update: "Modifier",
    details: "Détails",
    save: "Enregistrer",
    submit: "Soumettre pour examen",
    archive: "Archiver le produit",
    archived: "Produit archivé.",
    media: "Médias",
    upload: "Importer un média",
    uploading: "Importation…",
    category: "Catégorie",
    price: "Prix unitaire (unités mineures)",
    quantity: "Quantité de la commande",
    orderTotal: "Prix total de la commande (unités mineures)",
    approvalCode: "Code d’approbation entrepôt",
    mediaLimit: "Jusqu’à 4 médias",
    currency: "Devise",
    type: "Type de produit",
    material: "Matières",
    method: "Méthode de production",
    use: "Usage prévu",
    dimensions: "Dimensions",
    weight: "Poids (grammes)",
    region: "Région d’origine",
    locale: "Langue",
    name: "Nom",
    description: "Description",
    story: "Histoire du produit",
    context: "Contexte culturel",
    empty: "Aucun brouillon de produit.",
    saved: "Brouillon enregistré.",
    submitted: "Produit soumis pour examen.",
    invalid: "Complétez le nom et la description dans la langue sélectionnée avant l’envoi.",
    workshop: "Atelier",
    selectWorkshop: "Sélectionner un atelier",
    product: "Produit",
    status: "Statut",
    tablePrice: "Prix",
    mediaCount: "Médias",
    confirmArchive: "Archiver ce produit ?",
    archiveHelp: "Les produits archivés conservent leur historique et ne sont plus vendables.",
    close: "Fermer",
  },
  ar: {
    title: "إنشاء المنتجات",
    create: "مسودة جديدة",
    update: "تحديث",
    details: "التفاصيل",
    save: "حفظ المسودة",
    submit: "إرسال للمراجعة",
    archive: "أرشفة المنتج",
    archived: "تمت أرشفة المنتج.",
    media: "الوسائط",
    upload: "رفع وسائط",
    uploading: "جارٍ الرفع…",
    category: "الفئة",
    price: "السعر الوحدي (الوحدات الصغرى)",
    quantity: "كمية الطلب",
    orderTotal: "السعر الإجمالي للطلب (الوحدات الصغرى)",
    approvalCode: "رمز موافقة المستودع",
    mediaLimit: "حتى 4 وسائط",
    currency: "العملة",
    type: "نوع المنتج",
    material: "المواد",
    method: "طريقة الإنتاج",
    use: "الاستخدام المقصود",
    dimensions: "الأبعاد",
    weight: "الوزن (غرام)",
    region: "منطقة المنشأ",
    locale: "اللغة",
    name: "الاسم",
    description: "الوصف",
    story: "قصة المنتج",
    context: "السياق الثقافي",
    empty: "لا توجد مسودات منتجات.",
    saved: "تم حفظ المسودة.",
    submitted: "تم إرسال المنتج للمراجعة.",
    invalid: "أكمل الاسم والوصف باللغة المحددة قبل الإرسال.",
    workshop: "الورشة",
    selectWorkshop: "اختر ورشة",
    product: "المنتج",
    status: "الحالة",
    tablePrice: "السعر",
    mediaCount: "الوسائط",
    confirmArchive: "أرشفة هذا المنتج؟",
    archiveHelp: "تحتفظ المنتجات المؤرشفة بسجلها ولا تعود متاحة للبيع.",
    close: "إغلاق",
  },
  es: {
    title: "Creación de productos",
    create: "Nuevo borrador",
    update: "Actualizar",
    details: "Detalles",
    save: "Guardar borrador",
    submit: "Enviar a revisión",
    archive: "Archivar producto",
    archived: "Producto archivado.",
    media: "Medios",
    upload: "Subir medio",
    uploading: "Subiendo…",
    category: "Categoría",
    price: "Precio unitario (unidades menores)",
    quantity: "Cantidad del pedido",
    orderTotal: "Precio total del pedido (unidades menores)",
    approvalCode: "Código de aprobación del almacén",
    mediaLimit: "Hasta 4 medios",
    currency: "Moneda",
    type: "Tipo de producto",
    material: "Materiales",
    method: "Método de producción",
    use: "Uso previsto",
    dimensions: "Dimensiones",
    weight: "Peso (gramos)",
    region: "Región de origen",
    locale: "Idioma",
    name: "Nombre",
    description: "Descripción",
    story: "Historia del producto",
    context: "Contexto cultural",
    empty: "Aún no hay borradores de productos.",
    saved: "Borrador guardado.",
    submitted: "Producto enviado a revisión.",
    invalid:
      "Completa el nombre y la descripción en el idioma seleccionado antes de enviarlo.",
    workshop: "Taller",
    selectWorkshop: "Selecciona un taller",
    product: "Producto",
    status: "Estado",
    tablePrice: "Precio",
    mediaCount: "Medios",
    confirmArchive: "¿Archivar este producto?",
    archiveHelp: "Los productos archivados conservan su historial y dejan de estar disponibles para la venta.",
    close: "Cerrar",
  },
} as const;

export function ProductWorkspace({ locale }: { locale: ProductLocale }) {
  const text = labels[locale];
  const common = commonCopy(locale);
  const [items, setItems] = useState<ProductDraft[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [workshops, setWorkshops] = useState<OwnedWorkshop[]>([]);
  const [selected, setSelected] = useState<ProductDraft>();
  const [detailsItem, setDetailsItem] = useState<ProductDraft>();
  const [pendingArchive, setPendingArchive] = useState<ProductDraft>();
  const [dialog, setDialog] = useState<"editor" | "details" | "archive" | null>(null);
  const [editorLocale, setEditorLocale] = useState<ProductLocale>(locale);
  const [form, setForm] = useState<ProductInput>(() => blankProduct(locale));
  const [busy, setBusy] = useState(false);
  const mediaInput = useRef<HTMLInputElement>(null);

  const closeDialog = (force = false) => {
    if (busy && !force) return;
    setDialog(null);
    setDetailsItem(undefined);
    setPendingArchive(undefined);
  };

  const refreshWorkshops = async () => {
    try {
      const response = await listOwnedWorkshops();
      setWorkshops(response.items);
    } catch (error) {
      toast.error((error as Error).message);
    }
  };

  useEffect(() => {
    void Promise.all([
      listOwnedProducts(),
      listOwnedWorkshops(),
      catalogue("en"),
    ])
      .then(([products, ownedWorkshops, data]) => {
        setItems(products.items);
        setWorkshops(ownedWorkshops.items);
        setCategories(data.categories);
      })
      .catch((error: Error) => toast.error(error.message));
  }, []);

  const select = (item?: ProductDraft) => {
    setSelected(item);
    const language = item?.translations.find((translation) => translation.locale === locale)?.locale ?? item?.translations[0]?.locale ?? locale;
    setEditorLocale(language);
    setForm(
      item
        ? {
            ...item,
            plannedQuantity: item.plannedQuantity || 1,
            orderTotalMinor: item.orderTotalMinor || item.priceMinor,
            translations: [item.translations.find((translation) => translation.locale === language) ?? blankTranslation(language)],
          }
        : blankProduct(locale),
    );
  };
  const openEditor = (item?: ProductDraft) => {
    select(item);
    void refreshWorkshops();
    setDialog("editor");
  };
  const change = <K extends keyof ProductInput>(
    key: K,
    value: ProductInput[K],
  ) => setForm((current) => ({ ...current, [key]: value }));
  const changeTranslation = (
    language: ProductLocale,
    key: keyof ProductTranslation,
    value: string,
  ) =>
    setForm((current) => ({
      ...current,
      translations: [{ ...(current.translations[0] ?? blankTranslation(language)), locale: language, [key]: value }],
    }));
  const changeLanguage = (language: ProductLocale) => {
    setEditorLocale(language);
    setForm((current) => ({
      ...current,
      translations: [{ ...(current.translations[0] ?? blankTranslation(language)), locale: language }],
    }));
  };
  const save = async () => {
    setBusy(true);
    try {
      const saved = selected
        ? await updateProduct(selected.id, form)
        : await createProduct(form);
      setItems((current) => [
        saved,
        ...current.filter((item) => item.id !== saved.id),
      ]);
      select(saved);
      toast.success(text.saved);
    } catch (error) {
      toast.error((error as Error).message);
    } finally {
      setBusy(false);
    }
  };
  const archive = async () => {
    const item = pendingArchive ?? selected;
    if (!item || item.status === "ARCHIVED") return;
    setBusy(true);
    try {
      const archived = await archiveProduct(item.id);
      setItems((current) =>
        current.map((item) => (item.id === archived.id ? archived : item)),
      );
      if (selected?.id === archived.id) select(archived);
      closeDialog(true);
      toast.success(text.archived);
    } catch (error) {
      toast.error((error as Error).message);
    } finally {
      setBusy(false);
    }
  };
  const submit = async () => {
    if (
      form.translations.some(
        (item) => !item.name.trim() || !item.description.trim(),
      )
    ) {
      toast.error(text.invalid);
      return;
    }
    if (!selected) {
      toast.error(text.saved);
      return;
    }
    setBusy(true);
    try {
      const submitted = await submitProduct(selected.id);
      setItems((current) =>
        current.map((item) => (item.id === submitted.id ? submitted : item)),
      );
      select(submitted);
      toast.success(text.submitted);
    } catch (error) {
      toast.error((error as Error).message);
    } finally {
      setBusy(false);
    }
  };
  const upload = async (file: File) => {
    if (!selected) return;
    if (selected.media.length >= 4) {
      toast.error(text.mediaLimit);
      return;
    }
    setBusy(true);
    try {
      const media = await uploadProductMedia(selected.id, file, file.name);
      const next = { ...selected, media: [...selected.media, media] };
      setSelected(next);
      setForm((current) => ({ ...current, media: next.media }));
      toast.success(text.upload);
    } catch (error) {
      toast.error((error as Error).message);
    } finally {
      setBusy(false);
    }
  };
  const removeMedia = async (mediaId: string) => {
    if (!selected || selected.status === "ARCHIVED") return;
    const response = await deleteProductMedia(selected.id, mediaId);
    if (!response.ok) {
      toast.error(text.invalid);
      return;
    }
    const media = selected.media.filter((item) => item.id !== mediaId);
    const next = { ...selected, media };
    setSelected(next);
    setForm((current) => ({ ...current, media }));
  };
  const activeTranslation = form.translations[0] ?? blankTranslation(editorLocale);
  const activeLanguage = activeTranslation.locale;

  return (
    <>
    <section className="mt-12 border-t border-border pt-8">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-xs uppercase tracking-widest text-primary">
            AISHA
          </p>
          <h2 className="mt-2 font-serif text-3xl">{text.title}</h2>
        </div>
        <Button type="button" variant="outline" onClick={() => openEditor()}>
          {text.create}
        </Button>
      </div>
      <div className="mt-6">
        <div className="overflow-x-auto rounded-lg border border-border xl:col-span-2">
          {items.length ? <table className="w-full min-w-[1080px] text-sm">
            <thead className="bg-muted/40"><tr><th className="px-4 py-3 text-start font-medium">{text.details}</th><th className="px-4 py-3 text-start font-medium">{text.product}</th><th className="px-4 py-3 text-start font-medium">{text.workshop}</th><th className="px-4 py-3 text-start font-medium">{text.status}</th><th className="px-4 py-3 text-start font-medium">{text.tablePrice}</th><th className="px-4 py-3 text-start font-medium">{text.quantity}</th><th className="px-4 py-3 text-start font-medium">{text.orderTotal}</th><th className="px-4 py-3 text-start font-medium">{text.mediaCount}</th><th className="px-4 py-3 text-start font-medium">{text.approvalCode}</th></tr></thead>
            <tbody>{items.map((item) => {
              const name = item.translations.find((translation) => translation.locale === locale)?.name || item.translations[0]?.name || "—";
              return <tr className="border-t border-border align-top" key={item.id}>
                <td className="px-4 py-3"><div className="flex flex-wrap gap-2"><Button type="button" variant="outline" className="min-h-10 px-3" onClick={() => { setDetailsItem(item); setDialog("details"); }}>{text.details}</Button><Button type="button" variant="outline" className="min-h-10 px-3" onClick={() => openEditor(item)}>{text.update}</Button>{item.status !== "ARCHIVED" && <Button type="button" variant="destructive" className="min-h-10 px-3" disabled={busy} onClick={() => { setPendingArchive(item); setDialog("archive"); }}>{text.archive}</Button>}</div></td>
                <td className="px-4 py-3 font-medium">{name}</td><td className="px-4 py-3">{item.workshopName || "—"}</td><td className="px-4 py-3"><StatusBadge status={item.status} locale={locale} /></td><td className="px-4 py-3">{item.priceMinor} {item.currency}</td><td className="px-4 py-3">{item.plannedQuantity || 1}</td><td className="px-4 py-3">{item.orderTotalMinor || item.priceMinor} {item.currency}</td><td className="px-4 py-3">{item.media.length}/4</td><td className="px-4 py-3 font-mono text-xs">{item.productCode || "—"}</td>
              </tr>;
            })}</tbody>
          </table> : <p className="p-4 text-sm text-muted-foreground">{text.empty}</p>}
        </div>
        {dialog === "editor" && <Modal label={selected ? text.update : text.create} closeLabel={text.close} onClose={closeDialog} panelClassName="w-full p-6 sm:p-10 lg:max-w-6xl"><div className="space-y-5">
          <div className="grid gap-4 sm:grid-cols-2">
            <Field label={text.workshop}>
              <Combobox
                options={workshops
                  .filter((workshop) => workshop.status === "ACTIVE")
                  .map((workshop) => ({ value: workshop.id, label: workshop.name }))}
                value={form.workshopId}
                onChange={(next) => change("workshopId", next)}
                placeholder={text.selectWorkshop}
                ariaLabel={text.workshop}
              />
            </Field>
            <Field label={text.category}>
              <Combobox
                options={categories
                  .filter((item) => item.id)
                  .map((item) => ({ value: item.id ?? "", label: item.name.en }))}
                value={form.categoryId}
                onChange={(next) => change("categoryId", next)}
                placeholder={common.chooseCategory}
                ariaLabel={text.category}
              />
            </Field>
            <Field label={text.type}>
              <Combobox
                options={[
                  { value: "ARTISAN_SPECIFIC", label: common.artisanSpecific },
                  { value: "STANDARD_TRADITIONAL", label: common.standardTraditional },
                ]}
                value={form.productType}
                onChange={(next) => change("productType", next as ProductInput["productType"])}
                ariaLabel={text.type}
              />
            </Field>
            <Field label={text.price}>
              <input
                className="auth-input"
                type="number"
                min="1"
                value={form.priceMinor}
                onChange={(event) =>
                  change("priceMinor", Number(event.target.value))
                }
              />
            </Field>
            <Field label={text.quantity}>
              <input
                className="auth-input"
                type="number"
                min="1"
                value={form.plannedQuantity}
                onChange={(event) =>
                  change("plannedQuantity", Number(event.target.value))
                }
              />
            </Field>
            <Field label={text.orderTotal}>
              <input
                className="auth-input"
                type="number"
                min="1"
                value={form.orderTotalMinor}
                onChange={(event) =>
                  change("orderTotalMinor", Number(event.target.value))
                }
              />
            </Field>
            <Field label={text.currency}>
              <input
                className="auth-input"
                maxLength={3}
                value={form.currency}
                onChange={(event) =>
                  change("currency", event.target.value.toUpperCase())
                }
              />
            </Field>
            {(
              [
                ["materials", text.material],
                ["productionMethod", text.method],
                ["intendedUse", text.use],
                ["dimensions", text.dimensions],
                ["regionOfOrigin", text.region],
              ] as const
            ).map(([key, label]) => (
              <Field label={label} key={key}>
                <input
                  className="auth-input"
                  value={form[key]}
                  onChange={(event) => change(key, event.target.value)}
                />
              </Field>
            ))}
          </div>
          <Field label={text.locale}>
            <Combobox
              options={languageOptions}
              value={activeLanguage}
              onChange={(next) => changeLanguage(next as ProductLocale)}
              ariaLabel={text.locale}
            />
          </Field>
          <fieldset
            className="space-y-3 border border-border p-4"
            dir={activeLanguage === "ar" ? "rtl" : "ltr"}
          >
            <legend className="px-2 text-sm font-medium">
              {text.locale}: {activeLanguage}
            </legend>
            <Field label={text.name}>
              <input
                className="auth-input"
                value={activeTranslation.name}
                onChange={(event) =>
                  changeTranslation(activeLanguage, "name", event.target.value)
                }
              />
            </Field>
            <Field label={text.description}>
              <textarea
                className="auth-input min-h-24"
                value={activeTranslation.description}
                onChange={(event) =>
                  changeTranslation(activeLanguage, "description", event.target.value)
                }
              />
            </Field>
            <Field label={text.story}>
              <textarea
                className="auth-input min-h-20"
                value={activeTranslation.story}
                onChange={(event) =>
                  changeTranslation(activeLanguage, "story", event.target.value)
                }
              />
            </Field>
            <Field label={text.context}>
              <textarea
                className="auth-input min-h-20"
                value={activeTranslation.culturalContext}
                onChange={(event) =>
                  changeTranslation(activeLanguage, "culturalContext", event.target.value)
                }
              />
            </Field>
          </fieldset>
          <div className="border border-border p-4">
            <h3 className="font-medium">{text.media}</h3>
            <p className="mt-1 text-sm text-muted-foreground">{text.mediaLimit}</p>
            <div className="mt-4 flex flex-wrap gap-3">
              {form.media?.map((item) => (
                <div
                  key={item.id}
                  className="flex items-center gap-2 border border-border px-3 py-2 text-sm"
                >
                  <span>{item.originalFilename}</span>
                  <button
                    type="button"
                    className="text-destructive"
                    onClick={() => void removeMedia(item.id)}
                  >
                    ×
                  </button>
                </div>
              ))}
              <input
                ref={mediaInput}
                className="sr-only"
                type="file"
                accept="image/jpeg,image/png,image/webp,video/mp4"
                disabled={!selected || busy}
                onChange={(event) => {
                  const file = event.target.files?.[0];
                  event.currentTarget.value = "";
                  if (file) void upload(file);
                }}
              />
              <Button
                type="button"
                variant="outline"
                disabled={!selected || busy || selected.status === "ARCHIVED" || selected.media.length >= 4}
                onClick={() => mediaInput.current?.click()}
              >
                {busy ? text.uploading : text.upload}
              </Button>
            </div>
          </div>
          <div className="flex flex-wrap gap-3">
            <Button
              type="button"
              disabled={busy || selected?.status === "ARCHIVED"}
              onClick={() => void save()}
            >
              {text.save}
            </Button>
            <Button
              type="button"
              disabled={
                busy || !selected || selected.status === "PENDING_REVIEW"
              }
              variant="outline"
              onClick={() => void submit()}
            >
              {text.submit}
            </Button>
            <Button
              type="button"
              disabled={
                busy ||
                !selected ||
                selected.status === "ARCHIVED" ||
                selected.status === "PENDING_REVIEW"
              }
              variant="outline"
              onClick={() => void archive()}
            >
              {text.archive}
            </Button>
          </div>
        </div></Modal>}
      </div>
      </section>
      {dialog === "details" && detailsItem && <Modal label={text.details} closeLabel={text.close} onClose={closeDialog} panelClassName="w-full p-6 sm:p-10 lg:max-w-3xl"><div className="min-h-full"><p className="text-xs uppercase tracking-widest text-primary">{text.title}</p><h2 className="mt-2 font-serif text-3xl">{detailsItem.translations.find((translation) => translation.locale === locale)?.name || detailsItem.translations[0]?.name || "—"}</h2><dl className="mt-8 divide-y divide-border border-y border-border"><ProductDetail label={text.workshop} value={detailsItem.workshopName} /><ProductDetail label={text.status} value={detailsItem.status} /><ProductDetail label={text.tablePrice} value={`${detailsItem.priceMinor} ${detailsItem.currency}`} /><ProductDetail label={text.material} value={detailsItem.materials} /><ProductDetail label={text.method} value={detailsItem.productionMethod} /><ProductDetail label={text.use} value={detailsItem.intendedUse} /></dl><h3 className="mt-8 font-medium">{text.media}</h3><div className="mt-4 grid gap-4 sm:grid-cols-2">{detailsItem.media.length ? detailsItem.media.map((media) => <div className="border border-border p-3" key={media.id}>{media.url && media.mediaType.startsWith("image/") && <div role="img" aria-label={media.altText || media.originalFilename} className="aspect-video w-full bg-muted bg-cover bg-center" style={{ backgroundImage: `url("${media.url}")` }} />}{media.url && media.mediaType.startsWith("video/") && <video className="aspect-video w-full object-cover" controls preload="metadata" src={media.url} />}{!media.url && <div className="flex aspect-video items-center justify-center bg-muted text-xs text-muted-foreground">{media.mediaType}</div>}<a className="mt-2 block truncate text-sm underline" href={media.url} target="_blank" rel="noreferrer">{media.originalFilename}</a></div>) : <p className="text-sm text-muted-foreground">—</p>}</div><div className="mt-8"><Button type="button" variant="outline" onClick={() => closeDialog()}>{text.close}</Button></div></div></Modal>}
      {dialog === "archive" && pendingArchive && <Modal label={text.archive} closeLabel={text.close} onClose={closeDialog} panelClassName="w-full max-w-lg p-6 sm:p-10"><div><h2 className="font-serif text-3xl">{text.confirmArchive}</h2><p className="mt-3 text-muted-foreground">{pendingArchive.translations.find((translation) => translation.locale === locale)?.name || pendingArchive.translations[0]?.name || "—"}</p><p className="mt-2 text-sm text-muted-foreground">{text.archiveHelp}</p><div className="mt-8 flex flex-wrap gap-3"><Button type="button" variant="destructive" disabled={busy} onClick={() => void archive()}>{busy ? "…" : text.archive}</Button><Button type="button" variant="outline" disabled={busy} onClick={() => closeDialog()}>{text.close}</Button></div></div></Modal>}
    </>
  );
}

function ProductDetail({ label, value }: { label: string; value?: string }) {
  return <div className="grid gap-2 py-4 sm:grid-cols-[10rem_1fr]"><dt className="text-sm text-muted-foreground">{label}</dt><dd>{value || "—"}</dd></div>;
}

function Field({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <label className="block">
      <span className="mb-1 block text-sm">{label}</span>
      {children}
    </label>
  );
}
