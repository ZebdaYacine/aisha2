import { cookies } from "next/headers";
import { notFound, redirect } from "next/navigation";

import { isLocale } from "@/core/lib/i18n";

export const metadata = { robots: { index: false, follow: false } };

export default async function AdminDashboardPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const jar = await cookies();
  if (!jar.has("aisha_access") && !jar.has("aisha_refresh"))
    redirect(`/${locale}/login`);
  redirect(`/${locale}/admin/users`);
}
