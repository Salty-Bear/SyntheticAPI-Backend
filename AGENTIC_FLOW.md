# Agentic Flow for Synthetic Data Generation

This document describes the agentic flow implementation for generating synthetic data using multiple AI agents.

## Overview

The system uses a multi-agent architecture powered by Large Language Models (LLMs) to generate high-quality synthetic data. Two specialized agents work in sequence:

1. **Synthetic Data Agent** - Generates realistic synthetic data based on the provided schema and requirements
2. **Edge Case Agent** - Analyzes the schema and generates comprehensive edge case test data

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  Agent Orchestrator                      │
│                                                          │
│  ┌────────────────────┐      ┌─────────────────────┐   │
│  │ Synthetic Data     │──────>│  Edge Case          │   │
│  │ Agent              │      │  Agent              │   │
│  │                    │      │                     │   │
│  │ - Generate data    │      │ - Analyze schema    │   │
│  │ - Follow schema    │      │ - Find boundaries   │   │
│  │ - Match count      │      │ - Generate edge     │   │
│  │                    │      │   cases             │   │
│  └────────────────────┘      └─────────────────────┘   │
│                                                          │
└─────────────────────────────────────────────────────────┘
                        │
                        ▼
            ┌─────────────────────┐
            │  Combined Output     │
            │  - Synthetic Data    │
            │  - Edge Case Data    │
            │  - Metadata          │
            └─────────────────────┘
```

## Agent Details

### 1. Synthetic Data Agent

**Purpose:** Generate realistic synthetic data that adheres to the provided schema.

**Capabilities:**
- Generates data matching the exact schema structure
- Creates diverse and realistic data
- Follows specific user instructions
- Produces exactly the requested number of records
- Supports multiple data types (JSON, CSV, XML, SQL, YAML)

**Output:**
```json
{
  "data": [...array of generated records...],
  "metadata": {
    "agent": "SyntheticDataAgent",
    "records_requested": 10
  }
}
```

### 2. Edge Case Agent

**Purpose:** Identify and generate comprehensive edge case test data.

**Capabilities:**
- Identifies all possible edge cases for the schema
- Tests boundary conditions (min/max values)
- Covers unusual inputs and error scenarios
- Generates data with:
  - Empty/null values
  - Special characters
  - Extremely long/short values
  - Invalid but parseable data
  - Unusual but valid combinations

**Output:**
```json
{
  "edge_cases_identified": [
    {
      "category": "Boundary Values",
      "description": "Testing min/max constraints",
      "scenarios": ["minimum values", "maximum values"]
    }
  ],
  "edge_case_data": [...array of edge case records...]
}
```

## Configuration

### Environment Variables

Add the following to your `.env` file:

```bash
# OpenAI API Configuration
OPENAI_API_KEY=your_openai_api_key_here

