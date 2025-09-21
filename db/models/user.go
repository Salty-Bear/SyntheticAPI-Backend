package models

// User represents a user entity in the system.
// Simplified user model with essential fields for authentication and basic user management.
type User struct {
	Id         string `bson:"id"`          // Unique identifier for the user (UUID)
	Name       string `bson:"name"`        // Display name of the user
	Email      string `bson:"email"`       // Email address of the user
	Phone      string `bson:"phone"`       // Phone number of the user
	Enabled    bool   `bson:"enabled"`     // Whether the user account is active
	ProfilePic string `bson:"profile_pic"` // URL or path to the user's profile picture
}

// UserModel provides database access patterns and field mappings for User entities.
// It provides database collection name and field key mappings for MongoDB operations.
type UserModel struct {
	IdKey         string // BSON field key for user ID
	NameKey       string // BSON field key for user name
	EmailKey      string // BSON field key for user email
	PhoneKey      string // BSON field key for user phone
	EnabledKey    string // BSON field key for enabled status
	ProfilePicKey string // BSON field key for profile picture
}

// Name returns the MongoDB collection name for users.
// This implements the DbCollection interface.
func (u UserModel) Name() string {
	return "users"
}

// DbName returns the database name for the application.
// This implements the DbCollection interface.
func (u UserModel) DbName() string {
	return "syntra"
}

// GetUserModel returns a properly initialized UserModel with all field mappings.
// This function provides a singleton pattern for accessing user model operations.
//
// Returns a UserModel instance with all BSON field keys mapped to their respective field names.
func GetUserModel() UserModel {
	return UserModel{
		IdKey:         "id",
		NameKey:       "name",
		EmailKey:      "email",
		PhoneKey:      "phone",
		EnabledKey:    "enabled",
		ProfilePicKey: "profile_pic",
	}
}