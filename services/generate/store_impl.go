package generate

import (
	"context"
	"time"

	"github.com/Aryaman/syntra/db"
	"github.com/Aryaman/syntra/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// StoreImpl implements the GenerateStore interface with MongoDB operations
type StoreImpl struct {
	db db.DB
}

// NewStore creates a new instance of GenerateStore with MongoDB implementation
func NewStore(database db.DB) GenerateStore {
	return &StoreImpl{
		db: database,
	}
}

// CreateGenerate creates a new generate task in MongoDB
func (s *StoreImpl) CreateGenerate(ctx context.Context, generate *models.Generate) error {
	md := models.GetGenerateModel()
	_, err := s.db.InsertOne(ctx, md, generate)
	return err
}

// GetGenerateByID retrieves a generate task by its ID from MongoDB
func (s *StoreImpl) GetGenerateByID(ctx context.Context, generateID string) (*models.Generate, error) {
	md := models.GetGenerateModel()
	filter := bson.M{md.IdKey: generateID}

	result := s.db.FindOne(ctx, md, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}

	var generate models.Generate
	if err := result.Decode(&generate); err != nil {
		return nil, err
	}

	return &generate, nil
}

// GetGenerates retrieves all generate tasks with optional filtering
func (s *StoreImpl) GetGenerates(ctx context.Context, status, dataType string, enabled *bool, userId string) ([]*models.Generate, error) {
	filter := bson.M{}
	md := models.GetGenerateModel()

	if userId != "" {
		filter[md.UserIdKey] = userId
	}

	if status != "" {
		filter[md.StatusKey] = status
	}

	if dataType != "" {
		filter[md.DataTypeKey] = dataType
	}

	if enabled != nil {
		filter[md.EnabledKey] = *enabled
	}

	cursor, err := s.db.Find(ctx, md, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var generates []*models.Generate
	for cursor.Next(ctx) {
		var generate models.Generate
		if err := cursor.Decode(&generate); err != nil {
			return nil, err
		}
		generates = append(generates, &generate)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return generates, nil
}

// GetGeneratesByUser retrieves all generate tasks for a specific user with optional filtering
func (s *StoreImpl) GetGeneratesByUser(ctx context.Context, userId string, status, dataType string, enabled *bool) ([]*models.Generate, error) {
	return s.GetGenerates(ctx, status, dataType, enabled, userId)
}

// UpdateGenerate updates an existing generate task in MongoDB
func (s *StoreImpl) UpdateGenerate(ctx context.Context, generateID string, updates map[string]interface{}) error {
	md := models.GetGenerateModel()
	filter := bson.M{md.IdKey: generateID}

	// Add updated timestamp if not already provided
	if _, exists := updates["updated_at"]; !exists {
		updates["updated_at"] = time.Now()
	}

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

// DeleteGenerate deletes a generate task by its ID from MongoDB
func (s *StoreImpl) DeleteGenerate(ctx context.Context, generateID string) error {
	md := models.GetGenerateModel()
	filter := bson.M{md.IdKey: generateID}

	result, err := s.db.DeleteOne(ctx, md, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// GenerateExists checks if a generate task exists by ID in MongoDB
func (s *StoreImpl) GenerateExists(ctx context.Context, generateID string) (bool, error) {
	md := models.GetGenerateModel()
	filter := bson.M{md.IdKey: generateID}

	result := s.db.FindOne(ctx, md, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, result.Err()
	}

	return true, nil
}

// NameExists checks if a generate task exists by name for a specific user in MongoDB
func (s *StoreImpl) NameExists(ctx context.Context, name string, userId string) (bool, error) {
	md := models.GetGenerateModel()
	filter := bson.M{
		md.NameKey:   name,
		md.UserIdKey: userId,
	}

	result := s.db.FindOne(ctx, md, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, result.Err()
	}

	return true, nil
}
