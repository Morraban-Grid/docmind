package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docmind/go-user-service/internal/config"
	"github.com/docmind/go-user-service/internal/handler"
	"github.com/docmind/go-user-service/internal/infrastructure"
	"github.com/docmind/go-user-service/internal/middleware"
	"github.com/docmind/go-user-service/internal/repository"
	"github.com/docmind/go-user-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Note: These tests require a running PostgreSQL database
// Set DATABASE_URL environment variable before running

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize database
	db, err := infrastructure.InitDB(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

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
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.ErrorHandler())

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

	return router
}

func TestAuthenticationFlow(t *testing.T) {
	router := setupTestRouter()

	// Test user registration
	t.Run("Register new user", func(t *testing.T) {
		registerReq := map[string]string{
			"email":    "test@example.com",
			"password": "testpassword123",
			"name":     "Test User",
		}
		body, _ := json.Marshal(registerReq)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should succeed or return conflict if user already exists
		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusConflict)
	})

	// Test user login
	t.Run("Login with valid credentials", func(t *testing.T) {
		loginReq := map[string]string{
			"email":    "test@example.com",
			"password": "testpassword123",
		}
		body, _ := json.Marshal(loginReq)

		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"])
		assert.NotNil(t, response["user"])
	})

	// Test login with invalid credentials
	t.Run("Login with invalid credentials", func(t *testing.T) {
		loginReq := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(loginReq)

		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestValidation(t *testing.T) {
	router := setupTestRouter()

	t.Run("Register with invalid email", func(t *testing.T) {
		registerReq := map[string]string{
			"email":    "invalid-email",
			"password": "testpassword123",
			"name":     "Test User",
		}
		body, _ := json.Marshal(registerReq)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Register with short password", func(t *testing.T) {
		registerReq := map[string]string{
			"email":    "test2@example.com",
			"password": "short",
			"name":     "Test User",
		}
		body, _ := json.Marshal(registerReq)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
