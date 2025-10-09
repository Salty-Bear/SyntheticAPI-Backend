package sdk

import "time"

// Generate represents a generate task entity for service layer operations
type Generate struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	DataType    string      `json:"data_type"`   // json, csv, xml, sql, etc.
	Count       int         `json:"count"`       // Number of data records to generate
	Schema      interface{} `json:"schema"`      // JSON schema for data structure
	Format      string      `json:"format"`      // Output format specifications
	Status      string      `json:"status"`      // active, completed, failed, pending
	Enabled     bool        `json:"enabled"`     // Whether the generate task is enabled
	UserId      string      `json:"user_id"`     // ID of the user who owns this generate task
	CreatedAt   time.Time   `json:"created_at"`  // Timestamp when the generate task was created
	UpdatedAt   time.Time   `json:"updated_at"`  // Timestamp when the generate task was last updated
	OutputData  interface{} `json:"output_data"` // Generated data result
}

// CreateGenerateRequest represents the request body for creating a generate task
type CreateGenerateRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	DataType    string      `json:"data_type" validate:"required,oneof=json csv xml sql yaml"`
	Count       int         `json:"count" validate:"required,min=1,max=10000"`
	Schema      interface{} `json:"schema"`
	Format      string      `json:"format"`
	UserId      string      `json:"user_id" validate:"required"`
}

// UpdateGenerateRequest represents the request body for updating a generate task
type UpdateGenerateRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	DataType    string      `json:"data_type"`
	Count       *int        `json:"count"`
	Schema      interface{} `json:"schema"`
	Format      string      `json:"format"`
	Status      string      `json:"status"`
	Enabled     *bool       `json:"enabled"`
	OutputData  interface{} `json:"output_data"`
}

// ExecuteGenerateRequest represents the request body for executing a generate task
type ExecuteGenerateRequest struct {
	UserId string `json:"user_id" validate:"required"`
}

// GenerateResponse represents the response structure for generate operations
type GenerateResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Generates interface{} `json:"generates,omitempty"`
}

// GenerateQuery represents query parameters for filtering generate tasks
type GenerateQuery struct {
	Status   string `json:"status"`
	DataType string `json:"data_type"`
	Enabled  *bool  `json:"enabled"`
	UserId   string `json:"user_id"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

// GenerateErrorResponse represents error response structure
type GenerateErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ExecuteGenerateResponse represents the response structure for execute operations
type ExecuteGenerateResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
