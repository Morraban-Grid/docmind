package middleware

import (
	"net/http"
	"strings"

	"github.com/docmind/go-user-service/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtManager *infrastructure.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "TOKEN_MISSING",
					"message": "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// Check for Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "TOKEN_INVALID",
					"message": "Authorization header must be in format: Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token and extract claims
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "TOKEN_INVALID",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// Attach user_id to request context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
