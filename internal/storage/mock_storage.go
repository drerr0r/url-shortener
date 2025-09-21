package storage

import (
	"github.com/drerr0r/url-shortener/internal/models"
)

// MockStorage реализация Storage для тестов
type MockStorage struct {
	urls map[string]*models.URL
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		urls: make(map[string]*models.URL),
	}
}

// GetURLCount возвращает количество URL (для тестов)
func (m *MockStorage) GetURLCount() int {
	return len(m.urls)
}

// GetAllURLs возвращает все URL (для тестов)
func (m *MockStorage) GetAllURLs() map[string]*models.URL {
	return m.urls
}

func (m *MockStorage) SaveURL(url *models.URL) error {
	m.urls[url.ShortCode] = url
	return nil
}

func (m *MockStorage) GetURL(shortCode string) (*models.URL, error) {
	url, exists := m.urls[shortCode]
	if !exists {
		return nil, ErrNotFound
	}
	return url, nil
}

func (m *MockStorage) GetURLByOriginal(originalURL string) (*models.URL, error) {
	for _, url := range m.urls {
		if url.OriginalURL == originalURL {
			return url, nil
		}
	}
	return nil, nil
}

func (m *MockStorage) URLExists(shortCode string) (bool, error) {
	_, exists := m.urls[shortCode]
	return exists, nil
}

func (m *MockStorage) DeleteURL(shortCode string) error {
	delete(m.urls, shortCode)
	return nil
}

func (m *MockStorage) GetURLs(limit, offset int) ([]*models.URL, error) {
	var result []*models.URL
	count := 0
	for _, url := range m.urls {
		if count >= offset {
			result = append(result, url)
		}
		count++
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

func (m *MockStorage) GetURLsCount() (int, error) {
	return len(m.urls), nil
}
