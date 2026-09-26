import { notFound } from "next/navigation";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
import { WishlistList } from "@/features/wishlist/components/wishlist-list";

export default async function AccountWishlistPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const copy = storeCopy(locale);
  return <><h1 className="font-serif text-5xl">{copy.wishlist}</h1><WishlistList locale={locale} copy={copy} /></>;
}
