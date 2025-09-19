package providers

import (
	"github.com/Aryaman/pub-sub/config"
	"github.com/Aryaman/pub-sub/services/pubsub"
)

type Service struct {
	PubSub pubsub.PubSub
}

func NewServicesWithConfig(cnf config.AppConfig) *Service {
	pubsubSvc := pubsub.NewService(cnf.Server.MaxQueue, cnf.Server.MaxMessages)
	return &Service{PubSub: pubsubSvc}
}

func NewServices() *Service {
	return NewServicesWithConfig(config.AppConfig{})
}
