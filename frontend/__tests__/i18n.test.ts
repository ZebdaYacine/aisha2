import { dictionary, direction, isLocale } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";

describe("international foundation", () => {
  it("supports the four documented locales and Arabic RTL", () => {
    expect(["en", "fr", "ar", "es"].every(isLocale)).toBe(true);
    expect(direction("ar")).toBe("rtl");
    expect(direction("fr")).toBe("ltr");
  });

  it("provides complete storefront foundation copy for every locale", () => {
    for (const locale of ["en", "fr", "ar", "es"] as const) {
      const messages = dictionary(locale);
      expect(messages.home.heroTitle).toBeTruthy();
      expect(messages.navigation.skipToContent).toBeTruthy();
      expect(messages.footer.story).toBeTruthy();
    }
  });

  it("does not fall back to English for public-storefront interaction copy", () => {
    const keys = ["featuredBody", "madeToOrderBody", "newsletter", "subscribe", "featuredSort", "priceLow", "priceHigh", "apply", "tryAgain", "searchHint", "previous", "next", "trustOrigin", "trustQuality", "artisanStory", "quality", "origin", "delivery"] as const;
    const english = storeCopy("en");
    for (const locale of ["fr", "ar", "es"] as const) {
      const copy = storeCopy(locale);
      for (const key of keys) {
        expect(copy[key]).toBeTruthy();
        expect(copy[key]).not.toBe(english[key]);
      }
    }
  });
});
