import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom';

import { useAuth } from '../../features/auth/AuthContext';
import { BrandMark } from '../../shared/components/BrandMark';
import { PageShell } from '../../shared/components/PageShell';
import { CABINET_COPY } from './content';

/**
 * Каркас личного кабинета: шапка, табы «Моя прогулка» / «Профиль», выход.
 * Разметка перенесена из cabinet-*.html дизайн-макета.
 */
export function CabinetLayout() {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleExit = () => {
    logout();
    navigate('/login');
  };

  return (
    <PageShell page="cabinet">
      <header className="global-header shell">
        <Link className="brand" to="/" aria-label="Майская прогулка — главная">
          <BrandMark />
          <span className="brand-name">
            майская
            <br />
            прогулка.
          </span>
        </Link>
        <div className="header-links">
          <a href={CABINET_COPY.externalLinks.faq} target="_blank" rel="noopener noreferrer">
            {CABINET_COPY.externalLinks.help}
          </a>
        </div>
      </header>

      <nav className="cabinet-nav shell" aria-label="Личный кабинет">
        <div className="cabinet-nav-tabs" role="tablist" aria-label="Раздел кабинета">
          <NavLink
            to="journey"
            id="tab-journey"
            role="tab"
            aria-controls="cabinet-content"
            className={({ isActive }) => `cabinet-tab${isActive ? ' active' : ''}`}
          >
            {CABINET_COPY.page.journeyTitle}
          </NavLink>
          <NavLink
            to="profile"
            id="tab-profile"
            role="tab"
            aria-controls="cabinet-content"
            className={({ isActive }) => `cabinet-tab${isActive ? ' active' : ''}`}
          >
            {CABINET_COPY.page.profileTitle}
          </NavLink>
        </div>
        <button type="button" className="exit-link" onClick={handleExit}>
          {CABINET_COPY.actions.exit}
        </button>
      </nav>

      <main className="shell">
        <Outlet />
      </main>

      <footer className="site-footer shell">
        <p>Майская прогулка · Екатеринбург</p>
        <div className="site-footer-links">
          <a href={CABINET_COPY.externalLinks.privacy} target="_blank" rel="noopener noreferrer">
            Политика конфиденциальности
          </a>
          <a href={CABINET_COPY.externalLinks.consent} target="_blank" rel="noopener noreferrer">
            Обработка персональных данных
          </a>
        </div>
      </footer>
    </PageShell>
  );
}
