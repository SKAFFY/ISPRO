# ISPRO-app

Учебный проект по курсу "Инструменты современной разработки программного обеспечения" (ИТМО, 2025).

## Лабораторные работы

- **Lab1**: Регулярные выражения
- **Lab2**: OpenAPI сервис
- **Lab3**: Метрики
- **Lab4**: Журналирование
- **Lab5**: Распределенная трассировка
- **Lab6**: CI/CD

---

## Общие команды

| Команда | Описание |
|---------|----------|
| `make docker-up` | Запуск PostgreSQL, VictoriaMetrics, Grafana |
| `make docker-down` | Остановка всех Docker сервисов |
| `make migrate-up` | Применение миграций |
| `make migrate-down` | Откат миграций |
| `make start` | Сборка и запуск сервера |
| `make generate-api` | Генерация REST API из OpenAPI |
| `make di-generate` | Генерация DI контейнера |
| `make swagger` | Swagger UI на http://localhost:8081 |

---

## Lab1: Регулярные выражения

### Описание

Демонстрация простых и сложных регулярных выражений.

### Простое выражение

Валидация формата ссылки mindmap: `source:target:description`

Паттерн: `^[^:]+:[^:]+:[^:]+$`

- Три поля, разделённых двоеточием
- Каждое поле содержит хотя бы один символ

### Сложное выражение

Валидация пароля с использованием lookahead:

Паттерн: `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[a-zA-Z\d!@#$%^&*]{8,}$`

Требования:
- Минимум 8 символов
- Минимум 1 строчная буква (`(?=.*[a-z])`)
- Минимум 1 заглавная буква (`(?=.*[A-Z])`)
- Минимум 1 цифра (`(?=.*\d])`)
- Минимум 1 специальный символ (`(?=.*[!@#$%^&*])`)

### Запуск

```bash
make run-regexp
```

### Тесты

```bash
go test ./...
```

### Бенчмарки

```bash
go test -bench=. ./...
```

---

## Lab2: Сервис на OpenAPI

### Описание

REST API сервис для mindmap-приложения с использованием OpenAPI Specification и API-first подхода.

### Технологии

- Go
- OpenAPI 2.0 (Swagger)
- go-openapi (генерация кода)
- PostgreSQL

### API Endpoints

#### Entries

| Метод | Путь | Описание |
|-------|------|----------|
| GET | /entries | Список всех записей |
| POST | /entries | Создать запись |
| GET | /entries/{id} | Получить запись по ID |
| PUT | /entries/{id} | Обновить запись |
| DELETE | /entries/{id} | Удалить запись |

#### Links

| Метод | Путь | Описание |
|-------|------|----------|
| GET | /links | Список всех связей |
| POST | /links | Создать связь |
| GET | /links/{id} | Получить связь по ID |
| PUT | /links/{id} | Обновить связь |
| DELETE | /links/{id} | Удалить связь |

### Примеры запросов

```bash
# Создать запись
curl -X POST http://localhost:8080/entries \
  -H "Content-Type: application/json" \
  -d '{"title": "First Entry", "content": "Content here"}'

# Ответ:
# {"content":"Content here","created_at":"2026-05-25T21:57:26.130Z","id":1,"title":"First Entry","updated_at":"2026-05-25T21:57:26.130Z"}

# Получить все записи
curl http://localhost:8080/entries

# Ответ:
# [{"content":"Content here","created_at":"2026-05-25T21:57:26.130Z","id":1,"title":"First Entry","updated_at":"2026-05-25T21:57:26.130Z"}]

# Создать связь между записями
curl -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"source_id": 1, "target_id": 2}'

# Ответ:
# {"id":1,"source_id":1,"target_id":2}
```

### Структура проекта

```
cmd/server/main.go     # Точка входа
api/openapi.yaml       # OpenAPI спецификация
internal/
  ├── models/          # Модели данных (сгенерированы)
  ├── restapi/         # REST API обработчики (сгенерированы)
  ├── entities/        # Бизнес-логика
  └── di/              # Dependency Injection
```

---

## Lab3: Метрики

### Описание

Добавление метрик в сервис с использованием Prometheus client и визуализация в Grafana.

### Технологии

- Go + prometheus/client_golang
- VictoriaMetrics (Time Series DB)
- Grafana (визуализация)

### Метрики

#### Стандартные метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `http_requests_total` | Counter | Всего HTTP запросов |
| `http_request_duration_seconds` | Histogram | Время выполнения запроса |

#### Продуктовые метрики

| Метрика | Тип | Описание |
|---------|-----|----------|
| `entries_created_total` | Counter | Всего создано записей |
| `entries_deleted_total` | Counter | Всего удалено записей |
| `links_created_total` | Counter | Всего создано связей |
| `links_deleted_total` | Counter | Всего удалено связей |

### Доступ к сервисам

- **Приложение**: http://localhost:8080
- **Метрики**: http://localhost:8080/metrics
- **Grafana**: http://localhost:3000 (admin/admin)
- **VictoriaMetrics**: http://localhost:8428

