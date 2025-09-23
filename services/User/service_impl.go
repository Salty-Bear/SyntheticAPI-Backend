package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Aryaman/syntra/sdk"
	"github.com/google/uuid"
)

// ServiceImpl implements the UserService interface
type ServiceImpl struct {
	store UserStore
}

// NewService creates a new instance of UserService
func NewService(store UserStore) UserService {
	return &ServiceImpl{
		store: store,
	}
}

// Create creates a new user or returns existing user if email already exists
func (s *ServiceImpl) Create(ctx context.Context, user *sdk.User) (*sdk.User, error) {
	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Check if user with email already exists
	exists, err := s.store.EmailExists(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		// Get existing user by email and return it
		existingUser, err := s.store.GetUserByEmail(ctx, user.Email)
		if err != nil {
			return nil, fmt.Errorf("error retrieving existing user: %w", err)
		}

		// Convert database model to SDK type
		return fromModelToSdk(existingUser), nil
	}

	// Create database model from SDK user
	dbUser := fromSdkToModel(*user)

	// Save user to database
	if err := s.store.CreateUser(ctx, &dbUser); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Convert back to SDK type to return
	return fromModelToSdk(&dbUser), nil
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
	return fromModelToSdk(user), nil
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
		sdkUsers = append(sdkUsers, fromModelToSdk(user))
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
	return fromModelToSdk(updatedUser), nil
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

	return &sdk.FirebaseUser{
		UID:   tokenInfo.Sub,
		Email: tokenInfo.Email,
		Name:  tokenInfo.Name,
	}, nil
}
