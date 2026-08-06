import { NextRequest, NextResponse } from "next/server";

import { catalogue } from "@/features/catalogue/api";
import { isLocale } from "@/core/lib/i18n";

export async function GET(request: NextRequest) {
  const locale = request.nextUrl.searchParams.get("locale") ?? "en";
  if (!isLocale(locale)) {
    return NextResponse.json({ code: "VALIDATION_ERROR", message: "Unsupported locale." }, { status: 400 });
  }

  try {
    return NextResponse.json(await catalogue(locale));
  } catch (error) {
    console.error("Catalogue proxy failed", error);
    return NextResponse.json(
      { code: "CATALOGUE_UNAVAILABLE", message: "The catalogue is temporarily unavailable." },
      { status: 503 },
    );
  }
}
