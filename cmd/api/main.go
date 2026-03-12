package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/config"
	"github.com/nielwyn/inventory-system/internal/database"
	"github.com/nielwyn/inventory-system/internal/handlers"
	"github.com/nielwyn/inventory-system/internal/middleware"
	"github.com/nielwyn/inventory-system/internal/repository"
	"github.com/nielwyn/inventory-system/internal/service"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"github.com/nielwyn/inventory-system/pkg/validator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Encoding); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Go Inventory System API")

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	db, err := database.New(cfg.Database.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Run database migrations
	if err := db.AutoMigrate(); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Register custom validators
	validator.RegisterCustomValidations()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	inventoryRepo := repository.NewInventoryRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	inventoryService := service.NewInventoryService(inventoryRepo)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(authService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	// Setup router
	router := setupRouter(healthHandler, authHandler, inventoryHandler, authService)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped")
}

// setupRouter configures all routes and middleware
func setupRouter(
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	inventoryHandler *handlers.InventoryHandler,
	authService service.AuthService,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Health check endpoints (no authentication required)
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	// Metrics endpoint (Prometheus)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth endpoints (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Inventory endpoints (protected)
		inventory := v1.Group("/inventory")
		inventory.Use(middleware.Auth(authService))
		{
			inventory.POST("/items", inventoryHandler.CreateItem)
			inventory.GET("/items", inventoryHandler.GetAllItems)
			inventory.GET("/items/:id", inventoryHandler.GetItemByID)
			inventory.PUT("/items/:id", inventoryHandler.UpdateItem)
			inventory.DELETE("/items/:id", inventoryHandler.DeleteItem)
		}
	}

	return router
}
