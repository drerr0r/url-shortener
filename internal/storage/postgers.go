// internal/storage/postgres.go

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/drerr0r/url-shortener/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

// PostgresStorage реализация Storage для PostgreSQL
type PostgresStorage struct {
	db *sqlx.DB
}

// NewPostgresStorage создает новый экземпляр PostgresStorage
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	// Подключаемся к базе данных
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &PostgresStorage{db: db}, nil
}

// CreateURL сохраняет новую сокращенную ссылку в базе данных
func (s *PostgresStorage) CreateURL(ctx context.Context, url *models.URL) error {
	query := `
		INSERT INTO urls (original_url, short_code, created_at, click_count)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`

	// выполняем запрос с возвратом ID
	err := s.db.QueryRowContext(ctx, query,
		url.OriginalURL,
		url.ShortCode,
		time.Now(),
		0,
	).Scan(&url.ID)

	if err != nil {
		return fmt.Errorf("failed to create URL: %w", err)
	}
	return nil
}

// GetURLByShortCode возвращает URL по короткому коду
func (s *PostgresStorage) GetURLByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	var url models.URL

	query := `SELECT id, original_url, short_code, created_at, click_count FROM urls WHERE short_code = $1`

	err := s.db.GetContext(ctx, &url, query, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	return &url, nil
}

// IncrementClickCount увеличивает счетчик кликов
func (s *PostgresStorage) IncrementClickCount(ctx context.Context, id int64) error {
	query := `UPDATE urls SET click_count = click_count + 1 WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// GetURLStats возвращает статистику по ссылке
func (s *PostgresStorage) GetURLStats(ctx context.Context, shortCode string) (*models.URLStats, error) {
	var stats models.URLStats

	query := `
	SELECT short_code, original_url, created_at, click_count
	FROM urls
	WHERE short_code = $1
	`

	err := s.db.GetContext(ctx, &stats, query, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get URL stats: %w", err)
	}
	return &stats, nil
}

// Close закрывает соединение с базой
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

// Ошибки хранилища
var (
	ErrNotFound = fmt.Errorf("not found")
)
