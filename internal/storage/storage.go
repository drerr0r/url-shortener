package storage

import (
	"errors"

	"github.com/drerr0r/url-shortener/internal/models"
)

// 🟡 ДОБАВЛЕНО: Определяем ошибку для отсутствующих записей
var ErrNotFound = errors.New("record not found")

// Storage интерфейс для работы с хранилищем URL
type Storage interface {
	SaveURL(url *models.URL) error
	GetURL(shortCode string) (*models.URL, error)
	GetURLByOriginal(originalURL string) (*models.URL, error) // 🟡 ДОБАВЛЕНО: отсутствующий метод
	URLExists(shortCode string) (bool, error)
	DeleteURL(shortCode string) error
	GetURLs(limit, offset int) ([]*models.URL, error)
	GetURLsCount() (int, error)
}
