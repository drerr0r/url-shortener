# Makefile для управления проектом URL Shortener

# Применение
BINARY_NAME=url-shortener
DOCKER_COMPOSE=docker_compose
GO=go
MIGRATIONS_DIR=migrations

.PHONY: help build run test clean migrate-up migrate-down docker-up docker-do wn docker-build

help: ## Показать помощь по командам
	@echo "Доступные команды"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {print " \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Собрать приложение
	$(GO) build -o $(BINARY_NAME) ./cmd/server

run: ## Запустить риложение локально
	$(GO) run ./cmd/server

test: ## Запустить тесты
	$(GO) test ./... -v

clean: ## Очистить скомпилированые файлы
	$(GO) clean
	rm -f $(BINARY_NAME)

migrate-up: ## Применить миграции базы данных
	goose -dir $(MIGRATIONS_DIR) postgres "postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable" up

migrate-down: ## Откатить последнюю миграцию
	goose -dir $(MIGRATIONS_DIR) postgres "postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable" down

migrate-status: ## Показать статус миграций
	goose -dir $(MIGRATIONS_DIR) postgres "postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable" status

docker-up: ## Запустить контейнер Docker
	$(DOCKER_COMPOSE) up -d

docker-down: ## Остановить контейнер Docker
	$(DOCKER_COMPOSE) down

docker-build: ## Собрать Docker образ
	$(DOCKER_COMPOSE) build

docker-logs: ## Показать логи контейнеров
	$(DOCKER_COMPOSE) logs -f

docker-restart: ## Перезапустить контейнеры
	$(DOCKER_COMPOSE) restart

lint: ## Запустить линтеры
	golangci-lint run

vendor: ##Скачать зависимости в vendor
	$(GO) mod vendor

swagger: ## Генерировать Swagger документацию (нежно добавить swagger)
	swag init -g cmd/server/main.go	

.PHONY: deploy
deploy: docker-build docker-up migrate-up ## Полный деплой: собрать запустить, применить миграции
	

