package user

import (
	"github.com/Aryaman/syntra/db/models"
	"github.com/Aryaman/syntra/sdk"
)

// fromSdkToModel converts an SDK User to a database model User
func fromSdkToModel(user sdk.User) models.User {
	return models.User{
		Id:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		Enabled:    user.Enabled,
		ProfilePic: user.ProfilePic,
	}
}

// fromModelToSdk converts a database model User to an SDK User
func fromModelToSdk(user *models.User) *sdk.User {
	if user == nil {
		return nil
	}

	return &sdk.User{
		ID:         user.Id,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		Enabled:    user.Enabled,
		ProfilePic: user.ProfilePic,
		// Note: CreatedBy, UpdatedBy are not stored in the current models.User struct
		// but are part of the SDK User struct. These would need to be added to models.User
		// if they are required for persistence.
	}
}
