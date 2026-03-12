package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/docmind/go-user-service/internal/domain"
	"github.com/docmind/go-user-service/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

// ErrorHandler handles errors and returns standardized JSON responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Check if it's an AppError
			if appErr, ok := err.(*domain.AppError); ok {
				c.JSON(appErr.StatusCode, gin.H{"error": appErr})
				return
			}

			// Default to internal server error
			infrastructure.LogError("unhandled error", err, map[string]interface{}{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			})

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": domain.NewInternalError("An unexpected error occurred"),
			})
		}
	}
}

// RecoveryMiddleware recovers from panics and logs them
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				infrastructure.Logger.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": domain.NewInternalError("An unexpected error occurred"),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
