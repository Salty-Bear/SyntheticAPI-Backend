package config

// LLM contains the configuration for LLM services
type LLM struct {
	OpenAIAPIKey string
	Model        string
	Temperature  float64
	MaxTokens    int
}
