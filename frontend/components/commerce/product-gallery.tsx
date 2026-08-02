"use client";
import { ChevronLeft, ChevronRight, Expand } from "lucide-react";
import Image from "next/image";
import { useState } from "react";
import { Modal } from "@/components/ui/modal";

export function ProductGallery({ images, alt, previous, next, zoomLabel = alt, closeLabel = alt }: { images: string[]; alt: string; previous: string; next: string; zoomLabel?: string; closeLabel?: string }) {
  const [index, setIndex] = useState(0);
  const [zoom, setZoom] = useState(false);
  const move = (delta: number) => setIndex((current) => (current + delta + images.length) % images.length);
  return <div>
    <div className="gallery-media bg-muted">
      <Image src={images[index]} alt={`${alt} ${index + 1}`} fill priority sizes="(max-width:1024px) 100vw,60vw" className="pointer-events-none object-cover" />
      <button onClick={() => setZoom(true)} className="absolute end-3 top-3 z-10 grid size-11 place-items-center bg-background" aria-label={zoomLabel}><Expand aria-hidden="true" size={18} /></button>
      {images.length > 1 && <><button onClick={() => move(-1)} aria-label={previous} className="absolute start-3 top-1/2 z-10 grid size-11 -translate-y-1/2 place-items-center bg-background"><ChevronLeft className="rtl:-scale-x-100" /></button><button onClick={() => move(1)} aria-label={next} className="absolute end-3 top-1/2 z-10 grid size-11 -translate-y-1/2 place-items-center bg-background"><ChevronRight className="rtl:-scale-x-100" /></button></>}
    </div>
    <div className="mt-3 flex gap-3">{images.map((image, itemIndex) => <button key={itemIndex} onClick={() => setIndex(itemIndex)} className={`gallery-thumb bg-muted ${itemIndex === index ? "ring-2 ring-foreground" : ""}`} aria-label={`${alt} ${itemIndex + 1}`}><Image src={image} alt="" fill sizes="80px" className="pointer-events-none object-cover" /></button>)}</div>
    {zoom && <Modal label={zoomLabel} closeLabel={closeLabel} onClose={() => setZoom(false)} panelClassName="grid place-items-center bg-foreground/95 p-6"><div className="relative h-full w-full max-w-5xl"><Image src={images[index]} alt={alt} fill sizes="100vw" className="object-contain" /></div></Modal>}
  </div>;
}
