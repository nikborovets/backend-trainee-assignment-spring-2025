# Backend Trainee Assignment - Сервис ПВЗ

> Тестовое задание для прохождения отбора на стажировку в Авито (Весна 2025).

Проект сервиса для работы с Пунктами Выдачи Заказов (ПВЗ) для Авито, реализуемый с использованием принципов Чистой Архитектуры на Go.

## Структура проекта

Проект следует принципам Чистой Архитектуры с четкими слоями:
- **Entities**: Бизнес-объекты и правила
- **Use Cases**: Бизнес-логика приложения
- **Interfaces**: Адаптеры между слоями
- **Infrastructure**: Реализация всей инфраструктуры

## 🚀 Быстрый старт

1. **Клонируй репозиторий:**
   ```sh
   git clone https://github.com/nikborovets/backend-trainee-assignment-spring-2025.git
   cd backend-trainee-assignment-spring-2025
   ```

2. **Запусти всё через Docker:**
   ```sh
   docker compose up --build -d
   ```

3. **Проверь, что сервис жив:**
   ```sh
   curl -i http://localhost:8080/ping
   # Должно вернуть {"message":"pong"}
   ```

4. **Swagger/OpenAPI:**  
   Описание API — в файле `swagger.yaml` (можно открыть в Swagger Editor).

5. **Получить тестовый JWT:**
   ```sh
   curl -X POST http://localhost:8080/dummyLogin -H 'Content-Type: application/json' -d '{"role":"moderator"}'
   curl -X POST http://localhost:8080/dummyLogin -H 'Content-Type: application/json' -d '{"role":"pvz_staff"}'
   ```

6. **Примеры запросов:**
   - Создать ПВЗ:
     ```sh
     curl -X POST http://localhost:8080/pvz/ -H 'Authorization: Bearer <TOKEN>' -H 'Content-Type: application/json' -d '{"city":"Москва"}'
     ```
   - Получить список ПВЗ:
     ```sh
     curl -X GET http://localhost:8080/pvz/ -H 'Authorization: Bearer <TOKEN>'
     ```
   - Создать приёмку:
     ```sh
     curl -X POST http://localhost:8080/receptions/ -H 'Authorization: Bearer <TOKEN>' -H 'Content-Type: application/json' -d '{"pvzId":"<PVZ_ID>"}'
     ```
   - Добавить товар:
     ```sh
     curl -X POST http://localhost:8080/products/ -H 'Authorization: Bearer <TOKEN>' -H 'Content-Type: application/json' -d '{"type":"электроника","pvzId":"<PVZ_ID>"}'
     ```
   - Закрыть приёмку:
     ```sh
     curl -X POST http://localhost:8080/pvz/<PVZ_ID>/close_last_reception -H 'Authorization: Bearer <TOKEN>'
     ```
   - Удалить последний товар:
     ```sh
     curl -X POST http://localhost:8080/pvz/<PVZ_ID>/delete_last_product -H 'Authorization: Bearer <TOKEN>'
     ```

7. **Остановить сервис:**
   ```sh
   docker compose down
   ```

## Порты

- **HTTP API**: 8080
- **gRPC**: 3000
- **Prometheus**: 9000

## Запуск тестов и покрытие

