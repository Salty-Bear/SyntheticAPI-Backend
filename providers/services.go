package providers

import (
	"github.com/Aryaman/syntra/config"
	"github.com/Aryaman/syntra/db"
	user "github.com/Aryaman/syntra/services/User"
	"github.com/Aryaman/syntra/services/generate"
	"github.com/Aryaman/syntra/services/llm"
	"github.com/Aryaman/syntra/services/pubsub"
	"github.com/Aryaman/syntra/services/tunnel"
)

type Service struct {
	PubSub   pubsub.PubSub
	User     user.UserService
	Tunnel   tunnel.TunnelService
	Generate generate.GenerateService
	LLM      llm.LLMService
}

func NewServicesWithConfig(cnf config.AppConfig, database db.DB) *Service {
	pubsubSvc := pubsub.NewService(cnf.Server.MaxQueue, cnf.Server.MaxMessages)

	// Create user store and service
	userStore := user.NewStore(database)
	userSvc := user.NewService(userStore)

	// Create tunnel store and service
	tunnelStore := tunnel.NewStore(database)
	tunnelSvc := tunnel.NewService(tunnelStore)

	// Create LLM service
	llmSvc, err := llm.NewLangChainService(cnf.LLM)
	if err != nil {
		panic(err) // In production, handle this error more gracefully
	}

	// Create generate store and service with LLM
	generateStore := generate.NewStore(database)
	generateSvc := generate.NewService(generateStore, llmSvc)

	return &Service{
		PubSub:   pubsubSvc,
		User:     userSvc,
		Tunnel:   tunnelSvc,
		Generate: generateSvc,
		LLM:      llmSvc,
	}
}

func NewServices(database db.DB) *Service {
	return NewServicesWithConfig(config.AppConfig{}, database)
}
