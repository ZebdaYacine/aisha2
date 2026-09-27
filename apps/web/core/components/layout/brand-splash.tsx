"use client";

import { useEffect, useState } from "react";

import type { Locale } from "@/core/lib/i18n";

const copy = {
  en: { eyebrow: "A story made by hand", lines: ["Rooted in place", "Shaped by memory", "Carried forward"], loading: "Opening the collection" },
  fr: { eyebrow: "Une histoire faite à la main", lines: ["Ancrée dans un lieu", "Façonnée par la mémoire", "Transmise"], loading: "Ouverture de la collection" },
  ar: { eyebrow: "حكاية صنعت باليد", lines: ["متجذرة في المكان", "تشكلها الذاكرة", "نحملها إلى الأمام"], loading: "جارٍ فتح المجموعة" },
  es: { eyebrow: "Una historia hecha a mano", lines: ["Arraigada en un lugar", "Formada por la memoria", "Llevada hacia adelante"], loading: "Abriendo la colección" },
} as const;

const SPLASH_SEEN_KEY = "aisha:splash-seen-at";

export function BrandSplash({ locale }: { locale: Locale }) {
  const text = copy[locale];
  const [visible, setVisible] = useState(false);
  const [progress, setProgress] = useState(0);

  useEffect(() => {
    let frame = 0;
    let finish: number | undefined;
    let show = true;

    try {
      const seenAt = Number(window.localStorage.getItem(SPLASH_SEEN_KEY));
      show = !seenAt;
    } catch {
      // Continue with the splash when storage is unavailable.
    }

    if (!show) return;

    const started = performance.now();
    const duration = 2200;
    const animate = (now: number) => {
      const next = Math.min(100, ((now - started) / duration) * 100);
      setProgress(next);
      if (next < 100) frame = window.requestAnimationFrame(animate);
      else {
        finish = window.setTimeout(() => {
          try { window.localStorage.setItem(SPLASH_SEEN_KEY, String(Date.now())); } catch { /* no-op */ }
          setVisible(false);
        }, 260);
      }
    };
    const begin = window.setTimeout(() => {
      setVisible(true);
      frame = window.requestAnimationFrame(animate);
    }, 0);
    return () => {
      window.clearTimeout(begin);
      window.cancelAnimationFrame(frame);
      if (finish) window.clearTimeout(finish);
    };
  }, []);

  if (!visible) return null;
  const line = text.lines[Math.min(text.lines.length - 1, Math.floor(progress / 34))];
  return (
    <div className="brand-splash" role="status" aria-live="polite" aria-label={text.loading}>
      <div className="brand-splash__grain" aria-hidden="true" />
      <div className="brand-splash__content">
        <p className="brand-splash__eyebrow">{text.eyebrow}</p>
        <div className="brand-splash__mark" aria-hidden="true"><span>A</span><span>I</span><span>S</span><span>H</span><span>A</span></div>
        <p className="brand-splash__story" key={line}>{line}</p>
        <div className="brand-splash__progress" role="progressbar" aria-valuemin={0} aria-valuemax={100} aria-valuenow={Math.round(progress)}>
          <span style={{ transform: `scaleX(${progress / 100})` }} />
        </div>
        <div className="brand-splash__meta"><span>{text.loading}</span><span>{Math.round(progress)}%</span></div>
      </div>
    </div>
  );
}
