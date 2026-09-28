"use client";

import { Bell, CheckCheck, X } from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import { formatFullDateTime } from "@/core/lib/format";
import type { Locale } from "@/core/lib/i18n";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

type NotificationItem = {
  id: string;
  eventType: string;
  titleKey: string;
  bodyKey: string;
  readAt?: string | null;
  createdAt: string;
};

type NotificationResponse = {
  items?: NotificationItem[];
  unreadCount?: number;
};

const copy: Record<Locale, { label: string; title: string; empty: string; markAll: string; close: string; loading: string }> = {
  en: { label: "Notifications", title: "Notifications", empty: "You are all caught up.", markAll: "Mark all as read", close: "Close", loading: "Loading…" },
  fr: { label: "Notifications", title: "Notifications", empty: "Tout est à jour.", markAll: "Tout marquer comme lu", close: "Fermer", loading: "Chargement…" },
  ar: { label: "الإشعارات", title: "الإشعارات", empty: "لا توجد إشعارات جديدة.", markAll: "تحديد الكل كمقروء", close: "إغلاق", loading: "جارٍ التحميل…" },
  es: { label: "Notificaciones", title: "Notificaciones", empty: "Todo está al día.", markAll: "Marcar todo como leído", close: "Cerrar", loading: "Cargando…" },
};

