package config

// LLM contains the configuration for LLM services
type LLM struct {
	Provider     string  // "openai" or "gemini"
	OpenAIAPIKey string
	GeminiAPIKey string
	Model        string
	Temperature  float64
	MaxTokens    int
}
