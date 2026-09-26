"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

import { LocalizedLink } from "@/core/components/shared/localized-link";
import { Button } from "@/core/components/ui/button";
import { Combobox, MultiCombobox } from "@/core/components/ui/combobox";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { formatFullDate, formatFullDateTime } from "@/core/lib/format";
import type { Locale } from "@/core/lib/i18n";
import { hasCapability } from "@/features/auth/types";
import { useAuth } from "@/features/auth/viewmodel/auth-context";

type User = {
  id: string;
  email: string;
  displayName: string;
  status: string;
  roles: string[];
  createdAt: string;
};
type AuditEvent = {
  id: string;
  eventType: string;
  targetType: string;
  targetId?: string;
  reason?: string;
  occurredAt: string;
};
type Page<T> = { items: T[]; page: number; pageSize: number; total: number };
type Labels = {
  dashboard: string;
  welcome: string;
  profile: string;
  logout: string;
  applications: string;
  moderation: string;
  totalUsers: string;
  visibleUsers: string;
  auditEvents: string;
  users: string;
  email: string;
  status: string;
  roles: string;
  created: string;
  actions: string;
  details: string;
  close: string;
  noUsers: string;
  previous: string;
  next: string;
  pageOf: string;
  userDetails: string;
  save: string;
  reason: string;
  update: string;
  load: string;
  saved: string;
  statusSaved: string;
  filters: string;
  targetType: string;
  eventType: string;
  search: string;
  noAudit: string;
  loggingOut: string;
};

const pageSize = 10;
const roleOptions = [
  ["customer", "Customer"],
  ["artisan", "Artisan"],
  ["moderator", "Moderator"],
  ["warehouse_agent", "Warehouse agent"],
  ["administrator", "Administrator"],
].map(([value, label]) => ({ value, label }));
const accountStatusOptions = ["ACTIVE", "SUSPENDED", "DISABLED"].map((value) => ({
  value,
  label: value.replaceAll("_", " "),
}));
const labels: Record<Locale, Labels> = {
  en: {
    dashboard: "Administration dashboard",
    welcome: "Control centre",
    profile: "Manage profile",
    logout: "Sign out",
    applications: "Artisan applications",
    moderation: "Product moderation",
    totalUsers: "Total users",
    visibleUsers: "Rows on this page",
    auditEvents: "Audit events",
    users: "User accounts",
    email: "Email",
    status: "Status",
    roles: "Roles",
    created: "Created",
    actions: "Actions",
    details: "Details",
    close: "Close details",
    noUsers: "No user accounts found.",
    previous: "Previous",
    next: "Next",
    pageOf: "Page {page} of {pages}",
    userDetails: "User details",
    save: "Save roles",
    reason: "Reason for status change",
    update: "Update status",
    load: "Unable to load administration data",
    saved: "Roles updated",
    statusSaved: "Account status updated",
    filters: "Audit filters",
    targetType: "Target type",
    eventType: "Event type",
    search: "Search audit",
    noAudit: "No audit events match these filters.",
    loggingOut: "Signing out…",
  },
  fr: {
    dashboard: "Tableau de bord d’administration",
    welcome: "Centre de contrôle",
    profile: "Gérer le profil",
    logout: "Déconnexion",
    applications: "Candidatures d’artisans",
    moderation: "Modération des produits",
    totalUsers: "Utilisateurs au total",
    visibleUsers: "Lignes de cette page",
    auditEvents: "Événements d’audit",
    users: "Comptes utilisateurs",
    email: "E-mail",
    status: "Statut",
    roles: "Rôles",
    created: "Créé le",
    actions: "Actions",
    details: "Détails",
    close: "Fermer les détails",
    noUsers: "Aucun compte utilisateur trouvé.",
    previous: "Précédent",
    next: "Suivant",
    pageOf: "Page {page} sur {pages}",
    userDetails: "Détails de l’utilisateur",
    save: "Enregistrer les rôles",
    reason: "Motif du changement de statut",
    update: "Modifier le statut",
    load: "Impossible de charger les données d’administration",
    saved: "Rôles mis à jour",
    statusSaved: "Statut du compte mis à jour",
    filters: "Filtres d’audit",
    targetType: "Type de cible",
    eventType: "Type d’événement",
    search: "Rechercher dans l’audit",
    noAudit: "Aucun événement ne correspond à ces filtres.",
    loggingOut: "Déconnexion…",
  },
  ar: {
    dashboard: "لوحة تحكم الإدارة",
    welcome: "مركز التحكم",
    profile: "إدارة الملف الشخصي",
    logout: "تسجيل الخروج",
    applications: "طلبات الحرفيين",
    moderation: "مراجعة المنتجات",
    totalUsers: "إجمالي المستخدمين",
    visibleUsers: "صفوف هذه الصفحة",
    auditEvents: "أحداث التدقيق",
    users: "حسابات المستخدمين",
    email: "البريد الإلكتروني",
    status: "الحالة",
    roles: "الأدوار",
    created: "تاريخ الإنشاء",
    actions: "الإجراءات",
    details: "التفاصيل",
    close: "إغلاق التفاصيل",
    noUsers: "لم يتم العثور على حسابات.",
    previous: "السابق",
    next: "التالي",
    pageOf: "الصفحة {page} من {pages}",
    userDetails: "تفاصيل المستخدم",
    save: "حفظ الأدوار",
    reason: "سبب تغيير الحالة",
    update: "تحديث الحالة",
    load: "تعذر تحميل بيانات الإدارة",
    saved: "تم تحديث الأدوار",
    statusSaved: "تم تحديث حالة الحساب",
    filters: "فلاتر التدقيق",
    targetType: "نوع الهدف",
    eventType: "نوع الحدث",
    search: "بحث التدقيق",
    noAudit: "لا توجد أحداث تدقيق مطابقة.",
    loggingOut: "جارٍ تسجيل الخروج…",
  },
  es: {
    dashboard: "Panel de administración",
    welcome: "Centro de control",
    profile: "Gestionar perfil",
    logout: "Cerrar sesión",
    applications: "Solicitudes de artesanos",
    moderation: "Moderación de productos",
    totalUsers: "Usuarios totales",
    visibleUsers: "Filas de esta página",
    auditEvents: "Eventos de auditoría",
    users: "Cuentas de usuario",
    email: "Correo electrónico",
    status: "Estado",
    roles: "Roles",
    created: "Creado",
    actions: "Acciones",
    details: "Detalles",
    close: "Cerrar detalles",
    noUsers: "No se encontraron cuentas.",
    previous: "Anterior",
    next: "Siguiente",
    pageOf: "Página {page} de {pages}",
    userDetails: "Detalles del usuario",
    save: "Guardar roles",
    reason: "Motivo del cambio de estado",
    update: "Actualizar estado",
    load: "No se pudieron cargar los datos de administración",
    saved: "Roles actualizados",
    statusSaved: "Estado de la cuenta actualizado",
    filters: "Filtros de auditoría",
    targetType: "Tipo de objetivo",
    eventType: "Tipo de evento",
    search: "Buscar auditoría",
    noAudit: "No hay eventos que coincidan con estos filtros.",
    loggingOut: "Cerrando sesión…",
  },
};

