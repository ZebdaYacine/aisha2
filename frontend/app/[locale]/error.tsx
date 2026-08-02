"use client";

import { usePathname } from "next/navigation";

import { ErrorState } from "@/components/feedback/error-state";
import { Container } from "@/components/layout/container";
import { Button } from "@/components/ui/button";
import { isLocale } from "@/lib/i18n";
import { storeCopy } from "@/lib/store-copy";

export default function ErrorPage({ reset }: { error: Error & { digest?: string }; reset: () => void }) {
  const segment = usePathname().split("/")[1];
  const locale = isLocale(segment) ? segment : "en";
  const copy = storeCopy(locale);
  return <Container><ErrorState title={copy.errorTitle} body={copy.errorBody} action={<Button onClick={reset}>{copy.retry}</Button>} /></Container>;
}
