package storage

import (
	"database/sql"

	"github.com/drerr0r/url-shortener/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

// SaveURL сохраняет URL в базу данных
func (s *PostgresStorage) SaveURL(url *models.URL) error {
	query := `INSERT INTO urls (original_url, short_code) VALUES ($1, $2)`
	_, err := s.db.Exec(query, url.OriginalURL, url.ShortCode)
	return err
}

// GetURL возвращает URL по короткому коду
func (s *PostgresStorage) GetURL(shortCode string) (*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = $1`
	var url models.URL
	err := s.db.Get(&url, query, shortCode)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &url, err
}

// 🟡 ДОБАВЛЕНО: Реализация отсутствующего метода
// GetURLByOriginal возвращает URL по оригинальному URL
func (s *PostgresStorage) GetURLByOriginal(originalURL string) (*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE original_url = $1`
	var url models.URL
	err := s.db.Get(&url, query, originalURL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &url, err
}

// URLExists проверяет существование URL
func (s *PostgresStorage) URLExists(shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
	var exists bool
	err := s.db.Get(&exists, query, shortCode)
	return exists, err
}

// DeleteURL удаляет URL по короткому коду
func (s *PostgresStorage) DeleteURL(shortCode string) error {
	query := `DELETE FROM urls WHERE short_code = $1`
	_, err := s.db.Exec(query, shortCode)
	return err
}

// GetURLs возвращает все URL с пагинацией
func (s *PostgresStorage) GetURLs(limit, offset int) ([]*models.URL, error) {
	query := `SELECT id, original_url, short_code, created_at FROM urls ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	var urls []*models.URL
	err := s.db.Select(&urls, query, limit, offset)
	return urls, err
}

// GetURLsCount возвращает количество URL в базе
func (s *PostgresStorage) GetURLsCount() (int, error) {
	query := `SELECT COUNT(*) FROM urls`
	var count int
	err := s.db.Get(&count, query)
	return count, err
}
