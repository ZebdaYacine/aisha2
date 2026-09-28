import { cookies } from "next/headers";
import { notFound, redirect } from "next/navigation";
import { AccountSidebar } from "@/features/account";
import { Container } from "@/core/components/layout/container";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
export default async function AccountLayout({ children, params }: { children: React.ReactNode; params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) notFound(); const jar = await cookies(); if (!jar.has("aisha_access") && !jar.has("aisha_refresh")) redirect(`/${locale}/login`); const copy = storeCopy(locale); return <Container className="py-8 sm:py-12"><div className="grid min-w-0 gap-8 lg:grid-cols-[15rem_minmax(0,1fr)] lg:gap-10"><AccountSidebar locale={locale} copy={copy} mode="customer"/><main className="min-w-0">{children}</main></div></Container>; }
