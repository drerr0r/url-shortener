# URL Shortener

Микросервис для сокращения URL с использованием Go, PostgreSQL и Gin.

## 🚀 Быстрый старт

### Требования
- Go 1.25+
- PostgreSQL 12+
- Docker (опционально)

### Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/drerr0r/url-shortener.git
cd url-shortener
Установите зависимости:

bash
go mod download
Настройте переменные окружения:

bash
cp .env.example .env
# Отредактируйте .env файл под вашу конфигурацию
Запустите базу данных:

bash
docker-compose up postgres -d
Запустите приложение:

bash
go run ./cmd/server
Приложение автоматически применит миграции базы данных при первом запуске.

🌐 Веб-интерфейс
Откройте в браузере: http://localhost:8080/

Введите URL в поле ввода

Нажмите "Сократить"

Скопируйте короткую ссылку

🔌 API Endpoints
Сокращение URL
bash
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
Перенаправление
bash
curl -I http://localhost:8080/abc123
Статистика
bash
curl http://localhost:8080/api/v1/stats/abc123
Health Check
bash
curl http://localhost:8080/health
🐛 Устранение проблем
Ошибка "relation 'urls' does not exist"
Приложение автоматически создает таблицы при первом запуске.

Ошибки подключения к БД
Проверьте настройки в .env файле:

DB_HOST, DB_PORT, DB_NAME

DB_USER, DB_PASSWORD

📝 Переменные окружения (.env)
env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=urlshortener
DB_USER=postgres
DB_PASSWORD=password
DB_SSLMODE=disable
APP_BASE_URL=http://localhost:8080
🛠️ Команды разработки
bash
# Тесты
go test ./... -v

# Сборка
go build -o url-shortener ./cmd/server

# Форматирование
go fmt ./...

# Запуск с Docker
docker-compose up --build
Приложение готово к использованию! 🎉