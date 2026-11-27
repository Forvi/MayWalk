import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';

import { useAuth } from '../features/auth/AuthContext';

/** Гостевой маршрут: уже вошедшего участника сразу перенаправляет в кабинет. */
export function GuestRoute({ children }: { children: ReactNode }) {
  const { user } = useAuth();

  if (user) {
    return <Navigate to="/cabinet/journey" replace />;
  }

  return <>{children}</>;
}
