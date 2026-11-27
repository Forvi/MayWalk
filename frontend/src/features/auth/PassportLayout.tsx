import { NavLink } from 'react-router-dom';

import type { ReactNode } from 'react';

import { DialogCloseButton, ModalDialog } from '../../shared/components/ModalDialog';
import { PageShell } from '../../shared/components/PageShell';

type PassportLayoutProps = {
  children: ReactNode;
  /** Заголовок вкладки, которая подсвечена в табах (login/register). */
  activeTab: 'login' | 'register';
  vkDialogOpen: boolean;
  onCloseVkDialog: () => void;
};

/**
 * Постерная графика «маршрута»: в оригинальной концепции — фоновые PNG
 * из дизайн-макета (задаются стилями .poster-art). SVG-карта относилась
 * к авторской концепции и удалена вместе с ней.
 */
function PosterArt() {
  return <div className="poster-art" aria-hidden="true" />;
}

/**
 * «Паспортная» страница входа: шапка с брендом, постер, бланк-паспорт с формой.
 * Разметка из auth-original-*.html дизайн-макета (оригинальная концепция).
 */
export function PassportLayout({
  children,
  activeTab,
  vkDialogOpen,
  onCloseVkDialog,
}: PassportLayoutProps) {
  return (
    <PageShell page="auth">
      <header className="site-header shell">
        <NavLink className="brand" to="/" aria-label="Майская прогулка — главная">
          <span className="brand-symbol">
            <svg viewBox="0 0 48 48" fill="none" aria-hidden="true">
              <path
                d="M7 37c1-11 16-8 18-19 2-9 11-9 15-13M17 43c0-7 9-8 12-12"
                stroke="currentColor"
                strokeWidth="5"
                strokeLinecap="round"
              />
              <circle cx="9" cy="12" r="6" fill="currentColor" />
            </svg>
          </span>
          <span className="brand-copy">
            майская
            <br />
            прогулка.
          </span>
        </NavLink>
        <nav className="primary-nav" aria-label="Основная навигация">
          <a href="https://mayprogulka.ru/">Главная</a>
          <a href="https://mayprogulka.ru/about/uchast/">О прогулке</a>
          <a href="https://mayprogulka.ru/about/distance/">Маршруты</a>
        </nav>
        <a className="help-link" href="https://mayprogulka.ru/faq">
          Нужна помощь? <span className="help-icon" aria-hidden="true">?</span>
        </a>
      </header>

      <main className="access-stage shell">
        <section className="poster" aria-label="Майская прогулка в Екатеринбурге">
          <h2 className="poster-title">
            <span>Пора</span>
            <span>гулять</span>
          </h2>
          <p className="poster-description">
            Лесные тропы, городские маршруты и день, который хочется провести на ногах.
          </p>
          <PosterArt />
          <div className="poster-note">
            <span>
              <strong>У каждого свой темп.</strong>
              <br />
              У всех — своя Майская прогулка.
            </span>
          </div>
        </section>

        <section className="passport" aria-labelledby="form-title">
          <div className="auth-tabs" role="tablist" aria-label="Вход или регистрация">
            <NavLink
              to="/login"
              id="tab-login"
              role="tab"
              tabIndex={activeTab === 'login' ? 0 : -1}
              className={`auth-tab${activeTab === 'login' ? ' active' : ''}`}
              aria-selected={activeTab === 'login'}
              aria-controls="auth-content"
            >
              Вход
            </NavLink>
            <NavLink
              to="/register"
              id="tab-register"
              role="tab"
              className={`auth-tab${activeTab === 'register' ? ' active' : ''}`}
              aria-selected={activeTab === 'register'}
              aria-controls="auth-content"
              tabIndex={activeTab === 'register' ? 0 : -1}
            >
              Регистрация
            </NavLink>
          </div>
          <div id="auth-content" data-auth-panel="" role="tabpanel" aria-labelledby={`tab-${activeTab}`}>
            {children}
          </div>
        </section>
      </main>

      <footer className="site-footer shell">
        <p>© Майская прогулка, 2026</p>
        <div>
          <a href="https://mayprogulka.ru/about/pdpol/" target="_blank" rel="noopener noreferrer">
            Политика конфиденциальности
          </a>
          <a href="https://mayprogulka.ru/about/pdsogl/" target="_blank" rel="noopener noreferrer">
            Обработка персональных данных
          </a>
        </div>
        <span>Встретимся на старте.</span>
      </footer>

      <ModalDialog
        className="vk-dialog"
        open={vkDialogOpen}
        onClose={onCloseVkDialog}
        titleId="vk-dialog-title"
        descriptionId="vk-dialog-description"
      >
        <DialogCloseButton label="Закрыть окно" onClick={onCloseVkDialog} />
        <h2 id="vk-dialog-title">Вход через ВКонтакте</h2>
        <p id="vk-dialog-description">
          Для входа в личный кабинет используйте свой аккаунт ВКонтакте.
        </p>
        <button type="button" className="primary-button" onClick={onCloseVkDialog}>
          Понятно
        </button>
      </ModalDialog>
    </PageShell>
  );
}

/** Кнопка «Войти через ВК» — общий элемент паспортных страниц. */
export function VkButton({ onClick }: { onClick: () => void }) {
  return (
    <button className="vk-button" type="button" onClick={onClick}>
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path
          fill="currentColor"
          d="M12.63 17.7C6.8 17.7 3.47 13.7 3.33 7.03h2.92c.1 4.9 2.26 6.98 3.97 7.41V7.03h2.75v4.23c1.68-.18 3.44-2.11 4.04-4.23h2.75c-.46 2.61-2.36 4.54-3.71 5.33 1.35.65 3.5 2.34 4.32 5.34h-3.03c-.65-2-2.26-3.54-4.37-3.75v3.75z"
        />
      </svg>
      Войти через ВК
    </button>
  );
}