### Примеры PromQL запросов

```promql
# RPS (запросов в секунду)
sum(rate(http_requests_total[5m])) by (method, path)

# p95 latency
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (method, path, le))

# p99 latency
histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (method, path, le))

# Всего создано записей
entries_created_total

# Скорость создания записей (в секунду)
rate(entries_created_total[5m])

# Сколько сейчас активных горутин
go_goroutines
```

### Дашборд

В Grafana автоматически подключается дашборд `ISPRO App Metrics` с панелями:
- HTTP Requests Rate (RPS)
- HTTP Latency (p95, p99)
- Product Metrics (stat панели)
- Product Metrics (time series)

---

## Lab4: Журналирование

### Описание

Добавление структурированного журналирования (логов) в сервис с использованием встроенного пакета `log/slog` (Go 1.21+) и сбор логов в Loki с визуализацией в Grafana.

### Технологии

- Go + `log/slog` (структурированное логирование в JSON)
- Loki (хранение логов)
- Grafana (визуализация)

### Доступ к сервисам

- **Loki**: http://localhost:3100
- **Grafana Logs Dashboard**: http://localhost:3000 (admin/admin)

### Примеры LogQL запросов

```logql
# Все логи приложения
{service="ispro-app"}

# Только ошибки
{service="ispro-app"} |= `"level":"error"`

# Только предупреждения
{service="ispro-app"} |= `"level":"warn"`

# Поиск по конкретному ID записи
{service="ispro-app"} |= `"id":1`

# Поиск по title
{service="ispro-app"} |= "Title"

# Подсчёт логов по уровням
sum by (level) (count_over_time({service="ispro-app"}[5m]))
```

### Структура логов

Каждое сообщение — JSON строка с полями:

| Поле | Описание |
|------|----------|
| `time` | Время события (RFC3339Nano) |
| `level` | Уровень логирования (info, warn, error) |
| `service` | Имя сервиса (ispro-app) |
| `msg` | Текстовое сообщение |
| `id` | ID сущности (если применимо) |
| `title` | Заголовок записи (если применимо) |
| `error` | Детали ошибки (если есть) |

### Дашборд

В Grafana автоматически подключается дашборд `ISPRO App Logs` с панелью:
- Application Logs — все логи приложения с возможностью фильтрации по уровню, поиску по тексту и времени

---

## Lab5: Распределенная трассировка

### Описание

Добавление распределённой трассировки в сервис с использованием OpenTelemetry и визуализация в Jaeger. Каждый HTTP запрос отслеживается как trace с вложенными span'ами.

### Технологии

- Go + OpenTelemetry SDK
- Jaeger (сбор, хранение и визуализация трассировок)
- Grafana (просмотр через Jaeger datasource)

### Доступ к сервисам

- **Jaeger UI**: http://localhost:16686
- **Grafana**: http://localhost:3000 (admin/admin)

### Использование

1. Отправить HTTP запрос к API:
   ```bash
   curl http://localhost:8080/entries
   ```
2. Открыть Jaeger UI: http://localhost:16686
3. В поиске выбрать сервис `ispro-app`
4. Нажать "Find Traces"
5. Кликнуть на trace для просмотра span'ов

### Структура span'ов

| Span | Описание |
|------|----------|
| `server` | Корневой span — обработка входящего HTTP запроса |
| `HTTP {method} {path}` | Дочерний span — маршрутизация и бизнес-логика |

### Дашборд

В Jaeger UI доступны:
- Поиск трассировок по сервису, времени, тегам
- Детальный просмотр span'ов с длительностью
- Flame graph (диаграмма Ганта)
- Сравнение трассировок

В Grafana через datasource `Jaeger` можно просматривать трассировки из Explore.

---

## Lab6: CI/CD

### Описание

Настройка непрерывной интеграции (Continuous Integration) с помощью GitHub Actions. При пуше в ветку `lab6` автоматически выполняются: линтинг кода, тестирование, сборка бинарного файла и сборка Docker образа с публикацией в GitHub Container Registry.

### Технологии

- GitHub Actions
- golangci-lint (статический анализ кода)
- Go test (юнит-тесты с race detector и coverage)
- Docker + GitHub Container Registry (ghcr.io)

### Этапы пайплайна

| Этап | Описание |
|------|----------|
| `lint` | Статический анализ кода через golangci-lint |
| `test` | Запуск тестов с race detector и подсчётом coverage |
| `build` | Компиляция бинарного файла |
| `docker` | Сборка Docker образа и публикация в ghcr.io (зависит от lint, test, build) |

### Триггер

Workflow запускается автоматически при пуше в ветку `lab6`.

### Docker образ

Образ публикуется в GitHub Container Registry по адресу:
```
ghcr.io/skaffy/ispro-app:lab6
ghcr.io/skaffy/ispro-app:<commit-sha>
```

### Сервисы

На этапе тестирования поднимается PostgreSQL 16 в качестве сервиса для прогона миграций перед тестами.
