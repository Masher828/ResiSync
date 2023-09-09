package user_constants

import "time"

const (
	AccessTokenToUserKey = "acessToken:%s"
	UserToAccessTokenKey = "user:accessToken:%d"

	SessionExpiryTime = time.Hour * 24
)
