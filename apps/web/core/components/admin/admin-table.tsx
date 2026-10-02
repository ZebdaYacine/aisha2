import type { ReactNode } from "react";

import { cn } from "@/core/lib/utils";

/** Shared frame for administrator table tabs. Keeps actions visible while tables scroll horizontally. */
export function AdminTablePanel({
  eyebrow,
  title,
  description,
  action,
  summary,
  children,
  className,
}: {
  eyebrow?: string;
  title: string;
  description?: string;
  action?: ReactNode;
  summary?: ReactNode;
  children: ReactNode;
  className?: string;
}) {
  return (
    <section className={cn("space-y-4", className)}>
      <div className="flex flex-col gap-4 rounded border border-border bg-card p-4 sm:flex-row sm:items-start">
        {action && <div className="order-first shrink-0">{action}</div>}
        <div className="min-w-0 flex-1">
          {eyebrow && <p className="text-xs uppercase tracking-[0.18em] text-primary">{eyebrow}</p>}
          <h2 className="mt-1 font-serif text-2xl sm:text-3xl">{title}</h2>
          {description && <p className="mt-2 max-w-2xl text-sm text-muted-foreground">{description}</p>}
        </div>
        {summary && <span className="ms-auto shrink-0 pt-2 text-sm text-muted-foreground">{summary}</span>}
      </div>
      {children}
    </section>
  );
}

export function AdminTableScroll({ children, className }: { children: ReactNode; className?: string }) {
  return <div className={cn("w-full max-w-full touch-pan-x overflow-x-auto overscroll-x-contain rounded border border-border", className)}>{children}</div>;
}

export const adminTableClass = "w-full min-w-[48rem] text-start text-sm";
export const adminTableHeadClass = "border-b border-border bg-muted/40";
export const adminTableCellClass = "px-4 py-4 align-top";
