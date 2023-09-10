package user_models

type Resident struct {
	Id                int64
	FirstName         string
	LastName          string
	EmailId           string
	Phone             string
	CountryCode       string
	IsActive          bool
	Password          string
	Salt              string
	LastLoginOn       int64
	ProfilePictureUrl string
	DeletedOn         int64
	UpdatedOn         int64
	CreatedOn         int64
}

func (u *Resident) GetUserDTO() *ResidentDTO {
	var user ResidentDTO
	user.Id = u.Id
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.EmailId = u.EmailId
	user.Phone = u.Phone
	user.IsActive = u.IsActive
	user.ProfilePictureUrl = u.ProfilePictureUrl
	user.LastLoginOn = u.LastLoginOn

	return &user
}

type ResidentDTO struct {
	Id                int64  `json:"id,omitempty"`
	FirstName         string `json:"firstName,omitempty" validate:"required"`
	LastName          string `json:"lastName,omitempty"`
	EmailId           string `json:"emailId,omitempty" validate:"required"`
	Phone             string `json:"contact,omitempty"`
	IsActive          bool   `json:"isActive,omitempty"`
	Password          string `json:"password,omitempty" validate:"required"`
	LastLoginOn       int64  `json:"lastLoggedIn,omitempty"`
	ProfilePictureUrl string `json:"profile_picture_url,omitempty"`

	AccessToken string `json:"accessToken,omitempty"`

	Salt string `json:"-"`
}

func (dto *ResidentDTO) IsValid() bool {
	return len(dto.EmailId) > 0 && len(dto.Password) > 0 && len(dto.FirstName) > 0
}

func (dto *ResidentDTO) GetUser() *Resident {
	var user Resident
	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.EmailId = dto.EmailId
	user.Phone = dto.Phone
	user.IsActive = dto.IsActive
	user.ProfilePictureUrl = dto.ProfilePictureUrl

	return &user
}

func (dto *ResidentDTO) TableName() string {
	return "residents"
}
