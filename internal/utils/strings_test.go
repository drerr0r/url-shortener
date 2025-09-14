// internal/utils/strings_test.go

package utils

import (
	"testing"
)

// TestGenerateRandomString tests the GenerateRandomString function
// Этот тест проверяет основную функциональность генерации случайных строк
func TestGenerateRandomString(t *testing.T) {
	// Определяем тестовые случаи
	// Каждый случай имеет name, input и expected output
	tests := []struct {
		name    string // Название теста для читаемости
		length  int    // Входные данные: длинна строки
		wantErr bool   // Ожидаем ли мы ошибку?
	}{
		{"Valid lenght 6", 6, false},   // Корректная длина - ошибки быть не должно
		{"Valid lenght 10", 10, false}, // Другая корректная длина
		{"Zero lenght", 0, true},       // Нулевая длина - должна быть ошибка
		{"Negative lenght", -1, true},  // Отрицательная длина - ошибка
	}

	// Итерируем по всем тестовым случаям
	for _, tt := range tests {
		// t.Run создает подтест для каждого случая
		// Это позволяет видеть какие конкретно тесты проходят/падают
		t.Run(tt.name, func(t *testing.T) {
			// Вызываем тестируемую функцию
			result, err := GenerateRandomString(tt.length)

			// Проверяем, соответствует ли ошибка ожиданиям
			// (err != nil) != tt.wantErr означает:
			// - Если хотим ошибку (tt.wantErr = true), то err должен быть != nil
			// - Если не хотим ошибку (tt.wantErr = false), то укк должен быть = nil
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Если ошибки не ожидалось, то проверяем длинну результата
			if !tt.wantErr && len(result) != tt.length {
				t.Errorf("Expected lenght %d, got %d", tt.length, len(result))
			}
		})
	}

}

// TestGenerateRandomStringUniquness тестирует уникальность генерируемых строк
// Важно убедиться, что функция не генерирует дубликаты при многократном вызове
func TestGenerateRandomStringUniquness(t *testing.T) {
	// Используем [map] для отслеживания уникальности
	generated := make(map[string]bool)

	// Генерируем 100 строк и проверяем на уникальность
	for i := 0; i < 100; i++ {
		str, err := GenerateRandomString(8)
		if err != nil {
			t.Fatalf("GenerateRandomString failed: %v", err) // Fatlf для остановки теста
		}

		// Если строка уже была сгенерирована - это ошибка
		if generated[str] {
			t.Errorf("Duplicate string generated: %s", str)
		}
		generated[str] = true // Отмечаем строку как сгенерированую
	}

	// Дополнительная проверка: убедимся, что сгенерировано 100 уникальных строк
	// Это проверяет, что все итерации цикла отработали корректно
	if len(generated) != 100 {
		t.Errorf("Expected 100 unique strings, got %d", len(generated))
	}
}

// TestGenerateRandomStringCharset тестирует, что строка содержит только разрешенные символы
// Проверяем, что генератор использует только разрешенные символы
func TestGenerateRandomStringCharset(t *testing.T) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result, err := GenerateRandomString(20)
	if err != nil {
		t.Fatalf("GenerateRandomString failed: %v", err)
	}

	// Проверяем каждый стимвол в результате
	for _, char := range result {
		// Если символ не найден в разрешенном наборе - это ошибка
		if !containRune(charset, char) {
			t.Errorf("Invalid character in result: %c", char)
		}
	}
}

// containRune вспомогательная функция для проверки наличия символа в строке
// Используется для проверки валидности сгенерированных символов
func containRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
