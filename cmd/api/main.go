package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rare-merch-exchange/internal/delivery/http/v1"
	"rare-merch-exchange/internal/pkg/config"
	"rare-merch-exchange/internal/repository"
	"rare-merch-exchange/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	"go.uber.org/zap"

	_ "rare-merch-exchange/docs"
)

// @title Rare Merch Exchange API
// @version 1.0
// @description API для биржи раритетного мерча с аукционами, лотереями и конкурсами
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Инициализируем конфигурацию
	cfg := config.Load()

	// Инициализируем логгер
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// Подключаемся к базе данных
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Подключаемся к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Проверяем подключение к Redis
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	// Инициализируем репозитории
	repos := repository.NewRepositories(db, rdb)

	// Инициализируем use cases
	useCases := usecase.NewUseCases(repos, cfg, logger)

	// Инициализируем HTTP handlers
	handlers := v1.NewHandlers(useCases, logger)

	// Настраиваем Gin
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Swagger документация
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// Настраиваем маршруты
	v1.SetupRoutes(router, handlers)

	// Запускаем сервер
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.String("port", cfg.Server.Port))

	// Ожидаем сигнал для завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Даём время на завершение активных запросов
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
