package providers

import (
	"github.com/Aryaman/syntra/config"
	"github.com/Aryaman/syntra/db"
	"github.com/Aryaman/syntra/services/pubsub"
	"github.com/Aryaman/syntra/services/user"
)

type Service struct {
	PubSub pubsub.PubSub
	User   user.UserService
}

func NewServicesWithConfig(cnf config.AppConfig, database db.DB) *Service {
	pubsubSvc := pubsub.NewService(cnf.Server.MaxQueue, cnf.Server.MaxMessages)

	// Create user store and service
	userStore := user.NewStore(database)
	userSvc := user.NewService(userStore, cnf.Firebase.ProjectID)

	return &Service{
		PubSub: pubsubSvc,
		User:   userSvc,
	}
}

func NewServices(database db.DB) *Service {
	return NewServicesWithConfig(config.AppConfig{}, database)
}
