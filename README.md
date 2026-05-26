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
make run
```

### Тесты

```bash
make test
```

### Бенчмарки

```bash
make bench
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

### Запуск

```bash
# Запуск БД
make docker-up

# Применение миграций
make migrate-up

# Запуск сервера
make run
```

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

### Генерация кода из OpenAPI

```bash
make generate
```

Код генерируется с помощью `swagger generate server` по спецификации `api/openapi.yaml`.

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