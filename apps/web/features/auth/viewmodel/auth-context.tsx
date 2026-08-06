"use client";

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";

import { fetchCurrentUser, revokeSession } from "../api";
import type { AuthUser } from "../types";

export type AuthStatus = "loading" | "authenticated" | "unauthenticated";

type AuthContextValue = {
  user: AuthUser | null;
  status: AuthStatus;
  refresh: () => Promise<AuthUser | null>;
  setUser: (user: AuthUser | null) => void;
  updateUser: (changes: Partial<AuthUser>) => void;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [status, setStatus] = useState<AuthStatus>("loading");

  const setAuthenticatedUser = useCallback((nextUser: AuthUser | null) => {
    setUser(nextUser);
    setStatus(nextUser ? "authenticated" : "unauthenticated");
  }, []);

  const refresh = useCallback(async () => {
    try {
      const currentUser = await fetchCurrentUser();
      setAuthenticatedUser(currentUser);
      return currentUser;
    } catch {
      setAuthenticatedUser(null);
      return null;
    }
  }, [setAuthenticatedUser]);

  useEffect(() => {
    const controller = new AbortController();

    void fetchCurrentUser(controller.signal)
      .then((currentUser) => {
        setAuthenticatedUser(currentUser);
      })
      .catch(() => {
        if (!controller.signal.aborted) {
          setAuthenticatedUser(null);
        }
      });

    return () => controller.abort();
  }, [setAuthenticatedUser]);

  const updateUser = useCallback((changes: Partial<AuthUser>) => {
    setUser((currentUser) => (currentUser ? { ...currentUser, ...changes } : currentUser));
  }, []);

  const logout = useCallback(async () => {
    try {
      await revokeSession();
    } catch {
      // The local state must still be cleared when the network is unavailable.
    } finally {
      setAuthenticatedUser(null);
    }
  }, [setAuthenticatedUser]);

  const value = useMemo(
    () => ({ user, status, refresh, setUser: setAuthenticatedUser, updateUser, logout }),
    [logout, refresh, setAuthenticatedUser, status, updateUser, user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) throw new Error("useAuth must be used inside AuthProvider");
  return value;
}

export function useOptionalAuth() {
  return useContext(AuthContext);
}
