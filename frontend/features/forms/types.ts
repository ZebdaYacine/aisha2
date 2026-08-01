export type FormSubmission<TField extends string> =
  | { ok: true }
  | { ok: false; code?: "SERVICE_UNAVAILABLE" | "VALIDATION_FAILED" | "UNKNOWN"; fieldErrors?: Partial<Record<TField, string>> };

export type FormSubmitter<TValues, TField extends string = Extract<keyof TValues, string>> = (values: TValues) => Promise<FormSubmission<TField>>;
