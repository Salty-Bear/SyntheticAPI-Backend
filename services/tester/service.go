package tester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TestCase represents a single API test case
type TestCase struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Method   string            `json:"method"`
	Path     string            `json:"path"`
	Headers  map[string]string `json:"headers"`
	Body     interface{}       `json:"body,omitempty"`
	Expected ExpectedResult    `json:"expected"`
	Timeout  time.Duration     `json:"timeout"`
}

// ExpectedResult represents expected test results
type ExpectedResult struct {
	StatusCode   int               `json:"status_code"`
	Headers      map[string]string `json:"headers,omitempty"`
	BodyContains []string          `json:"body_contains,omitempty"`
	Schema       interface{}       `json:"schema,omitempty"`
}

// TestResult represents the result of a test case execution
type TestResult struct {
	TestCaseID   string            `json:"test_case_id"`
	Success      bool              `json:"success"`
	StatusCode   int               `json:"status_code"`
	ResponseTime time.Duration     `json:"response_time"`
	Headers      map[string]string `json:"headers"`
	Body         string            `json:"body"`
	Error        string            `json:"error,omitempty"`
	Timestamp    time.Time         `json:"timestamp"`
}

// TestSuite represents a collection of test cases
type TestSuite struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	BaseURL       string            `json:"base_url"`
	TestCases     []TestCase        `json:"test_cases"`
	GlobalHeaders map[string]string `json:"global_headers,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}

// TestExecution represents a complete test execution
type TestExecution struct {
	ID          string          `json:"id"`
	SuiteID     string          `json:"suite_id"`
	TunnelID    string          `json:"tunnel_id"`
	Status      ExecutionStatus `json:"status"`
	Results     []TestResult    `json:"results"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     time.Time       `json:"end_time"`
	Duration    time.Duration   `json:"duration"`
	PassedTests int             `json:"passed_tests"`
	FailedTests int             `json:"failed_tests"`
	TotalTests  int             `json:"total_tests"`
}

// ExecutionStatus represents the status of test execution
type ExecutionStatus string

const (
	StatusPending   ExecutionStatus = "pending"
	StatusRunning   ExecutionStatus = "running"
	StatusCompleted ExecutionStatus = "completed"
	StatusFailed    ExecutionStatus = "failed"
	StatusCancelled ExecutionStatus = "cancelled"
)

// Service defines the interface for the testing service
type Service interface {
	// CreateTestSuite creates a new test suite
	CreateTestSuite(ctx context.Context, suite *TestSuite) error

	// GetTestSuite retrieves a test suite by ID
	GetTestSuite(ctx context.Context, suiteID string) (*TestSuite, error)

	// ExecuteTestSuite runs all tests in a suite against a tunnel
	ExecuteTestSuite(ctx context.Context, suiteID, tunnelURL string) (*TestExecution, error)

	// GetTestExecution retrieves a test execution by ID
	GetTestExecution(ctx context.Context, executionID string) (*TestExecution, error)

	// GenerateTestCases automatically generates test cases for API endpoints
	GenerateTestCases(ctx context.Context, baseURL string) ([]TestCase, error)
}

// ServiceImpl implements the Service interface
type ServiceImpl struct {
	client *http.Client
	store  Store
}

// Store defines the interface for test data storage
type Store interface {
	SaveTestSuite(ctx context.Context, suite *TestSuite) error
	GetTestSuite(ctx context.Context, suiteID string) (*TestSuite, error)
	SaveTestExecution(ctx context.Context, execution *TestExecution) error
	GetTestExecution(ctx context.Context, executionID string) (*TestExecution, error)
	UpdateTestExecution(ctx context.Context, execution *TestExecution) error
}

// NewService creates a new testing service
func NewService(store Store) Service {
	return &ServiceImpl{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		store: store,
	}
}

// CreateTestSuite creates a new test suite
func (s *ServiceImpl) CreateTestSuite(ctx context.Context, suite *TestSuite) error {
	suite.CreatedAt = time.Now()
	return s.store.SaveTestSuite(ctx, suite)
}

// GetTestSuite retrieves a test suite by ID
func (s *ServiceImpl) GetTestSuite(ctx context.Context, suiteID string) (*TestSuite, error) {
	return s.store.GetTestSuite(ctx, suiteID)
}

