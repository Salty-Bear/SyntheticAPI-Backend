package sdk

// Tunnel represents a tunnel entity for service layer operations
type Tunnel struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"` // http, https, tcp, udp
	Status      string `json:"status"`   // active, inactive, pending
	Enabled     bool   `json:"enabled"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
}

// CreateTunnelRequest represents the request body for creating a tunnel
type CreateTunnelRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint" validate:"required"`
	Port        int    `json:"port" validate:"required,min=1,max=65535"`
	Protocol    string `json:"protocol" validate:"required,oneof=http https tcp udp"`
	ProjectID   string `json:"project_id" validate:"required"`
}

// UpdateTunnelRequest represents the request body for updating a tunnel
type UpdateTunnelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
	Port        *int   `json:"port"`
	Protocol    string `json:"protocol"`
	Status      string `json:"status"`
	Enabled     *bool  `json:"enabled"`
}

// TunnelResponse represents the response structure for tunnel operations
type TunnelResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Tunnels interface{} `json:"tunnels,omitempty"`
}

// TunnelQuery represents query parameters for filtering tunnels
type TunnelQuery struct {
	ProjectID string `json:"project_id"`
	Status    string `json:"status"`
	Protocol  string `json:"protocol"`
	Enabled   *bool  `json:"enabled"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

// TunnelErrorResponse represents error response structure
type TunnelErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
