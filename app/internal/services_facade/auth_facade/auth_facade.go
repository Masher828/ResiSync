package auth_facade

import (
	"ResiSync/app/internal/models"
	user_models "ResiSync/app/internal/models/user"
	"ResiSync/app/internal/services/auth_service"
	"ResiSync/app/internal/services/user_service.go"
	"ResiSync/pkg/api"
	pkg_models "ResiSync/pkg/models"
	postgres_db "ResiSync/shared/database"
	shared_errors "ResiSync/shared/errors"

	"go.uber.org/zap"
)

func SignUp(requestContext pkg_models.ResiSyncRequestContext, userDto *user_models.ResidentDTO) error {
	span := api.AddTrace(&requestContext, "info", "SignUp")
	defer span.End()

	log := requestContext.Log

	user, err := user_service.GetNewUserObject(requestContext, userDto)
	if err != nil {
		log.Error("error while creating user object", zap.Error(err))
		return err
	}

	if err := user_service.ValidateUser(requestContext, user); err != nil {
		return err
	}

	err = postgres_db.SaveOrUpdate(requestContext, user)
	if err != nil {
		log.Error("error while saving user", zap.Error(err))
		return err
	}

	return nil
}

func SignIn(requestContext pkg_models.ResiSyncRequestContext, userDto *user_models.ResidentDTO) error {
	span := api.AddTrace(&requestContext, "info", "SignIn")
	defer span.End()

	log := requestContext.Log

	isAuthenticated, err := auth_service.Authenticate(requestContext, userDto, userDto.Password)
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

	err = auth_service.InitUserSession(requestContext, userDto)
	if err != nil {
		log.Error("error while initializing session for user", zap.Error(err))
		return err
	}

	go auth_service.UpdateLastLogIn(requestContext, userDto.Id)

	return nil
}

func VerifyOtp(requestContext pkg_models.ResiSyncRequestContext, verifyRequest *models.ResetPassword) bool {
	span := api.AddTrace(&requestContext, "info", "VerifyOtp")
	defer span.End()

	if verifyRequest.Method == "email" {
		return user_service.VerifyEmailOtp(requestContext, verifyRequest.Contact, verifyRequest.Otp)
	}

	return false
}

func ResetPassword(requestContext pkg_models.ResiSyncRequestContext, verifyRequest *models.ResetPassword) error {
	span := api.AddTrace(&requestContext, "info", "ResetPassword")
	defer span.End()

	log := requestContext.Log
	if !VerifyOtp(requestContext, verifyRequest) {
		return shared_errors.ErrInvalidOTP
	}

	var condition string

	if verifyRequest.Method == "email" {
		condition = "email_id = ?"
	} else {
		return nil
	}

	var user user_models.Resident
	err := postgres_db.GetDataWithCriteria(requestContext, &user, condition)
	if err != nil {
		log.Error("error while getting user details",
			zap.String("method", verifyRequest.Method), zap.String("contact", verifyRequest.Contact), zap.Error(err))
		return err
	}

	err = auth_service.UpdatePassword(requestContext, &user, verifyRequest.Password)
	if err != nil {
		log.Error("error while updating password",
			zap.Int64("user id", user.Id), zap.Error(err))
		return err
	}
	return nil

}
