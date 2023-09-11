package user_facade

import (
	user_constants "ResiSync/app/internal/constants/user"
	user_models "ResiSync/app/internal/models/user"
	"ResiSync/app/internal/services/user_service.go"
	"ResiSync/pkg/api"
	pkg_constants "ResiSync/pkg/constants"
	pkg_models "ResiSync/pkg/models"
	postgres_db "ResiSync/shared/database"
	shared_utils "ResiSync/shared/utils"
	"fmt"
	"mime/multipart"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

func UpdateUserProfile(requestContext pkg_models.ResiSyncRequestContext, user *user_models.ResidentDTO) error {
	span := api.AddTrace(&requestContext, "info", "UpdateUserProfile")
	defer span.End()
	log := requestContext.Log

	if err := user_service.ValidateUser(requestContext, user.GetUser()); err != nil {
		return err
	}

	err := user_service.UpdateUserProfile(requestContext, user)
	if err != nil {
		log.Error("error while updating user profile", zap.Error(err), zap.Int64("user id", user.Id))
		return err
	}

	return nil
}

func RemoveProfilePicture(requestContext pkg_models.ResiSyncRequestContext, user *user_models.Resident) error {

	span := api.AddTrace(&requestContext, "info", "RemoveProfilePicture")
	defer span.End()
	log := requestContext.Log

	err := postgres_db.GetWithFields(requestContext, &user, "profile_picture_url")
	if err != nil {
		log.Error("Error while fetching user",
			zap.Int64("user Id", user.Id), zap.Error(err))
		return err
	}

	if len(user.ProfilePictureUrl) > 0 {
		err = shared_utils.DeleteObjectFromS3(requestContext, viper.GetString(pkg_constants.AWSS3Bucket), user.ProfilePictureUrl)
		if err != nil {
			log.Error("Error while deleting user profile picture",
				zap.Int64("user Id", user.Id), zap.Error(err), zap.String("key", user.ProfilePictureUrl))
			return err
		}
	}
	return nil
}

func UpdateProfilePicture(requestContext pkg_models.ResiSyncRequestContext, file multipart.File, header *multipart.FileHeader) error {
	span := api.AddTrace(&requestContext, "info", "UpdateProfilePicture")
	defer span.End()
	log := requestContext.Log

	user := user_models.Resident{Id: requestContext.GetUserContext().ID}

	err := RemoveProfilePicture(requestContext, &user)
	if err != nil {
		log.Error("error while deleting old profile picture", zap.Error(err))
	}

	key, err := user_service.UpdateProfilePictureInS3(requestContext, file, header)
	if err != nil {
		log.Error("Error while creating s3 object for profile picture",
			zap.Int64("user Id", user.Id), zap.Error(err))
		return err
	}

	err = user_service.UpdateProfilePicture(requestContext, key)
	if err != nil {
		log.Error("Error while updating profile picture url in db",
			zap.Int64("user Id", user.Id), zap.Error(err))
		return err
	}

	return nil
}

func SendOTP(requestContext pkg_models.ResiSyncRequestContext, method, contact string) error {
	span := api.AddTrace(&requestContext, "info", "SendOTP")
	defer span.End()
	log := requestContext.Log
	var err error = nil
	var key string
	otp := shared_utils.GenerateOTP()

	if method == "email" {
		err = user_service.SendEmailOtp(requestContext, contact, otp)
		key = fmt.Sprintf(user_constants.EmailOtpKey, contact)
	} else {
		return nil
	}

	if err != nil {
		log.Error("Error while sending otp", zap.String("method", method), zap.String("contact", contact))
		return err
	}

	redisDb := api.ApplicationContext.Redis

	redisDb.Set(requestContext.Context, key, otp, user_constants.OTPExpiry).Err()

	return nil
}
