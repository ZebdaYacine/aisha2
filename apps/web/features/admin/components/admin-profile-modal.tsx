/* eslint-disable react-hooks/set-state-in-effect */
"use client";

import { useEffect, useState } from "react";
import { toast } from "sonner";

import { Button } from "@/core/components/ui/button";
import { submitJSON } from "@/core/forms/api";
import { formCopy } from "@/core/lib/form-copy";
import type { Locale } from "@/core/lib/i18n";
import { useEscapeKey } from "@/core/hooks/use-escape-key";
import { useOptionalAuth } from "@/features/auth/viewmodel/auth-context";

type Profile = {
  id?: string;
  displayName: string;
  email: string;
  phone?: string;
  roles?: string[];
};

type Copy = {
  title: string;
  personal: string;
  name: string;
  email: string;
  phone: string;
  password: string;
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
  save: string;
  changePassword: string;
  close: string;
  loading: string;
  saved: string;
  passwordChanged: string;
  mismatch: string;
  unavailable: string;
  adminOnly: string;
};

const copy: Record<Locale, Copy> = {
  en: {
    title: "Manage administrator profile",
    personal: "Personal information",
    name: "Name",
    email: "Email",
    phone: "Phone",
    password: "Change password",
    currentPassword: "Current password",
    newPassword: "New password",
    confirmPassword: "Confirm new password",
    save: "Save profile",
    changePassword: "Change password",
    close: "Close",
    loading: "Loading profile…",
    saved: "Profile updated",
    passwordChanged: "Password changed. Sign in again if this session expires.",
    mismatch: "Passwords do not match.",
    unavailable: "Could not save your changes.",
    adminOnly: "Administrator account",
  },
  fr: {
    title: "Gérer le profil administrateur",
    personal: "Informations personnelles",
    name: "Nom",
    email: "E-mail",
    phone: "Téléphone",
    password: "Modifier le mot de passe",
    currentPassword: "Mot de passe actuel",
    newPassword: "Nouveau mot de passe",
    confirmPassword: "Confirmer le nouveau mot de passe",
    save: "Enregistrer le profil",
    changePassword: "Modifier le mot de passe",
    close: "Fermer",
    loading: "Chargement du profil…",
    saved: "Profil mis à jour",
    passwordChanged:
      "Mot de passe modifié. Reconnectez-vous si cette session expire.",
    mismatch: "Les mots de passe ne correspondent pas.",
    unavailable: "Impossible d’enregistrer les modifications.",
    adminOnly: "Compte administrateur",
  },
  ar: {
    title: "إدارة ملف المدير",
    personal: "المعلومات الشخصية",
    name: "الاسم",
    email: "البريد الإلكتروني",
    phone: "الهاتف",
    password: "تغيير كلمة المرور",
    currentPassword: "كلمة المرور الحالية",
    newPassword: "كلمة المرور الجديدة",
    confirmPassword: "تأكيد كلمة المرور الجديدة",
    save: "حفظ الملف",
    changePassword: "تغيير كلمة المرور",
    close: "إغلاق",
    loading: "جارٍ تحميل الملف…",
    saved: "تم تحديث الملف",
    passwordChanged:
      "تم تغيير كلمة المرور. سجّل الدخول مجدداً إذا انتهت هذه الجلسة.",
    mismatch: "كلمتا المرور غير متطابقتين.",
    unavailable: "تعذر حفظ التغييرات.",
    adminOnly: "حساب مدير",
  },
  es: {
    title: "Gestionar perfil de administrador",
    personal: "Información personal",
    name: "Nombre",
    email: "Correo electrónico",
    phone: "Teléfono",
    password: "Cambiar contraseña",
    currentPassword: "Contraseña actual",
    newPassword: "Nueva contraseña",
    confirmPassword: "Confirmar nueva contraseña",
    save: "Guardar perfil",
    changePassword: "Cambiar contraseña",
    close: "Cerrar",
    loading: "Cargando perfil…",
    saved: "Perfil actualizado",
    passwordChanged:
      "Contraseña cambiada. Vuelve a iniciar sesión si esta sesión caduca.",
    mismatch: "Las contraseñas no coinciden.",
    unavailable: "No se pudieron guardar los cambios.",
    adminOnly: "Cuenta de administrador",
  },
};

