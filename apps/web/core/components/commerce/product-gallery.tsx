"use client";
import { ChevronLeft, ChevronRight, Expand } from "lucide-react";
import Image from "next/image";
import { useState, type PointerEvent } from "react";
import { Modal } from "@/core/components/ui/modal";

export function ProductGallery({ images, alt, previous, next, zoomLabel = alt, closeLabel = alt }: { images: string[]; alt: string; previous: string; next: string; zoomLabel?: string; closeLabel?: string }) {
  const [index, setIndex] = useState(0);
  const [zoom, setZoom] = useState(false);
  const [zoomed, setZoomed] = useState(false);
  const [origin, setOrigin] = useState("50% 50%");
  const move = (delta: number) => setIndex((current) => (current + delta + images.length) % images.length);
  const trackPointer = (event: PointerEvent<HTMLDivElement>) => {
    const bounds = event.currentTarget.getBoundingClientRect();
    const x = bounds.width ? Math.min(100, Math.max(0, ((event.clientX - bounds.left) / bounds.width) * 100)) : 50;
    const y = bounds.height ? Math.min(100, Math.max(0, ((event.clientY - bounds.top) / bounds.height) * 100)) : 50;
    setOrigin(`${x}% ${y}%`);
  };
  return <div>
    <div className="gallery-media bg-muted">
      <Image src={images[index]} alt={`${alt} ${index + 1}`} fill priority sizes="(max-width:1024px) 100vw,60vw" className="pointer-events-none object-cover" />
      <button onClick={() => setZoom(true)} className="absolute end-3 top-3 z-10 grid size-11 place-items-center bg-background" aria-label={zoomLabel}><Expand aria-hidden="true" size={18} /></button>
      {images.length > 1 && <><button onClick={() => move(-1)} aria-label={previous} className="absolute start-3 top-1/2 z-10 grid size-11 -translate-y-1/2 place-items-center bg-background"><ChevronLeft className="rtl:-scale-x-100" /></button><button onClick={() => move(1)} aria-label={next} className="absolute end-3 top-1/2 z-10 grid size-11 -translate-y-1/2 place-items-center bg-background"><ChevronRight className="rtl:-scale-x-100" /></button></>}
    </div>
    <div className="mt-3 flex gap-3">{images.map((image, itemIndex) => <button key={itemIndex} onClick={() => setIndex(itemIndex)} className={`gallery-thumb bg-muted ${itemIndex === index ? "ring-2 ring-foreground" : ""}`} aria-label={`${alt} ${itemIndex + 1}`}><Image src={image} alt="" fill sizes="80px" className="pointer-events-none object-cover" /></button>)}</div>
    {zoom && <Modal label={zoomLabel} closeLabel={closeLabel} onClose={() => { setZoom(false); setZoomed(false); }} panelClassName="flex min-h-full w-full items-center justify-center bg-foreground/95 p-5 sm:p-10"><div className="flex max-h-full max-w-full items-center justify-center overflow-auto"><div className="relative max-h-[calc(100svh-5rem)] max-w-[92vw] cursor-zoom-in overflow-hidden" onPointerMove={trackPointer} onClick={() => setZoomed((current) => !current)} style={{ cursor: zoomed ? "zoom-out" : "zoom-in" }}><Image src={images[index]} alt={alt} width={1600} height={1200} sizes="92vw" className="max-h-[calc(100svh-5rem)] w-auto max-w-[92vw] object-contain transition-transform duration-150 ease-out" style={{ transformOrigin: origin, transform: zoomed ? "scale(2)" : "scale(1)" }} /></div></div></Modal>}
  </div>;
}
