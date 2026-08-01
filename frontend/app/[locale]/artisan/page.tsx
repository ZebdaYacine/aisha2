import{notFound}from"next/navigation";import{ArtisanWorkspace}from"@/components/artisan/workspace";import{isLocale}from"@/lib/i18n";
export default async function ArtisanPage({params}:{params:Promise<{locale:string}>}){const{locale}=await params;if(!isLocale(locale))notFound();return <ArtisanWorkspace locale={locale}/>}
