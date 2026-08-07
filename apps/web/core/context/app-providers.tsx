"use client";

import type { ReactNode } from "react";
import { Toaster } from "sonner";

import { CartDrawer } from "@/features/cart/components/containers/cart-drawer";
import { AuthProvider } from "@/features/auth/viewmodel/auth-context";
import { CartProvider } from "@/features/cart/viewmodel/cart-context";
import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";

export function AppProviders({ children, locale, copy }: { children: ReactNode; locale: Locale; copy: StoreCopy }) {
  return (
    <AuthProvider>
      <CartProvider>
        {children}
        <CartDrawer locale={locale} copy={copy} />
        <Toaster position="top-center" richColors closeButton />
      </CartProvider>
    </AuthProvider>
  );
}
