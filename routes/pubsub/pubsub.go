package pubsub

import (
	"context"

	"github.com/Aryaman/pub-sub/providers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
)

// HandleWebSocket handles WebSocket connections for pub/sub operations
func HandleWebSocket(c *websocket.Conn) {
	log.Debug("received websocket connection request")
	ctx := context.Background()
	pr := providers.GetProvidersWS()
	pr.S.PubSub.HandleWebSocket(ctx, c)
}

// CreateTopic creates a new topic
func CreateTopic(c *fiber.Ctx) error {
	log.Debug("received create topic request")
	pr := providers.GetProviders(c)
	err := pr.S.PubSub.CreateTopic(c.Context(), c)
	if err != nil {
		log.Errorw("failed to create topic", "error", err)
		return err
	}
	log.Debug("topic created successfully")
	return nil
}

// DeleteTopic deletes an existing topic and disconnects all subscribers
func DeleteTopic(c *fiber.Ctx) error {
	log.Debug("received delete topic request")
	pr := providers.GetProviders(c)
	err := pr.S.PubSub.DeleteTopic(c.Context(), c)
	if err != nil {
		log.Errorw("failed to delete topic", "error", err)
		return err
	}
	log.Debug("topic deleted successfully")
	return nil
}

// ListTopics returns all available topics with subscriber counts
func ListTopics(c *fiber.Ctx) error {
	log.Debug("received list topics request")
	pr := providers.GetProviders(c)
	err := pr.S.PubSub.ListTopics(c.Context(), c)
	if err != nil {
		log.Errorw("failed to list topics", "error", err)
		return err
	}
	log.Debug("topics listed successfully")
	return nil
}

// Health returns system health information
func Health(c *fiber.Ctx) error {
	log.Debug("received health check request")
	pr := providers.GetProviders(c)
	err := pr.S.PubSub.Health(c.Context(), c)
	if err != nil {
		log.Errorw("health check failed", "error", err)
		return err
	}
	log.Debug("health check successful")
	return nil
}

// Stats returns system statistics including topic and message counts
func Stats(c *fiber.Ctx) error {
	log.Debug("received stats request")
	pr := providers.GetProviders(c)
	err := pr.S.PubSub.Stats(c.Context(), c)
	if err != nil {
		log.Errorw("failed to get stats", "error", err)
		return err
	}
	log.Debug("stats retrieved successfully")
	return nil
}
