package tunnel

import (
	"context"

	"github.com/Aryaman/syntra/db/models"
)

// TunnelStore defines the interface for tunnel data access operations
// This interface abstracts the database operations for tunnel management
type TunnelStore interface {
	// CreateTunnel creates a new tunnel in the database
	CreateTunnel(ctx context.Context, tunnel *models.Tunnel) error

	// GetTunnelByID retrieves a tunnel by its ID
	GetTunnelByID(ctx context.Context, tunnelID string) (*models.Tunnel, error)

	// GetTunnels retrieves all tunnels with optional filtering
	GetTunnels(ctx context.Context, projectID, status, protocol string, enabled *bool) ([]*models.Tunnel, error)

	// UpdateTunnel updates an existing tunnel
	UpdateTunnel(ctx context.Context, tunnelID string, updates map[string]interface{}) error

	// DeleteTunnel deletes a tunnel by its ID
	DeleteTunnel(ctx context.Context, tunnelID string) error

	// TunnelExists checks if a tunnel exists by ID
	TunnelExists(ctx context.Context, tunnelID string) (bool, error)

	// NameExists checks if a tunnel exists by name within a project
	NameExists(ctx context.Context, projectID, name string) (bool, error)
}
