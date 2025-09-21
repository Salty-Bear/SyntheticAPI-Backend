package sdk

// FirebaseUser represents the decoded Firebase token information
type FirebaseUser struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone"`
	ProjectId  string `json:"project_id" validate:"required"`
	ProfilePic string `json:"profile_pic"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Enabled    *bool  `json:"enabled"`
	ProfilePic string `json:"profile_pic"`
}

// UserResponse represents the response structure for user operations
type UserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Users   interface{} `json:"users,omitempty"`
}

// User represents a user entity for service layer operations
type User struct {
	ID         string `json:"id"`
	ProjectId  string `json:"project_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Enabled    bool   `json:"enabled"`
	ProfilePic string `json:"profile_pic"`
	CreatedBy  string `json:"created_by"`
	UpdatedBy  string `json:"updated_by"`
}

// UserQuery represents query parameters for filtering users
type UserQuery struct {
	ProjectID string `json:"project_id"`
	Enabled   *bool  `json:"enabled"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

// ErrorResponse represents error response structure
type UserErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}