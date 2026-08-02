import { notFound } from "next/navigation";
import { ProfileForm } from "@/components/forms/profile-form";
import { isLocale } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";

export default async function ProfilePage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const copy = storeCopy(locale);
  return <section className="max-w-2xl">
    <p className="text-xs uppercase tracking-widest text-primary">{copy.account}</p><h1 className="mt-3 font-serif text-5xl">{copy.profile}</h1>
    <ProfileForm locale={locale} copy={copy} />
  </section>;
}
