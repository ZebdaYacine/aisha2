import { dictionary, direction, isLocale } from "@/lib/i18n";

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
});
