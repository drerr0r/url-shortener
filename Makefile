# Makefile для управления проектом URL Shortener

# Переменные
BINARY_NAME=url-shortener
DOCKER_COMPOSE=docker-compose
GO=go
MIGRATIONS_DIR=migrations

.PHONY: help build run clean migrate-up migrate-down docker-up docker-down docker-build docker-logs docker-restart lint vendor swagger deploy test test-cover test-cover-html test-handlers test-utils test-storage

help: ## Показать помощь по командам
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Собрать приложение
	$(GO) build -o $(BINARY_NAME) ./cmd/server

run: ## Запустить приложение локально
	$(GO) run ./cmd/server

test: ## Запустить тесты
	$(GO) test ./... -v

clean: ## Очистить скомпилированные файлы
	$(GO) clean
	rm -f $(BINARY_NAME)

# 🟡 ИСПРАВЛЕНО: Использование переменной окружения для DB_URL
# БЫЛО: goose -dir $(MIGRATIONS_DIR) postgres "postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable" up
migrate-up: ## Применить миграции базы данных
	goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" up

migrate-down: ## Откатить последнюю миграцию
	goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" down

migrate-status: ## Показать статус миграций
	goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" status

docker-up: ## Запустить контейнеры Docker
	$(DOCKER_COMPOSE) up -d

docker-down: ## Остановить контейнеры Docker
	$(DOCKER_COMPOSE) down

docker-build: ## Собрать Docker образ
	$(DOCKER_COMPOSE) build

docker-logs: ## Показать логи контейнеров
	$(DOCKER_COMPOSE) logs -f

docker-restart: ## Перезапустить контейнеры
	$(DOCKER_COMPOSE) restart

lint: ## Запустить линтеры
	golangci-lint run

vendor: ## Скачать зависимости в vendor
	$(GO) mod vendor

swagger: ## Генерировать Swagger документацию
	swag init -g cmd/server/main.go

deploy: docker-build docker-up migrate-up ## Полный деплой: собрать, запустить, применить миграции

test-cover: ## Run tests with coverage report
	$(GO) test ./... -cover

test-cover-html: ## Generate HTML coverage report
	$(GO) test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser"

test-handlers: ## Run only handler tests
	$(GO) test ./internal/handlers/ -v

test-utils: ## Run only utils tests
	$(GO) test ./internal/utils/ -v

test-storage: ## Run only storage tests
	$(GO) test ./internal/storage/ -v