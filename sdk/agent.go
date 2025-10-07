package sdk

// AgentInput represents input to an agent
type AgentInput struct {
	Prompt         string                 `json:"prompt"`
	Schema         interface{}            `json:"schema"`
	Count          int                    `json:"count"`
	DataType       string                 `json:"data_type"`
	Format         string                 `json:"format"`
	Context        map[string]interface{} `json:"context,omitempty"`
	PreviousOutput interface{}            `json:"previous_output,omitempty"`
}

// AgentOutput represents output from an agent
type AgentOutput struct {
	Data     interface{}            `json:"data"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Error    error                  `json:"error,omitempty"`
}

// OrchestrationResult contains the combined results from all agents
type OrchestrationResult struct {
	SyntheticData interface{}            `json:"synthetic_data"`
	EdgeCaseData  interface{}            `json:"edge_case_data"`
	Metadata      map[string]interface{} `json:"metadata"`
	Success       bool                   `json:"success"`
	Error         string                 `json:"error,omitempty"`
}

// AgentType represents the type of agent
type AgentType string

const (
	AgentTypeSynthetic AgentType = "synthetic"
	AgentTypeEdgeCase  AgentType = "edge_case"
)

// AgentConfig holds configuration for an agent
type AgentConfig struct {
	Type        AgentType `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
}
