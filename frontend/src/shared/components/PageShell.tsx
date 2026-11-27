import type { ReactNode } from 'react';

/**
 * Семейства страниц: «сцена» (лендинг и кабинет) и «паспорт» (вход).
 * Задают активный файл стилей: stage-original.css или auth-original.css.
 */
export type PageKind = 'registration' | 'cabinet' | 'auth';

type PageShellProps = {
  /** Семейство страницы: выбирает файл стилей. */
  page: PageKind;
  children: ReactNode;
};

/**
 * Корневая обёртка страницы: вешает data-page,
 * по которому скин применяет стили своего семейства.
 * Дизайн сайта — единственная, оригинальная концепция дизайнера
 * (зелёная), выбранная командой; переключение тем удалено.
 */
export function PageShell({ page, children }: PageShellProps) {
  return (
    <div data-page={page} className="page-shell">
      {children}
    </div>
  );
}
