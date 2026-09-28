/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { hasCapability } from "@/features/auth/types";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";
import type { Locale } from "@/core/lib/i18n";
import {
  AdminTablePanel,
  AdminTableScroll,
  adminTableCellClass,
  adminTableClass,
  adminTableHeadClass,
} from "./admin-table";

type Item = {
  id: string;
  publicDisplayName: string;
  workshopName: string;
  wilaya: string;
  status: string;
  reviewReason?: string;
};
type DocumentItem = {
  id: string;
  documentType: string;
  originalFilename: string;
  mediaType: string;
  sizeBytes: number;
  url?: string;
};
type MediaItem = {
  id: string;
  mediaKind: string;
  originalFilename: string;
  mediaType: string;
  sizeBytes: number;
  url?: string;
};

export function ArtisanReview({ locale = "en" }: { locale?: Locale }) {
  const auth = useOptionalAuth();
  const canDecide = auth ? hasCapability(auth.user, "admin.artisan_applications.write") : true;
  const [items, setItems] = useState<Item[]>([]);
  const [selected, setSelected] = useState<Item | null>(null);
  const [documents, setDocuments] = useState<DocumentItem[]>([]);
  const [media, setMedia] = useState<MediaItem[]>([]);
  const [filesLoading, setFilesLoading] = useState(false);
  const [reason, setReason] = useState("");
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 10;
  useEscapeKey(() => setSelected(null), selected !== null);
  const load = useCallback(
    async (nextPage = page) => {
      setLoading(true);
      const response = await fetch(
        `/api/admin/artisan-applications?status=SUBMITTED&page=${nextPage}&pageSize=${pageSize}`,
      );
      if (response.ok) {
        const data = (await response.json()) as {
          items?: Item[];
          page?: number;
          total?: number;
        };
        setItems(data.items ?? []);
        setPage(data.page ?? nextPage);
        setTotal(data.total ?? 0);
      } else toast.error("Unable to load applications");
      setLoading(false);
    },
    [page],
  );
  useEffect(() => {
    void load();
  }, [load]);
  const decide = async (decision: string) => {
    if (!selected) return;
    const response = await fetch(
      `/api/admin/artisan-applications/${selected.id}/${decision}`,
      {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ reason }),
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
  const openDetails = async (item: Item) => {
    setSelected(item);
    setDocuments([]);
    setMedia([]);
    setFilesLoading(true);
    try {
      const [documentsResponse, mediaResponse] = await Promise.all([
        fetch(`/api/admin/artisan-applications/${item.id}/documents`),
        fetch(`/api/admin/artisan-applications/${item.id}/media`),
      ]);
      if (!documentsResponse.ok || !mediaResponse.ok) throw new Error();
      setDocuments((await documentsResponse.json()) as DocumentItem[]);
      setMedia((await mediaResponse.json()) as MediaItem[]);
    } catch {
      toast.error("Unable to load application media");
    } finally {
      setFilesLoading(false);
    }
  };
  if (loading)
    return (
      <p
        className="border border-border p-5 text-sm text-muted-foreground"
        role="status"
      >
        Loading applications…
      </p>
    );
  return (
    <>
      <AdminTablePanel
        eyebrow="Review queue"
        title="Submitted applications"
        summary={`${total} application(s)`}
      >
      <AdminTableScroll>
        <table className={adminTableClass}>
          <thead className={adminTableHeadClass}>
            <tr>
              <th className="px-4 py-3 font-medium">Applicant</th>
              <th className="px-4 py-3 font-medium">Workshop</th>
              <th className="px-4 py-3 font-medium">Wilaya</th>
              <th className="px-4 py-3 font-medium">Status</th>
              <th className="px-4 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {items.map((item) => (
              <tr className="transition-colors hover:bg-muted/30" key={item.id}>
                <td className={adminTableCellClass}>{item.publicDisplayName}</td>
                <td className={adminTableCellClass}>{item.workshopName}</td>
                <td className={adminTableCellClass}>{item.wilaya}</td>
                <td className={adminTableCellClass}>
                  <StatusBadge status={item.status} locale={locale} />
                </td>
                <td className={adminTableCellClass}>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => {
                      setReason("");
                      void openDetails(item);
                    }}
                  >
                    Details
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {items.length === 0 && (
          <p className="p-6 text-sm text-muted-foreground">
            No submitted applications.
          </p>
        )}
      </AdminTableScroll>
      <div className="mt-4 flex items-center justify-between gap-3">
        <Button
          type="button"
          variant="outline"
          disabled={loading || page <= 1}
          onClick={() => void load(page - 1)}
        >
          Previous
        </Button>
        <span className="text-sm text-muted-foreground">
          Page {page} of {Math.max(1, Math.ceil(total / pageSize))}
        </span>
        <Button
          type="button"
          variant="outline"
          disabled={loading || page >= Math.max(1, Math.ceil(total / pageSize))}
          onClick={() => void load(page + 1)}
        >
          Next
        </Button>
      </div>
      </AdminTablePanel>
      {selected && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center overflow-x-auto bg-foreground/50 p-4"
          role="dialog"
          aria-modal="true"
          aria-labelledby="artisan-application-detail"
        >
          <div className="max-h-[calc(100svh-2rem)] w-full max-w-xl overflow-x-auto overflow-y-auto border border-border bg-background p-6">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-widest text-primary">
                  Application details
                </p>
                <h2
                  id="artisan-application-detail"
                  className="mt-2 font-serif text-3xl"
                >
                  {selected.publicDisplayName}
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
                <dt className="text-muted-foreground">Workshop</dt>
                <dd>{selected.workshopName}</dd>
              </div>
              <div>
                <dt className="text-muted-foreground">Wilaya</dt>
                <dd>{selected.wilaya}</dd>
              </div>
              <div>
                <dt className="text-muted-foreground">Status</dt>
                <dd><StatusBadge status={selected.status} locale={locale} /></dd>
              </div>
              {selected.reviewReason && (
                <div>
                  <dt className="text-muted-foreground">Previous reason</dt>
                  <dd>{selected.reviewReason}</dd>
                </div>
              )}
            </dl>
            <section className="mt-6 border-t border-border pt-5" aria-labelledby="artisan-application-media">
              <h3 id="artisan-application-media" className="font-medium">Uploaded media</h3>
              {filesLoading && <p className="mt-3 text-sm text-muted-foreground" role="status">Loading media…</p>}
              {!filesLoading && documents.length === 0 && media.length === 0 && (
                <p className="mt-3 text-sm text-muted-foreground">No uploaded media.</p>
              )}
              {!filesLoading && (documents.length > 0 || media.length > 0) && (
                <div className="mt-4 grid gap-4 sm:grid-cols-2">
                  {documents.map((item) => (
                    <MediaFileCard key={`document-${item.id}`} label={item.documentType} item={item} />
                  ))}
                  {media.map((item) => (
                    <MediaFileCard key={`media-${item.id}`} label={item.mediaKind} item={item} />
                  ))}
                </div>
              )}
            </section>
            {canDecide && <>
              <label className="mt-6 block text-sm">
                <span className="mb-2 block">Decision reason</span>
                <textarea
                  className="auth-input min-h-24"
                  value={reason}
                  onChange={(event) => setReason(event.target.value)}
                />
              </label>
              <div className="mt-5 flex flex-wrap gap-2">
                <Button type="button" onClick={() => void decide("approve")}>
                  Approve
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => void decide("request-changes")}
                >
                  Request changes
                </Button>
                <Button
                  type="button"
                  variant="destructive"
                  onClick={() => void decide("reject")}
                >
                  Reject
                </Button>
              </div>
            </>}
          </div>
        </div>
      )}
    </>
  );
}

