/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useCallback, useEffect, useState } from "react";
import { Eye, Plus } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { Combobox } from "@/core/components/ui/combobox";
import { FileDropzone } from "@/core/components/ui/file-dropzone";
import { IconAction } from "@/core/components/ui/icon-action";
import { Modal } from "@/core/components/ui/modal";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { UploadProgress } from "@/core/components/ui/upload-progress";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import type { Locale } from "@/core/lib/i18n";

type ReceptionStatus = "RECEIVED_PENDING_INSPECTION" | "INSPECTED";
type Evidence = {
  id: string;
  originalFilename: string;
  mediaType: string;
  sizeBytes: number;
  url?: string;
};
type Inspection = {
  acceptedQuantity: number;
  rejectedQuantity: number;
  quarantinedQuantity: number;
  damagedQuantity: number;
  reason: string;
};
type Reception = {
  id: string;
  productId: string;
  productCode: string;
  productName: string;
  artisanName: string;
  artisanPhone?: string;
  workshopName: string;
  receivedQuantity: number;
  referenceKey: string;
  parcelReference: string;
  supplierName: string;
  notes: string;
  status: ReceptionStatus;
  receivedAt: string;
  inspection?: Inspection;
  evidence: Evidence[];
};
type ValidatedProduct = {
  productCode: string;
  productName: string;
  productStatus: string;
  artisanName: string;
  artisanPhone: string;
  workshopId: string;
  workshopName: string;
  priceMinor: number;
  currency: string;
  availableQuantity: number;
};

const statuses: Array<"" | ReceptionStatus> = [
  "RECEIVED_PENDING_INSPECTION",
  "INSPECTED",
];

