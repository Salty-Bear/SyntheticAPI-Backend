package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

type UserContextKey struct {
	key string
}

var userContextKey = UserContextKey{"user_id"}

// RequireAuth middleware extracts user ID from X-User-ID header
// and adds it to the context
func RequireAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get user ID from header
		userID := c.Get("X-User-ID")

		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "User ID header required",
			})
		}

		// Add user ID to context
		c.Locals(userContextKey, userID)

		return c.Next()
	}
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(c *fiber.Ctx) string {
	userID, ok := c.Locals(userContextKey).(string)
	if !ok {
		return ""
	}
	return userID
}
