import { notFound } from "next/navigation";
import { LiveOrderDetail } from "@/features/order";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
export default async function OrderPage({params}:{params:Promise<{locale:string;id:string}>}){const{locale,id}=await params;if(!isLocale(locale))notFound();return <LiveOrderDetail locale={locale} copy={storeCopy(locale)} id={id}/>;}
