import { cookies } from "next/headers";
import { notFound, redirect } from "next/navigation";

import { Container } from "@/core/components/layout/container";
import { isLocale } from "@/core/lib/i18n";
import { AdminDashboardShell, AdminOrderOperations } from "@/features/admin";

export const metadata = { robots: { index: false, follow: false } };

export default async function AdminOrdersPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const jar = await cookies();
  if (!jar.has("aisha_access") && !jar.has("aisha_refresh")) redirect(`/${locale}/login`);
  return <Container className="py-12"><AdminDashboardShell locale={locale} section="orders"><AdminOrderOperations locale={locale} /></AdminDashboardShell></Container>;
}