const titles: Record<Locale, Record<string, string>> = {
  en: {
    "notifications.artisanApplication.title": "Artisan application update",
    "notifications.verification.title": "Verification update",
    "notifications.membership.title": "Membership update",
    "notifications.workshop.title": "Workshop update",
    "notifications.product.title": "Product update",
    "notifications.order.title": "Order update",
    "notifications.account.title": "Account update",
    "notifications.inventory.title": "Inventory update",
    "notifications.warehouse.title": "Warehouse update",
    "notifications.artisanApplication.submitted.title": "Artisan application submitted",
    "notifications.artisanApplication.approved.title": "Artisan application approved",
    "notifications.artisanApplication.changesRequested.title": "Changes requested for artisan application",
    "notifications.artisanApplication.rejected.title": "Artisan application rejected",
    "notifications.verification.submitted.title": "Verification submitted",
    "notifications.verification.verified.title": "Artisan verification approved",
    "notifications.verification.changesRequested.title": "Changes requested for verification",
    "notifications.verification.rejected.title": "Artisan verification rejected",
    "notifications.membership.active.title": "Artisan membership activated",
    "notifications.membership.suspended.title": "Artisan membership suspended",
    "notifications.membership.closed.title": "Artisan membership closed",
    "notifications.workshop.created.title": "Workshop created",
    "notifications.workshop.statusChanged.title": "Workshop status changed",
    "notifications.workshop.deleted.title": "Workshop deleted",
    "notifications.product.draftCreated.title": "Product draft created",
    "notifications.product.submitted.title": "Product submitted for review",
    "notifications.product.moderation.title": "Product moderation decision",
    "notifications.product.autoActivated.title": "Product activated",
    "notifications.product.archived.title": "Product archived",
    "notifications.order.checkout.title": "Order placed",
    "notifications.order.cancelled.title": "Order cancelled",
    "notifications.order.returnRecorded.title": "Order return recorded",
    "notifications.warehouse.received.title": "Warehouse reception recorded",
    "notifications.warehouse.inspected.title": "Warehouse inspection completed",
    "notifications.warehouse.inventoryAccepted.title": "Inventory accepted",
    "notifications.inventory.adjusted.title": "Inventory adjusted",
    "notifications.account.statusChanged.title": "Account status changed",
    "notifications.account.rolesChanged.title": "Account roles changed",
  },
  fr: {
    "notifications.artisanApplication.title": "Mise à jour de la demande artisan",
    "notifications.verification.title": "Mise à jour de la vérification",
    "notifications.membership.title": "Mise à jour de l'adhésion",
    "notifications.workshop.title": "Mise à jour de l'atelier",
    "notifications.product.title": "Mise à jour du produit",
    "notifications.order.title": "Mise à jour de la commande",
    "notifications.account.title": "Mise à jour du compte",
    "notifications.inventory.title": "Mise à jour du stock",
    "notifications.warehouse.title": "Mise à jour de l'entrepôt",
    "notifications.artisanApplication.submitted.title": "Demande artisan soumise",
    "notifications.artisanApplication.approved.title": "Demande artisan approuvée",
    "notifications.artisanApplication.changesRequested.title": "Modifications demandées pour la demande artisan",
    "notifications.artisanApplication.rejected.title": "Demande artisan rejetée",
    "notifications.verification.submitted.title": "Vérification soumise",
    "notifications.verification.verified.title": "Vérification artisan approuvée",
    "notifications.verification.changesRequested.title": "Modifications demandées pour la vérification",
    "notifications.verification.rejected.title": "Vérification artisan rejetée",
    "notifications.membership.active.title": "Adhésion artisan activée",
    "notifications.membership.suspended.title": "Adhésion artisan suspendue",
    "notifications.membership.closed.title": "Adhésion artisan clôturée",
    "notifications.workshop.created.title": "Atelier créé",
    "notifications.workshop.statusChanged.title": "Statut de l'atelier modifié",
    "notifications.workshop.deleted.title": "Atelier supprimé",
    "notifications.product.draftCreated.title": "Brouillon de produit créé",
    "notifications.product.submitted.title": "Produit soumis pour examen",
    "notifications.product.moderation.title": "Décision de modération du produit",
    "notifications.product.autoActivated.title": "Produit activé",
    "notifications.product.archived.title": "Produit archivé",
    "notifications.order.checkout.title": "Commande passée",
    "notifications.order.cancelled.title": "Commande annulée",
    "notifications.order.returnRecorded.title": "Retour de commande enregistré",
    "notifications.warehouse.received.title": "Réception d'entrepôt enregistrée",
    "notifications.warehouse.inspected.title": "Inspection d'entrepôt terminée",
    "notifications.warehouse.inventoryAccepted.title": "Stock accepté",
    "notifications.inventory.adjusted.title": "Stock ajusté",
    "notifications.account.statusChanged.title": "Statut du compte modifié",
    "notifications.account.rolesChanged.title": "Rôles du compte modifiés",
  },
  ar: {
    "notifications.artisanApplication.title": "تحديث طلب الحرفي",
    "notifications.verification.title": "تحديث التحقق",
    "notifications.membership.title": "تحديث العضوية",
    "notifications.workshop.title": "تحديث الورشة",
    "notifications.product.title": "تحديث المنتج",
    "notifications.order.title": "تحديث الطلب",
    "notifications.account.title": "تحديث الحساب",
    "notifications.inventory.title": "تحديث المخزون",
    "notifications.warehouse.title": "تحديث المستودع",
    "notifications.artisanApplication.submitted.title": "تم إرسال طلب الحرفي",
    "notifications.artisanApplication.approved.title": "تمت الموافقة على طلب الحرفي",
    "notifications.artisanApplication.changesRequested.title": "مطلوب تعديل طلب الحرفي",
    "notifications.artisanApplication.rejected.title": "تم رفض طلب الحرفي",
    "notifications.verification.submitted.title": "تم إرسال التحقق",
    "notifications.verification.verified.title": "تمت الموافقة على تحقق الحرفي",
    "notifications.verification.changesRequested.title": "مطلوب تعديل التحقق",
    "notifications.verification.rejected.title": "تم رفض تحقق الحرفي",
    "notifications.membership.active.title": "تم تفعيل عضوية الحرفي",
    "notifications.membership.suspended.title": "تم تعليق عضوية الحرفي",
    "notifications.membership.closed.title": "تم إغلاق عضوية الحرفي",
    "notifications.workshop.created.title": "تم إنشاء الورشة",
    "notifications.workshop.statusChanged.title": "تغيرت حالة الورشة",
    "notifications.workshop.deleted.title": "تم حذف الورشة",
    "notifications.product.draftCreated.title": "تم إنشاء مسودة المنتج",
    "notifications.product.submitted.title": "تم إرسال المنتج للمراجعة",
    "notifications.product.moderation.title": "قرار مراجعة المنتج",
    "notifications.product.autoActivated.title": "تم تفعيل المنتج",
    "notifications.product.archived.title": "تمت أرشفة المنتج",
    "notifications.order.checkout.title": "تم إنشاء الطلب",
    "notifications.order.cancelled.title": "تم إلغاء الطلب",
    "notifications.order.returnRecorded.title": "تم تسجيل إرجاع الطلب",
    "notifications.warehouse.received.title": "تم تسجيل استلام المستودع",
    "notifications.warehouse.inspected.title": "اكتمل فحص المستودع",
    "notifications.warehouse.inventoryAccepted.title": "تم قبول المخزون",
    "notifications.inventory.adjusted.title": "تم تعديل المخزون",
    "notifications.account.statusChanged.title": "تغيرت حالة الحساب",
    "notifications.account.rolesChanged.title": "تغيرت أدوار الحساب",
  },
  es: {
    "notifications.artisanApplication.title": "Actualización de la solicitud artesanal",
    "notifications.verification.title": "Actualización de verificación",
    "notifications.membership.title": "Actualización de membresía",
    "notifications.workshop.title": "Actualización del taller",
    "notifications.product.title": "Actualización del producto",
    "notifications.order.title": "Actualización del pedido",
    "notifications.account.title": "Actualización de la cuenta",
    "notifications.inventory.title": "Actualización del inventario",
    "notifications.warehouse.title": "Actualización del almacén",
    "notifications.artisanApplication.submitted.title": "Solicitud artesanal enviada",
    "notifications.artisanApplication.approved.title": "Solicitud artesanal aprobada",
    "notifications.artisanApplication.changesRequested.title": "Cambios solicitados para la solicitud artesanal",
    "notifications.artisanApplication.rejected.title": "Solicitud artesanal rechazada",
    "notifications.verification.submitted.title": "Verificación enviada",
    "notifications.verification.verified.title": "Verificación artesanal aprobada",
    "notifications.verification.changesRequested.title": "Cambios solicitados para la verificación",
    "notifications.verification.rejected.title": "Verificación artesanal rechazada",
    "notifications.membership.active.title": "Membresía artesanal activada",
    "notifications.membership.suspended.title": "Membresía artesanal suspendida",
    "notifications.membership.closed.title": "Membresía artesanal cerrada",
    "notifications.workshop.created.title": "Taller creado",
    "notifications.workshop.statusChanged.title": "Estado del taller cambiado",
    "notifications.workshop.deleted.title": "Taller eliminado",
    "notifications.product.draftCreated.title": "Borrador de producto creado",
    "notifications.product.submitted.title": "Producto enviado para revisión",
    "notifications.product.moderation.title": "Decisión de moderación del producto",
    "notifications.product.autoActivated.title": "Producto activado",
    "notifications.product.archived.title": "Producto archivado",
    "notifications.order.checkout.title": "Pedido realizado",
    "notifications.order.cancelled.title": "Pedido cancelado",
    "notifications.order.returnRecorded.title": "Devolución del pedido registrada",
    "notifications.warehouse.received.title": "Recepción de almacén registrada",
    "notifications.warehouse.inspected.title": "Inspección de almacén completada",
    "notifications.warehouse.inventoryAccepted.title": "Inventario aceptado",
    "notifications.inventory.adjusted.title": "Inventario ajustado",
    "notifications.account.statusChanged.title": "Estado de la cuenta cambiado",
    "notifications.account.rolesChanged.title": "Roles de la cuenta cambiados",
  },
};

