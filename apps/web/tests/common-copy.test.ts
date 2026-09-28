import { commonCopy, statusLabel } from "@/core/lib/common-copy";

describe("shared locale copy", () => {
  it("provides translated operational labels for every supported locale", () => {
    for (const locale of ["en", "fr", "ar", "es"] as const) {
      const copy = commonCopy(locale);
      expect(copy.loading).toBeTruthy();
      expect(copy.details).toBeTruthy();
      expect(copy.standardTraditional).toBeTruthy();
      expect(copy.chooseCategory).toBeTruthy();
    }
  });

  it("renders workflow statuses in the selected language", () => {
    expect(statusLabel("PENDING_REVIEW", "fr")).toBe("En attente d’examen");
    expect(statusLabel("PENDING_REVIEW", "ar")).toBe("قيد المراجعة");
    expect(statusLabel("PENDING_REVIEW", "es")).toBe("Pendiente de revisión");
  });
});
