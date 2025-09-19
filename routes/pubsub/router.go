package pubsub

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Use controller functions from pubsub.go
func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")
	v1.Post("/topics", CreateTopic)
	v1.Delete("/topics/:name", DeleteTopic)
	v1.Get("/topics", ListTopics)
	v1.Get("/health", Health)
	v1.Get("/stats", Stats)
	v1.Get("/ws", websocket.New(HandleWebSocket))
}
