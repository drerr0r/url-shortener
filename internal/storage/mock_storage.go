// internal/storage/mock_storage.go

package storage

import (
	"context"
	"fmt"

	"github.com/drerr0r/url-shortener/internal/models"
)

// MockStorage - имитация хранилища данных для тестирования
// Вместо реальной базы данных используем [map] в памяти
// Это позволяет тестировать ллгику без внешних зависимостей
type MockStorage struct {
	URLs map[string]*models.URL // Храним URL по их короткому коду
}

// NewMockStorage - создает новый экземпляр mock хранилища
// Конструктор иннициализирует пустую [map] для хранения URL
func NewMockStorage() *MockStorage {
	return &MockStorage{
		URLs: make(map[string]*models.URL),
	}
}

// CreateURL сохраняет URL память [map]
// Имитирует поведение реального хранилища при создании записи
func (m *MockStorage) CreateURL(ctx context.Context, url *models.URL) error {
	// Проверяем не существует ли уже такой short code
	// Это имитирует constraint violation в реальной БД
	if _, exists := m.URLs[url.ShortCode]; exists {
		return fmt.Errorf("short code already exists")
	}
	// Сохраняем URL в [map]
	m.URLs[url.ShortCode] = url
	return nil
}

// GetURLByShortCode извлекает URL по short code из memory map
// Возвращает ErrNotFound если URL не существует
func (m *MockStorage) GetURLByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	url, exists := m.URLs[shortCode]
	if !exists {
		return nil, ErrNotFound // Возвращаем стандартную ошибку "not found"
	}
	return url, nil
}

// IncrementClickCount увеличивает счетчик кликов для конкретного url
// Ищем URL по ID и увеличиваем счетчик - имитация UPDATE запроса
func (m *MockStorage) IncrementClickCount(ctx context.Context, id int64) error {
	// Итерируем по всем URL чтобы найти по ID
	// В реальной БД это был бы UPDATE с WHERE по id
	for _, url := range m.URLs {
		if url.ID == id {
			url.ClickCount++
			return nil
		}
	}
	return ErrNotFound // Если URL с таким ID не найден
}

// GetURLStats возвращает статистику по short code
// Преобразуем модель URL в модель URLStats для ответа API
func (m *MockStorage) GetURLStats(ctx context.Context, shortCode string) (*models.URLStats, error) {
	url, exists := m.URLs[shortCode]
	if !exists {
		return nil, ErrNotFound
	}

	// Преобразуем URL в URLStats (только нужные поля для статистики)
	return &models.URLStats{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		CreatedAt:   url.CreatedAt,
		ClickCount:  url.ClickCount,
	}, nil
}

// Close заглушка для интерфейса - ничего не делает в mock
// В реальном хранилище здесь бы закрывалось соединение с БД
func (m *MockStorage) Close() error {
	return nil
}
