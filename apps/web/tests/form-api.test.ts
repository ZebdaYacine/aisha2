import { submitJSON } from "@/core/forms/api";

const response = (status: number, body: unknown) =>
  ({
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  }) as Response;

describe("form API mapping", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("maps backend field errors and status", async () => {
    globalThis.fetch = jest
      .fn()
      .mockResolvedValue(
        response(400, {
          error: {
            code: "VALIDATION_ERROR",
            fields: { email: "INVALID_EMAIL" },
          },
        }),
      ) as unknown as typeof fetch;
    await expect(
      submitJSON("/api/auth/login", { email: "bad" }),
    ).resolves.toEqual({
      ok: false,
      status: 400,
      code: "VALIDATION_FAILED",
      fieldErrors: { email: "INVALID_EMAIL" },
    });
  });

  it("maps rate limiting without empty field errors", async () => {
    globalThis.fetch = jest
      .fn()
      .mockResolvedValue(
        response(429, { error: { code: "RATE_LIMITED" } }),
      ) as unknown as typeof fetch;
    await expect(
      submitJSON("/api/auth/login", { email: "a@b.com" }),
    ).resolves.toEqual({ ok: false, status: 429, code: "RATE_LIMITED" });
  });

  it("keeps successful authentication data for the state store", async () => {
    const user = { id: "user-1", email: "a@b.com", displayName: "Amina" };
    globalThis.fetch = jest
      .fn()
      .mockResolvedValue(response(200, { user })) as unknown as typeof fetch;
    await expect(
      submitJSON("/api/auth/login", { email: user.email }),
    ).resolves.toEqual({ ok: true, data: { user } });
  });
});
