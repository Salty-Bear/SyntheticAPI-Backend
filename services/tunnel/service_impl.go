package tunnel

import (
	"context"
	"fmt"

	"github.com/Aryaman/syntra/sdk"
	"github.com/google/uuid"
)

// ServiceImpl implements the TunnelService interface
type ServiceImpl struct {
	store TunnelStore
}

// NewService creates a new instance of TunnelService
func NewService(store TunnelStore) TunnelService {
	return &ServiceImpl{
		store: store,
	}
}

// Create creates a new tunnel or returns existing tunnel if name already exists in project
func (s *ServiceImpl) Create(ctx context.Context, tunnel *sdk.Tunnel) (*sdk.Tunnel, error) {
	// Generate ID if not provided
	if tunnel.ID == "" {
		tunnel.ID = uuid.New().String()
	}

	// Check if tunnel with name already exists in the project
	exists, err := s.store.NameExists(ctx, tunnel.ProjectID, tunnel.Name)
	if err != nil {
		return nil, fmt.Errorf("error checking tunnel name existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("tunnel with name '%s' already exists in project '%s'", tunnel.Name, tunnel.ProjectID)
	}

	// Set default values
	if tunnel.Status == "" {
		tunnel.Status = "inactive"
	}
	if tunnel.Protocol == "" {
		tunnel.Protocol = "http"
	}

	// Create database model from SDK tunnel
	dbTunnel := fromSdkToModel(*tunnel)

	// Save tunnel to database
	if err := s.store.CreateTunnel(ctx, &dbTunnel); err != nil {
		return nil, fmt.Errorf("error creating tunnel: %w", err)
	}

	// Convert back to SDK type to return
	return fromModelToSdk(&dbTunnel), nil
}

// Get retrieves a tunnel by ID
func (s *ServiceImpl) Get(ctx context.Context, id string) (*sdk.Tunnel, error) {
	if id == "" {
		return nil, fmt.Errorf("tunnel ID is required")
	}

	tunnel, err := s.store.GetTunnelByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tunnel: %w", err)
	}

	if tunnel == nil {
		return nil, fmt.Errorf("tunnel not found")
	}

	// Convert database model to SDK type
	return fromModelToSdk(tunnel), nil
}

// List retrieves all tunnels with optional filtering and pagination
func (s *ServiceImpl) List(ctx context.Context, query *sdk.TunnelQuery) ([]*sdk.Tunnel, error) {
	var (
		projectID string
		status    string
		protocol  string
		enabled   *bool
	)

	if query != nil {
		projectID = query.ProjectID
		status = query.Status
		protocol = query.Protocol
		enabled = query.Enabled
	}

	tunnels, err := s.store.GetTunnels(ctx, projectID, status, protocol, enabled)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tunnels: %w", err)
	}

	// Apply pagination if specified
	if query != nil && query.Limit > 0 {
		start := query.Offset
		end := start + query.Limit

		if start >= len(tunnels) {
			return []*sdk.Tunnel{}, nil
		}

		if end > len(tunnels) {
			end = len(tunnels)
		}

		tunnels = tunnels[start:end]
	}

	// Convert database models to SDK types
	var sdkTunnels []*sdk.Tunnel
	for _, tunnel := range tunnels {
		sdkTunnels = append(sdkTunnels, fromModelToSdk(tunnel))
	}

	return sdkTunnels, nil
}

// Update updates an existing tunnel
func (s *ServiceImpl) Update(ctx context.Context, tunnel *sdk.Tunnel) (*sdk.Tunnel, error) {
	if tunnel.ID == "" {
		return nil, fmt.Errorf("tunnel ID is required")
	}

	// Check if tunnel exists
	exists, err := s.store.TunnelExists(ctx, tunnel.ID)
	if err != nil {
		return nil, fmt.Errorf("error checking tunnel existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("tunnel not found")
	}

	// Build updates map
	updates := make(map[string]interface{})
	if tunnel.Name != "" {
		updates["name"] = tunnel.Name
	}
	if tunnel.Description != "" {
		updates["description"] = tunnel.Description
	}
	if tunnel.Endpoint != "" {
		updates["endpoint"] = tunnel.Endpoint
	}
	if tunnel.Port > 0 {
		updates["port"] = tunnel.Port
	}
	if tunnel.Protocol != "" {
		updates["protocol"] = tunnel.Protocol
	}
	if tunnel.Status != "" {
		updates["status"] = tunnel.Status
	}
	updates["enabled"] = tunnel.Enabled

	// Update tunnel
	if err := s.store.UpdateTunnel(ctx, tunnel.ID, updates); err != nil {
		return nil, fmt.Errorf("error updating tunnel: %w", err)
	}

	// Retrieve and return updated tunnel
	updatedTunnel, err := s.store.GetTunnelByID(ctx, tunnel.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving updated tunnel: %w", err)
	}

	// Convert to SDK type
	return fromModelToSdk(updatedTunnel), nil
}

// Delete deletes a tunnel by ID
func (s *ServiceImpl) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("tunnel ID is required")
	}

	// Check if tunnel exists
	exists, err := s.store.TunnelExists(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking tunnel existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("tunnel not found")
	}

	// Delete tunnel
	if err := s.store.DeleteTunnel(ctx, id); err != nil {
		return fmt.Errorf("error deleting tunnel: %w", err)
	}

	return nil
}

// GetByProjectID retrieves all tunnels for a specific project
func (s *ServiceImpl) GetByProjectID(ctx context.Context, projectID string) ([]*sdk.Tunnel, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	query := &sdk.TunnelQuery{
		ProjectID: projectID,
	}

	return s.List(ctx, query)
}
