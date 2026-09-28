"use client";

import { X } from "lucide-react";
import { useEffect, useRef, type ReactNode } from "react";

import { cn } from "@/core/lib/utils";

export function Modal({ children, closeLabel, label, onClose, panelClassName }: { children: ReactNode; closeLabel: string; label: string; onClose: () => void; panelClassName?: string }) {
  const panelRef = useRef<HTMLDivElement>(null);
  const onCloseRef = useRef(onClose);

  useEffect(() => {
    onCloseRef.current = onClose;
  }, [onClose]);

  useEffect(() => {
    const previousFocus = document.activeElement as HTMLElement | null;
    const panel = panelRef.current;
    panel?.focus();
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape" && !event.defaultPrevented) onCloseRef.current();
      if (event.key !== "Tab" || !panel) return;
      const focusable = [...panel.querySelectorAll<HTMLElement>('a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])')];
      if (!focusable.length) return;
      const first = focusable[0];
      const last = focusable[focusable.length - 1];
      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first.focus();
      }
    };
    document.body.style.overflow = "hidden";
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.body.style.overflow = "";
      document.removeEventListener("keydown", onKeyDown);
      previousFocus?.focus();
    };
  }, []);

  return <div className="fixed inset-0 z-[100] overscroll-contain overflow-x-hidden overflow-y-auto bg-foreground/60" onMouseDown={onClose}>
    <div ref={panelRef} role="dialog" aria-modal="true" aria-label={label} tabIndex={-1} className={cn("relative min-h-full w-full max-w-full overflow-x-auto overflow-y-auto bg-background outline-none [overflow-wrap:anywhere]", panelClassName)} onMouseDown={(event) => event.stopPropagation()}>
      <button type="button" onClick={onClose} aria-label={closeLabel} className="absolute end-4 top-4 z-20 grid size-11 place-items-center bg-background text-foreground"><X aria-hidden="true" size={20} /></button>
      {children}
    </div>
  </div>;
}
