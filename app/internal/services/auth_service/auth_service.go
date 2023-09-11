package auth_service

import (
	user_models "ResiSync/app/internal/models/user"
	"ResiSync/app/internal/services/user_service.go"
	user_utils "ResiSync/app/internal/utils"
	"ResiSync/pkg/api"
	pkg_constants "ResiSync/pkg/constants"
	pkg_models "ResiSync/pkg/models"
	"ResiSync/pkg/security"
	postgres_db "ResiSync/shared/database"
	shared_utils "ResiSync/shared/utils"
	"encoding/json"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Authenticate(requestContext pkg_models.ResiSyncRequestContext, userDto *user_models.ResidentDTO, password string) (bool, error) {
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

func InitUserSession(requestContext pkg_models.ResiSyncRequestContext, userSession *user_models.ResidentDTO) error {
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

func UpdateLastLogIn(requestContext pkg_models.ResiSyncRequestContext, id int64) error {
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

func LogOut(requestContext pkg_models.ResiSyncRequestContext) {
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

func UpdatePassword(requestContext pkg_models.ResiSyncRequestContext, user *user_models.Resident, newPassword string) error {

	span := api.AddTrace(&requestContext, "info", "UpdatePassword")
	defer span.End()
	log := requestContext.Log

	err := user_service.HashUserPassword(requestContext, user, newPassword)
	if err != nil {
		log.Error("error while hashing the password", zap.Int64("user id", user.Id), zap.Error(err))
		return err
	}

	err = postgres_db.UpdateWithFields(requestContext, &user, "password", "salt")
	if err != nil {
		log.Error("error while updating the password in postgres", zap.Int64("user id", user.Id), zap.Error(err))
		return err
	}

	return nil

}