export function AdminProfileModal({
  locale,
  open,
  onClose,
}: {
  locale: Locale;
  open: boolean;
  onClose: () => void;
}) {
  useEscapeKey(onClose, open);
  const text = copy[locale];
  const messages = formCopy(locale);
  const auth = useOptionalAuth();
  const [profile, setProfile] = useState<Profile>({
    displayName: "",
    email: "",
    phone: "",
  });
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [passwordSaving, setPasswordSaving] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!open) return;
    let active = true;
    setLoading(true);
    setError("");
    void fetch("/api/customer/profile", { cache: "no-store" })
      .then(async (response) => {
        if (!response.ok) throw new Error();
        const data = (await response.json()) as Profile;
        if (active) setProfile({ ...data, phone: data.phone ?? "" });
      })
      .catch(() => active && setError(text.unavailable))
      .finally(() => active && setLoading(false));
    return () => {
      active = false;
    };
  }, [open, text.unavailable]);

  if (!open) return null;

  const saveProfile = async (event: React.FormEvent) => {
    event.preventDefault();
    setSaving(true);
    setError("");
    const result = await submitJSON(
      "/api/customer/profile",
      { displayName: profile.displayName, phone: profile.phone ?? "" },
      "PATCH",
    );
    if (result.ok) {
      auth?.updateUser({ displayName: profile.displayName });
      toast.success(text.saved);
    } else {
      setError(
        result.code === "SERVICE_UNAVAILABLE"
          ? messages.unavailable
          : text.unavailable,
      );
    }
    setSaving(false);
  };

  const changePassword = async (event: React.FormEvent) => {
    event.preventDefault();
    if (newPassword !== confirmPassword) {
      setError(text.mismatch);
      return;
    }
    if (newPassword.length < 12) {
      setError(messages.password);
      return;
    }
    setPasswordSaving(true);
    setError("");
    const result = await submitJSON(
      "/api/customer/password",
      { currentPassword, newPassword },
      "PATCH",
    );
    if (result.ok) {
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
      toast.success(text.passwordChanged);
    } else {
      setError(
        result.code === "SERVICE_UNAVAILABLE"
          ? messages.unavailable
          : text.unavailable,
      );
    }
    setPasswordSaving(false);
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center overflow-x-auto bg-foreground/50 p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="admin-profile-title"
    >
      <div className="max-h-[calc(100svh-2rem)] w-full max-w-2xl overflow-x-auto overflow-y-auto border border-border bg-background p-6 shadow-xl">
        <div className="flex items-start justify-between gap-4 border-b border-border pb-5">
          <div>
            <p className="text-xs uppercase tracking-widest text-primary">
              {text.adminOnly}
            </p>
            <h2 id="admin-profile-title" className="mt-2 font-serif text-3xl">
              {text.title}
            </h2>
          </div>
          <Button type="button" variant="ghost" onClick={onClose}>
            {text.close}
          </Button>
        </div>
        {loading ? (
          <p className="py-8 text-sm text-muted-foreground" role="status">
            {text.loading}
          </p>
        ) : (
          <>
            {error && (
              <p
                className="mt-5 border border-destructive p-3 text-sm text-destructive"
                role="alert"
              >
                {error}
              </p>
            )}
            <form
              className="mt-6 space-y-4"
              onSubmit={(event) => void saveProfile(event)}
            >
              <h3 className="font-serif text-2xl">{text.personal}</h3>
              <Field label={text.name}>
                <input
                  className="auth-input"
                  value={profile.displayName}
                  onChange={(event) =>
                    setProfile({ ...profile, displayName: event.target.value })
                  }
                  required
                  minLength={2}
                />
              </Field>
              <Field label={text.email}>
                <input className="auth-input" value={profile.email} readOnly />
              </Field>
              <Field label={text.phone}>
                <input
                  className="auth-input"
                  type="tel"
                  value={profile.phone ?? ""}
                  onChange={(event) =>
                    setProfile({ ...profile, phone: event.target.value })
                  }
                />
              </Field>
              <Button type="submit" disabled={saving}>
                {text.save}
              </Button>
            </form>
            <form
              className="mt-10 space-y-4 border-t border-border pt-6"
              onSubmit={(event) => void changePassword(event)}
            >
              <h3 className="font-serif text-2xl">{text.password}</h3>
              <Field label={text.currentPassword}>
                <input
                  className="auth-input"
                  type="password"
                  value={currentPassword}
                  onChange={(event) => setCurrentPassword(event.target.value)}
                  required
                />
              </Field>
              <Field label={text.newPassword}>
                <input
                  className="auth-input"
                  type="password"
                  value={newPassword}
                  onChange={(event) => setNewPassword(event.target.value)}
                  required
                  minLength={12}
                />
              </Field>
              <Field label={text.confirmPassword}>
                <input
                  className="auth-input"
                  type="password"
                  value={confirmPassword}
                  onChange={(event) => setConfirmPassword(event.target.value)}
                  required
                  minLength={12}
                />
              </Field>
              <Button type="submit" variant="outline" disabled={passwordSaving}>
                {text.changePassword}
              </Button>
            </form>
          </>
        )}
      </div>
    </div>
  );
}

function Field({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <label className="block">
      <span className="mb-2 block text-sm">{label}</span>
      {children}
    </label>
  );
}
