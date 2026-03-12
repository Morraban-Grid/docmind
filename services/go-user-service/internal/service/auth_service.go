package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/docmind/go-user-service/internal/domain"
	"github.com/docmind/go-user-service/internal/infrastructure"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   domain.UserRepository
	jwtManager *infrastructure.JWTManager
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo domain.UserRepository, jwtManager *infrastructure.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return domain.NewValidationError("Email is required", nil)
	}

	// RFC 5322 simplified email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return domain.NewValidationError("Invalid email format", nil)
	}

	return nil
}

// ValidatePassword validates password requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return domain.NewValidationError("Password must be at least 8 characters", nil)
	}
	return nil
}

// Register registers a new user
func (s *AuthService) Register(req RegisterRequest) (*domain.User, error) {
	// Validate email
	if err := ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	// Validate password
	if err := ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Normalize email to lowercase
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check if email already exists
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		infrastructure.LogAuthAttempt(email, false, "email already exists")
		return nil, &domain.AppError{
			Code:       "EMAIL_ALREADY_EXISTS",
			Message:    "Email already registered",
			StatusCode: 409,
		}
	}

	// Create user
	user := &domain.User{
		Email:        email,
		PasswordHash: req.Password, // Will be hashed in repository
		Name:         strings.TrimSpace(req.Name),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, &domain.AppError{
				Code:       "EMAIL_ALREADY_EXISTS",
				Message:    "Email already registered",
				StatusCode: 409,
			}
		}
		infrastructure.LogError("failed to create user", err, map[string]interface{}{
			"email": email,
		})
		return nil, domain.NewInternalError("Failed to create user")
	}

	infrastructure.LogUserCreated(user.UserID, user.Email)
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		infrastructure.LogAuthAttempt(email, false, "user not found")
		return nil, domain.NewAuthenticationError("Invalid email or password")
	}

	// Verify password
	if repo, ok := s.userRepo.(*repository.PostgresUserRepository); ok {
		if err := repo.VerifyPassword(user.PasswordHash, req.Password); err != nil {
			infrastructure.LogAuthAttempt(email, false, "invalid password")
			return nil, domain.NewAuthenticationError("Invalid email or password")
		}
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.UserID, user.Email)
	if err != nil {
		infrastructure.LogError("failed to generate token", err, map[string]interface{}{
			"user_id": user.UserID,
		})
		return nil, domain.NewInternalError("Failed to generate authentication token")
	}

	infrastructure.LogAuthAttempt(email, true, "")

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}
