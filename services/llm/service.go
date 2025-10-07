package llm

import (
	"context"

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
