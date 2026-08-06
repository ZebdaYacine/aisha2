import { NextRequest } from "next/server";import{authenticatedBackend,proxyResponse}from"@/core/lib/server/backend";
export async function GET(){return proxyResponse(await authenticatedBackend("/artisan-applications/me"));}
export async function POST(r:NextRequest){return proxyResponse(await authenticatedBackend("/artisan-applications",{method:"POST",body:JSON.stringify(await r.json())}));}
