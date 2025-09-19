package connector

import (
	"context"
	"sync"
	"time"
)

// TunnelStatus represents the status of a tunnel
type TunnelStatus string

const (
	TunnelStatusActive     TunnelStatus = "active"
	TunnelStatusInactive   TunnelStatus = "inactive"
	TunnelStatusConnecting TunnelStatus = "connecting"
	TunnelStatusError      TunnelStatus = "error"
	TunnelStatusTerminated TunnelStatus = "terminated"
)

// Tunnel represents a tunnel connection
type Tunnel struct {
	ID           string       `json:"id"`
	UserID       string       `json:"user_id"`
	LocalPort    int          `json:"local_port"`
	PublicURL    string       `json:"public_url"`
	Status       TunnelStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	LastActivity time.Time    `json:"last_activity"`
	ExpiresAt    time.Time    `json:"expires_at"`
	Analytics    *Analytics   `json:"analytics,omitempty"`
}

// Analytics represents tunnel usage analytics
type Analytics struct {
	RequestCount     int64           `json:"request_count"`
	ResponseTimes    []time.Duration `json:"response_times"`
	StatusCodes      map[int]int     `json:"status_codes"`
	Endpoints        map[string]int  `json:"endpoints"`
	LastRequest      time.Time       `json:"last_request"`
	BytesTransferred int64           `json:"bytes_transferred"`
}

// Store defines the interface for tunnel storage operations
type Store interface {
	// CreateTunnel creates a new tunnel record
	CreateTunnel(ctx context.Context, tunnel *Tunnel) error

	// GetTunnel retrieves a tunnel by ID
	GetTunnel(ctx context.Context, tunnelID string) (*Tunnel, error)

	// GetTunnelByUserID retrieves active tunnels for a user
	GetTunnelsByUserID(ctx context.Context, userID string) ([]*Tunnel, error)

	// UpdateTunnelStatus updates the status of a tunnel
	UpdateTunnelStatus(ctx context.Context, tunnelID string, status TunnelStatus) error

	// UpdateTunnelActivity updates the last activity time
	UpdateTunnelActivity(ctx context.Context, tunnelID string) error

	// UpdateTunnelAnalytics updates tunnel analytics
	UpdateTunnelAnalytics(ctx context.Context, tunnelID string, analytics *Analytics) error

	// DeleteTunnel removes a tunnel record
	DeleteTunnel(ctx context.Context, tunnelID string) error

	// GetExpiredTunnels returns tunnels that have expired
	GetExpiredTunnels(ctx context.Context) ([]*Tunnel, error)

	// CleanupExpiredTunnels removes expired tunnel records
	CleanupExpiredTunnels(ctx context.Context) error
}

// MemoryStore implements Store interface using in-memory storage
type MemoryStore struct {
	tunnels map[string]*Tunnel
	mutex   sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() Store {
	return &MemoryStore{
		tunnels: make(map[string]*Tunnel),
	}
}
