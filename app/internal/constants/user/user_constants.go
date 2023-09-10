package user_constants

var GetUserProfileFields = []string{"id", "first_name", "last_name", "email_id", "profile_picture_url",
	"phone", "last_login_on", "updated_on", "created_on"}

var UpdateUserProfileFields = []string{"first_name", "last_name", "phone"}

var ProfilePictureS3Folder = "/profile_pictures/"
