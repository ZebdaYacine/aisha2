"use client";

import { FileText, UploadCloud, X } from "lucide-react";
import { useRef, useState } from "react";

type FileDropzoneProps = {
  accept: string;
  label: string;
  disabled?: boolean;
  inputRef?: React.RefObject<HTMLInputElement | null>;
  file?: File | null;
  onFileChange?: (file: File | undefined) => void;
  resetAfterChange?: boolean;
};

export function FileDropzone({
  accept,
  label,
  disabled = false,
  inputRef,
  file,
  onFileChange,
  resetAfterChange = false,
}: FileDropzoneProps) {
  const internalRef = useRef<HTMLInputElement>(null);
  const [dragging, setDragging] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | undefined>(file ?? undefined);
  const displayedFile = file === undefined ? selectedFile : file ?? undefined;
  const input = inputRef ?? internalRef;

  const choose = (next?: File) => {
    if (!next || disabled) return;
    setSelectedFile(resetAfterChange ? undefined : next);
    onFileChange?.(next);
    if (resetAfterChange && input.current) input.current.value = "";
  };

  const clear = () => {
    setSelectedFile(undefined);
    onFileChange?.(undefined);
    if (input.current) input.current.value = "";
  };

  return (
    <div
      className={`rounded-2xl border-2 border-dashed p-4 transition-colors sm:p-5 ${dragging ? "border-primary bg-primary/10" : "border-border bg-muted/20 hover:border-primary/50 hover:bg-primary/5"} ${disabled ? "cursor-not-allowed opacity-60" : ""}`}
      onDragEnter={(event) => { event.preventDefault(); if (!disabled) setDragging(true); }}
      onDragOver={(event) => event.preventDefault()}
      onDragLeave={(event) => { if (event.currentTarget === event.target) setDragging(false); }}
      onDrop={(event) => {
        event.preventDefault();
        setDragging(false);
        choose(event.dataTransfer.files?.[0]);
      }}
    >
      <input
        ref={input}
        className="sr-only"
        type="file"
        accept={accept}
        aria-label={label}
        disabled={disabled}
        onChange={(event) => choose(event.target.files?.[0])}
      />
      <button
        type="button"
        disabled={disabled}
        onClick={() => input.current?.click()}
        className="flex w-full items-center gap-3 text-start outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
      >
        <span className="grid size-11 shrink-0 place-items-center rounded-xl bg-primary/10 text-primary">
          <UploadCloud size={21} aria-hidden="true" />
        </span>
        <span className="min-w-0">
          <span className="block text-sm font-semibold">{label}</span>
          <span className="mt-1 block text-xs text-muted-foreground">Drag and drop or browse from your device</span>
        </span>
      </button>
      {displayedFile && (
        <div className="mt-4 flex items-center gap-3 rounded-xl border border-border bg-background p-3">
          <FileText className="shrink-0 text-primary" size={19} aria-hidden="true" />
          <span className="min-w-0 flex-1">
            <span className="block truncate text-sm font-medium" title={displayedFile.name}>{displayedFile.name}</span>
            <span className="mt-0.5 block text-xs text-muted-foreground">{formatBytes(displayedFile.size)}</span>
          </span>
          {!disabled && (
            <button type="button" onClick={clear} className="rounded-full p-1 text-muted-foreground hover:bg-muted hover:text-foreground" aria-label={`Remove ${displayedFile.name}`}>
              <X size={16} aria-hidden="true" />
            </button>
          )}
        </div>
      )}
    </div>
  );
}

function formatBytes(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
