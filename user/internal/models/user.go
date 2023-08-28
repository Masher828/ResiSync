package userModels

type Resident struct {
	Id          int64
	FirstName   string
	LastName    string
	EmailId     string
	Phone       string
	IsActive    bool
	Password    string
	Salt        string
	LastLoginOn int64
	DeletedOn   int64
	UpdatedOn   int64
	CreatedOn   int64
}

func (u *Resident) GetUserDTO() *ResidentDTO {
	var user ResidentDTO
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.EmailId = u.EmailId
	user.Phone = u.Phone
	user.IsActive = u.IsActive
	user.LastLoginOn = u.LastLoginOn

	return &user
}

type ResidentDTO struct {
	Id          int64  `json:"id,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	EmailId     string `json:"emailId,omitempty"`
	Phone       string `json:"phone,omitempty"`
	IsActive    bool   `json:"isActive,omitempty"`
	Password    string `json:"password,omitempty"`
	LastLoginOn int64  `json:"lastLoggedIn,omitempty"`

	AccessToken string `json:"accessToken,omitempty"`

	Salt string `json:"-"`
}

func (dto *ResidentDTO) GetUser() *Resident {
	var user Resident
	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.EmailId = dto.EmailId
	user.Phone = dto.Phone
	user.IsActive = dto.IsActive

	return &user
}

func (dto *ResidentDTO) TableName() string {
	return "residents"
}
