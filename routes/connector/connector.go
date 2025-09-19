package connector

import (
	"github.com/Aryaman/pub-sub/services/connector"
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for tunnel operations
type Handler struct {
	service connector.Service
	proxy   *connector.TunnelProxy
}

// Service interface for dependency injection
type Service = connector.Service

// NewHandler creates a new handler instance
func NewHandler(service connector.Service) *Handler {
	return &Handler{
		service: service,
		proxy:   connector.NewTunnelProxy(service),
	}
}

// RegisterRoutes registers all connector routes
func RegisterRoutes(app fiber.Router, handler *Handler, wsHandler *WebSocketHandler) {
	// REST API routes
	app.Post("/create", handler.CreateTunnel)
	app.Get("/:tunnelId", handler.GetTunnel)
	app.Get("/user/:userId", handler.GetUserTunnels)
	app.Delete("/:tunnelId", handler.TerminateTunnel)

	// WebSocket and proxy routes
	RegisterWebSocketRoutes(app, wsHandler)
}

// CreateTunnelRequest represents the request payload for creating a tunnel
type CreateTunnelRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	LocalPort int    `json:"local_port" validate:"required,min=1,max=65535"`
}

// CreateTunnel creates a new tunnel
func (h *Handler) CreateTunnel(c *fiber.Ctx) error {
	var req CreateTunnelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate port
	if err := connector.ValidatePort(req.LocalPort); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid port",
			"details": err.Error(),
		})
	}

	tunnel, err := h.service.CreateTunnel(c.Context(), req.UserID, req.LocalPort)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create tunnel",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    tunnel,
		"message": "Tunnel created successfully",
	})
}

// GetTunnel retrieves tunnel information
func (h *Handler) GetTunnel(c *fiber.Ctx) error {
	tunnelID := c.Params("tunnelId")
	if tunnelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tunnel ID is required",
		})
	}

	tunnel, err := h.service.GetTunnel(c.Context(), tunnelID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Tunnel not found",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    tunnel,
	})
}

// GetUserTunnels retrieves all tunnels for a user
func (h *Handler) GetUserTunnels(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	tunnels, err := h.service.GetUserTunnels(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve tunnels",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    tunnels,
		"count":   len(tunnels),
	})
}

// TerminateTunnel terminates an active tunnel
func (h *Handler) TerminateTunnel(c *fiber.Ctx) error {
	tunnelID := c.Params("tunnelId")
	if tunnelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tunnel ID is required",
		})
	}

	err := h.service.TerminateTunnel(c.Context(), tunnelID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to terminate tunnel",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Tunnel terminated successfully",
	})
}

// ProxyRequest forwards requests through the tunnel
func (h *Handler) ProxyRequest(c *fiber.Ctx) error {
	return h.proxy.ProxyHandler(c)
}
