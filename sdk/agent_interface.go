package sdk

import "context"

// Agent represents a specialized AI agent
type Agent interface {
	Execute(ctx context.Context, input *AgentInput) (*AgentOutput, error)
	Name() string
}
