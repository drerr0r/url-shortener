package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
	originalEnv := map[string]string{}
	keys := []string{
		"SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT",
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_SSLMODE",
		"DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "DB_CONN_MAX_LIFETIME",
		"APP_BASE_URL", "APP_SHORT_CODE_LENGTH",
	}

	for _, key := range keys {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key) // Очищаем для теста
	}
	defer func() {
		// Восстанавливаем оригинальные значения
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	t.Run("Default values", func(t *testing.T) {
		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "8080", cfg.ServerPort)
		assert.Equal(t, 10*time.Second, cfg.ServerReadTimeout)
		assert.Equal(t, 10*time.Second, cfg.ServerWriteTimeout)
		assert.Equal(t, 60*time.Second, cfg.ServerIdleTimeout)

		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, "5432", cfg.DBPort)
		assert.Equal(t, "urlshortener", cfg.DBName)
		assert.Equal(t, "postgres", cfg.DBUser)
		assert.Equal(t, "password", cfg.DBPassword)
		assert.Equal(t, "disable", cfg.DBSSLMode)
		assert.Equal(t, 25, cfg.DBMaxOpenConns)
		assert.Equal(t, 25, cfg.DBMaxIdleConns)
		assert.Equal(t, 5*time.Minute, cfg.DBConnMaxLifetime)

		assert.Equal(t, "http://localhost:8080", cfg.AppBaseURL)
		assert.Equal(t, 6, cfg.AppShortCodeLength)
	})

	t.Run("Custom values", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "9090")
		os.Setenv("DB_HOST", "test-host")
		os.Setenv("DB_NAME", "test-db")
		os.Setenv("DB_USER", "test-user")
		os.Setenv("DB_PASSWORD", "test-password")

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "9090", cfg.ServerPort)
		assert.Equal(t, "test-host", cfg.DBHost)
		assert.Equal(t, "test-db", cfg.DBName)
		assert.Equal(t, "test-user", cfg.DBUser)
		assert.Equal(t, "test-password", cfg.DBPassword)
	})

	t.Run("Invalid duration", func(t *testing.T) {
		os.Setenv("SERVER_READ_TIMEOUT", "invalid-duration")

		cfg, err := LoadConfig()
		assert.NoError(t, err) // Должен вернуть значение по умолчанию, а не ошибку
		assert.Equal(t, 10*time.Second, cfg.ServerReadTimeout)
	})

	t.Run("Invalid integer", func(t *testing.T) {
		os.Setenv("DB_MAX_OPEN_CONNS", "not-a-number")

		cfg, err := LoadConfig()
		assert.NoError(t, err) // Должен вернуть значение по умолчанию, а не ошибку
		assert.Equal(t, 25, cfg.DBMaxOpenConns)
	})
}

func TestGetDSN(t *testing.T) {
	cfg := &Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "password",
		DBName:     "urlshortener",
		DBSSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=postgres password=password dbname=urlshortener sslmode=disable"
	assert.Equal(t, expected, cfg.GetDSN())
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				ServerPort: "8080",
				DBHost:     "localhost",
				DBName:     "testdb",
				DBUser:     "user",
			},
			wantErr: false,
		},
		{
			name: "Missing server port",
			config: &Config{
				ServerPort: "",
				DBHost:     "localhost",
				DBName:     "testdb",
				DBUser:     "user",
			},
			wantErr: true,
		},
		{
			name: "Missing DB host",
			config: &Config{
				ServerPort: "8080",
				DBHost:     "",
				DBName:     "testdb",
				DBUser:     "user",
			},
			wantErr: true,
		},
		{
			name: "Missing DB name",
			config: &Config{
				ServerPort: "8080",
				DBHost:     "localhost",
				DBName:     "",
				DBUser:     "user",
			},
			wantErr: true,
		},
		{
			name: "Missing DB user",
			config: &Config{
				ServerPort: "8080",
				DBHost:     "localhost",
				DBName:     "testdb",
				DBUser:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
