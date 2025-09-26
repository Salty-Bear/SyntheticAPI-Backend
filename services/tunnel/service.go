package tunnel

import (
	"context"

	"github.com/Aryaman/syntra/sdk"
)

// TunnelService defines the interface for tunnel management operations
// This service handles tunnel CRUD operations with validation and business logic
type TunnelService interface {
	// Create creates a new tunnel or returns existing tunnel if name already exists
	Create(ctx context.Context, tunnel *sdk.Tunnel) (*sdk.Tunnel, error)

	// Get retrieves a tunnel by ID
	Get(ctx context.Context, id string) (*sdk.Tunnel, error)

	// List retrieves all tunnels with optional filtering and pagination
	List(ctx context.Context, query *sdk.TunnelQuery) ([]*sdk.Tunnel, error)

	// Update updates an existing tunnel
	Update(ctx context.Context, tunnel *sdk.Tunnel) (*sdk.Tunnel, error)

	// Delete deletes a tunnel by ID
	Delete(ctx context.Context, id string) error
}
