package config

import (
	"fmt"
	"os"
	"time"
)

// Config представляет конфигурацию приложения
type Config struct {
	ServerPort         string        `mapstructure:"SERVER_PORT"`
	ServerReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ServerIdleTimeout  time.Duration `mapstructure:"SERVER_IDLE_TIMEOUT"`

	DBHost            string        `mapstructure:"DB_HOST"`
	DBPort            string        `mapstructure:"DB_PORT"`
	DBName            string        `mapstructure:"DB_NAME"`
	DBUser            string        `mapstructure:"DB_USER"`
	DBPassword        string        `mapstructure:"DB_PASSWORD"`
	DBSSLMode         string        `mapstructure:"DB_SSLMODE"`
	DBMaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`

	AppBaseURL         string `mapstructure:"APP_BASE_URL"`
	AppShortCodeLength int    `mapstructure:"APP_SHORT_CODE_LENGTH"`
}

// 🟡 ИСПРАВЛЕНО: Добавляем функцию LoadConfig
// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		ServerReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		ServerWriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		ServerIdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),

		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBName:            getEnv("DB_NAME", "urlshortener"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "password"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
		DBConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),

		AppBaseURL:         getEnv("APP_BASE_URL", "http://localhost:8080"),
		AppShortCodeLength: getEnvAsInt("APP_SHORT_CODE_LENGTH", 6),
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsDuration получает переменную окружения как time.Duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// validateConfig проверяет валидность конфигурации
func validateConfig(cfg *Config) error {
	if cfg.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}
	if cfg.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if cfg.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if cfg.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}

	// Дополнительные проверки могут быть добавлены здесь

	return nil
}

// 🟡 ДОБАВЛЕНО: Вспомогательная функция для создания DSN
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}
