package utils

import (
	"strings"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"Length 6", 6},
		{"Length 8", 8},
		{"Length 12", 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 🟡 ИСПРАВЛЕНО: Убрали лишнюю переменную ошибки
			result := GenerateRandomString(tt.length)
			if len(result) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(result))
			}

			// Проверяем, что строка содержит только допустимые символы
			for _, char := range result {
				if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_", char) {
					t.Errorf("Invalid character in random string: %c", char)
				}
			}
		})
	}
}

func TestGenerateRandomString_Unique(t *testing.T) {
	// 🟡 ИСПРАВЛЕНО: Убрали лишнюю переменную ошибки
	str1 := GenerateRandomString(6)
	str2 := GenerateRandomString(6)

	if str1 == str2 {
		t.Errorf("Generated identical strings: %s and %s", str1, str2)
	}
}

func TestIsValidShortCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"Valid code", "abc123", true},
		{"Valid with dash", "abc-123", true},
		{"Valid with underscore", "abc_123", true},
		{"Too short", "abc", false},
		{"Too long", "abcdefghijklm", false},
		{"Invalid chars", "abc@123", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidShortCode(tt.code)
			if result != tt.expected {
				t.Errorf("IsValidShortCode(%q) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected string
	}{
		{"Shorter than limit", "hello", 10, "hello"},
		{"Exactly limit", "hello", 5, "hello"},
		{"Longer than limit", "hello world", 5, "hello"},
		{"Zero length", "hello", 0, ""},
		{"Empty string", "", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateString(tt.input, tt.length)
			if result != tt.expected {
				t.Errorf("TruncateString(%q, %d) = %q, expected %q", tt.input, tt.length, result, tt.expected)
			}
		})
	}
}
