import { cookies } from "next/headers";
import { notFound, redirect } from "next/navigation";
import { AccountSidebar } from "@/features/account";
import { Container } from "@/core/components/layout/container";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
export default async function AccountLayout({ children, params }: { children: React.ReactNode; params: Promise<{ locale: string }> }) { const { locale } = await params; if (!isLocale(locale)) notFound(); const jar = await cookies(); if (!jar.has("aisha_access") && !jar.has("aisha_refresh")) redirect(`/${locale}/login`); const copy = storeCopy(locale); return <Container className="py-12"><div className="grid gap-10 lg:grid-cols-[15rem_1fr]"><AccountSidebar locale={locale} copy={copy}/><main>{children}</main></div></Container>; }
