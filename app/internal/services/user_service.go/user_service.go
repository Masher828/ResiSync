package user_service

import (
	user_constants "ResiSync/app/internal/constants/user"
	user_models "ResiSync/app/internal/models/user"
	"fmt"

	"ResiSync/pkg/api"
	pkg_constants "ResiSync/pkg/constants"
	pkg_models "ResiSync/pkg/models"

	"ResiSync/pkg/security"
	postgres_db "ResiSync/shared/database"
	shared_errors "ResiSync/shared/errors"
	shared_utils "ResiSync/shared/utils"
	"bytes"
	"mime/multipart"
	"strconv"
	"time"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetNewUserObject(requestContext pkg_models.ResiSyncRequestContext, userDto *user_models.ResidentDTO) (*user_models.Resident, error) {
	span := api.AddTrace(&requestContext, "info", "GetNewUser")
	defer span.End()

	log := requestContext.Log

	user := userDto.GetUser()

	now := shared_utils.NowInUTC().UnixNano()

	user.CreatedOn = now
	user.DeletedOn = now
	user.IsActive = true

	var err error

	err = HashUserPassword(requestContext, user, userDto.Password)
	if err != nil {
		log.Error("Error while hashing password", zap.Error(err))
		return nil, err
	}

	user.Id, err = postgres_db.GetSequenceId(requestContext)
	if err != nil {
		log.Error("Error while getting UserId", zap.Error(err))
		return nil, err
	}

	return user, err
}

func GetUserProfile(requestContext pkg_models.ResiSyncRequestContext) (*user_models.ResidentDTO, error) {
	span := api.AddTrace(&requestContext, "info", "GetUserProfile")
	defer span.End()
	log := requestContext.Log

	userContext := requestContext.GetUserContext()

	user := user_models.Resident{Id: userContext.ID}

	err := postgres_db.GetWithFields(requestContext, &user, user_constants.GetUserProfileFields...)
	if err != nil {
		log.Error("Error while fetching profile", zap.Int64("user id", userContext.ID), zap.Error(err))
		return nil, err
	}

	user.ProfilePictureUrl = shared_utils.GetPresignedS3Url(requestContext, viper.GetString(pkg_constants.AWSS3Bucket), user.ProfilePictureUrl, time.Minute*5)
	return user.GetUserDTO(), nil
}

func UpdateUserProfile(requestContext pkg_models.ResiSyncRequestContext, user *user_models.ResidentDTO) error {
	span := api.AddTrace(&requestContext, "info", "UpdateUserProfile")
	defer span.End()
	log := requestContext.Log

	err := postgres_db.UpdateWithFields(requestContext, &user, user_constants.UpdateUserProfileFields...)
	if err != nil {
		log.Error("Error while updating profile", zap.Int64("user id", requestContext.GetUserContext().ID), zap.Error(err))
		return err
	}
	return nil
}

func ValidateUser(requestContext pkg_models.ResiSyncRequestContext, user *user_models.Resident) error {
	span := api.AddTrace(&requestContext, "info", "ValidateUser")
	defer span.End()

	if !shared_utils.IsValidEmail(user.EmailId) {
		return shared_errors.ErrInvalidEmail
	}

	if !shared_utils.IsValidContact(user.Phone, user.CountryCode) {
		return shared_errors.ErrInvalidContact
	}

	return nil
}

func UpdateProfilePictureInS3(requestContext pkg_models.ResiSyncRequestContext, file multipart.File, header *multipart.FileHeader) (string, error) {
	span := api.AddTrace(&requestContext, "info", "GetProfilePictureS3Object")
	defer span.End()

	log := requestContext.Log

	userId := requestContext.GetUserContext().ID
	if len(header.Header["Content-Type"]) == 0 {
		log.Error("content type not present in image", zap.Int64("userId", userId))
		return "", shared_errors.ErrInvalidPayload
	}

	buffer := make([]byte, header.Size)

	_, err := file.Read(buffer)
	if err != nil {
		log.Error("Error while reading file", zap.Int64("userId", userId), zap.Error(err))
		return "", err
	}

	fileBytes := bytes.NewReader(buffer)

	fileType := header.Header["Content-Type"][0]

	pathToProfile := user_constants.ProfilePictureS3Folder + strings.ReplaceAll(header.Filename, " ", "_") + "-" + strconv.FormatInt(userId, 10)

	params := s3.PutObjectInput{
		Bucket:        aws.String(viper.GetString(pkg_constants.AWSS3Bucket)),
		Key:           aws.String(pathToProfile),
		Body:          fileBytes,
		ContentLength: aws.Int64(header.Size),
		ContentType:   aws.String(fileType),
	}

	s3Session := api.ApplicationContext.S3Session
	_, err = s3Session.PutObject(&params)
	if err != nil {
		log.Error("Error while uploading profile picture to s3",
			zap.Int64("user Id", userId), zap.Error(err))
		return "", err
	}
	return "", nil
}

func UpdateProfilePicture(requestContext pkg_models.ResiSyncRequestContext, profilePictureUrl string) error {
	span := api.AddTrace(&requestContext, "info", "UpdateProfilePicture")
	defer span.End()

	log := requestContext.Log

	user := user_models.Resident{
		Id:                requestContext.GetUserContext().ID,
		ProfilePictureUrl: profilePictureUrl,
	}

	err := postgres_db.UpdateWithFields(requestContext, &user, "profile_picture_url")
	if err != nil {
		log.Error("Error while updating profile pic url",
			zap.Int64("user id", user.Id), zap.Error(err))
		return err
	}

	return nil
}

func SendEmailOtp(requestContext pkg_models.ResiSyncRequestContext, emailId, otp string) error {

	span := api.AddTrace(&requestContext, "info", "SendEmailOtp")
	defer span.End()

	log := requestContext.Log
	if len(emailId) == 0 {
		log.Error("error no email is provided")
		return nil
	}

	redisDb := api.ApplicationContext.Redis

	key := fmt.Sprintf(user_constants.EmailOtpKey, emailId)
	count, _ := redisDb.Exists(requestContext.Context, key).Result()
	if count > 0 {
		return nil
	}

	// go pkg_utils.SendEmail(emailId, "Reset Password OTP", "your otp is "+otp)
	return nil
}

func VerifyEmailOtp(requestContext pkg_models.ResiSyncRequestContext, emailId, enteredOtp string) bool {
	span := api.AddTrace(&requestContext, "info", "VerifyEmailOtp")
	defer span.End()

	log := requestContext.Log

	redisDb := api.ApplicationContext.Redis

	key := fmt.Sprintf(user_constants.EmailOtpKey, emailId)
	actualOtp, err := redisDb.Get(requestContext.Context, key).Result()
	if err != nil {
		log.Error("error while getting otp", zap.Error(err), zap.String("emailId", emailId))
	}

	return actualOtp == enteredOtp
}

func HashUserPassword(requestContext pkg_models.ResiSyncRequestContext, user *user_models.Resident, password string) error {
	span := api.AddTrace(&requestContext, "info", "HashUserPassword")
	defer span.End()

	log := requestContext.Log
	var err error = nil
	user.Password, user.Salt, err = security.Hashpassword(requestContext, 16, password)
	if err != nil {
		log.Error("Error while hashing password", zap.Error(err))
		return err
	}
	return nil
}
