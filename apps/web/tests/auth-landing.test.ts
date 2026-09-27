import { landingPathForUser } from "@/features/auth/types";

describe("authentication landing paths", () => {
  it("sends administrators to the admin workspace", () => {
    expect(
      landingPathForUser(
        {
          id: "1",
          email: "admin@example.test",
          displayName: "Admin",
          customerEnabled: true,
          artisanStatus: "NOT_STARTED",
          artisanEnabled: false,
          capabilities: ["admin.audit.read"],
        },
        "en",
      ),
    ).toBe("/en/admin");
  });

  it("sends artisans to the artisan workspace", () => {
    expect(
      landingPathForUser(
        {
          id: "2",
          email: "artisan@example.test",
          displayName: "Artisan",
          customerEnabled: true,
          artisanStatus: "ACTIVE",
          artisanEnabled: true,
          capabilities: ["artisan.account.read"],
        },
        "fr",
      ),
    ).toBe("/fr/artisan");
  });

  it("keeps customers in the account workspace", () => {
    expect(
      landingPathForUser(
        {
          id: "3",
          email: "customer@example.test",
          displayName: "Customer",
          customerEnabled: true,
          artisanStatus: "NOT_STARTED",
          artisanEnabled: false,
          capabilities: ["customer.account.read"],
        },
        "ar",
      ),
    ).toBe("/ar/account");
  });

  it("keeps a suspended artisan in the customer account", () => {
    expect(
      landingPathForUser(
        {
          id: "4",
          email: "suspended@example.test",
          displayName: "Suspended",
          customerEnabled: true,
          artisanStatus: "SUSPENDED",
          artisanEnabled: false,
          capabilities: ["customer.account.read"],
        },
        "en",
      ),
    ).toBe("/en/account");
  });
});
