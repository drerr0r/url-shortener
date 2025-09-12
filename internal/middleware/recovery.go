// internal/middleware/recovery.go

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RecoveryMiddleware обрабатывает паники и предотвращает падение приложения
// Возвращает клиенту 500 ошибку вместо прекращения работы сервера
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Логируем ошибку
				log.Error().
					Interface("error", err).
					Str("path", c.Request.URL.Path).
					Msg("recovered from panic")

				// Отправляем клиенту ошибку 500
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})

				// Прерываем цепочку обработчиков
				c.Abort()
			}
		}()

		// Продолжаем обработку запроса
		c.Next()
	}
}
