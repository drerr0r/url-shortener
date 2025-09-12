# URL Shortener Service

Микросервис для сокращения URL ссылок с REST API, PostgreSQL и Docker.

## Функциональность

- Создание коротких ссылок
- Перенаправление по коротким ссылкам
- Статистика кликов
- RESTful API
- Docker контейнеризация
- CI/CD с GitHub Actions

## Быстрый старт

### Локальная разработка

```bash
# Клонировать репозиторий
git clone <your-repo>
cd url-shortener

# Запустить контейнеры
make docker-up

# Применить миграции
make migrate-up

# Запустить приложение
make run