"use client";

import { usePathname } from "next/navigation";

import { ErrorState } from "@/core/components/feedback/error-state";
import { Container } from "@/core/components/layout/container";
import { Button } from "@/core/components/ui/button";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";

export default function ErrorPage({ reset }: { error: Error & { digest?: string }; reset: () => void }) {
  const segment = usePathname().split("/")[1];
  const locale = isLocale(segment) ? segment : "en";
  const copy = storeCopy(locale);
  return <Container><ErrorState title={copy.errorTitle} body={copy.errorBody} action={<Button onClick={reset}>{copy.retry}</Button>} /></Container>;
}
