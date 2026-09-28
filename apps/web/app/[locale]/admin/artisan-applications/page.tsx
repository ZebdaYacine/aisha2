import { cookies } from "next/headers";
import { notFound, redirect } from "next/navigation";

import { Container } from "@/core/components/layout/container";
import { isLocale } from "@/core/lib/i18n";
import {
  AdminDashboardShell,
  ArtisanMembershipReview,
  ArtisanReview,
} from "@/features/admin";

export const metadata = { robots: { index: false, follow: false } };
export default async function AdminArtisanApplicationsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const jar = await cookies();
  if (!jar.has("aisha_access") && !jar.has("aisha_refresh"))
    redirect(`/${locale}/login`);
  const copy = {
    eyebrow: locale === "fr" ? "Demandes d’artisans" : locale === "ar" ? "طلبات الحرفيين" : locale === "es" ? "Solicitudes de artesanos" : "Artisan applications",
    title: locale === "fr" ? "File de vérification" : locale === "ar" ? "قائمة المراجعة" : locale === "es" ? "Cola de revisión" : "Review queue",
  };
  return (
    <Container className="py-12">
      <AdminDashboardShell locale={locale} section="artisan-applications">
        <section className="space-y-8">
          <div>
            <p className="text-xs uppercase tracking-widest text-primary">
              {copy.eyebrow}
            </p>
            <h2 className="mt-3 font-serif text-4xl">{copy.title}</h2>
          </div>
          <ArtisanReview locale={locale} />
          <ArtisanMembershipReview locale={locale} />
        </section>
      </AdminDashboardShell>
    </Container>
  );
}
