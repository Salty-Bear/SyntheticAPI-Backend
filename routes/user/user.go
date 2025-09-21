package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Aryaman/pub-sub/providers"
	"github.com/Aryaman/pub-sub/sdk"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

// Create creates a new user
func Create(c *fiber.Ctx) error {
	log.Debug("received create user request")
	payload := new(sdk.User)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}
	log.Debug("parsed create user request")

	// Generate ID if not provided
	if payload.ID == "" {
		payload.ID = uuid.New().String()
	}

	pr := providers.GetProviders(c)
	err := pr.S.User.Create(c.Context(), payload)
	if err != nil {
		message := fmt.Sprintf("failed to create user. %v", err)
		log.Errorw("failed to create user", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}
	log.Debug("user created successfully")

	return c.Status(http.StatusCreated).JSON(sdk.UserResponse{
		Success: true,
		Message: "User created successfully",
		Data:    payload,
	})
}

// Get retrieves a user by ID
func Get(c *fiber.Ctx) error {
	log.Debug("received get user request")
	userID := c.Params("id")

	pr := providers.GetProviders(c)
	user, err := pr.S.User.Get(c.Context(), userID)
	if err != nil {
		message := fmt.Sprintf("failed to get user. %v", err)
		log.Errorw("failed to get user", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("user retrieved successfully")
	return c.JSON(sdk.UserResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// List retrieves all users with optional filtering and pagination
func List(c *fiber.Ctx) error {
	log.Debug("received list users request")

	// Parse query parameters
	query := &sdk.UserQuery{
		ProjectID: c.Query("project_id"),
	}

	// Parse enabled parameter
	enabledStr := c.Query("enabled")
	if enabledStr != "" {
		enabledBool, err := strconv.ParseBool(enabledStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
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
			return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
				Success: false,
				Message: fmt.Sprintf("invalid page parameter. %v", err),
			})
		}
		query.Page = page
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
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
	users, err := pr.S.User.List(c.Context(), query)
	if err != nil {
		message := fmt.Sprintf("failed to list users. %v", err)
		log.Errorw("failed to list users", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("users retrieved successfully")
	return c.JSON(sdk.UserResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Users:   users,
	})
}

// Update updates an existing user
func Update(c *fiber.Ctx) error {
	log.Debug("received update user request")
	userID := c.Params("id")

	payload := new(sdk.User)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(sdk.UserResponse{
			Success: false,
			Message: fmt.Sprintf("invalid request. %v", err),
		})
	}
	log.Debug("parsed update user request")

	// Set the user ID from the URL parameter
	payload.ID = userID

	pr := providers.GetProviders(c)
	updatedUser, err := pr.S.User.Update(c.Context(), payload)
	if err != nil {
		message := fmt.Sprintf("failed to update user. %v", err)
		log.Errorw("failed to update user", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("user updated successfully")
	return c.JSON(sdk.UserResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    updatedUser,
	})
}

// Delete deletes a user by ID
func Delete(c *fiber.Ctx) error {
	log.Debug("received delete user request")
	userID := c.Params("id")

	pr := providers.GetProviders(c)
	err := pr.S.User.Delete(c.Context(), userID)
	if err != nil {
		message := fmt.Sprintf("failed to delete user. %v", err)
		log.Errorw("failed to delete user", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}

	log.Debug("user deleted successfully")
	return c.JSON(sdk.UserResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
