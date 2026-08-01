import { MapPin } from "lucide-react";
import { notFound } from "next/navigation";
import { Button } from "@/components/ui/button";
import { isLocale } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";

export default async function AddressesPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const copy = storeCopy(locale);
  return <section><p className="text-xs uppercase tracking-widest text-primary">{copy.account}</p><h1 className="mt-3 font-serif text-5xl">{copy.addresses}</h1><article className="mt-10 max-w-xl border border-border bg-card p-6"><MapPin aria-hidden className="text-primary" /><h2 className="mt-5 font-serif text-2xl">{copy.address}</h2><address className="mt-3 not-italic text-muted-foreground">Amina Benali<br />Alger, Algeria</address><Button type="button" variant="outline" className="mt-6">{copy.addresses}</Button></article></section>;
}
