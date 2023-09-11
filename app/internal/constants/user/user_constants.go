package user_constants

import "time"

var GetUserProfileFields = []string{"id", "first_name", "last_name", "email_id", "profile_picture_url",
	"phone", "last_login_on", "updated_on", "created_on"}

var UpdateUserProfileFields = []string{"first_name", "last_name", "phone"}

const ProfilePictureS3Folder = "/profile_pictures/"

var EmailOtpKey = "EmailOTP:%s"

var OTPExpiry = time.Minute * 5
