package tester

import (
	"time"

	"github.com/Aryaman/pub-sub/services/tester"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for testing operations
type Handler struct {
	service tester.Service
}

// NewHandler creates a new handler instance
func NewHandler(service tester.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers all testing routes
func RegisterRoutes(app fiber.Router, handler *Handler) {
	app.Post("/suites", handler.CreateTestSuite)
	app.Get("/suites/:suiteId", handler.GetTestSuite)
	app.Post("/suites/:suiteId/execute", handler.ExecuteTestSuite)
	app.Get("/executions/:executionId", handler.GetTestExecution)
	app.Post("/generate", handler.GenerateTestCases)
}

// CreateTestSuiteRequest represents the request to create a test suite
type CreateTestSuiteRequest struct {
	Name          string            `json:"name" validate:"required"`
	BaseURL       string            `json:"base_url" validate:"required,url"`
	TestCases     []tester.TestCase `json:"test_cases" validate:"required"`
	GlobalHeaders map[string]string `json:"global_headers,omitempty"`
}

// ExecuteTestSuiteRequest represents the request to execute a test suite
type ExecuteTestSuiteRequest struct {
	TunnelURL string `json:"tunnel_url" validate:"required,url"`
}

// GenerateTestCasesRequest represents the request to generate test cases
type GenerateTestCasesRequest struct {
	BaseURL string `json:"base_url" validate:"required,url"`
}

// CreateTestSuite creates a new test suite
func (h *Handler) CreateTestSuite(c *fiber.Ctx) error {
	var req CreateTestSuiteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	suite := &tester.TestSuite{
		ID:            uuid.New().String(),
		Name:          req.Name,
		BaseURL:       req.BaseURL,
		TestCases:     req.TestCases,
		GlobalHeaders: req.GlobalHeaders,
		CreatedAt:     time.Now(),
	}

	// Set IDs for test cases if not provided
	for i := range suite.TestCases {
		if suite.TestCases[i].ID == "" {
			suite.TestCases[i].ID = uuid.New().String()
		}
		// Set default timeout if not provided
		if suite.TestCases[i].Timeout == 0 {
			suite.TestCases[i].Timeout = 30 * time.Second
		}
	}

	err := h.service.CreateTestSuite(c.Context(), suite)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create test suite",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    suite,
		"message": "Test suite created successfully",
	})
}

// GetTestSuite retrieves a test suite by ID
func (h *Handler) GetTestSuite(c *fiber.Ctx) error {
	suiteID := c.Params("suiteId")
	if suiteID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Suite ID is required",
		})
	}

	suite, err := h.service.GetTestSuite(c.Context(), suiteID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Test suite not found",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    suite,
	})
}

// ExecuteTestSuite executes a test suite against a tunnel
func (h *Handler) ExecuteTestSuite(c *fiber.Ctx) error {
	suiteID := c.Params("suiteId")
	if suiteID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Suite ID is required",
		})
	}

	var req ExecuteTestSuiteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Execute test suite asynchronously
	execution, err := h.service.ExecuteTestSuite(c.Context(), suiteID, req.TunnelURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to execute test suite",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    execution,
		"message": "Test execution started",
	})
}

// GetTestExecution retrieves a test execution by ID
func (h *Handler) GetTestExecution(c *fiber.Ctx) error {
	executionID := c.Params("executionId")
	if executionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Execution ID is required",
		})
	}

	execution, err := h.service.GetTestExecution(c.Context(), executionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Test execution not found",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    execution,
	})
}

// GenerateTestCases generates test cases for an API
func (h *Handler) GenerateTestCases(c *fiber.Ctx) error {
	var req GenerateTestCasesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	testCases, err := h.service.GenerateTestCases(c.Context(), req.BaseURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate test cases",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    testCases,
		"count":   len(testCases),
		"message": "Test cases generated successfully",
	})
}
