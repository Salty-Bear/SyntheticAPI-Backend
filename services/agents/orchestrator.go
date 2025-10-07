package agents

import (
	"context"
	"fmt"

	"github.com/Aryaman/syntra/sdk"
	"github.com/Aryaman/syntra/services/llm"
	"github.com/gofiber/fiber/v2/log"
)

// AgentOrchestrator coordinates multiple agents to complete a task
type AgentOrchestrator struct {
	llmService llm.LLMService
	agents     map[string]sdk.Agent
}

// NewAgentOrchestrator creates a new agent orchestrator
func NewAgentOrchestrator(llmService llm.LLMService) *AgentOrchestrator {
	return &AgentOrchestrator{
		llmService: llmService,
		agents: map[string]sdk.Agent{
			"synthetic": NewSyntheticDataAgent(llmService),
			"edge_case": NewEdgeCaseAgent(llmService),
		},
	}
}

// ExecuteAgenticFlow runs the multi-agent workflow
func (o *AgentOrchestrator) ExecuteAgenticFlow(ctx context.Context, input *sdk.AgentInput) (*sdk.OrchestrationResult, error) {
	log.Info("Starting agentic flow orchestration")

	result := &sdk.OrchestrationResult{
		Metadata: make(map[string]interface{}),
		Success:  false,
	}

	// Step 1: Execute Synthetic Data Agent
	log.Info("Step 1: Executing Synthetic Data Agent")
	syntheticAgent := o.agents["synthetic"]
	syntheticOutput, err := syntheticAgent.Execute(ctx, input)
	if err != nil {
		errMsg := fmt.Sprintf("Synthetic data agent failed: %v", err)
		log.Error(errMsg)
		result.Error = errMsg
		return result, err
	}

	result.SyntheticData = syntheticOutput.Data
	result.Metadata["synthetic_data_metadata"] = syntheticOutput.Metadata
	log.Info("Step 1: Completed - Synthetic data generated successfully")

	// Step 2: Execute Edge Case Agent with context from synthetic data
	log.Info("Step 2: Executing Edge Case Agent")
	edgeCaseInput := input
	edgeCaseInput.PreviousOutput = syntheticOutput.Data
	edgeCaseInput.Context = map[string]interface{}{
		"synthetic_data_generated": true,
		"synthetic_metadata":       syntheticOutput.Metadata,
	}

	edgeCaseAgent := o.agents["edge_case"]
	edgeCaseOutput, err := edgeCaseAgent.Execute(ctx, edgeCaseInput)
	if err != nil {
		// Edge case generation is not critical, log but don't fail
		log.Warnf("Edge case agent failed (non-critical): %v", err)
		result.Metadata["edge_case_error"] = err.Error()
		// Still mark as success since synthetic data was generated
		result.Success = true
		return result, nil
	}

	result.EdgeCaseData = edgeCaseOutput.Data
	result.Metadata["edge_case_metadata"] = edgeCaseOutput.Metadata
	log.Info("Step 2: Completed - Edge case data generated successfully")

	// Mark orchestration as successful
	result.Success = true
	result.Metadata["orchestration_complete"] = true
	result.Metadata["total_agents_executed"] = 2

	log.Info("Agentic flow orchestration completed successfully")
	return result, nil
}

// CombineResults combines synthetic and edge case data into a single output
func (o *AgentOrchestrator) CombineResults(result *sdk.OrchestrationResult) interface{} {
	combined := map[string]interface{}{
		"synthetic_data": result.SyntheticData,
		"edge_case_data": result.EdgeCaseData,
		"metadata": map[string]interface{}{
			"generation_metadata": result.Metadata,
			"success":             result.Success,
		},
	}

	// If we have edge case data with the expected structure, extract the actual data
	if edgeCaseMap, ok := result.EdgeCaseData.(map[string]interface{}); ok {
		if edgeCaseArray, ok := edgeCaseMap["edge_case_data"]; ok {
			combined["edge_case_data"] = edgeCaseArray
		}
		if scenarios, ok := edgeCaseMap["edge_cases_identified"]; ok {
			combined["edge_case_scenarios"] = scenarios
		}
	}

	return combined
}
