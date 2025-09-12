// internal/utils/strings.go

package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomString генерирует случайную строку заданной длины
// Используется для создания коротких кодов ссылок
func GenerateRandomString(lenght int) (string, error) {
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
