// internal/storage/storage.go

package storage

import (
	"context"

	"github.com/drerr0r/url-shortener/internal/models"
)

// Storage интерфейс определяет методы для работы с хранилищем данных
type Storage interface {
	// CreateURL сохраняет новую сокращенную ссылку
	CreateURL(ctx context.Context, url *models.URL) error

	// GetURLByShortCode возвращает оригинальный URL по короткому коду
	GetURLByShortCode(ctx context.Context, shortCode string) (*models.URL, error)

	// IncrementClickCount увеличивает счетчик кликов для ссылки
	IncrementClickCount(ctx context.Context, id int64) error

	// GetURLStats возвращает статистику по сокращенной ссылке
	GetURLStats(ctx context.Context, shortCode string) (*models.URLStats, error)

	// Close закрывает соединение с хранилищем
	Close() error
}
