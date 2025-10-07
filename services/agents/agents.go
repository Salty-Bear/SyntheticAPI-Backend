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
	log.Infof("[%s] Starting execution for %d records", a.Name(), input.Count)

	// Build the prompt for synthetic data generation
	systemPrompt := `You are an expert synthetic data generator. Your task is to generate realistic, high-quality synthetic data that matches the provided schema and requirements.

Rules:
1. Generate data that strictly adheres to the provided schema structure
2. Ensure data is realistic and diverse
3. Follow any specific instructions in the user prompt
4. Return ONLY valid JSON without any additional text or explanations
5. Generate exactly the number of records requested`

	// Prepare schema string
	schemaStr := "No specific schema provided. Generate generic data."
	if input.Schema != nil {
		schemaBytes, err := json.Marshal(input.Schema)
		if err == nil {
			schemaStr = string(schemaBytes)
		}
	}

	userPrompt := fmt.Sprintf(`Generate %d synthetic data records with the following requirements:

Prompt: %s

Schema: %s

Data Type: %s
Format: %s

Return the data as a JSON array of objects. Each object should match the schema structure.`,
		input.Count,
		input.Prompt,
		schemaStr,
		input.DataType,
		input.Format,
	)

	messages := []sdk.LLMMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// Call LLM with JSON mode
	response, err := a.llmService.ChatWithJSON(ctx, messages)
	if err != nil {
		log.Errorf("[%s] LLM call failed: %v", a.Name(), err)
		return nil, fmt.Errorf("failed to generate synthetic data: %w", err)
	}

	// Parse the JSON response
	var generatedData interface{}
	if err := json.Unmarshal([]byte(response), &generatedData); err != nil {
		log.Errorf("[%s] Failed to parse response: %v", a.Name(), err)
		return nil, fmt.Errorf("failed to parse generated data: %w", err)
	}

	log.Infof("[%s] Successfully generated synthetic data", a.Name())

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
	log.Infof("[%s] Starting edge case analysis and generation", a.Name())

	// Build the prompt for edge case generation
	systemPrompt := `You are an expert at identifying and generating edge case scenarios. Your task is to analyze the provided schema and generate comprehensive edge case test data that covers all boundary conditions, unusual inputs, and potential error scenarios.

Rules:
1. Identify ALL possible edge cases for the given schema
2. Generate data that tests boundaries (min/max values, empty strings, nulls, special characters, etc.)
3. Include corner cases that might break typical validations
4. Cover different data type edge cases (very long strings, negative numbers, special dates, etc.)
5. Return ONLY valid JSON without any additional text or explanations
6. Structure your response to include both the edge case scenarios and the generated data`

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

	userPrompt := fmt.Sprintf(`Analyze the schema and generate comprehensive edge case test data:

Original Prompt: %s

Schema: %s

Data Type: %s
Format: %s

Previously Generated Data (for context): %s

Generate a JSON response with the following structure:
{
  "edge_cases_identified": [
    {
      "category": "category name",
      "description": "what edge case this tests",
      "scenarios": ["scenario 1", "scenario 2"]
    }
  ],
  "edge_case_data": [
    // Array of objects matching the schema but with edge case values
  ]
}

Generate at least 15-20 edge case records covering various scenarios like:
- Boundary values (min/max)
- Empty or null values
- Special characters and encoding
- Extremely long/short values
- Invalid but parseable data
- Unusual but valid combinations`,
		input.Prompt,
		schemaStr,
		input.DataType,
		input.Format,
		previousDataStr,
	)

	messages := []sdk.LLMMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	// Call LLM with JSON mode
	response, err := a.llmService.ChatWithJSON(ctx, messages)
	if err != nil {
		log.Errorf("[%s] LLM call failed: %v", a.Name(), err)
		return nil, fmt.Errorf("failed to generate edge cases: %w", err)
	}

	// Parse the JSON response
	var edgeCaseResponse interface{}
	if err := json.Unmarshal([]byte(response), &edgeCaseResponse); err != nil {
		log.Errorf("[%s] Failed to parse response: %v", a.Name(), err)
		return nil, fmt.Errorf("failed to parse edge case data: %w", err)
	}

	log.Infof("[%s] Successfully generated edge case data", a.Name())

	return &sdk.AgentOutput{
		Data: edgeCaseResponse,
		Metadata: map[string]interface{}{
			"agent":             a.Name(),
			"analysis_complete": true,
		},
	}, nil
}
