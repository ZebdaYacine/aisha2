"use client";

import { Check, ChevronDown, X } from "lucide-react";
import { useEffect, useId, useMemo, useRef, useState } from "react";

import { cn } from "@/core/lib/utils";

export type ComboboxOption = { value: string; label: string };

type ComboboxProps = {
  options: ComboboxOption[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  emptyMessage?: string;
  ariaLabel?: string;
  disabled?: boolean;
  allowCustom?: boolean;
  displayValue?: (option: ComboboxOption) => string;
  className?: string;
  id?: string;
};

export function Combobox({
  options,
  value,
  onChange,
  placeholder = "Select an option",
  emptyMessage = "No matches found",
  ariaLabel,
  disabled,
  allowCustom = false,
  displayValue,
  className,
  id: inputId,
}: ComboboxProps) {
  const generatedId = useId();
  const id = inputId ?? generatedId;
  const rootRef = useRef<HTMLDivElement>(null);
  const selected = options.find((option) => option.value === value);
  const selectedDisplay = selected ? displayValue?.(selected) ?? selected.label : value;
  const [query, setQuery] = useState(selectedDisplay);
  const [open, setOpen] = useState(false);
  const filtered = useMemo(() => {
    const normalized = query.trim().toLocaleLowerCase();
    if (!normalized || normalized === selected?.label.toLocaleLowerCase() || normalized === selectedDisplay.toLocaleLowerCase()) {
      return options;
    }
    return options.filter((option) =>
      option.label.toLocaleLowerCase().includes(normalized),
    );
  }, [options, query, selected?.label, selectedDisplay]);

  useEffect(() => {
    const close = (event: MouseEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, []);

  const choose = (option: ComboboxOption) => {
    onChange(option.value);
    setQuery(displayValue?.(option) ?? option.label);
    setOpen(false);
  };

  return (
    <div ref={rootRef} className={cn("relative min-w-0", open && "z-50", className)}>
      <div className="relative">
        <input
          id={id}
          role="combobox"
          aria-label={ariaLabel}
          aria-controls={`${id}-options`}
          aria-expanded={open}
          aria-autocomplete="list"
          className="auth-input min-w-0 pe-10"
          value={query}
          placeholder={placeholder}
          disabled={disabled}
          onFocus={() => {
            setQuery(selected?.label ?? value);
            setOpen(true);
          }}
          onChange={(event) => {
            setQuery(event.target.value);
            setOpen(true);
            if (allowCustom) onChange(event.target.value);
          }}
          onKeyDown={(event) => {
            if (event.key === "Escape") setOpen(false);
            if (event.key === "Enter" && filtered[0]) {
              event.preventDefault();
              choose(filtered[0]);
            }
          }}
          onBlur={() => {
            window.setTimeout(() => {
              if (!allowCustom && !options.some((option) => option.value === value)) {
                setQuery(selected?.label ?? "");
              }
              setOpen(false);
            }, 120);
          }}
        />
        {value && allowCustom ? (
          <button
            type="button"
            className="absolute inset-y-0 end-8 flex w-8 items-center justify-center text-muted-foreground hover:text-foreground"
            aria-label="Clear selection"
            onMouseDown={(event) => event.preventDefault()}
            onClick={() => {
              onChange("");
              setQuery("");
            }}
          >
            <X size={15} />
          </button>
        ) : null}
        <ChevronDown
          aria-hidden
          className="pointer-events-none absolute end-3 top-1/2 -translate-y-1/2 text-muted-foreground"
          size={16}
        />
      </div>
      {open && !disabled ? (
        <ul
          id={`${id}-options`}
          role="listbox"
          className="absolute inset-x-0 top-[calc(100%+0.35rem)] z-[60] max-h-[min(16rem,40vh)] touch-pan-y overscroll-contain overflow-y-auto rounded-md border border-border bg-background p-1 shadow-lg"
        >
          {filtered.length ? (
            filtered.map((option) => (
              <li key={option.value} role="option" aria-selected={option.value === value}>
                <button
                  type="button"
                  className="flex min-h-11 w-full items-center justify-between gap-3 px-3 py-2 text-start text-sm hover:bg-muted"
                  onMouseDown={(event) => event.preventDefault()}
                  onClick={() => choose(option)}
                >
                  <span>{option.label}</span>
                  {option.value === value ? <Check size={15} /> : null}
                </button>
              </li>
            ))
          ) : (
            <li className="px-3 py-2 text-sm text-muted-foreground">{emptyMessage}</li>
          )}
        </ul>
      ) : null}
    </div>
  );
}

type MultiComboboxProps = Omit<ComboboxProps, "value" | "onChange" | "allowCustom"> & {
  value: string[];
  onChange: (value: string[]) => void;
};

export function MultiCombobox({
  options,
  value,
  onChange,
  placeholder = "Select options",
  emptyMessage,
  ariaLabel,
  disabled,
  className,
}: MultiComboboxProps) {
  const id = useId();
  const rootRef = useRef<HTMLDivElement>(null);
  const [query, setQuery] = useState("");
  const [open, setOpen] = useState(false);
  const filtered = options.filter((option) =>
    option.label.toLocaleLowerCase().includes(query.trim().toLocaleLowerCase()),
  );
  const labels = value
    .map((item) => options.find((option) => option.value === item)?.label ?? item)
    .filter(Boolean);

  useEffect(() => {
    const close = (event: MouseEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, []);

  return (
    <div ref={rootRef} className={cn("relative", className)}>
      <div className="flex min-h-12 flex-wrap items-center gap-1 border border-input bg-background px-2 py-1 focus-within:border-foreground">
        {labels.map((label, index) => (
          <button
            type="button"
            key={`${value[index]}-${label}`}
            className="inline-flex items-center gap-1 bg-muted px-2 py-1 text-xs"
            onClick={() => onChange(value.filter((_, itemIndex) => itemIndex !== index))}
            aria-label={`Remove ${label}`}
          >
            {label}<X size={12} />
          </button>
        ))}
        <input
          id={id}
          role="combobox"
          aria-label={ariaLabel}
          aria-controls={`${id}-options`}
          aria-expanded={open}
          aria-autocomplete="list"
          className="min-w-24 flex-1 bg-transparent px-2 py-2 text-sm outline-none"
          value={query}
          placeholder={labels.length ? "Add another" : placeholder}
          disabled={disabled}
          onFocus={() => setOpen(true)}
          onChange={(event) => {
            setQuery(event.target.value);
            setOpen(true);
          }}
          onKeyDown={(event) => {
            if (event.key === "Escape") setOpen(false);
            if (event.key === "Backspace" && !query && value.length) {
              onChange(value.slice(0, -1));
            }
          }}
        />
        <ChevronDown aria-hidden className="me-1 text-muted-foreground" size={16} />
      </div>
      {open && !disabled ? (
        <ul id={`${id}-options`} role="listbox" className="absolute inset-x-0 top-[calc(100%+0.35rem)] z-50 max-h-60 overflow-y-auto border border-border bg-background p-1 shadow-lg">
          {filtered.length ? filtered.map((option) => {
            const active = value.includes(option.value);
            return <li key={option.value} role="option" aria-selected={active}>
              <button type="button" className="flex w-full items-center justify-between gap-3 px-3 py-2 text-start text-sm hover:bg-muted" onMouseDown={(event) => event.preventDefault()} onClick={() => { onChange(active ? value.filter((item) => item !== option.value) : [...value, option.value]); setQuery(""); }}>
                {option.label}{active ? <Check size={15} /> : null}
              </button>
            </li>;
          }) : <li className="px-3 py-2 text-sm text-muted-foreground">{emptyMessage ?? "No matches found"}</li>}
        </ul>
      ) : null}
    </div>
  );
}
