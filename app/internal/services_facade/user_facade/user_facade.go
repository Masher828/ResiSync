package user_facade

import (
	user_models "ResiSync/app/internal/models"
	"ResiSync/app/internal/services/user_service.go"
	"ResiSync/pkg/api"
	"ResiSync/pkg/models"

	"go.uber.org/zap"
)

func UpdateUserProfile(requestContext models.ResiSyncRequestContext, user *user_models.Resident) error {
	span := api.AddTrace(&requestContext, "info", "UpdateUserProfile")
	defer span.End()
	log := requestContext.Log

	if err := user_service.ValidateUser(requestContext, user); err != nil {
		return err
	}

	user.Id = requestContext.GetUserContext().ID

	err := user_service.UpdateUserProfile(requestContext, user)
	if err != nil {
		log.Error("error while updating user profile", zap.Error(err), zap.Int64("user id", user.Id))
		return err
	}

	return nil
}
