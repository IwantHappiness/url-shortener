# URL Shortener Backend

## API overview

В текущей реализации у сервиса есть два типа HTTP-маршрутов:

- версионированный JSON API: `/api/v1/...`
- публичный редирект по короткому коду: `/{shortURL}`

### Route map

| Method   | Path                 | Description                     |
| -------- | -------------------- | ------------------------------- |
| `POST`   | `/api/v1/users`      | Create user                     |
| `GET`    | `/api/v1/users`      | Get users list                  |
| `GET`    | `/api/v1/users/{id}` | Get user by id                  |
| `PATCH`  | `/api/v1/users/{id}` | Partially update user           |
| `DELETE` | `/api/v1/users/{id}` | Delete user                     |
| `POST`   | `/api/v1/urls`       | Create short URL                |
| `GET`    | `/api/v1/urls`       | Get URLs list                   |
| `GET`    | `/api/v1/urls/{id}`  | Get URL by id                   |
| `PATCH`  | `/api/v1/urls/{id}`  | Regenerate short code           |
| `DELETE` | `/api/v1/urls/{id}`  | Delete URL                      |
| `GET`    | `/api/v1/stats`      | Get statistics for a short code |
| `GET`    | `/{shortURL}`        | Redirect to original URL        |

## Common conventions

- Все `POST`/`PATCH`, кроме `PATCH /api/v1/urls/{id}`, работают с JSON body.
- Временные поля сериализуются в формате `ГГГГ-ММ-ДД` (например `"2024-06-28"`).
  - Параметры `from`/`to` в `/api/v1/stats` принимаются в том же формате.
- Все успешные ответы — JSON, кроме:
  - `GET /{shortURL}` → `307 Temporary Redirect`
  - `DELETE ...` → `204 No Content`
- `short_url` в JSON — это именно короткий код, а не полный абсолютный URL.
  - Пример: `"short_url": "V1StGXR8_Z5j"`
  - Публичная короткая ссылка строится как `<host>/<short_url>`

## Common error format

При ошибке API возвращает JSON вида:

```json
{
  "message": "failed to get user",
  "error": "user with id='42': not found"
}
```

### Status codes

| Status                      | Meaning                                                           |
| --------------------------- | ----------------------------------------------------------------- |
| `400 Bad Request`           | Некорректный JSON, невалидные path/query params, ошибки валидации |
| `404 Not Found`             | Сущность не найдена                                               |
| `409 Conflict`              | Конфликт при конкурентном обновлении                              |
| `500 Internal Server Error` | Любая иная серверная/инфраструктурная ошибка                      |

> Важно: некоторые DB-constraint ошибки в текущей реализации не мапятся в `409`, поэтому могут возвращаться как `500`.

## Response object schemas

### User object

```json
{
  "id": 1,
  "version": 1,
  "nickname": "alice",
  "email": "alice@example.com",
  "created_at": "2026-06-28T12:00:00Z"
}
```

| Field        | Type               | Description      |
| ------------ | ------------------ | ---------------- |
| `id`         | `integer`          | User ID          |
| `version`    | `integer`          | Версия записи    |
| `nickname`   | `string`           | Ник пользователя |
| `email`      | `string`           | Email            |
| `created_at` | `string (RFC3339)` | Время создания   |

### URL object

```json
{
  "id": 1,
  "version": 1,
  "user_id": 1,
  "original_url": "https://example.com",
  "short_url": "V1StGXR8_Z5j",
  "created_at": "2026-06-28T12:00:00Z"
}
```

| Field          | Type               | Description    |
| -------------- | ------------------ | -------------- |
| `id`           | `integer`          | URL ID         |
| `version`      | `integer`          | Версия записи  |
| `user_id`      | `integer`          | ID владельца   |
| `original_url` | `string`           | Исходный URL   |
| `short_url`    | `string`           | Короткий код   |
| `created_at`   | `string (RFC3339)` | Время создания |

### Stats object

```json
{
  "short_url": "V1StGXR8_Z5j",
  "original_url": "https://example.com",
  "created_at": "2026-06-28T12:00:00Z",
  "total_clicks": 3,
  "unique_ips": 2,
  "last_clicked_at": "2026-06-28T13:30:00Z"
}
```

| Field             | Type               | Description                                                       |
| ----------------- | ------------------ | ----------------------------------------------------------------- |
| `short_url`       | `string`           | Короткий код                                                      |
| `original_url`    | `string`           | Исходный URL                                                      |
| `created_at`      | `string (RFC3339)` | Время создания короткой ссылки                                    |
| `total_clicks`    | `integer`          | Число кликов в выбранном диапазоне                                |
| `unique_ips`      | `integer`          | Количество уникальных IP в выбранном диапазоне                    |
| `last_clicked_at` | `string (RFC3339)` | Последний клик в диапазоне; поле отсутствует, если кликов не было |

