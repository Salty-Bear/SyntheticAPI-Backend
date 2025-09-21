package user

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers user management routes
func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")

	// Create user
	v1.Post("/", Create)

	// Get All users
	v1.Get("/", List)

	// Get User by ID
	v1.Get("/:id", Get)

	// Update User by ID
	v1.Put("/:id", Update)

	// Delete User by ID
	v1.Delete("/:id", Delete)
}
