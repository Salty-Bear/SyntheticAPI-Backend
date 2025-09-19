package connector

import (
	"context"
	"fmt"
	"time"
)

// CreateTunnel creates a new tunnel record
func (ms *MemoryStore) CreateTunnel(ctx context.Context, tunnel *Tunnel) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if _, exists := ms.tunnels[tunnel.ID]; exists {
		return fmt.Errorf("tunnel with ID %s already exists", tunnel.ID)
	}

	ms.tunnels[tunnel.ID] = tunnel
	return nil
}

// GetTunnel retrieves a tunnel by ID
func (ms *MemoryStore) GetTunnel(ctx context.Context, tunnelID string) (*Tunnel, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	tunnel, exists := ms.tunnels[tunnelID]
	if !exists {
		return nil, fmt.Errorf("tunnel with ID %s not found", tunnelID)
	}

	return tunnel, nil
}

// GetTunnelsByUserID retrieves active tunnels for a user
func (ms *MemoryStore) GetTunnelsByUserID(ctx context.Context, userID string) ([]*Tunnel, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	var userTunnels []*Tunnel
	for _, tunnel := range ms.tunnels {
		if tunnel.UserID == userID && tunnel.Status == TunnelStatusActive {
			userTunnels = append(userTunnels, tunnel)
		}
	}

	return userTunnels, nil
}

// UpdateTunnelStatus updates the status of a tunnel
func (ms *MemoryStore) UpdateTunnelStatus(ctx context.Context, tunnelID string, status TunnelStatus) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	tunnel, exists := ms.tunnels[tunnelID]
	if !exists {
		return fmt.Errorf("tunnel with ID %s not found", tunnelID)
	}

	tunnel.Status = status
	tunnel.LastActivity = time.Now()
	return nil
}

// UpdateTunnelActivity updates the last activity time
func (ms *MemoryStore) UpdateTunnelActivity(ctx context.Context, tunnelID string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	tunnel, exists := ms.tunnels[tunnelID]
	if !exists {
		return fmt.Errorf("tunnel with ID %s not found", tunnelID)
	}

	tunnel.LastActivity = time.Now()
	return nil
}

// UpdateTunnelAnalytics updates tunnel analytics
func (ms *MemoryStore) UpdateTunnelAnalytics(ctx context.Context, tunnelID string, analytics *Analytics) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	tunnel, exists := ms.tunnels[tunnelID]
	if !exists {
		return fmt.Errorf("tunnel with ID %s not found", tunnelID)
	}

	tunnel.Analytics = analytics
	tunnel.LastActivity = time.Now()
	return nil
}

// DeleteTunnel removes a tunnel record
func (ms *MemoryStore) DeleteTunnel(ctx context.Context, tunnelID string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if _, exists := ms.tunnels[tunnelID]; !exists {
		return fmt.Errorf("tunnel with ID %s not found", tunnelID)
	}

	delete(ms.tunnels, tunnelID)
	return nil
}

// GetExpiredTunnels returns tunnels that have expired
func (ms *MemoryStore) GetExpiredTunnels(ctx context.Context) ([]*Tunnel, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	var expiredTunnels []*Tunnel
	now := time.Now()

	for _, tunnel := range ms.tunnels {
		if now.After(tunnel.ExpiresAt) {
			expiredTunnels = append(expiredTunnels, tunnel)
		}
	}

	return expiredTunnels, nil
}

// CleanupExpiredTunnels removes expired tunnel records
func (ms *MemoryStore) CleanupExpiredTunnels(ctx context.Context) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	now := time.Now()
	for tunnelID, tunnel := range ms.tunnels {
		if now.After(tunnel.ExpiresAt) {
			delete(ms.tunnels, tunnelID)
		}
	}

	return nil
}
