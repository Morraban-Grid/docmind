package domain

import (
	"net/http"
)

// AppError represents a standardized application error
type AppError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// Common error constructors
func NewValidationError(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

func NewAuthenticationError(message string) *AppError {
	return &AppError{
		Code:       "INVALID_CREDENTIALS",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewAuthorizationError(message string) *AppError {
	return &AppError{
		Code:       "ACCESS_DENIED",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}