## Users API

### `POST /api/v1/users`

Создать пользователя.

#### Request body

```json
{
  "nickname": "alice",
  "email": "alice@example.com"
}
```

| Field      | Type     | Required | Validation                                              |
| ---------- | -------- | -------- | ------------------------------------------------------- |
| `nickname` | `string` | yes      | от 3 до 20 символов                                     |
| `email`    | `string` | yes      | валидный email; практический лимит БД — до 254 символов |

#### Success response

- `201 Created`
- body: `User object`

#### Typical errors

- `400` — невалидный JSON / `nickname` / `email`
- `500` — прочие ошибки БД/сервера

> В текущей реализации конфликт уникальности по `nickname`/`email` не документирован как отдельный `409` и может прийти как `500`.

---

### `GET /api/v1/users`

Получить список пользователей.

#### Query params

| Param    | Type      | Required | Description        |
| -------- | --------- | -------- | ------------------ |
| `limit`  | `integer` | no       | должен быть `>= 0` |
| `offset` | `integer` | no       | должен быть `>= 0` |

#### Success response

- `200 OK`
- body: массив `User object`

```json
[
  {
    "id": 1,
    "version": 1,
    "nickname": "alice",
    "email": "alice@example.com",
    "created_at": "2026-06-28T12:00:00Z"
  }
]
```

#### Notes

- Сортировка: `id ASC`
- Если данных нет, возвращается пустой массив `[]`

#### Typical errors

- `400` — `limit`/`offset` не integer или `< 0`
- `500` — серверная ошибка

---

### `GET /api/v1/users/{id}`

Получить пользователя по ID.

#### Path params

| Param | Type      | Required |
| ----- | --------- | -------- |
| `id`  | `integer` | yes      |

#### Success response

- `200 OK`
- body: `User object`

#### Typical errors

- `400` — `id` не integer
- `404` — пользователь не найден
- `500` — серверная ошибка

---

### `PATCH /api/v1/users/{id}`

Частично обновить пользователя.

#### Path params

| Param | Type      | Required |
| ----- | --------- | -------- |
| `id`  | `integer` | yes      |

#### Request body

Все поля опциональны.

```json
{
  "nickname": "alice_new",
  "email": "alice_new@example.com"
}
```

| Field      | Type     | Required | Notes                                                 |
| ---------- | -------- | -------- | ----------------------------------------------------- |
| `nickname` | `string` | no       | если передан, не может быть `null`; длина 3..20       |
| `email`    | `string` | no       | если передан, не может быть `null`; должен быть email |

#### Notes

- Отсутствующее поле означает `не менять`.
- Явное `null` запрещено и даёт `400`.
- Пустой body приведёт к ошибке декодирования JSON; для no-op patch используйте `{}`.
- `PATCH {}` допустим и приведёт к обновлению версии записи (`version`) без изменения данных.
- В текущей transport-валидации для `email` дополнительно применяется ограничение длины `3..20` символов. README фиксирует текущее поведение кода, а не желаемое поведение.

#### Success response

- `200 OK`
- body: `User object`

#### Typical errors

- `400` — невалидный `id`, JSON, `null`, невалидные поля
- `404` — пользователь не найден
- `409` — конкурентный конфликт при обновлении
- `500` — прочие ошибки БД/сервера

---

### `DELETE /api/v1/users/{id}`

Удалить пользователя.

#### Path params

| Param | Type      | Required |
| ----- | --------- | -------- |
| `id`  | `integer` | yes      |

#### Success response

- `204 No Content`
- body отсутствует

#### Notes

- Удаление пользователя каскадно удаляет его URL-ы.
- Клики по этим URL-ам также удаляются каскадно через связь с таблицей `urls`.

#### Typical errors

- `400` — `id` не integer
- `404` — пользователь не найден
- `500` — серверная ошибка

## URLs API

### `POST /api/v1/urls`

Создать короткую ссылку.

#### Request body

```json
{
  "url": "https://example.com",
  "user_id": 1
}
```

| Field     | Type      | Required | Validation                                                               |
| --------- | --------- | -------- | ------------------------------------------------------------------------ |
| `url`     | `string`  | yes      | валидный URL, схема только `http`/`https`, непустой host, длина `< 2048` |
| `user_id` | `integer` | yes      | обязательный ID владельца                                                |

#### Success response

- `201 Created`
- body: `URL object`

#### Notes

- `short_url` генерируется автоматически.
- Текущая длина генерируемого короткого кода — `12` символов.

#### Typical errors

- `400` — невалидный JSON / `url` / `user_id`
- `404` — пользователь с `user_id` не найден
- `409` — конфликт короткого кода
- `500` — прочие ошибки БД/сервера

