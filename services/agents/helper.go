package agents

// SystemPrompts contains all system prompts for agents
var SystemPrompts = struct {
	SyntheticDataAgent string
	EdgeCaseAgent      string
}{
	SyntheticDataAgent: `You are an expert synthetic data generator. Your task is to generate realistic, high-quality synthetic data that matches the provided schema and requirements.

Rules:
1. Generate data that strictly adheres to the provided schema structure
2. Ensure data is realistic and diverse
3. Follow any specific instructions in the user prompt
4. Return ONLY valid JSON without any additional text or explanations
5. Generate exactly the number of records requested
6. Use realistic values appropriate to the field names and types
7. Ensure data relationships are logical and consistent
8. Vary the data to represent realistic distribution and patterns`,

	EdgeCaseAgent: `You are an expert at identifying and generating edge case scenarios. Your task is to analyze the provided schema and generate comprehensive edge case test data that covers all boundary conditions, unusual inputs, and potential error scenarios.

Rules:
1. Identify ALL possible edge cases for the given schema
2. Generate data that tests boundaries (min/max values, empty strings, nulls, special characters, etc.)
3. Include corner cases that might break typical validations
4. Cover different data type edge cases (very long strings, negative numbers, special dates, etc.)
5. Return ONLY valid JSON without any additional text or explanations
6. Structure your response to include both the edge case scenarios and the generated data
7. Consider security implications (SQL injection patterns, XSS attempts, etc.)
8. Test internationalization edge cases (unicode, RTL text, etc.)`,
}

// UserPromptTemplates contains templates for user prompts
var UserPromptTemplates = struct {
	SyntheticData string
	EdgeCase      string
}{
	SyntheticData: `Generate %d synthetic data records with the following requirements:

Prompt: %s

Schema: %s

Data Type: %s
Format: %s

Return the data as a JSON array of objects. Each object should match the schema structure.
Ensure diversity in the generated data while maintaining realism.`,

	EdgeCase: `Analyze the schema and generate comprehensive edge case test data:

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
- Boundary values (min/max for numbers, length limits for strings)
- Empty or null values (empty strings, null fields, missing required fields)
- Special characters and encoding (unicode, emojis, control characters, newlines)
- Extremely long/short values (1-character strings, 10000+ character strings)
- Invalid but parseable data (malformed emails, unusual phone numbers)
- Unusual but valid combinations (leap year dates, timezone edge cases)
- Security edge cases (SQL injection patterns, XSS attempts, path traversal)
- Numeric edge cases (zero, negative, very large numbers, decimals with many places)
- Type confusion (string representations of numbers, booleans as strings)`,
}

// ErrorMessages contains standard error messages
var ErrorMessages = struct {
	NoMessages         string
	GenerationFailed   string
	ParseFailed        string
	EdgeCaseGenFailed  string
	EdgeCaseParseFailed string
}{
	NoMessages:         "no messages provided",
	GenerationFailed:   "failed to generate synthetic data",
	ParseFailed:        "failed to parse generated data",
	EdgeCaseGenFailed:  "failed to generate edge cases",
	EdgeCaseParseFailed: "failed to parse edge case data",
}

// LogMessages contains standard log messages
var LogMessages = struct {
	AgentStarting        string
	AgentCompleted       string
	LLMCallFailed        string
	ParseResponseFailed  string
}{
	AgentStarting:       "[%s] Starting execution for %d records",
	AgentCompleted:      "[%s] Successfully generated %s",
	LLMCallFailed:       "[%s] LLM call failed: %v",
	ParseResponseFailed: "[%s] Failed to parse response: %v",
}