# LLM Model Settings (optional, defaults shown)
LLM_MODEL=gpt-4o-mini          # Model to use
LLM_TEMPERATURE=0.7            # Creativity level (0.0-1.0)
LLM_MAX_TOKENS=4096            # Maximum tokens per request
```

### Supported Models

- `gpt-4o-mini` (default, cost-effective)
- `gpt-4o` (more capable, higher cost)
- `gpt-4-turbo`
- `gpt-3.5-turbo`

## API Usage

### Create a Generate Task

```bash
POST /v1/generate
```

**Request Body:**
```json
{
  "name": "User Data Generator",
  "description": "Generate realistic user profile data",
  "data_type": "json",
  "count": 10,
  "schema": {
    "id": "string",
    "name": "string",
    "email": "string",
    "age": "integer",
    "address": {
      "street": "string",
      "city": "string",
      "zipcode": "string"
    }
  },
  "format": "pretty",
  "enabled": true
}
```

### Execute the Task

```bash
POST /v1/generate/:id/execute
```

This will trigger the agentic flow:
1. Synthetic Data Agent generates the main dataset
2. Edge Case Agent generates edge cases based on the schema
3. Results are combined and returned

**Response:**
```json
{
  "success": true,
  "message": "Generate task executed successfully",
  "data": {
    "id": "task-id",
    "name": "User Data Generator",
    "status": "completed",
    "output_data": {
      "synthetic_data": [...],
      "edge_case_data": [...],
      "edge_case_scenarios": [...],
      "metadata": {
        "generation_metadata": {
          "orchestration_complete": true,
          "total_agents_executed": 2
        },
        "success": true
      }
    }
  }
}
```

## Example Use Cases

### 1. E-commerce Product Data

```json
{
  "name": "Product Catalog",
  "description": "Generate product data with prices, descriptions, and inventory",
  "data_type": "json",
  "count": 50,
  "schema": {
    "product_id": "string",
    "name": "string",
    "description": "string",
    "price": "float",
    "stock": "integer",
    "category": "string"
  }
}
```

### 2. User Authentication Testing

```json
{
  "name": "Auth Test Data",
  "description": "Generate user credentials with various edge cases",
  "data_type": "json",
  "count": 20,
  "schema": {
    "username": "string",
    "email": "string",
    "password": "string",
    "role": "enum[admin,user,guest]"
  }
}
```

### 3. Financial Transaction Data

```json
{
  "name": "Transaction Records",
  "description": "Generate financial transaction data",
  "data_type": "json",
  "count": 100,
  "schema": {
    "transaction_id": "string",
    "amount": "decimal",
    "currency": "string",
    "timestamp": "datetime",
    "status": "enum[pending,completed,failed]"
  }
}
```

## Benefits of Agentic Approach

1. **Higher Quality Data**: AI agents understand context and generate more realistic data
2. **Comprehensive Coverage**: Edge case agent ensures all boundary conditions are tested
3. **Flexibility**: Natural language prompts allow for complex requirements
4. **Intelligent Edge Cases**: AI identifies edge cases you might not think of
5. **Maintainability**: Adding new agent capabilities is straightforward

## Code Structure

```
services/
├── llm/
│   ├── service.go          # LLM service interface
│   └── openai.go          # OpenAI implementation
├── agents/
│   ├── agents.go          # Agent implementations
│   └── orchestrator.go    # Multi-agent orchestration
└── generate/
    └── service_impl.go    # Uses agents for execution
```

## Error Handling

The agentic flow includes robust error handling:

- If Synthetic Data Agent fails → Task marked as failed
- If Edge Case Agent fails → Task still succeeds with synthetic data only
- LLM API errors → Detailed error messages returned
- Invalid configurations → Validation before execution

## Performance Considerations

- **Token Usage**: Monitor LLM token consumption based on schema complexity
- **Rate Limits**: OpenAI API has rate limits per minute/day
- **Async Execution**: Consider implementing async execution for large datasets
- **Caching**: Future improvement to cache similar requests

## Future Enhancements

1. **Additional Agents**:
   - Validation Agent (validates generated data)
   - Format Conversion Agent (converts between formats)
   - Privacy Agent (ensures PII compliance)

2. **Enhanced Orchestration**:
   - Parallel agent execution where possible
   - Agent result validation and retry logic
   - Dynamic agent selection based on requirements

3. **Advanced Features**:
   - Streaming responses for large datasets
   - Custom agent configurations per task
   - Human-in-the-loop for agent feedback

## Troubleshooting

### Missing API Key
```
Error: OpenAI API key not configured
Solution: Set OPENAI_API_KEY in your .env file
```

### Rate Limit Errors
```
Error: Rate limit exceeded
Solution: Reduce request frequency or upgrade OpenAI plan
```

### Invalid Schema
```
Error: Failed to parse generated data
Solution: Ensure schema is well-defined and valid JSON
```

## Support

For issues or questions:
1. Check the logs for detailed error messages
2. Verify OpenAI API key and configuration
3. Review the schema structure
4. Check OpenAI service status

## License

This feature is part of the Syntra project and follows the same license terms.
