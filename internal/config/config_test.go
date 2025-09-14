// internal/config/config_test.go

package config

import (
	"os"
	"testing"
)

func TestLoadDefaultValues(t *testing.T) {
	// Простой тест без очистки env
	cfg := Load()

	// Проверяем что значения загружаются (не важно какие)
	if cfg.Server.Port == "" {
		t.Error("Server port should not be empty")
	}

	if cfg.Database.Host == "" {
		t.Error("Database host should not be empty")
	}

	if cfg.App.ShortCodeLength <= 0 {
		t.Error("Short code length should be positive")
	}
}

func TestLoadEnvOverride(t *testing.T) {
	// Тестируем переопределение одной переменной
	originalPort := os.Getenv("SERVER_PORT")
	defer os.Setenv("SERVER_PORT", originalPort)

	os.Setenv("SERVER_PORT", "9090")

	cfg := Load()

	if cfg.Server.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Server.Port)
	}
}
