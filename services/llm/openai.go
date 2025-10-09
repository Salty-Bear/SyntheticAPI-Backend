package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Aryaman/syntra/config"
	"github.com/Aryaman/syntra/sdk"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// LangChainService implements LLMService using LangChain Go
type LangChainService struct {
	llm         llms.Model
	apiKey      string
	model       string
	temperature float64
	maxTokens   int
}

// NewLangChainService creates a new LangChain-based LLM service
func NewLangChainService(cfg config.LLM) (LLMService, error) {
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	// Create OpenAI LLM with LangChain
	llm, err := openai.New(
		openai.WithToken(cfg.OpenAIAPIKey),
		openai.WithModel(cfg.Model),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenAI LLM: %w", err)
	}

	return &LangChainService{
		llm:         llm,
		apiKey:      cfg.OpenAIAPIKey,
		model:       cfg.Model,
		temperature: cfg.Temperature,
		maxTokens:   cfg.MaxTokens,
	}, nil
}

// Chat sends a chat completion request
func (s *LangChainService) Chat(ctx context.Context, messages []sdk.LLMMessage) (string, error) {
	return s.chat(ctx, messages, false)
}

// ChatWithJSON sends a chat completion request expecting JSON response
func (s *LangChainService) ChatWithJSON(ctx context.Context, messages []sdk.LLMMessage) (string, error) {
	return s.chat(ctx, messages, true)
}

// chat internal method to handle chat requests
func (s *LangChainService) chat(ctx context.Context, messages []sdk.LLMMessage, jsonMode bool) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("no messages provided")
	}

	// Convert SDK messages to LangChain format
	var msgContents []llms.MessageContent
	for _, msg := range messages {
		var role llms.ChatMessageType
		switch msg.Role {
		case "system":
			role = llms.ChatMessageTypeSystem
		case "user":
			role = llms.ChatMessageTypeHuman
		case "assistant":
			role = llms.ChatMessageTypeAI
		default:
			role = llms.ChatMessageTypeHuman
		}

		msgContents = append(msgContents, llms.MessageContent{
			Role: role,
			Parts: []llms.ContentPart{
				llms.TextPart(msg.Content),
			},
		})
	}

	// Set up call options
	var opts []llms.CallOption
	if s.temperature > 0 {
		opts = append(opts, llms.WithTemperature(s.temperature))
	}
	if s.maxTokens > 0 {
		opts = append(opts, llms.WithMaxTokens(s.maxTokens))
	}

	// Add JSON mode if requested
	if jsonMode {
		opts = append(opts, llms.WithJSONMode())
	}

	// Generate response
	response, err := s.llm.GenerateContent(ctx, msgContents, opts...)
	if err != nil {
		log.Errorf("LangChain LLM generation failed: %v", err)
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	// Extract the text content from the response
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	content := response.Choices[0].Content

	// Log token information if available
	// Note: langchain-go doesn't expose token usage directly in ContentResponse
	// If you need detailed usage tracking, consider using the OpenAI SDK directly
	log.Debugf("LLM API response generated successfully")

	return content, nil
}

// GenerateWithRequest sends a request with full configuration
func (s *LangChainService) GenerateWithRequest(ctx context.Context, req *sdk.LLMRequest) (*sdk.LLMResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// Use custom parameters if provided, otherwise use defaults
	temperature := s.temperature
	if req.Temperature > 0 {
		temperature = req.Temperature
	}

	maxTokens := s.maxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	// Temporarily update service config for this request
	origTemp := s.temperature
	origMax := s.maxTokens
	s.temperature = temperature
	s.maxTokens = maxTokens

	// Call appropriate chat method
	var content string
	var err error
	if req.JSONMode {
		content, err = s.ChatWithJSON(ctx, req.Messages)
	} else {
		content, err = s.Chat(ctx, req.Messages)
	}

	// Restore original config
	s.temperature = origTemp
	s.maxTokens = origMax

	if err != nil {
		return nil, err
	}

	// Create response (note: token counts are approximate without direct API access)
	response := &sdk.LLMResponse{
		Content: content,
		Model:   s.model,
	}

	return response, nil
}

// Helper function to parse JSON response safely
func ParseJSONResponse(content string, target interface{}) error {
	if err := json.Unmarshal([]byte(content), target); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return nil
}
