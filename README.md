# MayWalk

## Структура проекта

```
MayWalk/
├── api/                    # OpenAPI спецификация
├── backend/                # Go backend
├── docker-compose.yml      # Контейнеры
├── Makefile                # Команды
└── .env                    # Переменные окружения
```

## Быстрый старт

### 1. Настройка переменных окружения

```bash
cp .env.example .env     # Нужно изменить значения на валидные
```

### 2. Запуск через Docker Compose

**Полный стек (backend + swagger + инфраструктура):**
```bash
docker compose --profile dev up -d
```

**Только инфраструктура (postgres + redis + redisinsight):**
```bash
docker compose --profile infra up -d
```

**Только Swagger UI:**
```bash
docker compose --profile docs up -d
```

**Production (backend + инфраструктура, без swagger):**
```bash
docker compose --profile prod up -d
```

### 3. Запуск через Make (локально)

**Инфраструктура через Docker:**
```bash
docker compose --profile infra up -d
```

**Backend:**
```bash
make back-generate     # Генерация типов
make back-run          # Запуск без сборки
make back-build        # Сборка
```

**Swagger UI:**
```bash
make docs              # Запуск Swagger UI на :8081
```

**Другие команды:**
```bash
make back-test         # Тесты
make back-generate     # Генерация типов из OpenAPI
make back-migrate      # Применение миграций
make help              # Список всех команд
```

### 4. Прямой запуск (без Make)

```bash
# Инфраструктура
docker compose --profile infra up -d

# Backend
cd backend && go run ./cmd/server/main.go

# Swagger UI
cd backend && go run ./cmd/docs/main.go
```

## Порты

| Сервис | Порт |
|---|---|
| Backend API | :8080 |
| Swagger UI | :8081 |
| PostgreSQL | :5432 |
| Redis | :6379 |
| RedisInsight | :5540 |

## Полезные ссылки

- Swagger UI: http://localhost:8081/swagger/
- RedisInsight: http://localhost:5540
- Health check: http://localhost:8080/health
