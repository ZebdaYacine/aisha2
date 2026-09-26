import Image from "next/image";
import { notFound } from "next/navigation";

import { AuthForm } from "@/features/auth";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export default async function RegisterPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();

  const copy = storeCopy(locale);
  const isArabic = locale === "ar";

  return (
    <div className="grid min-h-[calc(100svh-6rem)] lg:grid-cols-2" dir="ltr">
      <div className={`relative hidden lg:block ${isArabic ? "lg:order-1" : "lg:order-2"}`}>
        <Image
          src="/images/aisha/Algeria art.jpg"
          alt=""
          fill
          sizes="50vw"
          className="object-cover"
        />
      </div>
      <section
        dir={isArabic ? "rtl" : "ltr"}
        className={`flex items-center justify-center px-5 py-16 ${isArabic ? "lg:order-2" : "lg:order-1"}`}
      >
        <div className="w-full max-w-md">
          <p className="text-xs uppercase tracking-widest text-primary">AISHA</p>
          <h1 className="mt-4 font-serif text-5xl">{copy.register}</h1>
          <div className="mt-9">
            <AuthForm mode="register" locale={locale} copy={copy} />
          </div>
        </div>
      </section>
    </div>
  );
}
