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