package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Aryaman/syntra/sdk"
	"github.com/Aryaman/syntra/services/llm"
	"github.com/gofiber/fiber/v2/log"
)

// SyntheticDataAgent generates synthetic data based on schema and prompt
type SyntheticDataAgent struct {
	llmService llm.LLMService
}

// NewSyntheticDataAgent creates a new synthetic data generation agent
func NewSyntheticDataAgent(llmService llm.LLMService) sdk.Agent {
	return &SyntheticDataAgent{
		llmService: llmService,
	}
}

func (a *SyntheticDataAgent) Name() string {
	return "SyntheticDataAgent"
}

func (a *SyntheticDataAgent) Execute(ctx context.Context, input *sdk.AgentInput) (*sdk.AgentOutput, error) {
	log.Infof(LogMessages.AgentStarting, a.Name(), input.Count)

	// Prepare schema string
	schemaStr := "No specific schema provided. Generate generic data."
	if input.Schema != nil {
		schemaBytes, err := json.Marshal(input.Schema)
		if err == nil {
			schemaStr = string(schemaBytes)
		}
	}

	userPrompt := fmt.Sprintf(UserPromptTemplates.SyntheticData,
		input.Count,
		input.Prompt,
		schemaStr,
		input.DataType,
		input.Format,
	)

	messages := []sdk.LLMMessage{
		{Role: "system", Content: SystemPrompts.SyntheticDataAgent},
		{Role: "user", Content: userPrompt},
	}

	// Call LLM with JSON mode
	response, err := a.llmService.ChatWithJSON(ctx, messages)
	if err != nil {
		log.Errorf(LogMessages.LLMCallFailed, a.Name(), err)
		return nil, fmt.Errorf("%s: %w", ErrorMessages.GenerationFailed, err)
	}

	// Parse the JSON response
	var generatedData interface{}
	if err := json.Unmarshal([]byte(response), &generatedData); err != nil {
		log.Errorf(LogMessages.ParseResponseFailed, a.Name(), err)
		return nil, fmt.Errorf("%s: %w", ErrorMessages.ParseFailed, err)
	}

	log.Infof(LogMessages.AgentCompleted, a.Name(), "synthetic data")

	return &sdk.AgentOutput{
		Data: generatedData,
		Metadata: map[string]interface{}{
			"agent":             a.Name(),
			"records_requested": input.Count,
		},
	}, nil
}

// EdgeCaseAgent analyzes and generates edge case data
type EdgeCaseAgent struct {
	llmService llm.LLMService
}

// NewEdgeCaseAgent creates a new edge case generation agent
func NewEdgeCaseAgent(llmService llm.LLMService) sdk.Agent {
	return &EdgeCaseAgent{
		llmService: llmService,
	}
}

func (a *EdgeCaseAgent) Name() string {
	return "EdgeCaseAgent"
}

func (a *EdgeCaseAgent) Execute(ctx context.Context, input *sdk.AgentInput) (*sdk.AgentOutput, error) {
	log.Infof(LogMessages.AgentStarting, a.Name(), 0)

	// Prepare schema string
	schemaStr := "No specific schema provided."
	if input.Schema != nil {
		schemaBytes, err := json.Marshal(input.Schema)
		if err == nil {
			schemaStr = string(schemaBytes)
		}
	}

	// Include previously generated data for context
	previousDataStr := "No previous data available."
	if input.PreviousOutput != nil {
		prevBytes, err := json.Marshal(input.PreviousOutput)
		if err == nil {
			previousDataStr = string(prevBytes)
		}
	}

	userPrompt := fmt.Sprintf(UserPromptTemplates.EdgeCase,
		input.Prompt,
		schemaStr,
		input.DataType,
		input.Format,
		previousDataStr,
	)

	messages := []sdk.LLMMessage{
		{Role: "system", Content: SystemPrompts.EdgeCaseAgent},
		{Role: "user", Content: userPrompt},
	}

	// Call LLM with JSON mode
	response, err := a.llmService.ChatWithJSON(ctx, messages)
	if err != nil {
		log.Errorf(LogMessages.LLMCallFailed, a.Name(), err)
		return nil, fmt.Errorf("%s: %w", ErrorMessages.EdgeCaseGenFailed, err)
	}

	// Parse the JSON response
	var edgeCaseResponse interface{}
	if err := json.Unmarshal([]byte(response), &edgeCaseResponse); err != nil {
		log.Errorf(LogMessages.ParseResponseFailed, a.Name(), err)
		return nil, fmt.Errorf("%s: %w", ErrorMessages.EdgeCaseParseFailed, err)
	}

	log.Infof(LogMessages.AgentCompleted, a.Name(), "edge case data")

	return &sdk.AgentOutput{
		Data: edgeCaseResponse,
		Metadata: map[string]interface{}{
			"agent":             a.Name(),
			"analysis_complete": true,
		},
	}, nil
}
