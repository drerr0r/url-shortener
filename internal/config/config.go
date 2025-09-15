//internal/config/config.go

package config

import (
	"os"
	"strconv"
	"time"
)

// Config хранит все конфигурационные параметры приложения
type Config struct {
	Server struct {
		Port         string        // Порт сервера
		ReadTimeout  time.Duration // Таймаут чтения запроса
		WriteTimeout time.Duration // Таймаут записи ответа
		IdleTimeout  time.Duration // Таймаут бездействия
	}
	Database struct {
		Host            string        // Хост базы данных
		Port            string        // Порт базы данных
		Name            string        // Имя базы данных
		User            string        // Пользователь БД
		Password        string        // Пароль БД
		SSLMode         string        // Режим SSL
		MaxOpenConns    int           // Максимальное число открытых соединений
		MaxIdleConns    int           // Максимальное число idle соединений
		ConnMaxLifetime time.Duration // Максимальное время жизни соединения
	}
	App struct {
		BaseURL         string // Базовый URL для коротких ссылок
		ShortCodeLength int    // Длинна короткого кода
	}
}

// Load загружает конфигурацию из переменных окружения со значениями по умолчанию
func Load() *Config {
	var cfg Config

	// Настройки сервера
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.ReadTimeout = getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second)
	cfg.Server.WriteTimeout = getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second)
	cfg.Server.IdleTimeout = getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second)

	// Настройка базы данных
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.Name = getEnv("DB_NAME", "urlshortener")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")
	cfg.Database.MaxOpenConns = getIntEnv("DB_MAX_OPEN_CONNS", 25)
	cfg.Database.MaxIdleConns = getIntEnv("DB_MAX_IDLE_CONNS", 25)
	cfg.Database.ConnMaxLifetime = getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	// Настройки приложения
	cfg.App.BaseURL = getEnv("APP_BASE_URL", "http://localhost:8080")
	cfg.App.ShortCodeLength = getIntEnv("APP_SHORT_CODE_LENGHT", 6)

	return &cfg

}

// Вспомогательные функции для чтения переменных окружения
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
