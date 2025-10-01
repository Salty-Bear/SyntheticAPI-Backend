package models

import "time"

// Generate represents a data generation task entity in the system.
// Generate tasks allow users to create synthetic data with specified schemas and formats.
type Generate struct {
	Id          string      `bson:"id"`          // Unique identifier for the generate task (UUID)
	Name        string      `bson:"name"`        // Display name of the generate task
	Description string      `bson:"description"` // Description of the generate task purpose
	DataType    string      `bson:"data_type"`   // Type of data to generate (json, csv, xml, sql, etc.)
	Count       int         `bson:"count"`       // Number of data records to generate
	Schema      interface{} `bson:"schema"`      // JSON schema for data structure
	Format      string      `bson:"format"`      // Output format specifications
	Status      string      `bson:"status"`      // Current status (active, completed, failed, pending)
	Enabled     bool        `bson:"enabled"`     // Whether the generate task is enabled
	UserId      string      `bson:"user_id"`     // ID of the user who owns this generate task
	CreatedBy   string      `bson:"created_by"`  // ID of the user who created this generate task
	UpdatedBy   string      `bson:"updated_by"`  // ID of the user who last updated this generate task
	CreatedAt   time.Time   `bson:"created_at"`  // Timestamp when the generate task was created
	UpdatedAt   time.Time   `bson:"updated_at"`  // Timestamp when the generate task was last updated
	OutputData  interface{} `bson:"output_data"` // Generated data result
}

// GenerateModel provides database access patterns and field mappings for Generate entities.
// It provides database collection name and field key mappings for MongoDB operations.
type GenerateModel struct {
	IdKey          string // BSON field key for generate task ID
	NameKey        string // BSON field key for generate task name
	DescriptionKey string // BSON field key for generate task description
	DataTypeKey    string // BSON field key for data type
	CountKey       string // BSON field key for count
	SchemaKey      string // BSON field key for schema
	FormatKey      string // BSON field key for format
	StatusKey      string // BSON field key for status
	EnabledKey     string // BSON field key for enabled status
	UserIdKey      string // BSON field key for user ID
	CreatedByKey   string // BSON field key for created by
	UpdatedByKey   string // BSON field key for updated by
	CreatedAtKey   string // BSON field key for created at
	UpdatedAtKey   string // BSON field key for updated at
	OutputDataKey  string // BSON field key for output data
}

// Name returns the MongoDB collection name for generate tasks.
// This implements the DbCollection interface.
func (g GenerateModel) Name() string {
	return "generates"
}

// DbName returns the database name for the application.
// This implements the DbCollection interface.
func (g GenerateModel) DbName() string {
	return "syntra"
}

// GetGenerateModel returns a properly initialized GenerateModel with all field mappings.
// This function provides a singleton pattern for accessing generate model operations.
//
// Returns a GenerateModel instance with all BSON field keys mapped to their respective field names.
func GetGenerateModel() GenerateModel {
	return GenerateModel{
		IdKey:          "id",
		NameKey:        "name",
		DescriptionKey: "description",
		DataTypeKey:    "data_type",
		CountKey:       "count",
		SchemaKey:      "schema",
		FormatKey:      "format",
		StatusKey:      "status",
		EnabledKey:     "enabled",
		UserIdKey:      "user_id",
		CreatedByKey:   "created_by",
		UpdatedByKey:   "updated_by",
		CreatedAtKey:   "created_at",
		UpdatedAtKey:   "updated_at",
		OutputDataKey:  "output_data",
	}
}
