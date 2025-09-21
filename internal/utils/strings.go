// internal/utils/strings.go

package utils

import (
	"crypto/rand" // 🔴 ИСПРАВЛЕНО: Заменен math/rand на crypto/rand
	"encoding/base64"
)

// GenerateRandomString генерирует случайную строку заданной длины
// 🔴 ИСПРАВЛЕНО: Заменен небезопасный math/rand на криптографически безопасный crypto/rand
func GenerateRandomString(length int) string {
	b := make([]byte, length)

	// БЫЛО: rand.Read(b) // math/rand - предсказуемо и небезопасно
	// СТАЛО: использование crypto/rand для генерации криптографически безопасных случайных значений
	_, err := rand.Read(b)
	if err != nil {
		// В продакшене следует использовать proper error handling
		panic("failed to generate random string: " + err.Error())
	}

	return base64.URLEncoding.EncodeToString(b)[:length]
}

// IsValidShortCode проверяет валидность короткого кода
func IsValidShortCode(code string) bool {
	if len(code) < 4 || len(code) > 12 {
		return false
	}

	for _, char := range code {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}

	return true
}

// TruncateString обрезает строку до указанной длины
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}
