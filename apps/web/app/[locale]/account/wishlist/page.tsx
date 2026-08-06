import { notFound } from "next/navigation";
import { EmptyState } from "@/core/components/feedback/empty-state";
import { ButtonLink } from "@/core/components/ui/button";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export default async function AccountWishlistPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const copy = storeCopy(locale);
  return <><h1 className="font-serif text-5xl">{copy.wishlist}</h1><EmptyState title={copy.wishlist} body={copy.noResults} action={<ButtonLink href={`/${locale}/products`}>{copy.continueShopping}</ButtonLink>} /></>;
}
