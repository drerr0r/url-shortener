// cmd/server/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drerr0r/url-shortener/internal/config"
	"github.com/drerr0r/url-shortener/internal/handlers"
	"github.com/drerr0r/url-shortener/internal/middleware"
	"github.com/drerr0r/url-shortener/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	fmt.Printf("Server Port: '%s'\n", cfg.Server.Port)
	fmt.Printf("Database config: %s@%s:%s/%s\n",
		cfg.Database.User,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)

	// Иннициализируем логгер
	// (можно добавить настройку zerolog или другого логгера)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Инициализируем хранилище
	storage, err := storage.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	// Настраиваем Gin router
	router := gin.Default()

	// Добавляем Middleware
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Инициализируем обработчики
	urlHandler := handlers.NewURLHandler(storage, cfg)

	// Настраиваем маршруты API
	api := router.Group("/api/v1")
	{
		api.POST("/urls", urlHandler.CreateShortURL)
		api.GET("/urls/:shortCode/stats", urlHandler.GetURLStats)
	}

	// Маршрут для перенаправления
	router.GET("/:shortCode", urlHandler.RedirectToOriginalURL)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Создаем HTTP сервер с настройками таймаутов
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Ожидаем сигналов дял graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Panicln("Shutting down server...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}

	log.Panicln("Server exited")

}
