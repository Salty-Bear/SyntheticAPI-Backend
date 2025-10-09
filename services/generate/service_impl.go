package generate

import (
	"context"
	"fmt"
	"time"

	"github.com/Aryaman/syntra/sdk"
	"github.com/Aryaman/syntra/services/agents"
	"github.com/Aryaman/syntra/services/llm"
	"github.com/google/uuid"
)

// ServiceImpl implements the GenerateService interface
type ServiceImpl struct {
	store        GenerateStore
	llmService   llm.LLMService
	orchestrator *agents.AgentOrchestrator
}

// NewService creates a new instance of GenerateService
func NewService(store GenerateStore, llmService llm.LLMService) GenerateService {
	return &ServiceImpl{
		store:        store,
		llmService:   llmService,
		orchestrator: agents.NewAgentOrchestrator(llmService),
	}
}

// Create creates a new generate task or returns existing generate task if name already exists
func (s *ServiceImpl) Create(ctx context.Context, generate *sdk.Generate) (*sdk.Generate, error) {
	// Generate ID if not provided
	if generate.ID == "" {
		generate.ID = uuid.New().String()
	}

	// Validate required fields
	if generate.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Check if generate task with name already exists for this user
	exists, err := s.store.NameExists(ctx, generate.Name, generate.UserId)
	if err != nil {
		return nil, fmt.Errorf("error checking generate task name existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("generate task with name '%s' already exists for this user", generate.Name)
	}

	// Set default values
	if generate.Status == "" {
		generate.Status = "pending"
	}
	if generate.DataType == "" {
		generate.DataType = "json"
	}
	if generate.Count <= 0 {
		generate.Count = 10
	}

	// Set timestamps
	now := time.Now()
	generate.CreatedAt = now
	generate.UpdatedAt = now

	// Set created by if not provided
	if generate.CreatedBy == "" {
		generate.CreatedBy = generate.UserId
	}

	// Create database model from SDK generate
	dbGenerate := fromSdkToModel(*generate)

	// Save generate task to database
	if err := s.store.CreateGenerate(ctx, &dbGenerate); err != nil {
		return nil, fmt.Errorf("error creating generate task: %w", err)
	}

	// Convert back to SDK type to return
	return fromModelToSdk(&dbGenerate), nil
}

// Get retrieves a generate task by ID
func (s *ServiceImpl) Get(ctx context.Context, id string, userId string) (*sdk.Generate, error) {
	if id == "" {
		return nil, fmt.Errorf("generate task ID is required")
	}
	if userId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	generate, err := s.store.GetGenerateByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving generate task: %w", err)
	}

	if generate == nil {
		return nil, fmt.Errorf("generate task not found")
	}

	// Check if the generate task belongs to the user
	if generate.UserId != userId {
		return nil, fmt.Errorf("generate task not found")
	}

	// Convert database model to SDK type
	return fromModelToSdk(generate), nil
}

// List retrieves all generate tasks with optional filtering and pagination
func (s *ServiceImpl) List(ctx context.Context, query *sdk.GenerateQuery) ([]*sdk.Generate, error) {
	var (
		status   string
		dataType string
		enabled  *bool
		userId   string
	)

	if query != nil {
		status = query.Status
		dataType = query.DataType
		enabled = query.Enabled
		userId = query.UserId
	}

	if userId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	generates, err := s.store.GetGeneratesByUser(ctx, userId, status, dataType, enabled)
	if err != nil {
		return nil, fmt.Errorf("error retrieving generate tasks: %w", err)
	}

	// Apply pagination if specified
	if query != nil && query.Limit > 0 {
		start := query.Offset
		end := start + query.Limit

		if start >= len(generates) {
			return []*sdk.Generate{}, nil
		}

		if end > len(generates) {
			end = len(generates)
		}

		generates = generates[start:end]
	}

	// Convert database models to SDK types
	var sdkGenerates []*sdk.Generate
	for _, generate := range generates {
		sdkGenerates = append(sdkGenerates, fromModelToSdk(generate))
	}

	return sdkGenerates, nil
}

