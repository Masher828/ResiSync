package userModels

type UserLoginContext struct {
	Id        int64
	Name      string `gorm:"name"`
	EmailId   string `gorm:"unique"`
	CreatedOn int64
}

// func (u *UserLoginContext) GetUserLoginContext(user interface{}) (*UserLoginContext) {
// 	reflect
// }
