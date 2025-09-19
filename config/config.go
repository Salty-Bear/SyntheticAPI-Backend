package config

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server     Server
	Deployment Deployment
}

func NewAppConfig() *AppConfig {
	cnf := &AppConfig{}
	cnf.Load()
	return cnf
}

type keyType struct {
	key string
}

var configKey = keyType{"config"}

func (a *AppConfig) Handle(c *fiber.Ctx) error {
	c.Locals(configKey, a)
	return c.Next()
}

func GetAppConfig(c *fiber.Ctx) AppConfig {
	return c.Locals(configKey).(AppConfig)
}

func (a *AppConfig) Load() {
	/*
	 * load env file
	 * load each config one by one
	 */
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file found. Using default environment values")
	}
	a.LoadServerConfig()
	a.LoadDeploymentConfig()
}

// LoadServerConfig load server config
func (a *AppConfig) LoadServerConfig() {
	// load the default values
	// then load from env variables
	a.Server.Host = "localhost"
	a.Server.Port = "3000"
	a.Server.MaxQueue = 100
	a.Server.MaxMessages = 100

	host := os.Getenv("SERVER_HOST")
	if host != "" {
		a.Server.Host = host
	}
	port := os.Getenv("SERVER_PORT")
	if port != "" {
		a.Server.Port = port
	}
	if mq := os.Getenv("MAX_QUEUE"); mq != "" {
		if v, err := strconv.Atoi(mq); err == nil {
			a.Server.MaxQueue = v
		}
	}
	if mm := os.Getenv("MAX_MESSAGES"); mm != "" {
		if v, err := strconv.Atoi(mm); err == nil {
			a.Server.MaxMessages = v
		}
	}
}

// LoadDeploymentConfig loads the deployment config
func (a *AppConfig) LoadDeploymentConfig() {
	// load the default values
	// then load from env variables
	a.Deployment.Environment = "development"
	a.Deployment.Name = "Pub Sub Service"

	environment := os.Getenv("DEPLOYMENT_ENVIRONMENT")
	if environment != "" {
		a.Deployment.Environment = environment
	}

	name := os.Getenv("DEPLOYMENT_NAME")
	if name != "" {
		a.Deployment.Name = name
	}
}
