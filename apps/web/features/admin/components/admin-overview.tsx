"use client";

import { useEffect, useState } from "react";
import { ArrowUpRight, ClipboardCheck, FileClock, PackageCheck, ShoppingBag } from "lucide-react";

import { LocalizedLink } from "@/core/components/shared/localized-link";
import type { Locale } from "@/core/lib/i18n";

type PageResponse = { total?: number };

type Copy = {
  eyebrow: string;
  title: string;
  description: string;
  openOrders: string;
  pendingReviews: string;
  users: string;
  auditEvents: string;
  attention: string;
  attentionBody: string;
  recentActivity: string;
  recentActivityBody: string;
  applications: string;
  applicationsBody: string;
  moderation: string;
  moderationBody: string;
  orders: string;
  ordersBody: string;
  audit: string;
  viewQueue: string;
  viewAll: string;
  loading: string;
  activityWorkspace: string;
  activityUsers: string;
  activityOrders: string;
  now: string;
  today: string;
};

const copy: Record<Locale, Copy> = {
  en: { eyebrow: "ADMINISTRATION / CONTROL CENTER", title: "Good morning", description: "A clear view of what needs your attention across the marketplace.", openOrders: "Open orders", pendingReviews: "Pending reviews", users: "User accounts", auditEvents: "Audit events", attention: "Needs your attention", attentionBody: "Queues that need a decision or follow-up.", recentActivity: "Recent activity", recentActivityBody: "Latest changes across your workspace.", applications: "Artisan applications awaiting review", applicationsBody: "Review evidence and make a reasoned decision.", moderation: "Products need moderation", moderationBody: "Approve, reject, or request changes.", orders: "Orders in progress", ordersBody: "Monitor fulfilment and customer context.", audit: "Audit log", viewQueue: "View queue", viewAll: "View all", loading: "Loading…", activityWorkspace: "Administration workspace", activityUsers: "User accounts", activityOrders: "Orders", now: "Now", today: "Today" },
  fr: { eyebrow: "ADMINISTRATION / CENTRE DE CONTRÔLE", title: "Bonjour", description: "Une vue claire des sujets qui nécessitent votre attention.", openOrders: "Commandes ouvertes", pendingReviews: "Révisions en attente", users: "Comptes utilisateurs", auditEvents: "Événements d’audit", attention: "À votre attention", attentionBody: "Files nécessitant une décision ou un suivi.", recentActivity: "Activité récente", recentActivityBody: "Derniers changements dans votre espace.", applications: "Candidatures d’artisans à examiner", applicationsBody: "Vérifiez les preuves et prenez une décision motivée.", moderation: "Produits à modérer", moderationBody: "Approuver, rejeter ou demander des modifications.", orders: "Commandes en cours", ordersBody: "Suivez l’exécution et le contexte client.", audit: "Journal d’audit", viewQueue: "Voir la file", viewAll: "Tout voir", loading: "Chargement…", activityWorkspace: "Espace d’administration", activityUsers: "Comptes utilisateurs", activityOrders: "Commandes", now: "Maintenant", today: "Aujourd’hui" },
  ar: { eyebrow: "الإدارة / مركز التحكم", title: "صباح الخير", description: "نظرة واضحة على ما يحتاج إلى اهتمامك في المنصة.", openOrders: "الطلبات المفتوحة", pendingReviews: "المراجعات المعلقة", users: "حسابات المستخدمين", auditEvents: "أحداث التدقيق", attention: "يحتاج إلى اهتمامك", attentionBody: "قوائم تحتاج إلى قرار أو متابعة.", recentActivity: "النشاط الأخير", recentActivityBody: "آخر التغييرات في مساحة العمل.", applications: "طلبات حرفيين بانتظار المراجعة", applicationsBody: "راجع الأدلة واتخذ قراراً مع ذكر السبب.", moderation: "منتجات تحتاج إلى مراجعة", moderationBody: "اعتمد أو ارفض أو اطلب تعديلات.", orders: "الطلبات قيد التنفيذ", ordersBody: "تابع التنفيذ وسياق العميل.", audit: "سجل التدقيق", viewQueue: "عرض القائمة", viewAll: "عرض الكل", loading: "جار التحميل…", activityWorkspace: "مساحة إدارة المنصة", activityUsers: "حسابات المستخدمين", activityOrders: "الطلبات", now: "الآن", today: "اليوم" },
  es: { eyebrow: "ADMINISTRACIÓN / CENTRO DE CONTROL", title: "Buenos días", description: "Una vista clara de lo que necesita tu atención en el marketplace.", openOrders: "Pedidos abiertos", pendingReviews: "Revisiones pendientes", users: "Cuentas de usuario", auditEvents: "Eventos de auditoría", attention: "Necesita tu atención", attentionBody: "Colas que necesitan una decisión o seguimiento.", recentActivity: "Actividad reciente", recentActivityBody: "Últimos cambios en tu espacio de trabajo.", applications: "Solicitudes de artesanos pendientes", applicationsBody: "Revisa las pruebas y decide con un motivo.", moderation: "Productos necesitan moderación", moderationBody: "Aprueba, rechaza o solicita cambios.", orders: "Pedidos en curso", ordersBody: "Supervisa el cumplimiento y el contexto del cliente.", audit: "Registro de auditoría", viewQueue: "Ver cola", viewAll: "Ver todo", loading: "Cargando…", activityWorkspace: "Espacio de administración", activityUsers: "Cuentas de usuario", activityOrders: "Pedidos", now: "Ahora", today: "Hoy" },
};

