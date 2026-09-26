/* eslint-disable @typescript-eslint/no-explicit-any */
"use client";

import { useEffect, useState } from "react";
import { Files, PackageOpen, Store } from "lucide-react";

import { Skeleton } from "@/core/components/ui/skeleton";
import { StatusBadge } from "@/core/components/ui/status-badge";
import { ApplicationForm } from "./application-form";
import { ArtisanMediaPanel } from "@/features/artisan/components/artisan-media-panel";
import { ProductWorkspace } from "@/features/product";
import { WorkshopsPanel } from "@/features/artisan/components/workshops-panel";
import { hasCapability } from "@/features/auth/types";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

const labels = {
  en: { heading: "Artisan workspace", status: "Application status", none: "Start your artisan application", suspended: "Seller actions are disabled while artisan membership is suspended.", workshops: "Workshops", workshopsDescription: "Manage your production locations and workshop visibility.", privateFiles: "Private files", privateFilesDescription: "Manage verification documents and profile media.", productAuthoring: "Product authoring", productAuthoringDescription: "Create, edit, submit, and archive products." },
  fr: { heading: "Espace artisan", status: "Statut de la candidature", none: "Commencez votre candidature artisan", suspended: "Les actions de vendeur sont désactivées pendant la suspension de l’adhésion artisan.", workshops: "Ateliers", workshopsDescription: "Gérez vos lieux de production et leur visibilité.", privateFiles: "Fichiers privés", privateFilesDescription: "Gérez vos documents de vérification et vos médias.", productAuthoring: "Création de produits", productAuthoringDescription: "Créez, modifiez, soumettez et archivez vos produits." },
  ar: { heading: "مساحة الحرفي", status: "حالة الطلب", none: "ابدأ طلب الانضمام كحرفي", suspended: "تم تعطيل إجراءات البيع أثناء تعليق عضوية الحرفي.", workshops: "الورشات", workshopsDescription: "إدارة أماكن الإنتاج وظهورها.", privateFiles: "الملفات الخاصة", privateFilesDescription: "إدارة وثائق التحقق ووسائط الملف الشخصي.", productAuthoring: "إنشاء المنتجات", productAuthoringDescription: "إنشاء المنتجات وتعديلها وإرسالها وأرشفتها." },
  es: { heading: "Espacio de artesano", status: "Estado de la solicitud", none: "Inicia tu solicitud de artesano", suspended: "Las acciones de vendedor están desactivadas mientras la membresía artesanal está suspendida.", workshops: "Talleres", workshopsDescription: "Gestiona tus lugares de producción y su visibilidad.", privateFiles: "Archivos privados", privateFilesDescription: "Gestiona documentos de verificación y medios del perfil.", productAuthoring: "Creación de productos", productAuthoringDescription: "Crea, edita, envía y archiva productos." },
} as const;

type WorkspaceTab = "workshops" | "private-files" | "products";

