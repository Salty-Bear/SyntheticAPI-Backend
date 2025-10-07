package generate

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Aryaman/syntra/middlewares"
	"github.com/Aryaman/syntra/providers"
	"github.com/Aryaman/syntra/sdk"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

// Create creates a new generate task
func Create(c *fiber.Ctx) error {
	log.Debug("received create generate task request")

	// Get user ID from context
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.GenerateResponse{
			Success: false,
			Message: "user authentication required",
		})
	}

	payload := new(sdk.Generate)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.GenerateResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}
	log.Debug("parsed create generate task request")

	// Generate ID if not provided
	if payload.ID == "" {
		payload.ID = uuid.New().String()
	}

	// Set user ID from context
	payload.UserId = userID
	payload.CreatedBy = userID
	payload.UpdatedBy = userID

	pr := providers.GetProviders(c)
	createdGenerate, err := pr.S.Generate.Create(c.Context(), payload)
	if err != nil {
		message := fmt.Sprintf("failed to create generate task. %v", err)
		log.Errorw("failed to create generate task", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.GenerateResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("generate task created successfully")

	return c.Status(http.StatusCreated).JSON(sdk.GenerateResponse{
		Success: true,
		Message: "Generate task created successfully",
		Data:    createdGenerate,
	})
}

// Get retrieves a generate task by ID
func Get(c *fiber.Ctx) error {
	log.Debug("received get generate task request")

	// Get user ID from context
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.GenerateResponse{
			Success: false,
			Message: "user authentication required",
		})
	}

	generateID := c.Params("id")
	userId := userID

	pr := providers.GetProviders(c)
	generate, err := pr.S.Generate.Get(c.Context(), generateID, userId)
	if err != nil {
		message := fmt.Sprintf("failed to get generate task. %v", err)
		log.Errorw("failed to get generate task", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.GenerateResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("generate task retrieved successfully")
	return c.JSON(sdk.GenerateResponse{
		Success: true,
		Message: "Generate task retrieved successfully",
		Data:    generate,
	})
}

// List retrieves all generate tasks with optional filtering and pagination
func List(c *fiber.Ctx) error {
	log.Debug("received list generate tasks request")

	// Get user ID from context
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.GenerateResponse{
			Success: false,
			Message: "user authentication required",
		})
	}

	userId := userID

	// Parse query parameters
	query := &sdk.GenerateQuery{
		Status:   c.Query("status"),
		DataType: c.Query("data_type"),
		UserId:   userId,
	}

	// Parse enabled parameter
	enabledStr := c.Query("enabled")
	if enabledStr != "" {
		enabledBool, err := strconv.ParseBool(enabledStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(sdk.GenerateResponse{
				Success: false,
				Message: fmt.Sprintf("invalid enabled parameter. %v", err),
			})
		}
		query.Enabled = &enabledBool
	}

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(sdk.GenerateResponse{
				Success: false,
				Message: fmt.Sprintf("invalid page parameter. %v", err),
			})
		}
		query.Page = page
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(sdk.GenerateResponse{
				Success: false,
				Message: fmt.Sprintf("invalid limit parameter. %v", err),
			})
		}
		query.Limit = limit
		// Calculate offset based on page and limit
		query.Offset = (query.Page - 1) * query.Limit
		if query.Offset < 0 {
			query.Offset = 0
		}
	}

	pr := providers.GetProviders(c)
	generates, err := pr.S.Generate.List(c.Context(), query)
	if err != nil {
		message := fmt.Sprintf("failed to list generate tasks. %v", err)
		log.Errorw("failed to list generate tasks", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.GenerateResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("generate tasks retrieved successfully")
	return c.JSON(sdk.GenerateResponse{
		Success:   true,
		Message:   "Generate tasks retrieved successfully",
		Generates: generates,
	})
}

// Update updates an existing generate task
func Update(c *fiber.Ctx) error {
	log.Debug("received update generate task request")

	// Get user ID from context
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.GenerateResponse{
			Success: false,
			Message: "user authentication required",
		})
	}

	generateID := c.Params("id")

	payload := new(sdk.Generate)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.GenerateResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}
	log.Debug("parsed update generate task request")

	// Set user ID and generate task ID
	payload.UserId = userID
	payload.UpdatedBy = userID
	payload.ID = generateID

	pr := providers.GetProviders(c)
	updatedGenerate, err := pr.S.Generate.Update(c.Context(), payload)
	if err != nil {
		message := fmt.Sprintf("failed to update generate task. %v", err)
		log.Errorw("failed to update generate task", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.GenerateResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("generate task updated successfully")
	return c.JSON(sdk.GenerateResponse{
		Success: true,
		Message: "Generate task updated successfully",
		Data:    updatedGenerate,
	})
}

// Delete deletes a generate task by ID
func Delete(c *fiber.Ctx) error {
	log.Debug("received delete generate task request")

	// Get user ID from context
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.GenerateResponse{
			Success: false,
			Message: "user authentication required",
		})
	}

	generateID := c.Params("id")
	userId := userID

	pr := providers.GetProviders(c)
	err := pr.S.Generate.Delete(c.Context(), generateID, userId)
	if err != nil {
		message := fmt.Sprintf("failed to delete generate task. %v", err)
		log.Errorw("failed to delete generate task", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.GenerateResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("generate task deleted successfully")
	return c.JSON(sdk.GenerateResponse{
		Success: true,
		Message: "Generate task deleted successfully",
	})
}

// Execute executes a generate task to produce synthetic data
func Execute(c *fiber.Ctx) error {
	log.Debug("received execute generate task request")

	// Get user ID from context
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.GenerateResponse{
			Success: false,
			Message: "user authentication required",
		})
	}

	generateID := c.Params("id")

	pr := providers.GetProviders(c)
	executedGenerate, err := pr.S.Generate.Execute(c.Context(), generateID, userID)
	if err != nil {
		message := fmt.Sprintf("failed to execute generate task. %v", err)
		log.Errorw("failed to execute generate task", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.GenerateResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("generate task executed successfully")
	return c.JSON(sdk.GenerateResponse{
		Success: true,
		Message: "Generate task executed successfully",
		Data:    executedGenerate,
	})
}
