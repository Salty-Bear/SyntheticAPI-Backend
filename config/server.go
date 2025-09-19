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
