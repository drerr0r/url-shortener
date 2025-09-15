// internal/handlers/url_handler_test.go

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drerr0r/url-shortener/internal/config"

	"github.com/drerr0r/url-shortener/internal/models"
	"github.com/drerr0r/url-shortener/internal/storage"
	"github.com/gin-gonic/gin"
)

// TestCreateShortURL тестирует CreateShortURL handler
// Этот тест проверяет корректное создание короткой ссылки через HTTP API
func TestCreateShortURL(t *testing.T) {
	// Setup: создаем mock storage и конфиг
	mockStorage := storage.NewMockStorage()
	cfg := &config.Config{}
	cfg.App.BaseURL = "http://localhost:8080"
	cfg.App.ShortCodeLength = 6

	handler := NewURLHandler(mockStorage, cfg)

	// Создаем Gin router для тестирования
	// Gin.TestMode отключает логгирование для чистого вывода тестов
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/urls", handler.CreateShortURL)

	// Test case: валидный URL
	requestBody := `{"url": "https://example.com"}`
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Записываем response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус код - должен быть 201 Created
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Парсим JSON
	var response models.CreateURLResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Проверяем, что short URL содержит base URL
	expectedPrefix := cfg.App.BaseURL + "/"
	if response.ShortURL[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Short URL should start with '%s', got '%s'", expectedPrefix, response.ShortURL)
	}

	// Проверяем, что original URL сохранился
	if response.OriginalURL != "https://example.com" {
		t.Errorf("Expected original URL 'https://example.com', got '%s'", response.OriginalURL)
	}

	//  Проверяем, что URL сохранился в  storage
	if len(mockStorage.URLs) != 1 {
		t.Errorf("Expected 1 URL in storage, got %d", len(mockStorage.URLs))
	}
}

// TestCreateShortURLInvalidURL tests invalid URL handling
// Проверяем обработку некоторых URL - API должен возвращать 400 Bad Request
func TestCreateShortURLInvalidURL(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	cfg := &config.Config{}
	handler := NewURLHandler(mockStorage, cfg)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/urls", handler.CreateShortURL)

	// Test case: не валидный URL
	requestBody := `{"url": "invalid-url"}`
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Должны получать 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid URL, got %d", w.Code)
	}

	// Проверяем, что в storage ничего не добавилось
	if len(mockStorage.URLs) != 0 {
		t.Errorf("Expected 0 URLs in storage for invalid request, got %d", len(mockStorage.URLs))
	}
}

// TestCreaqteShortURLEmptyBody проверяет обработку пустого тела запроса
func TestCreaqteShortURLEmptyBody(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	cfg := &config.Config{}
	handler := NewURLHandler(mockStorage, cfg)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/urls", handler.CreateShortURL)

	// Test case: пустое тело запроса
	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//  Должны получить 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expectedd status 400 for empty body, got %d", w.Code)
	}
}

// TestRedirectToOriginalURL tests redirection functionality
// Тестируем перенаправление по короткой ссылке - основной функционал сервиса
func TestRedirectToOriginalURL(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	cfg := &config.Config{}
	handler := NewURLHandler(mockStorage, cfg)

	// Предварительно сохраняем URL в storage
	// Это имитирует ситуацию, когда URL уже был создан ранее
	testURL := &models.URL{
		ID:          1,
		OriginalURL: "https://example.com",
		ShortCode:   "abc123",
		ClickCount:  0,
	}
	// Сохраняем через метод CreateURL
	ctx := context.Background()
	err := mockStorage.CreateURL(ctx, testURL)
	if err != nil {
		t.Fatalf("Failed to create test URL: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:shortCode", handler.RedirectToOriginalURL)

	// Важно: используем с параметром
	req, _ := http.NewRequest("GET", "/abc123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Должны получить redirect 302 Found
	if w.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Проверяем Location header - должен содержать оригинальный URL
	location := w.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("Expected Location header 'https://example.com', got '%s'", location)
	}

	// Проверяем, что счетчик кликов увеличился
	// Нужно получить обновленные данные из storage
	updatedURL, err := mockStorage.GetURLByShortCode(ctx, "abc123")
	if err != nil {
		t.Fatalf("Failed to get updated URL: %v", err)
	}

	if updatedURL.ClickCount != 1 {
		t.Errorf("Click count should be 1, got %d", updatedURL.ClickCount)
	}
}

// TestRedirectToOriginalURLNotFound проверяет обработку ситуации, когда короткий код не найден
func TestRedirectToOriginalURL_NotFound(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	cfg := &config.Config{}
	handler := NewURLHandler(mockStorage, cfg)

	// Storage пустой - URL не существует
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:shortCode", handler.RedirectToOriginalURL)

	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Должны получить 404 Not Found
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent short code, got %d", w.Code)
	}
}

// TestGetURLStats тестируем получение статистики по короткой ссылке
func TestGetURLStats(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	cfg := &config.Config{}
	handler := NewURLHandler(mockStorage, cfg)

	// Предварительно сохраняем URL со статистикой
	ctx := context.Background()
	testURL := &models.URL{
		ID:          1,
		OriginalURL: "https://example.com",
		ShortCode:   "stats123",
		ClickCount:  15,
	}
	err := mockStorage.CreateURL(ctx, testURL)
	if err != nil {
		t.Fatalf("Failed to create test URL: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/urls/:shortCode/stats", handler.GetURLStats)

	req, _ := http.NewRequest("GET", "/api/v1/urls/stats123/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Должны получить 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Парсим JSON
	var stats models.URLStats
	if err := json.Unmarshal(w.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Fail to parse stats response: %v", err)
	}

	// Проверяем данные статистики
	if stats.ShortCode != "stats123" {
		t.Errorf("Expected short code 'stats123', got '%s'", stats.ShortCode)
	}
	if stats.OriginalURL != "https://example.com" {
		t.Errorf("Expected original URL 'https://example.com', got '%s'", stats.OriginalURL)
	}
	if stats.ClickCount != 15 {
		t.Errorf("Expected click count 15, got %d", stats.ClickCount)
	}
}

// TestGetURLStatsNotFound тестируем получение статистики для несуществующей ссылки
func TestGetURLStatsNotFound(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	cfg := &config.Config{}
	handler := NewURLHandler(mockStorage, cfg)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/urls/:shortCode/stats", handler.GetURLStats)

	req, _ := http.NewRequest("GET", "/api/v1/urls/nonexistent/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Должны получить 404 Not Found
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent stats, got %d", w.Code)
	}
}
