import { Link } from 'react-router-dom';

/**
 * Страница «не найдено»: путь вне маршрутов приложения.
 * Оформление минимальное и автономное — не зависит от скинов страниц.
 */
export function NotFoundPage() {
  return (
    <div className="nf-shell">
      <main>
        <h1>Такой тропы нет</h1>
        <p>Страница не найдена — вернитесь на размеченные маршруты.</p>
        <nav>
          <Link to="/">На главную</Link>
          <Link to="/login">К входу</Link>
        </nav>
      </main>
    </div>
  );
}
