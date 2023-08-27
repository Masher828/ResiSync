package userService

type User struct {
	Id           int64 `gorm:"primaryKey"`
	FirstName    string
	LastName     string
	EmailId      string `gorm:"unique"`
	Phone        string
	IsActive     bool
	Password     string
	Salt         string
	LastLoggedIn int64
	DeletedAt    int64
	UpdatedAt    int64 `gorm:"autoUpdateTime:nano"`
	CreatedAt    int64 `gorm:"autoCreateTime:nano"`
}

func (user *User) SignIn() {

}
