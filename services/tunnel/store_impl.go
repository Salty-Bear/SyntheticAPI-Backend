package tunnel

import (
	"context"
	"time"

	"github.com/Aryaman/syntra/db"
	"github.com/Aryaman/syntra/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// StoreImpl implements the TunnelStore interface with MongoDB operations
type StoreImpl struct {
	db db.DB
}

// NewStore creates a new instance of TunnelStore with MongoDB implementation
func NewStore(database db.DB) TunnelStore {
	return &StoreImpl{
		db: database,
	}
}

// CreateTunnel creates a new tunnel in MongoDB
func (s *StoreImpl) CreateTunnel(ctx context.Context, tunnel *models.Tunnel) error {
	md := models.GetTunnelModel()
	_, err := s.db.InsertOne(ctx, md, tunnel)
	return err
}

// GetTunnelByID retrieves a tunnel by its ID from MongoDB
func (s *StoreImpl) GetTunnelByID(ctx context.Context, tunnelID string) (*models.Tunnel, error) {
	md := models.GetTunnelModel()
	filter := bson.M{md.IdKey: tunnelID}

	result := s.db.FindOne(ctx, md, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Err()
	}

	var tunnel models.Tunnel
	if err := result.Decode(&tunnel); err != nil {
		return nil, err
	}

	return &tunnel, nil
}

// GetTunnels retrieves all tunnels with optional filtering
func (s *StoreImpl) GetTunnels(ctx context.Context, projectID, status, protocol string, enabled *bool) ([]*models.Tunnel, error) {
	filter := bson.M{}
	md := models.GetTunnelModel()

	if projectID != "" {
		filter[md.ProjectIdKey] = projectID
	}

	if status != "" {
		filter[md.StatusKey] = status
	}

	if protocol != "" {
		filter[md.ProtocolKey] = protocol
	}

	if enabled != nil {
		filter[md.EnabledKey] = *enabled
	}

	cursor, err := s.db.Find(ctx, md, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tunnels []*models.Tunnel
	for cursor.Next(ctx) {
		var tunnel models.Tunnel
		if err := cursor.Decode(&tunnel); err != nil {
			return nil, err
		}
		tunnels = append(tunnels, &tunnel)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tunnels, nil
}

// UpdateTunnel updates an existing tunnel in MongoDB
func (s *StoreImpl) UpdateTunnel(ctx context.Context, tunnelID string, updates map[string]interface{}) error {
	md := models.GetTunnelModel()
	filter := bson.M{md.IdKey: tunnelID}

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

// DeleteTunnel deletes a tunnel by its ID from MongoDB
func (s *StoreImpl) DeleteTunnel(ctx context.Context, tunnelID string) error {
	md := models.GetTunnelModel()
	filter := bson.M{md.IdKey: tunnelID}

	result, err := s.db.DeleteOne(ctx, md, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// TunnelExists checks if a tunnel exists by ID in MongoDB
func (s *StoreImpl) TunnelExists(ctx context.Context, tunnelID string) (bool, error) {
	md := models.GetTunnelModel()
	filter := bson.M{md.IdKey: tunnelID}

	result := s.db.FindOne(ctx, md, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, result.Err()
	}

	return true, nil
}

// NameExists checks if a tunnel exists by name within a project in MongoDB
func (s *StoreImpl) NameExists(ctx context.Context, projectID, name string) (bool, error) {
	md := models.GetTunnelModel()
	filter := bson.M{
		md.ProjectIdKey: projectID,
		md.NameKey:      name,
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
