package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Aryaman/syntra/db/models"
	"github.com/Aryaman/syntra/sdk"
)

// ServiceImpl implements the UserService interface
type ServiceImpl struct {
	store     UserStore
	projectID string
}

// NewService creates a new instance of UserService
func NewService(store UserStore, projectID string) UserService {
	return &ServiceImpl{
		store:     store,
		projectID: projectID,
	}
}

// Create creates a new user
func (s *ServiceImpl) Create(ctx context.Context, user *sdk.User) error {
	// Check if user with email already exists
	exists, err := s.store.EmailExists(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user with this email already exists")
	}

	// Create database model from SDK user
	dbUser := &models.User{
		Id:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		Enabled:    user.Enabled,
		ProfilePic: user.ProfilePic,
	}

	// Save user to database
	if err := s.store.CreateUser(ctx, dbUser); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// Get retrieves a user by ID
func (s *ServiceImpl) Get(ctx context.Context, id string) (*sdk.User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Convert database model to SDK type
	sdkUser := &sdk.User{
		ID:         user.Id,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		Enabled:    user.Enabled,
		ProfilePic: user.ProfilePic,
	}

	return sdkUser, nil
}

// List retrieves all users with optional filtering and pagination
func (s *ServiceImpl) List(ctx context.Context, query *sdk.UserQuery) ([]*sdk.User, error) {
	var enabled *bool
	if query != nil {
		enabled = query.Enabled
	}

	projectID := ""
	if query != nil {
		projectID = query.ProjectID
	}

	users, err := s.store.GetUsers(ctx, projectID, enabled)
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	// Apply pagination if specified
	if query != nil && query.Limit > 0 {
		start := query.Offset
		end := start + query.Limit

		if start >= len(users) {
			return []*sdk.User{}, nil
		}

		if end > len(users) {
			end = len(users)
		}

		users = users[start:end]
	}

	// Convert database models to SDK types
	var sdkUsers []*sdk.User
	for _, user := range users {
		sdkUser := &sdk.User{
			ID:         user.Id,
			Name:       user.Name,
			Email:      user.Email,
			Phone:      user.Phone,
			Enabled:    user.Enabled,
			ProfilePic: user.ProfilePic,
		}
		sdkUsers = append(sdkUsers, sdkUser)
	}

	return sdkUsers, nil
}

// Update updates an existing user
func (s *ServiceImpl) Update(ctx context.Context, user *sdk.User) (*sdk.User, error) {
	if user.ID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Check if user exists
	exists, err := s.store.UserExists(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Build updates map
	updates := make(map[string]interface{})
	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Phone != "" {
		updates["phone"] = user.Phone
	}
	updates["enabled"] = user.Enabled
	if user.ProfilePic != "" {
		updates["profile_pic"] = user.ProfilePic
	}

	// Update user
	if err := s.store.UpdateUser(ctx, user.ID, updates); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	// Retrieve and return updated user
	updatedUser, err := s.store.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving updated user: %w", err)
	}

	// Convert to SDK type
	sdkUser := &sdk.User{
		ID:         updatedUser.Id,
		Name:       updatedUser.Name,
		Email:      updatedUser.Email,
		Phone:      updatedUser.Phone,
		Enabled:    updatedUser.Enabled,
		ProfilePic: updatedUser.ProfilePic,
	}

	return sdkUser, nil
}

// Delete deletes a user by ID
func (s *ServiceImpl) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user ID is required")
	}

	// Check if user exists
	exists, err := s.store.UserExists(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Delete user
	if err := s.store.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

// VerifyFirebaseToken validates Firebase authentication token
// This is a simplified implementation using Google's tokeninfo endpoint
// For production, consider using Firebase Admin SDK
func (s *ServiceImpl) VerifyFirebaseToken(ctx context.Context, token string) (*sdk.FirebaseUser, error) {
	// Use Google's tokeninfo endpoint to verify the ID token
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token: status %d", resp.StatusCode)
	}

	var tokenInfo struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Aud   string `json:"aud"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %w", err)
	}

	// Verify audience (project ID)
	if s.projectID != "" && tokenInfo.Aud != s.projectID {
		return nil, fmt.Errorf("invalid audience: expected %s, got %s", s.projectID, tokenInfo.Aud)
	}

	return &sdk.FirebaseUser{
		UID:   tokenInfo.Sub,
		Email: tokenInfo.Email,
		Name:  tokenInfo.Name,
	}, nil
}
