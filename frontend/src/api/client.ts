import createClient from 'openapi-fetch';

import type { components, paths } from './schema';

/**
 * Типизированный клиент MayWalk API.
 *
 * Типы `paths` и `components` генерируются из api/openapi.yml командой
 * `npm run generate:api` (make front-generate) — это те же контракты,
 * что использует Go-бэкенд (oapi-codegen). После изменения спеки достаточно
 * перегенерировать файл, и TypeScript подсветит все несовпадения на фронте.
 *
 * Авторизация — сессии в cookie (session_id), поэтому credentials: include:
 * без него браузер не пришлёт cookie при кросс-доменных запросах к бэкенду.
 */
export const apiClient = createClient<paths>({
  baseUrl: import.meta.env.VITE_API_URL ?? 'http://localhost:8080',
  credentials: 'include',
});

export type UserResponse = components['schemas']['UserResponse'];
export type ErrorResponse = components['schemas']['ErrorResponse'];
export type RegisterUserRequest = components['schemas']['RegisterUserRequest'];
export type LoginUserRequest = components['schemas']['LoginUserRequest'];
