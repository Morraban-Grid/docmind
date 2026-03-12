package main

import (
	"log"
	"os"

	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/config"
	httphandler "github.com/Morraban-Grid/docmind/services/go-user-service/internal/handler/http"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/infrastructure"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/middleware"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/repository/postgres"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/service"
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

	// Initialize MinIO client
	minioClient, err := infrastructure.NewMinIOClient(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ROOT_USER"),
		os.Getenv("MINIO_ROOT_PASSWORD"),
		os.Getenv("MINIO_BUCKET"),
		os.Getenv("MINIO_USE_SSL") == "true",
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	documentRepo := postgres.NewDocumentRepository(db)

	// Initialize JWT manager
	jwtManager := infrastructure.NewJWTManager(cfg.JWTSecret)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo)
	documentService := service.NewDocumentService(documentRepo, minioClient)

	// Initialize handlers
	authHandler := httphandler.NewAuthHandler(authService)
	userHandler := httphandler.NewUserHandler(userService)
	documentHandler := httphandler.NewDocumentHandler(documentService)
	healthHandler := httphandler.NewHealthHandler()

	// Setup router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Health check endpoint
	router.GET("/health", healthHandler.Health)

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

	// Document routes (protected)
	documents := router.Group("/api/documents")
	documents.Use(middleware.AuthMiddleware(jwtManager))
	{
		documents.POST("", documentHandler.UploadDocument)
		documents.GET("", documentHandler.ListDocuments)
		documents.GET("/:id", documentHandler.GetDocument)
		documents.GET("/:id/download", documentHandler.DownloadDocument)
		documents.DELETE("/:id", documentHandler.DeleteDocument)
	}

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	infrastructure.LogInfo("Starting DocMind Go Service", map[string]interface{}{
		"port": port,
	})

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
