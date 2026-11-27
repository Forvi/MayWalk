import type { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';

import { useAuth } from '../features/auth/AuthContext';

/**
 * Закрытый маршрут: пускает вошедшего участника, гостя отправляет на /login,
 * сохраняя исходный адрес в state.from (после входа вернёт обратно).
 */
export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const location = useLocation();

  if (!user) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  return <>{children}</>;
}
