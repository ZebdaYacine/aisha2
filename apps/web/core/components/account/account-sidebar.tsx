"use client";
import { useRouter } from "next/navigation";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { LocalizedLink } from "@/core/components/shared/localized-link";
import { Button } from "@/core/components/ui/button";
import { useAuth } from "@/features/auth/viewmodel/auth-context";
export function AccountSidebar({ locale, copy }: { locale: Locale; copy: StoreCopy }) { const router = useRouter(); const { logout } = useAuth(); const links = [[copy.overview, "/account"], [copy.profile, "/account/profile"], [copy.addresses, "/account/addresses"], [copy.orders, "/account/orders"], [copy.wishlist, "/account/wishlist"]]; const signOut = async () => { await logout(); router.replace(`/${locale}/login`); router.refresh(); }; return <aside><nav aria-label={copy.account} className="flex gap-2 overflow-x-auto border-b border-border pb-4 lg:flex-col lg:border-b-0 lg:border-e lg:pe-8">{links.map(([label, href]) => <LocalizedLink className="min-h-11 whitespace-nowrap px-3 py-3 hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" locale={locale} href={href} key={label}>{label}</LocalizedLink>)}<Button type="button" variant="ghost" onClick={() => void signOut()}>{copy.logout ?? "Logout"}</Button></nav></aside>; }
