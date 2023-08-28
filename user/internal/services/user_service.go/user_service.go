package userService

import (
	"ResiSync/pkg/api"
	"ResiSync/pkg/models"
	"ResiSync/pkg/security"
	postgres_db "ResiSync/shared/database"
	shared_utils "ResiSync/shared/utils"
	user_constants "ResiSync/user/internal/constants"
	userModels "ResiSync/user/internal/models"
	user_utils "ResiSync/user/internal/utils"
	"encoding/json"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetNewUserObject(requestContext models.ResiSyncRequestContext, userDto *userModels.ResidentDTO) (*userModels.Resident, error) {
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

func Authenticate(requestContext models.ResiSyncRequestContext, userDto *userModels.ResidentDTO, password string) (bool, error) {
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

func InitUserSession(requestContext models.ResiSyncRequestContext, userSession *userModels.ResidentDTO) error {
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

	err = redisDB.Set(requestContext.Context, key, userSessionBytes, user_constants.SessionExpiryTime).Err()
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

	var user = userModels.Resident{Id: id, LastLoginOn: shared_utils.NowInUTC().UnixNano()}

	err := postgres_db.SaveOrUpdate(requestContext, &user)
	if err != nil {
		log.Error("error while updating user last logged in time", zap.Error(err))
		return err
	}

	return nil
}
