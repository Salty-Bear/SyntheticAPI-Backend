package connector

import (
	"context"
	"time"
)

// Service defines the interface for tunnel service operations
type Service interface {
	// CreateTunnel creates a new tunnel for exposing localhost
	CreateTunnel(ctx context.Context, userID string, localPort int) (*Tunnel, error)

	// GetTunnel retrieves tunnel information by ID
	GetTunnel(ctx context.Context, tunnelID string) (*Tunnel, error)

	// GetUserTunnels retrieves all active tunnels for a user
	GetUserTunnels(ctx context.Context, userID string) ([]*Tunnel, error)

	// TerminateTunnel terminates an active tunnel
	TerminateTunnel(ctx context.Context, tunnelID string) error

	// HandleTunnelRequest forwards requests through the tunnel
	HandleTunnelRequest(ctx context.Context, tunnelID string, requestData []byte) ([]byte, error)

	// UpdateTunnelAnalytics updates analytics for a tunnel
	UpdateTunnelAnalytics(ctx context.Context, tunnelID string, statusCode int, responseTime time.Duration, endpoint string, bytesTransferred int64) error

	// CleanupExpiredTunnels removes expired tunnels
	CleanupExpiredTunnels(ctx context.Context) error
}

// TunnelConfig represents configuration for tunnel creation
type TunnelConfig struct {
	MaxTunnelsPerUser int           `json:"max_tunnels_per_user"`
	TunnelExpiration  time.Duration `json:"tunnel_expiration"`
	AllowedPorts      []int         `json:"allowed_ports"`
	BaseURL           string        `json:"base_url"`
}
