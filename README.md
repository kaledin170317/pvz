# PVZ Management Backend

Сервис для управления пунктами выдачи заказов (ПВЗ), приёмкой товаров и учётом продукции.

---

## 📦 Основной функционал

- Регистрация и логин пользователей (JWT)
- Создание и фильтрация ПВЗ
- Приёмка товаров: создание, закрытие, добавление/удаление товаров
- gRPC метод получения всех ПВЗ без авторизации
- OpenAPI спецификация REST API
- CI/CD pipeline (GitHub Actions)

---

## 🧠 Условие задачи

Файл [`Backend-trainee-assignment-spring-2025.md`](./docs/Backend-trainee-assignment-spring-2025.md)

---

## 📁 Структура проекта

```
.
├── api/
│   ├── swagger.yaml
│   ├── oapi-codegen.yaml
│   └── pvz.proto
├── build/
│   ├── Dockerfile
│   └── docker-compose.yml
├── cmd/
│   └── pvzapp/
├── docs/
│   └── Backend-trainee-assignment-spring-2025.md
├── init/
│   ├── generate/
│   └── migrations/
├── internal/
│   ├── adapters/
│   │   ├── api/
│   │   │   ├── grpc/
│   │   │   └── rest/
│   │   └── db/
│   │       ├── mocks/
│   │       └── postgreSQL/
│   ├── app/
│   ├── config/
│   ├── domain/
│   │   ├── models/
│   │   ├── models_gen/
│   │   └── usecases/
│   └── usecases_impl/
├── tests/
│   ├── integration/
│   └── unit/
├── go.mod
└── README.md
```

---

## ⚙️ Переменные окружения по умолчанию

| Переменная      | Назначение      | Значение по умолчанию |
|-----------------|-----------------|-----------------------|
| `DB_HOST`       | Хост PostgreSQL | `localhost`           |
| `DB_PORT`       | Порт PostgreSQL | `5432`                |
| `DB_USER`       | Пользователь БД | `postgres`            |
| `DB_PASSWORD`   | Пароль БД       | `password`            |
| `DB_NAME`       | Название БД     | `pvz`                 |
| `APP_PORT`      | REST API порт   | `8080`                |
| `APP_GRPC_PORT` | gRPC API порт   | `3000`                |
| `JWT_SECRET`    | Секрет для JWT  | `jwt-secret`          |

---
## Юнит-тесты

```bash
go test ./tests/unit/... -v
```
---

## Интеграционные тесты


```bash
go test ./tests/integration/... -v -count=1
```

> Убедись, что PostgreSQL поднят на `localhost:5432` перед запуском!

---
## 🔧 Генерация моков

```bash
go generate ./init/generate/...
```

---

## ✨ Генерация DTO по Swagger

```bash
oapi-codegen -config api/oapi-codegen.yaml api/swagger.yaml
```


---

## 🛰️ gRPC

- Интерфейс описан в `api/pvz.proto`
- Сервер работает на порту `:3000`
- Не требует авторизации

---

## 🐳 Docker

```bash
docker-compose -f build/docker-compose.yml up --build
```

Пример `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=pvz
APP_PORT=8080
JWT_SECRET=jwt-secret
```

---

## 📋 CI/CD

- GitHub Actions: `.github/workflows/ci.yml`
- Генерация моков, юнит и интеграционные тесты, сборка контейнера

---

## 🗃️ Схема базы данных

Файл: [`docs/product.png`](./docs/product.png)

![db schema](docs/product.png)
