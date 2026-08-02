import { LockKeyhole } from "lucide-react";
import { notFound } from "next/navigation";
import { OrderSummary } from "@/components/commerce/order-summary";
import { CheckoutForm } from "@/components/forms/checkout-form";
import { Container } from "@/components/layout/container";
import { Breadcrumbs } from "@/components/shared/breadcrumbs";
import { isLocale } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";

export default async function CheckoutPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params; if (!isLocale(locale)) notFound(); const copy = storeCopy(locale);
  return <Container className="py-10"><Breadcrumbs locale={locale} items={[{ label: copy.cart, href: "/cart" }, { label: copy.checkout }]} /><div className="mt-8 flex items-center justify-between border-b border-border pb-6"><h1 className="font-serif text-4xl">{copy.checkout}</h1><p className="flex items-center gap-2 text-sm"><LockKeyhole size={16} />{copy.secureCheckout}</p></div><div className="mt-10 grid gap-12 lg:grid-cols-[1fr_24rem]"><CheckoutForm copy={copy} locale={locale} /><OrderSummary locale={locale} copy={copy} checkoutAction={false} /></div></Container>;
}
