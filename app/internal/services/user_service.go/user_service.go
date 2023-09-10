package user_service

import (
	user_constants "ResiSync/app/internal/constants/user"
	user_models "ResiSync/app/internal/models"
	user_utils "ResiSync/app/internal/utils"
	"ResiSync/pkg/api"
	pkg_constants "ResiSync/pkg/constants"
	"ResiSync/pkg/models"
	"ResiSync/pkg/security"
	postgres_db "ResiSync/shared/database"
	shared_errors "ResiSync/shared/errors"
	shared_utils "ResiSync/shared/utils"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"strconv"
	"time"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetNewUserObject(requestContext models.ResiSyncRequestContext, userDto *user_models.ResidentDTO) (*user_models.Resident, error) {
	span := api.AddTrace(&requestContext, "info", "GetNewUser")
	defer span.End()

	log := requestContext.Log

	user := userDto.GetUser()

	now := shared_utils.NowInUTC().UnixNano()

	user.CreatedOn = now
	user.DeletedOn = now
	user.IsActive = true

	var err error

	user.Password, user.Salt, err = security.Hashpassword(requestContext, 16, userDto.Password)
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

func Authenticate(requestContext models.ResiSyncRequestContext, userDto *user_models.ResidentDTO, password string) (bool, error) {
	span := api.AddTrace(&requestContext, "info", "Authenticate")
	defer span.End()

	log := requestContext.Log

	userDto.Password = ""

	err := postgres_db.GetData(requestContext, &userDto)
	if err != nil {
		log.Error("error while querying data for user", zap.Error(err))
		return false, err
	}

	passwordMatching := security.ComparePassword(userDto.Password, userDto.Salt, password)

	return passwordMatching, nil
}

func InitUserSession(requestContext models.ResiSyncRequestContext, userSession *user_models.ResidentDTO) error {
	span := api.AddTrace(&requestContext, "info", "InitUserSession")
	defer span.End()

	log := requestContext.Log

	userSession.AccessToken = uuid.New().String()

	redisDB := api.ApplicationContext.Redis

	userSessionBytes, err := json.Marshal(userSession)
	if err != nil {
		log.Error("error while marshalling user sesison", zap.Error(err))
		return err
	}

	key := user_utils.GetAccessTokenToUserKey(userSession.AccessToken)

	err = redisDB.Set(requestContext.Context, key, userSessionBytes, pkg_constants.SessionExpiryTime).Err()
	if err != nil {
		log.Error("error while creating access token to user key", zap.Error(err))
		return err
	}

	key = user_utils.GetUserToAccessTokenKey(userSession.Id)
	err = redisDB.LPush(requestContext.Context, key, 0, userSession.AccessToken).Err()
	if err != nil {
		log.Error("error while creating user to access token key", zap.Error(err))
		return err
	}

	return nil
}

func UpdateLastLogIn(requestContext models.ResiSyncRequestContext, id int64) error {
	span := api.AddTrace(&requestContext, "info", "UpdateLastLogIn")
	defer span.End()

	log := requestContext.Log

	var user = user_models.Resident{Id: id, LastLoginOn: shared_utils.NowInUTC().UnixNano()}

	err := postgres_db.UpdateWithFields(requestContext, &user, "last_login_on")
	if err != nil {
		log.Error("error while updating user last logged in time", zap.Error(err))
		return err
	}

	return nil
}

func LogOut(requestContext models.ResiSyncRequestContext) {
	span := api.AddTrace(&requestContext, "info", "LogOut")
	defer span.End()
	log := requestContext.Log

	userContext := requestContext.GetUserContext()

	redisDB := api.ApplicationContext.Redis

	key := user_utils.GetAccessTokenToUserKey(userContext.AccessToken)

	if err := redisDB.Del(requestContext.Context, key).Err(); err != nil {
		log.Error("error while removing user access token", zap.Error(err), zap.Int64("userId", userContext.ID))
	}

	key = user_utils.GetUserToAccessTokenKey(userContext.ID)
	if err := redisDB.LRem(requestContext.Context, key, 1, userContext.AccessToken).Err(); err != nil {
		log.Error("error while removing user access token ", zap.Error(err), zap.Int64("userId", userContext.ID))
	}
}

func GetUserProfile(requestContext models.ResiSyncRequestContext) (*user_models.ResidentDTO, error) {
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

func UpdateUserProfile(requestContext models.ResiSyncRequestContext, user *user_models.ResidentDTO) error {
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

func ValidateUser(requestContext models.ResiSyncRequestContext, user *user_models.Resident) error {
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

func UpdateProfilePictureInS3(requestContext models.ResiSyncRequestContext, file multipart.File, header *multipart.FileHeader) (string, error) {
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

func UpdateProfilePicture(requestContext models.ResiSyncRequestContext, profilePictureUrl string) error {
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
