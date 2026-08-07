import{notFound}from"next/navigation";import{ArtisanWorkspace}from"@/features/artisan";import{isLocale}from"@/core/lib/i18n";
export default async function ArtisanPage({params}:{params:Promise<{locale:string}>}){const{locale}=await params;if(!isLocale(locale))notFound();return <ArtisanWorkspace locale={locale}/>}
