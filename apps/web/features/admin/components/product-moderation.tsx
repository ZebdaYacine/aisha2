/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useCallback, useEffect, useState } from "react";
import { Eye } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { IconAction } from "@/core/components/ui/icon-action";
import { Combobox } from "@/core/components/ui/combobox";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import type { Locale } from "@/core/lib/i18n";
import { statusLabel } from "@/core/lib/common-copy";
import {
  AdminTablePanel,
  AdminTableScroll,
  adminTableCellClass,
  adminTableClass,
  adminTableHeadClass,
} from "@/core/components/admin/admin-table";

type ProductStatus =
  | "PENDING_REVIEW"
  | "CHANGES_REQUESTED"
  | "APPROVED"
  | "ACTIVE"
  | "SUSPENDED"
  | "ARCHIVED";
type Item = {
  submissionId: string;
  productId: string;
  productCode?: string;
  productName: string;
  productStatus: ProductStatus;
  priceMinor: number;
  currency: string;
  artisanName: string;
  workshopName: string;
  media: MediaItem[];
};
type MediaItem = {
  id: string;
  mediaKind: string;
  originalFilename: string;
  mediaType: string;
  sizeBytes: number;
  altText: string;
  visibility: string;
  url?: string;
};
const statusOptions: ProductStatus[] = [
  "PENDING_REVIEW",
  "CHANGES_REQUESTED",
  "APPROVED",
  "ACTIVE",
  "SUSPENDED",
  "ARCHIVED",
];

