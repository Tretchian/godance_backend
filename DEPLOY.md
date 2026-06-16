# Деплой на VPS

Бэкенд разворачивается через Docker Compose. Снаружи смотрит только **Caddy**,
который автоматически выпускает TLS-сертификат (Let's Encrypt) — это обязательно,
т.к. фронт на Vercel ходит по HTTPS (иначе браузер заблокирует запросы как
mixed-content).

```
[Vercel фронт] ──HTTPS──▶ [Caddy :443] ──▶ [app:8080 (Gin)] ──▶ [postgres]
                                       └──▶ [worker]        └──▶ [minio]
```

## 0. Предусловия

- VPS с Linux (Ubuntu/Debian) и root/sudo.
- Доменное имя и доступ к его DNS.
- Установленные Docker и Docker Compose plugin:
  ```bash
  curl -fsSL https://get.docker.com | sh
  ```

## 1. DNS

Создайте A-запись на IP вашего VPS:

| Запись | Тип | Значение |
| --- | --- | --- |
| `api.ваш-домен` | A | IP VPS |
| `storage.ваш-домен` *(опц., для видео)* | A | IP VPS |

Дождитесь распространения (`dig api.ваш-домен` должен вернуть IP VPS) — иначе
Caddy не сможет выпустить сертификат.

## 2. Файрвол

Наружу нужны только SSH и HTTP(S). Postgres/MinIO/app забиндены на `127.0.0.1`
и недоступны извне.
```bash
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

## 3. Код и конфиг

```bash
git clone <repo-url> godance && cd godance
git checkout claude/optimistic-goodall-0qufca   # или ветка/тег релиза
cp .env.example .env
```

Отредактируйте `.env` — обязательно поменяйте:

| Переменная | Значение |
| --- | --- |
| `DB_PASSWORD` | надёжный пароль |
| `JWT_SECRET` | длинная случайная строка (`openssl rand -hex 32`) |
| `MINIO_ACCESS_KEY` / `MINIO_SECRET_KEY` | свои ключи |
| `API_DOMAIN` | `api.ваш-домен` |
| `ACME_EMAIL` | ваш email |
| `CORS_ALLOWED_ORIGINS` | прод-домен фронта (Vercel), напр. `https://dancefeedbackplatform.vercel.app` |

### Видео (MinIO) — выберите одно

**A. Без видео (проще, для демо):** закомментируйте `MINIO_ENDPOINT` в `.env`.
Бэкенд возьмёт stub-хранилище — видео-эндпоинты вернут фиктивные URL, остальное
работает полностью.

**B. С видео (публичный MinIO):**
- задайте `STORAGE_DOMAIN=storage.ваш-домен` в `.env`;
- раскомментируйте блок `{$STORAGE_DOMAIN}` в `Caddyfile`;
- выставьте `MINIO_PUBLIC_ENDPOINT=https://storage.ваш-домен`
  (а `MINIO_ENDPOINT=http://minio:9000` оставьте — это внутренний адрес).

## 4. Запуск

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

Поднимутся: `caddy`, `app`, `worker`, `postgres`, `minio`. Схема БД создаётся
автоматически (AutoMigrate при старте `app`).

## 5. Проверка

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
curl https://api.ваш-домен/api/v1/competitions   # должен вернуть JSON-список
```

Логи:
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f app worker caddy
```

## 6. Подключение фронта

В Vercel (Project Settings → Environment Variables):
```
VITE_API_BASE_URL=https://api.ваш-домен/api/v1
VITE_USE_MOCKS=false
```
И убедитесь, что `.env.production` в репозитории фронта не форсит `VITE_USE_MOCKS=true`.

## 7. Обновление версии

```bash
git pull
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

## Шпаргалка по неполадкам

| Симптом | Причина / решение |
| --- | --- |
| Caddy не выпускает сертификат | DNS ещё не указывает на VPS, или закрыт порт 80/443 |
| `failed to init MinIO storage: ... required vars are empty` | задайте `MINIO_ACCESS_KEY/SECRET/BUCKET` или уберите `MINIO_ENDPOINT` (stub) |
| Фронт: CORS-ошибка | домен фронта не в `CORS_ALLOWED_ORIGINS` |
| Фронт: mixed content | `VITE_API_BASE_URL` должен быть `https://`, не `http://` |
| presigned-загрузка видео падает с SignatureDoesNotMatch | `MINIO_PUBLIC_ENDPOINT` должен совпадать с доменом, по которому ходит браузер |
| app не стартует, ждёт БД | норма при первом старте — есть ретраи; смотрите логи postgres |
