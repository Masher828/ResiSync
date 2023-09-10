package authfacade

import (
	userModels "ResiSync/app/internal/models"
	userService "ResiSync/app/internal/services/user_service.go"
	"ResiSync/pkg/api"
	"ResiSync/pkg/models"
	postgres_db "ResiSync/shared/database"
	shared_errors "ResiSync/shared/errors"

	"go.uber.org/zap"
)

func SignUp(requestContext models.ResiSyncRequestContext, userDto *userModels.ResidentDTO) error {
	span := api.AddTrace(&requestContext, "info", "SignUp")
	defer span.End()

	log := requestContext.Log

	user, err := userService.GetNewUserObject(requestContext, userDto)
	if err != nil {
		log.Error("error while creating user object", zap.Error(err))
		return err
	}

	if err := userService.ValidateUser(requestContext, user); err != nil {
		return err
	}

	err = postgres_db.SaveOrUpdate(requestContext, user)
	if err != nil {
		log.Error("error while saving user", zap.Error(err))
		return err
	}

	return nil
}

func SignIn(requestContext models.ResiSyncRequestContext, userDto *userModels.ResidentDTO) error {
	span := api.AddTrace(&requestContext, "info", "SignIn")
	defer span.End()

	log := requestContext.Log

	isAuthenticated, err := userService.Authenticate(requestContext, userDto, userDto.Password)
	if err != nil {
		log.Error("error while authentication for user", zap.Error(err))
		return err
	}

	if !isAuthenticated {
		return shared_errors.ErrInvalidCredentials
	} else {
		userDto.Salt = ""
		userDto.Password = ""
	}

	err = userService.InitUserSession(requestContext, userDto)
	if err != nil {
		log.Error("error while initializing session for user", zap.Error(err))
		return err
	}

	go userService.UpdateLastLogIn(requestContext, userDto.Id)

	return nil
}
