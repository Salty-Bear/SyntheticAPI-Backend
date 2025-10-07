package llm

import (
	"context"
	"fmt"

	"github.com/Aryaman/syntra/config"
	"github.com/Aryaman/syntra/sdk"
)

// LLMService defines the interface for LLM operations
type LLMService interface {
	// Chat sends a chat completion request
	Chat(ctx context.Context, messages []sdk.LLMMessage) (string, error)

	// ChatWithJSON sends a chat completion request and expects JSON response
	ChatWithJSON(ctx context.Context, messages []sdk.LLMMessage) (string, error)

	// GenerateWithRequest sends a request with full configuration
	GenerateWithRequest(ctx context.Context, req *sdk.LLMRequest) (*sdk.LLMResponse, error)
}

// NewLLMService creates a new LLM service based on the provider configuration
func NewLLMService(cfg config.LLM) (LLMService, error) {
	switch cfg.Provider {
	case "gemini":
		return NewGeminiService(cfg)
	case "openai":
		return NewLangChainService(cfg)
	default:
		// Default to OpenAI if no provider specified
		if cfg.OpenAIAPIKey != "" {
			return NewLangChainService(cfg)
		} else if cfg.GeminiAPIKey != "" {
			return NewGeminiService(cfg)
		}
		return nil, fmt.Errorf("no LLM provider configured. Please set LLM_PROVIDER to 'openai' or 'gemini'")
	}
}
