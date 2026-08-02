import Link from "next/link";import type { ComponentProps } from "react";import type { Locale } from "@/lib/i18n";
export function LocalizedLink({locale,href,...props}:{locale:Locale;href:string}&Omit<ComponentProps<typeof Link>,"href">){return <Link href={`/${locale}${href}`} {...props}/>}
