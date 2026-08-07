import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

import { isLocale, locales, type Locale } from "@/core/lib/i18n";

function preferredLocale(request: NextRequest): Locale {
  const cookieLocale = request.cookies.get("NEXT_LOCALE")?.value;
  if (cookieLocale && isLocale(cookieLocale)) {
    return cookieLocale;
  }

  const accepted = request.headers.get("accept-language")?.split(",").map((part) => part.split(";")[0].trim().slice(0, 2));
  const matched = accepted?.find((value): value is Locale => isLocale(value));
  return matched ?? locales[0];
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const firstSegment = pathname.split("/")[1];

  if (isLocale(firstSegment)) {
    return NextResponse.next();
  }

  const url = request.nextUrl.clone();
  url.pathname = `/${preferredLocale(request)}${pathname === "/" ? "" : pathname}`;
  return NextResponse.redirect(url);
}

export const config = {
  matcher: ["/((?!api|_next|.*\\..*).*)"],
};
