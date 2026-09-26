/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { hasCapability } from "@/features/auth/types";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

type Application = { id: string; publicDisplayName: string; workshopName: string; wilaya: string; status: string; membershipStatus: string };
type Media = { id: string; mediaKind: string; originalFilename: string; mediaType: string; sizeBytes: number; url?: string };

export function ArtisanMembershipReview() {
  const auth = useOptionalAuth();
  const canManageMembership = auth ? hasCapability(auth.user, "admin.artisan_applications.write") : true;
  const [applications, setApplications] = useState<Application[]>([]);
  const [selected, setSelected] = useState<Application | null>(null);
  const [selectedMedia, setSelectedMedia] = useState<Media[]>([]);
  const [mediaLoading, setMediaLoading] = useState(false);
  const [reasons, setReasons] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);
  useEscapeKey(() => setSelected(null), selected !== null);

  const load = useCallback(async () => {
    const response = await fetch("/api/admin/artisan-applications?status=APPROVED&pageSize=100");
    if (!response.ok) {
      toast.error("Unable to load artisan membership review");
      setLoading(false);
      return;
    }
    setApplications((await response.json()).items ?? []);
    setLoading(false);
  }, []);

  useEffect(() => {
    if (canManageMembership) void load();
  }, [canManageMembership, load]);

  useEffect(() => {
    if (!selected) {
      setSelectedMedia([]);
      return;
    }
    let cancelled = false;
    setMediaLoading(true);
    void fetch(`/api/admin/artisan-applications/${selected.id}/media`, { cache: "no-store" })
      .then(async (response) => {
        if (!response.ok) throw new Error("Unable to load artisan media");
        return (await response.json()) as Media[];
      })
      .then((items) => { if (!cancelled) setSelectedMedia(items); })
      .catch((error) => { if (!cancelled) toast.error(error instanceof Error ? error.message : "Unable to load artisan media"); })
      .finally(() => { if (!cancelled) setMediaLoading(false); });
    return () => { cancelled = true; };
  }, [selected]);

  if (!canManageMembership) return null;

  const changeMembership = async (id: string, status: string) => {
    const response = await fetch(`/api/admin/artisan-memberships/${id}/status`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ status, reason: reasons[id] ?? "" }),
    });
    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      toast.error(body.error?.message ?? "Unable to change membership status");
      return;
    }
    toast.success("Membership status updated");
    setSelected(null);
    await load();
  };

  if (loading) return <section className="mt-10 border-t border-border pt-8" aria-busy="true">Loading membership review…</section>;

  return <section className="mt-10 space-y-8 border-t border-border pt-8" aria-labelledby="membership-review-title">
    <h2 id="membership-review-title" className="font-serif text-3xl">Membership management</h2>
    <div className="space-y-4">
      <h3 className="font-medium">Approved artisans</h3>
      <div className="overflow-x-auto border border-border">
        <table className="w-full min-w-[44rem] text-start text-sm">
          <thead className="border-b border-border bg-muted/40"><tr><th className="px-4 py-3 font-medium">Artisan</th><th className="px-4 py-3 font-medium">Shop</th><th className="px-4 py-3 font-medium">Status</th><th className="px-4 py-3 text-end font-medium">Actions</th></tr></thead>
          <tbody className="divide-y divide-border">{applications.map((item) => <tr key={item.id}><td className="px-4 py-4">{item.publicDisplayName}<span className="block text-xs text-muted-foreground">{item.wilaya}</span></td><td className="px-4 py-4">{item.workshopName || "—"}</td><td className="px-4 py-4"><StatusBadge status={item.membershipStatus} /></td><td className="px-4 py-4"><div className="flex flex-wrap justify-end gap-2"><Button type="button" variant="outline" disabled={item.membershipStatus === "ACTIVE"} onClick={() => void changeMembership(item.id, "ACTIVE")}>Activate</Button><Button type="button" variant="destructive" disabled={item.membershipStatus !== "ACTIVE"} onClick={() => void changeMembership(item.id, "SUSPENDED")}>Suspend</Button><Button type="button" variant="ghost" onClick={() => setSelected(item)}>Details</Button></div></td></tr>)}</tbody>
        </table>
        {applications.length === 0 && <p className="p-6 text-sm text-muted-foreground">No approved artisans found.</p>}
      </div>
    </div>
    {selected && <div className="fixed inset-0 z-50 flex items-center justify-center overflow-x-auto overflow-y-auto bg-foreground/50 p-4" role="dialog" aria-modal="true" aria-labelledby="approved-artisan-details"><div className="max-h-[calc(100svh-2rem)] w-full max-w-2xl overflow-x-auto overflow-y-auto border border-border bg-background p-6"><div className="flex items-start justify-between gap-4"><div><p className="text-xs uppercase tracking-widest text-primary">Approved artisan</p><h3 id="approved-artisan-details" className="mt-2 font-serif text-2xl">{selected.publicDisplayName}</h3></div><Button type="button" variant="ghost" onClick={() => setSelected(null)}>Close</Button></div><dl className="mt-6 grid gap-4 text-sm sm:grid-cols-2"><div><dt className="text-muted-foreground">Shop</dt><dd className="mt-1">{selected.workshopName || "—"}</dd></div><div><dt className="text-muted-foreground">Membership status</dt><dd className="mt-1"><StatusBadge status={selected.membershipStatus} /></dd></div></dl><section className="mt-6 border-t border-border pt-5" aria-labelledby="artisan-user-media-title"><h4 id="artisan-user-media-title" className="font-medium">User media</h4>{mediaLoading ? <p className="mt-3 text-sm text-muted-foreground">Loading media…</p> : selectedMedia.length === 0 ? <p className="mt-3 text-sm text-muted-foreground">No uploaded media.</p> : <div className="mt-3 grid gap-4 sm:grid-cols-2">{selectedMedia.map((item) => <article key={item.id} className="overflow-hidden border border-border">{item.url && item.mediaType.startsWith("image/") ? <img src={item.url} alt={item.originalFilename || item.mediaKind} className="aspect-video w-full object-cover" /> : item.url && item.mediaType.startsWith("video/") ? <video src={item.url} controls className="aspect-video w-full bg-muted object-cover" /> : <div className="flex aspect-video items-center justify-center bg-muted text-xs text-muted-foreground">{item.mediaType}</div>}<div className="p-3 text-xs"><p className="truncate font-medium">{item.originalFilename || item.mediaKind}</p><p className="mt-1 text-muted-foreground">{item.mediaType} · {formatBytes(item.sizeBytes)}</p>{item.url && <a className="mt-2 inline-block underline" href={item.url} target="_blank" rel="noreferrer">Open media</a>}</div></article>)}</div>}</section><label className="mt-6 block text-sm"><span className="mb-2 block">Reason for suspension</span><input className="auth-input" value={reasons[selected.id] ?? ""} onChange={(event) => setReasons((current) => ({ ...current, [selected.id]: event.target.value }))} /></label><div className="mt-5 flex flex-wrap justify-end gap-2"><Button type="button" variant="outline" disabled={selected.membershipStatus === "ACTIVE"} onClick={() => void changeMembership(selected.id, "ACTIVE")}>Activate</Button><Button type="button" variant="destructive" disabled={selected.membershipStatus !== "ACTIVE"} onClick={() => void changeMembership(selected.id, "SUSPENDED")}>Suspend</Button></div></div></div>}
  </section>;
}

function formatBytes(value: number) { if (value < 1024) return `${value} B`; if (value < 1024 * 1024) return `${Math.round(value / 1024)} KB`; return `${(value / (1024 * 1024)).toFixed(1)} MB`; }
