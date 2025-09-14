// internal/utils/strings.go

package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateRandomString генерирует случайную строку заданной длины
// Используется для создания коротких кодов ссылок
// Важно: длина должна быть положительным числом
func GenerateRandomString(lenght int) (string, error) {
	if lenght <= 0 {
		return "", fmt.Errorf("lenght must be positive, got %d", lenght)
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, lenght)
	for i := 0; i < lenght; i++ {
		// Генерируем случайный индекс в диапазоне charset
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
