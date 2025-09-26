package tunnel

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers tunnel management routes
func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")

	// Create tunnel
	v1.Post("/", Create)

	// Get All tunnels
	v1.Get("/", List)

	// Get tunnel by ID
	v1.Get("/:id", Get)

	// Update tunnel by ID
	v1.Put("/:id", Update)

	// Delete tunnel by ID
	v1.Delete("/:id", Delete)
}
