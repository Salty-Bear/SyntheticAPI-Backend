package config

type Server struct {
	Host        string
	Port        string
	MaxQueue    int
	MaxMessages int
}

type Deployment struct {
	Environment string
	Name        string
}

type Firebase struct {
	ProjectID string
}

type Database struct {
	MongoURL string
}
