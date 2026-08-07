import { HomeHero } from "../components/home-hero";
import { HomeSections } from "../components/home-sections";
import { dictionary, type Locale } from "@/core/lib/i18n";
import { catalogue } from "@/features/catalogue/api";

export async function HomeView({ locale }: { locale: Locale }) {
  const messages = dictionary(locale);
  const data = await catalogue(locale);

  return (
    <>
      <HomeHero locale={locale} messages={messages} />
      <HomeSections locale={locale} messages={messages} {...data} />
    </>
  );
}
