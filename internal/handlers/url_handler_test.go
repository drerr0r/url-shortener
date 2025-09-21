package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drerr0r/url-shortener/internal/models"
	"github.com/drerr0r/url-shortener/internal/storage"
	"github.com/gin-gonic/gin"
)

// TestShortenURLHandler тестирует ShortenURLHandler
func TestShortenURLHandler(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/shorten", handler.ShortenURLHandler)

	requestBody := `{"url": "https://example.com"}`
	req, _ := http.NewRequest("POST", "/api/v1/shorten", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response ShortenResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.ShortURL == "" {
		t.Error("Short URL should not be empty")
	}

	count := mockStorage.GetURLCount()
	if count != 1 {
		t.Errorf("Expected 1 URL in storage, got %d", count)
	}
}

// TestShortenURLHandlerInvalidURL tests invalid URL handling
func TestShortenURLHandlerInvalidURL(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/shorten", handler.ShortenURLHandler)

	requestBody := `{"url": "invalid-url"}`
	req, _ := http.NewRequest("POST", "/api/v1/shorten", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid URL, got %d. Body: %s", w.Code, w.Body.String())
	}

	count := mockStorage.GetURLCount()
	if count != 0 {
		t.Errorf("Expected 0 URLs in storage for invalid request, got %d", count)
	}
}

// TestShortenURLHandlerEmptyBody проверяет обработку пустого тела запроса
func TestShortenURLHandlerEmptyBody(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/shorten", handler.ShortenURLHandler)

	req, _ := http.NewRequest("POST", "/api/v1/shorten", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty body, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// TestRedirectHandler tests redirection functionality
func TestRedirectHandler(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	testURL := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "abc123",
	}
	err := mockStorage.SaveURL(testURL)
	if err != nil {
		t.Fatalf("Failed to create test URL: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:shortCode", handler.RedirectHandler)

	req, _ := http.NewRequest("GET", "/abc123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d. Body: %s", w.Code, w.Body.String())
	}

	location := w.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("Expected Location header 'https://example.com', got '%s'", location)
	}
}

// TestRedirectHandlerNotFound проверяет обработку ситуации, когда короткий код не найден
func TestRedirectHandlerNotFound(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:shortCode", handler.RedirectHandler)

	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent short code, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// TestRedirectHandlerInvalidCode проверяет обработку невалидного короткого кода
func TestRedirectHandlerInvalidCode(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/:shortCode", handler.RedirectHandler)

	req, _ := http.NewRequest("GET", "/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid short code, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// TestGetURLStatsHandler тестируем получение статистики по короткой ссылке
func TestGetURLStatsHandler(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	testURL := &models.URL{
		OriginalURL: "https://example.com",
		ShortCode:   "stats123",
	}
	err := mockStorage.SaveURL(testURL)
	if err != nil {
		t.Fatalf("Failed to create test URL: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/stats/:shortCode", handler.GetURLStatsHandler)

	req, _ := http.NewRequest("GET", "/api/v1/stats/stats123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var url models.URL
	if err := json.Unmarshal(w.Body.Bytes(), &url); err != nil {
		t.Fatalf("Failed to parse stats response: %v", err)
	}

	if url.ShortCode != "stats123" {
		t.Errorf("Expected short code 'stats123', got '%s'", url.ShortCode)
	}
	if url.OriginalURL != "https://example.com" {
		t.Errorf("Expected original URL 'https://example.com', got '%s'", url.OriginalURL)
	}
}

// TestGetURLStatsHandlerNotFound тестируем получение статистики для несуществующей ссылки
func TestGetURLStatsHandlerNotFound(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/stats/:shortCode", handler.GetURLStatsHandler)

	req, _ := http.NewRequest("GET", "/api/v1/stats/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent stats, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// TestGetURLStatsHandlerInvalidCode тестируем получение статистики для невалидного кода
func TestGetURLStatsHandlerInvalidCode(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/stats/:shortCode", handler.GetURLStatsHandler)

	req, _ := http.NewRequest("GET", "/api/v1/stats/inv@lid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid short code, got %d. Body: %s", w.Code, w.Body.String())
	}
}

// TestDuplicateURL тестирует создание дубликата URL
func TestDuplicateURL(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	handler := NewURLHandler(mockStorage)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/shorten", handler.ShortenURLHandler)

	requestBody1 := `{"url": "https://example.com"}`
	req1, _ := http.NewRequest("POST", "/api/v1/shorten", bytes.NewBufferString(requestBody1))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Errorf("First request should succeed, got %d", w1.Code)
	}

	var response1 ShortenResponse
	json.Unmarshal(w1.Body.Bytes(), &response1)
	firstShortCode := response1.ShortURL

	requestBody2 := `{"url": "https://example.com"}`
	req2, _ := http.NewRequest("POST", "/api/v1/shorten", bytes.NewBufferString(requestBody2))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Duplicate request should return 200, got %d", w2.Code)
	}

	var response2 ShortenResponse
	json.Unmarshal(w2.Body.Bytes(), &response2)

	if response2.ShortURL != firstShortCode {
		t.Errorf("Duplicate URL should return same short code, got %s vs %s", response2.ShortURL, firstShortCode)
	}

	count := mockStorage.GetURLCount()
	if count != 1 {
		t.Errorf("Should have only 1 URL in storage for duplicates, got %d", count)
	}
}