function MediaFileCard({
  label,
  item,
}: {
  label: string;
  item: DocumentItem | MediaItem;
}) {
  const isImage = item.mediaType.startsWith("image/");
  const isVideo = item.mediaType.startsWith("video/");
  return (
    <article className="border border-border p-3">
      {item.url && isImage && (
        <a href={item.url} target="_blank" rel="noreferrer" className="block aspect-video overflow-hidden bg-muted">
          <span className="block h-full w-full bg-cover bg-center" style={{ backgroundImage: `url("${item.url}")` }} />
        </a>
      )}
      {item.url && isVideo && <video className="aspect-video w-full bg-muted object-cover" controls preload="metadata" src={item.url} />}
      {!isImage && !isVideo && <div className="flex aspect-video items-center justify-center bg-muted text-xs text-muted-foreground">{item.mediaType}</div>}
      <p className="mt-3 truncate text-sm font-medium">{item.originalFilename || label}</p>
      <p className="mt-1 text-xs text-muted-foreground">{label} · {item.mediaType} · {formatBytes(item.sizeBytes)}</p>
      {item.url && <a className="mt-2 inline-block text-xs underline" href={item.url} target="_blank" rel="noreferrer">Open media</a>}
    </article>
  );
}

function formatBytes(value: number) {
  if (value < 1024) return `${value} B`;
  if (value < 1024 * 1024) return `${Math.round(value / 1024)} KB`;
  return `${(value / (1024 * 1024)).toFixed(1)} MB`;
}
