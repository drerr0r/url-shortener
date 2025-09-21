package handlers

import (
	"net/http"
	"net/url"

	"github.com/drerr0r/url-shortener/internal/models"
	"github.com/drerr0r/url-shortener/internal/storage"
	"github.com/drerr0r/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type URLHandler struct {
	storage storage.Storage
}

func NewURLHandler(storage storage.Storage) *URLHandler {
	return &URLHandler{storage: storage}
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// ShortenURLHandler обрабатывает запрос на сокращение URL
func (h *URLHandler) ShortenURLHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1024)

	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if !isValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	existingURL, err := h.storage.GetURLByOriginal(req.URL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check existing URL")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if existingURL != nil {
		c.JSON(http.StatusOK, ShortenResponse{ShortURL: existingURL.ShortCode})
		return
	}

	shortCode := utils.GenerateRandomString(6)

	if !utils.IsValidShortCode(shortCode) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate valid short code"})
		return
	}

	urlModel := &models.URL{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
	}

	if err := h.storage.SaveURL(urlModel); err != nil {
		log.Error().Err(err).Msg("Failed to save URL")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL"})
		return
	}

	c.JSON(http.StatusCreated, ShortenResponse{ShortURL: shortCode})
}

// 🟡 ИСПРАВЛЕНО: Переименовали метод для соответствия вызовам в main.go
// RedirectHandler обрабатывает перенаправление по короткому URL
func (h *URLHandler) RedirectHandler(c *gin.Context) {
	shortCode := c.Param("shortCode")

	if !utils.IsValidShortCode(shortCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid short code format"})
		return
	}

	url, err := h.storage.GetURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusFound, url.OriginalURL)
}

// 🟡 ИСПРАВЛЕНО: Переименовали метод для соответствия вызовам в main.go
// GetURLStatsHandler возвращает статистику по URL
func (h *URLHandler) GetURLStatsHandler(c *gin.Context) {
	shortCode := c.Param("shortCode")

	if !utils.IsValidShortCode(shortCode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid short code format"})
		return
	}

	url, err := h.storage.GetURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, url)
}

func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}
