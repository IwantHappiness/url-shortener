# URL Shortener

Сервис для создания коротких ссылок с отслеживанием статистики переходов.

## Архитектура

```
┌──────────┐      ┌──────────┐      ┌───────────┐
│ Frontend │ ───→ │  Nginx   │ ───→ │  Backend  │ ───→ PostgreSQL
│  (React) │      │ (прокси) │      │   (Go)    │
└──────────┘      └──────────┘      └───────────┘
```

- **Backend** — Go 1.26, `net/http` (ServeMux), PostgreSQL, pgx
- **Frontend** — React 19, TypeScript, Vite, nginx
- **Infrastructure** — Docker Compose, migrate/migrate

## Быстрый старт

### Docker (полный стек)

```bash
# Из корня проекта
docker compose up -d backend frontend
```

После запуска:

- Фронтенд: http://127.0.0.1:3000
- Backend API: http://127.0.0.1:8080
- Swagger UI: http://127.0.0.1:8080/swagger/

### Локальная разработка

**1. Инфраструктура (PostgreSQL):**

```bash
make env-up
make migrate-up
```

**2. Backend:**

```bash
make run-app
```

**3. Frontend (отдельный терминал):**

```bash
cd frontend
npm install
npm run dev
```

Фронтенд будет на http://127.0.0.1:5173, Vite сам проксирует `/api` на бекенд.

> **Внимание:** В режиме Vite (`npm run dev`) редирект по короткой ссылке работает через кастомный middleware, который сначала пробует бекенд. Если бекенд недоступен — Vite отдаёт SPA.

## API

### Route map

| Method   | Path                 | Description                        |
| -------- | -------------------- | ---------------------------------- |
| `POST`   | `/api/v1/users`      | Создать пользователя               |
| `GET`    | `/api/v1/users`      | Список пользователей               |
| `GET`    | `/api/v1/users/{id}` | Пользователь по ID                 |
| `PATCH`  | `/api/v1/users/{id}` | Частичное обновление               |
| `DELETE` | `/api/v1/users/{id}` | Удалить пользователя               |
| `POST`   | `/api/v1/urls`       | Создать короткую ссылку            |
| `GET`    | `/api/v1/urls`       | Список ссылок                      |
| `GET`    | `/api/v1/urls/{id}`  | Ссылка по ID                       |
| `PATCH`  | `/api/v1/urls/{id}`  | Перегенерировать короткий код      |
| `DELETE` | `/api/v1/urls/{id}`  | Удалить ссылку                     |
| `GET`    | `/api/v1/stats`      | Статистика по короткому коду       |
| `GET`    | `/{shortURL}`        | Редирект (307) на оригинальный URL |

### Формат даты

Все даты в JSON ответах возвращаются в формате `ГГГГ-ММ-ДД`:

```json
{ "created_at": "2024-06-28" }
```

Параметры `from`/`to` в `/api/v1/stats` принимаются в том же формате.

### Формат ошибок

```json
{
  "message": "failed to get user",
  "error": "user with id='42': not found"
}
```

### Swagger

Полная OpenAPI-спецификация доступна:

- Swagger UI: `http://127.0.0.1:8080/swagger/`
- JSON: `http://127.0.0.1:8080/swagger/doc.json`

Генерация через `make swagger-gen`.

## Makefile

| Команда              | Описание                              |
| -------------------- | ------------------------------------- |
| `make env-up`        | Запустить PostgreSQL                  |
| `make env-down`      | Остановить PostgreSQL                 |
| `make migrate-up`    | Накатить миграции                     |
| `make migrate-down`  | Откатить миграции                     |
| `make run-app`       | Запустить бекенд локально             |
| `make swagger-gen`   | Сгенерировать OpenAPI-спецификацию    |
| `make up`            | Запустить backend + frontend в Docker |
| `make down`          | Остановить Docker-сервисы             |
| `make clean-up`      | Очистить volume-файлы (pgdata)        |
| `make clean-up-logs` | Очистить логовые файлы                |

## Структура проекта

```
url_shortener/
├── docker-compose.yml       # Полный стек (postgres, backend, frontend)
├── Makefile                 # Основные команды
├── backend/
│   ├── Dockerfile
│   ├── docker-compose.yml   # Инфраструктурные сервисы (legacy)
│   ├── Makefile
│   ├── cmd/api/main.go      # Точка входа
│   ├── internal/
│   │   ├── core/            # Общие компоненты (config, logger, errors, domain)
│   │   └── features/
│   │       ├── users/       # CRUD пользователей
│   │       ├── urls/        # CRUD коротких ссылок
│   │       ├── redirect/    # Редирект {shortURL} → оригинал
│   │       └── statistics/  # Статистика переходов
│   ├── migrations/
│   └── docs/                # Swagger-спецификация
└── frontend/
    ├── Dockerfile
    ├── nginx.conf
    ├── vite.config.ts
    └── src/
        ├── api/client.ts    # HTTP-клиент к бекенду
        ├── components/
        │   ├── UsersSection.tsx
        │   ├── UrlsSection.tsx
        │   └── StatsSection.tsx
        └── types/index.ts
```

## Детальная документация

- [Backend API контракты](./backend/README.md)
