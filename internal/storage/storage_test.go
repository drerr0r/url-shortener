// internal/storage/storage_test.go

package storage

import (
	"context"
	"testing"
	"time"

	"github.com/drerr0r/url-shortner/internal/models"
)

// TestMockStorageCreateURL тестируем базовое создание ГКД - позитивный сценарий
func TestMockStorageCreateURL(t *testing.T) {
	mockStorage := NewMockStorage()
	ctx := context.Background()

	testURL := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
		CreatedAt:   time.Now(),
		ClickCount:  0,
	}

	err := mockStorage.CreateURL(ctx, testURL)
	if err != nil {
		t.Fatalf("CreateURL failed: %v", err)
	}

	// Проверяе, что URL сохранился
	if len(mockStorage.URLs) != 1 {
		t.Errorf("Expected 1 URL in storage, got %d", len(mockStorage.URLs))
	}
}

// TestMockStorageCreateURLDuplicate тестируем создание дубликата - должен возвращаться error
func TestMockStorageCreateURLDuplicate(t *testing.T) {
	mockStorage := NewMockStorage()
	ctx := context.Background()

	testURL := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "duplicate",
		CreatedAt:   time.Now(),
		ClickCount:  0,
	}

	// Первое создание - должно работать
	err := mockStorage.CreateURL(ctx, testURL)
	if err != nil {
		t.Fatalf("First CreateURL failed: %v", err)
	}

	// Второе создание с тем же short code - должно вернуть ошибку
	err = mockStorage.CreateURL(ctx, testURL)
	if err == nil {
		t.Errorf("Expected error for duplicate short code, got nil")
	}

}

// TestMockStorageGetByShortCode тестируем получение существующего URL
func TestMockStorageGetByShortCode(t *testing.T) {
	mockStorage := NewMockStorage()
	ctx := context.Background()

	testURL := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
		CreatedAt:   time.Now(),
		ClickCount:  5,
	}

	mockStorage.URLs["test123"] = testURL

	retrievedURL, err := mockStorage.GetURLByShortCode(ctx, "test123")
	if err != nil {
		t.Fatalf("GetURLByShortCode failed: %v", err)
	}

	if retrievedURL.OriginalURL != "https://example.com" {
		t.Errorf("Expected OriginalURL 'https://example.com', got '%s'", retrievedURL.OriginalURL)
	}

	if retrievedURL.ClickCount != 5 {
		t.Errorf("Expected ClickCount 5, got %d", retrievedURL.ClickCount)
	}
}

// TestMockStorageGetURLByShortCodeNotFound тестируем получение несуществующего URL
func TestMockStorageGetURLByShortCodeNotFound(t *testing.T) {
	mockStorage := NewMockStorage()
	ctx := context.Background()

	_, err := mockStorage.GetURLByShortCode(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestMockStorageIncrementClickCount тестируем увеличение счетчика кликов
func TestMockStorageIncrementClickCount(t *testing.T) {
	mockStorage := NewMockStorage()
	ctx := context.Background()

	testURL := &models.URL{
		ID:          1,
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
		CreatedAt:   time.Now(),
		ClickCount:  10,
	}

	mockStorage.URLs["test123"] = testURL

	err := mockStorage.IncrementClickCount(ctx, 1)
	if err != nil {
		t.Fatalf("IncrementClickCount failed: %v", err)
	}

	if testURL.ClickCount != 11 {
		t.Errorf("Expected ClickCount 11, got %d", testURL.ClickCount)
	}

}

// TestMockStorageIncrementClickCountNotFound тестируем увеличение счетчика для несуществующего URL
func TestMockStorageIncrementClickCountNotFound(t *testing.T) {
	mockStorage := NewMockStorage()
	ctx := context.Background()

	err := mockStorage.IncrementClickCount(ctx, 999)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestMockStorageGetURLStats тестируем получение статистики URL
func TestMockStorageGetURLStats(t *testing.T) {
	mockStorage := NewMockStorage()
	ctx := context.Background()

	createdTime := time.Now()
	testURL := &models.URL{
		ID:          1,
		OriginalURL: "https://example.com",
		ShortCode:   "test123",
		CreatedAt:   createdTime,
		ClickCount:  42,
	}

	mockStorage.URLs["test123"] = testURL

	stats, err := mockStorage.GetURLStats(ctx, "test123")
	if err != nil {
		t.Fatalf("GetURLStats failed: %v", err)
	}

	if stats.ShortCode != "test123" {
		t.Errorf("Expected ShortCode 'test123', got '%s'", stats.ShortCode)
	}

	if stats.OriginalURL != "https://example.com" {
		t.Errorf("Expected OriginalURL 'https://example.com', got '%s'", stats.OriginalURL)
	}

	if stats.ClickCount != 42 {
		t.Errorf("Expected ClickCount 42, got %d", stats.ClickCount)
	}

}

// TestMockStorageGetURLStatsNotFound тестируем получение статистики для несуществующего URL
func TestMockStorage_GetURLStats_NotFound(t *testing.T) {
	mockStorage := NewMockStorage()
	ctx := context.Background()

	_, err := mockStorage.GetURLStats(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestMockStorageClose тестируем закрытие хранилища (заглушка для mock)
func TestMockStorage_Close(t *testing.T) {
	mockStorage := NewMockStorage()

	err := mockStorage.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
