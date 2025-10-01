package generate

import (
	"context"

	"github.com/Aryaman/syntra/db/models"
)

// GenerateStore defines the interface for generate task data access operations
// This interface abstracts the database operations for generate task management
type GenerateStore interface {
	// CreateGenerate creates a new generate task in the database
	CreateGenerate(ctx context.Context, generate *models.Generate) error

	// GetGenerateByID retrieves a generate task by its ID
	GetGenerateByID(ctx context.Context, generateID string) (*models.Generate, error)

	// GetGenerates retrieves all generate tasks with optional filtering
	GetGenerates(ctx context.Context, status, dataType string, enabled *bool, userId string) ([]*models.Generate, error)

	// UpdateGenerate updates an existing generate task
	UpdateGenerate(ctx context.Context, generateID string, updates map[string]interface{}) error

	// DeleteGenerate deletes a generate task by its ID
	DeleteGenerate(ctx context.Context, generateID string) error

	// GenerateExists checks if a generate task exists by ID
	GenerateExists(ctx context.Context, generateID string) (bool, error)

	// NameExists checks if a generate task exists by name for a specific user
	NameExists(ctx context.Context, name string, userId string) (bool, error)

	// GetGeneratesByUser retrieves all generate tasks for a specific user
	GetGeneratesByUser(ctx context.Context, userId string, status, dataType string, enabled *bool) ([]*models.Generate, error)
}
