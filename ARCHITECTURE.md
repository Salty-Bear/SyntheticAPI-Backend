# Architecture Overview

## Project Structure

```
Syntra/
├── sdk/                           # Shared types and interfaces
│   ├── agent_interface.go        # Agent interface definition
│   ├── agent.go                  # Agent input/output types
│   ├── llm.go                    # LLM types (Message, Request, Response)
│   ├── generate.go               # Generate service types
│   ├── user.go                   # User types
│   ├── tunnel.go                 # Tunnel types
│   └── pubsub_topic.go          # PubSub types
│
├── services/                      # Business logic services
│   ├── agents/
│   │   ├── agents.go             # SyntheticDataAgent & EdgeCaseAgent
│   │   └── orchestrator.go      # Multi-agent orchestration
│   ├── llm/
│   │   ├── service.go            # LLMService interface
│   │   └── openai.go            # LangChain-Go implementation
│   ├── generate/
│   │   └── service_impl.go      # Generate service implementation
│   └── ...
│
└── providers/
    └── services.go               # Service initialization & DI
```

## Data Flow

```
1. HTTP Request → Routes
2. Routes → Generate Service
3. Generate Service → Agent Orchestrator
4. Agent Orchestrator → Individual Agents
5. Agents → LLM Service (LangChain-Go)
6. LLM Service → OpenAI API
7. Response flows back through the chain
```

## Agent Flow

```
┌─────────────────────────────────────────────────┐
│          Generate Service (Execute)              │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│         Agent Orchestrator                       │
│  (ExecuteAgenticFlow)                           │
└──────┬──────────────────────────────────────────┘
       │
       ├──► Step 1: SyntheticDataAgent
       │    ├── Generates realistic synthetic data
       │    ├── Uses LangChain-Go LLM
       │    └── Returns AgentOutput with data
       │
       └──► Step 2: EdgeCaseAgent
            ├── Analyzes schema for edge cases
            ├── Receives context from Step 1
            ├── Uses LangChain-Go LLM
            └── Returns AgentOutput with edge cases
```

## LangChain Integration

```
┌─────────────────────────────────────────────────┐
│              Agent (Execute)                     │
└────────────────┬────────────────────────────────┘
                 │
                 │ sdk.LLMMessage[]
                 ▼
┌─────────────────────────────────────────────────┐
│        LangChainService (ChatWithJSON)          │
│  - Converts SDK messages to LangChain format    │
│  - Sets options (temperature, max_tokens)       │
│  - Enables JSON mode                            │
└────────────────┬────────────────────────────────┘
                 │
                 │ llms.MessageContent[]
                 ▼
┌─────────────────────────────────────────────────┐
│    langchaingo/llms/openai                      │
│  (GenerateContent)                              │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
          OpenAI API Response
```

## Key Interfaces

### Agent Interface (sdk/agent_interface.go)
```go
type Agent interface {
    Execute(ctx context.Context, input *AgentInput) (*AgentOutput, error)
    Name() string
}
```

### LLMService Interface (services/llm/service.go)
```go
type LLMService interface {
    Chat(ctx context.Context, messages []sdk.LLMMessage) (string, error)
    ChatWithJSON(ctx context.Context, messages []sdk.LLMMessage) (string, error)
    GenerateWithRequest(ctx context.Context, req *sdk.LLMRequest) (*sdk.LLMResponse, error)
}
```

## Benefits of This Architecture

1. **Separation of Concerns**: SDK types are separate from business logic
2. **Type Safety**: Strong typing with Go interfaces
3. **Testability**: Easy to mock interfaces for testing
4. **Extensibility**: Easy to add new agents or LLM providers
5. **Maintainability**: Clear structure and dependencies
6. **Reusability**: SDK types can be used across multiple services

## LangChain-Go Advantages

1. **Battle-tested**: Well-maintained library with community support
2. **Multiple Providers**: Easy to switch between OpenAI, Anthropic, etc.
3. **Advanced Features**: Built-in support for chains, agents, tools
4. **JSON Mode**: Native support for structured output
5. **Token Management**: Automatic token counting and optimization
