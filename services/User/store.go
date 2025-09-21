package user

import (
	"context"

	"github.com/Aryaman/pub-sub/db/models"
)

// UserStore defines the interface for user data access operations
// This interface abstracts the database operations for user management
type UserStore interface {
	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user *models.User) error

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, userID string) (*models.User, error)

	// GetUserByEmail retrieves a user by their email address
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// GetUsers retrieves all users with optional filtering
	GetUsers(ctx context.Context, projectID string, enabled *bool) ([]*models.User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error

	// DeleteUser deletes a user by their ID
	DeleteUser(ctx context.Context, userID string) error

	// UserExists checks if a user exists by ID
	UserExists(ctx context.Context, userID string) (bool, error)

	// EmailExists checks if a user exists by email
	EmailExists(ctx context.Context, email string) (bool, error)
}