Для корректного подсчёта покрытия production-кода при тестах в test/* используйте:

```sh
# Запуск всех тестов с подробным выводом
go test -v ./test/...

# Покрытие production-кода (всех пакетов internal) тестами из test/*
go test -coverpkg=./internal/... ./test/...

go test -coverpkg=./internal/... -coverprofile=cover.out ./test/...
go tool cover -func=cover.out

go tool cover -html=cover.out
```

> Если тесты будут перенесены в internal/*, достаточно будет обычного go test -cover ./... для покрытия.

**P.S.**
- Все переменные окружения и настройки — в `.env` (по умолчанию всё работает из коробки).
- Для CI/CD, Prometheus, gRPC — см. отдельные секции в README.
- Если что-то не работает — смотри логи: `docker compose logs --tail=100 app`
- Для Swagger UI: https://editor.swagger.io/ (загрузи swagger.yaml)

## Прогресс выполнения

### Базовый функционал

- [x] **Настройка проекта**
  - [x] Инициализация Git-репозитория
  - [x] Создание .gitignore
  - [x] Создание PUML-диаграммы архитектуры проекта (ЧА)
  - [x] Настройка структуры проекта в соответствии с чистой архитектурой
    - [x] Создание директорий: cmd/service, internal/entities, internal/usecases, internal/interfaces, internal/infrastructure, configs, test
    - [x] Инициализация go.mod (Go module)
    - [x] Настройка структуры проекта в соответствии с чистой архитектурой
  - [x] Настройка Docker и Docker Compose

- [x] **Domain Layer (Entities)**
  - [x] Определение моделей данных (User, PVZ, Reception, Product)
  - [x] Определение правил бизнес-логики

- [x] **Use Cases**
  - [x] CreatePVZ (создание ПВЗ)
  - [x] CreateReception (создание приёмки)
  - [x] AddProduct (добавление товара)
  - [x] DeleteLastProductUseCase (удаление последнего товара LIFO)
  - [x] CloseReceptionUseCase (закрытие приёмки)
  - [x] ListPVZsUseCase (листинг ПВЗ с фильтрами)
  - [x] DummyLoginUseCase (выдача токена по роли, без пароля)
  - [x] LoginUseCase (логин по email+пароль)
  - [x] RegisterUseCase (регистрация пользователя)
  - [x] Получение данных с фильтрацией (есть фильтры по дате и пагинация в /pvz)

- [x] **Interfaces Layer**
  - [x] Определение репозиториев (интерфейсы)
  - [x] Определение контроллеров
  - [x] Определение DTO
    - [x] Реализованы все основные DTO: UserDTO, PVZDTO, FullPVZDTO, ReceptionDTO, ProductDTO, ReceptionWithProductsDTO, RegisterRequest, LoginRequest, ListParams, AddProductRequest (internal/interfaces/dto.go)
    - [x] Все репозиторные интерфейсы: UserRepository, PVZRepository, ReceptionRepository, ProductRepository (internal/interfaces/repository.go)
    - [x] Все контроллеры: AuthController, PVZController, ReceptionController, ProductController (internal/interfaces/*_controller.go) — строго по .puml и ТЗ

- [x] **Infrastructure Layer**
  - [x] **Базы данных**
    - [x] Настройка PostgreSQL (миграции, FK, индексы, ограничения)
    - [x] Реализованы все PG-репозитории: User, PVZ, Reception, Product (internal/infrastructure/repositories)
    - [x] Используется Squirrel, без ORM, только чистый SQL
    - [x] Интеграционные тесты для каждого репозитория (test/infrastructure/repositories), покрытие >75%
    - [x] Структура файлов: всё разнесено по подпапкам, нет каши
  - [x] **HTTP сервер**
    - [x] Настройка HTTP сервера (Gin)
    - [x] Реализация эндпоинтов авторизации
      - [x] POST /dummyLogin
      - [x] POST /register
      - [x] POST /login
    - [x] Реализация эндпоинтов ПВЗ
      - [x] POST /pvz
      - [x] GET /pvz
      - [x] POST /pvz/{pvzId}/close_last_reception
      - [x] POST /pvz/{pvzId}/delete_last_product
    - [x] Реализация эндпоинтов приемки
      - [x] POST /receptions
    - [x] Реализация эндпоинтов товаров
      - [x] POST /products

- [x] **Тесты**
  - [x] Unit-тесты (покрытие > 75%)
  - [x] Интеграционный тест сценария работы с ПВЗ и товарами (ручной прогон — все сценарии проверены, curl-скрипты есть)

### Дополнительные задания

- [x] **Авторизация**
  - [x] Полная реализация авторизации (регистрация и логин)

- [ ] **GRPC**
  - [ ] Настройка gRPC сервера
  - [ ] Реализация метода GetPVZList

- [ ] **Метрики Prometheus**
  - [ ] Сбор технических метрик (количество запросов, время ответа)
  - [ ] Сбор бизнес-метрик (количество ПВЗ, приемок, товаров)
  - [ ] Настройка эндпоинта /metrics

- [ ] **Логирование**
  - [ ] Настройка централизованного логирования
  - [ ] Логирование ключевых операций

- [ ] **Кодогенерация**
  - [ ] Настройка генерации DTO по OpenAPI схеме
