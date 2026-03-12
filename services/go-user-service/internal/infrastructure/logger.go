package infrastructure

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

// InitLogger initializes the structured logger
func InitLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// LogAuthAttempt logs authentication attempts
func LogAuthAttempt(email string, success bool, reason string) {
	if success {
		Logger.Info("authentication attempt succeeded",
			"email", email,
			"success", true,
		)
	} else {
		Logger.Warn("authentication attempt failed",
			"email", email,
			"success", false,
			"reason", reason,
		)
	}
}

// LogUserCreated logs user creation
func LogUserCreated(userID, email string) {
	Logger.Info("user created",
		"user_id", userID,
		"email", email,
	)
}

// LogUserUpdated logs user updates
func LogUserUpdated(userID string) {
	Logger.Info("user updated",
		"user_id", userID,
	)
}

// LogUserDeleted logs user deletion
func LogUserDeleted(userID string) {
	Logger.Info("user deleted",
		"user_id", userID,
	)
}

// LogError logs errors with context
func LogError(message string, err error, context map[string]interface{}) {
	args := []interface{}{"error", err.Error()}
	for k, v := range context {
		args = append(args, k, v)
	}
	Logger.Error(message, args...)
}
