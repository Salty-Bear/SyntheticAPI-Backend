package generate

import (
	"github.com/Aryaman/syntra/db/models"
	"github.com/Aryaman/syntra/sdk"
)

// fromSdkToModel converts an SDK Generate to a database model Generate
func fromSdkToModel(generate sdk.Generate) models.Generate {
	return models.Generate{
		Id:          generate.ID,
		Name:        generate.Name,
		Description: generate.Description,
		DataType:    generate.DataType,
		Count:       generate.Count,
		Schema:      generate.Schema,
		Format:      generate.Format,
		Status:      generate.Status,
		Enabled:     generate.Enabled,
		UserId:      generate.UserId,
		CreatedAt:   generate.CreatedAt,
		UpdatedAt:   generate.UpdatedAt,
		OutputData:  generate.OutputData,
	}
}

// fromModelToSdk converts a database model Generate to an SDK Generate
func fromModelToSdk(generate *models.Generate) *sdk.Generate {
	if generate == nil {
		return nil
	}

	return &sdk.Generate{
		ID:          generate.Id,
		Name:        generate.Name,
		Description: generate.Description,
		DataType:    generate.DataType,
		Count:       generate.Count,
		Schema:      generate.Schema,
		Format:      generate.Format,
		Status:      generate.Status,
		Enabled:     generate.Enabled,
		UserId:      generate.UserId,
		CreatedBy:   generate.CreatedBy,
		UpdatedBy:   generate.UpdatedBy,
		CreatedAt:   generate.CreatedAt,
		UpdatedAt:   generate.UpdatedAt,
		OutputData:  generate.OutputData,
	}
}
