package generate

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers generate task management routes
func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")

	// Create generate task
	v1.Post("/", Create)

	// Get All generate tasks
	v1.Get("/", List)

	// Get generate task by ID
	v1.Get("/:id", Get)

	// Update generate task by ID
	v1.Put("/:id", Update)

	// Delete generate task by ID
	v1.Delete("/:id", Delete)

	// Execute generate task by ID
	v1.Post("/:id/execute", Execute)
}
