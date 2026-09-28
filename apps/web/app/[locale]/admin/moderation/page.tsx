import { cookies } from "next/headers";
import { notFound, redirect } from "next/navigation";

import { Container } from "@/core/components/layout/container";
import { isLocale } from "@/core/lib/i18n";
import { AdminDashboardShell, ProductModeration } from "@/features/admin";

export const metadata = { robots: { index: false, follow: false } };
export default async function AdminModerationPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const jar = await cookies();
  if (!jar.has("aisha_access") && !jar.has("aisha_refresh"))
    redirect(`/${locale}/login`);
  return (
    <Container className="py-12">
      <AdminDashboardShell locale={locale} section="moderation">
        <section>
            <p className="text-xs uppercase tracking-widest text-primary">
              {locale === "fr" ? "Modération des produits" : locale === "ar" ? "مراجعة المنتجات" : locale === "es" ? "Moderación de productos" : "Product moderation"}
            </p>
            <h2 className="mt-3 font-serif text-4xl">{locale === "fr" ? "File des soumissions" : locale === "ar" ? "قائمة الطلبات" : locale === "es" ? "Cola de envíos" : "Submission queue"}</h2>
          <div className="mt-8">
            <ProductModeration locale={locale} />
          </div>
        </section>
      </AdminDashboardShell>
    </Container>
  );
}
