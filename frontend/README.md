# MayWalk · frontend

Фронтенд сайта «Майская прогулка»: React 19 + TypeScript (strict) + Vite + React Router 7. Оформление — оригинальная (зелёная) концепция дизайнера; разметка и тексты сохранены дословно.

Живёт в монорепозитории MayWalk как директория `frontend/`; API-контракты общие с Go-бэкендом — типы генерируются из `api/openapi.yml`.

## Разработка (из корня репозитория)

```bash
make front-dev        # dev-сервер Vite на :5173
make front-lint       # ESLint
make front-test       # юнит-тесты (vitest)
make front-generate   # TS-типы из ../api/openapi.yml → src/api/schema.d.ts
make front-build      # docker-образ maywalk-frontend (многоэтапная сборка)
make front-run        # сборка + запуск контейнера на :8082
```

Внутри `frontend/` доступны прямые npm-скрипты: `dev`, `build`, `preview`, `lint`, `typecheck`, `test`, `generate:api`.

## OpenAPI-типы

`make front-generate` (или `npm run generate:api`) превращает `api/openapi.yml` в `src/api/schema.d.ts` — те же контракты, что и у бэкенда (oapi-codegen). Типизированный клиент — `src/api/client.ts` (openapi-fetch, cookie-сессии: `credentials: include`). Базовый URL бэкенда — через `VITE_API_URL` (по умолчанию `http://localhost:8080`, в docker-сети `http://backend:8080`).

## Docker

`frontend/Dockerfile` — многоэтапная сборка: этап 1 (node:22-alpine) делает `npm ci` и `npm run build` (строгий tsc + Vite → `dist/`); этап 2 (nginx:1.27-alpine) раздаёт готовую статику c SPA-роутингом (`nginx.conf`). Контекст сборки — корень репозитория, как у `backend/Dockerfile`.

## Структура (кратко)
```

## Маршруты

| Путь                | Страница                              | Доступ                        |
| ------------------- | ------------------------------------- | ----------------------------- |
| `/`                 | Лендинг: витрина + форма регистрации  | публичный                     |
| `/login`            | Вход («паспорт»)                      | только гость (`GuestRoute`)   |
| `/register`         | Регистрация («паспорт»)               | только гость (`GuestRoute`)   |
| `/recovery`         | Восстановление пароля                 | только гость (`GuestRoute`)   |
| `/cabinet/journey`  | ЛК · Моя прогулка (маршрут, анкета)   | **участник** (`ProtectedRoute`) |
| `/cabinet/profile` | ЛК · Профиль (логин, почта)           | **участник** (`ProtectedRoute`) |

Проверка авторизации: `src/app/ProtectedRoute.tsx` перекидывает гостя на `/login`, запоминая адрес возврата в `state.from`; `GuestRoute` входящего сразу ведёт в кабинет.

Авторизация — демо на `localStorage` (ключи `mp:users`, `mp:session`), без бэкенда: `src/features/auth/AuthContext.tsx`. Валидация полей и формулировки ошибок перенесены из `auth-state.js` дизайнера: `src/features/auth/validation.ts`.

## Структура

```
src/
├── main.tsx                — входная точка, подключение стилей
├── App.tsx                 — провайдер авторизации и карта маршрутов
├── app/
│   ├── GuestRoute.tsx      — гостевой гвард
│   ├── ProtectedRoute.tsx  — гвард авторизации (state.from)
│   └── NotFoundPage.tsx    — 404
├── features/
│   ├── auth/
│   │   ├── AuthContext.tsx — демо-авторизация (localStorage)
│   │   ├── validation.ts   — валидация из auth-state.js дизайнера
│   │   ├── PassportLayout.tsx — «паспортная» рамка: шапка+постер+бланк
│   │   └── pages/          — Login, Register, Recovery
│   ├── cabinet/
│   │   ├── content.ts      — маршруты/тексты/ссылки из v4/content.js
│   │   ├── CabinetLayout.tsx — шапка, табы, выход
│   │   ├── components/CabinetDialogs.tsx — диалоги: маршрут / убрать / черновик
│   │   └── pages/          — Journey (маршрут), Profile (данные)
│   └── landing/LandingPage.tsx — лендинг-регистрация
├── shared/components/       — PageShell (скоупы), BrandMark/FlowerMark, ModalDialog,
│                             PasswordField (переключатель видимости пароля)
├── styles/
│   ├── fonts.css           — @font-face (Onest, UnboundedOriginal, UnboundedLegacy)
│   ├── global.css          — каркас документа (.page-shell)
│   ├── controls.css        — exit-link, автономная страница 404
│   └── pages/              — скины страниц (см. ниже)
└── assets/
    ├── fonts/              — woff2/ttf + лицензии OFL
    └── images/             — фирменные PNG из CSS дизайнера
```

### Скины страниц (`src/styles/pages/`)

| Файл               | Скоуп CSS-правил                                          | Что оформляет                       |
| ------------------ | --------------------------------------------------------- | ------------------------------------ |
| `stage-original.css` | `:is([data-page="registration"], [data-page="cabinet"])` | лендинг и кабинет (UnboundedLegacy)  |
| `auth-original.css`  | `[data-page="auth"]`                                     | «паспорт» входа (Onest+UnboundedOriginal) |

Каждая страница рендерится внутри `<PageShell page=…>` (`shared/components/PageShell.tsx`), который вешает на корневой `div` атрибут `data-page` (`registration` / `cabinet` / `auth`). Скины загружаются вместе, но не конфликтуют: у дизайнера классы повторяются между семействами страниц, изоляцию обеспечивает скоуп.

## Что изменилось относительно исходника

Исходник: `~/Downloads/mayprogulka-all-in-one.html` и `~/Downloads/mayprogulka-pages-r4/` (демо-страницы, ~5.7 тыс. строк HTML).

- **Выбрана одна концепция**: фестивальный и авторский скины, переключатель тем, каталог демо-страниц и их шрифты (Manrope, Dela Gothic One, Oswald) удалены; «авторская» SVG-карта постера вырезана из разметки (в оригинале постер — PNG).
- **Шрифты** вынуты из base64-вставок (по ~250 КБ на страницу) в `src/assets/fonts/`; лицензии OFL сохранены.
- **Стили**: инлайн-`<style>` каждой страницы → отдельные скоупленные файлы; `@font-face`-дубликаты и мёртвые правила удалены; демо-хак `.passport form {display:none}` и глобальные утечки (`:root`, `body`, `html`, голые селекторы в `@media`) убраны.
- **Разметка и логика**: hash-роутинг (`#login`, `#journey`) → маршруты React Router; подмена `innerHTML` из `templates.js` → компоненты; стейт-машины `auth-state.js`/`cabinet-state.js` → хуки с сохранением сообщений и переходов; нативные `<dialog>` сохранены, обёртка возвращает фокус на открывшую кнопку.
- **Токены страниц**: сидят на корневом `div` (бывший `body`/`:root` дизайнера).
