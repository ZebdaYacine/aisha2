"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import type { Locale } from "@/core/lib/i18n";

type User = { id: string; email: string; displayName: string; status: string; roles: string[] };
type AuditEvent = { id: string; eventType: string; actorUserId?: string; targetType: string; targetId?: string; occurredAt: string };
const labels = { en: { users: "Users and roles", audit: "Audit events", save: "Save roles", load: "Unable to load administration data", saved: "Roles updated" }, fr: { users: "Utilisateurs et rôles", audit: "Événements d’audit", save: "Enregistrer les rôles", load: "Impossible de charger les données d’administration", saved: "Rôles mis à jour" }, ar: { users: "المستخدمون والأدوار", audit: "أحداث التدقيق", save: "حفظ الأدوار", load: "تعذر تحميل بيانات الإدارة", saved: "تم تحديث الأدوار" }, es: { users: "Usuarios y roles", audit: "Eventos de auditoría", save: "Guardar roles", load: "No se pudieron cargar los datos de administración", saved: "Roles actualizados" } } as const;

export function AdminControlPanel({ locale }: { locale: Locale }) {
  const text = labels[locale];
  const [users, setUsers] = useState<User[]>([]);
  const [audit, setAudit] = useState<AuditEvent[]>([]);
  const [roles, setRoles] = useState<Record<string, string>>({});
  useEffect(() => { void Promise.all([fetch("/api/admin/users?pageSize=100"), fetch("/api/admin/audit-events?pageSize=20")]).then(async ([usersResponse, auditResponse]) => { if (!usersResponse.ok || !auditResponse.ok) throw new Error(); const userData = await usersResponse.json(); const auditData = await auditResponse.json(); setUsers(userData.items ?? []); setAudit(auditData.items ?? []); }).catch(() => toast.error(text.load)); }, [text.load]);
  const updateRoles = async (id: string) => { const value = roles[id]; if (!value) return; const response = await fetch(`/api/admin/users/${encodeURIComponent(id)}/roles`, { method: "PATCH", headers: { "content-type": "application/json" }, body: JSON.stringify({ roles: value.split(",").map((role) => role.trim()).filter(Boolean) }) }); if (!response.ok) { toast.error(text.load); return; } const user = await response.json(); setUsers((current) => current.map((item) => item.id === user.id ? user : item)); toast.success(text.saved); };
  return <section className="mt-12 border-t border-border pt-8"><div className="grid gap-10 xl:grid-cols-2"><div><h2 className="font-serif text-3xl">{text.users}</h2><div className="mt-5 divide-y divide-border">{users.map((user) => <article className="py-4" key={user.id}><div className="flex flex-wrap justify-between gap-3"><div><strong>{user.displayName || user.email}</strong><p className="text-sm text-muted-foreground">{user.email} · {user.status}</p></div><Button type="button" variant="outline" onClick={() => void updateRoles(user.id)}>{text.save}</Button></div><input className="auth-input mt-3" aria-label={user.email} value={roles[user.id] ?? user.roles.join(", ")} onChange={(event) => setRoles((current) => ({ ...current, [user.id]: event.target.value }))}/></article>)}</div></div><div><h2 className="font-serif text-3xl">{text.audit}</h2><div className="mt-5 divide-y divide-border">{audit.map((event) => <article className="py-4 text-sm" key={event.id}><strong>{event.eventType}</strong><p className="text-muted-foreground">{event.targetType} {event.targetId ?? ""} · {new Date(event.occurredAt).toLocaleString(locale)}</p></article>)}</div></div></div></section>;
}
