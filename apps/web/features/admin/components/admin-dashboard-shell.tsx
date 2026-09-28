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
  | "users"
  | "workers"
  | "artisan-applications"
  | "moderation"
  | "warehouse"
  | "inventory"
  | "media"
  | "audit"
  | "categories"
  | "orders";

const labels: Record<
  Locale,
  {
    title: string;
    profile: string;
    logout: string;
    users: string;
    workers: string;
    applications: string;
    moderation: string;
    warehouse: string;
    inventory: string;
    media: string;
    audit: string;
    categories: string;
    orders: string;
    loggingOut: string;
  }
> = {
  en: {
    title: "Administration dashboard",
    profile: "Manage profile",
    logout: "Sign out",
    users: "Users",
    workers: "Workers",
    applications: "Artisan applications",
    moderation: "Product moderation",
    warehouse: "Warehouses",
    inventory: "Inventory",
    media: "Media",
    audit: "Event audit",
    categories: "Categories",
    orders: "Orders",
    loggingOut: "Signing out…",
  },
  fr: {
    title: "Tableau de bord d’administration",
    profile: "Gérer le profil",
    logout: "Déconnexion",
    users: "Utilisateurs",
    workers: "Personnel",
    applications: "Candidatures d’artisans",
    moderation: "Modération des produits",
    warehouse: "Entrepôts",
    inventory: "Inventaire",
    media: "Médias",
    audit: "Audit des événements",
    categories: "Catégories",
    orders: "Commandes",
    loggingOut: "Déconnexion…",
  },
  ar: {
    title: "لوحة تحكم الإدارة",
    profile: "إدارة الملف الشخصي",
    logout: "تسجيل الخروج",
    users: "المستخدمون",
    workers: "العمال",
    applications: "طلبات الحرفيين",
    moderation: "مراجعة المنتجات",
    warehouse: "المستودعات",
    inventory: "المخزون",
    media: "الوسائط",
    audit: "تدقيق الأحداث",
    categories: "التصنيفات",
    orders: "الطلبات",
    loggingOut: "جارٍ تسجيل الخروج…",
  },
  es: {
    title: "Panel de administración",
    profile: "Gestionar perfil",
    logout: "Cerrar sesión",
    users: "Usuarios",
    workers: "Trabajadores",
    applications: "Solicitudes de artesanos",
    moderation: "Moderación de productos",
    warehouse: "Almacenes",
    inventory: "Inventario",
    media: "Medios",
    audit: "Auditoría de eventos",
    categories: "Categorías",
    orders: "Pedidos",
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
  if (hasCapability(user, "admin.users.read")) links.push(["users", "/admin/users", text.users]);
  if (hasCapability(user, "admin.users.read")) links.push(["workers", "/admin/workers", text.workers]);
  if (hasCapability(user, "admin.artisan_applications.read")) links.push(["artisan-applications", "/admin/artisan-applications", text.applications]);
  if (hasCapability(user, "admin.product_moderation.read")) links.push(["moderation", "/admin/moderation", text.moderation]);
  if (hasCapability(user, "warehouse.read")) links.push(["warehouse", "/admin/warehouse", text.warehouse]);
  if (hasCapability(user, "inventory.read")) links.push(["inventory", "/admin/inventory", text.inventory]);
  if (hasCapability(user, "admin.media.read")) links.push(["media", "/admin/media", text.media]);
  if (hasCapability(user, "admin.audit.read")) links.push(["audit", "/admin/audit", text.audit]);
  if (hasCapability(user, "admin.users.write")) links.push(["categories", "/admin/categories", text.categories]);
  if (hasCapability(user, "admin.users.read")) links.push(["orders", "/admin/orders", text.orders]);
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
      <nav
        className="-mx-1 flex max-w-full flex-nowrap gap-2 overflow-x-auto px-1 pb-2"
        aria-label={text.title}
        role="tablist"
      >
        {links.map(([key, href, label]) => (
          <LocalizedLink
            key={key}
            locale={locale}
            href={href}
            role="tab"
            aria-selected={section === key}
            aria-current={section === key ? "page" : undefined}
            className={`shrink-0 rounded border px-4 py-3 text-sm transition-colors ${section === key ? "border-primary bg-secondary text-foreground shadow-sm" : "border-border bg-card hover:border-primary/60 hover:bg-muted"}`}
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
