package tunnel

import (
	"github.com/Aryaman/syntra/db/models"
	"github.com/Aryaman/syntra/sdk"
)

// fromSdkToModel converts an SDK Tunnel to a database model Tunnel
func fromSdkToModel(tunnel sdk.Tunnel) models.Tunnel {
	return models.Tunnel{
		Id:          tunnel.ID,
		Name:        tunnel.Name,
		Description: tunnel.Description,
		Endpoint:    tunnel.Endpoint,
		Port:        tunnel.Port,
		Protocol:    tunnel.Protocol,
		Status:      tunnel.Status,
		Enabled:     tunnel.Enabled,
	}
}

// fromModelToSdk converts a database model Tunnel to an SDK Tunnel
func fromModelToSdk(tunnel *models.Tunnel) *sdk.Tunnel {
	if tunnel == nil {
		return nil
	}

	return &sdk.Tunnel{
		ID:          tunnel.Id,
		Name:        tunnel.Name,
		Description: tunnel.Description,
		Endpoint:    tunnel.Endpoint,
		Port:        tunnel.Port,
		Protocol:    tunnel.Protocol,
		Status:      tunnel.Status,
		Enabled:     tunnel.Enabled,
		// Note: CreatedBy, UpdatedBy are not stored in the current models.Tunnel struct
		// but are part of the SDK Tunnel struct. These would need to be added to models.Tunnel
		// if they are required for persistence.
	}
}
