package models

// Tunnel represents a tunnel entity in the system.
// Tunnels provide secure network connections and endpoint management.
type Tunnel struct {
	Id          string `bson:"id"`          // Unique identifier for the tunnel (UUID)
	ProjectId   string `bson:"project_id"`  // Project ID this tunnel belongs to
	Name        string `bson:"name"`        // Display name of the tunnel
	Description string `bson:"description"` // Description of the tunnel purpose
	Endpoint    string `bson:"endpoint"`    // Target endpoint for the tunnel
	Port        int    `bson:"port"`        // Port number for the tunnel
	Protocol    string `bson:"protocol"`    // Protocol type (http, https, tcp, udp)
	Status      string `bson:"status"`      // Current status (active, inactive, pending)
	Enabled     bool   `bson:"enabled"`     // Whether the tunnel is enabled
}

// TunnelModel provides database access patterns and field mappings for Tunnel entities.
// It provides database collection name and field key mappings for MongoDB operations.
type TunnelModel struct {
	IdKey          string // BSON field key for tunnel ID
	ProjectIdKey   string // BSON field key for project ID
	NameKey        string // BSON field key for tunnel name
	DescriptionKey string // BSON field key for tunnel description
	EndpointKey    string // BSON field key for tunnel endpoint
	PortKey        string // BSON field key for tunnel port
	ProtocolKey    string // BSON field key for tunnel protocol
	StatusKey      string // BSON field key for tunnel status
	EnabledKey     string // BSON field key for enabled status
}

// Name returns the MongoDB collection name for tunnels.
// This implements the DbCollection interface.
func (t TunnelModel) Name() string {
	return "tunnels"
}

// DbName returns the database name for the application.
// This implements the DbCollection interface.
func (t TunnelModel) DbName() string {
	return "syntra"
}

// GetTunnelModel returns a properly initialized TunnelModel with all field mappings.
// This function provides a singleton pattern for accessing tunnel model operations.
//
// Returns a TunnelModel instance with all BSON field keys mapped to their respective field names.
func GetTunnelModel() TunnelModel {
	return TunnelModel{
		IdKey:          "id",
		ProjectIdKey:   "project_id",
		NameKey:        "name",
		DescriptionKey: "description",
		EndpointKey:    "endpoint",
		PortKey:        "port",
		ProtocolKey:    "protocol",
		StatusKey:      "status",
		EnabledKey:     "enabled",
	}
}
