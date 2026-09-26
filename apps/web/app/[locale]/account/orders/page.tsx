import { notFound } from "next/navigation";
import { LiveOrders } from "@/features/order";
import { isLocale } from "@/core/lib/i18n";
import { storeCopy } from "@/core/lib/store-copy";
export default async function OrdersPage({params}:{params:Promise<{locale:string}>}){const{locale}=await params;if(!isLocale(locale))notFound();const copy=storeCopy(locale);return <><p className="text-xs uppercase tracking-widest text-primary">{copy.account}</p><h1 className="mt-3 font-serif text-5xl">{copy.orderHistory}</h1><LiveOrders locale={locale} copy={copy}/></>}
