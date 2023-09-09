package user_utils

import (
	"ResiSync/app/internal/constants"
	"fmt"
)

func GetAccessTokenToUserKey(accessToken string) string {
	return fmt.Sprintf(constants.AccessTokenToUserKey, accessToken)
}

func GetUserToAccessTokenKey(accessToken int64) string {
	return fmt.Sprintf(constants.UserToAccessTokenKey, accessToken)
}