async function countFrom(path: string) {
  try {
    const response = await fetch(path, { cache: "no-store" });
    if (!response.ok) return null;
    const body = (await response.json()) as PageResponse;
    return typeof body.total === "number" ? body.total : null;
  } catch {
    return null;
  }
}

export function AdminOverview({ locale }: { locale: Locale }) {
  const text = copy[locale];
  const [stats, setStats] = useState({ orders: null as number | null, reviews: null as number | null, users: null as number | null, audit: null as number | null });

  useEffect(() => {
    let active = true;
    void Promise.all([
      countFrom("/api/admin/orders?page=1&pageSize=1"),
      countFrom("/api/admin/artisan-applications?status=SUBMITTED&page=1&pageSize=1"),
      countFrom("/api/admin/users?page=1&pageSize=1&role=user"),
      countFrom("/api/admin/audit-events?page=1&pageSize=1"),
    ]).then(([orders, reviews, users, audit]) => { if (active) setStats({ orders, reviews, users, audit }); });
    return () => { active = false; };
  }, []);

  const number = (value: number | null) => value === null ? text.loading : new Intl.NumberFormat(locale).format(value);
  const statCards = [
    { label: text.openOrders, value: number(stats.orders), tone: "blue", icon: ShoppingBag },
    { label: text.pendingReviews, value: number(stats.reviews), tone: "amber", icon: ClipboardCheck },
    { label: text.users, value: number(stats.users), tone: "purple", icon: PackageCheck },
    { label: text.auditEvents, value: number(stats.audit), tone: "green", icon: FileClock },
  ];

  return <div className="admin-overview">
    <header className="admin-overview__heading"><div><p className="admin-overview__eyebrow">{text.eyebrow}</p><h1>{text.title} <span>✦</span></h1><p>{text.description}</p></div></header>
    <div className="admin-overview__stats">{statCards.map(({ label, value, tone, icon: Icon }) => <section className={`admin-overview__stat admin-overview__stat--${tone}`} key={label}><div><p>{label}</p><strong>{value}</strong></div><Icon aria-hidden="true" size={18} /></section>)}</div>
    <div className="admin-overview__grid">
      <section className="admin-overview__card"><div className="admin-overview__card-heading"><div><h2>{text.attention}</h2><p>{text.attentionBody}</p></div><LocalizedLink locale={locale} href="/admin/artisan-applications">{text.viewQueue}<ArrowUpRight aria-hidden="true" size={15} /></LocalizedLink></div><div className="admin-overview__queue"><LocalizedLink locale={locale} href="/admin/artisan-applications" className="admin-overview__queue-item"><span className="admin-overview__queue-icon admin-overview__queue-icon--amber"><ClipboardCheck size={16} /></span><span><strong>{text.applications}</strong><small>{text.applicationsBody}</small></span><ArrowUpRight size={15} /></LocalizedLink><LocalizedLink locale={locale} href="/admin/moderation" className="admin-overview__queue-item"><span className="admin-overview__queue-icon admin-overview__queue-icon--blue"><PackageCheck size={16} /></span><span><strong>{text.moderation}</strong><small>{text.moderationBody}</small></span><ArrowUpRight size={15} /></LocalizedLink><LocalizedLink locale={locale} href="/admin/orders" className="admin-overview__queue-item"><span className="admin-overview__queue-icon admin-overview__queue-icon--purple"><ShoppingBag size={16} /></span><span><strong>{text.orders}</strong><small>{text.ordersBody}</small></span><ArrowUpRight size={15} /></LocalizedLink></div></section>
      <section className="admin-overview__card"><div className="admin-overview__card-heading"><div><h2>{text.recentActivity}</h2><p>{text.recentActivityBody}</p></div><LocalizedLink locale={locale} href="/admin/audit">{text.viewAll}<ArrowUpRight aria-hidden="true" size={15} /></LocalizedLink></div><div className="admin-overview__activity"><div><span>SM</span><p><strong>{text.activityWorkspace}</strong><small>{text.audit}</small></p><time>{text.now}</time></div><div><span>AC</span><p><strong>{text.activityUsers}</strong><small>{text.recentActivityBody}</small></p><time>{text.today}</time></div><div><span>AI</span><p><strong>{text.activityOrders}</strong><small>{text.openOrders}</small></p><time>{text.today}</time></div></div></section>
    </div>
  </div>;
}