export function ProductModeration({ locale = "en" }: { locale?: Locale }) {
  const [items, setItems] = useState<Item[]>([]);
  const [selected, setSelected] = useState<Item | null>(null);
  useEscapeKey(() => setSelected(null), selected !== null);
  const [reason, setReason] = useState("");
  const [status, setStatus] = useState<ProductStatus>("PENDING_REVIEW");
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 10;
  const pageCount = Math.max(1, Math.ceil(total / pageSize));

  const load = useCallback(
    async (nextPage = page, nextStatus = status) => {
      setLoading(true);
      try {
        const response = await fetch(
          `/api/admin/product-submissions?status=${nextStatus}&page=${nextPage}&pageSize=${pageSize}`,
          { cache: "no-store" },
        );
        if (!response.ok) throw new Error();
        const data = (await response.json()) as {
          items?: Item[];
          page?: number;
          total?: number;
        };
        setItems(data.items ?? []);
        setPage(data.page ?? nextPage);
        setTotal(data.total ?? 0);
      } catch {
        toast.error("Unable to load product submissions");
      } finally {
        setLoading(false);
      }
    },
    [page, status],
  );

  useEffect(() => {
    void load();
  }, [load]);

  const selectStatus = (nextStatus: ProductStatus) => {
    setStatus(nextStatus);
    setPage(1);
    void load(1, nextStatus);
  };

  const decide = async (action: string) => {
    if (!selected) return;
    if (
      ["REQUEST_CHANGES", "REJECT", "SUSPEND"].includes(action) &&
      !reason.trim()
    ) {
      toast.error("A reason is required");
      return;
    }
    const decisionID = ["ACTIVATE", "SUSPEND", "ARCHIVE"].includes(action)
      ? selected.productId
      : selected.submissionId;
    const response = await fetch(
      `/api/admin/product-submissions/${decisionID}/decision`,
      {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ action, reason }),
      },
    );
    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      toast.error(body.error?.message ?? "Unable to record decision");
      return;
    }
    toast.success("Decision recorded");
    setSelected(null);
    setReason("");
    await load();
  };

  const actions =
    selected?.productStatus === "PENDING_REVIEW"
      ? ([
          ["APPROVE", "Approve", "primary"],
          ["REQUEST_CHANGES", "Request changes", "outline"],
          ["REJECT", "Reject", "destructive"],
        ] as const)
      : selected?.productStatus === "APPROVED"
        ? ([
            ["ACTIVATE", "Activate", "primary"],
            ["ARCHIVE", "Archive", "destructive"],
          ] as const)
        : selected?.productStatus === "ACTIVE"
          ? ([
              ["SUSPEND", "Suspend", "destructive"],
              ["ARCHIVE", "Archive", "outline"],
            ] as const)
          : selected?.productStatus === "SUSPENDED"
            ? ([["ARCHIVE", "Archive", "destructive"]] as const)
            : ([] as const);

  return (
    <>
      <AdminTablePanel
        eyebrow="Product moderation"
        title="Submission queue"
        action={<label className="block text-sm"><span className="sr-only">{locale === "fr" ? "Statut de la file" : locale === "ar" ? "حالة القائمة" : locale === "es" ? "Estado de la cola" : "Queue status"}</span><Combobox className="min-w-56" options={statusOptions.map((option) => ({ value: option, label: statusLabel(option, locale) }))} value={status} onChange={(next) => selectStatus(next as ProductStatus)} ariaLabel={locale === "fr" ? "Statut de la file" : locale === "ar" ? "حالة القائمة" : locale === "es" ? "Estado de la cola" : "Queue status"} /></label>}
        summary={`${total} product(s)`}
      >
      {loading ? (
        <p
          className="border border-border p-5 text-sm text-muted-foreground"
          role="status"
        >
          Loading product submissions…
        </p>
      ) : (
        <>
          <AdminTableScroll>
            <table className={adminTableClass}>
              <thead className={adminTableHeadClass}>
                <tr>
                  <th className="px-4 py-3 font-medium">Product</th>
                  <th className="px-4 py-3 font-medium">Artisan</th>
                  <th className="px-4 py-3 font-medium">Status</th>
                  <th className="px-4 py-3 font-medium">Price</th>
                  <th className="px-4 py-3 font-medium">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {items.map((item) => (
                  <tr className="transition-colors hover:bg-muted/30" key={item.submissionId || item.productId}>
                    <td className={adminTableCellClass}>
                      <span>{item.productName || "Unnamed product"}</span>
                      {item.productCode && <p className="text-xs text-muted-foreground">{item.productCode}</p>}
                    </td>
                    <td className={adminTableCellClass}>
                      {item.artisanName}
                      <p className="text-xs text-muted-foreground">
                        {item.workshopName}
                      </p>
                    </td>
                    <td className={adminTableCellClass}>
                      <StatusBadge status={item.productStatus} locale={locale} />
                    </td>
                    <td className={adminTableCellClass}>
                      {item.priceMinor} {item.currency}
                    </td>
                    <td className={adminTableCellClass}>
                      <IconAction
                        icon={<Eye size={17} />}
                        label="Details"
                        onClick={() => {
                          setSelected(item);
                          setReason("");
                        }}
                      />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {items.length === 0 && (
              <p className="p-6 text-sm text-muted-foreground">
                No products in this queue.
              </p>
            )}
          </AdminTableScroll>
          <div className="mt-4 flex flex-col-reverse items-stretch gap-3 sm:flex-row sm:items-center sm:justify-between">
            <Button
              type="button"
              variant="outline"
              className="w-full sm:w-auto"
              disabled={loading || page <= 1}
              onClick={() => void load(page - 1)}
            >
              Previous
            </Button>
            <span className="text-sm text-muted-foreground">
              Page {page} of {pageCount}
            </span>
            <Button
              type="button"
              variant="outline"
              className="w-full sm:w-auto"
              disabled={loading || page >= pageCount}
              onClick={() => void load(page + 1)}
            >
              Next
            </Button>
          </div>
        </>
      )}
      </AdminTablePanel>
      {selected && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center overflow-x-auto bg-foreground/50 p-4"
          role="dialog"
          aria-modal="true"
          aria-labelledby="product-detail-title"
        >
          <div className="max-h-[calc(100svh-2rem)] w-full max-w-xl overflow-x-auto overflow-y-auto border border-border bg-background p-6">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-widest text-primary">
                  Product moderation
                </p>
                <h2
                  id="product-detail-title"
                  className="mt-2 font-serif text-3xl"
                >
                  {selected.productName || "Unnamed product"}
                </h2>
              </div>
              <Button
                type="button"
                variant="ghost"
                onClick={() => setSelected(null)}
              >
                Close
              </Button>
            </div>
            <dl className="mt-6 space-y-3 text-sm">
              <div>
                <dt className="text-muted-foreground">Artisan</dt>
                <dd>{selected.artisanName}</dd>
              </div>
              {selected.productCode && <div>
                <dt className="text-muted-foreground">Warehouse product code</dt>
                <dd>{selected.productCode}</dd>
              </div>}
              <div>
                <dt className="text-muted-foreground">Workshop</dt>
                <dd>{selected.workshopName}</dd>
              </div>
              <div>
                <dt className="text-muted-foreground">Status</dt>
                <dd><StatusBadge status={selected.productStatus} locale={locale} /></dd>
              </div>
              <div>
                <dt className="text-muted-foreground">Price</dt>
                <dd>
                  {selected.priceMinor} {selected.currency}
                </dd>
              </div>
            </dl>
            <section className="mt-6 border-t border-border pt-5" aria-labelledby="product-media-title">
              <h3 id="product-media-title" className="font-medium">Product media</h3>
              {selected.media?.length ? (
                <div className="mt-4 grid gap-4 sm:grid-cols-2">
                  {selected.media.map((item) => <ModerationMediaCard key={item.id} item={item} />)}
                </div>
              ) : (
                <p className="mt-3 text-sm text-muted-foreground">No product media uploaded.</p>
              )}
            </section>
            {actions.length > 0 && (
              <>
                <label className="mt-6 block text-sm">
                  <span className="mb-2 block">Decision reason</span>
                  <textarea
                    className="auth-input min-h-24"
                    value={reason}
                    onChange={(event) => setReason(event.target.value)}
                    placeholder="Required for changes, rejection, and suspension"
                  />
                </label>
                <div className="mt-5 flex flex-wrap gap-2">
                  {actions.map(([action, label, variant]) => (
                    <Button
                      key={action}
                      type="button"
                      variant={variant}
                      onClick={() => void decide(action)}
                    >
                      {label}
                    </Button>
                  ))}
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </>
  );
}

function ModerationMediaCard({ item }: { item: MediaItem }) {
  const isImage = item.mediaType.startsWith("image/");
  const isVideo = item.mediaType.startsWith("video/");
  return <article className="border border-border p-3">
    {item.url && isImage && <a href={item.url} target="_blank" rel="noreferrer" className="block aspect-video overflow-hidden bg-muted"><span className="block h-full w-full bg-cover bg-center" style={{ backgroundImage: `url("${item.url}")` }} /></a>}
    {item.url && isVideo && <video className="aspect-video w-full bg-muted object-cover" controls preload="metadata" src={item.url} />}
    {!isImage && !isVideo && <div className="flex aspect-video items-center justify-center bg-muted text-xs text-muted-foreground">{item.mediaType}</div>}
    <p className="mt-3 truncate text-sm font-medium">{item.originalFilename || item.mediaKind}</p>
    <p className="mt-1 text-xs text-muted-foreground">{item.mediaKind} · {item.mediaType} · {formatBytes(item.sizeBytes)}</p>
    {item.url && <a className="mt-2 inline-block text-xs underline" href={item.url} target="_blank" rel="noreferrer">Open media</a>}
  </article>;
}

function formatBytes(value: number) {
  if (value < 1024) return `${value} B`;
  if (value < 1024 * 1024) return `${Math.round(value / 1024)} KB`;
  return `${(value / (1024 * 1024)).toFixed(1)} MB`;
}
