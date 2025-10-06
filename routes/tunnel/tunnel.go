package tunnel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Aryaman/syntra/middlewares"
	"github.com/Aryaman/syntra/providers"
	"github.com/Aryaman/syntra/sdk"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

// Create creates a new tunnel
func Create(c *fiber.Ctx) error {
	log.Debug("received create tunnel request")

	// Get user ID from context
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.TunnelResponse{
			Success: false,
			Message: "user authentication required",
		})
	}

	payload := new(sdk.Tunnel)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.TunnelResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}
	log.Debug("parsed create tunnel request")

	// Generate ID if not provided
	if payload.ID == "" {
		payload.ID = uuid.New().String()
	}

	// Set user ID from context
	payload.CreatedBy = userID
	payload.UpdatedBy = userID

	pr := providers.GetProviders(c)
	createdTunnel, err := pr.S.Tunnel.Create(c.Context(), payload)
	if err != nil {
		message := fmt.Sprintf("failed to create tunnel. %v", err)
		log.Errorw("failed to create tunnel", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.TunnelResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("tunnel created successfully")

	return c.Status(http.StatusCreated).JSON(sdk.TunnelResponse{
		Success: true,
		Message: "Tunnel created successfully",
		Data:    createdTunnel,
	})
}

// Get retrieves a tunnel by ID
func Get(c *fiber.Ctx) error {
	log.Debug("received get tunnel request")
	tunnelID := c.Params("id")

	pr := providers.GetProviders(c)
	tunnel, err := pr.S.Tunnel.Get(c.Context(), tunnelID)
	if err != nil {
		message := fmt.Sprintf("failed to get tunnel. %v", err)
		log.Errorw("failed to get tunnel", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.TunnelResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("tunnel retrieved successfully")
	return c.JSON(sdk.TunnelResponse{
		Success: true,
		Message: "Tunnel retrieved successfully",
		Data:    tunnel,
	})
}

// List retrieves all tunnels with optional filtering and pagination
func List(c *fiber.Ctx) error {
	log.Debug("received list tunnels request")

	// Parse query parameters
	query := &sdk.TunnelQuery{
		Status:   c.Query("status"),
		Protocol: c.Query("protocol"),
	}

	// Parse enabled parameter
	enabledStr := c.Query("enabled")
	if enabledStr != "" {
		enabledBool, err := strconv.ParseBool(enabledStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(sdk.TunnelResponse{
				Success: false,
				Message: fmt.Sprintf("invalid enabled parameter. %v", err),
			})
		}
		query.Enabled = &enabledBool
	}

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(sdk.TunnelResponse{
				Success: false,
				Message: fmt.Sprintf("invalid page parameter. %v", err),
			})
		}
		query.Page = page
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(sdk.TunnelResponse{
				Success: false,
				Message: fmt.Sprintf("invalid limit parameter. %v", err),
			})
		}
		query.Limit = limit
		// Calculate offset based on page and limit
		query.Offset = (query.Page - 1) * query.Limit
		if query.Offset < 0 {
			query.Offset = 0
		}
	}

	pr := providers.GetProviders(c)
	tunnels, err := pr.S.Tunnel.List(c.Context(), query)
	if err != nil {
		message := fmt.Sprintf("failed to list tunnels. %v", err)
		log.Errorw("failed to list tunnels", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.TunnelResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("tunnels retrieved successfully")
	return c.JSON(sdk.TunnelResponse{
		Success: true,
		Message: "Tunnels retrieved successfully",
		Tunnels: tunnels,
	})
}

// Update updates an existing tunnel
func Update(c *fiber.Ctx) error {
	log.Debug("received update tunnel request")
	tunnelID := c.Params("id")

	payload := new(sdk.Tunnel)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.TunnelResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}
	log.Debug("parsed update tunnel request")

	// Set the tunnel ID from the URL parameter
	payload.ID = tunnelID

	pr := providers.GetProviders(c)
	updatedTunnel, err := pr.S.Tunnel.Update(c.Context(), payload)
	if err != nil {
		message := fmt.Sprintf("failed to update tunnel. %v", err)
		log.Errorw("failed to update tunnel", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.TunnelResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("tunnel updated successfully")
	return c.JSON(sdk.TunnelResponse{
		Success: true,
		Message: "Tunnel updated successfully",
		Data:    updatedTunnel,
	})
}

// Delete deletes a tunnel by ID
func Delete(c *fiber.Ctx) error {
	log.Debug("received delete tunnel request")
	tunnelID := c.Params("id")

	pr := providers.GetProviders(c)
	err := pr.S.Tunnel.Delete(c.Context(), tunnelID)
	if err != nil {
		message := fmt.Sprintf("failed to delete tunnel. %v", err)
		log.Errorw("failed to delete tunnel", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.TunnelResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("tunnel deleted successfully")
	return c.JSON(sdk.TunnelResponse{
		Success: true,
		Message: "Tunnel deleted successfully",
	})
}
