# Описание проекта

В этом репозитории выполняются лабораторные работы https://cs.petrsu.ru/~vadim/isrpo2025/

Ветки: лабораторные работы стартуют друг от друга.

Результаты в README в формате: общая документация, ЛР1, ЛР2...

## Стек

Go, digen, muonsoft/validation, go-openapi, PostgreSQL.

## Структура проекта

```
internal/
├── di/                              # DI контейнер
│   ├── container.go                 # СГЕНЕРИРОВАН
│   ├── internal/
│   │   ├── container.go             # СГЕНЕРИРОВАН
│   │   ├── definitions/container.go# ОПРЕДЕЛЕНИЯ (руками)
│   │   └── factories/container.go   # ФАБРИКИ (руками)
│   └── lookup/container.go          # СГЕНЕРИРОВАН
├── entities/{entity}/
│   ├── domain/                      # Модели, интерфейсы репозиториев, валидация
│   │   ├── repository.go            # Интерфейс репозитория
│   │   ├── {entity}.go              # Доменная модель с валидацией
│   │   └── validation/              # Кастомные валидаторы
│   ├── repository/                  # PostgreSQL реализация репозитория
│   ├── handler.go                   # HTTP обработчик
│   └── {create,get,list,update,delete}_{entity}.go  # Use cases
├── models/                          # СГЕНЕРИРОВАН из openapi.yaml
├── restapi/                         # СГЕНЕРИРОВАН из openapi.yaml
│   ├── operations/                  # HTTP обработчики и схемы
│   └── server.go                    # Сервер
└── migrations/                      # goose миграции
```

## Генерация кода

Проект использует два инструмента генерации:

| Инструмент | Источник | Генерирует | Команда |
|------------|----------|------------|---------|
| swagger | `api/openapi.yaml` | `internal/models/`, `internal/restapi/` | `make generate-api` |
| digen | `internal/di/internal/definitions/container.go` | `internal/di/container.go`, `internal/di/internal/container.go`, `internal/di/lookup/container.go` | `make di-generate` |

**Важно**: сгенерированный код НЕЛЬЗЯ удалять или править вручную.

## Как добавить новую сущность

1. Добавить определение в `api/openapi.yaml` (definitions + paths)
2. Запустить `make generate-api`
3. Создать доменную модель: `internal/entities/{entity}/domain/{entity}.go`
4. Создать интерфейс репозитория: `internal/entities/{entity}/domain/repository.go`
5. Реализовать репозиторий: `internal/entities/{entity}/repository/postgres.go`
6. Создать use cases: `create_{entity}.go`, `get_{entity}.go`, `list_{entity}.go`, `update_{entity}.go`, `delete_{entity}.go`
7. Создать handler: `internal/entities/{entity}/handler.go`
8. Добавить определение в DI: `internal/di/internal/definitions/container.go`
9. Запустить `make di-generate`
10. Зарегистрировать handler в `internal/di/internal/factories/container.go`

## Команды

| Команда | Описание |
|---------|----------|
| `make docker-up` | Запустить PostgreSQL в Docker |
| `make docker-down` | Остановить PostgreSQL |
| `make migrate-up` | Применить миграции |
| `make migrate-down` | Откатить миграции |
| `make generate-api` | Сгенерировать REST API из openapi.yaml |
| `make di-generate` | Сгенерировать DI контейнер |
| `make build-server` | Собрать приложение |
| `make start` | Собрать и запустить сервер |
| `make swagger` | Swagger UI на http://localhost:8081 |

## Запуск

```bash
make docker-up
make migrate-up
make start
```

Сервер запускается на http://localhost:8080