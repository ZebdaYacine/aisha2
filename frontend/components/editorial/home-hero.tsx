import { ArrowDown } from "lucide-react";
import Image from "next/image";

import { ButtonLink } from "@/components/ui/button";
import type { Locale, Messages } from "@/lib/i18n";

export function HomeHero({ locale, messages }: { locale: Locale; messages: Messages }) {
  return (
    <section className="grid min-h-[calc(100svh-6.5rem)] bg-muted lg:grid-cols-[0.9fr_1.1fr]">
      <div className="order-2 flex items-center bg-background px-4 py-16 sm:px-8 lg:order-1 lg:px-12 xl:px-20">
        <div className="max-w-3xl">
          <p className="latin-tracking mb-6 text-xs font-medium uppercase tracking-[0.2em] text-primary">
            {messages.home.heroEyebrow}
          </p>
          <h1 className="text-display-xl max-w-4xl text-balance">{messages.home.heroTitle}</h1>
          <p className="mt-8 max-w-xl text-base leading-7 text-muted-foreground sm:text-lg">
            {messages.home.heroBody}
          </p>
          <div className="mt-10 flex flex-col gap-3 sm:flex-row">
            <ButtonLink href={`/${locale}/#collection`}>{messages.home.explore}</ButtonLink>
            <ButtonLink href={`/${locale}/#artisans`} variant="outline">
              {messages.home.meetArtisans}
            </ButtonLink>
          </div>
        </div>
      </div>
      <div className="relative order-1 min-h-[52svh] overflow-hidden bg-muted lg:order-2 lg:min-h-full">
        <Image
          src="/images/aisha/Algeria art.jpg"
          alt={messages.home.imageAlt}
          fill
          priority
          sizes="(max-width: 1024px) 100vw, 55vw"
          className="object-cover object-[62%_center] lg:object-center"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-foreground/55 via-transparent to-transparent" />
        <div className="absolute inset-x-6 bottom-7 z-10 flex items-end justify-between text-background sm:inset-x-10 sm:bottom-10">
          <p className="max-w-xs font-serif text-2xl sm:text-3xl">{messages.home.heroEyebrow}</p>
          <ArrowDown aria-hidden="true" className="animate-none" strokeWidth={1.25} />
        </div>
      </div>
    </section>
  );
}
