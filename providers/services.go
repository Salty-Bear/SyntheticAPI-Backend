package providers

import (
	"github.com/Aryaman/pub-sub/config"
	"github.com/Aryaman/pub-sub/db"
	"github.com/Aryaman/pub-sub/services/pubsub"
	"github.com/Aryaman/pub-sub/services/user"
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
