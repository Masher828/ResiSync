package user_utils

import (
	pkg_constants "ResiSync/pkg/constants"
	"fmt"
)

func GetAccessTokenToUserKey(accessToken string) string {
	return fmt.Sprintf(pkg_constants.AccessTokenToUserFormatKey, accessToken)
}

func GetUserToAccessTokenKey(accessToken int64) string {
	return fmt.Sprintf(pkg_constants.UserToAccessTokenKey, accessToken)
}
