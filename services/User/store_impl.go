package user

import (
	"context"
	"time"

	"github.com/Aryaman/syntra/db"
	"github.com/Aryaman/syntra/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// StoreImpl implements the UserStore interface with MongoDB operations
type StoreImpl struct {
	db db.DB
}

// NewStore creates a new instance of UserStore with MongoDB implementation
func NewStore(database db.DB) UserStore {
	return &StoreImpl{
		db: database,
	}
}

// CreateUser creates a new user in MongoDB
func (s *StoreImpl) CreateUser(ctx context.Context, user *models.User) error {
	md := models.GetUserModel()
	_, err := s.db.InsertOne(ctx, md, user)
	return err
}

// GetUserByID retrieves a user by their ID from MongoDB
func (s *StoreImpl) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	md := models.GetUserModel()
	filter := bson.M{md.IdKey: userID}

	result := s.db.FindOne(ctx, md, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}

	var user models.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email address from MongoDB
func (s *StoreImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	md := models.GetUserModel()
	filter := bson.M{md.EmailKey: email}

	result := s.db.FindOne(ctx, md, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}

	var user models.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUsers retrieves all users with optional filtering
func (s *StoreImpl) GetUsers(ctx context.Context, projectID string, enabled *bool) ([]*models.User, error) {
	filter := bson.M{}
	md := models.GetUserModel()

	// Note: ProjectID field has been removed from the User model
	// If project-based filtering is needed, it should be implemented differently

	if enabled != nil {
		filter[md.EnabledKey] = *enabled
	}

	cursor, err := s.db.Find(ctx, md, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user in MongoDB
func (s *StoreImpl) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	md := models.GetUserModel()
	filter := bson.M{md.IdKey: userID}

	// Add updated timestamp
	updates["updated_at"] = time.Now()

	update := bson.M{"$set": updates}

	result, err := s.db.UpdateOne(ctx, md, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// DeleteUser deletes a user by their ID from MongoDB
func (s *StoreImpl) DeleteUser(ctx context.Context, userID string) error {
	md := models.GetUserModel()
	filter := bson.M{md.IdKey: userID}

	result, err := s.db.DeleteOne(ctx, md, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// UserExists checks if a user exists by ID in MongoDB
func (s *StoreImpl) UserExists(ctx context.Context, userID string) (bool, error) {
	md := models.GetUserModel()
	filter := bson.M{md.IdKey: userID}

	result := s.db.FindOne(ctx, md, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, result.Err()
	}

	return true, nil
}

// EmailExists checks if a user exists by email in MongoDB
func (s *StoreImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	md := models.GetUserModel()
	var usr models.User
	filter := bson.M{md.EmailKey: email}
	err := s.db.FindOne(ctx, md, filter).Decode(&usr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
