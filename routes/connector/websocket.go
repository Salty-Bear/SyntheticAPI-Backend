package connector

import (
	"log"
	"strconv"

	"github.com/Aryaman/pub-sub/services/connector"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// WebSocketHandler handles tunnel WebSocket connections
type WebSocketHandler struct {
	tunnelManager *connector.TunnelManager
	reverseProxy  *connector.ReverseProxy
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(tunnelManager *connector.TunnelManager, reverseProxy *connector.ReverseProxy) *WebSocketHandler {
	return &WebSocketHandler{
		tunnelManager: tunnelManager,
		reverseProxy:  reverseProxy,
	}
}

// RegisterWebSocketRoutes registers WebSocket and proxy routes
func RegisterWebSocketRoutes(app fiber.Router, handler *WebSocketHandler) {
	// WebSocket endpoint for CLI connections
	app.Get("/ws/connect", websocket.New(handler.ConnectTunnel))

	// Proxy endpoints for forwarding requests
	app.All("/proxy/:subdomain/*", handler.ProxyRequest)

	// Health check endpoint
	app.Get("/health/:subdomain", handler.HealthCheck)

	// Stats endpoint
	app.Get("/stats/:subdomain", handler.GetStats)

	// List active tunnels
	app.Get("/active", handler.GetActiveTunnels)
}

// ConnectTunnel handles new tunnel WebSocket connections
func (h *WebSocketHandler) ConnectTunnel(c *websocket.Conn) {
	// Get user ID from query parameters or headers
	userID := c.Query("user_id")
	if userID == "" {
		log.Printf("WebSocket connection rejected: missing user_id")
		c.Close()
		return
	}

	// Get local port from query parameters
	localPortStr := c.Query("local_port", "3000")
	localPort, err := strconv.Atoi(localPortStr)
	if err != nil {
		log.Printf("WebSocket connection rejected: invalid local_port")
		c.Close()
		return
	}

	log.Printf("New WebSocket connection from user %s for port %d", userID, localPort)

	// Handle the WebSocket connection
	h.tunnelManager.HandleWebSocket(c, userID)
}

// ProxyRequest forwards requests through tunnels
func (h *WebSocketHandler) ProxyRequest(c *fiber.Ctx) error {
	return h.reverseProxy.ProxyHandler(c)
}

// HealthCheck checks tunnel health
func (h *WebSocketHandler) HealthCheck(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain")
	if subdomain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "subdomain is required",
		})
	}

	healthy, err := h.reverseProxy.HealthCheck(subdomain)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   err.Error(),
			"healthy": false,
		})
	}

	return c.JSON(fiber.Map{
		"healthy":   healthy,
		"subdomain": subdomain,
	})
}

// GetStats returns tunnel statistics
func (h *WebSocketHandler) GetStats(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain")
	if subdomain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "subdomain is required",
		})
	}

	stats, err := h.reverseProxy.GetTunnelStats(subdomain)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// GetActiveTunnels returns all active tunnel connections
func (h *WebSocketHandler) GetActiveTunnels(c *fiber.Ctx) error {
	connections := h.tunnelManager.GetActiveConnections()

	tunnels := make([]map[string]interface{}, 0, len(connections))
	for _, conn := range connections {
		tunnels = append(tunnels, map[string]interface{}{
			"id":         conn.ID,
			"user_id":    conn.UserID,
			"subdomain":  conn.Subdomain,
			"local_port": conn.LocalPort,
			"created_at": conn.CreatedAt,
			"last_ping":  conn.LastPing,
			"public_url": "https://" + conn.Subdomain + ".syntra.dev",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    tunnels,
		"count":   len(tunnels),
	})
}
