package routes

import (
	"github.com/Aryaman/pub-sub/routes/connector"
	"github.com/Aryaman/pub-sub/routes/pubsub"
	"github.com/Aryaman/pub-sub/routes/tester"
	connectorService "github.com/Aryaman/pub-sub/services/connector"
	testerService "github.com/Aryaman/pub-sub/services/tester"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all main API routes
func RegisterRoutes(app *fiber.App) {
	// PubSub routes
	pubsub.RegisterRoutes(app.Group("/pubsub"))

	// Connector/Tunnel routes
	setupConnectorRoutes(app)

	// Testing routes
	setupTesterRoutes(app)
}

// setupConnectorRoutes initializes and registers connector routes
func setupConnectorRoutes(app *fiber.App) {
	// Initialize store and service
	store := connectorService.NewMemoryStore()
	config := connectorService.GenerateTunnelConfig("https://api.syntra.dev")
	service := connectorService.NewService(store, config)

	// Initialize tunnel manager and reverse proxy
	tunnelManager := connectorService.NewTunnelManager(store)
	reverseProxy := connectorService.NewReverseProxy(tunnelManager, "syntra.dev")

	// Create handlers
	handler := connector.NewHandler(service)
	wsHandler := connector.NewWebSocketHandler(tunnelManager, reverseProxy)

	// Register routes
	connectorGroup := app.Group("/tunnel")
	connector.RegisterRoutes(connectorGroup, handler, wsHandler)
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
