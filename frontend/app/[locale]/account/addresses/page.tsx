import { notFound } from "next/navigation";
import { AddressManager } from "@/components/account/address-manager";
import { isLocale } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";
export default async function AddressesPage({ params }: { params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) notFound(); const copy = storeCopy(locale); return <section><p className="text-xs uppercase tracking-widest text-primary">{copy.account}</p><h1 className="mt-3 font-serif text-5xl">{copy.addresses}</h1><AddressManager copy={copy} locale={locale}/></section>; }