// Update updates an existing generate task
func (s *ServiceImpl) Update(ctx context.Context, generate *sdk.Generate) (*sdk.Generate, error) {
	if generate.ID == "" {
		return nil, fmt.Errorf("generate task ID is required")
	}
	if generate.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Check if generate task exists and belongs to user
	existing, err := s.Get(ctx, generate.ID, generate.UserId)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("generate task not found")
	}

	// Build updates map
	updates := make(map[string]interface{})
	if generate.Name != "" {
		updates["name"] = generate.Name
	}
	if generate.Description != "" {
		updates["description"] = generate.Description
	}
	if generate.DataType != "" {
		updates["data_type"] = generate.DataType
	}
	if generate.Count > 0 {
		updates["count"] = generate.Count
	}
	if generate.Schema != nil {
		updates["schema"] = generate.Schema
	}
	if generate.Format != "" {
		updates["format"] = generate.Format
	}
	if generate.Status != "" {
		updates["status"] = generate.Status
	}
	updates["enabled"] = generate.Enabled
	if generate.OutputData != nil {
		updates["output_data"] = generate.OutputData
	}

	// Update timestamps
	updates["updated_at"] = time.Now()
	if generate.UpdatedBy != "" {
		updates["updated_by"] = generate.UpdatedBy
	} else {
		updates["updated_by"] = generate.UserId
	}

	// Update generate task
	if err := s.store.UpdateGenerate(ctx, generate.ID, updates); err != nil {
		return nil, fmt.Errorf("error updating generate task: %w", err)
	}

	// Retrieve and return updated generate task
	updatedGenerate, err := s.store.GetGenerateByID(ctx, generate.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving updated generate task: %w", err)
	}

	// Convert to SDK type
	return fromModelToSdk(updatedGenerate), nil
}

// Delete deletes a generate task by ID
func (s *ServiceImpl) Delete(ctx context.Context, id string, userId string) error {
	if id == "" {
		return fmt.Errorf("generate task ID is required")
	}
	if userId == "" {
		return fmt.Errorf("user_id is required")
	}

	// Check if generate task exists and belongs to user
	_, err := s.Get(ctx, id, userId)
	if err != nil {
		return err
	}

	// Delete generate task
	if err := s.store.DeleteGenerate(ctx, id); err != nil {
		return fmt.Errorf("error deleting generate task: %w", err)
	}

	return nil
}

// Execute executes a generate task to produce synthetic data using agentic flow
func (s *ServiceImpl) Execute(ctx context.Context, id string, userId string) (*sdk.Generate, error) {
	if id == "" {
		return nil, fmt.Errorf("generate task ID is required")
	}
	if userId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Get the generate task
	generate, err := s.Get(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	// Check if task is enabled
	if !generate.Enabled {
		return nil, fmt.Errorf("generate task is disabled")
	}

	// Update status to active
	generate.Status = "active"
	generate.UpdatedAt = time.Now()
	generate.UpdatedBy = userId

	// Save initial status update
	s.Update(ctx, generate)

	// Prepare prompt for the agents
	prompt := generate.Description
	if prompt == "" {
		prompt = fmt.Sprintf("Generate %s data with %d records", generate.DataType, generate.Count)
	}

	// Execute agentic flow
	agentInput := &sdk.AgentInput{
		Prompt:   prompt,
		Schema:   generate.Schema,
		Count:    generate.Count,
		DataType: generate.DataType,
		Format:   generate.Format,
	}

	// Run the multi-agent orchestration
	result, err := s.orchestrator.ExecuteAgenticFlow(ctx, agentInput)
	if err != nil {
		// If execution fails, set status to failed
		generate.Status = "failed"
		generate.OutputData = map[string]interface{}{
			"error": err.Error(),
		}
		s.Update(ctx, generate)
		return nil, fmt.Errorf("agentic flow execution failed: %w", err)
	}

	// Combine results from all agents
	outputData := s.orchestrator.CombineResults(result)
	generate.OutputData = outputData

	// Update status to completed
	generate.Status = "completed"

	// Update the generate task with the results
	updatedGenerate, err := s.Update(ctx, generate)
	if err != nil {
		// If update fails, set status to failed
		generate.Status = "failed"
		s.Update(ctx, generate)
		return nil, fmt.Errorf("error updating generate task with results: %w", err)
	}

	return updatedGenerate, nil
}

// generateSampleData creates sample data based on the generate task configuration
func (s *ServiceImpl) generateSampleData(generate *sdk.Generate) interface{} {
	// This is a simple implementation. In a real application, you would
	// implement sophisticated data generation based on the schema and format

	switch generate.DataType {
	case "json":
		var data []map[string]interface{}
		for i := 0; i < generate.Count; i++ {
			data = append(data, map[string]interface{}{
				"id":    i + 1,
				"name":  fmt.Sprintf("Generated Item %d", i+1),
				"value": fmt.Sprintf("sample_value_%d", i+1),
			})
		}
		return data
	case "csv":
		return "id,name,value\n1,Generated Item 1,sample_value_1\n2,Generated Item 2,sample_value_2"
	case "xml":
		return "<data><item id=\"1\"><name>Generated Item 1</name><value>sample_value_1</value></item></data>"
	default:
		return map[string]interface{}{
			"message": "Data generated successfully",
			"count":   generate.Count,
			"type":    generate.DataType,
		}
	}
}
