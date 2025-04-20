# PVZ Management Backend

Сервис для управления пунктами выдачи заказов (ПВЗ), приёмкой товаров и учётом продукции. Поддерживает регистрацию и авторизацию пользователей (модератор / сотрудник ПВЗ), создание ПВЗ, приёмок и добавление/удаление товаров в приёмке.

---

##  Стек технологий

- **Язык**: Go 1.24
- **База данных**: PostgreSQL
- **Фреймворк**: `net/http`, `gorilla/mux`
- **Авторизация**: JWT
- **Документация API**: OpenAPI 3.0 ([swagger.yaml](./swagger.yaml))
- **Тесты**:
    - Unit-тесты: с использованием `gomock`
    - Интеграционные тесты: с `httptest.NewServer`

---

## Структура проекта

```
.
├── api/                    # swagger
├── cmd/                    # Точка входа (main)
├── internal/
│   ├── adapters/
│   │   ├── api/            # Контроллеры REST API
│   │   └── db/             # Репозитории (PostgreSQL)
│   ├── app/                # Setup маршрутов
│   ├── domain/             # Интерфейсы и бизнес-модели
│   └── config/             # Конфигуратор приложения
├── init/migrations/        # SQL-миграции
├── tests/                  # Unit и Integration тесты
├── swagger.yaml            # OpenAPI спецификация
└── build/  # Сборка контейнера
```

---

## Переменные окружения

| Переменная     | Назначение         | Значение по умолчанию |
|----------------|--------------------|-----------------------|
| `DB_HOST`      | Хост PostgreSQL     | `localhost`           |
| `DB_PORT`      | Порт PostgreSQL     | `5432`                |
| `DB_USER`      | Пользователь БД     | `postgres`            |
| `DB_PASSWORD`  | Пароль БД           | `password`            |
| `DB_NAME`      | Название БД         | `pvz`                 |
| `APP_PORT`     | Порт API сервера     | `8080`                |
| `JWT_SECRET`   | Секрет для токенов  | `jwt-secret`          |

---

## Тестирование

### Юнит-тесты

```bash
go test ./tests/unit/... -v
```

### Интеграционные тесты

Интеграционные тесты используют `httptest.NewServer` и запускают приложение с реальной БД.

```bash
go test ./tests/integration/... -v -count=1
```

> Убедись, что PostgreSQL поднят на `localhost:5432` перед запуском!

---

##  Docker

###  Сборка

```bash
docker build -t pvz-service .
```

###  Запуск

```bash
docker run -p 8080:8080 --env-file .env pvz-service
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

##  Авторизация

Все защищённые маршруты требуют JWT. Токен можно получить через:

```http
POST /dummyLogin
{
  "role": "moderator"
}
```

Ответ:

```json
"eyJhbGciOiJIUzI1NiIsInR5cCI6..."
```

---

##  API документация

Файл [swagger.yaml](./api/swagger.yaml) описывает все маршруты API.

---

##  Основной функционал

- Регистрация и логин пользователей
- Создание ПВЗ
- Получение списка ПВЗ с фильтрацией
- Создание приёмки товаров
- Закрытие приёмки
- Добавление / удаление товаров

---

##  CI/CD

Проект поддерживает CI на GitHub Actions:
- генерация mock'ов
- запуск unit и интеграционных тестов
- сборка Docker-образа