// ExecuteTestSuite runs all tests in a suite against a tunnel
func (s *ServiceImpl) ExecuteTestSuite(ctx context.Context, suiteID, tunnelURL string) (*TestExecution, error) {
	suite, err := s.store.GetTestSuite(ctx, suiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test suite: %w", err)
	}

	execution := &TestExecution{
		ID:         generateExecutionID(),
		SuiteID:    suiteID,
		TunnelID:   extractTunnelID(tunnelURL),
		Status:     StatusRunning,
		Results:    make([]TestResult, 0, len(suite.TestCases)),
		StartTime:  time.Now(),
		TotalTests: len(suite.TestCases),
	}

	// Save initial execution state
	err = s.store.SaveTestExecution(ctx, execution)
	if err != nil {
		return nil, fmt.Errorf("failed to save test execution: %w", err)
	}

	// Execute tests
	for _, testCase := range suite.TestCases {
		result := s.executeTestCase(ctx, testCase, tunnelURL, suite.GlobalHeaders)
		execution.Results = append(execution.Results, result)

		if result.Success {
			execution.PassedTests++
		} else {
			execution.FailedTests++
		}
	}

	// Update execution status
	execution.EndTime = time.Now()
	execution.Duration = execution.EndTime.Sub(execution.StartTime)
	execution.Status = StatusCompleted

	err = s.store.UpdateTestExecution(ctx, execution)
	if err != nil {
		return nil, fmt.Errorf("failed to update test execution: %w", err)
	}

	return execution, nil
}

// executeTestCase executes a single test case
func (s *ServiceImpl) executeTestCase(ctx context.Context, testCase TestCase, baseURL string, globalHeaders map[string]string) TestResult {
	result := TestResult{
		TestCaseID: testCase.ID,
		Timestamp:  time.Now(),
	}

	// Build request URL
	url := baseURL + testCase.Path

	// Prepare request body
	var bodyReader io.Reader
	if testCase.Body != nil {
		bodyBytes, err := json.Marshal(testCase.Body)
		if err != nil {
			result.Error = fmt.Sprintf("Failed to marshal request body: %v", err)
			return result
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	startTime := time.Now()
	req, err := http.NewRequestWithContext(ctx, testCase.Method, url, bodyReader)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		return result
	}

	// Set headers
	for key, value := range globalHeaders {
		req.Header.Set(key, value)
	}
	for key, value := range testCase.Headers {
		req.Header.Set(key, value)
	}

	// Set content type for JSON body
	if testCase.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply timeout
	if testCase.Timeout > 0 {
		ctx, cancel := context.WithTimeout(ctx, testCase.Timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	// Execute request
	resp, err := s.client.Do(req)
	result.ResponseTime = time.Since(startTime)

	if err != nil {
		result.Error = fmt.Sprintf("Request failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	// Read response
	result.StatusCode = resp.StatusCode
	result.Headers = make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[key] = values[0]
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to read response body: %v", err)
		return result
	}
	result.Body = string(bodyBytes)

	// Validate result
	result.Success = s.validateTestResult(result, testCase.Expected)

	return result
}

// validateTestResult validates the test result against expectations
func (s *ServiceImpl) validateTestResult(result TestResult, expected ExpectedResult) bool {
	// Check status code
	if result.StatusCode != expected.StatusCode {
		return false
	}

	// Check required headers
	for key, expectedValue := range expected.Headers {
		if actualValue, exists := result.Headers[key]; !exists || actualValue != expectedValue {
			return false
		}
	}

	// Check body contains
	for _, expectedContent := range expected.BodyContains {
		if !contains(result.Body, expectedContent) {
			return false
		}
	}

	return true
}

// GetTestExecution retrieves a test execution by ID
func (s *ServiceImpl) GetTestExecution(ctx context.Context, executionID string) (*TestExecution, error) {
	return s.store.GetTestExecution(ctx, executionID)
}

// GenerateTestCases automatically generates test cases for API endpoints
func (s *ServiceImpl) GenerateTestCases(ctx context.Context, baseURL string) ([]TestCase, error) {
	// This is a simplified implementation
	// In practice, you might want to:
	// 1. Discover endpoints by crawling OpenAPI specs
	// 2. Generate synthetic data based on schemas
	// 3. Create edge cases and error scenarios

	commonTestCases := []TestCase{
		{
			ID:     "health-check",
			Name:   "Health Check",
			Method: "GET",
			Path:   "/health",
			Expected: ExpectedResult{
				StatusCode: 200,
			},
			Timeout: 5 * time.Second,
		},
		{
			ID:     "not-found",
			Name:   "404 Not Found",
			Method: "GET",
			Path:   "/nonexistent-endpoint",
			Expected: ExpectedResult{
				StatusCode: 404,
			},
			Timeout: 5 * time.Second,
		},
	}

	return commonTestCases, nil
}

// Helper functions

func generateExecutionID() string {
	return fmt.Sprintf("exec_%d", time.Now().UnixNano())
}

func extractTunnelID(tunnelURL string) string {
	// Extract tunnel ID from URL (simplified)
	// In practice, you might parse the subdomain or use other methods
	return "tunnel_id_from_url"
}

func contains(text, substring string) bool {
	return len(text) > 0 && len(substring) > 0 &&
		(text == substring ||
			(len(text) > len(substring) &&
				findSubstring(text, substring)))
}

func findSubstring(text, substring string) bool {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}
