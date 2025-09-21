package user

import (
	"context"
	"github.com/Aryaman/pub-sub/sdk"
)

// UserService defines the interface for user management operations
// This service handles user CRUD operations with Firebase authentication integration
type UserService interface {
	// Create creates a new user
	Create(ctx context.Context, user *sdk.User) error

	// Get retrieves a user by ID
	Get(ctx context.Context, id string) (*sdk.User, error)

	// List retrieves all users with optional filtering and pagination
	List(ctx context.Context, query *sdk.UserQuery) ([]*sdk.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *sdk.User) (*sdk.User, error)

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error

	// VerifyFirebaseToken validates Firebase authentication token
	VerifyFirebaseToken(ctx context.Context, token string) (*sdk.FirebaseUser, error)
}