package providers

import (
	"github.com/Aryaman/syntra/config"
	"github.com/Aryaman/syntra/db"
	"github.com/gofiber/fiber/v2"
)

type Provider struct {
	S  *Service
	DB db.DB
}

func InjectDefaultProviders(cnf config.AppConfig) (*Provider, error) {
	// Initialize MongoDB connection
	mongoConn, err := db.NewMongoConnection(cnf.Database.MongoURL)
	if err != nil {
		return nil, err
	}

	// Initialize services with database connection
	svcs := NewServicesWithConfig(cnf, mongoConn)
	return &Provider{
		S:  svcs,
		DB: mongoConn,
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
