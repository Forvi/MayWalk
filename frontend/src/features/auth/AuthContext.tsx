import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import type { ReactNode } from 'react';

import { isValidEmail } from './validation';

export type Profile = {
  login: string;
  email: string;
};

type StoredUser = Profile & {
  password: string;
  /** Идентификатор выбранного маршрута (r15/r18/r30plus3/r50) или null. */
  routeId: string | null;
};

const USERS_KEY = 'mp:users';
const SESSION_KEY = 'mp:session';

const DEFAULT_ROUTE_ID: string = 'r18';

function readJson<T>(key: string, fallback: T): T {
  try {
    const raw = window.localStorage.getItem(key);
    if (!raw) {
      return fallback;
    }
    return JSON.parse(raw) as T;
  } catch {
    return fallback;
  }
}

function writeJson(key: string, value: unknown): void {
  try {
    window.localStorage.setItem(key, JSON.stringify(value));
  } catch {
    // приватный режим браузера: данные проживут только до перезагрузки
  }
}

function sanitizeLogin(login: string): string {
  return login.trim().toLowerCase();
}

function loadUsers(): StoredUser[] {
  return readJson<StoredUser[]>(USERS_KEY, []);
}

function saveUsers(users: StoredUser[]): void {
  writeJson(USERS_KEY, users);
}

export type RegisterResult = { ok: true; user: Profile } | { ok: false; error: string };
export type LoginResult = { ok: true; user: Profile } | { ok: false; error: string };
export type UpdateProfileResult = { ok: true; user: Profile } | { ok: false; error: string };

type AuthContextValue = {
  /** Текущий участник или null, если вход не выполнен. */
  user: Profile | null;
  /** Идентификатор выбранного маршрута текущего участника (существование проверяется на месте). */
  routeId: string | null;
  register: (input: { login: string; email: string; password: string }) => RegisterResult;
  login: (input: { login: string; password: string }) => LoginResult;
  /** Демо-восстановление доступа: всегда сообщает об отправке письма. */
  recoverAccess: (input: { email: string }) => { message: string };
  logout: () => void;
  setRoute: (routeId: string | null) => void;
  updateProfile: (input: { login: string; email: string }) => UpdateProfileResult;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<Profile | null>(() => {
    const sessionLogin = readJson<string | null>(SESSION_KEY, null);
    if (!sessionLogin) {
      return null;
    }
    const found = loadUsers().find((candidate) => candidate.login === sessionLogin);
    return found ? { login: found.login, email: found.email } : null;
  });

  const [routeId, setRouteIdState] = useState<string | null>(() => {
    const sessionLogin = readJson<string | null>(SESSION_KEY, null);
    if (!sessionLogin) {
      return null;
    }
    const found = loadUsers().find((candidate) => candidate.login === sessionLogin);
    return found?.routeId ?? null;
  });

  useEffect(() => {
    document.title = user
      ? `${user.login} — Личный кабинет · Майская прогулка`
      : 'Майская прогулка';
  }, [user]);

  const register = useCallback(
    ({ login, email, password }: { login: string; email: string; password: string }): RegisterResult => {
      const normalizedLogin = sanitizeLogin(login);
      if (!normalizedLogin) {
        return { ok: false, error: 'Введите логин.' };
      }
      const normalizedEmail = email.trim();
      if (!isValidEmail(normalizedEmail)) {
        return { ok: false, error: 'Проверьте адрес электронной почты.' };
      }
      const users = loadUsers();
      if (users.some((candidate) => candidate.login === normalizedLogin)) {
        return { ok: false, error: 'Логин уже занят.' };
      }
      if (users.some((candidate) => candidate.email === normalizedEmail)) {
        return { ok: false, error: 'Почта уже используется.' };
      }
      const created: StoredUser = {
        login: normalizedLogin,
        email: normalizedEmail,
        password,
        routeId: DEFAULT_ROUTE_ID,
      };
      saveUsers([...users, created]);
      writeJson(SESSION_KEY, created.login);
      setUser({ login: created.login, email: created.email });
      setRouteIdState(created.routeId);
      return {
        ok: true,
        user: { login: created.login, email: created.email },
      };
    },
    [],
  );

  const login = useCallback(
    ({ login, password }: { login: string; password: string }): LoginResult => {
      const normalizedLogin = sanitizeLogin(login);
      const users = loadUsers();
      const found = users.find((candidate) => candidate.login === normalizedLogin);
      if (!found || found.password !== password) {
        return { ok: false, error: 'Неверный логин или пароль.' };
      }
      writeJson(SESSION_KEY, found.login);
      setUser({ login: found.login, email: found.email });
      setRouteIdState(found.routeId);
      return { ok: true, user: { login: found.login, email: found.email } };
    },
    [],
  );

  const recoverAccess = useCallback(({ email }: { email: string }) => {
    void email;
    // Демо-режим без бэкенда: адрес проверяется формой, письмо «отправляется» всегда.
    return { message: 'Письмо с инструкцией отправлено. Проверьте почту.' };
  }, []);

  const logout = useCallback(() => {
    try {
      window.localStorage.removeItem(SESSION_KEY);
    } catch {
      // localStorage может быть недоступен
    }
    setUser(null);
    setRouteIdState(null);
  }, []);

  const setRoute = useCallback((nextRouteId: string | null) => {
    setRouteIdState(nextRouteId);
    const sessionLogin = readJson<string | null>(SESSION_KEY, null);
    if (!sessionLogin) {
      return;
    }
    const users = loadUsers().map((candidate) =>
      candidate.login === sessionLogin ? { ...candidate, routeId: nextRouteId } : candidate,
    );
    saveUsers(users);
  }, []);

  const updateProfile = useCallback(
    ({ login, email }: { login: string; email: string }): UpdateProfileResult => {
      const normalizedLogin = sanitizeLogin(login);
      const normalizedEmail = email.trim();
      if (!user) {
        return { ok: false, error: 'Сначала войдите в личный кабинет.' };
      }
      if (!isValidEmail(normalizedEmail)) {
        return { ok: false, error: 'Введите корректный адрес электронной почты.' };
      }
      const users = loadUsers();
      if (users.some((candidate) => candidate.login === normalizedLogin && candidate.login !== user.login)) {
        return { ok: false, error: 'Логин уже занят.' };
      }
      if (users.some((candidate) => candidate.email === normalizedEmail && candidate.login !== user.login)) {
        return { ok: false, error: 'Почта уже используется.' };
      }
      const updated = users.map((candidate) =>
        candidate.login === user.login
          ? { ...candidate, login: normalizedLogin, email: normalizedEmail }
          : candidate,
      );
      saveUsers(updated);
      writeJson(SESSION_KEY, normalizedLogin);
      setUser({ login: normalizedLogin, email: normalizedEmail });
      return { ok: true, user: { login: normalizedLogin, email: normalizedEmail } };
    },
    [user],
  );

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      routeId,
      register,
      login,
      recoverAccess,
      logout,
      setRoute,
      updateProfile,
    }),
    [user, routeId, register, login, recoverAccess, logout, setRoute, updateProfile],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth должен вызываться внутри <AuthProvider>');
  }
  return context;
}
