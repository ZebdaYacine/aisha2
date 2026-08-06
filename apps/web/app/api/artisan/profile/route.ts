import { NextRequest } from "next/server";import{authenticatedBackend,proxyResponse}from"@/core/lib/server/backend";
export async function PATCH(r:NextRequest){return proxyResponse(await authenticatedBackend("/artisan/profile",{method:"PATCH",body:JSON.stringify(await r.json())}));}