export function AdminControlPanel({
  locale,
  showChrome = true,
  section = "overview",
}: {
  locale: Locale;
  showChrome?: boolean;
  section?: "overview" | "users" | "audit";
}) {
  const text = labels[locale];
  const router = useRouter();
  const { user, logout } = useAuth();
  const canReadAudit = hasCapability(user, "admin.audit.read");
  const canManageUsers = hasCapability(user, "admin.users.write");
  const [users, setUsers] = useState<User[]>([]);
  const [audit, setAudit] = useState<AuditEvent[]>([]);
  const [userPage, setUserPage] = useState(1);
  const [userTotal, setUserTotal] = useState(0);
  const [auditTotal, setAuditTotal] = useState(0);
  const [auditPage, setAuditPage] = useState(1);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [selectedAudit, setSelectedAudit] = useState<AuditEvent | null>(null);
  useEscapeKey(() => {
    if (selectedAudit) setSelectedAudit(null);
    else if (selectedUser) setSelectedUser(null);
  }, selectedAudit !== null || selectedUser !== null);
  const [roles, setRoles] = useState<Record<string, string[]>>({});
  const [statuses, setStatuses] = useState<Record<string, string>>({});
  const [reasons, setReasons] = useState<Record<string, string>>({});
  const [targetType, setTargetType] = useState("");
  const [eventType, setEventType] = useState("");
  const [loading, setLoading] = useState(true);
  const [loadingUsers, setLoadingUsers] = useState(false);
  const [loadingAudit, setLoadingAudit] = useState(false);
  const [busy, setBusy] = useState<string | null>(null);
  const [loggingOut, setLoggingOut] = useState(false);

  const loadUsers = useCallback(async (page: number) => {
    setLoadingUsers(true);
    try {
      const response = await fetch(
        `/api/admin/users?page=${page}&pageSize=${pageSize}`,
        { cache: "no-store" },
      );
      if (!response.ok) throw new Error();
      const data = (await response.json()) as Page<User>;
      setUsers(data.items ?? []);
      setUserPage(data.page ?? page);
      setUserTotal(data.total ?? 0);
      setSelectedUser(
        (current) =>
          data.items?.find((item) => item.id === current?.id) ?? null,
      );
    } finally {
      setLoadingUsers(false);
    }
  }, []);

  const loadAudit = useCallback(
    async (
      filters: { targetType?: string; eventType?: string } = {},
      page = auditPage,
    ) => {
      setLoadingAudit(true);
      try {
        const query = new URLSearchParams({
          page: String(page),
          pageSize: "10",
        });
        if (filters.targetType) query.set("targetType", filters.targetType);
        if (filters.eventType) query.set("eventType", filters.eventType);
        const response = await fetch(
          `/api/admin/audit-events?${query.toString()}`,
          { cache: "no-store" },
        );
        if (!response.ok) throw new Error();
        const data = (await response.json()) as Page<AuditEvent>;
        setAudit(data.items ?? []);
        setAuditTotal(data.total ?? 0);
        setAuditPage(data.page ?? page);
      } finally {
        setLoadingAudit(false);
      }
    },
    [auditPage],
  );

  useEffect(() => {
    let active = true;
    const run = async () => {
      try {
        await Promise.all([loadUsers(1), canReadAudit ? loadAudit() : Promise.resolve()]);
      } catch {
        if (active) toast.error(text.load);
      } finally {
        if (active) setLoading(false);
      }
    };
    void run();
    return () => {
      active = false;
    };
  }, [canReadAudit, loadAudit, loadUsers, text.load]);

  const pageCount = Math.max(1, Math.ceil(userTotal / pageSize));
  const pageLabel = text.pageOf
    .replace("{page}", String(userPage))
    .replace("{pages}", String(pageCount));
  const selectedStatus = selectedUser
    ? (statuses[selectedUser.id] ?? selectedUser.status)
    : "";
  const selectedReason = selectedUser ? (reasons[selectedUser.id] ?? "") : "";
  const selectedRoles = selectedUser
    ? (roles[selectedUser.id] ?? selectedUser.roles)
    : [];
  const visibleStats = useMemo(
    () => [
      { label: text.totalUsers, value: userTotal },
      { label: text.visibleUsers, value: users.length },
      { label: text.auditEvents, value: auditTotal },
    ],
    [
      auditTotal,
      text.auditEvents,
      text.totalUsers,
      text.visibleUsers,
      userTotal,
      users.length,
    ],
  );

  const updateRoles = async () => {
    if (!selectedUser || !selectedRoles.length) return;
    setBusy(`roles:${selectedUser.id}`);
    try {
      const response = await fetch(
        `/api/admin/users/${encodeURIComponent(selectedUser.id)}/roles`,
        {
          method: "PATCH",
          headers: { "content-type": "application/json" },
          body: JSON.stringify({
            roles: selectedRoles,
          }),
        },
      );
      if (!response.ok) throw new Error();
      const updated = (await response.json()) as User;
      setUsers((current) =>
        current.map((item) => (item.id === updated.id ? updated : item)),
      );
      setSelectedUser(updated);
      toast.success(text.saved);
    } catch {
      toast.error(text.load);
    } finally {
      setBusy(null);
    }
  };

  const updateStatus = async () => {
    if (!selectedUser) return;
    setBusy(`status:${selectedUser.id}`);
    try {
      const response = await fetch(
        `/api/admin/users/${encodeURIComponent(selectedUser.id)}/status`,
        {
          method: "POST",
          headers: { "content-type": "application/json" },
          body: JSON.stringify({
            status: selectedStatus,
            reason: selectedReason,
          }),
        },
      );
      const body = await response.json().catch(() => undefined);
      if (!response.ok) throw new Error(body?.error?.message ?? text.load);
      const updated = body as User;
      setUsers((current) =>
        current.map((item) => (item.id === updated.id ? updated : item)),
      );
      setSelectedUser(updated);
      toast.success(text.statusSaved);
      await loadAudit({ targetType, eventType });
    } catch (error) {
      toast.error((error as Error).message);
    } finally {
      setBusy(null);
    }
  };

  const signOut = async () => {
    if (loggingOut) return;
    setLoggingOut(true);
    try {
      await logout();
    } finally {
      router.replace(`/${locale}/login`);
      router.refresh();
      setLoggingOut(false);
    }
  };

  if (loading)
    return (
      <section className="mt-12 border-t border-border pt-8" aria-busy="true">
        {text.load}
      </section>
    );
  return (
    <section className="mt-8 space-y-8" aria-labelledby="admin-dashboard-title">
      {showChrome && (
        <>
          <header className="flex flex-col gap-5 border-b border-border pb-8 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <p className="text-xs uppercase tracking-[0.18em] text-primary">
                {text.dashboard}
              </p>
              <h1
                id="admin-dashboard-title"
                className="mt-3 font-serif text-5xl"
              >
                {text.welcome}
              </h1>
              <p className="mt-3 text-sm text-muted-foreground">
                {user?.displayName} · {user?.email}
              </p>
            </div>
            <div className="flex flex-wrap gap-2">
              <LocalizedLink
                className="inline-flex min-h-12 items-center border border-current px-5 text-sm hover:bg-foreground hover:text-background"
                locale={locale}
                href="/account/profile"
              >
                {text.profile}
              </LocalizedLink>
              <Button
                type="button"
                variant="ghost"
                disabled={loggingOut}
                onClick={() => void signOut()}
              >
                {loggingOut ? text.loggingOut : text.logout}
              </Button>
            </div>
          </header>
          <nav className="flex flex-wrap gap-2" aria-label={text.dashboard}>
            <LocalizedLink
              className="border border-border px-4 py-3 text-sm hover:bg-muted"
              locale={locale}
              href="/admin/artisan-applications"
            >
              {text.applications}
            </LocalizedLink>
            <LocalizedLink
              className="border border-border px-4 py-3 text-sm hover:bg-muted"
              locale={locale}
              href="/admin/moderation"
            >
              {text.moderation}
            </LocalizedLink>
          </nav>
        </>
      )}
      {section !== "audit" && (
        <>
          <div className="grid gap-4 md:grid-cols-3">
            {visibleStats.map((stat, index) => (
              <section className="relative overflow-hidden border border-border bg-card p-6 shadow-sm" key={stat.label}>
                <span className="absolute end-0 top-0 h-1 w-20 bg-primary" aria-hidden="true" />
                <p className="text-xs uppercase tracking-[0.18em] text-muted-foreground">
                  {stat.label}
                </p>
                <p className="mt-4 font-serif text-5xl">{new Intl.NumberFormat(locale).format(stat.value)}</p>
                <p className="mt-2 text-xs text-muted-foreground">
                  {index === 0 ? "All registered accounts" : index === 1 ? "Current page" : "Recorded activity"}
                </p>
              </section>
            ))}
          </div>
          <div className="grid gap-8 xl:grid-cols-[minmax(0,1.6fr)_minmax(20rem,0.8fr)]">
            <section className="min-w-0" aria-labelledby="admin-users-title">
              <div className="flex flex-wrap items-end justify-between gap-3">
                <div>
                  <p className="text-xs uppercase tracking-widest text-primary">
                    {text.users}
                  </p>
                  <h2
                    id="admin-users-title"
                    className="mt-2 font-serif text-3xl"
                  >
                    {text.users}
                  </h2>
                </div>
                <span className="text-sm text-muted-foreground">
                  {pageLabel}
                </span>
              </div>
              <div className="mt-5 overflow-x-auto border border-border">
                <table className="w-full min-w-[40rem] text-start text-sm">
                  <thead className="border-b border-border bg-muted/40">
                    <tr>
                      <th className="px-4 py-3 font-medium">{text.email}</th>
                      <th className="px-4 py-3 font-medium">{text.status}</th>
                      <th className="px-4 py-3 font-medium">{text.roles}</th>
                      <th className="px-4 py-3 font-medium">{text.created}</th>
                      <th className="px-4 py-3 font-medium">{text.actions}</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {users.map((item) => (
                      <tr className="align-top" key={item.id}>
                        <td className="px-4 py-4">
                          <strong>{item.displayName || item.email}</strong>
                          <p className="mt-1 text-xs text-muted-foreground">
                            {item.email}
                          </p>
                        </td>
                        <td className="px-4 py-4">
                          <StatusBadge status={item.status} />
                        </td>
                        <td className="px-4 py-4 text-xs">
                          {item.roles.join(", ") || "—"}
                        </td>
                        <td className="whitespace-nowrap px-4 py-4 text-xs text-muted-foreground">
                          {formatFullDate(item.createdAt, locale)}
                        </td>
                        <td className="px-4 py-4">
                          <Button
                            type="button"
                            variant="outline"
                            onClick={() => setSelectedUser(item)}
                          >
                            {text.details}
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {!loadingUsers && users.length === 0 && (
                  <p className="p-6 text-sm text-muted-foreground">
                    {text.noUsers}
                  </p>
                )}
                {loadingUsers && (
                  <p
                    className="p-6 text-sm text-muted-foreground"
                    role="status"
                  >
                    {text.load}
                  </p>
                )}
              </div>
              <div className="mt-4 flex items-center justify-between gap-3">
                <Button
                  type="button"
                  variant="outline"
                  disabled={loadingUsers || userPage <= 1}
                  onClick={() => void loadUsers(userPage - 1)}
                >
                  {text.previous}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  disabled={loadingUsers || userPage >= pageCount}
                  onClick={() => void loadUsers(userPage + 1)}
                >
                  {text.next}
                </Button>
              </div>
            </section>
            {selectedUser && (
              <aside
                className="fixed inset-0 z-50 overflow-x-auto overflow-y-auto bg-foreground/50 p-4"
                aria-labelledby="admin-user-details-title"
              >
                <div className="mx-auto mt-8 max-h-[calc(100svh-4rem)] max-w-xl overflow-x-auto overflow-y-auto border border-border bg-background p-5">
                  {selectedUser ? (
                    <>
                      <div className="flex items-start justify-between gap-3">
                        <div>
                          <p className="text-xs uppercase tracking-widest text-primary">
                            {text.details}
                          </p>
                          <h2
                            id="admin-user-details-title"
                            className="mt-2 font-serif text-2xl"
                          >
                            {text.userDetails}
                          </h2>
                        </div>
                        <Button
                          type="button"
                          variant="ghost"
                          onClick={() => setSelectedUser(null)}
                        >
                          {text.close}
                        </Button>
                      </div>
                      <dl className="mt-6 space-y-3 text-sm">
                        <div>
                          <dt className="text-muted-foreground">
                            {text.email}
                          </dt>
                          <dd className="mt-1 break-all">
                            {selectedUser.email}
                          </dd>
                        </div>
                        <div>
                          <dt className="text-muted-foreground">ID</dt>
                          <dd className="mt-1 break-all text-xs">
                            {selectedUser.id}
                          </dd>
                        </div>
                        <div>
                          <dt className="text-muted-foreground">
                            {text.created}
                          </dt>
                          <dd className="mt-1">
                            {formatFullDateTime(selectedUser.createdAt, locale)}
                          </dd>
                        </div>
                      </dl>
                      {canManageUsers && <div className="mt-6 space-y-4">
                        <label className="block text-sm">
                          <span className="mb-2 block">{text.roles}</span>
                          <MultiCombobox
                            options={roleOptions}
                            value={selectedRoles}
                            onChange={(next) =>
                              setRoles((current) => ({
                                ...current,
                                [selectedUser.id]: next,
                              }))
                            }
                            ariaLabel={text.roles}
                            placeholder="Choose roles"
                          />
                        </label>
                        <Button
                          type="button"
                          variant="outline"
                          disabled={busy === `roles:${selectedUser.id}`}
                          onClick={() => void updateRoles()}
                        >
                          {text.save}
                        </Button>
                        <label className="block text-sm">
                          <span className="mb-2 block">{text.status}</span>
                          <Combobox
                            options={accountStatusOptions}
                            value={selectedStatus}
                            onChange={(next) =>
                              setStatuses((current) => ({
                                ...current,
                                [selectedUser.id]: next,
                              }))
                            }
                            ariaLabel={text.status}
                          />
                        </label>
                        <label className="block text-sm">
                          <span className="mb-2 block">{text.reason}</span>
                          <input
                            className="auth-input"
                            value={selectedReason}
                            onChange={(event) =>
                              setReasons((current) => ({
                                ...current,
                                [selectedUser.id]: event.target.value,
                              }))
                            }
                          />
                        </label>
                        <Button
                          type="button"
                          variant="outline"
                          disabled={
                            busy === `status:${selectedUser.id}` ||
                            selectedStatus === selectedUser.status
                          }
                          onClick={() => void updateStatus()}
                        >
                          {text.update}
                        </Button>
                      </div>}
                    </>
                  ) : (
                    <p className="text-sm text-muted-foreground">
                      {text.details}
                    </p>
                  )}
                </div>
              </aside>
            )}
          </div>
        </>
      )}
      {section !== "users" && canReadAudit && (
        <section
          className="border-t border-border pt-8"
          aria-labelledby="admin-audit-title"
        >
          <div className="flex flex-wrap items-end justify-between gap-3">
            <div>
              <p className="text-xs uppercase tracking-widest text-primary">
                {text.auditEvents}
              </p>
              <h2 id="admin-audit-title" className="mt-2 font-serif text-3xl">
                {text.auditEvents}
              </h2>
            </div>
            <Button
              type="button"
              variant="outline"
              disabled={loadingAudit}
              onClick={() =>
                void loadAudit({ targetType, eventType }).catch(() =>
                  toast.error(text.load),
                )
              }
            >
              {text.search}
            </Button>
          </div>
          <div className="mt-4 grid gap-2 sm:grid-cols-2">
            <input
              className="auth-input"
              aria-label={text.targetType}
              placeholder={text.targetType}
              value={targetType}
              onChange={(event) => setTargetType(event.target.value)}
            />
            <input
              className="auth-input"
              aria-label={text.eventType}
              placeholder={text.eventType}
              value={eventType}
              onChange={(event) => setEventType(event.target.value)}
            />
          </div>
          <div className="mt-5 divide-y divide-border">
            {audit.map((event) => (
              <article
                className="flex items-start justify-between gap-4 py-4 text-sm"
                key={event.id}
              >
                <div>
                  <strong>{event.eventType}</strong>
                  <p className="text-muted-foreground">
                    {event.targetType} {event.targetId ?? ""} ·{" "}
                    {formatFullDateTime(event.occurredAt, locale)}
                  </p>
                  {event.reason && <p>{event.reason}</p>}
                </div>
                {section !== "overview" && (
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setSelectedAudit(event)}
                  >
                    {text.details}
                  </Button>
                )}
              </article>
            ))}
            {!loadingAudit && audit.length === 0 && (
              <p className="py-5 text-sm text-muted-foreground">
                {text.noAudit}
              </p>
            )}
          </div>
          {section !== "overview" && (
            <div className="mt-4 flex items-center justify-between gap-3">
              <Button
                type="button"
                variant="outline"
                disabled={loadingAudit || auditPage <= 1}
                onClick={() =>
                  void loadAudit(
                    { targetType, eventType },
                    auditPage - 1,
                  ).catch(() => toast.error(text.load))
                }
              >
                Previous
              </Button>
              <span className="text-sm text-muted-foreground">
                Page {auditPage} of {Math.max(1, Math.ceil(auditTotal / 10))}
              </span>
              <Button
                type="button"
                variant="outline"
                disabled={
                  loadingAudit ||
                  auditPage >= Math.max(1, Math.ceil(auditTotal / 10))
                }
                onClick={() =>
                  void loadAudit(
                    { targetType, eventType },
                    auditPage + 1,
                  ).catch(() => toast.error(text.load))
                }
              >
                Next
              </Button>
            </div>
          )}
        </section>
      )}
      {selectedAudit && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center overflow-x-auto bg-foreground/50 p-4"
          role="dialog"
          aria-modal="true"
          aria-labelledby="audit-detail-title"
        >
          <div className="max-h-[calc(100svh-2rem)] w-full max-w-xl overflow-x-auto overflow-y-auto border border-border bg-background p-6">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-widest text-primary">
                  {text.details}
                </p>
                <h2
                  id="audit-detail-title"
                  className="mt-2 font-serif text-2xl"
                >
                  {selectedAudit.eventType}
                </h2>
              </div>
              <Button
                type="button"
                variant="ghost"
                onClick={() => setSelectedAudit(null)}
              >
                {text.close}
              </Button>
            </div>
            <dl className="mt-6 space-y-3 text-sm">
              <div>
                <dt className="text-muted-foreground">{text.targetType}</dt>
                <dd>{selectedAudit.targetType}</dd>
              </div>
              <div>
                <dt className="text-muted-foreground">ID</dt>
                <dd className="break-all">
                  {selectedAudit.targetId ?? selectedAudit.id}
                </dd>
              </div>
              <div>
                <dt className="text-muted-foreground">{text.created}</dt>
                <dd>
                  {formatFullDateTime(selectedAudit.occurredAt, locale)}
                </dd>
              </div>
              {selectedAudit.reason && (
                <div>
                  <dt className="text-muted-foreground">{text.reason}</dt>
                  <dd>{selectedAudit.reason}</dd>
                </div>
              )}
            </dl>
          </div>
        </div>
      )}
    </section>
  );
}
