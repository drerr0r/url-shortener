package main

import (
	"log"
	"net/http"

	"github.com/drerr0r/url-shortener/internal/config"
	"github.com/drerr0r/url-shortener/internal/handlers"
	"github.com/drerr0r/url-shortener/internal/middleware"
	"github.com/drerr0r/url-shortener/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к базе данных
	// 🟡 ИСПРАВЛЕНО: Используем метод GetDSN из конфига
	db, err := sqlx.Connect("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	// Создание хранилища
	storage := storage.NewPostgresStorage(db)

	// Создание обработчиков
	urlHandler := handlers.NewURLHandler(storage)

	// Настройка роутера
	router := gin.Default()

	// Middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// Маршруты
	api := router.Group("/api/v1")
	{
		api.POST("/shorten", urlHandler.ShortenURLHandler)
		api.GET("/stats/:shortCode", urlHandler.GetURLStatsHandler)
	}

	router.GET("/:shortCode", urlHandler.RedirectHandler)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Запуск сервера
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout,
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
