package main

import (
	"log"
	"os"

	"github.com/docmind/go-user-service/internal/config"
	"github.com/docmind/go-user-service/internal/handler"
	"github.com/docmind/go-user-service/internal/infrastructure"
	"github.com/docmind/go-user-service/internal/middleware"
	"github.com/docmind/go-user-service/internal/repository"
	"github.com/docmind/go-user-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	infrastructure.InitLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := infrastructure.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)

	// Initialize JWT manager
	jwtManager := infrastructure.NewJWTManager(cfg.JWTSecret)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.ErrorHandler())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	users := router.Group("/api/users")
	users.Use(middleware.AuthMiddleware(jwtManager))
	{
		users.GET("/me", userHandler.GetMe)
		users.PUT("/me", userHandler.UpdateMe)
		users.DELETE("/me", userHandler.DeleteMe)
	}

	// Start server
	port := cfg.ServerPort
	infrastructure.Logger.Info("Starting DocMind Go Service", "port", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