const bodies: Record<Locale, Record<string, string>> = {
  en: {
    "notifications.artisanApplication.body": "A workflow step needs your attention.",
    "notifications.verification.body": "Your verification workflow has changed.",
    "notifications.membership.body": "Your artisan membership has changed.",
    "notifications.workshop.body": "A workshop workflow step was recorded.",
    "notifications.product.body": "A product workflow step was recorded.",
    "notifications.order.body": "An order workflow step was recorded.",
    "notifications.account.body": "Your account settings were updated.",
    "notifications.inventory.body": "An inventory workflow step was recorded.",
    "notifications.warehouse.body": "A warehouse workflow step was recorded.",
  },
  fr: {
    "notifications.artisanApplication.body": "Une étape du workflow nécessite votre attention.",
    "notifications.verification.body": "Votre workflow de vérification a changé.",
    "notifications.membership.body": "Votre adhésion artisan a changé.",
    "notifications.workshop.body": "Une étape du workflow de l'atelier a été enregistrée.",
    "notifications.product.body": "Une étape du workflow du produit a été enregistrée.",
    "notifications.order.body": "Une étape du workflow de la commande a été enregistrée.",
    "notifications.account.body": "Les paramètres de votre compte ont été mis à jour.",
    "notifications.inventory.body": "Une étape du workflow du stock a été enregistrée.",
    "notifications.warehouse.body": "Une étape du workflow de l'entrepôt a été enregistrée.",
  },
  ar: {
    "notifications.artisanApplication.body": "تحتاج إحدى مراحل سير العمل إلى انتباهك.",
    "notifications.verification.body": "تغيرت مرحلة التحقق الخاصة بك.",
    "notifications.membership.body": "تغيرت عضوية الحرفي الخاصة بك.",
    "notifications.workshop.body": "تم تسجيل مرحلة من سير عمل الورشة.",
    "notifications.product.body": "تم تسجيل مرحلة من سير عمل المنتج.",
    "notifications.order.body": "تم تسجيل مرحلة من سير عمل الطلب.",
    "notifications.account.body": "تم تحديث إعدادات حسابك.",
    "notifications.inventory.body": "تم تسجيل مرحلة من سير عمل المخزون.",
    "notifications.warehouse.body": "تم تسجيل مرحلة من سير عمل المستودع.",
  },
  es: {
    "notifications.artisanApplication.body": "Una etapa del flujo de trabajo requiere tu atención.",
    "notifications.verification.body": "Tu flujo de verificación ha cambiado.",
    "notifications.membership.body": "Tu membresía artesanal ha cambiado.",
    "notifications.workshop.body": "Se registró una etapa del flujo del taller.",
    "notifications.product.body": "Se registró una etapa del flujo del producto.",
    "notifications.order.body": "Se registró una etapa del flujo del pedido.",
    "notifications.account.body": "Se actualizaron los ajustes de tu cuenta.",
    "notifications.inventory.body": "Se registró una etapa del flujo del inventario.",
    "notifications.warehouse.body": "Se registró una etapa del flujo del almacén.",
  },
};

