"use client";

import { Monitor, Moon, Sun } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useTheme } from "next-themes";

import { Button } from "@/core/components/ui/button";
import type { Locale } from "@/core/lib/i18n";

const copy = {
  en: { label: "Theme", light: "Light", dark: "Dark", system: "System" },
  fr: { label: "Thème", light: "Clair", dark: "Sombre", system: "Système" },
  ar: { label: "المظهر", light: "فاتح", dark: "داكن", system: "النظام" },
  es: { label: "Tema", light: "Claro", dark: "Oscuro", system: "Sistema" },
} as const;

export function ThemeToggle({ locale }: { locale: Locale }) {
  const text = copy[locale];
  const { theme, setTheme } = useTheme();
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  const active = theme === "dark" ? "dark" : theme === "light" ? "light" : "system";
  const Icon = active === "dark" ? Moon : active === "light" ? Sun : Monitor;

  useEffect(() => {
    const close = (event: MouseEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) setOpen(false);
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", close);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", close);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, []);

  return (
    <div ref={rootRef} className="relative">
      <Button
        type="button"
        variant="ghost"
        className="size-11 min-h-11 px-0"
        aria-label={text.label}
        aria-haspopup="menu"
        aria-expanded={open}
        onClick={() => setOpen((current) => !current)}
      >
        <Icon aria-hidden size={18} strokeWidth={1.5} />
      </Button>
      {open ? (
        <div
          className="absolute end-0 top-[calc(100%+0.5rem)] z-50 min-w-32 border border-border bg-popover p-1 shadow-[var(--shadow-floating)]"
          role="menu"
          aria-label={text.label}
        >
          {([
            ["light", Sun, text.light],
            ["dark", Moon, text.dark],
            ["system", Monitor, text.system],
          ] as const).map(([value, OptionIcon, label]) => (
            <button
              type="button"
              role="menuitemradio"
              aria-checked={active === value}
              className="flex w-full items-center gap-3 px-3 py-2 text-start text-sm hover:bg-muted"
              key={value}
              onClick={() => {
                setTheme(value);
                setOpen(false);
              }}
            >
              <OptionIcon aria-hidden size={15} />
              {label}
            </button>
          ))}
        </div>
      ) : null}
    </div>
  );
}