export function ArtisanWorkspace({ locale }: { locale: keyof typeof labels }) {
  const auth = useOptionalAuth();
  const [profile, setProfile] = useState<any>();
  const [loaded, setLoaded] = useState(false);
  const [activeTab, setActiveTab] = useState<WorkspaceTab>("workshops");
  const t = labels[locale];
  const sellerEnabled =
    auth?.user?.artisanEnabled === true && hasCapability(auth.user, "artisan.account.read");

  useEffect(() => {
    void fetch("/api/artisan/application")
      .then(async (response) => {
        if (response.ok) setProfile(await response.json());
      })
      .finally(() => setLoaded(true));
  }, []);

  if (!loaded) {
    return (
      <div className="space-y-4" aria-busy="true">
        <Skeleton className="h-12 w-72" />
        <Skeleton className="h-80 w-full" />
      </div>
    );
  }

  const suspended = profile?.status === "SUSPENDED" || auth?.user?.artisanStatus === "SUSPENDED";
  const canEditApplication =
    !suspended &&
    (!profile || ["DRAFT", "CHANGES_REQUESTED", "REJECTED", "APPROVED"].includes(profile.status));

  return (
    <>
      <header className="border-b border-border pb-6">
        <p className="text-xs uppercase tracking-widest text-primary">AISHA</p>
        <h1 className="mt-3 font-serif text-4xl sm:text-5xl">{t.heading}</h1>
        {profile && (
          <p className="mt-4">
            <strong>{t.status}:</strong> <StatusBadge status={profile.status} />
          </p>
        )}
        {profile?.reviewReason && <p className="mt-2 border-s-2 border-primary ps-3">{profile.reviewReason}</p>}
        {suspended && <p className="mt-4 text-sm text-muted-foreground" role="status">{t.suspended}</p>}
      </header>

      {canEditApplication && profile?.status !== "APPROVED" && <div className="mt-8 max-w-3xl" id="workshops">
        <h2 className="mb-6 font-serif text-2xl">{profile ? t.status : t.none}</h2>
        <ApplicationForm locale={locale} profile={profile} />
      </div>}

      {profile?.status === "APPROVED" && (
        <>
          <div className="mt-8 min-w-0 overflow-hidden rounded-2xl border border-border bg-card shadow-sm" id="workshops">
            <div className="border-b border-border bg-muted/20 p-2" role="tablist" aria-label={t.heading}>
              <div className="grid gap-2 md:grid-cols-3">
                <WorkspaceTabButton icon={<Store aria-hidden="true" size={18} />} description={t.workshopsDescription} active={activeTab === "workshops"} id="workshops-tab" controls="workshops-panel" onClick={() => setActiveTab("workshops")}>{t.workshops}</WorkspaceTabButton>
                {profile && !suspended && sellerEnabled && <WorkspaceTabButton icon={<Files aria-hidden="true" size={18} />} description={t.privateFilesDescription} active={activeTab === "private-files"} id="private-files-tab" controls="private-files-panel" onClick={() => setActiveTab("private-files")}>{t.privateFiles}</WorkspaceTabButton>}
                {!suspended && sellerEnabled && <WorkspaceTabButton icon={<PackageOpen aria-hidden="true" size={18} />} description={t.productAuthoringDescription} active={activeTab === "products"} id="products-tab" controls="products-panel" onClick={() => setActiveTab("products")}>{t.productAuthoring}</WorkspaceTabButton>}
              </div>
            </div>
            <div className="min-w-0 p-1 sm:p-3">
              {activeTab === "workshops" && <div id="workshops-panel" role="tabpanel" aria-labelledby="workshops-tab"><WorkshopsPanel locale={locale} profile={profile} /></div>}
              {activeTab === "private-files" && !suspended && sellerEnabled && <div id="private-files-panel" role="tabpanel" aria-labelledby="private-files-tab"><ArtisanMediaPanel locale={locale} /></div>}
              {activeTab === "products" && !suspended && sellerEnabled && <div id="products-panel" role="tabpanel" aria-labelledby="products-tab"><ProductWorkspace locale={locale} /></div>}
            </div>
          </div>
        </>
      )}
    </>
  );
}

function WorkspaceTabButton({ active, id, controls, onClick, icon, description, children }: { active: boolean; id: string; controls: string; onClick: () => void; icon: React.ReactNode; description: string; children: React.ReactNode }) {
  return <button type="button" role="tab" id={id} aria-controls={controls} aria-selected={active} onClick={onClick} className={`flex min-h-20 items-start gap-3 rounded-xl border p-4 text-start transition-colors ${active ? "border-foreground bg-background text-foreground shadow-sm" : "border-transparent text-muted-foreground hover:border-border hover:bg-background/70 hover:text-foreground"}`}><span className={`mt-0.5 grid size-9 shrink-0 place-items-center rounded-full ${active ? "bg-foreground text-background" : "bg-muted text-foreground"}`}>{icon}</span><span><span className="block text-sm font-semibold">{children}</span><span className="mt-1 block text-xs leading-5 text-muted-foreground">{description}</span></span></button>;
}
