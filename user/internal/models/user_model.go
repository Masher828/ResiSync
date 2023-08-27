package userModels

type User struct {
	Id        int64
	Name      string `gorm:"name"`
	EmailId   string `gorm:"unique"`
	CreatedOn int64
}
