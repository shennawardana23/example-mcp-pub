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

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shennawardana23/example-mcp-pub/internal/app/controller"
	"github.com/shennawardana23/example-mcp-pub/internal/app/middleware"
	"github.com/shennawardana23/example-mcp-pub/internal/app/repository"
	"github.com/shennawardana23/example-mcp-pub/internal/app/service"
	"github.com/shennawardana23/example-mcp-pub/internal/config"
	"github.com/shennawardana23/example-mcp-pub/internal/database"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Internal Developer Portal API
// @version 1.0
// @description Production-ready Internal Developer Portal with Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger
	logger := setupLogger(cfg)
	logger.Info("Starting Internal Developer Portal API")

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Database connected successfully")

	// Initialize Redis
	var redisClient *database.RedisClient
	redisClient, err = database.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Warnf("Failed to connect to Redis: %v (continuing without cache)", err)
		redisClient = nil
	} else {
		defer redisClient.Close()
		logger.Info("Redis connected successfully")
	}

	// Setup Gin mode
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repositories
	serviceRepo := repository.NewPostgresServiceRepository(db.DB)
	userRepo := repository.NewPostgresUserRepository(db.DB)

	// Initialize services
	var redisClientInstance *redis.Client
	if redisClient != nil {
		redisClientInstance = redisClient.Client
	}
	serviceService := service.NewServiceService(serviceRepo, redisClientInstance)
	authService := service.NewAuthService(userRepo, &cfg.JWT)

	// Initialize controllers
	serviceController := controller.NewServiceController(serviceService)
	authController := controller.NewAuthController(authService)
	healthController := controller.NewHealthController(db, redisClient)

	// Setup router
	router := setupRouter(cfg, logger, authService, serviceController, authController, healthController)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Infof("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func setupLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set format
	if cfg.Log.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	return logger
}

func setupRouter(
	cfg *config.Config,
	logger *logrus.Logger,
	authService service.AuthService,
	serviceController *controller.ServiceController,
	authController *controller.AuthController,
	healthController *controller.HealthController,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())

	// Health check (no auth required)
	router.GET("/health", healthController.Health)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Auth
			protected.GET("/auth/me", authController.Me)

			// Services
			protected.GET("/services", serviceController.List)
			protected.GET("/services/:id", serviceController.GetByID)
			protected.POST("/services", serviceController.Create)
			protected.PUT("/services/:id", serviceController.Update)
			protected.DELETE("/services/:id", serviceController.Delete)
		}
	}

	return router
}
