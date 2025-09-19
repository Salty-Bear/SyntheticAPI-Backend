package providers

import (
	"github.com/Aryaman/pub-sub/config"
	"github.com/gofiber/fiber/v2"
)

type Provider struct {
	S *Service
}

func InjectDefaultProviders(cnf config.AppConfig) (*Provider, error) {
	svcs := NewServicesWithConfig(cnf)
	return &Provider{
		S: svcs,
	}, nil
}

type keyType struct {
	key string
}

var providerKey = keyType{"providers"}
var globalProvider *Provider

func (p Provider) Handle(c *fiber.Ctx) error {
	c.Locals(providerKey, p)
	globalProvider = &p
	return c.Next()
}

func GetProviders(c *fiber.Ctx) Provider {
	return c.Locals(providerKey).(Provider)
}

// Overload for websocket.Conn (for use in WebSocket handlers)
func GetProvidersWS() Provider {
	if globalProvider != nil {
		return *globalProvider
	}
	panic("Provider not initialized")
}
