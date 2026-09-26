"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { LocalizedLink } from "@/core/components/shared/localized-link";
import { Button } from "@/core/components/ui/button";
import type { Locale } from "@/core/lib/i18n";
import { hasCapability } from "@/features/auth/types";
import { useAuth } from "@/features/auth/viewmodel/auth-context";
import { AdminProfileModal } from "./admin-profile-modal";

type Section =
  | "overview"
  | "users"
  | "artisan-applications"
  | "moderation"
  | "warehouse"
  | "inventory"
  | "media"
  | "audit";
const labels: Record<
  Locale,
  {
    title: string;
    profile: string;
    logout: string;
    overview: string;
    users: string;
    applications: string;
    moderation: string;
    warehouse: string;
    inventory: string;
    media: string;
    audit: string;
    loggingOut: string;
  }
> = {
  en: {
    title: "Administration dashboard",
    profile: "Manage profile",
    logout: "Sign out",
    overview: "Overview",
    users: "Users",
    applications: "Artisan applications",
    moderation: "Product moderation",
    warehouse: "Warehouse",
    inventory: "Inventory",
    media: "Media",
    audit: "Event audit",
    loggingOut: "Signing out…",
  },
  fr: {
    title: "Tableau de bord d’administration",
    profile: "Gérer le profil",
    logout: "Déconnexion",
    overview: "Vue d’ensemble",
    users: "Utilisateurs",
    applications: "Candidatures d’artisans",
    moderation: "Modération des produits",
    warehouse: "Entrepôt",
    inventory: "Inventaire",
    media: "Médias",
    audit: "Audit des événements",
    loggingOut: "Déconnexion…",
  },
  ar: {
    title: "لوحة تحكم الإدارة",
    profile: "إدارة الملف الشخصي",
    logout: "تسجيل الخروج",
    overview: "نظرة عامة",
    users: "المستخدمون",
    applications: "طلبات الحرفيين",
    moderation: "مراجعة المنتجات",
    warehouse: "المستودع",
    inventory: "المخزون",
    media: "الوسائط",
    audit: "تدقيق الأحداث",
    loggingOut: "جارٍ تسجيل الخروج…",
  },
  es: {
    title: "Panel de administración",
    profile: "Gestionar perfil",
    logout: "Cerrar sesión",
    overview: "Resumen",
    users: "Usuarios",
    applications: "Solicitudes de artesanos",
    moderation: "Moderación de productos",
    warehouse: "Almacén",
    inventory: "Inventario",
    media: "Medios",
    audit: "Auditoría de eventos",
    loggingOut: "Cerrando sesión…",
  },
};

export function AdminDashboardShell({
  locale,
  section,
  children,
}: {
  locale: Locale;
  section: Section;
  children: React.ReactNode;
}) {
  const text = labels[locale];
  const router = useRouter();
  const { user, logout } = useAuth();
  const [profileOpen, setProfileOpen] = useState(false);
  const [loggingOut, setLoggingOut] = useState(false);
  const links: [Section, string, string][] = [];
  if (hasCapability(user, "admin.audit.read")) links.push(["overview", "/admin", text.overview]);
  if (hasCapability(user, "admin.users.read")) links.push(["users", "/admin/users", text.users]);
  if (hasCapability(user, "admin.artisan_applications.read")) links.push(["artisan-applications", "/admin/artisan-applications", text.applications]);
  if (hasCapability(user, "admin.product_moderation.read")) links.push(["moderation", "/admin/moderation", text.moderation]);
  if (hasCapability(user, "warehouse.read")) links.push(["warehouse", "/admin/warehouse", text.warehouse]);
  if (hasCapability(user, "inventory.read")) links.push(["inventory", "/admin/inventory", text.inventory]);
  if (hasCapability(user, "admin.media.read")) links.push(["media", "/admin/media", text.media]);
  if (hasCapability(user, "admin.audit.read")) links.push(["audit", "/admin/audit", text.audit]);
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
  return (
    <div className="space-y-8">
      <header className="flex flex-col gap-5 border-b border-border pb-8 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.18em] text-primary">
            {text.title}
          </p>
          <h1 className="mt-3 font-serif text-5xl">
            {
              text[
                section === "artisan-applications" ? "applications" : section
              ]
            }
          </h1>
          <p className="mt-3 text-sm text-muted-foreground">
            {user?.displayName} · {user?.email}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => setProfileOpen(true)}
          >
            {text.profile}
          </Button>
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
      <nav className="flex flex-wrap gap-2" aria-label={text.title}>
        {links.map(([key, href, label]) => (
          <LocalizedLink
            key={key}
            locale={locale}
            href={href}
            aria-current={section === key ? "page" : undefined}
            className={`border px-4 py-3 text-sm ${section === key ? "border-foreground bg-foreground text-background" : "border-border hover:bg-muted"}`}
          >
            {label}
          </LocalizedLink>
        ))}
      </nav>
      <main>{children}</main>
      <AdminProfileModal
        locale={locale}
        open={profileOpen}
        onClose={() => setProfileOpen(false)}
      />
    </div>
  );
}
