package generate

import (
	"context"

	"github.com/Aryaman/syntra/sdk"
)

// GenerateService defines the interface for generate task management operations
// This service handles generate CRUD operations with validation and business logic
type GenerateService interface {
	// Create creates a new generate task or returns existing generate task if name already exists
	Create(ctx context.Context, generate *sdk.Generate) (*sdk.Generate, error)

	// Get retrieves a generate task by ID
	Get(ctx context.Context, id string, userId string) (*sdk.Generate, error)

	// List retrieves all generate tasks with optional filtering and pagination
	List(ctx context.Context, query *sdk.GenerateQuery) ([]*sdk.Generate, error)

	// Update updates an existing generate task
	Update(ctx context.Context, generate *sdk.Generate) (*sdk.Generate, error)

	// Delete deletes a generate task by ID
	Delete(ctx context.Context, id string, userId string) error

	// Execute executes a generate task to produce synthetic data
	Execute(ctx context.Context, id string, userId string) (*sdk.Generate, error)
}
