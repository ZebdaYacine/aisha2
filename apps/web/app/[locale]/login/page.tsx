import Image from "next/image";
import { notFound } from "next/navigation";
import { AuthForm } from "@/core/components/forms/auth-form";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
export default async function LoginPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!isLocale(locale)) notFound();
  const copy = storeCopy(locale);
  return (
    <div className="grid min-h-[calc(100svh-6rem)] lg:grid-cols-2">
      <section className="flex items-center justify-center px-5 py-16">
        <div className="w-full max-w-md">
          <p className="text-xs uppercase tracking-widest text-primary">
            AISHA
          </p>
          <h1 className="mt-4 font-serif text-5xl">{copy.login}</h1>
          <p className="mt-4 text-muted-foreground">{copy.welcome}</p>
          <div className="mt-9">
            <AuthForm mode="login" locale={locale} copy={copy} />
          </div>
        </div>
      </section>
      <div className="relative hidden lg:block">
        <Image
          src="/images/aisha/Algeria art.jpg"
          alt=""
          fill
          sizes="50vw"
          className="object-cover"
        />
      </div>
    </div>
  );
}
