package connector

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ReverseProxy handles incoming requests and forwards them through tunnels
type ReverseProxy struct {
	tunnelManager *TunnelManager
	baseDomain    string
}

// NewReverseProxy creates a new reverse proxy instance
func NewReverseProxy(tunnelManager *TunnelManager, baseDomain string) *ReverseProxy {
	return &ReverseProxy{
		tunnelManager: tunnelManager,
		baseDomain:    baseDomain,
	}
}

// ProxyHandler handles incoming requests to subdomains
func (rp *ReverseProxy) ProxyHandler(c *fiber.Ctx) error {
	// Extract subdomain from Host header
	host := c.Get("Host")
	subdomain := rp.extractSubdomain(host)

	if subdomain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subdomain",
		})
	}

	// Create HTTP request from Fiber context
	req, err := rp.fiberToHTTPRequest(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process request",
		})
	}

	// Forward request through tunnel
	startTime := time.Now()
	resp, err := rp.tunnelManager.ForwardRequest(subdomain, req)
	responseTime := time.Since(startTime)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Tunnel unavailable",
			"details": err.Error(),
		})
	}

	// Copy response back to Fiber
	err = rp.httpResponseToFiber(c, resp, responseTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process response",
		})
	}

	return nil
}

// extractSubdomain extracts subdomain from host header
func (rp *ReverseProxy) extractSubdomain(host string) string {
	// Remove port if present
	if colonIndex := strings.LastIndex(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	// Check if it matches our base domain pattern
	suffix := "." + rp.baseDomain
	if !strings.HasSuffix(host, suffix) {
		return ""
	}

	// Extract subdomain
	subdomain := strings.TrimSuffix(host, suffix)
	if subdomain == "" || subdomain == rp.baseDomain {
		return ""
	}

	return subdomain
}

// fiberToHTTPRequest converts Fiber request to standard HTTP request
func (rp *ReverseProxy) fiberToHTTPRequest(c *fiber.Ctx) (*http.Request, error) {
	// Build URL
	scheme := "http"
	if c.Protocol() == "https" {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s%s", scheme, c.Get("Host"), c.OriginalURL())

	// Create request
	req, err := http.NewRequest(c.Method(), url, bytes.NewReader(c.Body()))
	if err != nil {
		return nil, err
	}

	// Copy headers
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})

	return req, nil
}

// httpResponseToFiber copies HTTP response to Fiber context
func (rp *ReverseProxy) httpResponseToFiber(c *fiber.Ctx, resp *http.Response, responseTime time.Duration) error {
	// Set status code
	c.Status(resp.StatusCode)

	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	// Add timing header
	c.Set("X-Tunnel-Response-Time", fmt.Sprintf("%dms", responseTime.Milliseconds()))

	// Copy response body if present
	if resp.Body != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return c.Send(body)
	}

	return nil
}

// HealthCheck checks if a tunnel is healthy
func (rp *ReverseProxy) HealthCheck(subdomain string) (bool, error) {
	connections := rp.tunnelManager.GetActiveConnections()

	for _, conn := range connections {
		if conn.Subdomain == subdomain {
			// Check if connection is recent and responsive
			if time.Since(conn.LastPing) < 2*time.Minute {
				return true, nil
			}
		}
	}

	return false, fmt.Errorf("tunnel not found or unhealthy")
}

// GetTunnelStats returns statistics for a tunnel
func (rp *ReverseProxy) GetTunnelStats(subdomain string) (map[string]interface{}, error) {
	connections := rp.tunnelManager.GetActiveConnections()

	for _, conn := range connections {
		if conn.Subdomain == subdomain {
			return map[string]interface{}{
				"tunnel_id":  conn.ID,
				"user_id":    conn.UserID,
				"subdomain":  conn.Subdomain,
				"local_port": conn.LocalPort,
				"created_at": conn.CreatedAt,
				"last_ping":  conn.LastPing,
				"uptime":     time.Since(conn.CreatedAt).Seconds(),
				"is_healthy": time.Since(conn.LastPing) < 2*time.Minute,
			}, nil
		}
	}

	return nil, fmt.Errorf("tunnel not found")
}
