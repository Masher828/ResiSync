package user_utils

import (
	user_constants "ResiSync/app/internal/constants"
	"fmt"
)

func GetAccessTokenToUserKey(accessToken string) string {
	return fmt.Sprintf(user_constants.AccessTokenToUserKey, accessToken)
}

func GetUserToAccessTokenKey(accessToken int64) string {
	return fmt.Sprintf(user_constants.UserToAccessTokenKey, accessToken)
}