---

### `GET /api/v1/urls`

Получить список коротких ссылок.

#### Query params

| Param     | Type      | Required | Description                             |
| --------- | --------- | -------- | --------------------------------------- |
| `user_id` | `integer` | no       | фильтр по владельцу, должен быть `>= 0` |
| `limit`   | `integer` | no       | должен быть `>= 0`                      |
| `offset`  | `integer` | no       | должен быть `>= 0`                      |

#### Success response

- `200 OK`
- body: массив `URL object`

```json
[
  {
    "id": 1,
    "version": 1,
    "user_id": 1,
    "original_url": "https://example.com",
    "short_url": "V1StGXR8_Z5j",
    "created_at": "2026-06-28T12:00:00Z"
  }
]
```

#### Notes

- Сортировка: `id ASC`
- Если данных нет, возвращается пустой массив `[]`

#### Typical errors

- `400` — `user_id`/`limit`/`offset` не integer или `< 0`
- `500` — серверная ошибка

---

### `GET /api/v1/urls/{id}`

Получить короткую ссылку по ID.

#### Path params

| Param | Type      | Required |
| ----- | --------- | -------- |
| `id`  | `integer` | yes      |

#### Success response

- `200 OK`
- body: `URL object`

#### Typical errors

- `400` — `id` не integer
- `404` — ссылка не найдена
- `500` — серверная ошибка

---

### `PATCH /api/v1/urls/{id}`

Перегенерировать короткий код для существующей ссылки.

#### Path params

| Param | Type      | Required |
| ----- | --------- | -------- |
| `id`  | `integer` | yes      |

#### Request body

Body не используется. Можно отправлять запрос без body.

#### Success response

- `200 OK`
- body: `URL object`

#### Notes

- Endpoint генерирует новый случайный `short_url` длиной `12` символов.
- При успешном обновлении старый короткий код перестаёт резолвиться.
- `version` увеличивается.

#### Typical errors

- `400` — `id` не integer
- `404` — ссылка не найдена
- `409` — конкурентный конфликт при обновлении
- `500` — прочие ошибки БД/сервера

---

### `DELETE /api/v1/urls/{id}`

Удалить короткую ссылку.

#### Path params

| Param | Type      | Required |
| ----- | --------- | -------- |
| `id`  | `integer` | yes      |

#### Success response

- `204 No Content`
- body отсутствует

#### Notes

- Клики по этой ссылке удаляются каскадно.

#### Typical errors

- `400` — `id` не integer
- `404` — ссылка не найдена
- `500` — серверная ошибка

## Redirect API

### `GET /{shortURL}`

Публичный переход по короткому коду.

#### Path params

| Param      | Type     | Required |
| ---------- | -------- | -------- |
| `shortURL` | `string` | yes      |

#### Success response

- `307 Temporary Redirect`
- `Location` указывает на `original_url`

#### Notes

- При успешном редиректе сервис пытается записать клик в таблицу `clicks`.
- IP для статистики берётся в следующем порядке:
  1. `X-Forwarded-For`
  2. `X-Real-IP`
  3. `RemoteAddr`
- Если запись клика не удалась, редирект всё равно выполняется.

#### Typical errors

- `400` — пустой `shortURL`
- `404` — короткий код не найден
- `500` — ошибка резолва короткого кода

## Statistics API

### `GET /api/v1/stats`

Получить статистику по короткому коду.

#### Query params

| Param       | Type               | Required | Description                                    |
| ----------- | ------------------ | -------- | ---------------------------------------------- |
| `short_url` | `string`           | yes      | короткий код, например `V1StGXR8_Z5j`          |
| `from`      | `string (RFC3339)` | no       | нижняя граница диапазона кликов, включительно  |
| `to`        | `string (RFC3339)` | no       | верхняя граница диапазона кликов, включительно |

#### Example request

```text
GET /api/v1/stats?short_url=V1StGXR8_Z5j&from=2026-06-28T00:00:00Z&to=2026-06-28T23:59:59Z
```

#### Success response

- `200 OK`
- body: `Stats object`

#### Notes

- `short_url` должен быть именно кодом, а не полным URL.
- Фильтр `from`/`to` применяется только к кликам, а не к самой ссылке.
- Если ссылка существует, но в выбранном диапазоне нет кликов:
  - `total_clicks = 0`
  - `unique_ips = 0`
  - `last_clicked_at` отсутствует
- Для работы статистики и записи кликов должна быть применена миграция `backend/migrations/000002_add_history.up.sql`.

#### Typical errors

- `400` — отсутствует `short_url`, неверный формат `from`/`to`, `from > to`
- `404` — короткий код не найден
- `500` — серверная ошибка
