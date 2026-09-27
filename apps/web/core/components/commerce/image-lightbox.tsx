"use client";

import Image from "next/image";
import { useState, type PointerEvent } from "react";

import { Modal } from "@/core/components/ui/modal";
import { cn } from "@/core/lib/utils";

type ImageLightboxProps = {
  src: string;
  alt: string;
  thumbnailClassName?: string;
  sizes?: string;
  closeLabel?: string;
};

/** A compact image trigger with a centered, cursor-position zoom viewer. */
export function ImageLightbox({
  src,
  alt,
  thumbnailClassName,
  sizes = "(max-width: 640px) 100vw, 25vw",
  closeLabel = "Close image",
}: ImageLightboxProps) {
  const [open, setOpen] = useState(false);
  const [zoomed, setZoomed] = useState(false);
  const [origin, setOrigin] = useState("50% 50%");

  const trackPointer = (event: PointerEvent<HTMLDivElement>) => {
    const bounds = event.currentTarget.getBoundingClientRect();
    const x = bounds.width ? Math.min(100, Math.max(0, ((event.clientX - bounds.left) / bounds.width) * 100)) : 50;
    const y = bounds.height ? Math.min(100, Math.max(0, ((event.clientY - bounds.top) / bounds.height) * 100)) : 50;
    setOrigin(`${x}% ${y}%`);
  };

  const show = () => {
    setOrigin("50% 50%");
    setZoomed(false);
    setOpen(true);
  };

  return (
    <>
      <button
        type="button"
        onClick={show}
        aria-label={alt}
        className={cn("relative block overflow-hidden bg-muted", thumbnailClassName)}
      >
        <Image src={src} alt={alt} fill sizes={sizes} className="object-cover" />
      </button>
      {open && (
        <Modal label={alt} closeLabel={closeLabel} onClose={() => setOpen(false)} panelClassName="flex min-h-full w-full items-center justify-center bg-foreground/95 p-5 sm:p-10">
          <div className="flex max-h-full max-w-full items-center justify-center overflow-auto">
            <div
              className="relative max-h-[calc(100svh-5rem)] max-w-[92vw] cursor-zoom-in overflow-hidden bg-foreground"
              onPointerMove={trackPointer}
              onClick={() => setZoomed((current) => !current)}
              style={{ cursor: zoomed ? "zoom-out" : "zoom-in" }}
            >
              <Image
                src={src}
                alt={alt}
                width={1600}
                height={1200}
                sizes="92vw"
                className="max-h-[calc(100svh-5rem)] w-auto max-w-[92vw] object-contain transition-transform duration-150 ease-out"
                style={{ transformOrigin: origin, transform: zoomed ? "scale(2)" : "scale(1)" }}
              />
            </div>
          </div>
        </Modal>
      )}
    </>
  );
}
