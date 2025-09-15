// internal/handlers/url_handler.go

package handlers

import (
	"net/http"
	"net/url"

	"github.com/drerr0r/url-shortener/internal/config"
	"github.com/drerr0r/url-shortener/internal/models"
	"github.com/drerr0r/url-shortener/internal/storage"
	"github.com/drerr0r/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
)

// URLHandler обрабатывает HTTP запросы связанные с сокращением URL
type URLHandler struct {
	storage storage.Storage
	cfg     *config.Config
}

// NewURLHandler создает новый экземпляр URLHandler
// storage: интерфейс для работы с хранилищем данных
// cfg: конфигурация приложения
func NewURLHandler(storage storage.Storage, cfg *config.Config) *URLHandler {
	return &URLHandler{
		storage: storage,
		cfg:     cfg,
	}
}

// CreateShortURL обрабатывает POST запрос для создания сокращенной ссылки
// @Summary Создать сокращенную ссылку
// @Description Принимает оригинальный URL и возвращает сокращенную версию
// @Tags urls
// @Accept json
// @Produce json
// @Param request body models.CreateURLRequest true "URL для сокращения"
// @Success  201 {object} models.CreateURLResponse
// @Failure 400 {object} map[string]string "Не верный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/urls [post]
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req models.CreateURLRequest

	// Парсим JSON тело запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Проверяем что URL валидный
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	// Генерируем уникальный короткий код
	shortCode, err := utils.GenerateRandomString(h.cfg.App.ShortCodeLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gen rate short code"})
		return
	}

	// Создаем обьект URL для сохранения в БД
	url := &models.URL{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
	}

	// Сохраняем в базу данных
	ctx := c.Request.Context()
	if err := h.storage.CreateURL(ctx, url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	// Формируем ответ
	response := models.CreateURLResponse{
		ShortURL:    h.cfg.App.BaseURL + "/" + shortCode,
		OriginalURL: req.URL,
	}

	c.JSON(http.StatusCreated, response)

}

// RedirectToOriginalURL обрабатывает GRT запрос и перенаправляет на оригинальный URL
// @Summary Перенаправление по короткой ссылке
// @Description Перенаправляет пользователя на оригинальный URL по короткому коду
// @Tags urls
// @Param shortCode path string true "Короткий код ссылки"
// @Success 302 "Перенаправление на оригинальный URL"
// @Failure 404 {object} map[string]string "Ссылка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /{shortCode} [get]
func (h *URLHandler) RedirectToOriginalURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Ищем URL в базе данных
	ctx := c.Request.Context()
	url, err := h.storage.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		if err == storage.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Увеличиваем счетчик кликов
	if err := h.storage.IncrementClickCount(ctx, url.ID); err != nil {
		// Логируем ошибку, но не прерываем перенаправление
		// Можно добавить логирование через log.Error()
	}

	// Перенаправляем на оригинальный URL
	c.Redirect(http.StatusFound, url.OriginalURL)
}

// GetURLStats возвращает статистику по короткой ссылке
// @Summary Получить статистику ссылки
// @Description Возвращает информацию о количестве кликов и дате создания
// @Tags urls
// @Produce json
// @Param shortCode path string true "Коротки код ссылки"
// @Success 200 {object} models.URLStats
// @Failure 404 {object} map[string]string "Ссылка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/urls/{shortCode}/stats [get]
func (h *URLHandler) GetURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	ctx := c.Request.Context()
	stats, err := h.storage.GetURLStats(ctx, shortCode)
	if err != nil {
		if err == storage.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
