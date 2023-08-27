package userRepository

import (
	"ResiSync/pkg/api"
	"ResiSync/pkg/models"
	userService "ResiSync/user/internal/services/user_service.go"

	"go.uber.org/zap"
)

func GetUser(requestContext models.ResiSyncRequestContext, user *userService.User) error {
	db := api.ApplicationContext.Postgres

	log := requestContext.Log
	err := db.Find(&user).Error
	if err != nil {
		log.Error("Error while getting user", zap.Any("user", user), zap.Error(err))
		return err
	}

	return nil
}
