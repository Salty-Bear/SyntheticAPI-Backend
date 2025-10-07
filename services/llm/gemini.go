package llm

import (
	"context"
	"fmt"

	"github.com/Aryaman/syntra/config"
	"github.com/Aryaman/syntra/sdk"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

// GeminiService implements LLMService using Google Gemini
type GeminiService struct {
	llm         llms.Model
	apiKey      string
	model       string
	temperature float64
	maxTokens   int
}

// NewGeminiService creates a new Gemini-based LLM service
func NewGeminiService(cfg config.LLM) (LLMService, error) {
	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("Gemini API key not configured")
	}

	// Create Gemini LLM with LangChain
	llm, err := googleai.New(
		context.Background(),
		googleai.WithAPIKey(cfg.GeminiAPIKey),
		googleai.WithDefaultModel(cfg.Model),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini LLM: %w", err)
	}

	return &GeminiService{
		llm:         llm,
		apiKey:      cfg.GeminiAPIKey,
		model:       cfg.Model,
		temperature: cfg.Temperature,
		maxTokens:   cfg.MaxTokens,
	}, nil
}

// Chat sends a chat completion request
func (s *GeminiService) Chat(ctx context.Context, messages []sdk.LLMMessage) (string, error) {
	return s.chat(ctx, messages, false)
}

// ChatWithJSON sends a chat completion request expecting JSON response
func (s *GeminiService) ChatWithJSON(ctx context.Context, messages []sdk.LLMMessage) (string, error) {
	return s.chat(ctx, messages, true)
}

// chat internal method to handle chat requests
func (s *GeminiService) chat(ctx context.Context, messages []sdk.LLMMessage, jsonMode bool) (string, error) {
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
	// Note: Gemini doesn't have a direct JSON mode like OpenAI, but we can guide it via prompts
	if jsonMode {
		// Append instruction to the last user message
		if len(msgContents) > 0 {
			lastIdx := len(msgContents) - 1
			if msgContents[lastIdx].Role == llms.ChatMessageTypeHuman {
				// Get the text content
				if len(msgContents[lastIdx].Parts) > 0 {
					if textPart, ok := msgContents[lastIdx].Parts[0].(llms.TextContent); ok {
						msgContents[lastIdx].Parts[0] = llms.TextPart(textPart.Text + "\n\nIMPORTANT: Respond ONLY with valid JSON. Do not include any explanation, markdown formatting, or additional text.")
					}
				}
			}
		}
	}

	// Generate response
	response, err := s.llm.GenerateContent(ctx, msgContents, opts...)
	if err != nil {
		log.Errorf("Gemini LLM generation failed: %v", err)
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	// Extract the text content from the response
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	content := response.Choices[0].Content

	log.Debugf("Gemini API response generated successfully")

	return content, nil
}

// GenerateWithRequest sends a request with full configuration
func (s *GeminiService) GenerateWithRequest(ctx context.Context, req *sdk.LLMRequest) (*sdk.LLMResponse, error) {
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

	// Create response
	response := &sdk.LLMResponse{
		Content: content,
		Model:   s.model,
	}

	return response, nil
}
