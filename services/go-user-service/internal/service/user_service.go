package service

import (
	"strings"

	"github.com/docmind/go-user-service/internal/domain"
	"github.com/docmind/go-user-service/internal/infrastructure"
)

// UserService handles user business logic
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, domain.NewNotFoundError("User not found")
		}
		infrastructure.LogError("failed to get user", err, map[string]interface{}{
			"user_id": userID,
		})
		return nil, domain.NewInternalError("Failed to retrieve user")
	}

	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(userID string, req UpdateUserRequest) (*domain.User, error) {
	// Get existing user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, domain.NewNotFoundError("User not found")
		}
		infrastructure.LogError("failed to get user for update", err, map[string]interface{}{
			"user_id": userID,
		})
		return nil, domain.NewInternalError("Failed to retrieve user")
	}

	// Update fields
	user.Name = strings.TrimSpace(req.Name)

	// Save to database
	if err := s.userRepo.UpdateUser(user); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, domain.NewNotFoundError("User not found")
		}
		infrastructure.LogError("failed to update user", err, map[string]interface{}{
			"user_id": userID,
		})
		return nil, domain.NewInternalError("Failed to update user")
	}

	infrastructure.LogUserUpdated(userID)
	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(userID string) error {
	if err := s.userRepo.DeleteUser(userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return domain.NewNotFoundError("User not found")
		}
		infrastructure.LogError("failed to delete user", err, map[string]interface{}{
			"user_id": userID,
		})
		return domain.NewInternalError("Failed to delete user")
	}

	infrastructure.LogUserDeleted(userID)
	return nil
}
