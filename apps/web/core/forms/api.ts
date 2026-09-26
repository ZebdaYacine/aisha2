import type { FormSubmission } from "./types";
type APIError = { error?: { code?: string; fields?: Record<string, string> } };
export async function submitJSON<T extends object>(
  url: string,
  values: T,
  method = "POST",
): Promise<FormSubmission<Extract<keyof T, string>>> {
  try {
    const response = await fetch(url, {
      method,
      headers: { "content-type": "application/json" },
      body: JSON.stringify(values),
    });
    if (response.ok) {
      const data = await response.json().catch(() => undefined);
      return data === undefined ? { ok: true } : { ok: true, data };
    }
    const data = (await response.json().catch(() => ({}))) as APIError;
    const result: FormSubmission<Extract<keyof T, string>> = {
      ok: false,
      status: response.status,
      code:
        response.status >= 500
          ? "SERVICE_UNAVAILABLE"
          : response.status === 429 || data.error?.code === "RATE_LIMITED"
            ? "RATE_LIMITED"
            : response.status === 422 || response.status === 400
              ? "VALIDATION_FAILED"
              : "UNKNOWN",
    };
    if (data.error?.fields) result.fieldErrors = data.error.fields;
    return result;
  } catch {
    return { ok: false, code: "SERVICE_UNAVAILABLE" };
  }
}
