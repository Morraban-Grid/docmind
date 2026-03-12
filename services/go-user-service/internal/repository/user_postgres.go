package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/docmind/go-user-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *PostgresUserRepository) CreateUser(user *domain.User) error {
	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING user_id, created_at, updated_at
	`

	err = r.db.QueryRow(query, user.Email, string(hashedPassword), user.Name).
		Scan(&user.UserID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (r *PostgresUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT user_id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *PostgresUserRepository) GetUserByID(userID string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT user_id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE user_id = $1
	`

	err := r.db.QueryRow(query, userID).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (r *PostgresUserRepository) UpdateUser(user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, updated_at = NOW()
		WHERE user_id = $2
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, user.Name, user.UserID).Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user from the database
func (r *PostgresUserRepository) DeleteUser(userID string) error {
	query := `DELETE FROM users WHERE user_id = $1`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// VerifyPassword checks if the provided password matches the stored hash
func (r *PostgresUserRepository) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
