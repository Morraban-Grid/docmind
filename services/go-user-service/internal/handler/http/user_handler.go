package http

import (
	"net/http"

	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/domain"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/infrastructure"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// GetMe retrieves the current authenticated user
// GET /api/users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetUser(userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}

// UpdateMe updates the current authenticated user
// PUT /api/users/me
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		infrastructure.LogError("invalid update request", err, nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.userService.UpdateUser(userID.(string), service.UpdateUserRequest{
		Email: req.Email,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}

// DeleteMe deletes the current authenticated user
// DELETE /api/users/me
func (h *UserHandler) DeleteMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.userService.DeleteUser(userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