function titleFor(item: NotificationItem, locale: Locale) {
  return titles[locale][item.titleKey] ?? item.titleKey.replace(/^notifications\./, "").replace(/\./g, " ");
}

function bodyFor(item: NotificationItem, locale: Locale) {
  return bodies[locale][item.bodyKey] ?? "A workflow update was recorded.";
}

function websocketURL() {
  if (process.env.NEXT_PUBLIC_API_WS_URL) return process.env.NEXT_PUBLIC_API_WS_URL;
  if (typeof window === "undefined") return "ws://localhost:8088/api/v1/notifications/ws";
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${window.location.hostname}:8088/api/v1/notifications/ws`;
}

export function NotificationCenter({ locale }: { locale: Locale }) {
  const auth = useOptionalAuth();
  const [open, setOpen] = useState(false);
  const [items, setItems] = useState<NotificationItem[]>([]);
  const [unread, setUnread] = useState(0);
  const [loading, setLoading] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);
  const retryRef = useRef<number | null>(null);
  const text = useMemo(() => copy[locale], [locale]);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const response = await fetch("/api/notifications?page=1&pageSize=20", { cache: "no-store" });
      if (!response.ok) return;
      const result = (await response.json()) as NotificationResponse;
      setItems(result.items ?? []);
      setUnread(result.unreadCount ?? 0);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (auth?.status !== "authenticated") {
      socketRef.current?.close();
      socketRef.current = null;
      const resetTimer = window.setTimeout(() => {
        setItems([]);
        setUnread(0);
      }, 0);
      return () => window.clearTimeout(resetTimer);
    }
    let cancelled = false;
    const loadTimer = window.setTimeout(() => { void load(); }, 0);
    const connect = async () => {
      const response = await fetch("/api/notifications/ws-ticket", { cache: "no-store" });
      if (!response.ok || cancelled) return;
      const { ticket } = (await response.json()) as { ticket?: string };
      if (!ticket || cancelled) return;
      const socket = new WebSocket(`${websocketURL()}?ticket=${encodeURIComponent(ticket)}`);
      socketRef.current = socket;
      socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as { type?: string; notification?: NotificationItem };
          if (message.type !== "notification" || !message.notification) return;
          setItems((current) => [message.notification as NotificationItem, ...current.filter((item) => item.id !== message.notification?.id)].slice(0, 20));
          setUnread((count) => count + 1);
        } catch {
          // Ignore malformed frames from a stale connection.
        }
      };
      socket.onclose = () => {
        if (!cancelled) retryRef.current = window.setTimeout(() => void connect(), 5000);
      };
    };
    void connect().catch(() => undefined);
    return () => {
      cancelled = true;
      window.clearTimeout(loadTimer);
      if (retryRef.current !== null) window.clearTimeout(retryRef.current);
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, [auth?.status, load]);

  const markRead = async (id: string) => {
    const item = items.find((value) => value.id === id);
    if (!item || item.readAt) return;
    setItems((current) => current.map((value) => value.id === id ? { ...value, readAt: new Date().toISOString() } : value));
    setUnread((count) => Math.max(0, count - 1));
    await fetch(`/api/notifications/${encodeURIComponent(id)}/read`, { method: "POST" });
  };

  const markAllRead = async () => {
    setItems((current) => current.map((item) => ({ ...item, readAt: item.readAt ?? new Date().toISOString() })));
    setUnread(0);
    await fetch("/api/notifications/read-all", { method: "POST" });
  };

  if (!auth?.user) return null;
  return (
    <div className="relative">
      <button type="button" className="relative flex min-h-11 min-w-11 items-center justify-center" aria-label={text.label} aria-expanded={open} onClick={() => setOpen((value) => !value)}>
        <Bell size={19} strokeWidth={1.5} />
        {unread > 0 && <span className="absolute end-0 top-1 min-w-4 rounded-full bg-foreground px-1 text-[0.625rem] text-background" aria-label={`${unread} unread`}>{unread > 99 ? "99+" : unread}</span>}
      </button>
      {open && <div className="absolute end-0 top-12 z-[70] w-[min(22rem,calc(100vw-2rem))] overflow-hidden rounded-md border border-border bg-background shadow-xl" role="dialog" aria-label={text.title}>
        <div className="flex items-center justify-between border-b border-border px-4 py-3">
          <h2 className="font-serif text-xl">{text.title}{unread > 0 ? ` (${unread})` : ""}</h2>
          <button type="button" className="grid size-8 place-items-center" aria-label={text.close} onClick={() => setOpen(false)}><X size={16} /></button>
        </div>
        <div className="max-h-[min(28rem,65vh)] overflow-y-auto">
          {loading ? <p className="px-4 py-8 text-sm text-muted-foreground">{text.loading}</p> : items.length === 0 ? <p className="px-4 py-8 text-sm text-muted-foreground">{text.empty}</p> : items.map((item) => <button key={item.id} type="button" className={`block w-full border-b border-border px-4 py-3 text-start transition-colors hover:bg-muted/40 ${item.readAt ? "" : "bg-muted/20"}`} onClick={() => void markRead(item.id)}>
            <span className="flex items-start justify-between gap-3"><strong className="text-sm">{titleFor(item, locale)}</strong>{!item.readAt && <span className="mt-1 size-2 shrink-0 rounded-full bg-primary" aria-label="Unread" />}</span>
            <span className="mt-1 block text-xs text-muted-foreground">{bodyFor(item, locale)}</span>
            <span className="mt-2 block text-[0.6875rem] text-muted-foreground">{formatFullDateTime(item.createdAt, locale)}</span>
          </button>)}
        </div>
        {unread > 0 && <button type="button" className="flex w-full items-center justify-center gap-2 border-t border-border px-4 py-3 text-xs font-medium" onClick={() => void markAllRead()}><CheckCheck size={15} />{text.markAll}</button>}
      </div>}
    </div>
  );
}
