package main

import (
	"log"
	"strconv"

	"ahadpos-go/internal/config"
	"ahadpos-go/internal/handlers"
	"ahadpos-go/internal/middleware"
	"ahadpos-go/internal/repository"
	"ahadpos-go/internal/services"
	"ahadpos-go/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	barangRepo := repository.NewBarangRepository()

	// Initialize services
	authService := services.NewAuthService(cfg, userRepo)
	barangService := services.NewBarangService(barangRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	barangHandler := handlers.NewBarangHandler(barangService)

	// Create router
	router := gin.Default()

	// Apply CORS middleware
	if cfg.CORS.Enabled {
		router.Use(middleware.CORSMiddleware(
			cfg.CORS.AllowedOrigins,
			cfg.CORS.AllowedMethods,
			cfg.CORS.AllowedHeaders,
			cfg.CORS.MaxAge,
		))
	}

	// Recovery middleware
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "AhadPOS API is running",
		})
	})

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth routes
		protected.GET("/auth/me", authHandler.GetCurrentUser)

		// Product routes
		protected.GET("/barang", barangHandler.GetBarangList)
		protected.GET("/barang/:id", barangHandler.GetBarang)
		protected.POST("/barang", barangHandler.CreateBarang)
		protected.PUT("/barang/:id", barangHandler.UpdateBarang)
		protected.DELETE("/barang/:id", barangHandler.DeleteBarang)
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireAdmin())
	{
		// Admin-specific routes can be added here
	}

	// Start server
	log.Printf("Starting server on port %d...", cfg.Server.Port)
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
