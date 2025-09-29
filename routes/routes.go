package routes

import (
	"github.com/Aryaman/syntra/providers"
	"github.com/Aryaman/syntra/routes/pubsub"
	"github.com/Aryaman/syntra/routes/tester"
	"github.com/Aryaman/syntra/routes/tunnel"
	"github.com/Aryaman/syntra/routes/user"
	testerService "github.com/Aryaman/syntra/services/tester"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all main API routes
func RegisterRoutes(app *fiber.App, provider *providers.Provider) {
	// PubSub routes
	pubsub.RegisterRoutes(app.Group("/pubsub"))

	// User routes
	user.RegisterRoutes(app.Group("/users"))

	// Tunnel routes
	tunnel.RegisterRoutes(app.Group("/tunnels"))

	// Testing routes
	setupTesterRoutes(app)
}

// setupTesterRoutes initializes and registers testing routes
func setupTesterRoutes(app *fiber.App) {
	// Initialize store and service
	store := testerService.NewMemoryStore()
	service := testerService.NewService(store)

	// Create handler
	handler := tester.NewHandler(service)

	// Register routes
	testerGroup := app.Group("/test")
	tester.RegisterRoutes(testerGroup, handler)
}
