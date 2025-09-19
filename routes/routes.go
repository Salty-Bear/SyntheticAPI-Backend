package routes

import (
	"github.com/Aryaman/pub-sub/routes/pubsub"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all main API routes
func RegisterRoutes(app *fiber.App) {
	pubsub.RegisterRoutes(app.Group("/pubsub"))
}
