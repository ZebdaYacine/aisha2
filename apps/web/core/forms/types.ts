export type FormSubmission<TField extends string> =
  | { ok: true; data?: unknown }
  | {
      ok: false;
      code?:
        | "SERVICE_UNAVAILABLE"
        | "VALIDATION_FAILED"
        | "RATE_LIMITED"
        | "UNKNOWN";
      status?: number;
      fieldErrors?: Partial<Record<TField, string>>;
    };

export type FormSubmitter<
  TValues,
  TField extends string = Extract<keyof TValues, string>,
> = (values: TValues) => Promise<FormSubmission<TField>>;
