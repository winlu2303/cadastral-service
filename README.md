# Cadastral Service

Микросервис для обработки кадастровых запросов с эмуляцией внешнего сервера.

# Особенности

- Обработка кадастровых запросов с координатами
- Эмуляция внешнего сервера с задержкой до 60 секунд
- Сохранение истории запросов в PostgreSQL
- REST API с валидацией данных
- JWT аутентификация (опционально)
- Документированный код
- Готовые Docker контейнеры
- Полная тестовая среда

# Требования

- Docker и Docker Compose
- Go 1.21+ (для локальной разработки)

# Доступность сервисов по адресам

- Основной API: http://localhost:8080
- Mock сервер: http://localhost:8081
- Adminer (админка БД): http://localhost:8082

# Клонируйте проект (если нужно)
git clone <your-repo>
cd cadastral-service

# Быстрый старт
make quickstart

# Проверьте, что всё работает
make health

# Протестируйте API
make test-api

# Начать разработку
make setup

# Запустить в Docker
make docker-up

# Проверить логи
make docker-logs-api

# Внести изменения, затем пересобрать
make docker-restart

# Протестировать изменения
make test-api

# Запустить тесты
make test

# Структура проекта

cadastral-service/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── mock-server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handler.go
│   │   ├── middleware.go
│   │   └── routes.go
│   ├── models/
│   │   └── models.go
│   ├── repository/
│   │   └── repository.go
│   ├── service/
│   │   └── service.go
│   └── config/
│       └── config.go
├── migrations/
│   └── 001_init.sql
├── pkg/
│   ├── database/
│   │   └── database.go
│   └── logger/
│       └── logger.go
├── docker-compose.yml
├── Dockerfile
├── Dockerfile.mock
├── Makefile
├── go.mod
├── go.sum
├── README.md
└── test/
    └── api_test.go

