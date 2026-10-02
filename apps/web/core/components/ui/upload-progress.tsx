"use client";

export function UploadProgress({ label, detail }: { label: string; detail?: string }) {
  return (
    <div className="mt-3 rounded-md border border-primary/20 bg-primary/5 px-3 py-2" role="status" aria-live="polite">
      <div className="flex items-center justify-between gap-3 text-xs">
        <span className="font-medium text-foreground">{label}</span>
      </div>
      <div className="mt-2 h-1 overflow-hidden rounded-full bg-muted" aria-hidden="true">
        <span className="upload-progress-bar block h-full w-2/5 rounded-full bg-primary" />
      </div>
      {detail && <p className="mt-2 text-[0.6875rem] text-muted-foreground">{detail}</p>}
    </div>
  );
}
