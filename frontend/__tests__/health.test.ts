/** @jest-environment node */

import { GET } from "@/app/api/health/route";

describe("frontend health route", () => {
  it("returns a successful health response", async () => {
    const response = GET();
    expect(response.status).toBe(200);
    await expect(response.json()).resolves.toEqual({ status: "ok" });
  });
});