export function WarehouseOperations({ locale = "en" }: { locale?: Locale }) {
  const text = {
    artisanPhone: locale === "fr" ? "Téléphone de l’artisan" : locale === "ar" ? "هاتف الحرفي" : locale === "es" ? "Teléfono del artesano" : "Artisan phone",
    phonePlaceholder: locale === "fr" ? "ex. 0550123456" : locale === "ar" ? "مثال: 0550123456" : locale === "es" ? "p. ej. 0550123456" : "e.g. 0550123456",
    findProducts: locale === "fr" ? "Rechercher des produits" : locale === "ar" ? "البحث عن المنتجات" : locale === "es" ? "Buscar productos" : "Find products",
    searching: locale === "fr" ? "Recherche…" : locale === "ar" ? "جارٍ البحث…" : locale === "es" ? "Buscando…" : "Searching…",
    workshop: locale === "fr" ? "Atelier" : locale === "ar" ? "الورشة" : locale === "es" ? "Taller" : "Workshop",
    selectWorkshop: locale === "fr" ? "Sélectionner un atelier" : locale === "ar" ? "اختر ورشة" : locale === "es" ? "Selecciona un taller" : "Select a workshop",
    searchArtisanFirst: locale === "fr" ? "Recherchez d’abord un artisan" : locale === "ar" ? "ابحث عن حرفي أولاً" : locale === "es" ? "Busca primero un artesano" : "Search an artisan first",
    productCode: locale === "fr" ? "Code ou nom du produit" : locale === "ar" ? "رمز المنتج أو اسمه" : locale === "es" ? "Código o nombre del producto" : "Product code or name",
    selectProduct: locale === "fr" ? "Sélectionner un produit approuvé" : locale === "ar" ? "اختر منتجاً معتمداً" : locale === "es" ? "Selecciona un producto aprobado" : "Select an approved product",
    selectWorkshopFirst: locale === "fr" ? "Sélectionnez d’abord un atelier" : locale === "ar" ? "اختر ورشة أولاً" : locale === "es" ? "Selecciona primero un taller" : "Select a workshop first",
  };
  const [items, setItems] = useState<Reception[]>([]);
  const [selected, setSelected] = useState<Reception | null>(null);
  const [createOpen, setCreateOpen] = useState(false);
  useEscapeKey(() => setSelected(null), selected !== null);
  const [status, setStatus] = useState<"" | ReceptionStatus>(
    "RECEIVED_PENDING_INSPECTION",
  );
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [validatedProducts, setValidatedProducts] = useState<ValidatedProduct[]>([]);
  const [artisanPhone, setArtisanPhone] = useState("");
  const [workshopId, setWorkshopId] = useState("");
  const [productsLoading, setProductsLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [createForm, setCreateForm] = useState({
    productCode: "",
    receivedQuantity: "",
    referenceKey: "",
    supplierName: "",
    parcelReference: "",
    notes: "",
  });
  const [inspection, setInspection] = useState({
    acceptedQuantity: "",
    rejectedQuantity: "0",
    quarantinedQuantity: "0",
    damagedQuantity: "0",
    reason: "",
  });
  const [evidence, setEvidence] = useState<File | null>(null);
  const pageSize = 10;
  const pageCount = Math.max(1, Math.ceil(total / pageSize));

  const loadValidatedProducts = useCallback(async (nextWorkshopId = workshopId) => {
    if (!artisanPhone.trim()) {
      toast.error("Enter the artisan phone number first");
      return;
    }
    setProductsLoading(true);
    try {
      const query = new URLSearchParams({ artisanPhone: artisanPhone.trim(), page: "1", pageSize: "100" });
      if (nextWorkshopId) query.set("workshopId", nextWorkshopId);
      const response = await fetch(`/api/warehouse/products?${query}`, { cache: "no-store" });
      if (!response.ok) throw new Error();
      const data = (await response.json()) as { items?: ValidatedProduct[] };
      setValidatedProducts(data.items ?? []);
      if (nextWorkshopId && !(data.items ?? []).some((item) => item.workshopId === nextWorkshopId)) {
        setWorkshopId("");
      }
      if (!(data.items ?? []).length) toast.info("No approved products found for this artisan");
    } catch {
      toast.error("Unable to find validated products");
    } finally {
      setProductsLoading(false);
    }
  }, [artisanPhone, workshopId]);

  const workshops = Array.from(new Map(validatedProducts.map((item) => [item.workshopId, item.workshopName])).entries()).map(([value, label]) => ({ value, label }));
  const workshopProducts = validatedProducts.filter((item) => !workshopId || item.workshopId === workshopId);
  const productOptions = workshopProducts
    .map((item) => ({ value: item.productCode, label: `${item.productCode} · ${item.productName} · ${item.availableQuantity} available` }));
  const selectedProduct = validatedProducts.find((item) => item.productCode === createForm.productCode);
  const inspectionTotal =
    Number(inspection.acceptedQuantity || 0) +
    Number(inspection.rejectedQuantity || 0) +
    Number(inspection.quarantinedQuantity || 0) +
    Number(inspection.damagedQuantity || 0);

  const load = useCallback(
    async (nextPage = page, nextStatus = status) => {
      setLoading(true);
      try {
        const query = new URLSearchParams({
          page: String(nextPage),
          pageSize: String(pageSize),
        });
        if (nextStatus) query.set("status", nextStatus);
        const response = await fetch(`/api/warehouse/receptions?${query}`, {
          cache: "no-store",
        });
        if (!response.ok) throw new Error();
        const data = (await response.json()) as {
          items?: Reception[];
          page?: number;
          total?: number;
        };
        setItems(data.items ?? []);
        setPage(data.page ?? nextPage);
        setTotal(data.total ?? 0);
      } catch {
        toast.error("Unable to load warehouse receptions");
      } finally {
        setLoading(false);
      }
    },
    [page, status],
  );

  useEffect(() => {
    void load();
  }, [load]);


  const createReception = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSaving(true);
    try {
      const response = await fetch("/api/warehouse/receptions", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          ...createForm,
          receivedQuantity: Number(createForm.receivedQuantity),
        }),
      });
      if (!response.ok) {
        const body = await response.json().catch(() => ({}));
        throw new Error(body.error?.message ?? "Unable to record reception");
      }
      const reception = (await response.json()) as Reception;
      setCreateForm({
        productCode: "",
        receivedQuantity: "",
        referenceKey: "",
        supplierName: "",
        parcelReference: "",
        notes: "",
      });
      setValidatedProducts([]);
      setArtisanPhone("");
      setWorkshopId("");
      setCreateOpen(false);
      selectItem(reception);
      toast.success("Reception recorded — complete the inspection to publish stock");
      await load(1, status);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Unable to record reception");
    } finally {
      setSaving(false);
    }
  };

  const uploadEvidence = async (): Promise<Evidence | null> => {
    if (!selected || !evidence) return null;
    const form = new FormData();
    form.set("file", evidence);
    const response = await fetch(`/api/warehouse/receptions/${selected.id}/evidence`, {
      method: "POST",
      body: form,
    });
    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      throw new Error(body.error?.message ?? "Unable to upload evidence");
    }
    return (await response.json()) as Evidence;
  };

  const inspect = async () => {
    if (!selected) return;
    setSaving(true);
    try {
      let selectedForInspection = selected;
      if (evidence) {
        const uploaded = await uploadEvidence();
        if (uploaded) {
          selectedForInspection = { ...selected, evidence: [...selected.evidence, uploaded] };
          setSelected(selectedForInspection);
          setEvidence(null);
        }
      }
      if (selectedForInspection.evidence.length === 0) {
        throw new Error("Select an evidence file before submitting the inspection");
      }
      const response = await fetch(`/api/warehouse/receptions/${selected.id}/inspect`, {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          acceptedQuantity: Number(inspection.acceptedQuantity),
          rejectedQuantity: Number(inspection.rejectedQuantity),
          quarantinedQuantity: Number(inspection.quarantinedQuantity),
          damagedQuantity: Number(inspection.damagedQuantity),
          reason: inspection.reason,
        }),
      });
      if (!response.ok) {
        const body = await response.json().catch(() => ({}));
        throw new Error(body.error?.message ?? "Unable to inspect reception");
      }
      toast.success("Inspection recorded");
      setSelected(null);
      await load();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Unable to inspect reception");
    } finally {
      setSaving(false);
    }
  };

  const selectItem = (item: Reception) => {
    setCreateOpen(false);
    setSelected(item);
    setInspection({
      acceptedQuantity: String(item.receivedQuantity),
      rejectedQuantity: "0",
      quarantinedQuantity: "0",
      damagedQuantity: "0",
      reason: "",
    });
  };

  return (
    <div className="space-y-6">
      <section className="rounded border border-border bg-card p-4 shadow-sm sm:p-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:gap-5">
          <Button type="button" className="order-first shrink-0" onClick={() => setCreateOpen(true)}>
            <Plus aria-hidden="true" size={17} />
            New reception
          </Button>
          <div className="min-w-0">
            <p className="text-xs uppercase tracking-[0.18em] text-primary">Warehouse intake</p>
            <h2 className="mt-1 font-serif text-3xl">Reception ledger</h2>
            <p className="mt-2 max-w-2xl text-sm text-muted-foreground">
              Record approved artisan stock, then inspect the batch before it becomes available to customers.
            </p>
          </div>
        </div>
        <div className="mt-6 flex flex-col gap-4 border-t border-border pt-5 sm:flex-row sm:items-end sm:justify-between">
          <label className="block w-full text-sm sm:max-w-xs">
            <span className="mb-2 block">Reception status</span>
            <Combobox
              className="w-full"
              options={[{ value: "", label: "All receptions" }, ...statuses.filter(Boolean).map((value) => ({ value, label: value.replaceAll("_", " ") }))]}
              value={status}
              onChange={(next) => {
                const nextStatus = next as "" | ReceptionStatus;
                setStatus(nextStatus);
                setPage(1);
                void load(1, nextStatus);
              }}
              ariaLabel="Reception status"
            />
          </label>
          <span className="text-sm text-muted-foreground">{total} reception(s)</span>
        </div>
      </section>

      {loading ? (
        <p className="border border-border p-5 text-sm text-muted-foreground" role="status">
          Loading receptions…
        </p>
      ) : (
        <>
          <section className="rounded border border-border bg-card p-3 shadow-sm sm:p-5">
            <div className="mb-4 flex items-end justify-between gap-3">
              <div>
                <p className="text-xs uppercase tracking-[0.18em] text-primary">Stock movement</p>
                <h3 className="mt-1 font-serif text-2xl">Recent receptions</h3>
              </div>
              <p className="hidden text-xs text-muted-foreground sm:block">Scroll horizontally to view all columns</p>
            </div>
          <div className="max-w-full touch-pan-x overflow-x-auto overscroll-x-contain rounded border border-border">
            <table className="w-full min-w-[60rem] text-start text-sm">
              <thead className="border-b border-border bg-muted/40">
                <tr>
                  <th className="px-4 py-3 font-medium">Product</th>
                  <th className="px-4 py-3 font-medium">Workshop</th>
                  <th className="px-4 py-3 font-medium">Quantity</th>
                  <th className="px-4 py-3 font-medium">Reference</th>
                  <th className="px-4 py-3 font-medium">Status</th>
                  <th className="px-4 py-3 font-medium">Details</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {items.map((item) => (
                  <tr key={item.id}>
                    <td className="px-4 py-4">{item.productName}<br /><span className="text-xs text-muted-foreground">{item.productCode || "Validated product"}</span></td>
                    <td className="px-4 py-4">{item.workshopName}</td>
                    <td className="px-4 py-4">{item.receivedQuantity}</td>
                    <td className="px-4 py-4">{item.referenceKey}</td>
                    <td className="px-4 py-4"><StatusBadge status={item.status} locale={locale} /></td>
                    <td className="px-4 py-4"><IconAction icon={<Eye size={17} />} label="Details" onClick={() => selectItem(item)} /></td>
                  </tr>
                ))}
                {items.length === 0 && <tr><td className="px-4 py-8 text-muted-foreground" colSpan={6}>No receptions found.</td></tr>}
              </tbody>
            </table>
          </div>
          <div className="mt-4 flex flex-col-reverse items-start gap-3 text-sm sm:flex-row sm:items-center sm:justify-between">
            <span>Page {page} of {pageCount}</span>
            <div className="flex gap-2">
              <Button type="button" variant="outline" disabled={page <= 1} onClick={() => void load(page - 1)}>Previous</Button>
              <Button type="button" variant="outline" disabled={page >= pageCount} onClick={() => void load(page + 1)}>Next</Button>
            </div>
          </div>
          </section>
        </>
      )}

      {createOpen && (
        <Modal
          label="New reception"
          closeLabel="Close new reception"
          onClose={() => setCreateOpen(false)}
          panelClassName="mx-auto my-4 min-h-0 max-h-[calc(100svh-2rem)] max-w-4xl overflow-y-auto border border-border p-4 shadow-2xl sm:p-6"
        >
          <form onSubmit={createReception}>
            <div className="mb-6 border-b border-border pb-5 pe-10">
              <p className="text-xs uppercase tracking-[0.18em] text-primary">Receive stock</p>
              <h2 className="mt-2 font-serif text-3xl">New reception</h2>
              <p className="mt-2 text-sm text-muted-foreground">
                Received units remain unavailable until an inspection is committed.
              </p>
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="rounded border border-primary/30 bg-secondary/30 p-4 md:col-span-2">
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <p className="text-sm font-medium">Find a validated product</p>
                    <p className="mt-1 text-xs text-muted-foreground">Search by artisan phone, choose the workshop, then select the moderator-approved product code.</p>
                  </div>
                  <span className="rounded-full bg-background px-3 py-1 text-[0.6875rem] uppercase tracking-wider text-muted-foreground">Step 1</span>
                </div>
                <div className="mt-4 grid gap-4 lg:grid-cols-[minmax(0,1fr)_auto]">
                  <label className="block text-sm">
                    <span className="mb-2 block">{text.artisanPhone}</span>
                    <input className="auth-input w-full" type="tel" aria-label={text.artisanPhone} value={artisanPhone} placeholder={text.phonePlaceholder} onChange={(event) => setArtisanPhone(event.target.value)} />
                  </label>
                  <Button className="self-end" type="button" variant="outline" disabled={productsLoading || !artisanPhone.trim()} onClick={() => void loadValidatedProducts("")}>
                    {productsLoading ? text.searching : text.findProducts}
                  </Button>
                </div>
                <div className="mt-4 grid gap-4 md:grid-cols-2">
                  <label className="block text-sm">
                    <span className="mb-2 block">{text.workshop}</span>
                    <Combobox className="w-full" options={workshops} value={workshopId} onChange={(next) => { setWorkshopId(next); setCreateForm({ ...createForm, productCode: "" }); }} placeholder={text.selectWorkshop} emptyMessage={text.searchArtisanFirst} ariaLabel={text.workshop} disabled={!workshops.length} />
                  </label>
                  <label className="block text-sm">
                    <span className="mb-2 block">{text.productCode}</span>
                    <Combobox className="w-full" options={productOptions} value={createForm.productCode} onChange={(next) => {
                      const product = validatedProducts.find((item) => item.productCode === next);
                      setCreateForm({ ...createForm, productCode: next, supplierName: product?.artisanName ?? createForm.supplierName });
                    }} placeholder={text.selectProduct} emptyMessage={text.selectWorkshopFirst} ariaLabel={text.productCode} disabled={!productOptions.length} />
                  </label>
                </div>
                <div className="mt-3 flex flex-wrap gap-2 text-xs text-muted-foreground">
                  <span>{workshopProducts.length} approved product(s) in this workshop</span>
                  {productsLoading && <span>Refreshing…</span>}
                </div>
                {selectedProduct && <p className="mt-3 rounded bg-background px-3 py-2 text-xs text-muted-foreground">{selectedProduct.artisanName} · {selectedProduct.artisanPhone} · {selectedProduct.workshopName} · <span className="font-medium text-foreground">{selectedProduct.productCode}</span> · {selectedProduct.availableQuantity} available</p>}
              </div>
              <div className="md:col-span-2">
                <p className="mb-1 text-sm font-medium">Reception details</p>
                <p className="text-xs text-muted-foreground">Add the quantity and traceability information for this delivery.</p>
              </div>
              {[
                ["receivedQuantity", "Received quantity", "number"],
                ["referenceKey", "Reception reference", "text"],
                ["supplierName", "Supplier or artisan", "text"],
                ["parcelReference", "Parcel or batch reference", "text"],
              ].map(([key, label, type]) => (
                <label key={key} className="block text-sm">
                  <span className="mb-2 block">{label}</span>
                  <input
                    className="auth-input w-full"
                    type={type}
                    min={type === "number" ? 1 : undefined}
                    required={key === "receivedQuantity" || key === "referenceKey"}
                    value={createForm[key as keyof typeof createForm]}
                    onChange={(event) => setCreateForm({ ...createForm, [key]: event.target.value })}
                  />
                </label>
              ))}
              <label className="block text-sm md:col-span-2">
                <span className="mb-2 block">Notes</span>
                <textarea className="auth-input min-h-24 w-full" value={createForm.notes} onChange={(event) => setCreateForm({ ...createForm, notes: event.target.value })} />
              </label>
            </div>
            <div className="mt-6 flex flex-col-reverse gap-3 border-t border-border pt-5 sm:flex-row sm:justify-end">
              <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>Cancel</Button>
              <Button type="submit" disabled={saving || !createForm.productCode}>
                {saving ? "Recording…" : "Record reception"}
              </Button>
            </div>
          </form>
        </Modal>
      )}

      {selected && (
        <div className="fixed inset-0 z-50 flex items-center justify-center overflow-x-auto overflow-y-auto bg-foreground/40 p-4" role="dialog" aria-modal="true">
          <div className="max-h-[calc(100svh-2rem)] w-full max-w-3xl overflow-x-auto overflow-y-auto bg-background p-6 shadow-xl">
            <div className="flex items-start justify-between gap-4">
              <div><p className="text-xs uppercase tracking-widest text-primary">Reception details</p><h2 className="mt-2 font-serif text-3xl">{selected.referenceKey}</h2></div>
              <Button type="button" variant="ghost" onClick={() => setSelected(null)}>Close</Button>
            </div>
            <dl className="mt-6 grid gap-4 text-sm md:grid-cols-2">
              <div><dt className="text-muted-foreground">Product</dt><dd>{selected.productName} ({selected.productCode || "Validated product"})</dd></div>
              <div><dt className="text-muted-foreground">Artisan / workshop</dt><dd>{selected.artisanName} / {selected.workshopName}</dd></div>
              <div><dt className="text-muted-foreground">Received quantity</dt><dd>{selected.receivedQuantity}</dd></div>
              <div><dt className="text-muted-foreground">Status</dt><dd><StatusBadge status={selected.status} locale={locale} /></dd></div>
            </dl>
            <div className="mt-6 border-t border-border pt-5">
              <h3 className="font-medium">Evidence</h3>
              <div className="mt-3 flex flex-wrap gap-2">{selected.evidence.map((file) => file.url ? <a className="underline" key={file.id} href={file.url} target="_blank" rel="noreferrer">{file.originalFilename || file.mediaType}</a> : <span key={file.id}>{file.originalFilename}</span>)}</div>
              {selected.status === "RECEIVED_PENDING_INSPECTION" && <div className="mt-4 space-y-3"><FileDropzone accept="image/*,application/pdf,video/mp4" label="Choose inspection evidence" file={evidence} onFileChange={(file) => setEvidence(file ?? null)} disabled={saving} />{saving && <UploadProgress label="Uploading evidence and submitting inspection…" />}</div>}
            </div>
            {selected.status === "RECEIVED_PENDING_INSPECTION" ? (
              <div className="mt-6 border-t border-border pt-5">
                <h3 className="font-medium">Inspect batch</h3>
                <div className="mt-4 grid gap-4 md:grid-cols-2">
                  {(["acceptedQuantity", "rejectedQuantity", "quarantinedQuantity", "damagedQuantity"] as const).map((key) => <label className="block text-sm" key={key}><span className="mb-2 block">{key.replace("Quantity", " quantity")}</span><input className="auth-input w-full" type="number" min={0} value={inspection[key]} onChange={(event) => setInspection({ ...inspection, [key]: event.target.value })} /></label>)}
                  <label className="block text-sm md:col-span-2"><span className="mb-2 block">Inspection reason (minimum 3 characters)</span><textarea className="auth-input min-h-24 w-full" required minLength={3} value={inspection.reason} onChange={(event) => setInspection({ ...inspection, reason: event.target.value })} /></label>
                </div>
                <p className="mt-3 text-xs text-muted-foreground">The four outcomes must total {selected.receivedQuantity}. At least one evidence file and a reason of at least 3 characters are required before committing.</p>
                <Button className="mt-4" type="button" disabled={saving || (selected.evidence.length === 0 && !evidence) || inspectionTotal !== selected.receivedQuantity || inspection.reason.trim().length < 3} onClick={() => void inspect()}>{saving ? "Submitting…" : "Submit inspection"}</Button>
              </div>
            ) : selected.inspection ? <div className="mt-6 border-t border-border pt-5"><h3 className="font-medium">Inspection outcome</h3><p className="mt-3 text-sm">Accepted {selected.inspection.acceptedQuantity} · Rejected {selected.inspection.rejectedQuantity} · Quarantined {selected.inspection.quarantinedQuantity} · Damaged {selected.inspection.damagedQuantity}</p><p className="mt-2 text-sm text-muted-foreground">{selected.inspection.reason}</p></div> : null}
          </div>
        </div>
      )}
    </div>
  );
}
