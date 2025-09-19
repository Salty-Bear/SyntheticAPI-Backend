package pubsub

import (
	"context"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type pubSubService struct {
	mu sync.Mutex
}

type PubSub interface {
	HandleWebSocket(ctx context.Context, c *websocket.Conn)
	CreateTopic(ctx context.Context, c *fiber.Ctx) error
	DeleteTopic(ctx context.Context, c *fiber.Ctx) error
	ListTopics(ctx context.Context, c *fiber.Ctx) error
	Health(ctx context.Context, c *fiber.Ctx) error
	Stats(ctx context.Context, c *fiber.Ctx) error
}
