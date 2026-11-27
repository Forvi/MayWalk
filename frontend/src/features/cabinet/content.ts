/**
 * Контент личного кабинета — перенос src/v4/content.js из дизайн-макета.
 * Маршруты, подписи действий и служебные ссылки сохранены дословно.
 */

export type RouteOption = {
  id: string;
  distance: string;
  title: string;
  description: string;
  kind: string;
  origin: string;
};

export const SEASON = 2026;

export const ROUTES: RouteOption[] = [
  {
    id: 'r15',
    distance: '15 км',
    title: 'Городской маршрут',
    description: 'Парки, набережная Исети и центр Екатеринбурга.',
    kind: 'Город',
    origin: 'https://mayprogulka.ru/about/distance/',
  },
  {
    id: 'r18',
    distance: '18 км',
    title: 'Шарташская кругосветка',
    description: 'Лесные тропы вокруг озера Шарташ.',
    kind: 'Лес и озеро',
    origin: 'https://mayprogulka.ru/about/distance/',
  },
  {
    id: 'r30plus3',
    distance: '30+3 км',
    title: 'Через Лиственную',
    description: 'Шарташ, Калиновский лесопарк и гора Лиственная.',
    kind: 'Лес',
    origin: 'https://mayprogulka.ru/about/distance/',
  },
  {
    id: 'r50',
    distance: '50 км',
    title: 'Классика',
    description: 'Длинная прогулка до Берёзовских песков.',
    kind: 'Большая дистанция',
    origin: 'https://mayprogulka.ru/about/distance/',
  },
];

export const CABINET_COPY = {
  page: {
    journeyTitle: 'Моя прогулка',
    profileTitle: 'Профиль',
    context: `Сезон ${SEASON}`,
    applicationStatus: 'Участие не оформлено',
    emptyTitle: 'Выберите маршрут',
    emptyBody: 'Подберите подходящую дистанцию.',
    participantLabel: 'Участник',
  },
  actions: {
    chooseRoute: 'Выбрать маршрут',
    changeRoute: 'Изменить маршрут',
    saveRoute: 'Сохранить выбор',
    removeRoute: 'Убрать выбор',
    cancel: 'Отменить',
    removeConfirm: 'Убрать маршрут',
    editProfile: 'Изменить данные',
    saveProfile: 'Сохранить',
    sample: 'Черновик анкеты',
    printSample: 'Распечатать черновик',
    close: 'Закрыть окно',
    organizerCabinet: 'Кабинет организатора',
    exit: 'Выйти',
  },
  routeDialog: {
    title: 'Выберите свой маршрут',
    description: `Маршруты сезона ${SEASON}. Выберите дистанцию по силам.`,
    requiredError: 'Выберите маршрут.',
    saved: 'Маршрут выбран.',
    removeTitle: 'Убрать выбранный маршрут?',
    removeBody: 'Вы сможете выбрать его снова. Данные профиля сохранятся.',
    removed: 'Выбор маршрута убран.',
  },
  preparation: {
    title: 'Подготовка к прогулке',
    body: 'После оформления участия распечатайте анкету и подпишите её перед стартом.',
    emptyNote: 'Сначала выберите маршрут.',
  },
  sample: {
    title: 'Черновик анкеты',
    watermark: 'ЧЕРНОВИК',
    note: 'Проверьте данные перед оформлением участия.',
    routeLabel: `Маршрут · Сезон ${SEASON}`,
  },
  profile: {
    title: 'Данные профиля',
    description: 'Укажите логин и электронную почту для связи.',
    loginLabel: 'Логин',
    emailLabel: 'Электронная почта',
    loginPlaceholder: 'Придумайте логин',
    emailPlaceholder: 'you@example.ru',
    invalidLogin: 'Введите логин.',
    invalidEmail: 'Введите корректный адрес электронной почты.',
    saved: 'Изменения применены.',
    dirty: 'Есть несохранённые изменения.',
  },
  externalLinks: {
    faq: 'https://mayprogulka.ru/faq',
    help: 'Помощь',
    organizerCabinet: 'https://mayprogulka.ru/form/',
    participants: 'https://mayprogulka.ru/about/uchast/',
    routes: 'https://mayprogulka.ru/about/distance/',
    privacy: 'https://mayprogulka.ru/about/pdpol/',
    consent: 'https://mayprogulka.ru/about/pdsogl/',
  },
} as const;

export function routeById(id: string | null): RouteOption | null {
  if (!id) {
    return null;
  }
  return ROUTES.find((route) => route.id === id) ?? null;
}
