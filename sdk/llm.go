package sdk

// LLMMessage represents a chat message in LLM conversation
type LLMMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// LLMRequest represents a request to the LLM service
type LLMRequest struct {
	Messages    []LLMMessage `json:"messages"`
	Temperature float64      `json:"temperature,omitempty"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	JSONMode    bool         `json:"json_mode,omitempty"`
}

// LLMResponse represents a response from the LLM service
type LLMResponse struct {
	Content        string `json:"content"`
	PromptTokens   int    `json:"prompt_tokens"`
	ResponseTokens int    `json:"response_tokens"`
	TotalTokens    int    `json:"total_tokens"`
	Model          string `json:"model"`
}

// LLMConfig holds configuration for LLM service
type LLMConfig struct {
	Provider    string  `json:"provider"` // openai, anthropic, etc.
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}
