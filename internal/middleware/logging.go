// internal/middleware/logging.go

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
)

// LoggingMiddleware добавляет логирование всех входящих запросов
// Это помогает отслеживать работу API и искать баги
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Засекаем время начала обработки запроса
		start := time.Now()

		// Продолжаем обработку запроса
		c.Next()

		// После обработки логируем информацию о запросе
		duration := time.Since(start)

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Int("status", c.Writer.Status()).
			Str("duration", duration.String()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("request processed")
	}
}
