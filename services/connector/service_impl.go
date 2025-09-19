package connector

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	store  Store
	config TunnelConfig
	client *http.Client
}

// NewService creates a new tunnel service instance
func NewService(store Store, config TunnelConfig) Service {
	return &ServiceImpl{
		store:  store,
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateTunnel creates a new tunnel for exposing localhost
func (s *ServiceImpl) CreateTunnel(ctx context.Context, userID string, localPort int) (*Tunnel, error) {
	// Check if port is allowed
	if !s.isPortAllowed(localPort) {
		return nil, fmt.Errorf("port %d is not allowed", localPort)
	}

	// Check user tunnel limit
	existingTunnels, err := s.store.GetTunnelsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing tunnels: %w", err)
	}

	if len(existingTunnels) >= s.config.MaxTunnelsPerUser {
		return nil, fmt.Errorf("maximum number of tunnels (%d) exceeded for user", s.config.MaxTunnelsPerUser)
	}

	// Generate tunnel ID and public URL
	tunnelID := uuid.New().String()
	publicURL := s.generatePublicURL(tunnelID)

	tunnel := &Tunnel{
		ID:           tunnelID,
		UserID:       userID,
		LocalPort:    localPort,
		PublicURL:    publicURL,
		Status:       TunnelStatusConnecting,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(s.config.TunnelExpiration),
		Analytics: &Analytics{
			StatusCodes: make(map[int]int),
			Endpoints:   make(map[string]int),
		},
	}

	err = s.store.CreateTunnel(ctx, tunnel)
	if err != nil {
		return nil, fmt.Errorf("failed to create tunnel: %w", err)
	}

	// Start tunnel connection (this would be where you establish the actual tunnel)
	go s.establishTunnel(ctx, tunnel)

	return tunnel, nil
}

// GetTunnel retrieves tunnel information by ID
func (s *ServiceImpl) GetTunnel(ctx context.Context, tunnelID string) (*Tunnel, error) {
	return s.store.GetTunnel(ctx, tunnelID)
}

// GetUserTunnels retrieves all active tunnels for a user
func (s *ServiceImpl) GetUserTunnels(ctx context.Context, userID string) ([]*Tunnel, error) {
	return s.store.GetTunnelsByUserID(ctx, userID)
}

// TerminateTunnel terminates an active tunnel
func (s *ServiceImpl) TerminateTunnel(ctx context.Context, tunnelID string) error {
	tunnel, err := s.store.GetTunnel(ctx, tunnelID)
	if err != nil {
		return fmt.Errorf("tunnel not found: %w", err)
	}

	if tunnel.Status == TunnelStatusTerminated {
		return fmt.Errorf("tunnel is already terminated")
	}

	err = s.store.UpdateTunnelStatus(ctx, tunnelID, TunnelStatusTerminated)
	if err != nil {
		return fmt.Errorf("failed to update tunnel status: %w", err)
	}

	// Here you would close the actual tunnel connection
	return nil
}

// HandleTunnelRequest forwards requests through the tunnel
func (s *ServiceImpl) HandleTunnelRequest(ctx context.Context, tunnelID string, requestData []byte) ([]byte, error) {
	tunnel, err := s.store.GetTunnel(ctx, tunnelID)
	if err != nil {
		return nil, fmt.Errorf("tunnel not found: %w", err)
	}

	if tunnel.Status != TunnelStatusActive {
		return nil, fmt.Errorf("tunnel is not active")
	}

	// Update last activity
	s.store.UpdateTunnelActivity(ctx, tunnelID)

	// Forward request to localhost (this is a simplified implementation)
	localURL := fmt.Sprintf("http://localhost:%d", tunnel.LocalPort)

	req, err := http.NewRequestWithContext(ctx, "POST", localURL, bytes.NewReader(requestData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return responseData, nil
}

// UpdateTunnelAnalytics updates analytics for a tunnel
func (s *ServiceImpl) UpdateTunnelAnalytics(ctx context.Context, tunnelID string, statusCode int, responseTime time.Duration, endpoint string, bytesTransferred int64) error {
	tunnel, err := s.store.GetTunnel(ctx, tunnelID)
	if err != nil {
		return err
	}

	if tunnel.Analytics == nil {
		tunnel.Analytics = &Analytics{
			StatusCodes: make(map[int]int),
			Endpoints:   make(map[string]int),
		}
	}

	tunnel.Analytics.RequestCount++
	tunnel.Analytics.ResponseTimes = append(tunnel.Analytics.ResponseTimes, responseTime)
	tunnel.Analytics.StatusCodes[statusCode]++
	tunnel.Analytics.Endpoints[endpoint]++
	tunnel.Analytics.LastRequest = time.Now()
	tunnel.Analytics.BytesTransferred += bytesTransferred

	return s.store.UpdateTunnelAnalytics(ctx, tunnelID, tunnel.Analytics)
}

// CleanupExpiredTunnels removes expired tunnels
func (s *ServiceImpl) CleanupExpiredTunnels(ctx context.Context) error {
	return s.store.CleanupExpiredTunnels(ctx)
}

// Helper methods

func (s *ServiceImpl) isPortAllowed(port int) bool {
	if len(s.config.AllowedPorts) == 0 {
		return true // Allow all ports if none specified
	}

	for _, allowedPort := range s.config.AllowedPorts {
		if port == allowedPort {
			return true
		}
	}
	return false
}

func (s *ServiceImpl) generatePublicURL(tunnelID string) string {
	// Generate a random subdomain for the tunnel
	subdomain := s.generateRandomString(8)
	return fmt.Sprintf("%s/%s/%s", s.config.BaseURL, subdomain, tunnelID)
}

func (s *ServiceImpl) generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *ServiceImpl) establishTunnel(ctx context.Context, tunnel *Tunnel) {
	// This is where you would establish the actual tunnel connection
	// For now, we'll just simulate the connection being established
	time.Sleep(2 * time.Second)

	err := s.store.UpdateTunnelStatus(ctx, tunnel.ID, TunnelStatusActive)
	if err != nil {
		s.store.UpdateTunnelStatus(ctx, tunnel.ID, TunnelStatusError)
	}
}
