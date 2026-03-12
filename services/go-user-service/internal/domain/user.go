package domain

import (
	"time"
)

// User represents a system user
type User struct {
	UserID       string    `json:"user_id" db:"user_id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(userID string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(userID string) error
}
