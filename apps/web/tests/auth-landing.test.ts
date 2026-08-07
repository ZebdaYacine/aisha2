import { landingPathForUser } from "@/features/auth/types";

describe("authentication landing paths", () => {
  it("sends administrators to the admin workspace", () => {
    expect(landingPathForUser({ id: "1", email: "admin@example.test", displayName: "Admin", roles: ["administrator"] }, "en")).toBe("/en/admin/artisan-applications");
  });

  it("sends artisans to the artisan workspace", () => {
    expect(landingPathForUser({ id: "2", email: "artisan@example.test", displayName: "Artisan", roles: ["artisan"] }, "fr")).toBe("/fr/artisan");
  });

  it("keeps customers in the account workspace", () => {
    expect(landingPathForUser({ id: "3", email: "customer@example.test", displayName: "Customer", roles: ["customer"] }, "ar")).toBe("/ar/account");
  });
});
