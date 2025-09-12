// internal/models/url.go

package models

import (
	"time"
)

// URL прредставляет модель данных для сокращенной ссылки
type URL struct {
	ID          int64     `db:"id" json:"id"`                     // Уникальный идентификатор
	OriginalURL string    `db:"original_url" json:"original_url"` // Оригинальный URL
	ShortCode   string    `db:"short_code" json:"short_code"`     // Сокращенный код
	CrreatedAt  time.Time `db:"created_at" json:"created_at"`     // Время создания
	ClickCount  int64     `db:"click_count" json:"click_count"`   // Счетчик кликов

}

type CreateURLRequest struct {
	URL string `json:"url" binding:"required,url"` // URL для сокращеня
}

// CreateURLResponse представляет ответ с созданной сокращенной ссылкой
type CreateURLResponse struct {
	ShortURL    string `json:"short_url"`    // Полная сокращенная ссылка
	OriginalURL string `json:"original_url"` // Оригинальный URL
}

// URLStats представляет статистику по сокращенной ссылке
type URLStats struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CrreatedAt  time.Time `json:"created_at"`
	ClickCount  int64     `json:"click_count"`
}
