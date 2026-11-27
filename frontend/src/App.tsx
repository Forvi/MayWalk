import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';

import { AuthProvider } from './features/auth/AuthContext';
import { GuestRoute } from './app/GuestRoute';
import { ProtectedRoute } from './app/ProtectedRoute';
import { LandingPage } from './features/landing/LandingPage';
import { LoginPage } from './features/auth/pages/LoginPage';
import { NotFoundPage } from './app/NotFoundPage';
import { RecoveryPage } from './features/auth/pages/RecoveryPage';
import { RegisterPage } from './features/auth/pages/RegisterPage';
import { CabinetLayout } from './features/cabinet/CabinetLayout';
import { JourneyPage } from './features/cabinet/pages/JourneyPage';
import { ProfilePage } from './features/cabinet/pages/ProfilePage';

/**
 * Корневой компонент: провайдер авторизации и карта маршрутов.
 * Оформление — оригинальная концепция дизайнера (зелёная), без переключения тем.
 *
 * Публичные страницы:
 *   /                              — лендинг с формой регистрации
 *   /login, /register, /recovery   — «паспортные» страницы входа (только для гостей)
 * Кабинет (требуется авторизация):
 *   /cabinet/journey               — «Моя прогулка»: выбор маршрута, черновик анкеты
 *   /cabinet/profile               — «Профиль»: логин и почта
 */
export function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route
            path="/login"
            element={
              <GuestRoute>
                <LoginPage />
              </GuestRoute>
            }
          />
          <Route
            path="/register"
            element={
              <GuestRoute>
                <RegisterPage />
              </GuestRoute>
            }
          />
          <Route
            path="/recovery"
            element={
              <GuestRoute>
                <RecoveryPage />
              </GuestRoute>
            }
          />
          <Route
            path="/cabinet"
            element={
              <ProtectedRoute>
                <CabinetLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Navigate to="journey" replace />} />
            <Route path="journey" element={<JourneyPage />} />
            <Route path="profile" element={<ProfilePage />} />
          </Route>
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
