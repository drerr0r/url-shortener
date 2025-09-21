package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drerr0r/url-shortener/internal/config"
	"github.com/drerr0r/url-shortener/internal/handlers"
	"github.com/drerr0r/url-shortener/internal/middleware"
	"github.com/drerr0r/url-shortener/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// applyMigrations автоматически применяет миграции базы данных при запуске
// Проверяет существование таблицы urls и применяет миграции если необходимо
func applyMigrations(db *sqlx.DB) error {
	// Проверяем существование таблицы urls
	var tableExists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'urls'
		)
	`).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	// Если таблица не существует, применяем миграции
	if !tableExists {
		log.Println("Database table 'urls' not found. Applying migrations...")

		// Применяем миграции с помощью goose
		if err := goose.SetDialect("postgres"); err != nil {
			return fmt.Errorf("failed to set dialect: %w", err)
		}

		if err := goose.Up(db.DB, "migrations"); err != nil {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}

		log.Println("Migrations applied successfully")
	} else {
		log.Println("Database table 'urls' already exists. Skipping migrations.")
	}

	return nil
}

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к базе данных
	db, err := sqlx.Connect("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	// Автоматическое применение миграций
	if err := applyMigrations(db); err != nil {
		log.Fatalf("Failed to apply database migrations: %v", err)
	}

	// Создание хранилища
	storage := storage.NewPostgresStorage(db)

	// Создание обработчиков
	urlHandler := handlers.NewURLHandler(storage)

	// Настройка роутера
	router := gin.Default()

	// 🔴 ДОБАВЛЕНО: Загрузка HTML шаблонов
	router.LoadHTMLGlob("templates/*")

	// Middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// 🔴 ДОБАВЛЕНО: Обработчик для главной страницы с HTML формой
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Маршруты API
	api := router.Group("/api/v1")
	{
		api.POST("/shorten", urlHandler.ShortenURLHandler)
		api.GET("/stats/:shortCode", urlHandler.GetURLStatsHandler)
	}

	router.GET("/:shortCode", urlHandler.RedirectHandler)

	// Health check с проверкой базы данных
	router.GET("/health", func(c *gin.Context) {
		// Проверка доступности базы данных
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "database unavailable",
				"error":  err.Error(),
			})
			return
		}

		// Проверка существования таблицы
		var tableExists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = 'urls'
			)
		`).Scan(&tableExists)

		if err != nil || !tableExists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "database schema error",
				"error":  "urls table does not exist",
			})
			return
		}

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
